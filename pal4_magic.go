package main

import (
    "encoding/json"
    "fmt"
    "github.com/gin-gonic/gin"
    "net/http"
    "time"
)
type Magic struct {
    Id int 			`json:"id"`
	Name string 		`json:"name"`		
	Description string  `json:"description,omitempty"`
	Attribute  string   `json:"attribute,omitempty"`
	AiCmdType int 		`json:"-"`
	Target int   	`json:"-"`
	Wuling int  `json:"-"`
	AttachedSkill int  `json:"attached_skill,omitempty"`
	ConsumedMP int    `json:"consumed_mp,omitempty"`
	AiCmdName string  `json:"ai_cmd_name,omitempty"`
	TargetName string  `json:"target_name,omitempty"`
	WulingName string    `json:"wuling_name,omitempty"`
	Animation string  `json:"animation,omitempty"`
	TargetEffect string `json:"target_effect,omitempty"`
	TargetBind string  `json:"target_bind,omitempty"`
}

func magic(c *gin.Context) {
	var (
        magics []Magic
        mt OneParam
		err error
		propSelect string
	)
	if err = c.ShouldBindJSON(&mt); err != nil {    
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
	}
	k:=fmt.Sprintf("pal4:magic:%s",mt.Type)
    val,err:=client.Get(k).Result()
    if err==nil {
		json.Unmarshal([]byte(val),&magics)
		c.IndentedJSON(http.StatusOK,magics)
        return
    }
    switch mt.Type {
    case "我方":
        propSelect=" attribute!=''"
    case "敌方":
        propSelect=" attribute=''"
    default:
        c.IndentedJSON(http.StatusNotFound, gin.H{"error": "参数错误，禁止查找！"})
        return
    }
    magicSql:=fmt.Sprintf(`
        select id,name,description,attribute,ai_cmd_type,target,wuling,attached_skill,consumed_mp,
        animation,target_ef,target_bind from Magic where %s
    `,propSelect)
    rows,_ := Db.Query(magicSql)
    for rows.Next() {
        magic := Magic{}
        rows.Scan(
            &magic.Id,&magic.Name,&magic.Description,&magic.Attribute,&magic.AiCmdType,&magic.Target,
            &magic.Wuling,&magic.AttachedSkill,&magic.ConsumedMP,&magic.Animation,&magic.TargetEffect,
            &magic.TargetBind,
        )
        magic.AiCmdName=getName("AiCommandType",magic.AiCmdType)
        magic.TargetName=getName("SkillTarget",magic.Target)
        magic.WulingName=getName("WuLing",magic.Wuling)
        magics = append(magics, magic)
    }
    rows.Close()
    if len(magics)==0 {
        c.IndentedJSON(http.StatusNotFound, gin.H{"error": "参数错误，什么也查不到！"})
        return
    }
	s,err:=json.Marshal(magics)
	if err!=nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		client.Set(k, string(s), 36*time.Hour)
		c.IndentedJSON(http.StatusOK,magics)
	}
}
