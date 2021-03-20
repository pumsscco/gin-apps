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
type RoleMan struct {
	Id    int   `json:"id"`
	Name   string  `json:"name"`
	Model int  `json:"model"`
	Own string  `json:"own,omitempty"`
	Fighters  string  `json:"fighters,omitempty"`
	OwnItem1 int  `json:"-"`
	OwnCount1 int  `json:"-"`
	OwnItem2 int `json:"-"`
	OwnCount2 int `json:"-"`
	Fighter1 int  `json:"-"`
	Fighter1Lvl int  `json:"-"`
	Fighter2 int  `json:"-"`
	Fighter2Lvl int  `json:"-"`
	Fighter3 int  `json:"-"`
	Fighter3Lvl int  `json:"-"`
	
}
func role(c *gin.Context) {
	var roles []RoleMan
	k:="crh:role"
	val,err:=client.Get(k).Result()
    if err==nil {
		json.Unmarshal([]byte(val),&roles)
		c.IndentedJSON(http.StatusOK,roles)
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
	s,err:=json.Marshal(roles)
	if err!=nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
	} else {
		client.Set(k, string(s), 36*time.Hour)
        c.IndentedJSON(http.StatusOK,roles)
	}
}