package main

import (
	_ "quickstart/models"
	_ "quickstart/routers"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		orm.Debug = true
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
