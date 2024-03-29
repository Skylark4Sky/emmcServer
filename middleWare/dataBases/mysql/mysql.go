package mysql

import (
	. "GoServer/utils/config"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"math/rand"
	"time"
)

var _db *gorm.DB

func init() {
	var err error

	if mysql, _ := GetMysql(); mysql != nil {
		rand.Seed(time.Now().Unix())

		connString := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", mysql.Name, mysql.Pwsd, mysql.Host, mysql.Port, mysql.Basedata)
		_db, err = gorm.Open("mysql", connString)

		if err != nil {
			fmt.Println("init MySQL db failed in %s, %s", connString, err)
			return
		}

		_db.LogMode(mysql.Debug)
		_db.SingularTable(true)
		//开启连接池
		_db.DB().SetMaxIdleConns(10)     //最大空闲连接
		_db.DB().SetMaxOpenConns(100)    //最大连接数
		_db.DB().SetConnMaxLifetime(120) //最大生存时间(s)
	}
}

func IsRecordNotFound(err error) bool {
	if err == gorm.ErrRecordNotFound {
		return true
	}
	return false
}

func ExecSQL() *gorm.DB {
	return _db
}

func CreateSQLAndRetLastID(entity interface{}) (uint64, error) {
	var id []uint64
	tx := _db.Begin()
	if err := tx.Create(entity).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Raw("select LAST_INSERT_ID() as id").Pluck("id", &id).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	return id[0], nil
}

func TXCreateSQLAndRetLastID(tx *gorm.DB, entity interface{}) (uint64, error) {
	var id []uint64
	if err := tx.Create(entity).Error; err != nil {
		tx.Rollback()
		return 0, err
	}
	if err := tx.Raw("select LAST_INSERT_ID() as id").Pluck("id", &id).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	return id[0], nil
}

func SQLClose() {
	fmt.Println("Close Mysql")
	if _db != nil {
		_db.Close()
	}
}
