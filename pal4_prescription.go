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
	ProductId int   `json:"attribute,omitempty"`
	NeedPotential int   `json:"attribute,omitempty"`
	NeedLing int  `json:"attribute,omitempty"`
	Skill int   `json:"attribute,omitempty"`
	Price int `json:"attribute,omitempty"`
	Effect1 string    `json:"-"`
	Effect2 string   `json:"-"`
	Effect3 string   `json:"-"`
	Effect string   `json:"-"`
	BuyScene string `json:"buy_scene,omitempty"`
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
	rows,err := Db.Query(prescriptionSql,typeId)
	if err!=nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return	
	}
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
