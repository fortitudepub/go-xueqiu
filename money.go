package main

import (
	"flag"
	"time"
	"fmt"
	"strconv"
)

func parseYMD(ymd string, loc *time.Location) time.Time {
	y, _ := strconv.ParseInt(ymd[0:4], 10, 32)
	m, _ := strconv.ParseInt(ymd[4:6], 10, 32)
	d, _ := strconv.ParseInt(ymd[6:8], 10, 32)
	return time.Date(int(y), time.Month(m), int(d), 0, 0, 0, 0, loc)
}

func main() {
	stockPtr := flag.String("stock", "SZ000001", "a stock number in shex/szex")
	flag.Parse()

	stockNo := *stockPtr

	secondsEastOfUTC := int((8 * time.Hour).Seconds())
	beijing := time.FixedZone("Beijing Time", secondsEastOfUTC)

	apiClient := NewApiClient()
	apiClient.SetAuthToken("9fe481846a05a6291aada0c2cbdb8529f6013d4d")

	dataMap := make(map[string]*DetailPerInterval)
	dataMap["20111231"] = &DetailPerInterval{}
	dataMap["20121231"] = &DetailPerInterval{}
	dataMap["20131231"] = &DetailPerInterval{}
	dataMap["20141231"] = &DetailPerInterval{}
	dataMap["20151231"] = &DetailPerInterval{}
	dataMap["20161231"] = &DetailPerInterval{}
	dataMap["20171231"] = &DetailPerInterval{}

	// 从此数据中获取年末时的股本情况
	list := apiClient.FetchShareChgData(stockNo)
	for k, v := range dataMap {
		kTime := parseYMD(k, beijing)
		nearestTime := time.Date(1, time.Month(1), 1, 0, 0, 0, 0, beijing)
		for _, item := range list.List {
			if len(item.Publishdate) != 8 {
				// bad year
				continue
			}
			pTime := parseYMD(item.Publishdate, beijing)
			if pTime.After(kTime) {
				continue
			}

			if pTime.After(nearestTime) {
				nearestTime = pTime
				v.shareCount = float64(item.Totalshare)
			}
		}
	}

	// 从此数据中获取年末时的股价情况
	stockList := apiClient.FetchStockList(stockNo)
	for k, v := range dataMap {
		kTime := parseYMD(k, beijing)
		nearestTime := time.Date(1, 1, 1, 0, 0, 0, 0, beijing)
		for _, item := range stockList.Chartlist {
			// conver msec to sec.
			closeTime := time.Unix(int64(item.TimeStamp/1000), 0).Add(time.Hour*8) // fix utc to beijing time.
			if closeTime.After(kTime) {
				continue
			}
			if closeTime.After(nearestTime) {
				nearestTime = closeTime
				v.closePrice = float64(item.Close)
			}
		}
	}

	// 从此数据中获取年度主营业务营收入及主营业务利润
	stockFinanceInfoList := apiClient.FetchStockFinanceInfoList(stockNo)
	for k, v := range dataMap {
		kTime := parseYMD(k, beijing)
		nearestTime := time.Date(1, time.Month(1), 1, 0, 0, 0, 0, beijing)
		for _, item := range stockFinanceInfoList.List {
			if len(item.ReportDate) != 8 {
				// bad year
				continue
			}
			pTime := parseYMD(item.ReportDate, beijing)
			if pTime.After(kTime) {
				continue
			}

			if pTime.After(nearestTime) {
				nearestTime = pTime
				v.mainBusiIncome = float64(item.MainBusiIncome)
				v.mainBusiProfit = float64(item.MainBusiProfit)
			}
		}
	}

	yearList := [...]string{"20111231", "20121231", "20131231", "20141231", "20151231", "20161231", "20171231"}

	// 进行数据计算，得到期望的数据集
	var lastYearMarketCap float64
	var lastYearMainBusiIncome float64
	var lastYearMainBusiProfit float64
	for _, year := range yearList {
		v := dataMap[year]
		v.marketCap = v.closePrice * v.shareCount
		if (lastYearMarketCap != 0) {
			v.investGainRate = (v.marketCap - lastYearMarketCap) / lastYearMarketCap
		}
		lastYearMarketCap = v.marketCap

		if (lastYearMainBusiIncome != 0) {
			v.mainBusiIncomeGrowRate = (v.mainBusiIncome - lastYearMainBusiIncome) / lastYearMainBusiIncome
		}
		lastYearMainBusiIncome = v.mainBusiIncome

		if (lastYearMainBusiProfit != 0) {
			v.mainBusiProfitGrowRate = (v.mainBusiProfit - lastYearMainBusiProfit) / lastYearMainBusiProfit
		}
		lastYearMainBusiProfit = v.mainBusiProfit
	}

	// go map range是无序的
	for _, year := range yearList {
		fmt.Printf("@ %v data %+v\n", year, dataMap[year])
	}
}
