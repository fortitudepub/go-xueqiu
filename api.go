package main

import (
	"fmt"
	"net/http"
	"encoding/json"
//        "bytes"
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
func (apiClient *ApiClient) FetchStockList(symbol string) *StockList {
	httpReq, err := http.NewRequest("GET",
		"https://xueqiu.com/stock/forchartk/stocklist.json?symbol="+symbol,
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

	var rsp StockList
	err = json.NewDecoder(resp.Body).Decode(&rsp)
	if err != nil  {
		fmt.Printf("decode share chg data failed, error %v", err)
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
