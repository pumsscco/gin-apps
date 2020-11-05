package main

import (
	"fmt"
	"reflect"
	"strings"
	"time"
	"encoding/json"
)

//两张物品表的合并结构，加上少量额外的字段，用于存储解析出来的数据
type EnemyInfo struct {
	Id, Level, Steal1, Steal1Num, Steal1Chance              int
	Name, ModelType                                         string //名称、模型类别
	Steal2, Steal2Num, Steal2Chance, DropItem, DropChance   int
	StolenSum, DropInfo, Fighters                           string
	FightType, Fighter1, Fighter1Lvl, Fighter2, Fighter2Lvl int
	Fighter3, Fighter3Lvl, Fighter4, Fighter4Lvl            int
}
type EnemyType struct {
	Type     string
	InfoList []EnemyInfo
}

//将敌人分为三组显示
func getEnemyType(t string) (enemys EnemyType) {
	//依据敌人类型，设置不同的缓存键名
    k:=fmt.Sprintf("et:%s",t)
    val,err:=client.Get(k).Result()
    if err==nil {
        json.Unmarshal([]byte(val),&enemys)
        return
    }
	enemys.Type = t
	sql := `
		select k3.id,k3.name,k1.name,level,steal1,steal1num,steal1chance,steal2,steal2num,steal2chance,
		drop_item,drop_chance,fight_type,f1,f1lvl,f2,f2lvl,f3,f3lvl,f4,f4lvl 
		from SPKindle3 k3 join SPKindle1 k1 on k3.k1id=k1.id where
	`
	switch t {
	case "男性":
		sql += " race=0"
	case "女性":
		sql += " race=5"
	case "其它":
		sql += " race!=0 and race!=5"
	}
	rows, _ := Db.Query(sql)
	for rows.Next() {
		e := EnemyInfo{}
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
		enemys.InfoList = append(enemys.InfoList, e)
	}
	rows.Close()
	as,err:=json.Marshal(enemys)
	client.Set(k, string(as), 12*time.Hour)
	if err!=nil {
		logger.Print(err)
	}
	return
}
