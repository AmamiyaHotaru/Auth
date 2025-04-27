package model

type Secret struct {
	ID              uint
	AccountName     string
	ServerName      string
	EncryptedSecret string `json:"-"`
	AccountType     uint   `json:"-"`
	Code            string
}
