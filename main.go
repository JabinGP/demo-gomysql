package main

import (
	"fmt"
	"log"

	"encoding/json"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// DbModel database table model
type DbModel struct {
	gorm.Model
	Title string
}

func main() {
	const (
		dbType   string = "mysql"
		dbHost   string = "***"
		dbPort   string = "3306"
		dbName   string = "demo"
		dbUser   string = "***"
		dbPasswd string = "***"
		dbParams string = "charset=utf8&parseTime=true"
	)
	var dbURL = fmt.Sprintf("%s:%s@(%s:%s)/%s?%s", dbUser, dbPasswd, dbHost, dbPort, dbName, dbParams)

	//  connect database
	db, err := gorm.Open(dbType, dbURL)
	if err != nil {
		log.Printf("Open mysql failed,err:%v\n", err)
		panic(err)
	}
	defer db.Close()

	// migrate the schema (auto create table)
	db.AutoMigrate(&DbModel{})

	// create
	db.Create(&DbModel{Title: "using gorm creat 1"})
	db.Create(&DbModel{Title: "using gorm creat 2"})

	// read
	var dbData1, dbData2 DbModel
	db.First(&dbData1, "title = ?", "using gorm creat 1") // find data with title
	db.First(&dbData2, "title = ?", "using gorm creat 2")

	var dbDataList []DbModel
	db.Find(&dbDataList) // find all

	// output as json style
	if jsonString, err := json.Marshal(dbData1); err == nil {
		log.Println(string(jsonString))
	}
	if jsonString, err := json.Marshal(dbData2); err == nil {
		log.Println(string(jsonString))
	}
	if jsonString, err := json.Marshal(dbDataList); err == nil {
		log.Println(string(jsonString))
	}

	// update
	db.Model(&dbData1).Update("title", "updating by gorm")
	log.Println("-----after update-----")
	db.Find(&dbDataList) // find all
	if jsonString, err := json.Marshal(dbDataList); err == nil {
		log.Println(string(jsonString))
	}

	// delete part
	db.Delete(&dbData2)
	log.Println("-----after delete part-----")
	db.Find(&dbDataList) // find all
	if jsonString, err := json.Marshal(dbDataList); err == nil {
		log.Println(string(jsonString))
	}

	// delete all
	db.Delete(&DbModel{})
	log.Println("-----after delete all-----")
	db.Find(&dbDataList) // find all
	if jsonString, err := json.Marshal(dbDataList); err == nil {
		log.Println(string(jsonString))
	}
}
