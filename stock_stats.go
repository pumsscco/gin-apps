package main
import (
	"net/http"
	"github.com/gin-gonic/gin"
	"time"
	"fmt"
	"encoding/json"
)
//一个参数，要么清仓，要么持仓，清仓统计总收益，持仓统计总成本
func statistic(c *gin.Context) {
	var (
		amount float32
		st OneParam
		err error
		cond string
	)
	if err = c.ShouldBindJSON(&st); err != nil {    
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    k:=fmt.Sprintf("stock:statistic:%s",st.Type)
    val,err:=client.Get(k).Result()
    if err==nil {
		json.Unmarshal([]byte(val),&amount)
		c.IndentedJSON(http.StatusOK,amount)
        return
	}
	switch st.Type {
	case "清仓收益":
		cond=`group by code having sum(volume)=0`
	case "持仓成本":
		cond=`group by code having sum(volume)!=0`
	default:
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "参数错误，禁止统计！"})
        return
	} 
	sql := fmt.Sprintf(`select sum(sum_amt)  whole from (select code, sum(amount)  sum_amt from stock %s) t`,cond)
	err = Db.QueryRow(sql).Scan(&amount)
	if err!=nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	s,err:=json.Marshal(amount)
	if err!=nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		client.Set(k, string(s), 36*time.Hour)
		c.IndentedJSON(http.StatusOK,amount)
	}
}