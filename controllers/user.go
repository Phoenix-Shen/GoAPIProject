package controllers

import (
	"encoding/json"
	"quickstart/models"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

// Operations about Users
type UserController struct {
	beego.Controller
}

// @Title CreateUser
// @Description create users
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {int} models.User.Id
// @Failure 403 body is empty
// @router / [post]
func (u *UserController) Post() {
	var user models.User
	json.Unmarshal(u.Ctx.Input.RequestBody, &user)
	logs.Info("请求体是", string(u.Ctx.Input.RequestBody))
	o := orm.NewOrm()
	_, err := o.Insert(&user)
	if err != nil {
		logs.Info(err.Error())
	}

	err = o.Read(&user)
	if err != nil {
		logs.Info(err.Error())
	}

	u.Data["json"] = map[string]int32{"uid": user.Id}
	u.ServeJSON()
}

// @Title GetAll
// @Description get all Users
// @Success 200 {object} models.User
// @router / [get]
func (u *UserController) GetAll() {
	o := orm.NewOrm()
	var users []*models.User
	_, err := o.QueryTable("user").All(&users)
	if err != nil {
		logs.Info(err.Error())
	}
	u.Data["json"] = users
	u.ServeJSON()
}

// @Title Get
// @Description get user by uid
// @Param	uid		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.User
// @Failure 403 :uid is empty
// @router /:uid [get]
func (u *UserController) Get() {
	uid, _ := u.GetInt32(":uid")
	o := orm.NewOrm()
	user := models.User{Id: uid}

	if uid != 0 {
		err := o.Read(user, "Id")
		if err != nil {
			u.Data["json"] = err.Error()
		} else {
			u.Data["json"] = user
		}
	}
	u.ServeJSON()
}

// @Title Update
// @Description update the user
// @Param	uid		path 	string	true		"The uid you want to update"
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {object} models.User
// @Failure 403 :uid is not int
// @router /:uid [put]
func (u *UserController) Put() {
	uid, _ := u.GetInt32(":uid")
	if uid != 0 {
		var user models.User
		json.Unmarshal(u.Ctx.Input.RequestBody, &user)
		user.Id = uid
		o := orm.NewOrm()
		uu, err := o.Update(&user)
		if err != nil {
			u.Data["json"] = err.Error()
		} else {
			u.Data["json"] = uu
		}
	}
	u.ServeJSON()
}

// @Title Delete
// @Description delete the user
// @Param	uid		path 	string	true		"The uid you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 uid is empty
// @router /:uid [delete]
func (u *UserController) Delete() {
	uid, _ := u.GetInt32(":uid")
	o := orm.NewOrm()
	_, err := o.Delete(&models.User{Id: uid})
	if err != nil {
		logs.Info(err.Error())
	}
	u.Redirect("http://127.0.0.1:8080/v1/user", 302)
	//u.Data["json"] = "delete success!"
	//u.ServeJSON()
}

// @Title Login
// @Description Logs user into the system
// @Param	username		query 	string	true		"The username for login"
// @Param	password		query 	string	true		"The password for login"
// @Success 200 {string} login success
// @Failure 403 user not exist
// @router /login [get]
func (u *UserController) Login() {
	username := u.GetString("username")
	password := u.GetString("password")
	var user models.User = models.User{Username: username}
	o := orm.NewOrm()
	err := o.Read(&user, "username")
	if err != nil {
		logs.Info(err.Error())
		u.Data["json"] = "不存在此用户"
	}

	if password != user.Password {
		u.Data["json"] = "密码错误 我叼你妈的"
	}

	u.Data["json"] = user
	u.ServeJSON()
}

// @Title logout
// @Description Logs out current logged in user session
// @Success 200 {string} logout success
// @router /logout [get]
func (u *UserController) Logout() {
	u.Data["json"] = "logout success"
	u.ServeJSON()
}
