package main

import (
    "encoding/json"
    "fmt"
    "github.com/gin-gonic/gin"
    "net/http"
    "time"
)
type Question struct {
	Problem string 		`json:"problem"`
	Answer1 string		`json:"answer1"`
	Answer2 string  	`json:"answer2"`
	Answer3 string  	`json:"answer3"`
    RightAnswer    int  `json:"right_answer"`
}
func question(c *gin.Context) {
    var (
        questions []Question
        qt OneParam
		err error
		db int
    )
	if err = c.ShouldBindJSON(&qt); err != nil {    
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    k:=fmt.Sprintf("pal4:question:%s",qt.Type)
    val,err:=client.Get(k).Result()
    if err==nil {
		json.Unmarshal([]byte(val),&questions)
		c.IndentedJSON(http.StatusOK,questions)
        return
    }
    switch qt.Type {
    case "仙剑历史":
        db=1
    case "仙剑故事":
        db=2
    case "仙剑世界":
        db=3
    default:
        c.IndentedJSON(http.StatusNotFound, gin.H{"error": "参数错误，禁止查找！"})
        return
    }
    questionSql:=`
        select question,answer1,answer2,answer3,right_answer from GameQuestion where db=?
    `
    rows,_ := Db.Query(questionSql,db)
    for rows.Next() {
        qu := Question{}
        rows.Scan(
            &qu.Problem,&qu.Answer1,&qu.Answer2,&qu.Answer3,&qu.RightAnswer,
        )
        questions = append(questions, qu)
    }
    rows.Close()
	s,err:=json.Marshal(questions)
	if err!=nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		client.Set(k, string(s), 36*time.Hour)
		c.IndentedJSON(http.StatusOK,questions)
	}
}