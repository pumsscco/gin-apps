package main

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "time"
    "encoding/json"
    _ "github.com/go-sql-driver/mysql"
)
type Combo struct {
	Id                                                     int
	Fighter1, Fighter1Lvl, Fighter2, Fighter2Lvl, RoleFlag int
	Role, Name, Fighter1Name, Fighter2Name                 string
}
func combo(c *gin.Context) {
	var  combos []Combo
	var c_infos []gin.H
	val, err := client.Get("crh:combo").Result()
	if err == nil {
		json.Unmarshal([]byte(val), &c_infos)
		c.IndentedJSON(http.StatusOK,c_infos)
		return
	}
	sql := "select id,f1id,f1lvl,f2id,f2lvl,cast(user_flag as unsigned) from SPComboMartial where user_flag!=0"
	rows, _ := Db.Query(sql)
	for rows.Next() {
		combo := Combo{}
		rows.Scan(
			&combo.Id, &combo.Fighter1, &combo.Fighter1Lvl, &combo.Fighter2, &combo.Fighter2Lvl, &combo.RoleFlag,
		)
		combo.Name = getName("SPFighterTable2", combo.Id)
		combo.Fighter1Name = getName("SPFighterTable2", combo.Fighter1)
		if combo.Fighter1Name == "" {
			combo.Fighter1Name = getName("SPFighterTable1", combo.Fighter1)
		}
		combo.Fighter2Name = getName("SPFighterTable2", combo.Fighter2)
		if combo.Fighter2Name == "" {
			combo.Fighter2Name = getName("SPFighterTable1", combo.Fighter2)
		}
		combo.Role = getValidRole(combo.RoleFlag)
		combos = append(combos, combo)
	}
	rows.Close()
	for _, v:= range combos {
		c_info:=gin.H{
			"name": v.Name,
			"fighter1_name": v.Fighter1Name,
			"fighter1_level": v.Fighter1Lvl,
			"fighter2_name": v.Fighter2Name,
			"fighter2_level": v.Fighter2Lvl,
			"role": v.Role,
		}
		c_infos=append(c_infos,c_info)
	}
	s, err := json.Marshal(c_infos)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		client.Set("crh:combo", string(s), 24*time.Hour)
		c.IndentedJSON(http.StatusOK,c_infos)
	}
}
