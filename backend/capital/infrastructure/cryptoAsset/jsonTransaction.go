package cryptoasset

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	domain "github.com/ozaki-physics/raison-me/capital/domain/cryptoAsset"
	infra "github.com/ozaki-physics/raison-me/capital/infrastructure"
)

// dataTransaction JSON から struct に変換する
type dataTransaction struct {
	Coin map[string]interface{} `json:"data"`
}

// coinTransaction JSON から struct に変換する
type coinTransaction struct {
	Symbol      string `json:"symbol"`
	Transaction []struct {
		Id        int     `json:"id"`
		Side      int     `json:"side"`
		Price     float64 `json:"price"`
		Size      float64 `json:"size"`
		Fee       int     `json:"fee"`
		Timestamp string  `json:"timestamp"`
	} `json:"transaction"`
}

// transactionInfra トランザクションインフラの実体
type transactionInfra struct {
	data map[string]coinTransaction
}

// CreateTransactionRepository トランザクションインフラのコンストラクタ
// 戻り値がインタフェースだから 実装を強制できる
func CreateTransactionRepository() domain.TransactionRepository {
	// 何度も JSON を読み込まなくていいように インスタンス変数に格納しておく

	persist := infra.NewConfig()
	// 取引履歴を読み込む
	bytes, err := os.ReadFile(persist.GetTransaction())
	if err != nil {
		log.Fatalln(err)
	}
	// JSON から struct にする
	var d dataTransaction
	if err := json.Unmarshal(bytes, &d); err != nil {
		log.Fatalln(err)
	}

	var coins = make(map[string]coinTransaction)
	// JSON の key が動的だから 再度 byte 化して struct 化する
	for key, value := range d.Coin {
		byteValue, err := json.Marshal(value)
		if err != nil {
			log.Fatalln(err)
		}
		var c coinTransaction
		if err := json.Unmarshal(byteValue, &c); err != nil {
			log.Fatalln(err)
		}

		if key != c.Symbol {
			log.Printf("key(%s) または 取引履歴のシンボル(%s)が間違っている", key, c.Symbol)
			continue
		}
		coins[c.Symbol] = c
	}

	return &transactionInfra{coins}
}

// func (t *transactionInfra) Create(dt *domain.Transaction) (*domain.Transaction, error) {}

func (t *transactionInfra) FindByID(id int) (*domain.Transaction, error) {
	for _, jc := range t.data {
		for _, jt := range jc.Transaction {
			if jt.Id == id {
				t, err := domain.ReconstructTransaction(
					jc.Symbol,
					jt.Side,
					jt.Price,
					jt.Size,
					jt.Fee,
					jt.Timestamp,
				)
				if err != nil {
					return nil, err
				}
				return t, nil
			}
		}
	}
	return nil, errors.New("トランザクションIDが存在しません")
}

func (t *transactionInfra) FindBySymbol(symbol string) (*[]domain.Transaction, error) {
	// インスタンス変数に格納してある JSON 情報から取り出す
	jc := t.data[symbol]

	var dt []domain.Transaction
	for _, jt := range jc.Transaction {
		t, err := domain.ReconstructTransaction(
			jc.Symbol,
			jt.Side,
			jt.Price,
			jt.Size,
			jt.Fee,
			jt.Timestamp,
		)
		if err != nil {
			return nil, err
		}
		dt = append(dt, *t)
	}

	return &dt, nil
}

// func (t *transactionInfra) Save(dt *domain.Transaction) (*domain.Transaction, error) {}

// func (t *transactionInfra) Delete(dt *domain.Transaction) error {}
