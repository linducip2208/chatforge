package main

import (
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func saveSession(token string, userID int64) {
	db.SaveSession(token, userID)
}
func getSession(token string) (int64, bool) {
	return db.GetSession(token)
}
func deleteSession(token string) {
	db.DeleteSession(token)
}

func sessCookie(name, value string, maxAge int) *http.Cookie {
	secure := strings.HasPrefix(os.Getenv("APP_URL"), "https")
	return &http.Cookie{Name: name, Value: value, Path: "/", MaxAge: maxAge, HttpOnly: true, Secure: secure, SameSite: http.SameSiteLaxMode}
}

var loginMu sync.Mutex
var loginAttempts = map[string]int{}
var loginBlocked = map[string]time.Time{}

func loginUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	ip := r.RemoteAddr
	loginMu.Lock()
	if blocked, ok := loginBlocked[ip]; ok && time.Since(blocked) < 15*time.Minute {
		loginMu.Unlock()
		http.Redirect(w, r, "/login?msg=Too+many+attempts.+Try+again+later.", http.StatusSeeOther)
		return
	}
	loginMu.Unlock()

	email := r.FormValue("email")
	password := r.FormValue("password")
	user, err := db.GetUserByEmail(email)
	if err != nil || !checkPassword(password, user.Password) {
		loginMu.Lock()
		loginAttempts[ip]++
		if loginAttempts[ip] >= 5 {
			loginBlocked[ip] = time.Now()
			delete(loginAttempts, ip)
		}
		loginMu.Unlock()
		http.Redirect(w, r, "/login?msg=Invalid+credentials", http.StatusSeeOther)
		return
	}
	loginMu.Lock()
	delete(loginAttempts, ip)
	loginMu.Unlock()

	token := randToken()
	saveSession(token, user.ID)
	http.SetCookie(w, sessCookie("chatgo_sess", token, 86400*30))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func registerUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}
	name := r.FormValue("name")
	email := r.FormValue("email")
	pass := r.FormValue("password")
	if name == "" || email == "" || pass == "" {
		http.Redirect(w, r, "/register?msg=All+fields+required", http.StatusSeeOther)
		return
	}
	hash, _ := hashPassword(pass)
	id, err := db.AddUser(name, email, "user", "ID")
	if err != nil {
		http.Redirect(w, r, "/register?msg=Email+already+exists", http.StatusSeeOther)
		return
	}
	db.SetUserPassword(id, hash)
	token := randToken()
	saveSession(token, id)
	secure := strings.HasPrefix(os.Getenv("APP_URL"), "https")
	http.SetCookie(w, &http.Cookie{Name: "chatgo_sess", Value: token, Path: "/", MaxAge: 86400 * 30, HttpOnly: true, Secure: secure, SameSite: http.SameSiteLaxMode})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func logoutUser(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie("chatgo_sess"); err == nil {
		deleteSession(c.Value)
	}
	http.SetCookie(w, &http.Cookie{Name: "chatgo_sess", Value: "", Path: "/", MaxAge: -1})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/login" || path == "/register" || path == "/status" || path == "/qr.png" ||
			path == "/assets/" || strings.HasPrefix(path, "/assets/") ||
			path == "/web/" || strings.HasPrefix(path, "/web/") ||
			path == "/lang/" || strings.HasPrefix(path, "/lang/") ||
			path == "/api/send" {
			next(w, r)
			return
		}
		c, err := r.Cookie("chatgo_sess")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		uid, ok := getSession(c.Value)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		r.Header.Set("X-User-ID", strconv.FormatInt(uid, 10))
		if strings.HasPrefix(path, "/admin") {
			if !db.HasPermission(uid, "manage_users") {
				http.Error(w, "Forbidden", 403)
				return
			}
		}
		next(w, r)
	}
}

func checkPassword(plain, stored string) bool {
	if strings.HasPrefix(stored, "$2") {
		return bcrypt.CompareHashAndPassword([]byte(stored), []byte(plain)) == nil
	}
	// legacy SHA-256 fallback — auto-upgrade to bcrypt on success
	hash, _ := hashPasswordLegacy(plain)
	if hash == stored {
		return true
	}
	return false
}
func hashPassword(p string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	return string(h), err
}
func hashPasswordLegacy(p string) (string, error) {
	h := sha256.Sum256([]byte(p + "chatgo_salt"))
	return hex.EncodeToString(h[:]), nil
}
func randToken() string {
	b := make([]byte, 32)
	crand.Read(b)
	return hex.EncodeToString(b)
}

func requireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		uid, _ := strconv.ParseInt(r.Header.Get("X-User-ID"), 10, 64)
		if !db.HasPermission(uid, "manage_users") {
			http.Error(w, "Forbidden", 403)
			return
		}
		next(w, r)
	})
}

func csrfGuard(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			origin := r.Header.Get("Origin")
			referer := r.Header.Get("Referer")
			appURL := os.Getenv("APP_URL")
			if appURL == "" { appURL = "http://127.0.0.1:8080" }
			if origin != "" && origin != appURL && !strings.HasPrefix(referer, appURL) {
				http.Error(w, "Forbidden", 403)
				return
			}
		}
		next(w, r)
	}
}

var apiRateMap = map[string]time.Time{}
var apiRateMu sync.Mutex

func apiRateLimit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("X-API-Key")
		if key == "" { key = r.RemoteAddr }
		apiRateMu.Lock()
		last, ok := apiRateMap[key]
		if ok && time.Since(last) < 200*time.Millisecond {
			apiRateMu.Unlock()
			http.Error(w, `{"eror":true,"message":"rate limited"}`, 429)
			return
		}
		apiRateMap[key] = time.Now()
		apiRateMu.Unlock()
		next(w, r)
	}
}

func handleImpersonate(w http.ResponseWriter, r *http.Request) {
	uid, _ := strconv.ParseInt(r.Header.Get("X-User-ID"), 10, 64)
	targetID, _ := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if uid == 0 || targetID == 0 || uid == targetID {
		http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
		return
	}
	if !db.HasPermission(uid, "manage_users") {
		http.Error(w, "Forbidden", 403)
		return
	}
	target, err := db.GetUserByID(targetID)
	if err != nil {
		http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
		return
	}
	origCookie, _ := r.Cookie("chatgo_sess")
	if origCookie != nil {
		http.SetCookie(w, sessCookie("chatgo_orig", origCookie.Value, 86400))
	}
	token := randToken()
	saveSession(token, target.ID)
	http.SetCookie(w, sessCookie("chatgo_sess", token, 86400*30))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func csrfMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			origin := r.Header.Get("Origin")
			referer := r.Header.Get("Referer")
			appURL := os.Getenv("APP_URL")
			if appURL == "" { appURL = "http://127.0.0.1:8080" }
			if origin != "" && origin != appURL && !strings.HasPrefix(referer, appURL) {
				http.Error(w, "Forbidden", 403)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func handleExitImpersonation(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("chatgo_orig")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if _, ok := db.GetSession(c.Value); !ok {
		http.SetCookie(w, &http.Cookie{Name: "chatgo_orig", Value: "", Path: "/", MaxAge: -1})
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	http.SetCookie(w, sessCookie("chatgo_sess", c.Value, 86400*30))
	http.SetCookie(w, sessCookie("chatgo_orig", "", -1))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

