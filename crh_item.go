package main

import (
    "encoding/json"
    "fmt"
    "github.com/gin-gonic/gin"
    "net/http"
    "reflect"
    "strings"
    "time"
)

//两张物品表的合并结构，加上少量额外的字段，用于存储解析出来的数据
type CItem struct {
	Id  int  		`json:"id"`
	Name string `json:"name"`
	Description string  `json:"description"`
	Attribute     string  `json:"attribute,omitempty"`
	Strength int  `json:"-"`
	Accuracy int   `json:"-"`
	CritRate     int    `json:"-"`
	Defence int   `json:"-"`
	Agility int `json:"-"`
	Quick          int    `json:"-"`
	Intelligence int 		`json:"-"`
	MP                 int    `json:"-"`
	Price int   		`json:"price,omitempty"`
	RoleFlag int 		`json:"-"`
	Func int 		`json:"function,omitempty"`
	FuncParam int    `json:"function_parameter,omitempty"`
	ValidRoles           string `json:"valid_roles,omitempty"`
}

func item(c *gin.Context) {
	var (
		it OneParam
		items []CItem
		err error
	)
	if err = c.ShouldBindJSON(&it); err != nil {    
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
	k:=fmt.Sprintf("crh:item:%s",it.Type)
    val,err:=client.Get(k).Result()
    if err==nil {
		json.Unmarshal([]byte(val),&items)
		c.IndentedJSON(http.StatusOK,items)
        return
    }
	sql := `
        select id,name,strength,accuracy,crit_rate,defence,agility,quick,intelligence,mp,price,
        cast(user_flag as unsigned),func,func_param from SPItemData where type=? and name not like "%保留%"
    `
	typeMap := map[string]int{
		"fan": 1, "sword": 2, "dagger": 3, "bow": 4,
		"armor": 6, "boots": 7, "ornament": 8,
		"kungfu": 10,
		"elixir": 12, "hidden-weapon": 12, "food": 12,
		"ingredient": 13,
	}
	switch it.Type {
	case "elixir":
		sql += " and id between 201 and 230"
	case "hidden-weapon":
		sql += " and id between 261 and 280"
	case "food":
		sql += " and id between 361 and 380"
	}
	rows, _ := Db.Query(sql, typeMap[it.Type])
	for rows.Next() {
		i := CItem{}
		rows.Scan(
			&i.Id, &i.Name, &i.Strength, &i.Accuracy, &i.CritRate, &i.Defence, &i.Agility, &i.Quick,
			&i.Intelligence, &i.MP, &i.Price, &i.RoleFlag, &i.Func, &i.FuncParam,
		)
		//合并属性，依次为力道、命中、暴击、守备、活泛、迅捷、才识、内力
		//+剔除食材、暗器、丹药
		if it.Type != "ingredient" && it.Type != "hidden-weapon" && it.Type != "elixir" {
			comments := [][]string{
				{"Strength", "力道"}, {"Accuracy", "命中"}, {"CritRate", "暴击"},
				{"Defence", "守备"}, {"Agility", "活泛"}, {"Quick", "迅捷"},
				{"Intelligence", "才识"}, {"MP", "内力"},
			}
			v := reflect.ValueOf(i)
			for _, f := range comments {
				if fv := v.FieldByName(f[0]).Int(); fv > 0 {
					i.Attribute += fmt.Sprintf("%s+%d ", f[1], fv)
				} else if fv < 0 {
					i.Attribute += fmt.Sprintf("%s%d ", f[1], fv)
				}
			}
			if it.Type != "food" {
				i.ValidRoles = getValidRole(i.RoleFlag)
			}
		}
		if it.Type == "kungfu" {
			i.Attribute = strings.Replace(i.Attribute, "+", ":", -1)
		}
		if it.Type == "food" {
			i.Attribute = strings.Replace(strings.Replace(strings.Replace(i.Attribute, "命中-1 ", "", -1), "力道", "生命", -1), "暴击", "力道", -1)
		}
		i.Description=getName("SPItemHelp",i.Id)
		i.Attribute = strings.TrimSuffix(i.Attribute, " ")
		items = append(items, i)
	}
	rows.Close()
	s,err:=json.Marshal(items)
	if err!=nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		client.Set(k, string(s), 36*time.Hour)
		c.IndentedJSON(http.StatusOK,items)
	}
}