package db

import (
	"auth/model"
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() {
	var err error
	// 使用纯 Go SQLite 库
	DB, err = sql.Open("sqlite", "file:data.db")
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	// 执行数据库初始化操作
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS secret (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		account_type INTEGER NOT NULL,
		account_name TEXT NOT NULL,
		server_name TEXT,
		encrypted_secret TEXT NOT NULL
	)`)
	if err != nil {
		log.Fatalf("创建表失败: %v", err)
	}

	log.Println("数据库初始化成功")

}

func InsertSecret(accountType int, accountName string, serverName string, encryptedSecret string) error {
	if DB == nil {
		return sql.ErrConnDone // 数据库未初始化
	}

	_, err := DB.Exec(`INSERT INTO secret (account_type, account_name, server_name, encrypted_secret) VALUES (?,?,?, ?)`,
		accountType, accountName, serverName, encryptedSecret)
	if err != nil {
		log.Printf("插入失败: %v\n", err)
		return err
	}

	log.Println("插入成功")
	return nil
}

func GetSecretsList() ([]model.Secret, error) {
	if DB == nil {
		return nil, sql.ErrConnDone // 数据库未初始化
	}
	rows, err := DB.Query("SELECT id,account_type,account_name, server_name, encrypted_secret FROM secret")

	if err != nil {
		return nil, err
	}

	var secrets []model.Secret
	for rows.Next() {
		var secret model.Secret
		err := rows.Scan(&secret.ID, &secret.AccountType, &secret.AccountName, &secret.ServerName, &secret.EncryptedSecret)
		if err != nil {
			return nil, err
		}
		secrets = append(secrets, secret)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return secrets, nil
}

func DeleteSecret(ids []int) error {
	if DB == nil {
		return sql.ErrConnDone // 数据库未初始化
	}

	// 创建占位符字符串
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?" // 每个 id 对应一个 ?
		args[i] = id          // 把 id 放到 args 中
	}

	// 将占位符拼接成 SQL
	query := fmt.Sprintf("DELETE FROM secret WHERE id IN (%s)", strings.Join(placeholders, ","))

	// 执行删除操作
	_, err := DB.Exec(query, args...)

	if err != nil {
		log.Printf("删除失败: %v\n", err)
		return err
	}

	return nil
}

func UpdateSecret(id int, accountName string, serverName string, accountType int) error {
	_, err := DB.Exec("UPDATE secret SET account_name = ?, server_name = ?, account_type = ? WHERE id = ?", accountName, serverName, accountType, id)

	if err != nil {
		log.Printf("编辑失败: %v\n", err)
		return err
	}

	return nil

}
