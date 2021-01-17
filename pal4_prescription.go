package main

import (
    "encoding/json"
    "fmt"
    "github.com/gin-gonic/gin"
    "net/http"
    "strings"
    "time"
)
type Prescription struct {
    Id       int     `json:"id"`
	Name string   `json:"name"`
	Description string   `json:"description"`
	Attribute  string  `json:"attribute,omitempty"`
	ProductId int   `json:"product_id,omitempty"`
	NeedPotential int   `json:"need_potential,omitempty"`
	NeedLing int  `json:"need_ling,omitempty"`
	Skill int   `json:"skill,omitempty"`
	Price int `json:"price,omitempty"`
	Effect1 string    `json:"-"`
	Effect2 string   `json:"-"`
	Effect3 string   `json:"-"`
	Effect string   `json:"-"`
	BuyScene string `json:"buy_scene,omitempty"`
}
func prescription(c *gin.Context) {
	var (
        prescriptions []Prescription
		pt TwoParam
		err error
        prescriptionSql string
        validPres=map[string]string{
            "剑":"熔铸图谱",
            "双剑":"熔铸图谱",
            "琴":"熔铸图谱",
            "头部防具":"熔铸图谱",
            "身体防具":"熔铸图谱",
            "足部防具":"熔铸图谱",
            "锻冶":"锻造图谱",
            "注灵":"注灵图谱",
        }
        valid bool
    )
	if err = c.ShouldBindJSON(&pt); err != nil {    
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    for k,v:=range validPres {
        if pt.Class==v && pt.Type==k {
            valid=true
            break
        }
    }
    if !valid {
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "非法参数"})
        return
    }
    k:=fmt.Sprintf("pal4:prescription:%s:%s",pt.Class,pt.Type)
    val,err:=client.Get(k).Result()
    if err==nil {
		json.Unmarshal([]byte(val),&prescriptions)
		c.IndentedJSON(http.StatusOK,prescriptions)
        return
    }
    typeId:=getId("PrescriptionType",pt.Class)
    if pt.Class=="熔铸图谱" {
        equipTypeId:=getId("EquipType",pt.Type)
        prescriptionSql=fmt.Sprintf(`
            select pra.id,pra.name,pra.description,attribute,product_id,need_potential,need_ling,pra.skill_id,pra.price,pra.ef2,pra.ef3,ef4 
            from Prescription pra join Equip ea on pra.product_id=ea.id where pra.type=? and ea.type=%d
        `,equipTypeId)
    } else {
        prescriptionSql=`
            select id,name,description,attribute,product_id,need_potential,need_ling,skill_id,price,ef2,ef3,ef4 
            from Prescription where type=?
        `
	}
	logger.Printf("prescriptionSql: %v; typeId: %v\n",prescriptionSql,typeId)
	rows,_ := Db.Query(prescriptionSql,typeId)
    for rows.Next() {
        pres := Prescription{}
        rows.Scan(
            &pres.Id,&pres.Name,&pres.Description,&pres.Attribute,
            &pres.ProductId,&pres.NeedPotential,&pres.NeedLing,&pres.Skill,
            &pres.Price,&pres.Effect1,&pres.Effect2,&pres.Effect3,
        )
        for _,eff:=range []string{pres.Effect1,pres.Effect2,pres.Effect3} {
            if eff!="" {
                pres.Effect+=fmt.Sprintf("%s ",eff)
            }
        }
        pres.Effect=strings.TrimSuffix(pres.Effect," ")
        pres.BuyScene=getBuyScene(pres.Id)
        prescriptions = append(prescriptions, pres)
    }
	rows.Close()
	s,err:=json.Marshal(prescriptions)
	if err!=nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		client.Set(k, string(s), 36*time.Hour)
		c.IndentedJSON(http.StatusOK,prescriptions)
	}
}
