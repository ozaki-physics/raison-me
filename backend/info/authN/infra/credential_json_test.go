package infra_test

import (
	"testing"

	"github.com/ozaki-physics/raison-me/info/authN/domain"
	"github.com/ozaki-physics/raison-me/info/authN/infra"
)

func Test_credentialRepoJSON_Fetch(t *testing.T) {
	// リポジトリは生成できる前提
	repo, jsonErr := infra.NewCredentialRepoJSON()
	if jsonErr != nil {
		t.Errorf(jsonErr.FullError())
	}

	// 実際にメソッドを動かす
	credentials, err := repo.Fetch()
	if err != nil {
		t.Errorf(err.FullError())
	}
	for _, v := range credentials {
		t.Errorf("%+v", v)
	}
}

func Test_credentialRepoJSON_FindByAccountId(t *testing.T) {
	// リポジトリは生成できる前提
	repo, jsonErr := infra.NewCredentialRepoJSON()
	if jsonErr != nil {
		t.Errorf(jsonErr.FullError())
	}

	// 実際にメソッドを動かす
	search, err := domain.ReNewAccountID("a-cmb9qloqtgph1nvh5ni0")
	if err != nil {
		t.Errorf(err.FullError())
	}

	credential, err := repo.FindByAccountId(search)
	if err != nil {
		t.Errorf(err.FullError())
	}
	t.Errorf("%+v\n", credential)
}

func Test_credentialRepoJSON_Insert(t *testing.T) {
	// リポジトリは生成できる前提
	repo, jsonErr := infra.NewCredentialRepoJSON()
	if jsonErr != nil {
		t.Errorf(jsonErr.FullError())
	}

	// 実際にメソッドを動かす
	uAID, _ := domain.NewAccountID()
	uID, _ := domain.NewUserID("DNA")
	uName, _ := domain.NewUserName("デオキシリボ核酸")
	u, _ := domain.NewUser(
		uAID,
		uID,
		uName,
	)
	pPID, _ := domain.NewPassID()
	pPass, _ := domain.NewPassword("deoxyribonucleic")
	pDate, _ := domain.NewDate()
	p, _ := domain.NewPass(
		pPID,
		uAID,
		pPass,
		pDate,
	)

	c, err := domain.NewCredential(*u, *p)
	if err != nil {
		t.Errorf(err.FullError())
	}
	credential, err := repo.Insert(*c)
	if err != nil {
		t.Errorf(err.FullError())
	}
	t.Errorf("%+v\n", credential)

}
