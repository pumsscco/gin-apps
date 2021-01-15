package main

import (
    "encoding/json"
    "fmt"
    "github.com/gin-gonic/gin"
    "net/http"
    "reflect"
    "strings"
    "time"
)
type Prescription struct {
    Id       int     `json:"id"`
	Name string   `json:"name"`
	Description string   `json:"description"`
	Attribute  string  `json:"attribute,omitempty"`
	ProductId int   `json:"attribute,omitempty"`
	NeedPotential int   `json:"attribute,omitempty"`
	NeedLing int  `json:"attribute,omitempty"`
	Skill int   `json:"attribute,omitempty"`
	Price int `json:"attribute,omitempty"`
	Effect1 string    `json:"attribute,omitempty"`
	Effect2 string   `json:"attribute,omitempty"`
	Effect3 string   `json:"attribute,omitempty"`
	Effect string   `json:"attribute,omitempty"`
	BuyScene string //粒子特效剑、双剑、琴、粒子特效组合、购买场景
    EaUrl string
}
func prescription(c *gin.Context) {
	var (
        prescriptions []Prescription
		pt PropType
		err error
		prescriptionSql string
	)
	if err = c.ShouldBindJSON(&pt); err != nil {    
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    k:=fmt.Sprintf("pal4:prescription:%s:%s",pt.Class,pt.Type)
    val,err:=client.Get(k).Result()
    if err==nil {
		json.Unmarshal([]byte(val),&equips)
		c.IndentedJSON(http.StatusOK,equips)
        return
    }
    
    typeId:=getId("PrescriptionType",pt.Class)
    if pt.Class=="熔铸图谱" {
        equipTypeId:=getId("EquipType",prescriptionType[1])
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
    rows,_ := Db.Query(prescriptionSql,typeId)
    for rows.Next() {
        prescription := Prescription{}
        rows.Scan(
            &prescription.Id,&prescription.Name,&prescription.Description,&prescription.Attribute,
            &prescription.ProductId,&prescription.NeedPotential,&prescription.NeedLing,&prescription.Skill,
            &prescription.Price,&prescription.Effect1,&prescription.Effect2,&prescription.Effect3,
        )
        for _,eff:=range []string{prescription.Effect1,prescription.Effect2,prescription.Effect3} {
            if eff!="" {
                prescription.Effect+=fmt.Sprintf("%s ",eff)
            }
        }
        prescription.Effect=strings.TrimSuffix(prescription.Effect," ")
        prescription.BuyScene=getBuyScene(prescription.Id)
        for k,v:=range routeMapEquip {
            if v==prescriptionType[1] {
                prescription.EaUrl=k
                break
            }
        }
        prescriptionList = append(prescriptionList, prescription)
    }
    rows.Close()
    prescriptions.PrescriptionList=prescriptionList
    prescriptions.Type=prescriptionType[1]
    return
}
