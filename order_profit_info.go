/***************************************************
 ** @Desc : This file for ...
 ** @Time : 2019/10/30 11:44
 ** @Author : yuebin
 ** @File : order_profit_info
 ** @Last Modified by : yuebin
 ** @Last Modified time: 2019/10/30 11:44
 ** @Software: GoLand
****************************************************/
package models

import (
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"strconv"
	"strings"
)

type OrderProfitInfo struct {
	Id              int
	MerchantName    string
	MerchantUid     string
	AgentName       string
	AgentUid        string
	PayProductCode  string
	PayProductName  string
	PayTypeCode     string
	PayTypeName     string
	Status          string
	MerchantOrderId string
	BankOrderId     string
	BankTransId     string
	OrderAmount     float64
	ShowAmount      float64
	FactAmount      float64
	UserInAmount    float64
	SupplierRate    float64
	PlatformRate    float64
	AgentRate       float64
	AllProfit       float64
	SupplierProfit  float64
	PlatformProfit  float64
	AgentProfit     float64
	UpdateTime      string
	CreateTime      string
}

const ORDER_PROFIT_INFO = "order_profit_info"

func GetOrderProfitByBankOrderId(bankOrderId string) OrderProfitInfo {
	o := orm.NewOrm()
	var orderProfit OrderProfitInfo
	_, err := o.QueryTable(ORDER_PROFIT_INFO).Filter("bank_order_id", bankOrderId).Limit(1).All(&orderProfit)
	if err != nil {
		logs.Error("GetOrderProfitByBankOrderId fail：", err)
	}
	return orderProfit
}

func GetOrderProfitLenByMap(params map[string]string) int {
	o := orm.NewOrm()
	qs := o.QueryTable(ORDER_PROFIT_INFO)
	for k, v := range params {
		if len(v) > 0 {
			qs = qs.Filter(k, v)
		}
	}
	cnt, _ := qs.Limit(-1).Count()
	return int(cnt)
}

func GetOrderProfitByMap(params map[string]string, display, offset int) []OrderProfitInfo {
	o := orm.NewOrm()
	var orderProfitInfoList []OrderProfitInfo
	qs := o.QueryTable(ORDER_PROFIT_INFO)
	for k, v := range params {
		if len(v) > 0 {
			qs = qs.Filter(k, v)
		}
	}
	_, err := qs.Limit(display, offset).OrderBy("-update_time").All(&orderProfitInfoList)
	if err != nil {
		logs.Error("get order by map fail: ", err)
	}
	return orderProfitInfoList
}

