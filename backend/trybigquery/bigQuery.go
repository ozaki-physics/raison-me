package trybigquery

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"cloud.google.com/go/bigquery"
	globalConfig "github.com/ozaki-physics/raison-me/share/config"
	"google.golang.org/api/iterator"
)

// BigQuery へ 接続 して クエリを 実行 する
// 環境変数 GOOGLE_APPLICATION_CREDENTIALS で JSON キー ファイル を 指定 して使う
// 自作 メソッド を 呼び出す 側
func getRecord() string {
	globalConfig := globalConfig.NewConfig()

	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, globalConfig.GetGCPProjectID())
	if err != nil {
		log.Fatalf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	rows, err := query(ctx, client)
	if err != nil {
		log.Fatal(err)
	}
	results, err := printResults(os.Stdout, rows)
	if err != nil {
		log.Fatal(err)
	}

	return results
}

// 取得 SQL 文
func query(ctx context.Context, client *bigquery.Client) (*bigquery.RowIterator, error) {

	query := client.Query(
		`
		SELECT
			name,
			count
		FROM
			` + "`babynames.names_2014`" + `
		WHERE
			gender = 'M'
		ORDER BY
			count DESC
		LIMIT
			5
		`,
	)
	return query.Read(ctx)
}

// 取得 した レコード の構造体
type BabyNamesRow struct {
	Name  string `bigquery:"name"`
	Count int    `bigquery:"count"`
}

// 取得 した レコード の出力
func printResults(w io.Writer, iter *bigquery.RowIterator) (string, error) {
	results := []BabyNamesRow{}
	for {
		var row BabyNamesRow
		err := iter.Next(&row)
		if err == iterator.Done {
			return fmt.Sprintf("%v", results), nil
		}
		if err != nil {
			return "", fmt.Errorf("error iterating through results: %v", err)
		}

		fmt.Fprintf(w, "name: %s, count: %d\n", row.Name, row.Count)
		results = append(results, row)
	}
}
