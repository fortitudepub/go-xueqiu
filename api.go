package main

import (
	"fmt"
	"net/http"
	"encoding/json"
//        "bytes"
)

const (
        STOCK_LIST_ITEM_DATA_TIMESTAMP_IDX = 0
        STOCK_LIST_ITEM_DATA_CLOSE_IDX = 5
)

type ApiClient struct {
	authToken string
	httpClient *http.Client
}

func NewApiClient() *ApiClient {
	httpClient := &http.Client{}
	apiClient := &ApiClient{
		httpClient: httpClient,
	}
	return apiClient
}

func (apiClient *ApiClient) SetAuthToken(token string) {
	apiClient.authToken = token
}

func (apiClient *ApiClient) FetchShareChgData(symbol string) *ShareChgList {
	// add a size in order to get 2011 data because default size may not feed us that much.
	httpReq, err := http.NewRequest("GET",
		"https://xueqiu.com/stock/f10/shareschg.json?symbol="+symbol+"&size=100",
		nil)
	httpReq.Header.Add("Content-type", "Application/json")
	authCookie := http.Cookie{
		Name:"xq_a_token",
		Value:apiClient.authToken,
	}
	httpReq.AddCookie(&authCookie)
	resp, err := apiClient.httpClient.Do(httpReq)
	if err != nil {
		fmt.Printf("fetch share chg data failed, error %v", err)
		return nil
	}
	defer resp.Body.Close()

	var rsp ShareChgList
	err = json.NewDecoder(resp.Body).Decode(&rsp)
	if err != nil  {
		fmt.Printf("decode share chg data failed, error %v", err)
		return nil
	}

	return &rsp
}

// 这里是不复权价格，我们用股份数与当时股价计算，要用到不复权价格
// 20200205: handle xueqiu api update.
func (apiClient *ApiClient) FetchStockList(symbol string, now int64) *StockList {
        // 10 * 365 means 10 years is enough for our analyse.
        url := fmt.Sprintf("https://stock.xueqiu.com/v5/stock/chart/kline.json?symbol=%s&begin=%v&period=day&type=normal&count=-3650&indicator=kline", symbol, now)
        fmt.Printf("stocklist url is:"+url+"\n")
	httpReq, err := http.NewRequest("GET", url, nil)
	httpReq.Header.Add("Content-type", "Application/json")
	authCookie := http.Cookie{
		Name:"xq_a_token",
		Value:apiClient.authToken,
	}
	httpReq.AddCookie(&authCookie)
	resp, err := apiClient.httpClient.Do(httpReq)
	if err != nil {
		fmt.Printf("fetch stock list data failed, error %v", err)
		return nil
	}
	defer resp.Body.Close()

	var rsp StockList
	err = json.NewDecoder(resp.Body).Decode(&rsp)
	if err != nil  {
		fmt.Printf("decode stock list data failed, error %v", err)
		return nil
	}

	return &rsp
}

func (apiClient *ApiClient) FetchStockFinanceInfoList(symbol string) *StockFinanceInfoList {
	httpReq, err := http.NewRequest("GET",
		"https://xueqiu.com/stock/f10/finmainindex.json?symbol="+symbol,
		nil)
	httpReq.Header.Add("Content-type", "Application/json")
	authCookie := http.Cookie{
		Name:"xq_a_token",
		Value:apiClient.authToken,
	}
	httpReq.AddCookie(&authCookie)
	resp, err := apiClient.httpClient.Do(httpReq)
	if err != nil {
		fmt.Printf("fetch stock finance info data failed, error %v", err)
		return nil
	}
	defer resp.Body.Close()

	var rsp StockFinanceInfoList
	err = json.NewDecoder(resp.Body).Decode(&rsp)
	if err != nil  {
		fmt.Printf("decode stock finance info data failed, error %v", err)
		return nil
	}

	return &rsp
}

func (apiClient *ApiClient) FetchStockBalanceSheetList(symbol string) *StockBalanceSheetList {
	httpReq, err := http.NewRequest("GET",
		"https://xueqiu.com/stock/f10/balsheet.json?symbol="+symbol,
		nil)
	httpReq.Header.Add("Content-type", "Application/json")
	authCookie := http.Cookie{
		Name:"xq_a_token",
		Value:apiClient.authToken,
	}
	httpReq.AddCookie(&authCookie)
	resp, err := apiClient.httpClient.Do(httpReq)
	if err != nil {
		fmt.Printf("fetch stock balance sheet data failed, error %v", err)
		return nil
	}
	defer resp.Body.Close()

	var rsp StockBalanceSheetList
	err = json.NewDecoder(resp.Body).Decode(&rsp)
	if err != nil  {
		fmt.Printf("decode stock balance sheet data failed, error %v", err)
		return nil
	}

	return &rsp
}