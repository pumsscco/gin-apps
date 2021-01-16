package main

import (
    "encoding/json"
    "fmt"
    "github.com/gin-gonic/gin"
    "net/http"
    "time"
)
type Method struct {
    Name,StolenEnemy,DropEnemy  string 
    BuyScene template.HTML
    HasNonePickup,HasPickup bool 
    PickList []Pick
}
//依据大类，子类与名称，获得各种得到该物品的渠道
func find(c *gin.Context) {
    switch cat {
    case "ea":
        id:=getId("Equip",thing)
        //购买
        method.BuyScene=template.HTML(getBuyScene(id))
        if method.BuyScene=="" && t!="p" {
            tName:=routeMapEquip[t]
            method.BuyScene=template.HTML(fmt.Sprintf(`<a href="/%s/%s\">%s配方</a>`,"pra",t,tName))
        }
        //拾取
        method.PickList=pickItem(id)
    case "pa":
        id:=getId("Property",thing)
        method.BuyScene=template.HTML(getBuyScene(id))
        method.StolenEnemy=getStolen(id)
        method.DropEnemy=getDrop(id)
        method.PickList=pickItem(id)
    case "pra":
        id:=getId("Prescription",thing)
        method.BuyScene=template.HTML(getBuyScene(id))
    }
    if method.BuyScene!="" || method.StolenEnemy!="" || method.DropEnemy!="" {
        method.HasNonePickup=true
    }
    if method.PickList!=nil {
        method.HasPickup=true
    }
    method.Name=thing
    return
}