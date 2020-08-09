package routers

import (
	"article-manage-system/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	//过滤校验
	beego.InsertFilter("/article/*",beego.BeforeExec,FilterMethod)

	beego.Router("/", &controllers.MainController{})

	//注册页面
	beego.Router("/register", &controllers.RegController{}, "get:ShowRegister;post:HandleRegister")
	beego.Router("/login", &controllers.LoginController{}, "get:ShowLogin;post:HandleLogin")
	//退出
	beego.Router("/logout", &controllers.LoginController{}, "get:HandleLogout")

	//显示文章列表
	beego.Router("/article/ShowArticles", &controllers.ArticleController{}, "get:ShowArticleList")
	//文章详情
	beego.Router("/article/ShowDetail", &controllers.ArticleController{}, "get:ShowArticleDetail")

	//编辑文章
	beego.Router("/article/UpdateArticle", &controllers.ArticleController{}, "get:ShowUpdateArticle;post:HandleUpdateArticle")

	//删除文章
	beego.Router("/article/DeleteArticle", &controllers.ArticleController{}, "get:DeleteArticle")

	//添加文章
	beego.Router("/article/AddArticle", &controllers.ArticleController{}, "get:ShowAddArticle;post:HandleAddArticle")
	//添加文章类型
	beego.Router("/article/AddArticleType", &controllers.ArticleController{}, "get:ShowAddArticleType;post:HandleAddType")
	//删除文章类型
	beego.Router("/article/DeleteArticleType", &controllers.ArticleController{}, "get:DeleteArticleType")
}


///过滤函数
var FilterMethod = func(ctx *context.Context) {
	userName := ctx.Input.Session("UserName")
	if userName == nil {
		ctx.Redirect(302, "/login")
		return
	}
}
