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
type Enemy struct {
	Id int 	   	`json:"id"`
	Name string   `json:"name"`
	ModelType  string   `json:"model_type"`
	Level int  		`json:"level,omitempty"`
	StolenSum string  `json:"stolen_sum,omitempty"`
	DropInfo string   `json:"drop_info,omitempty"`
	FightType int   `json:"fight_type"`
	Fighters string    `json:"fighters,omitempty"`
	Steal1 int  	`json:"-"`
	Steal1Num int   `json:"-"`
	Steal1Chance    int  `json:"-"`
	Steal2 int   `json:"-"`
	Steal2Num int   `json:"-"`
	Steal2Chance int   `json:"-"`
	DropItem int   `json:"-"`
	DropChance   int   `json:"-"`
	Fighter1 int    `json:"-"`
	Fighter1Lvl int    `json:"-"`
	Fighter2 int   `json:"-"`
	Fighter2Lvl int    `json:"-"`
	Fighter3 int      `json:"-"`
	Fighter3Lvl int     `json:"-"`
	Fighter4 int      `json:"-"`
	Fighter4Lvl int    `json:"-"`
}

func cEnemy(c *gin.Context) {
	var (
		et OneParam
		enemys []Enemy
		err error
	)
	if err = c.ShouldBindJSON(&et); err != nil {    
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    k:=fmt.Sprintf("crh:enemy:%s",et.Type)
    val,err:=client.Get(k).Result()
    if err==nil {
		json.Unmarshal([]byte(val),&enemys)
		c.IndentedJSON(http.StatusOK,enemys)
        return
    }
	sql := `
		select k3.id,k3.name,k1.name,level,steal1,steal1num,steal1chance,steal2,steal2num,steal2chance,
		drop_item,drop_chance,fight_type,f1,f1lvl,f2,f2lvl,f3,f3lvl,f4,f4lvl 
		from SPKindle3 k3 join SPKindle1 k1 on k3.k1id=k1.id where
	`
	switch et.Type {
	case "male":
		sql += " race=0"
	case "female":
		sql += " race=5"
	case "other":
		sql += " race!=0 and race!=5"
	}
	rows, _ := Db.Query(sql)
	for rows.Next() {
		e := Enemy{}
		rows.Scan(
			&e.Id, &e.Name, &e.ModelType, &e.Level, &e.Steal1, &e.Steal1Num, &e.Steal1Chance,
			&e.Steal2, &e.Steal2Num, &e.Steal2Chance, &e.DropItem, &e.DropChance,
			&e.FightType, &e.Fighter1, &e.Fighter1Lvl, &e.Fighter2, &e.Fighter2Lvl,
			&e.Fighter3, &e.Fighter3Lvl, &e.Fighter4, &e.Fighter4Lvl,
		)
		//分别合并偷窃、掉落、武功三类信息
		v := reflect.ValueOf(e)
		for i := 1; i <= 2; i++ {
			if fv := v.FieldByName(fmt.Sprintf("Steal%d", i)); fv.Int() > 0 {
				s := getName("SPItemData", int(fv.Int()))
				n := v.FieldByName(fmt.Sprintf("Steal%dNum", i)).Int()
				c := v.FieldByName(fmt.Sprintf("Steal%dChance", i)).Int()
				e.StolenSum += fmt.Sprintf("%s:%d:%d%% ", s, n, c)
			}
		}
		if e.DropItem > 0 {
			d := getName("SPItemData", e.DropItem)
			e.DropInfo = fmt.Sprintf("%s:%d%%", d, e.DropChance)
		}
		for i := 1; i <= 4; i++ {
			if fv := v.FieldByName(fmt.Sprintf("Fighter%d", i)); fv.Int() > 0 {
				f := ""
				if fv.Int() < 100 {
					f = getName("SPFighterTable1", int(fv.Int()))
				} else {
					f = getName("SPFighterTable2", int(fv.Int()))
				}
				l := v.FieldByName(fmt.Sprintf("Fighter%dLvl", i)).Int()
				e.Fighters += fmt.Sprintf("%s:%d ", f, l)
			}
		}
		e.StolenSum = strings.TrimSuffix(e.StolenSum, " ")
		e.Fighters = strings.TrimSuffix(e.Fighters, " ")
		enemys = append(enemys, e)
	}
	rows.Close()
	s,err:=json.Marshal(enemys)
	client.Set(k, string(s), 12*time.Hour)
	if err!=nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		client.Set(k, string(s), 36*time.Hour)
		c.IndentedJSON(http.StatusOK,enemys)
	}
}
