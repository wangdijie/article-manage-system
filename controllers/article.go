package controllers

import (
	"article-manage-system/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"math"
	"path"
	"time"
)

type ArticleController struct {
	beego.Controller
}

//显示文章列表
func (this *ArticleController) ShowArticleList() {
	om := orm.NewOrm()
	qs := om.QueryTable("Article")

	typeId, t_err := this.GetInt("select") //文章分类
	beego.Info(typeId)

	pageIndex, err := this.GetInt("pageIndex") //页索引
	if err != nil {
		pageIndex = 1
	}
	var totalCount int64                //总记录数
	pageSize := 5                       //每页数量
	start := (pageIndex - 1) * pageSize //起始位置计算

	if t_err != nil { //每页分类，查出所有的数量
		totalCount, _ = qs.Count()
	} else { //有分类，根据分类查
		totalCount, _ = qs.RelatedSel("ArticleType").Filter("ArticleType__Id", typeId).Count()
	}
	pageCount := math.Ceil(float64(totalCount) / float64(pageSize)) //总页数

	var articles []models.Article
	var types []models.ArticleType //获取文章的数据类型
	om.QueryTable("ArticleType").All(&types)

	//根据选中的类型查询相应类型文章
	if t_err != nil {
		qs.Limit(pageSize, start).RelatedSel("ArticleType").All(&articles)
	} else {
		qs.Limit(pageSize, start).RelatedSel("ArticleType").Filter("ArticleType__Id", typeId).All(&articles)
	}
	isFirstPage := false
	isLastPage := false
	if pageIndex == 1 {
		isFirstPage = true
	}
	if pageIndex == int(pageCount) {
		isLastPage = true
	}

	//文章的类型数据
	this.Data["ArticleTypes"] = types

	if typeId > 0 {
		this.Data["CurrentTypeId"] = typeId
	}

	this.Data["PageIndex"] = pageIndex
	this.Data["TotalCount"] = totalCount
	this.Data["PageCount"] = pageCount

	this.Data["IsFirstPage"] = isFirstPage //是否第一页
	this.Data["IsLastPage"] = isLastPage   //是否最后一页

	//文章的数据
	this.Data["Articles"] = articles

	this.Data["UserName"] = this.GetSession("UserName")
	this.Data["WebTitle"] = "文章列表"
	this.Layout = "layout.html"
	this.TplName = "index.html"
}

//添加文章
func (this *ArticleController) ShowAddArticle() {
	om := orm.NewOrm()

	//获取文章的数据类型
	var types []models.ArticleType
	om.QueryTable("ArticleType").All(&types)

	this.Data["ArticleTypes"] = types
	this.Data["WebTitle"] = "添加文章"

	//指定视图
	this.Layout = "layout.html"
	this.TplName = "add.html"
}

//处理添加文章
func (this *ArticleController) HandleAddArticle() {
	title := this.GetString("articleName")
	content := this.GetString("content")
	typeId, _ := this.GetInt("select")
	file, header, err := this.GetFile("uploadname")
	//defer file.Close()
	imgPath:=""

	if file!=nil{
		//判断图片的格式
		ext := path.Ext(header.Filename)
		if ext != ".jpg" && ext != ".png" && ext != "jpeg" {
			beego.Info("上传文件格式不正确")
			return
		}
		//判断文件大小
		if header.Size > 5000000 {
			beego.Info("上传太大，不允许上传")
			return
		}
		//文件名定义，不能重名
		filename := time.Now().Format("20160102150405") + ext
		imgPath = "/static/img/" + filename
		this.SaveToFile("uploadname", "."+imgPath)
		if err != nil {
			beego.Info("文件上传失败")
			return
		}
	}

	om := orm.NewOrm()
	article := models.Article{Title: title, Content: content, Img: imgPath, CreateTime: time.Now()}

	//获取文章类型对象
	articleType := models.ArticleType{Id: typeId}
	err = om.Read(&articleType)
	if err != nil {
		beego.Info("获取类型错误")
		return
	}
	//赋值给文章的类型  指针，赋地址
	article.ArticleType = &articleType

	n, _ := om.Insert(&article)
	if n > 0 {
		//插入成功 跳转
		this.Redirect("/article/ShowArticles", 302)
	} else {
		//指定视图布局
		this.Layout = "layout.html"
		this.TplName = "add.html"
	}
}

//删除文章
func (this *ArticleController) DeleteArticle(){
	id,_:=this.GetInt("id")
	om := orm.NewOrm()
	article := models.Article{Id: id}

	om.Delete(&article)

	this.Redirect("/article/ShowArticles", 302)
}

