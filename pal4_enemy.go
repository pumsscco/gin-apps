package main

import (
	"fmt"
	"reflect"
	"strings"
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
func enemy(c *gin.Context) {
	var (
		ecs []EnemyCommon
		ebs []EnemyBasic
		ers []EnemyResistance
		ess []EnemySkill
		eds []EnemyDrop
        et Type
		err error
		es interface{}
	)
	if err = c.ShouldBindJSON(&et); err != nil {    
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    k:=fmt.Sprintf("pal4:enemy:%s",et.Type)
    val,err:=client.Get(k).Result()
    if err==nil {
		json.Unmarshal([]byte(val),&es)
		c.IndentedJSON(http.StatusOK,es)
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
	case "八基本属性与五灵":     //basic
	case "抗性、反弹、受伤音效":     //resistance
	case "物理追加与技能":     //skill
	case "偷窃与掉落":     //drop
	}

}



func getEnemyBasic() (enemyBasics EnemyBasics) {
	sql := `
        select id,name,max_hp,rage,max_mp,physical,toughness,speed,lucky,will,water,fire,thunder,air,earth from Monster
    `
	enemyBasicList := []EnemyBasic{}
	rows, _ := Db.Query(sql)
	for rows.Next() {
		enemyBasic := EnemyBasic{}
		rows.Scan(
			&enemyBasic.Id, &enemyBasic.Name, &enemyBasic.HP, &enemyBasic.Rage, &enemyBasic.MP, &enemyBasic.Physical,
			&enemyBasic.Toughness, &enemyBasic.Speed, &enemyBasic.Lucky, &enemyBasic.Will, &enemyBasic.Water,
			&enemyBasic.Fire, &enemyBasic.Thunder, &enemyBasic.Air, &enemyBasic.Earth,
		)
		wulingAttribute := [][]string{
			{"Water", "水"}, {"Fire", "火"}, {"Thunder", "雷"},
			{"Air", "风"}, {"Earth", "土"},
		}
		v := reflect.ValueOf(enemyBasic)
		for _, f := range wulingAttribute {
			if fv := v.FieldByName(f[0]); fv.Int() > 0 {
				enemyBasic.Wuling += fmt.Sprintf("%s:%d ", f[1], fv.Int())
			}
		}
		enemyBasic.Wuling = strings.TrimSuffix(enemyBasic.Wuling, " ")
		enemyBasicList = append(enemyBasicList, enemyBasic)
	}
	rows.Close()
	enemyBasics.EnemyBasicList = enemyBasicList
	enemyBasics.Part = "八基本属性与五灵"
	return
}



func getEnemyResistance() (enemyResistances EnemyResistances) {
	sql := `
        select id,name,physical_extract,water_extract,fire_extract,thunder_extract,air_extract,earth_extract,
        physical_react,water_react,fire_react,thunder_react,air_react,earth_react,
        sound_wounded1,sound_wounded2,sound_wounded3 from Monster
    `
	enemyResistanceList := []EnemyResistance{}
	rows, _ := Db.Query(sql)
	for rows.Next() {
		enemyResistance := EnemyResistance{}
		rows.Scan(
			&enemyResistance.Id, &enemyResistance.Name, &enemyResistance.PhysicalExtract, &enemyResistance.WaterExtract,
			&enemyResistance.FireExtract, &enemyResistance.ThunderExtract, &enemyResistance.AirExtract,
			&enemyResistance.EarthExtract, &enemyResistance.PhysicalReact, &enemyResistance.WaterReact,
			&enemyResistance.FireReact, &enemyResistance.ThunderReact, &enemyResistance.AirReact,
			&enemyResistance.EarthReact, &enemyResistance.SoundWounded1, &enemyResistance.SoundWounded2,
			&enemyResistance.SoundWounded3,
		)
		reactAttribute := [][]string{
			{"PhysicalReact", "物理"}, {"WaterReact", "水"}, {"FireReact", "火"},
			{"ThunderReact", "雷"}, {"AirReact", "风"}, {"EarthReact", "土"},
		}
		v := reflect.ValueOf(&enemyResistance).Elem()
		for _, f := range reactAttribute {
			if fv := v.FieldByName(f[0]); fv.Float() > 0 {
				enemyResistance.React += fmt.Sprintf("%s:%s%% ", f[1], perDisp(float32(fv.Float()*100)))
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
		enemyResistance.React = strings.TrimSuffix(enemyResistance.React, " ")
		enemyResistanceList = append(enemyResistanceList, enemyResistance)
	}
	rows.Close()
	enemyResistances.EnemyResistanceList = enemyResistanceList
	enemyResistances.Part = "抗性、反弹、受伤音效"
	return
}



func getEnemySkill() (enemySkills EnemySkills) {
	sql := `
        select id,name,physical_additional,additional_critical,fend_off,additional_hitting,counterpunch_rate,
        skill1,skill2,skill3,skill4,skill5 from Monster
`
	enemySkillList := []EnemySkill{}
	rows, _ := Db.Query(sql)
	for rows.Next() {
		enemySkill := EnemySkill{}
		rows.Scan(
			&enemySkill.Id, &enemySkill.Name, &enemySkill.PhysicalAdditional, &enemySkill.AdditionalCritical,
			&enemySkill.FendOff, &enemySkill.AdditionalHitting, &enemySkill.CounterPunchRate, &enemySkill.Skill1,
			&enemySkill.Skill2, &enemySkill.Skill3, &enemySkill.Skill4, &enemySkill.Skill5,
		)
		additionalAttribute := [][]string{
			{"AdditionalCritical", "暴击"}, {"FendOff", "格挡"}, {"AdditionalHitting", "命中"}, {"CounterPunchRate", "反击"},
		}
		v := reflect.ValueOf(&enemySkill).Elem()
		for _, f := range additionalAttribute {
			if fv := v.FieldByName(f[0]); fv.Float() > 0 {
				enemySkill.AdditionalRate += fmt.Sprintf("%s:%s%% ", f[1], perDisp(float32(fv.Float()*100)))
			}
		}
		enemySkill.Skills = getSkills([]int{enemySkill.Skill1, enemySkill.Skill2, enemySkill.Skill3, enemySkill.Skill4, enemySkill.Skill5})
		enemySkillList = append(enemySkillList, enemySkill)
	}
	rows.Close()
	enemySkills.EnemySkillList = enemySkillList
	enemySkills.Part = "物理追加与技能"
	return
}



func getEnemyDrop() (enemyDrops EnemyDrops) {
	sql := `
        select id,name,stolen_property,stolen_number,stolen_money,drop1id,drop1rate,drop2id,drop2rate,
        drop3id,drop3rate,drop4id,drop4rate,max_drop_money,min_drop_money from Monster
    `
	enemyDropList := []EnemyDrop{}
	rows, _ := Db.Query(sql)
	for rows.Next() {
		enemyDrop := EnemyDrop{}
		rows.Scan(
			&enemyDrop.Id, &enemyDrop.Name, &enemyDrop.StolenProperty, &enemyDrop.StolenNumber,
			&enemyDrop.StolenMoney, &enemyDrop.Drop1ID, &enemyDrop.Drop1Rate, &enemyDrop.Drop2ID,
			&enemyDrop.Drop2Rate, &enemyDrop.Drop3ID, &enemyDrop.Drop3Rate, &enemyDrop.Drop4ID,
			&enemyDrop.Drop4Rate, &enemyDrop.MaxDropMoney, &enemyDrop.MinDropMoney,
		)
		if enemyDrop.StolenProperty != 0 {
			paName := getName("Property", enemyDrop.StolenProperty)
			if enemyDrop.StolenNumber > 1 {
				enemyDrop.Stolen = fmt.Sprintf("%s*%d", paName, enemyDrop.StolenNumber)
			} else {
				enemyDrop.Stolen = fmt.Sprintf("%s", paName)
			}
		} else if enemyDrop.StolenMoney != 0 {
			enemyDrop.Stolen = fmt.Sprintf("%s*%d", "金钱", enemyDrop.StolenMoney)
		}
		enemyDrop.Drop1 = getName("Property", enemyDrop.Drop1ID)
		enemyDrop.Drop2 = getName("Property", enemyDrop.Drop2ID)
		enemyDrop.Drop3 = getName("Property", enemyDrop.Drop3ID)
		enemyDrop.Drop4 = getName("Property", enemyDrop.Drop4ID)
		if enemyDrop.Drop1 != "" {
			enemyDrop.Drop1Per = fmt.Sprintf("%s%%", perDisp(float32(enemyDrop.Drop1Rate*100)))
		}
		if enemyDrop.Drop2 != "" {
			enemyDrop.Drop2Per = fmt.Sprintf("%s%%", perDisp(float32(enemyDrop.Drop2Rate*100)))
		}
		if enemyDrop.Drop3 != "" {
			enemyDrop.Drop3Per = fmt.Sprintf("%s%%", perDisp(float32(enemyDrop.Drop3Rate*100)))
		}
		if enemyDrop.Drop4 != "" {
			enemyDrop.Drop4Per = fmt.Sprintf("%s%%", perDisp(float32(enemyDrop.Drop4Rate*100)))
		}
		if enemyDrop.MaxDropMoney != enemyDrop.MinDropMoney {
			enemyDrop.DropMoney = fmt.Sprintf("%d~%d", enemyDrop.MinDropMoney, enemyDrop.MaxDropMoney)
		} else if enemyDrop.MaxDropMoney != 0 {
			enemyDrop.DropMoney = fmt.Sprintf("%d", enemyDrop.MaxDropMoney)
		}
		enemyDropList = append(enemyDropList, enemyDrop)
	}
	rows.Close()
	enemyDrops.EnemyDropList = enemyDropList
	enemyDrops.Part = "偷窃与掉落"
	return
}
