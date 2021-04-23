package data

import (
	"crypto/rand"
	"crypto/sha1"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var Db *sql.DB

// 初始化数据库
func init() {
	// 错误
	var err error
	// 连接数据库
	Db, err = sql.Open("postgres", "port=5432 user=postgres password=pygosuperman dbname=chitchat sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	// return
}

// 创建UUID
func createUUID() (uuid string) {
	// 创建16个字节的切片
	u := new([16]byte)
	// 随机读取
	_, err := rand.Read(u[:])
	if err != nil {
		log.Fatalln("无法创建UUID", err)
	}
	// 创建UUID的核心算法
	u[8] = (u[8] | 0x40) & 0x7F
	u[6] = (u[6] & 0xF) | (0x4 << 4)
	uuid = fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
	return
}

// 对文本进行hash
func Encrypt(plaintext string) (cryptext string) {
	cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
	return
}
