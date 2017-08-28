package controllers

import (
	"github.com/astaxie/beego"
	"net/url"
	"os"
	"io"
)

type AttachController struct {
	beego.Controller
}

func (this *AttachController) Get() {
	// 中文编码
	filePath,err:=url.QueryUnescape(this.Ctx.Request.RequestURI[1:])
	if err!=nil{
		this.Ctx.WriteString(err.Error())
		return
	}
	f,err :=os.Open(filePath)
	if err !=nil{
		this.Ctx.WriteString(err.Error())
		return
	}
	defer f.Close()

	// 输入流，输出流
	_,err = io.Copy(this.Ctx.ResponseWriter,f)
	if err !=nil{
		this.Ctx.WriteString(err.Error())
		return
	}
}
