package middleWare

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
	connString := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", MysqlConf().Name, MysqlConf().Pwsd, MysqlConf().Host, MysqlConf().Port, MysqlConf().Basedata)
	DBInstance, err = gorm.Open("mysql", connString)
	if err != nil {
		fmt.Errorf("init MySQL db failed in %s, %s", connString, err)
		return
	}

	DBInstance.LogMode(MysqlConf().Debug)
	DBInstance.SingularTable(true)
}

func SqlTime(t time.Time) string {
	return t.Format(SystemConf().Timeformat)
}

func IsRecordNotFound (err error) bool {
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

