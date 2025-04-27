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
	"time"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/pquerna/otp"

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

func (a *App) DeleteSecret(ids []uint) error {
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
		return fmt.Errorf("failed to decode image: %w", err)
	}

	// prepare BinaryBitmap
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return fmt.Errorf("failed to create BinaryBitmap: %w", err)
	}

	// decode image
	qrReader := qrcode.NewQRCodeReader()
	result, err := qrReader.Decode(bmp, nil)

	if err != nil {
		return fmt.Errorf("failed to decode QR code: %w", err)
	}

	fmt.Println(result)

	key, err := otp.NewKeyFromURL(result.GetText())
	if err != nil {
		return fmt.Errorf("failed to get otp info: %w", err)
	}

	fmt.Println("Account:", key.AccountName())
	fmt.Println("Issuer:", key.Issuer())
	fmt.Println("Secret:", key.Secret())

	// 添加解析出的验证码信息
	err = a.InsertSecret(key.AccountName(), key.Issuer(), key.Secret(), 0)
	if err != nil {
		return fmt.Errorf("failed to add account: %w", err)
	}

	return nil
}
