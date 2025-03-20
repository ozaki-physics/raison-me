package infra

import (
	"encoding/json"

	"github.com/ozaki-physics/raison-me/info/authN/domain"
	"github.com/ozaki-physics/raison-me/share"
)

const userPath = "./info/authN/infra/json/user.json"

type userRepoJSON struct {
	// 何度も JSON を読み込まなくていいように インスタンス変数に格納しておく
	// map にした理由は特に無い
	data map[domain.AccountID]domain.User
}

// 戻り値 が インタフェース だから 実装を強制できる
func NewUserRepoJSON() (domain.UserRepo, error) {
	var ju jsonUser
	// main.go からの path
	share.JsonToStruct(&ju, userPath)

	// JSON を インスタンス に保持させるため
	var d = make(map[domain.AccountID]domain.User)
	for _, user := range ju.User {
		ud, err := domain.ReNewUser(user.AccountID, user.Id, user.Name)
		if err != nil {
			return nil, err
		}
		d[ud.AccountID()] = *ud
	}

	return &userRepoJSON{d}, nil
}

func (urj *userRepoJSON) Insert(id domain.UserID, name domain.UserName) (*domain.User, error) {
	// TODO: 雑に作っている
	accountID, err := domain.ReNewAccountID("temp")
	if err != nil {
		return nil, err
	}

	u, err := domain.NewUser(accountID, id, name)
	if err != nil {
		return nil, err
	}

	urj.data[accountID] = *u
	if err := urj.save(); err != nil {
		return nil, err
	}

	return u, nil
}

func (urj *userRepoJSON) Fetch() ([]domain.User, error) {
	var users []domain.User
	for _, v := range urj.data {
		users = append(users, v)
	}
	return users, nil
}

func (urj *userRepoJSON) FindByAccountId(accountID domain.AccountID) (*domain.User, error) {
	ud := urj.data[accountID]
	// TODO: 存在しなかったときの扱いをどうしよう
	return &ud, nil
}

func (urj *userRepoJSON) FindById(id domain.UserID) (*domain.User, error) {
	for _, v := range urj.data {
		if v.ID() == id {
			return &v, nil
		}
	}
	// TODO: 存在しなかったときの扱いをどうしよう
	return nil, nil
}

func (urj *userRepoJSON) FindByName(name domain.UserName) (*domain.User, error) {
	for _, v := range urj.data {
		if v.Name() == name {
			return &v, nil
		}
	}
	// TODO: 存在しなかったときの扱いをどうしよう
	return nil, nil
}

func (urj *userRepoJSON) save() error {
	var j jsonUser
	for _, ud := range urj.data {
		// json で保存するために strct に プリミティブ型が要求されるため プリミティブ型にしないといけない
		// なぜか1個ずつやらないと構文エラーらしいから
		aID := ud.AccountID()
		uID := ud.ID()
		uName := ud.Name()
		u := user{
			AccountID: aID.Val(),
			Id:        uID.Val(),
			Name:      uName.Val(),
		}
		j.User = append(j.User, u)
	}

	// strct を ファイル保存用 の bytes に変換
	bytes, err := json.Marshal(j)
	if err != nil {
		return err
	}
	share.SaveJSON(bytes, userPath)
	return nil
}
