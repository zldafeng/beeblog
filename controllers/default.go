package controllers
//
//import (
//	"github.com/astaxie/beego"
//)
//
//type MainController struct {
//	beego.Controller
//}
//
//func (c *MainController) Get() {
//	c.Data["Website"] = "beego.me"
//	c.Data["Email"] = "astaxie@gmail.com"
//	//c.TplName = "index.tpl"
//	c.TplName = "index.tpl"
//
//	c.Data["TrueCond"] = true
//	c.Data["FalseCond"] = false
//
//	// 结构
//	type u struct {
//		Name string
//		Age  int
//		Sex  string
//	}
//	user := &u{
//		Name: "DaFeng",
//		Age:  20,
//		Sex:  "Man",
//	}
//	c.Data["User"] = user
//
//
//	// 数组
//	nums :=[]int{1,2,3,4,5,6,7,8,9,0}
//	c.Data["Nums"] = nums
//
//	// 模版变量
//	c.Data["TplVar"] = "hey boys"
//
//	// 模版函数
//	c.Data["Html"] = "<div>hello beeblog</div>"
//
//	c.Data["Pipe"] = "<div>hello beeblog</div>"
//
//	// 模版嵌套
//
//}
