package main

import (
    "encoding/json"
    "fmt"
    "github.com/gin-gonic/gin"
    "net/http"
    "time"
)
type Scene struct {
	Id int  `json:"id"`
	Scene string  `json:"scene"`
	Section string   `json:"section"`
	Name string  `json:"name"`
	Type int   `json:"-"`
	TypeName string  `json:"type_name"`
	EarthBall int  `json:"-"`
	IsEarthBall string  `json:"is_earthball"`
}
func scene(c *gin.Context) {
	var (
        scenes []Scene
        st OneParam
		err error
		sql string
    )
	if err = c.ShouldBindJSON(&st); err != nil {    
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    k:=fmt.Sprintf("pal4:scene:%s",st.Type)
    val,err:=client.Get(k).Result()
    if err==nil {
		json.Unmarshal([]byte(val),&scenes)
		c.IndentedJSON(http.StatusOK,scenes)
        return
	}
	switch st.Type {
	case "迷宫":
		sql=`select id,scene,section,name,type,earthball from Scene where scene regexp "^M"`
	case "城镇":
		sql=`select id,scene,section,name,type,earthball from Scene where scene regexp "^Q"`
	default:
        c.IndentedJSON(http.StatusNotFound, gin.H{"error": "参数错误，禁止查找！"})
        return
	}
    rows,_ := Db.Query(sql)
    for rows.Next() {
        sc:=Scene{}
		rows.Scan(&sc.Id,&sc.Scene,&sc.Section,&sc.Name,&sc.Type,&sc.EarthBall)
		switch sc.Type {
		case 0:
			sc.TypeName="城镇室外"
		case 1:
			sc.TypeName="城镇室内"
		case 2:
			sc.TypeName="迷宫"
		}
		if sc.EarthBall==0 {
			sc.IsEarthBall="城镇中，不能用土灵珠！"
		} else if  sc.EarthBall==1 {
			sc.IsEarthBall="迷宫中，可用土灵珠！"
		}
        scenes = append(scenes, sc)
    }
	rows.Close()
	s,err:=json.Marshal(scenes)
	if err!=nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		client.Set(k, string(s), 36*time.Hour)
		c.IndentedJSON(http.StatusOK,scenes)
	}
}