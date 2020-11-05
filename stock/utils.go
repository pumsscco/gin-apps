package main

import (
	"time"
)
//完整的交易记录结构
type Stock struct {
	Id                               int
	Date      time.Time      //日期
	Code, Name, Operation     string //代码、名称、操作
	Volume, Balance     int    //数量、变动后持股数量
	Price, Turnover, Amount,Brokerage,Stamps,TransferFee    float32    //均价、成交金额、发生金额、佣金、印花税、过户费
}
//公共统计结构
type Stats struct {
	Code,Name string
	FirstDealDay,LastDealDay,MaxBalanceDay time.Time
	HoldDays,MaxBalance,TransactionCount int
	//amount的总和，为正则是利润，为负则是成本
	Amount,TransactionFreq  float32
}
//依据可能的条件，获得代码与最新名称的映射
func getNameMap(cond string) map[string]string {
	//先拿代码列表
	sql:="select distinct code from stock "+cond
	codes:=[]string{}
	rows, _ := Db.Query(sql)
	for rows.Next() {
		c:=""
		rows.Scan(&c)
		codes = append(codes, c)
	}
	rows.Close()
	//再拿最新名字列表
	sql="select name from stock where code=? order by date desc"
	names:=make(map[string]string)
	for _,c:=range codes {
		name:=""
		Db.QueryRow(sql,c).Scan(&name)
		names[c]=name
	}
	return names
}
type ByProfitReverse []Clear
func (a ByProfitReverse) Len() int { return len(a) }
func (a ByProfitReverse) Less(i, j int) bool { return a[i].Amount > a[j].Amount }
func (a ByProfitReverse) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type ByCost []Hold
func (a ByCost) Len() int { return len(a) }
func (a ByCost) Less(i, j int) bool { return a[i].Amount < a[j].Amount }
func (a ByCost) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
/*
//以更好看的方式，显示全部的百分比数值
func perDisp(f float32) (fs string) {
	fs = fmt.Sprintf("%.2f", f)
	for {
		hasDot, TrailZero := strings.Contains(fs, "."), strings.HasSuffix(fs, "0")
		if !TrailZero || !hasDot {
			break
		} else {
			fs = strings.TrimSuffix(fs, "0")
			fs = strings.TrimSuffix(fs, ".")
		}
	}
	return
}*/