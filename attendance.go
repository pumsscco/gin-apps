package main

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "time"
    "encoding/json"
    "strings"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
)
//考勤表结构
type Attendance struct {
    //Id        int `json:"id"`
    CheckIn      time.Time `json:"checkin"`
    CheckOut     time.Time `json:"checkout"`
    Comments    string `json:"comments"`
}
type Duration struct {
	Name string  `json: "name"` 
}
//抓记录
func rec(c *gin.Context) {
    //先查一下输入的json有没有格式错误
    var d Duration
    if err := c.ShouldBindJSON(&d); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
    }
    //然后依据具体的post条件，来分析结果
    switch d.Name {
    //先处理只取最近一条的情况
    case "latest":
        var attn Attendance
        val,err:=client.Get("attendance:latest").Result()
        if err==nil {
            json.Unmarshal([]byte(val),&attn)
            a_info:=gin.H{
                "check_in": attn.CheckIn,
                "check_out": attn.CheckOut,
                "comments": attn.Comments,
            }
            c.IndentedJSON(http.StatusOK,a_info)
            return
        }
        err = Db.QueryRow("select checkin,checkout,comments from attendance order by checkin desc limit 1").Scan(&attn.CheckIn, &attn.CheckOut, &attn.Comments)
        if err!=nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		    return
        }
        as,err:=json.Marshal(attn)
        if err!=nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		    return
        } else {
            client.Set("attendance:latest", string(as), 2*time.Hour)
        }
        a_info:=gin.H{
            "check_in": attn.CheckIn,
            "check_out": attn.CheckOut,
            "comments": attn.Comments,
        }
        c.IndentedJSON(http.StatusOK,a_info)
        return
    //近一周与近一月的
    case "last-week", "last-month":
        var attns []Attendance
        val,err:=client.Get(fmt.Sprintf("attendance:%s",d.Name)).Result()
        if err==nil {
            json.Unmarshal([]byte(val),&attns)
            var a_infos []gin.H
            for _, v:= range attns {
                a_info:=gin.H{
                    "check_in": v.CheckIn,
                    "check_out": v.CheckOut,
                    "comments": v.Comments,
                }
                a_infos=append(a_infos,a_info)
            }
            c.IndentedJSON(http.StatusOK,a_infos)
        }
        dur:=strings.Split(d.Name,"-")[1]
        sql:=fmt.Sprintf("select checkin,checkout,comments from attendance where date_add(checkin,interval 1 %s)>now()",dur)
        rows,_ := Db.Query(sql)
        for rows.Next() {
            attn:=Attendance{}
            rows.Scan(&attn.CheckIn, &attn.CheckOut, &attn.Comments)
            attns=append(attns,attn)
        }
        rows.Close()
        as,err:=json.Marshal(attns)
        if err!=nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		    return
        } else {
            client.Set(fmt.Sprintf("attendance:%s",d.Name), string(as), 2*time.Hour)
        }
        var a_infos []gin.H
        for _, v:= range attns {
            a_info:=gin.H{
                "check_in": v.CheckIn,
                "check_out": v.CheckOut,
                "comments": v.Comments,
            }
            a_infos=append(a_infos,a_info)
        }
        c.IndentedJSON(http.StatusOK,a_infos)    
    }
}
//月度统计数据
/*func stat() (monthHours map[string]float64) {
    //先从redis取，如果有，就直接返回
    val,err:=client.Get("month-stats").Result()
    if err==nil {
        json.Unmarshal([]byte(val),&monthHours)
        return
    }
    //redis里没有，从数据库里获得
    monthHours=make(map[string]float64)
    mSql:=`SELECT year(checkin),month(checkin),
        sum(hour(timediff(checkout,checkin))+minute(timediff(checkout,checkin))/60)
        FROM attendance group by year(checkin),month(checkin)
    `
    rows,_ := Db.Query(mSql)
    for rows.Next() {
        //月统计映射表，键为年月组合，值为合计的小时数
        //monthHours['2019年11月']=217.5
        var year,month int
        var hours float64
        rows.Scan(&year,&month,&hours)
        yearMonth:=fmt.Sprintf("%d年%d月",year,month)
        monthHours[yearMonth]=hours
    }
    rows.Close()
    //然后再写入redis
    ms,_:=json.Marshal(monthHours)
    client.Set("month-stats",string(ms),2*time.Second)
    return
}
//周统计
func weekHour() (weekHours map[string]float64) {
    //先从redis取，如果有，就直接返回
    val,err:=client.Get("week-stats").Result()
    if err==nil {
        json.Unmarshal([]byte(val),&weekHours)
        return
    }
    //redis里没有，从数据库里获得
    weekHours=make(map[string]float64)
    wSql:=`SELECT year(checkin),week(checkin),
        sum(hour(timediff(checkout,checkin))+minute(timediff(checkout,checkin))/60)
        FROM attendance group by year(checkin),week(checkin)
    `
    rows,_ := Db.Query(wSql)
    for rows.Next() {
        //周统计映射表，键为年周组合，值为合计的小时数
        //weekHours['2019年17周']=242
        var year,week int
        var hours float64
        rows.Scan(&year,&week,&hours)
        yearWeek:=fmt.Sprintf("%d年%d周",year,week)
        weekHours[yearWeek]=hours
    }
    rows.Close()
    //然后再写入redis
    ws,_:=json.Marshal(weekHours)
    client.Set("week-stats",string(ws),2*time.Second)
    return
}
//年统计
func yearHour() (yearHours map[int]float64) {
    //先从redis取，如果有，就直接返回
    val,err:=client.Get("year-stats").Result()
    if err==nil {
        json.Unmarshal([]byte(val),&yearHours)
        return
    }
    //redis里没有，从数据库里获得
    yearHours=make(map[int]float64)
    wSql:=`SELECT year(checkin),
        sum(hour(timediff(checkout,checkin))+minute(timediff(checkout,checkin))/60)
        FROM attendance group by year(checkin)
    `
    rows,_ := Db.Query(wSql)
    for rows.Next() {
        //年统计，键为年，值为合计的小时数
        //yearHours[2019]=242
        var year int
        var hours float64
        rows.Scan(&year,&hours)
        yearHours[year]=hours
    }
    rows.Close()
    //然后再写入redis
    ys,_:=json.Marshal(yearHours)
    client.Set("year-stats",string(ys),2*time.Second)
    return
}
//创建新记录，返回数据为报错信息的字符串，如果成功，则返回空串
func (attn Attendance)NewAttn() (errInfo string) {
    //获得日期的字符串形式，固定为2020-02-20这样的形式
    tf := "2006-01-02"
    checkday:=attn.CheckIn.Format(tf)
    val,err:=client.Get(checkday).Result()
    if err==nil {
        if val=="true" {
            errInfo="依据缓存，此日考勤已录入！"
            logger.Print(errInfo)
            return
        }
    }
    if attn.CheckIn.After(time.Now())  || attn.CheckOut.After(time.Now()) {
        errInfo="上下班时间，均不能超过当前时间！"
        logger.Println(errInfo)
        return
    }
    if attn.CheckIn.Year()!=attn.CheckOut.Year() || attn.CheckIn.YearDay()!=attn.CheckOut.YearDay() {
        errInfo="考勤禁止跨日！"
        logger.Println(errInfo)
        return
    }
    if attn.CheckIn.Hour()>9 || attn.CheckOut.Hour()<17 || (attn.CheckOut.Hour()==17 && attn.CheckOut.Minute()<30) {
        if attn.Comments=="" {
            errInfo="上班时间为9：00～17：30，迟到、早退，请假等特殊情况，必须填写备注说明原因！"
            logger.Println()
            return
        }
    }
    var cnt int
    checksql:="select count(id) from attendance where  year(checkin)=? and dayofyear(checkin)=?"
    err = Db.QueryRow(checksql,attn.CheckIn.Year(),attn.CheckIn.YearDay()).Scan(&cnt)
    if err!=nil {
        errInfo=fmt.Sprintf("查询当日记录出错：%v\n",err)
        logger.Print(errInfo)
        return
    } else if cnt!=0 {
        errInfo="已经有当日记录！"
        logger.Print(errInfo)
        return
    }
    statement := "insert into attendance(checkin, checkout, comments) value (?,?,?)"
    _, err = Db.Exec(statement,attn.CheckIn,attn.CheckOut,attn.Comments)
    if err != nil {
        errInfo=fmt.Sprintf("无法创建新记录，错误：%v\n",err)
        logger.Print(errInfo)
        return
    }
    client.Set(checkday,"true",2*time.Second)
    return
}*/