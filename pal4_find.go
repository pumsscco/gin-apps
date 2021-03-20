package main

import (
    "encoding/json"
    "fmt"
    "github.com/gin-gonic/gin"
    "net/http"
    "time"
)
type Method struct {
    Name string
    StolenEnemy string
    DropEnemy  string 
    BuyScene string
    PickList []Pick
}
//依据大类，子类与名称，获得各种得到该物品的渠道
func find(c *gin.Context) {
    var (
        method Method
        mt TwoParam
        err error
    )
    if err = c.ShouldBindJSON(&mt); err != nil {    
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    k:=fmt.Sprintf("pal4:thing:%s:%s",mt.Class,mt.Type)
    val,err:=client.Get(k).Result()
    if err==nil {
		json.Unmarshal([]byte(val),&method)
		c.IndentedJSON(http.StatusOK,method)
        return
    }
    switch mt.Class {
    case "装备":
        id:=getId("Equip",mt.Type)
        if id==-1 {
            c.IndentedJSON(http.StatusNotFound, gin.H{"error": "参数错误，什么也查不到！"})
            return
        }
        //购买
        method.BuyScene=getBuyScene(id)
        //拾取
        method.PickList=pickItem(id)
    case "道具":
        id:=getId("Property",mt.Type)
        if id==-1 {
            c.IndentedJSON(http.StatusNotFound, gin.H{"error": "参数错误，什么也查不到！"})
            return
        }
        method.BuyScene=getBuyScene(id)
        method.StolenEnemy=getStolen(id)
        method.DropEnemy=getDrop(id)
        method.PickList=pickItem(id)
    case "配方":
        id:=getId("Prescription",mt.Type)
        if id==-1 {
            c.IndentedJSON(http.StatusNotFound, gin.H{"error": "参数错误，什么也查不到！"})
            return
        }
        method.BuyScene=getBuyScene(id)
    default:
        c.IndentedJSON(http.StatusNotFound, gin.H{"error": "参数错误，禁止查找！"})
        return
    }
    method.Name=mt.Type
    s,err:=json.Marshal(method)
	if err!=nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		client.Set(k, string(s), 36*time.Hour)
		c.IndentedJSON(http.StatusOK,method)
	}
}