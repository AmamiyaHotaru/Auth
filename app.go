package main

import (
	"auth/db"
	"auth/model"
	"auth/utils"
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/gif"  // 支持GIF格式
	_ "image/jpeg" // 支持JPEG格式
	_ "image/png"  // 支持PNG格式
	"log"
	"strings"
	"time"

	"github.com/pquerna/otp"

	gotp "auth/utils/otp_extractor"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/pquerna/otp/totp"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetSecretsList() []model.Secret {
	secrets, err := db.GetSecretsList()

	if err != nil {
		log.Println("数据库查询失败", err)
		return nil
	}

	for i := range secrets {
		decryptedSecret, err := utils.Decrypt(secrets[i].EncryptedSecret)
		if err != nil {
			log.Printf("解密失败: %v, AccountName: %s\n", err, secrets[i].AccountName)
			continue
		}
		var code string
		switch secrets[i].AccountType {
		case 0: // totp
			code, err = totp.GenerateCode(decryptedSecret, time.Now())
			if err != nil {
				log.Printf("生成 TOTP 验证码失败: %v, AccountName: %s\n", err, secrets[i].AccountName)
				continue
			}
		case 1: // steam
			code, err = utils.GenerateCode(string(decryptedSecret))
			if err != nil {
				log.Printf("生成 Steam 验证码失败: %v, AccountName: %s\n", err, secrets[i].AccountName)
				continue
			}
		default:
			continue
		}

		if err != nil {
			log.Println("生成验证码失败", err)
			continue
		}
		secrets[i].Code = code

	}

	return secrets
}

func (a *App) InsertSecret(accountName string, serverName string, secret string, accountType int) error {
	encryptedSecret, err := utils.Encrypt([]byte(secret))
	if err != nil {
		log.Println("加密失败", err)
		return err
	}

	log.Printf("accountName：%s,serverName:%s,accountType：%d,加密后密钥: %s\n", accountName, serverName, accountType, encryptedSecret)
	err = db.InsertSecret(accountType, accountName, serverName, encryptedSecret)
	if err != nil {
		log.Println("添加失败", err)
	}

	return nil

}

func (a *App) DeleteSecret(ids []int) error {
	err := db.DeleteSecret(ids)
	if err != nil {
		log.Println("删除失败", err)
	}
	return nil
}
func (a *App) RecognizeQRCode(imgBytes []byte) error {
	reader := bytes.NewReader(imgBytes)
	img, format, err := image.Decode(reader)
	if err != nil {
		log.Printf("图像解码失败: %v, 格式: %s\n", err, format)
		return fmt.Errorf("无法解码图像: %w", err)
	}

	// prepare BinaryBitmap
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return fmt.Errorf("无法创建 BinaryBitmap: %w", err)
	}

	// decode image
	qrReader := qrcode.NewQRCodeReader()
	result, err := qrReader.Decode(bmp, nil)

	if err != nil {
		return fmt.Errorf("解析二维码错误: %w,%s", err, "尝试单独截取二维码部分的图片进行解析")
	}

	fmt.Println(result)

	switch {
	case strings.HasPrefix(result.GetText(), "otpauth-migration://"):
		// 处理 谷歌认证器 otpauth-migration 前缀

		entries, err := gotp.ExtractOtpFromUrl(result.GetText())
		if err != nil {
			return fmt.Errorf("解析谷歌验证器导出码失败: %w", err)
		}

		// 打印提取的OTP信息
		for i, entry := range entries {
			otpType := 0
			fmt.Printf("\n%d. 密钥\n", i+1)
			fmt.Printf("名称:    %s\n", entry.Name)
			fmt.Printf("密钥:    %s\n", entry.Secret)
			if entry.Issuer != "" {
				fmt.Printf("发行者:  %s\n", entry.Issuer)
			}
			fmt.Printf("类型:    %s\n", entry.Type)
			if entry.Type == "hotp" {
				otpType = 2
				fmt.Printf("计数器:  %d\n", entry.Counter)
			}

			err = a.InsertSecret(entry.Name, entry.Issuer, entry.Secret, otpType)
			if err != nil {
				return fmt.Errorf("添加账户失败: %w", err)
			}
		}

	case strings.HasPrefix(result.GetText(), "otpauth://totp"):
		// 处理 totp 前缀
		key, err := otp.NewKeyFromURL(result.GetText())
		if err != nil {
			return fmt.Errorf("获取 OTP 信息失败: %w", err)
		}

		fmt.Println("Account:", key.AccountName())
		fmt.Println("Issuer:", key.Issuer())
		fmt.Println("Secret:", key.Secret())

		// 添加解析出的验证码信息
		err = a.InsertSecret(key.AccountName(), key.Issuer(), key.Secret(), 0)
		if err != nil {
			return fmt.Errorf("添加账户失败: %w", err)
		}

	//case strings.HasPrefix(result.GetText(), "yet-another-prefix://"):
	//	// 处理 yet-another-prefix 前缀
	default:
		return fmt.Errorf("二维码格式错误: %w", err)
	}

	return nil
}

func (a *App) UpdateSecret(id int, accountName string, serverName string, accountType int) error {
	err := db.UpdateSecret(id, accountName, serverName, accountType)
	if err != nil {
		log.Println("编辑失败", err)
	}
	return nil

}
