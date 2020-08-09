package controllers

import (
	"article-manage-system/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)


type RegController struct {
	beego.Controller
}

func (this * RegController) ShowRegister(){
	this.TplName="register.html"
}

//处理注册逻辑
func (this *RegController) HandleRegister(){
	//1.拿到浏览器传递的数据
	username:=this.GetString("userName")
	pwd:=this.GetString("password")
	if username == "" || pwd == ""{
		this.Data["errormsg"]="用户名或密码不能为空"
		this.TplName="register.html"
		return
	}
	//插入数据库
	om:=orm.NewOrm()
	userModel:=models.User{UserName:username,Password:pwd}
	n,err:=om.Insert(&userModel)
	if err!=nil{
		this.Data["errormsg"]="注册失败"
		this.TplName="register.html"
		return
	}
	if n>0{
		this.Redirect("/login",302)
	}else{
		this.TplName="register.html"
	}

}

type LoginController struct {
	beego.Controller
}

func (this *LoginController) ShowLogin(){
	username:=this.Ctx.GetCookie("userName")
	if username !=""{
		this.Data["userName"]=username
		this.Data["check"]="checked"
	}
	this.TplName="login.html"
}

func (this *LoginController) HandleLogin(){
	username:=this.GetString("userName")
	pwd:=this.GetString("password")
	if username=="" || pwd==""{
		beego.Info("用户名或密码不能为空")
		this.TplName="login.html"
		return
	}
	om:=orm.NewOrm()
	user:=models.User{UserName:username}
	//指定根据哪个字段查询
	err:=om.Read(&user,"UserName")
	if err!=nil{
		beego.Info("用户名错误")
		this.TplName = "login.html"
		return
	}
	//查询成功，判断用户的密码是否正确
	if pwd!=user.Password{
		beego.Info("密码错误")
		this.TplName = "login.html"
		return
	}
	isRemeber:=this.GetString("remeber")
	if isRemeber=="on"{
		//用户名cookie保存一小时
		this.Ctx.SetCookie("userName",username,time.Hour)
	}else{
		//清除cookie
		this.Ctx.SetCookie("userName", username, -1)
	}

	//设置session
	this.SetSession("UserName",username)

	//登录成功
	this.Redirect("/article/ShowArticles",302)
}

func (this *LoginController) HandleLogout(){
	//删除session
	this.DelSession("UserName")
	this.Redirect("/login",302)
}
