package main

//列出指定代码的股票的全部交易记录，按时间逆序
func getDealList(code string) (stocks []Stock) {
	sql := `
		select date,code,name,operation,volume,balance,
		price,turnover,amount,brokerage,stamps,transfer_fee from stock where code=? order by date desc
    `
	rows, _ := Db.Query(sql,code)
	for rows.Next() {
		s := Stock{}
		rows.Scan(
			&s.Date,&s.Code,&s.Name,&s.Operation,&s.Volume,&s.Balance,
			&s.Price,&s.Turnover,&s.Amount,&s.Brokerage,&s.Stamps,&s.TransferFee,
		)
		stocks = append(stocks, s)
	}
	rows.Close()
	return
}
func getHoldLastDeal() (stocks []Stock) {
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
		Db.QueryRow(outSql,c).Scan(
			&s.Date,&s.Code,&s.Name,&s.Operation,&s.Volume,&s.Balance,
			&s.Price,&s.Turnover,&s.Amount,&s.Brokerage,&s.Stamps,&s.TransferFee,
		)
		stocks=append(stocks,s)
	}
	return
}