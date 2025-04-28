package utils

import (
	"encoding/base32"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"

	pb "auth/utils/otp_extractor/proto"
	"google.golang.org/protobuf/proto"
)

// 创建 OTP 导出的结构体
type OtpEntry struct {
	Name    string `json:"name"`
	Secret  string `json:"secret"`
	Issuer  string `json:"issuer"`
	Type    string `json:"type"`
	Counter int64  `json:"counter,omitempty"`
	URL     string `json:"url"`
}

// 从URL提取OTP信息
func ExtractOtpFromUrl(otpURL string) ([]OtpEntry, error) {
	// 解析URL
	parsedURL, err := url.Parse(otpURL)
	if err != nil {
		return nil, fmt.Errorf("URL解析失败: %v", err)
	}

	// 获取data参数
	queryParams, err := url.ParseQuery(parsedURL.RawQuery)
	if err != nil {
		return nil, fmt.Errorf("解析查询参数失败: %v", err)
	}

	dataBase64 := queryParams.Get("data")
	if dataBase64 == "" {
		return nil, fmt.Errorf("URL中没有data参数")
	}

	// 替换空格为+，修正可能的base64编码问题
	dataBase64 = strings.ReplaceAll(dataBase64, " ", "+")

	// base64解码
	rawData, err := base64.StdEncoding.DecodeString(dataBase64)
	if err != nil {
		return nil, fmt.Errorf("Base64解码失败: %v", err)
	}

	// 解析protobuf - 使用生成的代码
	payload := &pb.MigrationPayload{}
	if err := proto.Unmarshal(rawData, payload); err != nil {
		return nil, fmt.Errorf("Protobuf解析失败: %v", err)
	}

	// 提取OTP信息
	var entries []OtpEntry
	for _, otp := range payload.OtpParameters {
		// 将二进制密钥转换为base32编码
		secret := base32.StdEncoding.EncodeToString(otp.Secret)
		// 移除base32编码的填充字符
		secret = strings.ReplaceAll(secret, "=", "")

		otpType := "totp"
		if otp.Type == pb.MigrationPayload_OTP_HOTP {
			otpType = "hotp"
		}

		// 构建标准otpauth URL
		otpAuthURL := BuildOtpURL(secret, otp)

		entry := OtpEntry{
			Name:   otp.Name,
			Secret: secret,
			Issuer: otp.Issuer,
			Type:   otpType,
			URL:    otpAuthURL,
		}

		if otp.Type == pb.MigrationPayload_OTP_HOTP {
			entry.Counter = otp.Counter
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

// 构建标准otpauth URL
func BuildOtpURL(secret string, otp *pb.MigrationPayload_OtpParameters) string {
	otpType := "totp"
	if otp.Type == pb.MigrationPayload_OTP_HOTP {
		otpType = "hotp"
	}

	params := url.Values{}
	params.Set("secret", secret)
	if otp.Issuer != "" {
		params.Set("issuer", otp.Issuer)
	}
	if otp.Type == pb.MigrationPayload_OTP_HOTP {
		params.Set("counter", fmt.Sprintf("%d", otp.Counter))
	}

	return fmt.Sprintf("otpauth://%s/%s?%s",
		otpType,
		url.QueryEscape(otp.Name),
		params.Encode())
}
