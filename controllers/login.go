package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

type LoginController struct {
	beego.Controller
}

func (this *LoginController) Get() {

	isExit := this.Input().Get("exit") == "true"
	if isExit{
		this.Ctx.SetCookie("uname","",-1,"/")
		this.Ctx.SetCookie("pwd","",-1,"/")
		this.Redirect("/",302)
		return
	}
	this.TplName="login.html"
}
func (this *LoginController) Post() {
	//this.Ctx.WriteString(fmt.Sprint(this.Input()))
	//return
	uname :=this.Input().Get("uname")
	pwd :=this.Input().Get("pwd")
	autoLogin :=this.Input().Get("autoLogin")=="on"

	if beego.AppConfig.String("uname")==uname &&
	beego.AppConfig.String("pwd")==pwd{
		maxAge :=0
		if autoLogin{
			maxAge = 1<<31 -1
		}

		this.Ctx.SetCookie("uname",uname,maxAge,"/")
		this.Ctx.SetCookie("pwd",pwd,maxAge,"/")
	}

	this.Redirect("/",302)
	return
}
// 没有beego.Context了，新api替换为context.Context
func checkAccount(ctx *context.Context)  bool{
	ck,err :=ctx.Request.Cookie("uname")
	if err!=nil{
		return false
	}
	uname :=ck.Value
	ck,err =ctx.Request.Cookie("pwd")
	if err!=nil{
		return false
	}
	pwd :=ck.Value

	return beego.AppConfig.String("uname")==uname &&
		beego.AppConfig.String("pwd")==pwd
}
