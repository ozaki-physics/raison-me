package infra

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ozaki-physics/raison-me/info/authN/domain"
)

type passRepoSQL struct {
	pool *pgxpool.Pool
}

// 戻り値 が インタフェース だから 実装を強制できる
func NewPassRepoSQL(pool *pgxpool.Pool) (domain.PassRepo, error) {
	return &passRepoSQL{pool}, nil
}

func (prs *passRepoSQL) Insert(ctx context.Context, pass domain.Pass) (*domain.Pass, error) {
	sql_statement := `
		INSERT INTO app.passwords (
			password_id
			, account_id
			, password_hash
			, iat
		)
			VALUES
		($1, $2, $3, $4)
		;
	`

	pID := pass.ID()
	aID := pass.AccountID()
	password := pass.Password()
	iat := pass.IssuedAt()

	_, err := prs.pool.Exec(ctx, sql_statement, pID.Val(), aID.Val(), password.HashedText(), iat.MyFormat())
	if err != nil {
		log.Printf("exec: %v\n", err)
		return nil, err
	}

	return &pass, nil
}

func (prs *passRepoSQL) FindByAccountId(ctx context.Context, accountID domain.AccountID) (*domain.Pass, error) {
	sql_statement := `
		SELECT
			password_id
			, account_id
			, password_hash
			, iat
		FROM app.passwords
		WHERE account_id = $1
		;
	`

	row := prs.pool.QueryRow(ctx, sql_statement, accountID.Val())

	var passwordID, accountIDStr, passwordStr string
	var iat time.Time
	if err := row.Scan(&passwordID, &accountIDStr, &passwordStr, &iat); err != nil {
		log.Printf("scan: %v\n", err)
		return nil, err
	}

	p, err := domain.ReNewPass(passwordID, accountIDStr, passwordStr, iat)
	if err != nil {
		log.Printf("ReNewPass: %v\n", err)
		return nil, err
	}
	return p, nil
}
