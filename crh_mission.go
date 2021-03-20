package main

import (
    "encoding/json"
    "github.com/gin-gonic/gin"
    "net/http"
    "time"
)

type CMission struct {
    Id int `json:"id"`
    Content string  `json:"content"`
}
func cMission(c *gin.Context)  {
	var missions []CMission
	k:="crh:mission"
    val,err:=client.Get(k).Result()
    if err==nil {
		json.Unmarshal([]byte(val),&missions)
		c.IndentedJSON(http.StatusOK,missions)
        return
    }
    sql:=`select id,content from SPMission where content!=""`
    rows,_ := Db.Query(sql)
    for rows.Next() {
        m := CMission{}
        rows.Scan(&m.Id,&m.Content)
        missions = append(missions, m)
    }
    rows.Close()
    s,err:=json.Marshal(missions)
	if err!=nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
	} else {
		client.Set(k, string(s), 36*time.Hour)
        c.IndentedJSON(http.StatusOK,missions)
	}
}