package hfvault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
)

// createHash hashes a key to sha256 and then string-hexes it
// Used to extend passwort length
func createHash(key string) string {
	hasher := sha256.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

// Encrypt takes data and encrypts it with a passphrase
func Encrypt(data []byte, passphrase string) ([]byte, error) {
	key := []byte(createHash(passphrase)[0:32])
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("new cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create gcm: %w", err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("read random nonce: %w", err)
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// Decrypt inverts Encrypt
func Decrypt(data []byte, passphrase string) ([]byte, error) {
	key := []byte(createHash(passphrase)[0:32])
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("new cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("new gcm: %w", err)
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("open gcm: %w", err)
	}
	return plaintext, nil
}

// EncryptToBase64 creates an AES-128 encrypted and then base64 encoded ciphertext
func EncryptToBase64(plaintext []byte, passphrase string) (string, error) {
	ciphertext, err := Encrypt(plaintext, passphrase)
	if err != nil {
		return "", fmt.Errorf("encrypt: %w", err)
	}
	enc := base64.StdEncoding.EncodeToString(ciphertext)
	return enc, nil
}

// DecryptFromBase64 takes a base64 encoded, AES-128 encrypted string and decodes it
func DecryptFromBase64(cipherstring, passphrase string) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(cipherstring)
	if err != nil {
		return nil, fmt.Errorf("base64dec failed: %w", err)
	}
	plaintext, err := Decrypt(ciphertext, passphrase)
	if err != nil {
		return nil, fmt.Errorf("base64dec failed: %w", err)
	}
	return plaintext, nil
}
