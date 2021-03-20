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
    CheckIn      time.Time `json:"checkin"  binding:"required"`
    CheckOut     time.Time `json:"checkout"  binding:"required"`
    Comments    string `json:"comments"`
}
type Duration struct {
	Name string  `json:"name" binding:"required"` 
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
            c.IndentedJSON(http.StatusOK,attn)
            return
        }
        err = Db.QueryRow("select checkin,checkout,comments from attendance order by checkin desc limit 1").Scan(&attn.CheckIn, &attn.CheckOut, &attn.Comments)
        if err!=nil {
            c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		    return
        }
        s,err:=json.Marshal(attn)
        if err!=nil {
            c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		    return
        } else {
            client.Set("attendance:latest", string(s), 2*time.Hour)
            c.IndentedJSON(http.StatusOK,attn)
            return
        }
    //近一周与近一月的
    case "last-week", "last-month":
        var attns []Attendance
        k:=fmt.Sprintf("attendance:%s",d.Name)
        val,err:=client.Get(k).Result()
        if err==nil {
            json.Unmarshal([]byte(val),&attns)
            c.IndentedJSON(http.StatusOK,attns)
            return
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
        s,err:=json.Marshal(attns)
        if err!=nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		    return
        } else {
            client.Set(k, string(s), 2*time.Hour)
            c.IndentedJSON(http.StatusOK,attns)
        }    
    }
}
//统计数据
func stat(c *gin.Context){
    var d Duration
    r:=make(map[string]float64)
    var sql string
    if err := c.ShouldBindJSON(&d); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    //先从redis取，如果有，就直接返回
    dk:=fmt.Sprintf("attendance:%s-stats",d.Name)
    val,err:=client.Get(dk).Result()
    if err==nil {
        json.Unmarshal([]byte(val),&r)
        c.IndentedJSON(http.StatusOK,r)
        return
    }
    //redis里没有，从数据库里获得
    switch d.Name {
    case "week", "month":
        sql=fmt.Sprintf(`SELECT year(checkin),%s(checkin),
            sum(hour(timediff(checkout,checkin))+minute(timediff(checkout,checkin))/60)
            FROM attendance group by year(checkin),%[1]s(checkin)
        `,d.Name)
    case "year":
        sql=`SELECT year(checkin),
            sum(hour(timediff(checkout,checkin))+minute(timediff(checkout,checkin))/60)
            FROM attendance group by year(checkin)
        `
    }
    rows,_ := Db.Query(sql)
    for rows.Next() {
        //月统计映射表，键为年月组合，值为合计的小时数
        //monthHours['2019年11月']=217.5
        var year,morw int
        var hours float64
        var k string
        switch d.Name {
        case "week":
            rows.Scan(&year,&morw,&hours)
            k=fmt.Sprintf("%d年%02d周",year,morw)
        case  "month":
            rows.Scan(&year,&morw,&hours)
            k=fmt.Sprintf("%d年%02d月",year,morw)
        case "year":
            rows.Scan(&year,&hours)
            k=fmt.Sprintf("%d年",year)
        }
        r[k]=hours
    }
    rows.Close()
    //然后再写入redis
    s,err:=json.Marshal(r)
    if err!=nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    } else {
        client.Set(dk,string(s),2*time.Hour)
        c.IndentedJSON(http.StatusOK,r)
    } 
}
//创建新记录，返回数据为报错信息字符串，如果成功，则返回成功信息
func add(c *gin.Context) {
    //先设置一个考勤记录的变形体，做为增加记录时的临时变量
    type AttnV struct {
        CheckIn      string `json:"checkin"  binding:"required"`
        CheckOut     string `json:"checkout"  binding:"required"`
        Comments    string `json:"comments"`
    }
    var attnV AttnV
    var attn Attendance 
    if err := c.ShouldBindJSON(&attnV); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    tfa := "2006-01-02 15:04"
    local, _ := time.LoadLocation("Local")
    attn.CheckIn,_ = time.ParseInLocation(tfa,attnV.CheckIn,local)
    attn.CheckOut,_ = time.ParseInLocation(tfa,attnV.CheckOut,local)
    attn.Comments=attnV.Comments
    //获得日期的字符串形式，固定为2020-02-20这样的形式
    tfd := "2006-01-02"
    checkday:=attn.CheckIn.Format(tfd)
    nk:=fmt.Sprintf("attendance:%s", checkday)
    val,err:=client.Get(nk).Result()
    if err==nil && val=="true"{
        c.JSON(http.StatusBadRequest, gin.H{"error": "依据缓存，此日考勤已录入！"})
        return
    }
    if attn.CheckIn.After(time.Now())  || attn.CheckOut.After(time.Now()) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "上下班时间，均不能超过当前时间！"})
        return
    }
    if attn.CheckIn.Year()!=attn.CheckOut.Year() || attn.CheckIn.YearDay()!=attn.CheckOut.YearDay() {
        c.JSON(http.StatusBadRequest, gin.H{"error": "考勤禁止跨日！"})
        return
    }
    //设置正常的考勤时间，此项不同公司有所不同，要注意！
    if (attn.CheckIn.Hour()>9 || attn.CheckOut.Hour()<18 || attn.CheckOut.Hour()==18 && attn.CheckOut.Minute()<30) && attn.Comments=="" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "上班时间为9：00～18：30，迟到、早退，请假等特殊情况，必须填写备注说明原因！"})
        return
    }
    var cnt int
    var errInfo string
    //检查是否已有当日记录
    checksql:="select count(id) from attendance where  year(checkin)=? and dayofyear(checkin)=?"
    err = Db.QueryRow(checksql,attn.CheckIn.Year(),attn.CheckIn.YearDay()).Scan(&cnt)
    if err!=nil {
        errInfo=fmt.Sprintf("查询当日记录出错：%v",err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": errInfo})
        return
    } else if cnt!=0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "已经有当日记录！"})
        return
    }
    //正式插入
    statement := "insert into attendance(checkin, checkout, comments) value (?,?,?)"
    _, err = Db.Exec(statement,attn.CheckIn,attn.CheckOut,attn.Comments)
    if err != nil {
        errInfo=fmt.Sprintf("无法创建新记录，错误：%v",err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": errInfo})
        return
    } else {
        client.Set(nk,"true",2*time.Hour)
        c.JSON(http.StatusOK, gin.H{"result": "恭喜！成功增加新考勤记录"})
    }
}