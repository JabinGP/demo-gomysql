package main

import (
	"fmt"
	"log"

	"encoding/json"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// DbModel 数据库表db_models（默认会加复数）对应的结构体
type DbModel struct {
	gorm.Model // 包含ID，创建时间，更新时间，删除时间（由于有这个，删除时会自动变成软删除而不是真的吧数据删掉）
	Title      string
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

	//  获取数据库连接实例
	db, err := gorm.Open(dbType, dbURL)
	if err != nil {
		log.Printf("Open mysql failed,err:%v\n", err)
		panic(err)
	}
	defer db.Close()

	// 自动根据结构体创建表或者添加表字段
	db.AutoMigrate(&DbModel{})

	// 插入数据
	db.Create(&DbModel{Title: "using gorm creat 1"})
	db.Create(&DbModel{Title: "using gorm creat 2"})

	// 读取数据，单个
	var dbData1, dbData2 DbModel
	db.First(&dbData1, "title = ?", "using gorm creat 1") // find data with title
	db.First(&dbData2, "title = ?", "using gorm creat 2")

	// 读取数据，多个
	var dbDataList []DbModel
	db.Find(&dbDataList)

	// 输出读取的数据
	if jsonString, err := json.Marshal(dbData1); err == nil {
		log.Println(string(jsonString))
	}
	if jsonString, err := json.Marshal(dbData2); err == nil {
		log.Println(string(jsonString))
	}
	if jsonString, err := json.Marshal(dbDataList); err == nil {
		log.Println(string(jsonString))
	}

	// 更新数据
	db.Model(&dbData1).Update("title", "updating by gorm")
	log.Println("-----after update-----")
	db.Find(&dbDataList) // find all
	if jsonString, err := json.Marshal(dbDataList); err == nil {
		log.Println(string(jsonString))
	}

	// 删除数据
	db.Delete(&dbData2)
	log.Println("-----after delete part-----")
	db.Find(&dbDataList) // find all
	if jsonString, err := json.Marshal(dbDataList); err == nil {
		log.Println(string(jsonString))
	}

	// 删除表，当没有限制条件的时候
	// gorm会删除整个表
	// 同时由于结构体中有gorm.Model，会根据其中的deletedAt来软删除数据
	db.Delete(&DbModel{})
	log.Println("-----after delete all-----")
	db.Find(&dbDataList) // find all
	if jsonString, err := json.Marshal(dbDataList); err == nil {
		log.Println(string(jsonString))
	}
}
