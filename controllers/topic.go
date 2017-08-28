package controllers

import (
	"github.com/astaxie/beego"
	"beeblog/models"
	"strings"
	"path"
)

type TopicController struct {
	beego.Controller
}

func (this *TopicController) Get() {
	this.Data["IsLogin"] = checkAccount(this.Ctx)
	this.Data["IsTopic"] = true
	this.TplName = "topic.html"

	topics, err := models.GetAllTopics("","",false)
	if err != nil {
		beego.Error(err)
	} else {
		this.Data["Topics"] = topics
	}
}
func (this *TopicController) Post() {
	if !checkAccount(this.Ctx) {
		this.Redirect("/login", 302)
		return
	}
	// 解析表单
	title := this.Input().Get("title")
	content := this.Input().Get("content")
	tid := this.Input().Get("tid")
	category := this.Input().Get("category")
	label := this.Input().Get("label")

	// 获取附件
	_,fh,err := this.GetFile("attachment")
	if err !=nil{
		beego.Error(err)
	}
	var attachment string
	if fh!=nil{
		// 保存
		attachment = fh.Filename
		beego.Info("附件名称：",attachment)
		err = this.SaveToFile("attachment",
			path.Join("attachment",attachment))
	}

	if len(tid) == 0 {
		err = models.AddTopic(title, category,label,content,attachment)
	} else {
		err = models.ModifyTopic(tid, title,category,label, content,attachment)
	}

	if err != nil {
		beego.Error(err)
	}
	this.Redirect("/topic", 302)
}
// 添加文章
func (this *TopicController) Add() {
	if !checkAccount(this.Ctx) {
		this.Redirect("/login", 302)
		return
	}
	this.TplName = "topic_add.html"
}
// 查看文章
func (this *TopicController) View() {
	this.TplName = "topic_view.html"
	tid := this.Ctx.Input.Param("0")

	// 获取tid
	//reqUrl := this.Ctx.Request.RequestURI
	//i :=strings.LastIndex(reqUrl,"/")
	//tid =reqUrl[i+1:]

	topic, err := models.GetTopic(tid)
	if err != nil {
		beego.Error(err)
		this.Redirect("/", 302)
		return
	}
	this.Data["Topic"] = topic
	this.Data["Labels"] = strings.Split(topic.Labels," ")

	replies,err :=models.GetAllReplies(tid)

	if err!=nil{
		beego.Error(err)
		return
	}
	this.Data["Replies"] = replies
	this.Data["IsLogin"] = checkAccount(this.Ctx)

}
// 修改文章
func (this *TopicController) Modify() {
	this.TplName = "topic_modify.html"

	// 两种方式操作:不同链接对应不同的取值
	// xxx/{{.Id}}
	// xxx/?tid={{.Id}}
	tid := this.Ctx.Input.Param("0")
	//tid := this.Input().Get("tid")

	topic, err := models.GetTopic(tid)
	if err != nil {
		beego.Error(err)
		this.Redirect("/", 302)
		return
	}

	this.Data["Topic"] = topic
	this.Data["Tid"] = tid
}
// 删除文章
func (this *TopicController) Delete() {

	if !checkAccount(this.Ctx) {
		this.Redirect("/login", 302)
		return
	}

	// 两种方式操作--同上个方法
	//err := models.DeleteTopic(this.Ctx.Input.Param("0"))
	err := models.DeleteTopic(this.Input().Get("tid"))
	if err != nil {
		beego.Error(err)
		return
	}
	this.Redirect("/topic", 302)
}