//文章详情
func (this *ArticleController) ShowArticleDetail() {
	id, err := this.GetInt("id")
	if err != nil {
		beego.Info("参数错误")
		return
	}

	om := orm.NewOrm()
	article := models.Article{Id: id}

	//om.Read(&article)

	//查询文章详情，并关联类型
	om.QueryTable("Article").RelatedSel("ArticleType").Filter("Id", id).One(&article)

	//更新阅读数
	article.ViewCount += 1
	om.Update(&article)

	//多对多插入浏览记录
	//第二个参数是article中的中的属性Users
	m2m := om.QueryM2M(&article, "Users")

	//当前登陆用户
	userName := this.GetSession("UserName")
	if userName == nil {
		this.Redirect("/login", 302)
		return
	}
	//当前用户
	user := models.User{UserName: userName.(string)}
	//查询用户信息
	om.Read(&user, "UserName")

	//插入操作（多对多），此时user_articles表中将有一条数据
	//当前文章插入一条用户浏览记录
	m2m.Add(user)

	//多对多查询操作，查询此篇文章所有的浏览用户
	var users []models.User
	om.QueryTable("User").Filter("Articles__Article__Id", id).Distinct().All(&users)

	this.Data["Readers"] = users
	this.Data["Article"] = article
	this.Data["WebTitle"] = article.Title

	//指定视图布局
	this.Layout = "layout.html"
	this.TplName = "content.html"
}

//显示更新文章
func (this *ArticleController) ShowUpdateArticle() {
	id, err := this.GetInt("id")
	if err != nil {
		beego.Info("参数错误")
		return
	}
	article := models.Article{Id: id}
	om := orm.NewOrm()
	err = om.Read(&article)
	if err != nil {
		beego.Info("文章不存在")
		this.Redirect("/article/ShowArticles", 302)
	}

	//获取文章的类型数据
	var types []models.ArticleType
	om.QueryTable("ArticleType").All(&types)

	this.Data["ArticleType"] = types
	this.Data["Article"] = article

	this.Layout = "layout.html"
	this.TplName = "update.html"
}

//处理更新文章
func (this *ArticleController) HandleUpdateArticle() {
	title := this.GetString("articleName")
	content := this.GetString("content")
	id, _ := this.GetInt("Id")
	filePath := UploadFile(&this.Controller, "uploadname")
	//数据校验
	if content == "" || title == "" || filePath == "" {
		beego.Info("请求错误")
		return
	}

	om := orm.NewOrm()
	article := models.Article{Id: id}
	error := om.Read(&article)
	if error != nil {
		beego.Info("要更新的数据不存在")
		return
	}

	article.Title = title
	article.Content = content
	if filePath != "NoImg" {
		article.Img = filePath
	}
	om.Update(&article)

	this.Redirect("/article/ShowArticles", 302)
}

//上传图片封装
func UploadFile(this *beego.Controller, filePath string) string {
	file, header, err := this.GetFile(filePath)
	defer file.Close()

	if header.Filename == "" {
		return "NoImg"
	}
	if err != nil {
		this.Data["errmsg"] = "文件上传失败"
		return ""
	}
	//判断文件格式
	ext := path.Ext(header.Filename)
	if ext != ".jpg" && ext != ".png" && ext != "jpeg" {
		this.Data["errmsg"] = "上传文件格式不正确"
		return ""
	}
	//判断文件大小
	if header.Size > 5000000 {
		this.Data["errmsg"] = "上传太大，不允许上传"
		return ""
	}

	//文件名定义，不能重名
	filename := time.Now().Format("20060102150405") + ext
	path := "/static/img/" + filename
	//存储
	this.SaveToFile("uploadname", "."+path)
	return path
}

//添加类型
func (this *ArticleController) ShowAddArticleType() {
	om := orm.NewOrm()
	var types []models.ArticleType
	om.QueryTable("ArticleType").All(&types)
	this.Data["Types"] = types
	this.Data["WebTitle"] = "添加文章类型"

	//指定视图布局
	this.Layout = "layout.html"
	this.TplName = "addType.html"
}

//处理添加类型
func (this *ArticleController) HandleAddType() {
	typeName := this.GetString("typeName")
	if typeName == "" {
		beego.Info("名称不能为空")
		return
	}
	om := orm.NewOrm()
	articleType := models.ArticleType{TypeName: typeName}
	_, err := om.Insert(&articleType)
	if err != nil {
		beego.Info("添加失败")
		this.TplName = "addType.html"
	}
	this.Redirect("/article/AddArticleType", 302)
}

//删除文章类型
func (this *ArticleController) DeleteArticleType(){
	id,_:=this.GetInt("id")

	om:=orm.NewOrm()
	articleCount,_:=om.QueryTable("Article").RelatedSel("ArticleType").Filter("ArticleType__Id",id).Count()

	if articleCount>0{
		this.Data["errmsg"]="类型有"+string(articleCount)+"篇文章，不能删除！"
	}else{
		om.Delete(&models.ArticleType{Id:id})
	}
	this.Redirect("/article/AddArticleType",302)
}
