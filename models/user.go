package models

import (
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Id       int32  `orm:"pk;auto"`
	Username string `orm:"unique"`
	Password string `orm:"size(8)"`
	//Profile  *Profile `orm:"rel(fk)"`
}

/*type Profile struct {
	Gender  string
	Age     int
	Address string
	Email   string
}*/

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:asdf@/student?charset=utf8")
	orm.RegisterModel(new(User))
	orm.RunSyncdb("default", false, false)
}
