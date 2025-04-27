package main

import (
	"auth/db"
	"auth/model"
	"auth/utils"
	"context"
	"github.com/pquerna/otp/totp"
	"log"
	"time"
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

	return nil

}
