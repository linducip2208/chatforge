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
			keyStr = "chatgo-32bytekey-xxxxxxxxxxxxxx!!!"
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
