package mysql

import (
	. "GoServer/utils"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"math/rand"
	"time"
)

var DBInstance *gorm.DB

func init() {
	var err error
	rand.Seed(time.Now().Unix())
	connString := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", GetMysql().Name, GetMysql().Pwsd, GetMysql().Host, GetMysql().Port, GetMysql().Basedata)
	DBInstance, err = gorm.Open("mysql", connString)

	if err != nil {
		fmt.Println("init MySQL db failed in %s, %s", connString, err)
		return
	}

	DBInstance.LogMode(GetMysql().Debug)
	DBInstance.SingularTable(true)
	//开启连接池
	DBInstance.DB().SetMaxIdleConns(10)     //最大空闲连接
	DBInstance.DB().SetMaxOpenConns(100)    //最大连接数
	DBInstance.DB().SetConnMaxLifetime(120) //最大生存时间(s)
}

func SqlTime(t time.Time) string {
	return t.Format(GetSystem().Timeformat)
}

func IsRecordNotFound(err error) bool {
	if err == gorm.ErrRecordNotFound {
		return true
	}
	return false
}

func DBClose() {
	fmt.Println("Close Mysql")
	if DBInstance != nil {
		DBInstance.Close()
	}
}
