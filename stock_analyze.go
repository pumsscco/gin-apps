package main

import (
	"github.com/gin-gonic/gin"
	"sort"
	"fmt"
	"encoding/json"
	"net/http"
	"time"
)
type Clear struct {
	Stats
	AvgDailyProfit float32  `json:"average_daily_profit"`
}
//清仓类
func clearance(c *gin.Context) {
	var (
        clears []Clear
		err error
	)
	k:=fmt.Sprintf("stock:clearance")
    val,err:=client.Get(k).Result()
    if err==nil {
		json.Unmarshal([]byte(val),&clears)
		c.IndentedJSON(http.StatusOK,clears)
        return
    }
	clearStocks:=getNameMap("group by code having sum(volume)=0")
	for k,v:=range clearStocks {
		clear:=Clear{}
		//第一步：先拿到净利润、首末交易日期、持股天数、日均利润
		sql:=`
			select sum(amount),min(date),max(date),datediff(max(date),min(date)),
			sum(amount)/datediff(max(date),min(date)) from stock where code=?
		`
		Db.QueryRow(sql,k).Scan(
			&clear.Amount,&clear.FirstDealDay,&clear.LastDealDay,&clear.HoldDays,&clear.AvgDailyProfit,
		)
		//第二步：获得买卖次数，并计算周买卖频率
		sql=`select count(id) from stock where code=? and operation in ('申购中签','证券买入','证券卖出')`
		Db.QueryRow(sql,k).Scan(&clear.TransactionCount)
		clear.TransactionFreq=float32(clear.TransactionCount)/float32(clear.HoldDays)*7
		//第三步，取出最高持仓量与相应日期
		sql=`select date,balance from stock where code=? and balance=(select max(balance) from stock where code=?)`
		Db.QueryRow(sql,k,k).Scan(&clear.MaxBalanceDay,&clear.MaxBalance)
		clear.Code=k
		clear.Name=v
		clears=append(clears,clear)
	}
	sort.Sort(ByProfitReverse(clears))
	if len(clears)==0 {
        c.IndentedJSON(http.StatusNotFound, gin.H{"error": "参数错误，什么也查不到！"})
        return
    }
	s,err:=json.Marshal(clears)
	if err!=nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		client.Set(k, string(s), 36*time.Hour)
		c.IndentedJSON(http.StatusOK,clears)
	}
}
type Hold struct {
	Stats
	Balance int `json:"Balance"`
	AvgCost float32 `json:"average_cost"`
}
//持仓类
func position(c *gin.Context) {
	var (
        holds []Hold
		err error
	)
	k:=fmt.Sprintf("stock:position")
    val,err:=client.Get(k).Result()
    if err==nil {
		json.Unmarshal([]byte(val),&holds)
		c.IndentedJSON(http.StatusOK,holds)
        return
    }
	holdStocks:=getNameMap("group by code having sum(volume)!=0")
	for k,v:=range holdStocks {
		hold:=Hold{}
		//第一步：先拿到总成本、首次及最近交易日期、持股天数、持仓量、计算平均成本
		sql:=`
			select sum(amount),min(date),max(date),datediff(curdate(),min(date)),
			sum(volume) from stock where code=?
		`
		Db.QueryRow(sql,k).Scan(
			&hold.Amount,&hold.FirstDealDay,&hold.LastDealDay,&hold.HoldDays,&hold.Balance,
		)
		hold.AvgCost=float32(hold.Amount)/float32(hold.Balance)
		//第二步：获得买卖次数，并计算周买卖频率
		sql=`select count(id) from stock where code=? and operation in ('申购中签','证券买入','证券卖出')`
		Db.QueryRow(sql,k).Scan(&hold.TransactionCount)
		
		hold.TransactionFreq=float32(hold.TransactionCount)/float32(hold.HoldDays)*7
		//第三步，取出最高持仓量与相应日期
		sql=`select date,balance from stock where code=? and balance=(select max(balance) from stock where code=?)`
		Db.QueryRow(sql,k,k).Scan(&hold.MaxBalanceDay,&hold.MaxBalance)
		hold.Code=k
		hold.Name=v
		holds=append(holds,hold)
	}
	sort.Sort(ByCost(holds))
	if len(holds)==0 {
        c.IndentedJSON(http.StatusNotFound, gin.H{"error": "参数错误，什么也查不到！"})
        return
    }
	s,err:=json.Marshal(holds)
	if err!=nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		client.Set(k, string(s), 36*time.Hour)
		c.IndentedJSON(http.StatusOK,holds)
	}
}
