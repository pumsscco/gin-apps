package main
import (
	"net/http"
	"time"
	"regexp"
	"strings"
	"math"
	"strconv"
	"unicode/utf8"
)
//新增一条记录
func createDeal(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    if err != nil {
        logger.Println(err, "Cannot parse form")
	}
	errors:=""
	//数据校验，第一步先拿原始串转换成相应类型
	date,_ := time.Parse("2006-01-02",r.PostFormValue("date"))
	code := r.PostFormValue("code")
	name := r.PostFormValue("name")
	operation := r.PostFormValue("operation")
	volume,_ := strconv.Atoi(r.PostFormValue("volume"))
	//balance,_ := strconv.Atoi(r.PostFormValue("balance"))
	price,_ := strconv.ParseFloat(r.PostFormValue("price"),32)
	turnover,_ := strconv.ParseFloat(r.PostFormValue("turnover"),32)
	amount,_ := strconv.ParseFloat(r.PostFormValue("amount"),32)
	brokerage,_ := strconv.ParseFloat(r.PostFormValue("brokerage"),32)
	stamps,_ := strconv.ParseFloat(r.PostFormValue("stamps"),32)
	transferFee,_ := strconv.ParseFloat(r.PostFormValue("transfer-fee"),32)
	//从日期开始校验全部字段，以及业务逻辑的一致性，确保即便前端被绕过，也不会出问题
	first:=time.Date(2007,10,24,0,0,0,0,time.Local)
	if date.Before(first) || date.After(time.Now()) {
		errors+="交易日期必须在2007年10月24日与今天之间\n"
	}
	if match, _ := regexp.MatchString(`(6|0|3)\d{5}`, code); !match {
		errors+="股票代码必须为6位数字，且以6/0/3之一开头\n"
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
	case "申购中签":
		if strings.HasPrefix(code,"6") {
			if !(volume>0 && volume<10000 && volume % 1000 ==0) {
				errors+="沪市申购1000股一个签位\n"
			}
		} else if strings.HasPrefix(code,"0") || strings.HasPrefix(code,"3") {
			if !(volume>0 && volume<10000 && volume % 500 ==0) {
				errors+="深市申购500股一个签位\n"
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
	case "证券买入","证券卖出","申购中签":
		if turnover!=math.Abs(float64(price)*float64(volume)) {
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
	if operation=="证券买入" || operation=="证券卖出" {
		if brokerage!=math.Max(math.Round(turnover*0.00025*100)/100,5) {
			errors+="佣金计算错误，请自查\n"
		}
	} else {
		if brokerage!=0 {
			errors+="佣金必须为0\n"
		}
	}
	//印花税计算
	if operation=="证券卖出" {
		if stamps!=math.Round(turnover/1000*100)/100 {
			errors+="印花税计算错误，请自查\n"
		}
	} else {
		if stamps!=0 {
			errors+="非卖出股票不收印花税\n"
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
	//default:
		if amount!=turnover {
			errors+="送红股及发股息时，发生金额与成交金额必定相同\n"
		}
	default:
		errors+="不合法的操作类型\n"
	}
	//过户费，如果录入不了，可能是券商计算有小的误差，此时只得手工录入了
	if operation=="证券买入" || operation=="证券卖出" {
		//由于过户费的计算中，厘进位到分存在取舍不一的情况，因此只要差距在正负1分钱，均算正常
		if brokerage!=math.Round(turnover*0.00002*100)/100 {
			errors+="过户费可能计算错误，请自查，如确为券商计算问题，请直接将数据插入数据库\n"
		}
	} else {
		if brokerage!=0 {
			errors+="过户费必须为0\n"
		}
	}
	//增加对余额，也即持仓量的校验，本次成交量与上回余额的和，应与本回余额相同
	if errors!="" {
		logger.Println("errors: ",errors)
	} else {
		logger.Println("表单合法")
	}
	//logger.Printf("date type:  %T ;   code type  :   %T   ;   operation type  : %T  ;    volume type:  %T   ;" ,date,code,operation,volume)
	/*
    //增加数据有效性校验，按说应该在前端做，但自己不写js，就改后端实现了
    //上班时间为9：00～17：30，如果不对，跳回重填
    if checkin.YearDay()!=checkout.YearDay() || checkin.Hour()>9 || checkout.Hour()<17 || (checkout.Hour()==17 && checkout.Minute()<30) {
        logger.Println("上班时间为9：00～17：30，禁止跨日，请自行检查")
        http.Redirect(w, r, "/new-attn", 302)    
	}
	statement := "insert into attendance(checkin, checkout, comments) value (?,?,?)"
    _, err = Db.Exec(statement,checkin,checkout,comments)
    if err != nil {
        logger.Println(err, "Cannot create deal record")
    }
    http.Redirect(writer, request, "/latest", 302)*/
}