package utils

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"time"
)

// Steam TOTP 可用的字符集
var steamChars = []byte("23456789BCDFGHJKMNPQRTVWXY")

// GenerateCode 根据 sharedSecret 和当前时间生成 Steam Guard 验证码
func GenerateCode(sharedSecret string) (string, error) {
	return GenerateCodeWithTime(sharedSecret, time.Now())
}

// GenerateCodeWithTime 根据 sharedSecret 和指定时间生成 Steam Guard 验证码（方便测试）
func GenerateCodeWithTime(sharedSecret string, t time.Time) (string, error) {
	if sharedSecret == "" {
		return "", errors.New("sharedSecret 不能为空")
	}

	// Base64 decode
	secretBytes, err := base64.StdEncoding.DecodeString(sharedSecret)
	if err != nil {
		return "", err
	}

	// 时间戳，单位30秒
	timeSlice := uint64(t.Unix() / 30)

	// 把时间戳转成8字节 big-endian
	timeBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(timeBytes, timeSlice)

	// 生成 HMAC-SHA1
	mac := hmac.New(sha1.New, secretBytes)
	mac.Write(timeBytes)
	hash := mac.Sum(nil)

	// 动态偏移量，取最后一字节的低4位
	offset := hash[len(hash)-1] & 0x0F

	// 取4字节并转为无符号整数
	code := binary.BigEndian.Uint32(hash[offset : offset+4])
	code &= 0x7FFFFFFF // 保证是正数

	// 生成5位 Steam 验证码
	var result [5]byte
	for i := 0; i < 5; i++ {
		result[i] = steamChars[code%uint32(len(steamChars))]
		code /= uint32(len(steamChars))
	}

	return string(result[:]), nil
}
