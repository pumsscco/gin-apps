package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
	"time"
)
//生成页面
func generateHTML(w http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}
	funcMap := template.FuncMap{ 
		"fdate": func(t time.Time) string { return t.Format("2006-01-02") },
	}
	t,_:=template.New("list.html").Funcs(funcMap).ParseFiles(files...)
	t.ExecuteTemplate(w, "layout", data)
}

//首页
func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	generateHTML(w, nil, "layout", "navbar", "index")
}

//交易记录
func dealList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		if r.Method=="GET" {
			names:=getNameMap("")
			generateHTML(w, &names, "layout", "navbar", "name-list")
		} else if r.Method=="POST" {
			deals := getDealList(r.PostFormValue("code"))
			generateHTML(w, &deals, "layout", "navbar", "deal-list")
		}
}
//持仓股票最新买卖记录
func holdLastDeal(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sums := getHoldLastDeal()
	generateHTML(w, &sums, "layout", "navbar", "deal-list")
}
//清仓记录
func clearance(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	clears := getClearance()
	generateHTML(w, &clears, "layout", "navbar", "clearance")
}
//持仓统计
func position(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	positions := getPosition()
	generateHTML(w, &positions, "layout", "navbar", "position")
}
//新增交易记录
func newDeal(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method=="GET" {
		generateHTML(w, nil, "layout", "navbar", "add")
	} else if r.Method=="POST" {
		createDeal(w,r)
	}
}
/*
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

//组合技
func comboSkill(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	combos := getCombo()
	generateHTML(w, &combos, "layout", "navbar", "combo")
}*/
