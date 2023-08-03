package controller

import (
	"goblog/common"
	"goblog/config"
	"goblog/dao"
	"goblog/models"
	"goblog/service"
	"net/http"
	"strconv"
	"strings"
)

func Index(w http.ResponseWriter, r *http.Request) {
	// 获取 indexTemplate，即模板引擎中的 "index" 模板
	indexTemplate := common.Template.Index

	// 解析表单数据: 将请求的Body部分（如果有）中的数据解析成表单键值对
	if err := r.ParseForm(); err != nil {
		indexTemplate.WriteError(w, err)
		return
	}

	// 获取请求路径
	path := r.URL.Path

	// 获取表单中的 "page" 参数
	page := r.Form.Get("page")
	if page == "" {
		page = "1"
	}

	// 从数据库获取所有博客分类
	categorys := dao.GetCategorys()

	// 将 "page" 参数转换为整数类型
	currentPage, _ := strconv.Atoi(page)

	// 去掉路径中的前缀 "/"
	slug := strings.TrimPrefix(path, "/")

	var post []models.PostMore
	var total int

	if slug != "" {
		// 请求的是自定义的路径，即文章的 slug
		// 查询文章时按照 slug 进行查询
		post, total = service.PostPageBySlug(currentPage, 10, slug)
	} else {
		// 请求的是默认路径
		// 查询文章时按照页码进行查询
		post, total = service.PostPage(currentPage, 10)
	}

	// 计算总页数
	pagesAll := ((total - 1) / 10) + 1
	pages := []int{}
	for i := 1; i <= pagesAll; i++ {
		pages = append(pages, i)
	}

	// 构建首页数据
	hd := models.HomeData{
		config.Cfg.Viewer,
		categorys,
		post,
		total,
		currentPage,
		pages,
		currentPage != pagesAll,
	}

	// 渲染模板并输出到 ResponseWriter
	indexTemplate.WriteData(w, hd)
}
