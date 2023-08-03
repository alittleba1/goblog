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

type HTMLApi struct{}

type Api struct{}

var HTML = &HTMLApi{}

var API = &Api{}

func (*HTMLApi) Category(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	id := strings.TrimPrefix(path, "/c/")
	cId, _ := strconv.Atoi(id)

	_ = r.ParseForm()
	page := r.Form.Get("page")
	if page == "" {
		page = "1"
	}
	currentPage, _ := strconv.Atoi(page)

	// 获取分类的博客数据
	post, total := service.PostPageByCategory(currentPage, 10, cId)
	pagesAll := ((total - 1) / 10) + 1
	pages := []int{}
	for i := 1; i <= pagesAll; i++ {
		pages = append(pages, i)
	}

	// 获取类名
	cName := dao.GetCategoryNameById(cId)
	// 获取分类数据
	categorys := dao.GetCategorys()

	hd := models.HomeData{
		Viewer:    config.Cfg.Viewer,
		Categorys: categorys,
		Posts:     post,
		Total:     total,
		Page:      currentPage,
		Pages:     pages,
		PageEnd:   currentPage != pagesAll,
	}
	var categoryData = &models.CategoryData{
		HomeData:     hd,
		CategoryName: cName,
	}

	common.Template.Category.WriteData(w, categoryData)
}
