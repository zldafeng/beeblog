package models

import (
	"time"
	"github.com/Unknown/com"
	_"github.com/mattn/go-sqlite3"
	"os"
	"path"
	"github.com/astaxie/beego/orm"
	"strconv"
	"fmt"
	"strings"
	"github.com/astaxie/beego"
)

const (
	_DB_NAME        = "data/beeblog.db"
	_SQLITE3_DRIVER = "sqlite3"
)

// 分类
type Category struct {
	Id              int64
	Title           string
	Created         time.Time `orm:"index"`
	Views           int64     `orm:"index"`
	TopicTime       time.Time `orm:"index"`
	TopicCount      int64
	TopicLastUserId int64
}

// 文章
type Topic struct {
	Id              int64
	Uid             int64
	Title           string
	Category        string
	Labels          string
	Content         string `orm:"size(5000)"`
	Attachment      string
	Created         time.Time `orm:"index"`
	Updated         time.Time `orm:"index"`
	Views           int64     `orm:"index"`
	Author          string
	ReplyTime       time.Time `orm:"index"`
	ReplyCount      int64
	ReplyLastUserId int64
}

// 评论
type Comment struct {
	Id      int64
	Tid     int64
	Name    string
	Content string `orm:"size(1000)"`
	Created time.Time `orm:"index"`
}

func RegisterDB() {
	if !com.IsExist(_DB_NAME) {
		os.MkdirAll(path.Dir(_DB_NAME), os.ModePerm)
		os.Create(_DB_NAME)
	}
	orm.RegisterModel(new(Category), new(Topic), new(Comment)) //注册模型
	orm.RegisterDriver(_SQLITE3_DRIVER, orm.DRSqlite)          //注册驱动。可有可无
	// 注册默认的数据库，强制要求默认数据库为default
	orm.RegisterDataBase("default", _SQLITE3_DRIVER, _DB_NAME, 10)
}

//添加分类
func AddCategory(name string) error {
	// 获取orm对象
	o := orm.NewOrm()
	// 创建Category对象
	cate := &Category{
		Title:     name,
		Created:   time.Now(),
		TopicTime: time.Now(),
	}
	// 查询操作
	qs := o.QueryTable("category")
	// 过滤--
	err := qs.Filter("title", name).One(cate)
	if err == nil {
		return err
	}
	// 插入数据
	_, err = o.Insert(cate)
	if err != nil {
		return err
	}
	return nil
}

// 获取所有分类
func GetAllCategories() ([] *Category, error) {
	// 获取orm对象
	o := orm.NewOrm()
	cates := make([]*Category, 0)
	// 查询操作
	qs := o.QueryTable("category")
	_, err := qs.All(&cates)
	return cates, err
}

// 删除分类
func DeleteAllCategories(id string) error {
	// 10进制 64位
	cid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}

	o := orm.NewOrm()

	cate := &Category{Id: cid}
	_, err = o.Delete(cate)
	if err!=nil{
		beego.Error(err)
	}
	return err

}

// 添加文章
func AddTopic(title, category, label, content ,attachment string) error {
	// 处理标签
	label = "$" + strings.Join(
		strings.Split(label, " "), "#$") + "#"

	// "beego orm"
	// [beego orm]
	// $beego#$orm#
	o := orm.NewOrm()

	topic := &Topic{
		Title:     title,
		Category:  category,
		Labels:    label,
		Content:   content,
		Attachment:attachment,
		Created:   time.Now(),
		Updated:   time.Now(),
		ReplyTime: time.Now(),
	}
	_, err := o.Insert(topic)
	if err != nil {
		return err
	}
	// 更新分类统计
	cate := new(Category)
	qs := o.QueryTable("category")
	err = qs.Filter("title", category).One(cate)
	if err == nil {
		// 不存在，则忽略更新操作
		topics, err :=GetAllTopics(category,"",true);
		if err != nil {
			beego.Error(err)
		}
		cate.TopicCount=int64(len(topics))
		//cate.TopicCount++
		_, err = o.Update(cate)
	}
	return err
}

// 获取所有文章
func GetAllTopics(cate, label string, isDesc bool) ([] *Topic, error) {
	// 获取orm对象
	o := orm.NewOrm()
	topics := make([]*Topic, 0)
	// 查询操作
	qs := o.QueryTable("topic")
	var err error
	if isDesc {
		// 点击分类筛选
		if len(cate) > 0 {
			qs = qs.Filter("category", cate)
		}
		// 点击标签筛选
		if len(label) > 0 {
			qs = qs.Filter("labels__contains", "$"+label+"#")
		}
		// 首页-按时间倒序
		_, err = qs.OrderBy("-created").All(&topics)
	} else {
		// 文章列表页
		_, err = qs.All(&topics)
	}

	return topics, err
}

