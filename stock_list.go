package main
import (
	"net/http"
	"github.com/gin-gonic/gin"
	"time"
	"fmt"
	"encoding/json"
)
//列出指定代码的股票的全部交易记录，按时间逆序
func list(c *gin.Context) {
	var (
		stocks []Stock
		code OneParam
		err error
	)
	if err = c.ShouldBindJSON(&code); err != nil {    
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    k:=fmt.Sprintf("stock:list:%s",code.Type)
    val,err:=client.Get(k).Result()
    if err==nil {
		json.Unmarshal([]byte(val),&stocks)
		c.IndentedJSON(http.StatusOK,stocks)
        return
    }
	sql := `
		select date,code,name,operation,volume,balance,
		price,turnover,amount,brokerage,stamps,transfer_fee from stock where code=? order by date desc
    `
	rows, _ := Db.Query(sql,code.Type)
	for rows.Next() {
		s := Stock{}
		rows.Scan(
			&s.Date,&s.Code,&s.Name,&s.Operation,&s.Volume,&s.Balance,
			&s.Price,&s.Turnover,&s.Amount,&s.Brokerage,&s.Stamps,&s.TransferFee,
		)
		stocks = append(stocks, s)
	}
	rows.Close()
	if len(stocks)==0 {
        c.IndentedJSON(http.StatusNotFound, gin.H{"error": "参数错误，什么也查不到！"})
        return
    }
	s,err:=json.Marshal(stocks)
	if err!=nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		client.Set(k, string(s), 36*time.Hour)
		c.IndentedJSON(http.StatusOK,stocks)
	}
}
//持仓股票的最近一次交易记录，出入各取一条
func holdLastDeal(c *gin.Context) {
	var (
        stocks []Stock
		err error
	)
	k:=fmt.Sprintf("stock:hld")
    val,err:=client.Get(k).Result()
    if err==nil {
		json.Unmarshal([]byte(val),&stocks)
		c.IndentedJSON(http.StatusOK,stocks)
        return
    }
	codeMap:=getNameMap("group by code having sum(volume)!=0")
	codes:=[]string{}
	for k, _ := range codeMap {
		codes=append(codes,k)
	}
	inSql:=`
		select date,code,name,operation,volume,balance,
		price,turnover,amount,brokerage,stamps,transfer_fee from stock 
		where code=?  and  operation in ("申购中签","证券买入","红股入账","股息红利税补") order by date desc limit 1
	`
	outSql:=`
		select date,code,name,operation,volume,balance,
		price,turnover,amount,brokerage,stamps,transfer_fee from stock 
		where code=?  and  operation in ("证券卖出","股息入账") order by date desc limit 1
	`
	for _,c:=range codes {
		s:=Stock{}
		Db.QueryRow(inSql,c).Scan(
			&s.Date,&s.Code,&s.Name,&s.Operation,&s.Volume,&s.Balance,
			&s.Price,&s.Turnover,&s.Amount,&s.Brokerage,&s.Stamps,&s.TransferFee,
		)
		stocks=append(stocks,s)
		//买入类可以不检查，但卖出类一定要检查，原因是可能没有卖出类操作
		err=Db.QueryRow(outSql,c).Scan(
			&s.Date,&s.Code,&s.Name,&s.Operation,&s.Volume,&s.Balance,
			&s.Price,&s.Turnover,&s.Amount,&s.Brokerage,&s.Stamps,&s.TransferFee,
		)
		if err==nil {
			stocks=append(stocks,s)
		}
	}
	if len(stocks)==0 {
        c.IndentedJSON(http.StatusNotFound, gin.H{"error": "参数错误，什么也查不到！"})
        return
    }
	s,err:=json.Marshal(stocks)
	if err!=nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		client.Set(k, string(s), 36*time.Hour)
		c.IndentedJSON(http.StatusOK,stocks)
	}
}