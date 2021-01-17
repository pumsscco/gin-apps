package main

import (
    "encoding/json"
    "fmt"
    "github.com/gin-gonic/gin"
    "net/http"
    "time"
)
type Stunt struct {
    Id int 				`json:"id"`
	Name 		string 		`json:"name"`
	Description string     `json:"description,omitempty"`
	Attribute  string   `json:"attribute,omitempty"`
	AiCmdType int   `json:"-"`
	Target int    `json:"-"`
	RequireLevel int  `json:"require_level,omitempty"`
	AttachedSkill int  `json:"attached_skill,omitempty"`
	ActType int     `json:"-"`
	AiCmdName string    `json:"ai_cmd_name,omitempty"`
	TargetName string   `json:"target_name,omitempty"`
	ActTypeName string `json:"act_type_name,omitempty"`
	Animation string   `json:"animation,omitempty"`
	TargetEffect string  `json:"target_effect,omitempty"`
	TargetBind string `json:"target_bind,omitempty"`
	TianheLvl int     `json:"-"`
	LingshaLvl int    `json:"-"`
	MengliLvl int  `json:"-"`
	ZiyingLvl int   `json:"-"`
    Role string   `json:"role,omitempty"`
}
func stunt(c *gin.Context) {
	var (
        stunts []Stunt
        stt OneParam
		err error
		propSelect string
	)
	if err = c.ShouldBindJSON(&stt); err != nil {    
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
	}
	k:=fmt.Sprintf("pal4:stunt:%s",stt.Type)
    val,err:=client.Get(k).Result()
    if err==nil {
		json.Unmarshal([]byte(val),&stunts)
		c.IndentedJSON(http.StatusOK,stunts)
        return
    }
    switch stt.Type {
    case "我方":
        propSelect=" attribute!=''"
    case "敌方":
		propSelect=" attribute=''"
	default:
        c.IndentedJSON(http.StatusNotFound, gin.H{"error": "参数错误，禁止查找！"})
        return
    }
    stuntSql:=fmt.Sprintf(`
        select id,name,description,attribute,ai_cmd_type,target,tianhe_lv_lmt,lingsha_lv_lmt,mengli_lv_lmt,ziying_lv_lmt lv_lmt,
        attached_skill,act_type,animation,target_ef,target_bind from Stunt where %s
    `,propSelect)
    rows,_ := Db.Query(stuntSql)
    for rows.Next() {
        st := Stunt{}
        rows.Scan(
            &st.Id,&st.Name,&st.Description,&st.Attribute,&st.AiCmdType,&st.Target,&st.TianheLvl,
            &st.LingshaLvl,&st.MengliLvl,&st.ZiyingLvl,&st.AttachedSkill,&st.ActType,&st.Animation,
            &st.TargetEffect,&st.TargetBind,
        )
        switch {
		case st.MengliLvl!=0:
            st.Role="柳梦璃"
        case st.ZiyingLvl!=0:
            st.Role="慕容紫英"
        case st.LingshaLvl!=0:
            st.Role="韩菱纱"
        case st.TianheLvl!=0:
            st.Role="云天河"
        }
        st.RequireLevel=st.TianheLvl|st.LingshaLvl|st.MengliLvl|st.ZiyingLvl
        st.AiCmdName=getName("AiCommandType",st.AiCmdType)
        st.TargetName=getName("SkillTarget",st.Target)
        st.ActTypeName=getName("ActType",st.ActType)
        stunts = append(stunts, st)
    }
    rows.Close()
	s,err:=json.Marshal(stunts)
	if err!=nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		client.Set(k, string(s), 36*time.Hour)
		c.IndentedJSON(http.StatusOK,stunts)
	}
}
