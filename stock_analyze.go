package main

import (
//	"fmt"
//	"reflect"
//	"net/http"
//	"time"
	"sort"
)
type Clear struct {
	Stats
	AvgDailyProfit float32
}
type Clears struct {
	Profits float32
	ClearList []Clear
}
//清仓类
func getClearance()(clears Clears) {
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
		clears.Profits+=clear.Amount
		clears.ClearList=append(clears.ClearList,clear)
	}
	sort.Sort(ByProfitReverse(clears.ClearList))
	return
}
type Hold struct {
	Stats
	Balance int
	AvgCost float32
}
type Holds struct {
	Costs float32
	HoldList []Hold
}
//持仓类
func getPosition()(holds Holds) {
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
		holds.Costs+=hold.Amount
		holds.HoldList=append(holds.HoldList,hold)
	}
	sort.Sort(ByCost(holds.HoldList))
	return
}
