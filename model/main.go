package model

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDB()  *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:@tcp(127.0.0.1:3306)/fiber"), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	return db
}

func CloseDB(db *gorm.DB) error {
	mysql, err := db.DB()

	if err != nil {
		panic(err)
	}

	return mysql.Close()
}

func CreateTable() {
	db  := ConnectDB()
	defer CloseDB(db)

	db.Exec(`CREATE TABLE IF  NOT EXISTS blogs (
		id INTEGER PRIMARY KEY AUTO_INCREMENT,
		title VARCHAR(255) NOT NULL,
		description VARCHAR(255) NOT NULL,
		createdAt DATETIME NOT NULL DEFAULT NOW(),
		updatedAt DATETIME NOT NULL DEFAULT NOW() ON UPDATE NOW()
	)`)
}