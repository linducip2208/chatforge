package main

import (
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"
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

func loginUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	email := r.FormValue("email")
	password := r.FormValue("password")
	user, err := db.GetUserByEmail(email)
	if err != nil || !checkPassword(password, user.Password) {
		http.Redirect(w, r, "/login?msg=Invalid+credentials", http.StatusSeeOther)
		return
	}
	token := randToken()
	saveSession(token, user.ID)
	http.SetCookie(w, &http.Cookie{Name: "chatgo_sess", Value: token, Path: "/", MaxAge: 86400 * 30})
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
	http.SetCookie(w, &http.Cookie{Name: "chatgo_sess", Value: token, Path: "/", MaxAge: 86400 * 30})
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
			u, err := db.GetUserByID(uid)
			if err != nil || u.Role != "admin" {
				http.Error(w, "Forbidden", 403)
				return
			}
		}
		next(w, r)
	}
}

func checkPassword(plain, stored string) bool {
	hash, _ := hashPassword(plain)
	return hash == stored
}
func hashPassword(p string) (string, error) {
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
		u, err := db.GetUserByID(uid)
		if err != nil || u.Role != "admin" {
			http.Error(w, "Forbidden", 403)
			return
		}
		next(w, r)
	})
}

