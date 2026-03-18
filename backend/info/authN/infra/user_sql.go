package infra

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ozaki-physics/raison-me/info/authN/domain"
)

type userRepoSQL struct {
	pool *pgxpool.Pool
}

// 戻り値 が インタフェース だから 実装を強制できる
func NewUserRepoSQL(pool *pgxpool.Pool) (domain.UserRepo, error) {
	return &userRepoSQL{pool}, nil
}

func (urs *userRepoSQL) Insert(ctx context.Context, user *domain.User) (*domain.User, error) {
	sql_statement := `
		INSERT INTO app.accounts (
			account_id
			, user_id
			, user_name
		)
			VALUES
		($1, $2, $3)
		;
	`

	aID := user.AccountID()
	uID := user.ID()
	uName := user.Name()

	_, err := urs.pool.Exec(ctx, sql_statement, aID.Val(), uID.Val(), uName.Val())
	if err != nil {
		log.Printf("exec: %v\n", err)
		return nil, err
	}

	return user, nil
}

func (urs *userRepoSQL) Fetch(ctx context.Context) ([]domain.User, error) {
	sql_statement := `
		SELECT
			account_id
			, user_id
			, user_name
		FROM app.accounts
		;
	`

	rows, err := urs.pool.Query(ctx, sql_statement)
	if err != nil {
		log.Printf("query: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var aID, uID, uName string
		if err := rows.Scan(&aID, &uID, &uName); err != nil {
			log.Printf("scan: %v\n", err)
			return nil, err
		}
		user, err := domain.ReNewUser(aID, uID, uName)
		if err != nil {
			log.Printf("renew user: %v\n", err)
			return nil, err
		}
		users = append(users, *user)
	}

	// 途中で ストリーム が エラー になることがある
	// rows.Next() が false を返すだけで エラーかどうかは分からない
	// rows.Err() で クエリの実行中に発生したエラーを検出する
	if err := rows.Err(); err != nil {
		log.Printf("rows: %v\n", err)
		return nil, err
	}

	return users, nil
}

func (urs *userRepoSQL) FindByAccountId(ctx context.Context, accountID domain.AccountID) (*domain.User, error) {
	sql_statement := `
		SELECT
			account_id
			, user_id
			, user_name
		FROM app.accounts
		WHERE account_id = $1
		;
	`

	row := urs.pool.QueryRow(ctx, sql_statement, accountID.Val())

	var aID, uID, uName string
	if err := row.Scan(&aID, &uID, &uName); err != nil {
		log.Printf("scan: %v\n", err)
		return nil, err
	}

	user, err := domain.ReNewUser(aID, uID, uName)
	if err != nil {
		log.Printf("renew user: %v\n", err)
		return nil, err
	}

	return user, nil
}

func (urs *userRepoSQL) FindById(ctx context.Context, id domain.UserID) (*domain.User, error) {
	sql_statement := `
		SELECT
			account_id
			, user_id
			, user_name
		FROM app.accounts
		WHERE user_id = $1
		;
	`

	row := urs.pool.QueryRow(ctx, sql_statement, id.Val())

	var aID, uID, uName string
	if err := row.Scan(&aID, &uID, &uName); err != nil {
		log.Printf("scan: %v\n", err)
		return nil, err
	}

	user, err := domain.ReNewUser(aID, uID, uName)
	if err != nil {
		log.Printf("ReNewUser: %v\n", err)
		return nil, err
	}

	return user, nil
}
