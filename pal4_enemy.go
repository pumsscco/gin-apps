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

//monster表分五组enemy结构展示
type EnemyCommon struct {
	Id              int    `json:"id"`
	Name string   `json:"name"`
	Icon string    `json:"icon"`
	Model string   `json:"model"`
	Description     string `json:"description"`
	IsBoss int          `json:"-"`
	Level int             `json:"level"`
	Experience int         `json:"experience"`
	Target int             `json:"-"`
	Range int    `json:"-"`
	Count   int    `json:"-"`
	StatsCount   string `json:"stats_name"`
	Boss     string `json:"boss"`
	TargetName string   `json:"target_name"`
	RangeName     string `json:"range_name"`
}
type EnemyBasic struct {
	Id int    `json:"id"`
	HP int   `json:"hp"`
	Rage int    `json:"rage"`
	MP      int      `json:"mp"`
	Name string    `json:"name"`
	Wuling      string    `json:"wuling"`
	Physical int   `json:"physical"`
	Toughness int   `json:"toughness"`
	Speed int   `json:"speed"`
	Lucky int   `json:"lucky"`
	Will int      `json:"will"`
	Water int   `json:"-"`
	Fire int   `json:"-"`
	Thunder int    `json:"-"`
	Air int    `json:"-"`
	Earth  int     `json:"-"`
}
type EnemyResistance struct {
	Id  int          `json:"id"`
	Name     string     `json:"name"`
	PhysicalExtract float32   `json:"-"`
	WaterExtract float32   `json:"-"`
	FireExtract     float32   `json:"-"`
	ThunderExtract float32    `json:"-"`
	AirExtract float32    `json:"-"`
	EarthExtract     float32    `json:"-"`
	PhysicalExtractPer string    `json:"physical_extract_per"`
	WaterExtractPer string   `json:"water_extract_per"`
	FireExtractPer string  `json:"fire_pxtract_per"`
	ThunderExtractPer string  `json:"thunder_extract_per"`
	AirExtractPer string  `json:"air_extract_per"`
	EarthExtractPer   string  `json:"earth_extract_per"`
	PhysicalReact float32   `json:"-"`
	WaterReact float32   `json:"-"`
	FireReact    float32 `json:"-"`
	ThunderReact float32  `json:"-"`
	AirReact float32   `json:"-"`
	EarthReact    float32 `json:"-"`
	React          string  `json:"react"`
	SoundWounded1 string   `json:"sound_wounded1"`
	SoundWounded2 string  `json:"sound_wounded2"`
	SoundWounded3         string  `json:"sound_wounded3"`
}
type EnemySkill struct {
	Id     int    `json:"id"`
	Name  string    `json:"name"`
	PhysicalAdditional  int     `json:"physical_additional"`
	AdditionalRate   string    `json:"additional_rate"`
	AdditionalCritical float32   `json:"-"`
	FendOff float32   `json:"-"`
	AdditionalHitting float32  `json:"-"`
	CounterPunchRate float32   `json:"-"`
	Skill1 int         `json:"-"`
	Skill2 int             `json:"-"`
	Skill3 int           `json:"-"`
	Skill4 int         `json:"-"`
	Skill5 int     `json:"-"`
	Skills       string    `json:"skills"`
}
type EnemyDrop struct {
	Id    int   `json:"id"`
	Name       string  `json:"name"`
	StolenProperty int   `json:"-"`
	StolenNumber int   `json:"-"`
	StolenMoney  int     `json:"-"`
	Drop1ID int  `json:"-"`
	Drop2ID int   `json:"-"`
	Drop3ID int   `json:"-"`
	Drop4ID  int     `json:"-"`
	Drop1 string   `json:"drop1"`
	Drop2 string   `json:"drop2"`
	Drop3 string   `json:"drop3"`
	Drop4  string  `json:"drop4"`
	Drop1Rate float32   `json:"-"`
	Drop2Rate float32   `json:"-"`
	Drop3Rate float32   `json:"-"`
	Drop4Rate float32 `json:"-"`
	Drop1Per string   `json:"drop1_per"`
	Drop2Per string   `json:"drop2_per"`
	Drop3Per string  `json:"drop3_per"`
	Drop4Per     string  `json:"drop4_per"`
	MaxDropMoney int   `json:"-"`
	MinDropMoney  int  `json:"-"`
	Stolen  string  `json:"stolen"`
	DropMoney  string  `json:"drop_money"`
}
func pEnemy(c *gin.Context) {
	var (
		ecs []EnemyCommon
		ebs []EnemyBasic
		ers []EnemyResistance
		ess []EnemySkill
		eds []EnemyDrop
        et OneParam
		err error
		enemys interface{}
	)
	if err = c.ShouldBindJSON(&et); err != nil {    
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    k:=fmt.Sprintf("pal4:enemy:%s",et.Type)
    val,err:=client.Get(k).Result()
    if err==nil {
		json.Unmarshal([]byte(val),&enemys)
		c.IndentedJSON(http.StatusOK,enemys)
        return
	}
	switch et.Type {
	case "日常属性":     //common
		sql := `
			select id,name,icon,model,description,is_boss,level,experience,physical_atk_target,physical_atk_range,count from Monster
		`
		rows, _ := Db.Query(sql)
		for rows.Next() {
			ec := EnemyCommon{}
			rows.Scan(
				&ec.Id, &ec.Name, &ec.Icon, &ec.Model,
				&ec.Description, &ec.IsBoss, &ec.Level, &ec.Experience,
				&ec.Target, &ec.Range, &ec.Count,
			)
			ec.Boss = getIsBoss(ec.IsBoss)
			ec.TargetName = getName("PhysicalAttackTarget", ec.Target)
			ec.RangeName = getName("PhysicalAttackType", ec.Range)
			if ec.Count == 1 {
				ec.StatsCount = "是"
			} else {
				ec.StatsCount = "否"
			}
			ecs = append(ecs, ec)
		}
		rows.Close()
		enemys=ecs
	case "八基本属性与五灵":     //basic
		sql := `
			select id,name,max_hp,rage,max_mp,physical,toughness,speed,lucky,will,water,fire,thunder,air,earth from Monster
		`
		rows, _ := Db.Query(sql)
		for rows.Next() {
			eb := EnemyBasic{}
			rows.Scan(
				&eb.Id, &eb.Name, &eb.HP, &eb.Rage, &eb.MP, &eb.Physical,
				&eb.Toughness, &eb.Speed, &eb.Lucky, &eb.Will, &eb.Water,
				&eb.Fire, &eb.Thunder, &eb.Air, &eb.Earth,
			)
			wulingAttribute := [][]string{
				{"Water", "水"}, {"Fire", "火"}, {"Thunder", "雷"},
				{"Air", "风"}, {"Earth", "土"},
			}
			v := reflect.ValueOf(eb)
			for _, f := range wulingAttribute {
				if fv := v.FieldByName(f[0]); fv.Int() > 0 {
					eb.Wuling += fmt.Sprintf("%s:%d ", f[1], fv.Int())
				}
			}
			eb.Wuling = strings.TrimSuffix(eb.Wuling, " ")
			ebs = append(ebs, eb)
		}
		rows.Close()
		enemys=ebs
	case "抗性、反弹、受伤音效":     //resistance
		sql := `
			select id,name,physical_extract,water_extract,fire_extract,thunder_extract,air_extract,earth_extract,
			physical_react,water_react,fire_react,thunder_react,air_react,earth_react,
			sound_wounded1,sound_wounded2,sound_wounded3 from Monster
		`
		rows, _ := Db.Query(sql)
		for rows.Next() {
			er := EnemyResistance{}
			rows.Scan(
				&er.Id, &er.Name, &er.PhysicalExtract, &er.WaterExtract,
				&er.FireExtract, &er.ThunderExtract, &er.AirExtract,
				&er.EarthExtract, &er.PhysicalReact, &er.WaterReact,
				&er.FireReact, &er.ThunderReact, &er.AirReact,
				&er.EarthReact, &er.SoundWounded1, &er.SoundWounded2,
				&er.SoundWounded3,
			)
			reactAttribute := [][]string{
				{"PhysicalReact", "物理"}, {"WaterReact", "水"}, {"FireReact", "火"},
				{"ThunderReact", "雷"}, {"AirReact", "风"}, {"EarthReact", "土"},
			}
			v := reflect.ValueOf(&er).Elem()
			for _, f := range reactAttribute {
				if fv := v.FieldByName(f[0]); fv.Float() > 0 {
					er.React += fmt.Sprintf("%s:%s%% ", f[1], perDisp(float32(fv.Float()*100)))
				}
			}
			extractAttribute := []string{
				"PhysicalExtract", "WaterExtract", "FireExtract",
				"ThunderExtract", "AirExtract", "EarthExtract",
			}
			for _, f := range extractAttribute {
				fv := v.FieldByName(f)
				fvp := v.FieldByName(f + "Per")
				if fv.Float() > 0 {
					fvp.SetString(fmt.Sprintf("%s%%", perDisp(float32(fv.Float()*100))))
				} else {
					fvp.SetString("0")
				}
			}
			er.React = strings.TrimSuffix(er.React, " ")
			ers= append(ers, er)
		}
		rows.Close()
		enemys=ers
	case "物理追加与技能":     //skill
		sql := `
			select id,name,physical_additional,additional_critical,fend_off,additional_hitting,counterpunch_rate,
			skill1,skill2,skill3,skill4,skill5 from Monster
		`
		rows, _ := Db.Query(sql)
		for rows.Next() {
			es := EnemySkill{}
			rows.Scan(
				&es.Id, &es.Name, &es.PhysicalAdditional, &es.AdditionalCritical,
				&es.FendOff, &es.AdditionalHitting, &es.CounterPunchRate, &es.Skill1,
				&es.Skill2, &es.Skill3, &es.Skill4, &es.Skill5,
			)
			additionalAttribute := [][]string{
				{"AdditionalCritical", "暴击"}, {"FendOff", "格挡"}, {"AdditionalHitting", "命中"}, {"CounterPunchRate", "反击"},
			}
			v := reflect.ValueOf(&es).Elem()
			for _, f := range additionalAttribute {
				if fv := v.FieldByName(f[0]); fv.Float() > 0 {
					es.AdditionalRate += fmt.Sprintf("%s:%s%% ", f[1], perDisp(float32(fv.Float()*100)))
				}
			}
			es.Skills = getSkills([]int{es.Skill1, es.Skill2, es.Skill3, es.Skill4, es.Skill5})
			ess = append(ess, es)
		}
		rows.Close()
		enemys=ess
	case "偷窃与掉落":     //drop
		sql := `
			select id,name,stolen_property,stolen_number,stolen_money,drop1id,drop1rate,drop2id,drop2rate,
			drop3id,drop3rate,drop4id,drop4rate,max_drop_money,min_drop_money from Monster
		`
		rows, _ := Db.Query(sql)
		for rows.Next() {
			ed := EnemyDrop{}
			rows.Scan(
				&ed.Id, &ed.Name, &ed.StolenProperty, &ed.StolenNumber,
				&ed.StolenMoney, &ed.Drop1ID, &ed.Drop1Rate, &ed.Drop2ID,
				&ed.Drop2Rate, &ed.Drop3ID, &ed.Drop3Rate, &ed.Drop4ID,
				&ed.Drop4Rate, &ed.MaxDropMoney, &ed.MinDropMoney,
			)
			if ed.StolenProperty != 0 {
				paName := getName("Property", ed.StolenProperty)
				if ed.StolenNumber > 1 {
					ed.Stolen = fmt.Sprintf("%s*%d", paName, ed.StolenNumber)
				} else {
					ed.Stolen = fmt.Sprintf("%s", paName)
				}
			} else if ed.StolenMoney != 0 {
				ed.Stolen = fmt.Sprintf("%s*%d", "金钱", ed.StolenMoney)
			}
			ed.Drop1 = getName("Property", ed.Drop1ID)
			ed.Drop2 = getName("Property", ed.Drop2ID)
			ed.Drop3 = getName("Property", ed.Drop3ID)
			ed.Drop4 = getName("Property", ed.Drop4ID)
			if ed.Drop1 != "" {
				ed.Drop1Per = fmt.Sprintf("%s%%", perDisp(float32(ed.Drop1Rate*100)))
			}
			if ed.Drop2 != "" {
				ed.Drop2Per = fmt.Sprintf("%s%%", perDisp(float32(ed.Drop2Rate*100)))
			}
			if ed.Drop3 != "" {
				ed.Drop3Per = fmt.Sprintf("%s%%", perDisp(float32(ed.Drop3Rate*100)))
			}
			if ed.Drop4 != "" {
				ed.Drop4Per = fmt.Sprintf("%s%%", perDisp(float32(ed.Drop4Rate*100)))
			}
			if ed.MaxDropMoney != ed.MinDropMoney {
				ed.DropMoney = fmt.Sprintf("%d~%d", ed.MinDropMoney, ed.MaxDropMoney)
			} else if ed.MaxDropMoney != 0 {
				ed.DropMoney = fmt.Sprintf("%d", ed.MaxDropMoney)
			}
			eds = append(eds, ed)
		}
		rows.Close()
		enemys=eds
    default:
        c.IndentedJSON(http.StatusNotFound, gin.H{"error": "参数错误，禁止查找！"})
        return
	}
	s,err:=json.Marshal(enemys)
	if err!=nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		client.Set(k, string(s), 36*time.Hour)
		c.IndentedJSON(http.StatusOK,enemys)
	}
}