package main

import (
	"fmt"
	"reflect"
	"strings"
	"time"
	"encoding/json"
)

//两张物品表的合并结构，加上少量额外的字段，用于存储解析出来的数据
type Item struct {
	Id                               int
	Name, Description, Attribute     string //名称、描述、属性综合
	Strength, Accuracy, CritRate     int    //力道、命中、暴击
	Defence, Agility, Quick          int    //守备、活泛、迅捷
	Intelligence, MP                 int    //才识、内力、要求等级
	Price, RoleFlag, Func, FuncParam int    //价格、角色标志位、调用函数
	ValidRoles                       string //适用角色
}
type Items struct {
	Type string
	IsEquip, HasEff, IsKungfu, OnlyPrice,IsFood bool
	ItemList                             []Item
}

//利用物品类型（中文），获得该类全部的原始数据，属性内容合并
func getItemType(t string) (items Items) {
	k:=fmt.Sprintf("it:%s",t)
    val,err:=client.Get(k).Result()
    if err==nil {
        json.Unmarshal([]byte(val),&items)
        return
    }
	items.Type=t
	sql := `
        select id,name,strength,accuracy,crit_rate,defence,agility,quick,intelligence,mp,price,
        cast(user_flag as unsigned),func,func_param from SPItemData where type=? and name not like "%保留%"
    `
	typeMap := map[string]int{
		"扇": 1, "剑": 2, "短剑": 3, "弓": 4,
		"盔甲": 6, "鞋": 7, "佩饰": 8,
		"武功": 10,
		"丹药": 12, "暗器": 12, "食物": 12,
		"食材": 13,
	}
	switch t {
	case "扇","弓","盔甲","鞋":
		items.IsEquip=true
	case "剑","短剑","佩饰":
		items.IsEquip=true
		items.HasEff=true
	case "武功":
		items.IsKungfu=true
	case "丹药":
		sql += " and id between 201 and 230"
		items.OnlyPrice=true
	case "暗器":
		sql += " and id between 261 and 280"
	case "食物":
		sql += " and id between 361 and 380"
		items.IsFood=true
	case "食材":
		items.OnlyPrice=true
	}
	rows, _ := Db.Query(sql, typeMap[t])
	for rows.Next() {
		//取出原始数据
		i := Item{}
		rows.Scan(
			&i.Id, &i.Name, &i.Strength, &i.Accuracy, &i.CritRate, &i.Defence, &i.Agility, &i.Quick,
			&i.Intelligence, &i.MP, &i.Price, &i.RoleFlag, &i.Func, &i.FuncParam,
		)
		//合并属性，依次为力道、命中、暴击、守备、活泛、迅捷、才识、内力
		//+剔除食材、暗器、丹药
		if t != "食材" && t != "暗器" && t != "丹药" {
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
			if t != "食物" {
				i.ValidRoles = getValidRole(i.RoleFlag)
			}
		}
		if t == "武功" {
			i.Attribute = strings.Replace(i.Attribute, "+", ":", -1)
		}
		if t == "食物" {
			i.Attribute = strings.Replace(strings.Replace(strings.Replace(i.Attribute, "命中-1 ", "", -1), "力道", "生命", -1), "暴击", "力道", -1)
		}
		i.Description=getName("SPItemHelp",i.Id)
		i.Attribute = strings.TrimSuffix(i.Attribute, " ")
		items.ItemList = append(items.ItemList, i)
	}
	rows.Close()
	as,err:=json.Marshal(items)
	client.Set(k, string(as), 12*time.Hour)
	if err!=nil {
		logger.Print(err)
	}
	return
}
