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
//将类别提取出后的结构组合，适用于自动依据类别填充模板
type Equip struct {
    Id int                 `json:"id"`
	Name string         `json:"name"`
	Description string  `json:"description"`
	Attribute     string   `json:"attribute,omitempty"`
	Model string  `json:"model,omitempty"`
	Texture     string   `json:"texture,omitempty"`
	LvLmt int   `json:"level_limit,omitempty"`
	Price int   `json:"price,omitempty"`
	LingCap          int        `json:"-"`
	Potential int   `json:"-"`
	MaxHP int   `json:"-"`
	AdditionalRage int `json:"-"`
	MaxMP int   `json:"-"`
	Physical int  `json:"-"`
	Toughness int  `json:"-"`
	Speed int  `json:"-"`
	Lucky int   `json:"-"`
	Will int `json:"-"`
	Water int  `json:"-"`
	Fire int  `json:"-"`
	Thunder int  `json:"-"`
	Air int  `json:"-"`
	Earth int         `json:"-"`
	WaterAdditional int  `json:"-"`
	FireAdditional int  `json:"-"`
	ThunderAdditional int  `json:"-"`
	AirAdditional int   `json:"-"`
	EarthAdditional int `json:"-"`
	PhysicalExtract float32   `json:"-"`
	WaterExtract float32  `json:"-"`
	FireExtract float32   `json:"-"`
	ThunderExtract float32   `json:"-"`
	AirExtract float32   `json:"-"`
	EarthExtract float32   `json:"-"`
	PhysicalReact float32   `json:"-"`
	WaterReact float32   `json:"-"`
	FireReact float32  `json:"-"`
	ThunderReact float32  `json:"-"`
	AirReact float32  `json:"-"`
	EarthReact float32 `json:"-"`
	AdditionalCritical float32   `json:"-"`
	FendOff float32   `json:"-"`
	AdditionalHitting float32   `json:"-"`
    Effect1 string   `json:"effect1,omitempty"`
    BuyScene string    `json:"buy_scene,omitempty"`
}

//依据装备类型的中文名，获得该类物品的全部属性
func equipment(c *gin.Context) {
    var (
        equips []Equip
        et OneParam
        err error
    )
	if err = c.ShouldBindJSON(&et); err != nil {    
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    k:=fmt.Sprintf("pal4:equip:%s",et.Type)
    val,err:=client.Get(k).Result()
    if err==nil {
		json.Unmarshal([]byte(val),&equips)
		c.IndentedJSON(http.StatusOK,equips)
        return
    }
    equipAttribute:=[][]string{
        {"MaxHP","精上限"},{"AdditionalRage","气"},{"MaxMP","神上限"},
        {"Physical","武"},{"Toughness","防"},{"Speed","速"},{"Lucky","运"},{"Will","灵"},
        {"Water","水"},{"Fire","火"},{"Thunder","雷"},{"Air","风"},{"Earth","土"},
        {"WaterAdditional","水伤害"},{"FireAdditional","火伤害"},{"ThunderAdditional","雷伤害"},
        {"AirAdditional","风伤害"},{"EarthAdditional","土伤害"},{"PhysicalExtract","物理吸收"},
        {"WaterExtract","水吸收"},{"FireExtract","火吸收"},{"ThunderExtract","雷吸收"},
        {"AirExtract","风吸收"},{"EarthExtract","土吸收"},{"PhysicalReact","物理反弹"},
        {"WaterReact","水反弹"},{"FireReact","火反弹"},{"ThunderReact","雷反弹"},
        {"AirReact","风反弹"},{"EarthReact","土反弹"},
        {"AdditionalCritical","暴击"},{"FendOff","格挡"},{"AdditionalHitting","命中"},
        {"LingCap","灵蕴"},{"Potential","潜力"},
    }
    typeId:=getId("EquipType",et.Type)
    //利用类型ID取原始数据
    equipSql:=`
        select id,name,description,model,texture,tianhe_lv_lmt|lingsha_lv_lmt|mengli_lv_lmt|ziying_lv_lmt lvl_lmt,price,
        ling_capacity,forge_potential,max_hp,additional_rage,max_mp,physical,toughness,speed,lucky,will,water,fire,
        thunder,air,earth,water_additional,fire_additional,thunder_additional,air_additional,earth_additional,
        physical_extract,water_extract,fire_extract,thunder_extract,air_extract,earth_extract,physical_react,
        water_react,fire_react,thunder_react,air_react,earth_react,additional_critical,fend_off,additional_hitting,ef1 
        from Equip where type=?
    `
    rows,_ := Db.Query(equipSql,typeId)
    for rows.Next() {
        equip := Equip{}
        rows.Scan(
            &equip.Id,&equip.Name,&equip.Description,&equip.Model,&equip.Texture,&equip.LvLmt,&equip.Price,
            &equip.LingCap,&equip.Potential,&equip.MaxHP,&equip.AdditionalRage,&equip.MaxMP,&equip.Physical,
            &equip.Toughness,&equip.Speed,&equip.Lucky,&equip.Will,&equip.Water,&equip.Fire,&equip.Thunder,
            &equip.Air,&equip.Earth,&equip.WaterAdditional,&equip.FireAdditional,&equip.ThunderAdditional,
            &equip.AirAdditional,&equip.EarthAdditional,&equip.PhysicalExtract,&equip.WaterExtract,&equip.FireExtract,
            &equip.ThunderExtract,&equip.AirExtract,&equip.EarthExtract,&equip.PhysicalReact,&equip.WaterReact,
            &equip.FireReact,&equip.ThunderReact,&equip.AirReact,&equip.EarthReact,&equip.AdditionalCritical,
            &equip.FendOff,&equip.AdditionalHitting,&equip.Effect1,
        )
        //利用反射合并属性
        v:=reflect.ValueOf(equip)
        for _,f:= range equipAttribute {
            fv:=v.FieldByName(f[0])
            t:=fv.Type()
            if t.String()=="int" && fv.Int()>0 {
                if f[1]=="灵蕴" || f[1]=="潜力" {
                    equip.Attribute+=fmt.Sprintf("%s:%d ",f[1],fv.Int())
                } else {
                    equip.Attribute+=fmt.Sprintf("%s+%d ",f[1],fv.Int())
                }
            } else if t.String()=="float32" && fv.Float()>0 {
                equip.Attribute+=fmt.Sprintf("%s+%s%% ",f[1],perDisp(float32(fv.Float()*100)))
            }
        }
        equip.Attribute=strings.TrimSuffix(equip.Attribute," ")
        buyScene:=getBuyScene(equip.Id)
        if buyScene=="" {
            equip.BuyScene="无法直接购买，可能配方表中有制造方法"
        } else {
            equip.BuyScene=buyScene
        }        
        equips = append(equips, equip)
    }
    rows.Close()
    if len(equips)==0 {
        c.IndentedJSON(http.StatusNotFound, gin.H{"error": "参数错误，什么也查不到！"})
        return
    }
	s,err:=json.Marshal(equips)
	if err!=nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		client.Set(k, string(s), 36*time.Hour)
		c.IndentedJSON(http.StatusOK,equips)
	}
}
