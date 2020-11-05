package main

import (
	"time"
	"encoding/json"
)

type Mission struct {
    Id int
    Content string
}
func getMission() (missions []Mission)  {
    val,err:=client.Get("mission").Result()
    if err==nil {
        json.Unmarshal([]byte(val),&missions)
        return
    }
    sql:=`select id,content from SPMission where content!=""`
    rows,_ := Db.Query(sql)
    for rows.Next() {
        m := Mission{}
        rows.Scan(&m.Id,&m.Content)
        missions = append(missions, m)
    }
    rows.Close()
    as,err:=json.Marshal(missions)
	client.Set("mission", string(as), 12*time.Hour)
	if err!=nil {
		logger.Print(err)
	}
    return
}