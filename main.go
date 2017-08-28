package main

import (
	"github.com/astaxie/beego"
	"beeblog/models"
	"beeblog/controllers"
	"github.com/astaxie/beego/orm"
	"os"
)

func init()  {
	models.RegisterDB()
}
func main() {
	//打印详细信息，开发模式
	orm.Debug = true
	// 自动建表1.表名。 2.若为true，每次都是删除重建。3.是否打印相关信息
	// 部署之后全部设置为false
	orm.RunSyncdb("default",false,true)

	// 注册首页路由
	beego.Router("/",&controllers.HomeController{})
	// 注册登录的路由
	beego.Router("/login",&controllers.LoginController{})
	// 注册分类的路由
	beego.Router("/category",&controllers.CategoryController{})
	// 注册文章的路由
	// 自动路由
	beego.AutoRouter(&controllers.TopicController{})
	beego.Router("/topic",&controllers.TopicController{})
	// 添加评论的路由
	beego.Router("/reply",&controllers.ReplyController{})
	beego.Router("/reply/add",&controllers.ReplyController{},"post:Add")
	beego.Router("/reply/delete",&controllers.ReplyController{},"get:Delete")

	// 创建附件目录
	os.Mkdir("attachment",os.ModePerm)
	//// 作为静态文件
	//beego.SetStaticPath("/attachment","attachment")
	// 作为单独一个控制器来处理
	beego.Router("/attachment/:all",&controllers.AttachController{})
	// 启动beego
	beego.Run()

}

