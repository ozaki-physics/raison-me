package infrastructure

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type quote struct {
	Currency struct {
		Price                    float64 `json:"price"`
		Volume_24h               float64 `json:"volume_24h"`
		Volume_change_24h        float64 `json:"volume_change_24h"`
		Percent_change_1h        float64 `json:"percent_change_1h"`
		Percent_change_24h       float64 `json:"percent_change_24h"`
		Percent_change_7d        float64 `json:"percent_change_7d"`
		Percent_change_30d       float64 `json:"percent_change_30d"`
		Percent_change_60d       float64 `json:"percent_change_60d"`
		Percent_change_90d       float64 `json:"percent_change_90d"`
		Market_cap               float64 `json:"market_cap"`
		Market_cap_dominance     float64 `json:"market_cap_dominance"`
		Fully_diluted_market_cap float64 `json:"fully_diluted_market_cap"`
		Last_updated             string  `json:"last_updated"`
	} `json:"JPY"`
}

type tags struct {
	Slug     string `json:"slug"`
	Name     string `json:"name"`
	Category string `json:"category"`
}

type quotesLatest struct {
	ID                            int     `json:"id"`
	Name                          string  `json:"name"`
	Symbol                        string  `json:"symbol"`
	Slug                          string  `json:"slug"`
	NumMarketPairs                int     `json:"num_market_pairs"`
	DateAdded                     string  `json:"date_added"`
	Tags                          []tags  `json:"tags"`
	MaxSupply                     int     `json:"max_supply"`
	CirculatingSupply             float64 `json:"circulating_supply"`
	TotalSupply                   float64 `json:"total_supply"`
	IsActive                      int     `json:"is_active"`
	Platform                      int     `json:"platform"`
	CmcRank                       int     `json:"cmc_rank"`
	IsFiat                        int     `json:"is_fiat"`
	SelfReportedCirculatingSupply int     `json:"self_reported_circulating_supply"`
	SelfReportedMarketCap         int     `json:"self_reported_market_cap"`
	LastUpdated                   string  `json:"last_updated"`
	Quote                         quote   `json:"quote"`
}

type responseQuotesLatest struct {
	Status responseStatus `json:"status"`
	// key が動的だから 単純に struct に変換できない map[string]interface{} を使う
	Data map[string]interface{} `json:"data"`
}

type usePrice struct {
	Symbol string
	CMCID  int
	Name   string
	Price  float64
}

func ExampleGetQuotesLatest(c *coinMarketCap) {
	var CMCIDs []int
	CMCIDs = append(CMCIDs, 1)
	CMCIDs = append(CMCIDs, 1027)
	tmp := c.getQuotesLatest(CMCIDs)
	fmt.Printf("%+v\n", tmp)
}

func (c *coinMarketCap) getQuotesLatest(CMCIDs []int) []usePrice {
	const entryURL = "/v2/cryptocurrency/quotes/latest"

	// クライアントの作成
	client := &http.Client{}
	req, err := http.NewRequest("GET", c.baseURL+entryURL, nil)
	if err != nil {
		log.Println(err)
	}

	// リクエストヘッダーの作成
	req.Header.Set("Accepts", "application/json")
	// req.Header.Set("Accept-Encoding", "deflate, gzip")
	req.Header.Set("X-CMC_PRO_API_KEY", c.Key)

	// クエリ文字列パラメータの作成
	q := url.Values{}
	id := makeQueryParmCMCID(&CMCIDs)
	if id == "" {
		log.Fatalln("必須パラメータが無い")
	}
	q.Add("id", id)
	q.Add("convert", "JPY")
	// q.Add("aux", "")
	req.URL.RawQuery = q.Encode()

	// リクエストする
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}

	// ステータスコードを確認
	log.Println(resp.Status)

	// ResponseBody を取り出す
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	// JSON にして出力してみる
	// 動的な key と status を削ったレスポンス
	var qls []quotesLatest

	var responseJSON responseQuotesLatest
	if err := json.Unmarshal(respBody, &responseJSON); err != nil {
		log.Println(err)
	}
	// 動的な部分を扱う 第1戻り値が 動的な key
	for _, value := range responseJSON.Data {
		// 動的な key の中身 struct にするために 一度 []byte にし直す
		// 通常だと value は map[string]interface{} 型になってしまう
		byteValue, err := json.Marshal(value)
		if err != nil {
			log.Fatalln(err)
		}
		// 一度 []byte にした部分を struct にする
		var ql quotesLatest
		if err := json.Unmarshal(byteValue, &ql); err != nil {
			log.Println(err)
		}
		qls = append(qls, ql)
	}

	return generateSymbolAndPrice(&qls)
}

func generateSymbolAndPrice(qls *[]quotesLatest) []usePrice {
	var priceMap []usePrice
	for _, ql := range *qls {
		tmp := usePrice{
			ql.Symbol,
			ql.ID,
			ql.Name,
			ql.Quote.Currency.Price,
		}
		priceMap = append(priceMap, tmp)
	}
	return priceMap
}

// makeQueryParmCMCID CMCID のスライスをカンマ区切りの1個の string にする
func makeQueryParmCMCID(CMCIDs *[]int) string {
	var queryParmID string
	if len(*CMCIDs) == 0 {
		return queryParmID
	}

	for i, id := range *CMCIDs {
		if i != 0 {
			queryParmID += ","
		}
		queryParmID += strconv.Itoa(id)
	}
	return queryParmID
}
