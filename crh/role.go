package main

import (
	"fmt"
	"reflect"
	"strings"
	"time"
	"encoding/json"
)
type RoleMan struct {
	Id                             int
	Name                                         string //名称、模型类别
	Model,OwnItem1,OwnCount1,OwnItem2,OwnCount2 int
	Fighter1,Fighter1Lvl,Fighter2,Fighter2Lvl,Fighter3, Fighter3Lvl int
	Own,Fighters  string
}
func getRoleMan() (roles []RoleMan) {
	val,err:=client.Get("role").Result()
    if err==nil {
        json.Unmarshal([]byte(val),&roles)
        return
    }
	sql := `
		select id,name,model,own_item1,own_count1,own_item2,own_count2,f1,f1lvl,f2,f2lvl,f3,f3lvl from SPRoleManDefault
	`
	rows, _ := Db.Query(sql)
	for rows.Next() {
		r := RoleMan{}
		rows.Scan(
			&r.Id, &r.Name, &r.Model, &r.OwnItem1, &r.OwnCount1, &r.OwnItem2, &r.OwnCount2,
			&r.Fighter1, &r.Fighter1Lvl, &r.Fighter2, &r.Fighter2Lvl, &r.Fighter3,&r.Fighter3Lvl,
		)
		v := reflect.ValueOf(r)
		for i := 1; i <= 2; i++ {
			if fv := v.FieldByName(fmt.Sprintf("OwnItem%d", i)); fv.Int() > 0 {
				o := getName("SPItemData", int(fv.Int()))
				n := v.FieldByName(fmt.Sprintf("OwnCount%d", i)).Int()
				r.Own += fmt.Sprintf("%s:%d ", o, n)
			}
		}
		for i := 1; i <= 3; i++ {
			if fv := v.FieldByName(fmt.Sprintf("Fighter%d", i)); fv.Int() > 0 {
				f := ""
				if fv.Int() < 100 {
					f = getName("SPFighterTable1", int(fv.Int()))
				} else {
					f = getName("SPFighterTable2", int(fv.Int()))
				}
				l := v.FieldByName(fmt.Sprintf("Fighter%dLvl", i)).Int()
				r.Fighters += fmt.Sprintf("%s:%d ", f, l)
			}
		}
		r.Own = strings.TrimSuffix(r.Own, " ")
		r.Fighters = strings.TrimSuffix(r.Fighters, " ")
		roles = append(roles, r)
	}
	rows.Close()
	as,err:=json.Marshal(roles)
	client.Set("role", string(as), 12*time.Hour)
	if err!=nil {
		logger.Print(err)
	}
	return
}
