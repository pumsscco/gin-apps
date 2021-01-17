package main

import (
    "encoding/json"
    "fmt"
    "github.com/gin-gonic/gin"
    "net/http"
    "time"
)
type PMission struct {
	Trunk int 	`json:"trunk"`
	QuestId int   `json:"quest_id"`
	Name string  `json:"name"`
	Picture string  `json:"picture,omitempty"`
	Description  string   `json:"description,omitempty"`
    StoryPer float32  `json:"story_per,omitempty"`
    StoryShow int `json:"-"`
    StoryShowName string `json:"story_show_name,omitempty"`
}
func pMission(c *gin.Context) {
	var (
        missions []PMission
        mt OneParam
		err error
		pattern string
	)
	if err = c.ShouldBindJSON(&mt); err != nil {    
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    k:=fmt.Sprintf("pal4:magic:%s",mt.Type)
    val,err:=client.Get(k).Result()
    if err==nil {
		json.Unmarshal([]byte(val),&missions)
		c.IndentedJSON(http.StatusOK,missions)
        return
    }
    switch mt.Type {
    case "主线":
        pattern=" depended_id<200"
    case "委托":
        pattern=" depended_id between 200 and 299"
    case "支线":
        pattern=" depended_id>=300"
    default:
        c.IndentedJSON(http.StatusNotFound, gin.H{"error": "参数错误，禁止查找！"})
        return
    }
    missionSql:=fmt.Sprintf(`
        select trunk,quest_id,name,picture,description,story_per,story_show from Mission where %s
    `,pattern)
    rows,_ := Db.Query(missionSql)
    for rows.Next() {
        mi := PMission{}
        rows.Scan(
            &mi.Trunk,&mi.QuestId,&mi.Name,&mi.Picture,
            &mi.Description,&mi.StoryPer,&mi.StoryShow,
        )
        if mi.StoryShow==1 {
            mi.StoryShowName="是"
        } else {
            mi.StoryShowName="否"
        }
        missions = append(missions, mi)
    }
    rows.Close()
    if len(missions)==0 {
        c.IndentedJSON(http.StatusNotFound, gin.H{"error": "参数错误，什么也查不到！"})
        return
    }
	s,err:=json.Marshal(missions)
	if err!=nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		client.Set(k, string(s), 36*time.Hour)
		c.IndentedJSON(http.StatusOK,missions)
	}
}