package cryptoasset

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"

	domain "github.com/ozaki-physics/raison-me/capital/domain/cryptoAsset"
)

// responseStatus レスポンスの共通部分
type responseStatus struct {
	Timestamp    string `json:"timestamp"`
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	Elapsed      int    `json:"elapsed"`
	CreditCount  int    `json:"credit_count"`
	Notice       string `json:"notice"`
}

// platform CMCIDMap と Metadata の共通部分
type platform struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Symbol       string `json:"symbol"`
	Slug         string `json:"slug"`
	TokenAddress string `json:"token_address"`
}

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
	ID                            int      `json:"id"`
	Name                          string   `json:"name"`
	Symbol                        string   `json:"symbol"`
	Slug                          string   `json:"slug"`
	NumMarketPairs                int      `json:"num_market_pairs"`
	DateAdded                     string   `json:"date_added"`
	Tags                          []tags   `json:"tags"`
	MaxSupply                     int      `json:"max_supply"`
	CirculatingSupply             float64  `json:"circulating_supply"`
	TotalSupply                   float64  `json:"total_supply"`
	IsActive                      int      `json:"is_active"`
	Platform                      platform `json:"platform"`
	CmcRank                       int      `json:"cmc_rank"`
	IsFiat                        int      `json:"is_fiat"`
	SelfReportedCirculatingSupply int      `json:"self_reported_circulating_supply"`
	SelfReportedMarketCap         int      `json:"self_reported_market_cap"`
	LastUpdated                   string   `json:"last_updated"`
	Quote                         quote    `json:"quote"`
}

type responseQuotesLatest struct {
	Status responseStatus `json:"status"`
	// key が動的だから 単純に struct に変換できない map[string]interface{} を使う
	Data map[string]interface{} `json:"data"`
}

// coinInfra コインインフラの実体
// 本当は 何も保持しなくてもいいが 必要なものを先に準備するニュアンスで用意する
type coinInfra struct {
	credential CredentialCmc
	cmcIds     CMCIds
}

// cmcPriceDto レスポンスを infra 層で使いやすい形に詰め替える
type cmcPriceDto struct {
	symbol string
	cmcId  int
	name   string
	price  float64
}

// CreateCoinRepository コインインフラのコンストラクタ
// 戻り値がインタフェースだから 実装を強制できる
// infra 層にも関わらず DI する
func CreateCoinRepository(cr CredentialCmc, cmcids CMCIds) domain.CoinRepository {
	// 何度も JSON を読み込まなくていいように インスタンス変数に格納しておく
	return &coinInfra{cr, cmcids}
}

func (c *coinInfra) FindBySymbol(symbol string) (*domain.Coin, error) {
	// ターゲットとなるコインが 手元で記録している cmcId 一覧に存在するか探す
	for s, cmcId := range c.cmcIds.SymbolAndId() {
		if s != symbol {
			continue
		}

		// CoinMarketCap にリクエストして 現在価格を取得する
		var cmcIds []int
		cmcIds = append(cmcIds, cmcId)
		// 現時点では if 文で symbol 1個 つまり 1個しか引数に渡してないから 1個しか返却されてないと仮定する
		price := c.getQuotesLatest(&cmcIds)[0].price

		dc, err := domain.ReconstructCoin(symbol, price)
		if err != nil {
			return nil, err
		}
		return dc, nil
	}
	return nil, errors.New("CoinMarketCap ID が存在しません")
}

// getQuotesLatest CoinMarketCap にリクエストして現在価格を取得する
// API 仕様として 複数のコインの情報を1リクエストで取得できるから 引数がスライスになっている
func (c *coinInfra) getQuotesLatest(cmcIdsParam *[]int) []cmcPriceDto {
	const entryURL = "/v2/cryptocurrency/quotes/latest"

	// クライアントの作成
	client := &http.Client{}
	req, err := http.NewRequest("GET", c.credential.BaseUrl()+entryURL, nil)
	if err != nil {
		log.Println(err)
	}
	// リクエストヘッダーの作成
	req.Header.Set("Accepts", "application/json")
	// req.Header.Set("Accept-Encoding", "deflate, gzip")
	req.Header.Set("X-CMC_PRO_API_KEY", c.credential.ApiKey())
	// クエリ文字列パラメータの作成
	q := url.Values{}
	id := makeQueryParmCMCID(cmcIdsParam)
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
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

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

// generateSymbolAndPrice レスポンスから infra 層で使いやすい形に詰め替える
func generateSymbolAndPrice(qls *[]quotesLatest) []cmcPriceDto {
	var priceMap []cmcPriceDto
	for _, ql := range *qls {
		tmp := cmcPriceDto{
			ql.Symbol,
			ql.ID,
			ql.Name,
			ql.Quote.Currency.Price,
		}
		priceMap = append(priceMap, tmp)
	}
	return priceMap
}
