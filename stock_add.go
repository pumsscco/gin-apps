package main
import (
	"net/http"
	"github.com/gin-gonic/gin"
	"time"
	"regexp"
	"strings"
	"math"
	"unicode/utf8"
	"fmt"
)
type StockBase struct {
	Id      int  `json:"-"`	
	Code string  `json:"code" binding:"required"`
	Name string  `json:"name" binding:"required"`
	Operation  string `json:"operation" binding:"required"`
	Volume int  `json:"volume"`
	Balance     uint  `json:"balance"`
	Price float32  `json:"price" binding:"required"`
	Turnover float32  `json:"turnover" binding:"required"`
	Amount float32  `json:"amount" binding:"required"`
	Brokerage float32  `json:"brokerage"`
	Stamps float32  `json:"stamps"`
	TransferFee    float32  `json:"transfer_fee"`
}
type Stock struct {
	Date    time.Time  `json:"date" binding:"required"`
	StockBase
}
type StockV struct {
	Date    string  `json:"date" binding:"required"`
	StockBase
}
//新增一条记录
func create(c *gin.Context) {
	var (
		stock Stock
		stockV StockV
		err error
		errors string
		match bool
	)
	if err = c.ShouldBindJSON(&stockV); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
	}
	nk:=fmt.Sprintf("stock:%s:%s:%s", stockV.Date,stockV.Code,stockV.Operation)
	val,err:=client.Get(nk).Result()
    if err==nil && val=="true"{
        c.JSON(http.StatusBadRequest, gin.H{"error": "依据缓存，该交易记录已录入！"})
        return
    }
	//数据校验，第一步先拿原始串转换成相应类型
	date,err := time.Parse("2006-01-02",stockV.Date)
	if err!=nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
	}
	code := stockV.Code
	name := stockV.Name
	operation := stockV.Operation
	volume := stockV.Volume
	balance := stockV.Balance
	price  := stockV.Price
	turnover:= stockV.Turnover
	amount:= stockV.Amount
	brokerage:= stockV.Brokerage
	stamps:= stockV.Stamps
	transferFee := stockV.TransferFee
	//操作类型优先检查，只有6种
	validOp:=[]string{
		"证券买入",
		"证券卖出",
		"申购中签",
		"红股入账",
		"股息入账",
		"股息红利税补",
		"配股缴款",
	}
	for _,vo:=range validOp {
		if operation==vo {
			match=true
		}
	}
	if !match {
		errors+="不合法的操作类型\n"
	}
	//从日期开始校验全部字段，以及业务逻辑的一致性，确保即便前端被绕过，也不会出问题
	first:=time.Date(2007,10,24,0,0,0,0,time.Local)
	if date.Before(first) || date.After(time.Now()) {
		errors+="交易日期必须在2007年10月24日与今天之间\n"
	}
	if match, _ = regexp.MatchString(`(6|0|3|1)\d{5}`, code); !match {
		errors+="股票代码必须为6位数字，且以6/0/3/1之一开头\n"
	}
	if len(name)<7 && (len(name)-utf8.RuneCountInString(name)>=4) {
		errors+="股票名称的长度不少于7，且中文不得少于2个（中文为3个字符）\n"
	}
	//成交数量通常在正负10000股之内，并且与操作严格匹配
	switch operation {
	case "证券买入":
		if !(volume>0 && volume<10000 && volume % 100 ==0) {
			errors+="买入时数量要为正，且必须整手（百股）的买\n"
		}
	case "证券卖出":
		if !(volume<0 && volume>-10000) {
			errors+="卖出时数量必须为负\n"
		}
	case "申购中签","配股缴款":
		if strings.HasPrefix(code,"6") {
			if !(volume>0 && volume<10000 && volume % 1000 ==0) {
				errors+="沪市申购1000股一个签位\n"
			}
		} else if strings.HasPrefix(code,"0") || strings.HasPrefix(code,"3") {
			if !(volume>0 && volume<10000 && volume % 500 ==0) {
				errors+="深市申购500股一个签位\n"
			}
		} else if strings.HasPrefix(code,"1") {
			if !(volume>0 && volume<10000 && volume % 10 ==0) {
				errors+="可转债申购10股一个签位\n"
			}
		}
	case "红股入账":
		if !(volume>0 && volume<10000) {
			errors+="送股必须是正数\n"
		}
	case "股息入账","股息红利税补":
		if volume!=0 {
			errors+="股息没有成交\n"
		}
	}
	//成交金额部分：如果为买卖，则成交金额为成交均价乘以成交数量的积的绝对值
	switch operation {
	case "证券买入","证券卖出","申购中签","配股缴款":
		if turnover!=float32(math.Abs(float64(price)*float64(volume))) {
			errors+="成交金额计算错误，请自查\n"
		}
	case "股息入账","股息红利税补":
		if turnover<=0 {
			errors+="成交金额必须大于0\n"
		}
	case "红股入账":
		if turnover!=0 {
			errors+="必须没有成交金额\n"
		}
	}
	//佣金计算
	match, _ = regexp.MatchString(`(6|0|3)\d{5}`, code)
	if match && strings.HasPrefix(operation,"证券") {
		if brokerage!=float32(math.Max(math.Round(float64(turnover)*0.00025*100)/100,5)) {
			errors+="佣金计算错误，请自查\n"
		}
	} else if strings.HasPrefix(code,"1") && operation=="证券卖出" {
		if brokerage==0 {
			errors+="可转债卖出佣金无法依规律计算出准确值，但必须不为0\n"
		}
	}
	//印花税计算
	if match && operation=="证券卖出" {
		if stamps!=float32(math.Round(float64(turnover)/1000*100)/100) {
			errors+="印花税计算错误，请自查\n"
		}
	} else {
		if stamps!=0 {
			errors+="非卖出股票，以及全部可转债交易，不收印花税\n"
		}
	}
	//发生金额计算
	switch operation {
	case "证券买入":
		if strings.HasPrefix(code,"6") {
			if amount!=-turnover-brokerage-transferFee {
				errors+="沪市买入时，发生金额为成交+佣金+过户费的结果的负值\n"
			}
		} else if strings.HasPrefix(code,"0") || strings.HasPrefix(code,"3") {
			if amount!=-turnover-brokerage {
				errors+="深市买入时，发生金额为成交+佣金的结果的负值（过户费虽计算，但并不扣减）\n"
			}
		}
	case "证券卖出":
		if strings.HasPrefix(code,"6") {
			if amount!=-turnover-brokerage-stamps-transferFee {
				errors+="沪市卖出时，发生金额为成交-佣金-印花税-过户费\n"
			}
		} else if strings.HasPrefix(code,"0") || strings.HasPrefix(code,"3") {
			if amount!=-turnover-brokerage-stamps {
				errors+="深市卖出时，发生金额为成交-佣金-印花税（过户费虽计算，但并不扣减）\n"
			}
		}
	case "申购中签","股息红利税补":
		if amount!=-turnover {
			errors+="申购及税补时发生金额为成交金额取负\n"
		}
	case "红股入账","股息入账":
		if amount!=turnover {
			errors+="送红股及发股息时，发生金额与成交金额必定相同\n"
		}
	}
	//过户费，如果录入不了，可能是券商计算有小的误差，此时只得手工录入了
	if match && strings.HasPrefix(operation,"证券") {
		//由于过户费的计算中，厘进位到分存在取舍不一的情况，因此只要差距在正负1分钱，均算正常
		if transferFee!=float32(math.Round(float64(turnover)*0.00002*100)/100) {
			errors+="过户费可能计算错误，请自查，如确为券商计算问题，请直接将数据插入数据库\n"
		}
	} else {
		if transferFee!=0 {
			errors+="过户费必须为0\n"
		}
	}
	//增加对余额，也即持仓量的校验，本次成交量与上回余额的和，应与本回余额相同
	chkSql:=`select balance from stock where code=? order by date desc limit 1`
	var lastBal int
	Db.QueryRow(chkSql,code).Scan(&lastBal)
	if int(balance)!=int(lastBal)+volume {
		errors+="最近持仓量与本次成交量的和，必须与本次交易后持仓量相同！\n"
	}
	if errors!="" {
		logger.Println("errors: ",errors)
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": errors})
        return
	}
	stock.Amount=amount
	stock.Code =  code
	stock.Name = name
	stock.Operation = operation
	stock.Volume = volume
	stock.Balance =  balance
	stock.Price  = price
	stock.Turnover=  turnover
	stock.Amount= amount
	stock.Brokerage= brokerage 
	stock.Stamps= stamps
	stock.TransferFee =  transferFee
	stock.Date=date
	statement := `insert into stock(date,code,name,operation,volume,balance,
		price,turnover,amount,brokerage,stamps,transfer_fee)
		value (?,?,?,?,?,?,?,?,?,?,?,?)
	`
	_, err = Db.Exec(statement,stock.Date,stock.Code,stock.Name,stock.Operation,stock.Volume,stock.Balance,
		stock.Price,stock.Turnover,stock.Amount,stock.Brokerage,stock.Stamps,stock.TransferFee)
    if err != nil {
        errors=fmt.Sprintf("无法创建新记录，错误：%v",err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": errors})
        return
    } else {
        client.Set(nk,"true",2*time.Hour)
        c.JSON(http.StatusOK, gin.H{"result": "恭喜！成功增加新交易记录"})
    }
}