// 获取文章内容
func GetTopic(tid string) (*Topic, error) {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return nil, err
	}

	o := orm.NewOrm()
	topic := new(Topic)

	qs := o.QueryTable("topic")
	err = qs.Filter("id", tidNum).One(topic)
	if err != nil {
		return nil, err
	}
	// 更新浏览次数
	topic.Views++
	_, err = o.Update(topic)

	topic.Labels = strings.Replace(strings.Replace(topic.Labels, "#", " ", -1),
		"$", "", -1)

	return topic, err
}

// 修改文章内容
func ModifyTopic(tid, title, category, label, content,attachment string) error {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}
	// 处理标签
	label = "$" + strings.Join(
		strings.Split(label, " "), "#$") + "#"
	var oldCate ,oldAttach string
	o := orm.NewOrm()
	topic := &Topic{Id: tidNum}
	if o.Read(topic) == nil {
		oldCate = topic.Category
		oldAttach = topic.Attachment
		topic.Title = title
		topic.Category = category
		topic.Labels = label
		topic.Content = content
		topic.Attachment = attachment
		topic.Updated = time.Now()
		_, err = o.Update(topic)
		if err != nil {
			return err
		}
	}
	// 删除旧的附件
	if len(oldAttach)>0{
		os.Remove(path.Join("attachment",oldAttach))
	}
	// 修改文章分类时候----同步更新分类统计
	// 1.更新旧的分类统计 --> -1
	if len(oldCate) > 0 {
		cate := new(Category)
		qs := o.QueryTable("category")
		// 如果根据旧的分类能找到文章
		err = qs.Filter("title", oldCate).One(cate)
		if err == nil {
			// 不存在，则忽略更新操作
			topics, err :=GetAllTopics(oldCate,"",true);
			if err != nil {
				beego.Error(err)
			}
			cate.TopicCount=int64(len(topics))
			//cate.TopicCount--
			_, err = o.Update(cate)
		}
	}
	// 2.同时更新新的分类统计 -->+1
	cate := new(Category)
	qs := o.QueryTable("category")
	err = qs.Filter("title", category).One(cate)
	if err == nil {
		// 不存在，则忽略更新操作
		topics, err :=GetAllTopics(category,"",true);
		if err != nil {
			beego.Error(err)
		}
		cate.TopicCount=int64(len(topics))
		//cate.TopicCount++
		_, err = o.Update(cate)
	}
	return nil
}

// 删除文章
func DeleteTopic(tid string) error {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}

	var oldCate string
	o := orm.NewOrm()
	topic := &Topic{Id: tidNum}
	if o.Read(topic) == nil {
		oldCate = topic.Category
		_, err = o.Delete(topic)
		if err != nil {
			return err
		}
	}
	// 更新分类统计
	if len(oldCate) > 0 {
		cate := new(Category)
		qs := o.QueryTable("category")
		err = qs.Filter("title", oldCate).One(cate)
		if err == nil {
			// 不存在，则忽略更新操作
			topics, err :=GetAllTopics(oldCate,"",true);
			if err != nil {
				beego.Error(err)
			}
			cate.TopicCount=int64(len(topics))
			//cate.TopicCount--
			_, err = o.Update(cate)
		}
	}
	return err
}

// 添加评论
func AddReply(tid, nickname, content string) error {
	tidNum, err := strconv.ParseInt(tid, 10, 64)

	fmt.Println("dafeng-models.go-----", tidNum)
	if err != nil {
		return err
	}

	o := orm.NewOrm()
	reply := &Comment{
		Tid:     tidNum,
		Name:    nickname,
		Content: content,
		Created: time.Now(),
	}
	_, err = o.Insert(reply)
	if err != nil {
		return err
	}

	topic := &Topic{Id: tidNum}
	if o.Read(topic) == nil {
		topic.ReplyTime = time.Now()
		topic.ReplyCount++
		_, err = o.Update(topic)
	}
	return err
}

// 查看所有评论
func GetAllReplies(tid string) (replies []*Comment, err error) {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return nil, err
	}
	replies = make([]*Comment, 0)
	o := orm.NewOrm()
	qs := o.QueryTable("comment")
	_, err = qs.Filter("tid", tidNum).All(&replies)
	return replies, err
}

// 删除评论
func DeleteReply(tid string) error {
	ridNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}

	o := orm.NewOrm()

	var tidNum int64
	reply := &Comment{Id: ridNum}
	if o.Read(reply) == nil {
		tidNum = reply.Tid
		_, err = o.Delete(reply)
		if err != nil {
			return err
		}
	}

	replies := make([]*Comment, 0)
	qs := o.QueryTable("comment")
	_, err = qs.Filter("tid", tidNum).OrderBy("-created").All(&replies)
	if err != nil {
		return err
	}
	topic := &Topic{Id: tidNum}
	if o.Read(topic) == nil {
		topic.ReplyTime = replies[0].Created
		topic.ReplyCount = int64(len(replies))
		_, err = o.Update(topic)
	}
	return err
}
