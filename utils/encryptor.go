package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	// 优先加载 .env.local，再加载 .env
	err := godotenv.Load(".env.local", ".env")
	if err != nil {
		log.Println("Warning: .env.local or .env file not found")
	}
}

// deriveKey 用主密码生成加密Key（内部私有）
func deriveKey() []byte {
	masterPassword := os.Getenv("MASTER_PASSWORD")
	if masterPassword == "" {
		panic("环境变量 MASTER_PASSWORD 未设置")
	}
	hash := sha256.Sum256([]byte(masterPassword))
	return hash[:]
}

// Encrypt 加密后返回 Base64 字符串
func Encrypt(plaintext []byte) (string, error) {
	key := deriveKey()

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesgcm.Seal(nonce, nonce, plaintext, nil)

	// 转成 base64 字符串，存储在数据库中
	encoded := base64.StdEncoding.EncodeToString(ciphertext)
	return encoded, nil
}

// Decrypt 解密 Base64 字符串
func Decrypt(ciphertextBase64 string) (string, error) {
	key := deriveKey()

	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesgcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("密文数据太短，无法解密")
	}

	nonce := ciphertext[:nonceSize]
	ciphertext = ciphertext[nonceSize:]

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
