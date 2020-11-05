package main

import (
	"encoding/json"
	"time"
)

type Combo struct {
	Id                                                     int
	Fighter1, Fighter1Lvl, Fighter2, Fighter2Lvl, RoleFlag int
	Role, Name, Fighter1Name, Fighter2Name                 string
}

func getCombo() (combos []Combo) {
	val, err := client.Get("combo").Result()
	if err == nil {
		json.Unmarshal([]byte(val), &combos)
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
	as, err := json.Marshal(combos)
	client.Set("combo", string(as), 12*time.Hour)
	if err != nil {
		logger.Print(err)
	}
	return
}
