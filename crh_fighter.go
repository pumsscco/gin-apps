package main

import (
    "fmt"
    "reflect"
    "strings"
    "html/template"
    "time"
	"encoding/json"
)
type UniqueFighter struct {
	Id int
	Name string
}
type FighterType struct {
    Type   string  `json: "type" binding:"required"`
}
type FighterName struct {
	FighterType
	Name string `json: "name" binding:"required"`
}
//neigong--->内功心法, lover--->情侣合技, rage--->怒技, common--->普通招式, combo--->组合技
func fighter(c *gin.Context) {
	var ft FighterType
	var fn FighterName
	var sql string
	var err error
	var ufs []UniqueFighter
	var uf_infos []gin.H
	k:=fmt.Sprintf("crh:ft:%s",ft.Type)
	if err = c.ShouldBindJSON(&ft); err == nil { //如果只给出了类型，就属于罗列武功列表       
		val,err:=client.Get(k).Result()
		if err==nil {
			json.Unmarshal([]byte(val),&ufs)
			for _, v:= range ufs {
                uf_info:=gin.H{
                    "id": v.Id,
                    "name": v.Name,
                }
                uf_infos=append(uf_infos,uf_info)
            }
            c.IndentedJSON(http.StatusOK,uf_infos)
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
		as,err:=json.Marshal(ufs)
		client.Set(k, string(as), 12*time.Hour)
		if err!=nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else {
			for _, v:= range ufs {
                uf_info:=gin.H{
                    "id": v.Id,
                    "name": v.Name,
                }
                uf_infos=append(uf_infos,uf_info)
            }
            c.IndentedJSON(http.StatusOK,uf_infos)
			return
		}
    } else if err = c.ShouldBindJSON(&fn); err == nil { //如果同时给出了名称，就属于要得到该武功修炼明细了
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})  //否则，按出错处理
	}   
}

type FighterLevel struct {
    //公共部分
    Lvl,Need int
    //内功专属
    Func,Param1,Param2,Param3,Param4 int 
    Params template.HTML
    //外功专属
    AdditionRate,BaseEffect,WorkRange,Times,MPConsume,AnimId int 
    RangeName string
}
type Fighters struct {
    Id int   //武功ID
    Name,Description string  //名称、描述
    IsNei,IsWai bool //内外功标识
    FighterLevels []FighterLevel
}
//依据武功名称，获得修炼详细数据
func getFighterLevel(table,fighter string) (fighters Fighters)  {
    k:=fmt.Sprintf("fl:%s",fighter)
    val,err:=client.Get(k).Result()
    if err==nil {
        json.Unmarshal([]byte(val),&fighters)
        return
    }
    fighters.Name=fighter
    fighters.Id=getId(table,fighter)
    fighters.Description=getName("SPFighterHelp",fighters.Id)
    sql:=""
    if table=="SPFighterTable1" {
        fighters.IsNei=true
        sql=`
            select lvl,need,func,param1,param2,param3,param4
            from SPFighterTable1 where name=?
        `
        rows,_ := Db.Query(sql,fighter)
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
                    p+=fmt.Sprintf("参数%d:&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;%d&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; ",i,fv.Int())
                }
            }
            p=strings.TrimSuffix(p," ")
            lvl.Params=template.HTML(p)
            fighters.FighterLevels = append(fighters.FighterLevels, lvl)
        }
    } else if table=="SPFighterTable2" {
        fighters.IsWai=true
        sql=`
        select lvl,need,addition_rate,base_effect,work_range,times,mp_consume,anim_id
        from SPFighterTable2 where name=?
        `
        rows,_ := Db.Query(sql,fighter)
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
    as,err:=json.Marshal(fighters)
	client.Set(k, string(as), 12*time.Hour)
	if err!=nil {
		logger.Print(err)
	}
    return
}

