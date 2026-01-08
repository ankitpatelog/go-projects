package config

import "database/sql"
import "time"
import "log"

var DB *sql.DB

func Connectmysql() {
	connectionstr :="root:root@123@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	db,err  := sql.Open("mysql",connectionstr)
	if err!=nil {
		log.Fatal("mysql connection failed")
		return
	}

	// Connection pool settings
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to connect to MySQL:", err)
	}

	DB = db
	log.Println("✅ MySQL connected")

}