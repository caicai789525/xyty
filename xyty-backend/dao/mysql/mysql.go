package mysql

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	model "ini/model/user_struct"
	"ini/services/parseyaml"
)

var DB *gorm.DB

func MysqlInit() {
	v := parseyaml.GetYaml()
	username := v.GetString("db.username")
	password := v.GetString("db.password")
	addr := v.GetString("db.addr")
	port := v.GetInt("db.port")
	dbname := v.GetString("db.dbname")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, addr, port, dbname)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	DB = db
	err = DB.AutoMigrate(&model.User{}, &model.UserCode{})
	if err != nil {
		panic(err)
		return
	}
}
