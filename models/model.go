package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type User struct {
	Id int
	UserName string
	Password string
	//多对多 关系生成表
	Articles []*Article `orm:"rel(m2m)"`
}

//用户（阅读人）User和文章Article是多对多关系
//一篇文章可被多人阅读，一个用户也可阅读多篇文章
type Article struct {
	Id int `orm:"pk,auto"`
	//允许标题为空
	Title string `orm:"size(128);null"`
	Content string `orm:"size(3000);column(content)"`
	Img string `orm:"size(512);null"`
	CreateTime time.Time `orm:"type(datetime);auto_now_add"`
	ViewCount int `orm:"default(0)"`
	ArticleType *ArticleType `orm:"rel(fk)"`
	Users []*User `orm:"reverse(many)"`
}

//类型 → 文章 1对多
type ArticleType struct {
	Id int
	TypeName string `orm:"size(20)"`
	Articles []*Article `orm:"reverse(many)"`
}

func init(){
	orm.RegisterDataBase(
		"default",
		"mysql",
		"root:123456@tcp(127.0.0.1:3306)/webarticles?charset=utf8")
	orm.RegisterModel(
		new(User),
		new(Article),
		new(ArticleType))
	orm.RunSyncdb("default",false,true)
}