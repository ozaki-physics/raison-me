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

func Main() {
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
	if err := printResults(os.Stdout, rows); err != nil {
		log.Fatal(err)
	}
}

// 取得 SQL
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
func printResults(w io.Writer, iter *bigquery.RowIterator) error {
	for {
		var row BabyNamesRow
		err := iter.Next(&row)
		if err == iterator.Done {
			return nil
		}
		if err != nil {
			return fmt.Errorf("error iterating through results: %v", err)
		}

		fmt.Fprintf(w, "name: %s, count: %d\n", row.Name, row.Count)
	}
}