func GetPlatformProfitByMap(params map[string]string) []PlatformProfit {

	o := orm.NewOrm()

	cond := "select merchant_name, agent_name, pay_product_name as supplier_name, pay_type_name, sum(fact_amount) as order_amount, count(1) as order_count, " +
		"sum(platform_profit) as platform_profit, sum(agent_profit) as agent_profit from " + ORDER_PROFIT_INFO + " where status='success' "
	flag := false
	for k, v := range params {
		if len(v) > 0 {
			if flag {
				cond += " and"
			}
			if strings.Contains(k, "create_time__gte") {
				cond = cond + " create_time>='" + v + "'"
			} else if strings.Contains(k, "create_time__lte") {
				cond = cond + " create_time<='" + v + "'"
			} else {
				cond = cond + " " + k + "='" + v + "'"
			}
			flag = true
		}
	}

	cond += " group by merchant_uid, agent_uid, pay_product_code, pay_type_code"

	var platformProfitList []PlatformProfit
	_, err := o.Raw(cond).QueryRows(&platformProfitList)
	if err != nil {
		logs.Error("get platform profit by map fail:", err)
	}

	return platformProfitList
}
func Todayallmonertlist(user string,params map[string]string)(todayall float64 ,todaypress int ,yesdatay float64,yesdataypress int,userinamount float64,today_prtessv float64) {



	//yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	//yesterdayStart := yesterday + " 00:00:00"
	//yesterdayEnd := yesterday + " 23:59:59"

	//ystart, _ := time.Parse("2006-01-02 15:04:05", yesterdayStart)
	//yend, _ := time.Parse("2006-01-02 15:04:05", yesterdayEnd)


	o := orm.NewOrm()

	cond := "select  sum(fact_amount) as order_amount,sum(all_profit) as today_prtess, count(1) as order_count, " +
		"sum(platform_profit) as platform_profit, sum(agent_profit) as agent_profit, sum(user_in_amount) as userinamountv from " + ORDER_PROFIT_INFO + " where status='success'"
	//ycond := "select  sum(fact_amount) as yorder_amount,sum(all_profit) as ytoday_prtess, count(1) as order_countv, " +
	//	"sum(platform_profit) as platform_profit, sum(agent_profit) as agent_profit, sum(user_in_amount) as userinamountv from " + ORDER_PROFIT_INFO + " where status='success' and merchant_uid =? and create_time BETWEEN ? AND ?"
	var maps []orm.Params
	//var ymaps []orm.Params
	for _, v := range params {
		if len(v) > 0 {
			cond = cond + "and "
			break
		}
	}
	flag := false
	for k, v := range params {
		if len(v) > 0 {

			if flag {
				cond += " and"
			}
			if strings.Contains(k, "create_time__gte") {
				cond = cond + " create_time>='" + v + "'"
			} else if strings.Contains(k, "create_time__lte") {
				cond = cond + " create_time<='" + v + "'"
			} else {
				cond = cond + " " + k + "='" + v + "'"
			}
			flag = true
		}
	}
	if params["pay_type_code"] != "" {
		if flag {
			cond = cond + " and "
		}
		cond = cond + " pay_type_code = '" + params["pay_type_code"] + "'"
	}

	cond += " group by merchant_uid, agent_uid, pay_product_code, pay_type_code"




	o.Raw(cond).Values(&maps)
	//o.Raw(ycond,user, ystart, yend).Values(&ymaps)
	//logs.Info("sql语句是=="+cond)
	//_, _ = o.Raw(cond).Values(&maps)
	allAmount := 0.00
	today_prtess :=0

	yesAmount := 0.00
	yes_prtess :=0
	userinamountv :=0.00

	if maps[0]["order_amount"] == nil{
		allAmount = 0.00
	}else {
		allAmount, _ = strconv.ParseFloat(maps[0]["order_amount"].(string), 64)
	}
	if maps[0]["userinamountv"] == nil{
		userinamountv = 0.00
	}else {
		userinamountv, _ = strconv.ParseFloat(maps[0]["userinamountv"].(string), 64)
	}

	if maps[0]["today_prtess"] == nil{
		today_prtessv = 0.00
	}else {
		today_prtessv, _ = strconv.ParseFloat(maps[0]["today_prtess"].(string), 64)
	}

	if maps[0]["order_count"] == nil{
		today_prtess = 0
	}else {
		today_prtess,_=strconv.Atoi(maps[0]["order_count"].(string))
	}


	//if ymaps[0]["yorder_amount"] == nil{
	//	yesAmount = 0.00
	//}else {
	//	yesAmount, _ = strconv.ParseFloat(ymaps[0]["yorder_amount"].(string), 64)
	//}
	//if ymaps[0]["order_countv"] == nil{
	//	yes_prtess = 0
	//}else {
	//	yes_prtess, _ = strconv.Atoi(ymaps[0]["order_countv"].(string))
	//}
	//
	//if ymaps[0]["userinamountv"] == nil{
	//	userinamountv = 0.00
	//}else {
	//	userinamountv, _ = strconv.ParseFloat(ymaps[0]["userinamountv"].(string), 64)
	//}

   fmt.Print(maps[0])
	//allAmount, _ = strconv.ParseFloat(maps[0]["order_amount"].(string), 64)
	//today_prtess, _ = strconv.ParseFloat(maps[0]["today_prtess"].(string), 64)
	//

	return allAmount,today_prtess,yesAmount,yes_prtess,userinamountv,today_prtessv
}