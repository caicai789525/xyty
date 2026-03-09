package main

import (
	"ini/dao/mysql"
	"ini/router"
	"log"
)

// @title oyhx API
// @version 1.0
// @description 芗音同韵API
// @termsOfService http://swagger.io/terms/
// @contact.name KitZhangYs
// @contact.email SJMbaiyang@163.com
// @host 8.146.198.169
// @BasePath /api/v1
func main() {
	mysql.MysqlInit()
	e := router.RouterInit()
	err := e.Run(":8080")
	if err != nil {
		log.Panic(err)
		return
	}
}
