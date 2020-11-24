package main

import (
    "encoding/json"
    "fmt"
    "github.com/gin-gonic/gin"
    "html/template"
    "net/http"
    "reflect"
    "strings"
    "time"
    _ "github.com/go-sql-driver/mysql"
)
type UniqueFighter struct {
	Id int  `json:"id"`
	Name string `json:"name"`
}
type FighterType struct {
    Type   string  `json:"type" binding:"required"`
}
type FighterName struct {
	FighterType
	Name string `json:"name" binding:"required"`
}
type FighterLevel struct {
    //公共部分
    Lvl int `json:"level"`
    Need int `json:"need,omitempty"`
    //内功专属
    Func int `json:"function,omitempty"`
    Param1 int `json:"-"`
    Param2 int  `json:"-"`
    Param3 int `json:"-"`
    Param4 int `json:"-"`
    Params template.HTML `json:"parameters,omitempty"`
    //外功专属
    AdditionRate int `json:"additional_rate,omitempty"`
    BaseEffect int `json:"base_effect,omitempty"`
    WorkRange int `json:"work_range,omitempty"`
    Times int `json:"times,omitempty"`
    MPConsume int `json:"mp_consume,omitempty"`
    AnimId int `json:"animation_id,omitempty"`
    RangeName string `json:"range_name,omitempty"`
}
type Fighters struct {
    Id int   `json:"id"` //武功ID
    Name string `json:"name"`
    Description string  `json:"description"` //名称、描述
    //IsNei,IsWai bool //内外功标识
    FighterLevels []FighterLevel `json:"pratice_levels"`
}
//neigong--->内功心法, lover--->情侣合技, rage--->怒技, common--->普通招式, combo--->组合技
func fighter(c *gin.Context) {
	var ft FighterType
    var fn FighterName
    var fighters Fighters
	var sql,k,table string
	var err error
    var ufs []UniqueFighter
    logger.Printf("c操作前分析: %v",*c)
    err = c.ShouldBindJSON(&fn)
    logger.Printf("c分析: %v",*c)
    if err == nil { //如果只给出了类型，就属于罗列武功列表   
        k=fmt.Sprintf("crh:fl:%s:%s",fn.Type,fn.Name)
        val,err:=client.Get(k).Result()
        if err==nil {
            json.Unmarshal([]byte(val),&fighters)
            c.IndentedJSON(http.StatusOK,fighters)
            return
        }
        fighters.Name=fn.Name
        if fn.Type=="neigong" {
            table="SPFighterTable1"
        } else {
            table="SPFighterTable2"
        }
        fighters.Id=getId(table,fn.Name)
        fighters.Description=getName("SPFighterHelp",fighters.Id)
        if table=="SPFighterTable1" {
            //fighters.IsNei=true
            sql=`
                select lvl,need,func,param1,param2,param3,param4
                from SPFighterTable1 where name=?
            `
            rows,_ := Db.Query(sql,fn.Name)
            for rows.Next() {
                lvl:=FighterLevel{}
                rows.Scan(&lvl.Lvl,&lvl.Need,&lvl.Func,
                    &lvl.Param1,&lvl.Param2,&lvl.Param3,&lvl.Param4,
                )
                p:=""
                v:=reflect.ValueOf(lvl)
                for i:=1;i<=4;i++ {
                    f:=fmt.Sprintf("Param%d",i)
                    if fv:=v.FieldByName(f); fv.Int()>0 {
                        p+=fmt.Sprintf("参数%d: %d  ",i,fv.Int())
                    }
                }
                p=strings.TrimSuffix(p," ")
                lvl.Params=template.HTML(p)
                fighters.FighterLevels = append(fighters.FighterLevels, lvl)
            }
        } else if table=="SPFighterTable2" {
            //fighters.IsWai=true
            sql=`
            select lvl,need,addition_rate,base_effect,work_range,times,mp_consume,anim_id
            from SPFighterTable2 where name=?
            `
            rows,_ := Db.Query(sql,fn.Name)
            for rows.Next() {
                lvl:=FighterLevel{}
                rows.Scan(
                    &lvl.Lvl,&lvl.Need,&lvl.AdditionRate,&lvl.BaseEffect,&lvl.WorkRange,
                    &lvl.Times,&lvl.MPConsume,&lvl.AnimId,
                )
                lvl.RangeName=getName("SPRange",lvl.WorkRange)
                fighters.FighterLevels = append(fighters.FighterLevels, lvl)
            }
        }
        s,err:=json.Marshal(fighters)
        if err!=nil {
            c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        } else {
            client.Set(k, string(s), 12*time.Hour)
            c.IndentedJSON(http.StatusOK,fighters)
        return
        }    
        
    } 
    err = c.ShouldBindJSON(&ft)
    logger.Printf("c二度分析: %v",*c)
    if err == nil { //如果同时给出了名称，就属于要得到该武功修炼明细了
        logger.Println("进入了武功类型处理")
        k=fmt.Sprintf("crh:ft:%s",ft.Type)
		val,err:=client.Get(k).Result()
		if err==nil {
			json.Unmarshal([]byte(val),&ufs)
            c.IndentedJSON(http.StatusOK,ufs)
			return
		}
		if ft.Type=="neigong" {
			sql=`select distinct id, name from SPFighterTable1 where name!="保留"`
		} else {
			sql=`select distinct id, name from SPFighterTable2 where name not regexp "組合|保留"`
			switch ft.Type {
				case "lover":
					sql+=` and id in (select id from SPFighterHelp where name regexp "合技")`
				case "rage":
					sql+=` and name regexp "之怒$"`
				case "common":
					sql+=` and name in (select name from SPFighterTable2 group by name having count(lvl)>=9) `
				case "combo":
					sql+=` and name regexp "[0-9]技$"`
			}
		}
		rows,_ := Db.Query(sql)
		for rows.Next() {
			var uf UniqueFighter
			rows.Scan(&uf)
			ufs = append(ufs,uf)
		}
        rows.Close()
        logger.Printf("武功列表：%v\n",ufs)
        s,err:=json.Marshal(ufs)
		if err!=nil {
            logger.Printf("类型序列化出了问题：%v\n",err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else {
            client.Set(k, string(s), 12*time.Hour)
            c.IndentedJSON(http.StatusOK,ufs)
			return
		}
	}
	c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})  //否则，按出错处理
}