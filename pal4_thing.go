package main

import (
    "encoding/json"
    "fmt"
    "github.com/gin-gonic/gin"
    "net/http"
    "time"
)
func thing(c *gin.Context) {
    var (
        things []string
        tht ThreeParam
        err error
        sql string
    )
    if err = c.ShouldBindJSON(&tht); err != nil {    
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    k:=fmt.Sprintf("pal4:thing:%s:%s:%s",tht.Class,tht.Type,tht.SubType)
    val,err:=client.Get(k).Result()
    if err==nil {
		json.Unmarshal([]byte(val),&things)
		c.IndentedJSON(http.StatusOK,things)
        return
    }
    switch tht.Class {
    case "装备":
        eaClass:=getId("EquipClass",tht.Type)
        eaType:=getId("EquipType",tht.SubType)
        matchSql:="select id from EquipRelation where class=? and type=?"
        var relId int
        err:=Db.QueryRow(matchSql,eaClass,eaType).Scan(&relId)
        if err!=nil {
            logger.Println("Category and Subclass mismatch!")
            c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        sql=fmt.Sprintf("select name from Equip where type=%d",eaType)
    case "道具":
        paType:=getId("PropertyClass",tht.Type)
        sql=fmt.Sprintf("select name from Property where type=%d",paType)
        switch tht.SubType {
        case "食物":
            sql+=` and model regexp "^SW"`
        case "其它恢复类":
            sql+=` and model not regexp "^SW"`
        case "香料":
            sql+=` and model regexp "^CX"`
        case "其它辅助类":
            sql+=` and model not regexp "^CX"`
        case "矿石":
            sql+=` and model regexp "^CK" and attribute="熔铸、锻冶的材料"`
        case "尸块":
            sql+=` and attribute="注灵的材料"`
        case "其它材料":
            sql+=` and (model regexp "^CQ" or attribute="")`
        default:
            c.IndentedJSON(http.StatusNotFound, gin.H{"error": "参数错误，禁止查找！"})
            return
        }
    case "配方":
        praType:=getId("PrescriptionType",tht.Type)
        switch tht.Type {
        case "熔铸图谱":
            eaType:=getId("EquipType",tht.SubType)
            sql=fmt.Sprintf(`
                select pra.name from Prescription pra join Equip ea on pra.product_id=ea.id where pra.type=%d and ea.type=%d
            `,praType,eaType)
        case "锻造图谱":
            if tht.SubType=="锻冶" {
                sql=fmt.Sprintf("select name from Prescription where type=%d",praType)
            } else {
                c.IndentedJSON(http.StatusNotFound, gin.H{"error": "参数错误，禁止查找！"})
                return
            }
        case "注灵图谱":
            if tht.SubType=="注灵" {
                sql=fmt.Sprintf("select name from Prescription where type=%d",praType)
            } else {
                c.IndentedJSON(http.StatusNotFound, gin.H{"error": "参数错误，禁止查找！"})
                return
            }
        default:
            c.IndentedJSON(http.StatusNotFound, gin.H{"error": "参数错误，禁止查找！"})
            return
        }
    default:
        c.IndentedJSON(http.StatusNotFound, gin.H{"error": "参数错误，禁止查找！"})
        return
    }
    rows,_ := Db.Query(sql)
    for rows.Next() {
        var name string
        rows.Scan(&name)
        things = append(things, name)
    }
    rows.Close()
    s,err:=json.Marshal(things)
	if err!=nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		client.Set(k, string(s), 36*time.Hour)
		c.IndentedJSON(http.StatusOK,things)
	}
}
