package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
)

//生成页面
func generateHTML(w http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}
	templates := template.Must(template.ParseFiles(files...))
	templates.ExecuteTemplate(w, "layout", data)
}

//首页
func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	generateHTML(w, nil, "layout", "navbar", "index")
}

//组合技
func comboSkill(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	combos := getCombo()
	generateHTML(w, &combos, "layout", "navbar", "combo")
}

//武功
func fighterPractice(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	routeMapfighter := map[string]string{
		"neigong": "内功心法", "lover": "情侣合技", "rage": "怒技", "common": "普通招式", "combo": "组合技",
	}
	fighterType := routeMapfighter[ps.ByName("Type")]
	table := ""
	if r.Method == "GET" {
		fighters := getFighterType(fighterType, ps.ByName("Type"))
		generateHTML(w, &fighters, "layout", "navbar", "fighter-form")
	} else if r.Method == "POST" {
		fighter := r.PostFormValue("fighter")
		if fighterType == "内功心法" {
			table = "SPFighterTable1"
		} else {
			table = "SPFighterTable2"
		}
		fighterPracs := getFighterLevel(table, fighter)
		generateHTML(w, &fighterPracs, "layout", "navbar", "fighter-detail")
	}
}

//道具
func itemType(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	routeMapItem := map[string]string{
		"fan": "扇", "sword": "剑", "dagger": "短剑", "bow": "弓",
		"armor": "盔甲", "boots": "鞋", "ornament": "佩饰",
		"kungfu": "武功",
		"elixir": "丹药", "hidden-weapon": "暗器", "food": "食物",
		"ingredient": "食材",
	}
	items := getItemType(routeMapItem[ps.ByName("Type")])
	generateHTML(w, &items, "layout", "navbar", "item")
}

//敌人明细
func enemyType(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	routeMapEnemy := map[string]string{
		"male": "男性", "female": "女性", "other": "其它",
	}
	enemys := getEnemyType(routeMapEnemy[ps.ByName("Type")])
	generateHTML(w, &enemys, "layout", "navbar", "enemy")
}
//任务
func mission(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	missions := getMission()
	generateHTML(w, &missions, "layout", "navbar", "mission")
}
//角色初始信息
func role(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	roles := getRoleMan()
	generateHTML(w, &roles, "layout", "navbar", "role")
}
/*
//收集详情
func findItem(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    cat:=ps.ByName("cat")
    t:=ps.ByName("Type")
    if r.Method=="GET" {
        things:=getThings(cat,t)
        generateHTML(w, &things, "layout", "navbar", "find/type")
    } else if r.Method=="POST" {
        thing:=r.PostFormValue("thing")
        methods:=findMethod(cat,t,thing)
        generateHTML(w, &methods, "layout", "navbar", "find/method")
    }
}*/
