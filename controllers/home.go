package controllers

import (
	"github.com/astaxie/beego"
	"beeblog/models"
	"fmt"
)

type HomeController struct {
	beego.Controller
}

func (this *HomeController) Get() {
	this.Data["IsHome"] = true
	this.TplName = "home.html"
	this.Data["IsLogin"] = checkAccount(this.Ctx)
	fmt.Println("dafeng-home.go",this.Input().Get("label"))
	topics ,err:=models.GetAllTopics(
		this.Input().Get("cate"),this.Input().Get("label"),true)
	if err!=nil{
		beego.Error(err)
	}else{
		this.Data["Topics"] = topics
		for _,topic :=range topics{
			topic.Created.Date()
		}
	}

	categories,err :=models.GetAllCategories()

	if err!=nil{
		beego.Error(err)
	}
	this.Data["Categories"] = categories
}
