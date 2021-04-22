package controllers

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"os"
	"quickstart/models"
	"quickstart/mongodb"
	"strconv"
	"time"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func writeTxt(id int32) {
	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		logs.Info("读取文件错误：", err)
	}

	file.WriteString(strconv.FormatInt(int64(id), 10))

	err = file.Close()
	if err != nil {
		logs.Info("写入文件失败：", err.Error())
	}
}

//var uri string = "mongodb+srv://root:asdf@ssk.3hxej.mongodb.net/GOAPIPROJDB?retryWrites=true&w=majority&authSource=admin"
var collection *mongo.Collection
var Id int32 = 0
var filepath = "conf/config.txt"

//获取连接
func init() {
	collection, _ = mongodb.ConnectToMongoDB(mongodb.Uri, "GOAPIPROJ", 10*time.Second, 20, "Users")
	//初始化ID自增操作

	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		logs.Info("读取文件错误：", err)
		return
	}

	reader := bufio.NewReader(file)
	str, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		logs.Info("读取文件时候出现了错误" + err.Error())
		return
	}
	Id64, _ := strconv.ParseInt(str, 10, 32)
	Id = int32(Id64)
	err = file.Close()
	if err != nil {
		logs.Info("关闭文件时候出现了错误" + err.Error())
		return
	}

	logs.Info("读取上次的ID成功，id是" + strconv.FormatInt(Id64, 10))
}

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
	//取得请求体里面的user
	var user models.User
	json.Unmarshal(u.Ctx.Input.RequestBody, &user)

	logs.Info("请求体是", string(u.Ctx.Input.RequestBody))

	//时代变了 不用这个代码了
	/*o := orm.NewOrm()
	_, err := o.Insert(&user)
	if err != nil {
		logs.Info(err.Error())
	}

	err = o.Read(&user)
	if err != nil {
		logs.Info(err.Error())
	}*/

	//完成userID自增
	user.Id = Id
	Id = Id + 1

	users := mongodb.Read(collection, bson.M{"username": user.Username})
	if users != nil {
		logs.Info("duplicate name!")
		u.Data["json"] = "duplicate name!"
		u.ServeJSON()
		Id = Id - 1
		return
	}
	//保存自增的ID
	writeTxt(Id)
	result := mongodb.Create(collection, []interface{}{user})

	u.Data["json"] = map[string]interface{}{"objID": result}
	u.ServeJSON()
}

// @Title GetAll
// @Description get all Users
// @Success 200 {object} models.User
// @router / [get]
func (u *UserController) GetAll() {

	//var users []*models.User
	/*
		o := orm.NewOrm()
		_, err := o.QueryTable("user").All(&users)
		if err != nil {
			logs.Info(err.Error())
		}*/
	var users []*models.User
	cursor, _ := collection.Find(context.TODO(), bson.D{})
	cursor.All(context.TODO(), &users)
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
	user := models.User{Id: uid}
	/*时代变了不用这个代码了
	o := orm.NewOrm()
	if uid != 0 {
		err := o.Read(user, "Id")
		if err != nil {
			u.Data["json"] = err.Error()
		} else {
			u.Data["json"] = user
		}
	}
	u.ServeJSON()*/
	cursor, err := collection.Find(context.TODO(), bson.M{"id": user.Id})
	cursor.All(context.TODO(), &user)
	if err != nil {
		u.Data["json"] = err.Error()
	} else {
		u.Data["json"] = user
	}
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
		/*
			o := orm.NewOrm()
			uu, err := o.Update(&user)
			if err != nil {
				u.Data["json"] = err.Error()
			} else {
				u.Data["json"] = uu
			}
		*/
		updateResult, err := collection.UpdateOne(context.TODO(), bson.M{"id": user.Id}, &user)
		if err != nil {
			u.Data["json"] = err.Error()
		} else {
			u.Data["json"] = "操作成功，修改的行数是：" + strconv.Itoa(int(updateResult.ModifiedCount))
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
	/*
		o := orm.NewOrm()
		_, err := o.Delete(&models.User{Id: uid})
		if err != nil {
			logs.Info(err.Error())
		}
	*/
	deleteResult, err := collection.DeleteOne(context.TODO(), bson.M{"id": uid})
	if err != nil {
		u.Data["json"] = err.Error()
	} else {
		u.Data["json"] = "删除的条数是" + strconv.FormatInt(deleteResult.DeletedCount, 10)
	}
	//u.Redirect("http://127.0.0.1:8080/v1/user", 302)
	//u.Data["json"] = "delete success!"
	u.ServeJSON()
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
	/*o := orm.NewOrm()
	err := o.Read(&user, "username")
	if err != nil {
		logs.Info(err.Error())
		u.Data["json"] = "不存在此用户"
	}

	if password != user.Password {
		u.Data["json"] = "密码错误 我叼你妈的"
	}
	*/
	err := collection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		logs.Info(err.Error())
		u.Data["json"] = err.Error() + "不存在此用户，滚！"
	} else {
		if user.Password != password {
			u.Data["json"] = "密码错误，你再输一次试试？"
		} else {
			u.Data["json"] = user
		}
	}

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

// @Title deleteall
// @Description 走人辣
// @Success 200 {string} 可以跑路了
// @router /deleteall [get]
func (u *UserController) Deleteall() {
	mongodb.S删库跑路(collection)
	u.Data["json"] = "可以跑路了"
	u.ServeJSON()
}

var length int64 = 10

// @Title 让张展玮的鸡儿减少1cm
// @Description 让张展玮的鸡儿减少1cm
// @Success 200 {string} 让张展玮的鸡儿减少1cm
// @router /decline [get]
func (u *UserController) Decline() {
	length = length - 1
	u.Data["json"] = "张展玮的鸡儿长度减少了1cm，现在的长度是:" + strconv.FormatInt(length, 10) + "cm"
	u.ServeJSON()
}
