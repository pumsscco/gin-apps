package main

import (
    "encoding/json"
    "fmt"
    "github.com/gin-gonic/gin"
    "net/http"
    "time"
)
type Things struct {
    Cat,Type   string  //物品的大类与细类
    ThingList []string   //物品列表
}
func thing(c *gin.Context) {
    sql:=""
    switch cat {
    case "ea":
        eaType:=getId("EquipType",routeMapEquip[thingType])
        sql=fmt.Sprintf("select name from Equip where type=%d",eaType)
    case "pa":
        paType:=getId("PropertyClass",routeMapProperty[thingType][0])
        sql=fmt.Sprintf("select name from Property where type=%d",paType)
        switch routeMapProperty[thingType][1] {
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
        }
    case "pra":
        praType:=getId("PrescriptionType",routeMapPrescription[thingType][0])
        switch routeMapPrescription[thingType][0] {
        case "熔铸图谱":
            eaType:=getId("EquipType",routeMapPrescription[thingType][1])
            sql=fmt.Sprintf(`
                select pra.name from Prescription pra join Equip ea on pra.product_id=ea.id where pra.type=%d and ea.type=%d
            `,praType,eaType)
        default:
            sql=fmt.Sprintf("select name from Prescription where type=%d",praType)
        }
    }
    rows,_ := Db.Query(sql)
    for rows.Next() {
        var name string
        rows.Scan(&name)
        things.ThingList = append(things.ThingList, name)
    }
    things.Type=thingType
    things.Cat=cat
    return
}
