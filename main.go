package main

import (
	_ "article-manage-system/routers"
	"github.com/astaxie/beego"
)

func main() {
	//绑定视图函数--上一页
	beego.AddFuncMap("ShowPrePage",HandlePrePage)
	//绑定视图函数--下一页
	beego.AddFuncMap("ShowNextPage",HandleNextPage)

	beego.Run()
}

//处理上一页
func HandlePrePage(index int)int{
	return  index-1
}

//处理下一页
func HandleNextPage(index int)int{
	return  index+1
}