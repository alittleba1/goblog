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
	indexTemplate := common.Template.Index

	// 解析表单数据: 将请求的Body部分（如果有）中的数据解析成表单键值对
	if err := r.ParseForm(); err != nil {
		indexTemplate.WriteError(w, err)
		return
	}

	// 获取请求路径
	path := r.URL.Path

	// 获取page参数
	page := r.Form.Get("page")
	if page == "" {
		page = "1"
	}

	// 转换为整数类型
	currentPage, _ := strconv.Atoi(page)

	var post []models.PostMore
	var total int

	slug := strings.TrimPrefix(path, "/")
	if slug != "" {
		// 按照 slug 查询文章
		post, total = service.PostPageBySlug(currentPage, 10, slug)
	} else {
		// 按照页码查询文章(默认路径)
		post, total = service.PostPage(currentPage, 10)
	}

	// 计算总页数
	pagesAll := ((total - 1) / 10) + 1
	// 创建空整数切片
	pages := []int{}
	for i := 1; i <= pagesAll; i++ {
		pages = append(pages, i)
	}

	// 获取分类数据
	categorys := dao.GetCategorys()

	// 构建首页数据
	hd := models.HomeData{
		Viewer:    config.Cfg.Viewer,
		Categorys: categorys,
		Posts:     post,
		Total:     total,
		Page:      currentPage,
		Pages:     pages,
		PageEnd:   currentPage != pagesAll,
	}

	// 渲染模板并输出到 ResponseWriter
	indexTemplate.WriteData(w, hd)
}
