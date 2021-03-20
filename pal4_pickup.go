package main

import (
    "encoding/json"
    "fmt"
    "github.com/gin-gonic/gin"
    "net/http"
    "time"
    "reflect"
)
type PItem struct {
    ItemId string  `json:"item_id"`
    Model string   `json:"model,omitempty"`
    Texture  string   `json:"texture"`
    SectionName string  `json:"section_name"`
    Apperance string  `json:"apperance"`
    Outer string `json:"outer,omitempty"`
    CoorX float32   `json:"coor_x"`
    CoorY float32  `json:"coor_y"`
    CoorZ float32 `json:"coor_z"`
    EquipId int  `json:"-"`
    PropertyId int   `json:"-"`
    ItemNum int  `json:"-"`
    EquipName string  `json:"-"`
    PropertyName string `json:"-"`
    Money   int `json:"-"`
    Item1Id int   `json:"-"`
    Item1Num int  `json:"-"`
    Item2Id int  `json:"-"`
    Item2Num int   `json:"-"`
    Item3Id int    `json:"-"`
    Item3Num int   `json:"-"`
    Item4Id int  `json:"-"`
    Item4Num int   `json:"-"`
    Item5Id int   `json:"-"`
    Item5Num int   `json:"-"`
    Item6Id int   `json:"-"`
    Item6Num int   `json:"-"`
    ItemMoney int   `json:"-"`
    ItemAll string `json:"item_all,omitempty"`
}

func pickup(c *gin.Context) {
    var (
        pItems []PItem
        ss TwoParam     // scene and section
        err error
        id int
    )
    if err = c.ShouldBindJSON(&ss); err != nil {    
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    k:=fmt.Sprintf("pal4:pickup:%s:%s",ss.Class,ss.Type)
    val,err:=client.Get(k).Result()
    if err==nil {
		json.Unmarshal([]byte(val),&pItems)
		c.IndentedJSON(http.StatusOK,pItems)
        return
    }
    chkSql:=`select id from Scene where scene=? and section=?`
    err=Db.QueryRow(chkSql,ss.Class,ss.Type).Scan(&id)
    if err!=nil {
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    sql:=`
        select item_id,model,texture,coor_x,coor_y,coor_z,equip_id,property_id,item_num,money,item1id,
        item1num,item2id,item2num,item3id,item3num,item4id,item4num,item5id,item5num,item6id,item6num,item_money 
        from SceneItem where scene=? and section=?
    `
    rows,_ := Db.Query(sql,ss.Class,ss.Type)
    for rows.Next() {
        it:=PItem{}
        rows.Scan(
            &it.ItemId,&it.Model,&it.Texture,
            &it.CoorX,&it.CoorY,&it.CoorZ,&it.EquipId,&it.PropertyId,
            &it.ItemNum,&it.Money,&it.Item1Id,&it.Item1Num,
            &it.Item2Id,&it.Item2Num,&it.Item3Id,&it.Item3Num,
            &it.Item4Id,&it.Item4Num,&it.Item5Id,&it.Item5Num,
            &it.Item6Id,&it.Item6Num,&it.ItemMoney,
        )
        it.SectionName=getSectionName(ss.Class,ss.Type)
        it.Apperance=getApperance(it.Model)
        switch {
        case it.EquipId!=0:
                it.EquipName=getName("Equip",it.EquipId)
                if it.ItemNum>1 {
                    it.ItemAll+=fmt.Sprintf("%s*%d ",it.EquipName,it.ItemNum)
                } else if it.ItemNum==1 {
                    it.ItemAll+=fmt.Sprintf("%s ",it.EquipName)
                }
        case it.PropertyId!=0:
                it.PropertyName=getName("Property",it.PropertyId)
                if it.ItemNum>1 {
                    it.ItemAll+=fmt.Sprintf("%s*%d ",it.PropertyName,it.ItemNum)
                } else if it.ItemNum==1 {
                    it.ItemAll+=fmt.Sprintf("%s ",it.PropertyName)
                }
        case it.Money!=0:
                it.ItemAll+=fmt.Sprintf("金钱*%d ",it.Money)
        }
        //只有在非单独的装备、道具或是钱时，并且物品数量不为0时，才可能是宝箱类
        if it.EquipId==0 && it.PropertyId==0 && it.Money==0 && it.ItemNum==0 {
            v:=reflect.ValueOf(it)
            for i:=1;i<=6;i++ {
                f:=fmt.Sprintf("Item%dId",i)
                if fv:=v.FieldByName(f); fv.Int()>0 {
                    //先查道具，再查装备
                    itName:=getName("Property",int(fv.Int()))
                    if itName=="" {
                        itName=getName("Equip",int(fv.Int()))
                    }
                    if itName!="" {
                        fn:=fmt.Sprintf("Item%dNum",i)
                        fvn:=v.FieldByName(fn).Int()
                        if fvn>1 {
                            it.ItemAll+=fmt.Sprintf("%s*%d ",itName,fvn)
                        } else if fvn==1 {
                            it.ItemAll+=fmt.Sprintf("%s ",itName)
                        }
                    }
                }
            }
        }
        it.Outer=getOuter(ss.Class,ss.Type)
        if it.SectionName!="" {
            pItems = append(pItems,it)
        }
    }
    rows.Close()
	s,err:=json.Marshal(pItems)
	if err!=nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		client.Set(k, string(s), 36*time.Hour)
		c.IndentedJSON(http.StatusOK,pItems)
	}
}