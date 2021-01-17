package main

import (
    "encoding/json"
    "fmt"
    "github.com/gin-gonic/gin"
    "net/http"
    "time"
)
type Upgrade struct {
	Level int    `json:"level"`
	Experience int        `json:"experience"`
	MaxHP int              `json:"max_hp"`
	MaxMP int   `json:"max_mp"`
	Physical int       `json:"physical"`
	Toughness int        `json:"toughness"`
	Speed int       `json:"speed"`
	Lucky int      `json:"lucky"`
	Will  int    `json:"will"`
    FendOff  float32    `json:"-"`
    FendOffPer  string    `json:"fend_off_per"`
}

func upgrade(c *gin.Context) {
	var (
        upgrades []Upgrade
        role OneParam
		err error
	)
	if err = c.ShouldBindJSON(&role); err != nil {    
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
	}
	k:=fmt.Sprintf("pal4:stunt:%s",role.Type)
    val,err:=client.Get(k).Result()
    if err==nil {
		json.Unmarshal([]byte(val),&upgrades)
		c.IndentedJSON(http.StatusOK,upgrades)
        return
    }
    roleId:=getId("Role",role.Type)
    upgradeSql:=`
        select level,experience,max_hp,max_mp,physical,toughness,speed,lucky,will,fend_off from UpgradeData where role_id=?
    `
    rows,_ := Db.Query(upgradeSql,roleId)
    for rows.Next() {
        up := Upgrade{}
        rows.Scan(
            &up.Level,&up.Experience,&up.MaxHP,&up.MaxMP,&up.Physical,&up.Toughness,&up.Speed,
            &up.Lucky,&up.Will,&up.FendOff,
        )
        up.FendOffPer=fmt.Sprintf("%s%%",perDisp(float32(up.FendOff*100)))
        upgrades = append(upgrades, up)
    }
	rows.Close()
	if len(upgrades)==0 {
        c.IndentedJSON(http.StatusNotFound, gin.H{"error": "参数错误，什么也查不到！"})
        return
    }
	s,err:=json.Marshal(upgrades)
	if err!=nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		client.Set(k, string(s), 36*time.Hour)
		c.IndentedJSON(http.StatusOK,upgrades)
	}
}
