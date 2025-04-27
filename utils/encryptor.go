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
	"path/filepath"

	"github.com/joho/godotenv"
)

// 这个变量将在编译时通过ldflags注入
var embeddedMasterPassword string

func init() {
	// 尝试多种方式加载环境变量:
	// 1. 先尝试从当前目录加载
	err1 := godotenv.Load(".env.local", ".env")

	// 2. 再尝试从可执行文件所在目录加载
	execPath, err := os.Executable()
	if err == nil {
		execDir := filepath.Dir(execPath)
		err2 := godotenv.Load(filepath.Join(execDir, ".env.local"), filepath.Join(execDir, ".env"))
		if err1 != nil && err2 != nil {
			log.Println("Warning: Failed to load .env files from current and executable directories")
		}
	} else {
		log.Println("Warning: Failed to determine executable path:", err)
		if err1 != nil {
			log.Println("Warning: Also failed to load .env files from current directory")
		}
	}
}

// deriveKey 用主密码生成加密Key（内部私有）
func deriveKey() []byte {
	// 优先使用环境变量
	masterPassword := os.Getenv("MASTER_PASSWORD")

	// 如果环境变量为空，使用编译时嵌入的密码
	if masterPassword == "" && embeddedMasterPassword != "" {
		masterPassword = embeddedMasterPassword
	}

	// 如果仍然为空，则报错
	if masterPassword == "" {
		panic("未能获取主密码，环境变量MASTER_PASSWORD未设置且编译时未嵌入密码")
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
