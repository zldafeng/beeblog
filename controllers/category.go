package controllers

import (
	"github.com/astaxie/beego"
	"beeblog/models"
	"fmt"
)

type CategoryController struct {
	beego.Controller
}

func (this *CategoryController) Get() {

	//op := this.GetString("op")
	op := this.Input().Get("op")
	fmt.Println("dafeng--category.go=====op",op)
	switch op {
	case "add":
		name := this.Input().Get("name")
		if len(name) == 0 {
			break
		}
		err := models.AddCategory(name)
		if err != nil {
			beego.Error(err)
		}
		this.Redirect("/category", 302)
		return
	case "del":
		id := this.Input().Get("id")
		fmt.Println("dafeng--category.go=====id",id)
		if len(id) == 0 {
			break
		}
		err :=models.DeleteAllCategories(id)
		if err!=nil{
			beego.Error(err)
		}

		this.Redirect("/category",302)
		return
	}

	this.Data["IsCategory"] = true
	this.TplName = "category.html"
	this.Data["IsLogin"] = checkAccount(this.Ctx)
	var err error
	// this.Data 已经声明类型，不能更改，故不能用:=
	this.Data["Categories"], err = models.GetAllCategories()
	if err != nil {
		beego.Error(err)
	}
}
