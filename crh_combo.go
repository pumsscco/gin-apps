package main

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "time"
    "encoding/json"
)
type Combo struct {
	Id  int		`json:"-"`
	Role string 		`json:"role"`
	Name string 		`json:"name"`
	Fighter1Name string 		`json:"fighter1_name"`
	Fighter1Lvl int 		`json:"fighter1_level"`
	Fighter2Name string		`json:"fighter2_name"`
	Fighter2Lvl int 		`json:"fighter2_level"`
	Fighter1 int		`json:"-"`
	Fighter2 int		`json:"-"`
	RoleFlag int		`json:"-"`	
}
func combo(c *gin.Context) {
	var  combos []Combo
    val, err := client.Get("crh:combo").Result()
    if err == nil {
        json.Unmarshal([]byte(val), &combos)
        c.IndentedJSON(http.StatusOK,combos)
        return
	}
	sql := "select id,f1id,f1lvl,f2id,f2lvl,cast(user_flag as unsigned) from SPComboMartial where user_flag!=0"
    rows, _ := Db.Query(sql)
    for rows.Next() {
        cb := Combo{}
        rows.Scan(
            &cb.Id, &cb.Fighter1, &cb.Fighter1Lvl, &cb.Fighter2, &cb.Fighter2Lvl, &cb.RoleFlag,
        )
        cb.Name = getName("SPFighterTable2", cb.Id)
        cb.Fighter1Name = getName("SPFighterTable2", cb.Fighter1)
        if cb.Fighter1Name == "" {
            cb.Fighter1Name = getName("SPFighterTable1", cb.Fighter1)
        }
        cb.Fighter2Name = getName("SPFighterTable2", cb.Fighter2)
        if cb.Fighter2Name == "" {
            cb.Fighter2Name = getName("SPFighterTable1", cb.Fighter2)
        }
        cb.Role = getValidRole(cb.RoleFlag)
        combos = append(combos, cb)
	}
	rows.Close()
	s, err := json.Marshal(combos)
    if err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    } else {
        client.Set("crh:combo", string(s), 36*time.Hour)
        c.IndentedJSON(http.StatusOK,combos)
    }
}
