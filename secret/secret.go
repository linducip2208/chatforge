package secret

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"os"
	"sync"
)

var (
	key  []byte
	once sync.Once
)

func ensureKey() {
	once.Do(func() {
		keyStr := os.Getenv("CHATGO_ENC_KEY")
		if keyStr == "" {
			keyStr = readOrCreateKeyFile()
		}
		k := []byte(keyStr)
		if len(k) < 32 {
			padded := make([]byte, 32)
			copy(padded, k)
			k = padded
		}
		key = k[:32]
	})
}

// Encrypt plaintext → base64-encoded ciphertext.
func Encrypt(plain string) (string, error) {
	ensureKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plain), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt base64 ciphertext → plaintext.
func Decrypt(encoded string) (string, error) {
	ensureKey()
	ciphertext, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", err
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plain, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}

func readOrCreateKeyFile() string {
	keyFile := "storage/enc.key"
	if data, err := os.ReadFile(keyFile); err == nil && len(data) >= 32 {
		return string(data[:32])
	}
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		panic("secret: failed to generate encryption key: " + err.Error())
	}
	os.MkdirAll("storage", 0700)
	if err := os.WriteFile(keyFile, key, 0600); err != nil {
		panic("secret: failed to write encryption key: " + err.Error())
	}
	return string(key)
}
