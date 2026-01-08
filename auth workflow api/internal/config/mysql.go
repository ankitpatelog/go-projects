package config

import "database/sql"
import "log"

var DB *sql.DB

func connectmysql()  {
	connectionstr :="root:root@123@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	db,err  := sql.Open("mysql",connectionstr)
	if err!=nil {
		log.Fatal("mysql connection failed")
		return
	}

	if db.Ping()!=nil {
		log.Fatal("connection mysql faild")
	}

	DB=db
}