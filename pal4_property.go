package main

import (
    "encoding/json"
    "fmt"
    "github.com/gin-gonic/gin"
    "net/http"
    "time"
)
type Property struct {
    Id              int  		`json:"id"`
	Name string 				`json:"name"`
	Description string   	`json:"description"`
	Attribute  string  		`json:"attribute,omitempty"`
	Model string    `json:"model,omitempty"`
	Texture     string   `json:"texture,omitempty"`
    Level int    `json:"-"`
    LevelName string `json:"level_name,omitempty"`
	Price string   			`json:"price,omitempty"`
	AttachedSkill          int    `json:"attached_skill,omitempty"`
    BuyScene string `json:"buy_scene,omitempty"`
}

//依据装备类型的中文名，获得该类物品的全部属性
func property(c *gin.Context) {
	var (
        properties []Property
        pt TwoParam
        err error
        validProp=map[string]string{
            "食物":"恢复",
            "其它恢复类":"恢复",
            "攻击类":"攻击",
            "香料":"辅助",
            "其它辅助类":"辅助",
            "矿石":"材料",
            "尸块":"材料",
            "其它材料":"材料",
            "剧情类":"剧情",
        }
        valid bool
    )
	if err = c.ShouldBindJSON(&pt); err != nil {    
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    for k,v:=range validProp {
        if pt.Class==v && pt.Type==k {
            valid=true
            break
        }
    }
    if !valid {
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "非法参数"})
        return
    }
    k:=fmt.Sprintf("pal4:property:%s:%s",pt.Class,pt.Type)
    val,err:=client.Get(k).Result()
    if err==nil {
		json.Unmarshal([]byte(val),&properties)
		c.IndentedJSON(http.StatusOK,properties)
        return
    }
    //先依名称取类型ID
    typeId:=getId("PropertyClass",pt.Class)
    //利用类型ID取原始数据
    propertySql:=`
        select id,name,description,attribute,model,texture,property_level,price,attached_skill 
        from Property where type=?
    `
    switch pt.Type {
    case "食物":
        propertySql+=` and model regexp "^SW"`
    case "其它恢复类":
        propertySql+=` and model not regexp "^SW"`
    case "香料":
        propertySql+=` and model regexp "^CX"`
    case "其它辅助类":
        propertySql+=` and model not regexp "^CX"`
    case "矿石":
        propertySql+=` and model regexp "^CK" and attribute="熔铸、锻冶的材料"`
    case "尸块":
        propertySql+=` and attribute="注灵的材料"`
    case "其它材料":
        propertySql+=` and (model regexp "^CQ" or attribute="")`
    }
    rows,_ := Db.Query(propertySql,typeId)
    for rows.Next() {
        prop := Property{}
        rows.Scan(
            &prop.Id,&prop.Name,&prop.Description,&prop.Attribute,&prop.Model,
            &prop.Texture,&prop.Level,&prop.Price,&prop.AttachedSkill,
        )
        if pt.Type=="其它恢复类" {
			prop.LevelName=getName("PropertyLevel",prop.Level)
		}
        prop.BuyScene=getBuyScene(prop.Id)
        properties = append(properties, prop)
    }
    rows.Close()
	s,err:=json.Marshal(properties)
	if err!=nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		client.Set(k, string(s), 36*time.Hour)
		c.IndentedJSON(http.StatusOK,properties)
	}
}
