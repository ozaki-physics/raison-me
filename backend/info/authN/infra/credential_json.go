package infra

import (
	"encoding/json"

	"github.com/ozaki-physics/raison-me/info/authN/domain"
	"github.com/ozaki-physics/raison-me/share"
)

type jsonUser struct {
	User []user `json:"user"`
}

// struct を切り出したのは 保存時に json 用の struct に変換する必要があるから
type user struct {
	AccountID string `json:"account_id"`
	Id        string `json:"id"`
	Name      string `json:"name"`
}

type jsonPass struct {
	Pass []pass `json:"pass"`
}

// struct を切り出したのは 保存時に json 用の struct に変換する必要があるから
type pass struct {
	PassID    string `json:"pass_id"`
	AccountID string `json:"account_id"`
	Password  string `json:"password"`
	IssuedAt  string `json:"issued_at"`
}

type credentialRepoJSON struct {
	// 何度も JSON を読み込まなくていいように インスタンス変数に格納しておく
	// map にした理由は特に無い
	dataUser map[domain.AccountID]domain.User
	dataPass map[domain.AccountID]domain.Pass
}

// 戻り値 が インタフェース だから 実装を強制できる
func NewCredentialRepoJSON() (domain.CredentialRepo, InfraError) {
	s := NewStoragePath()

	var ju jsonUser
	userJSONErr := share.JsonToStruct(&ju, s.userJSON)
	if userJSONErr != nil {
		return nil, WrapInfraError("JSONの取得に失敗しました", userJSONErr)
	}
	// JSON を インスタンス に保持させるため
	var dul = make(map[domain.AccountID]domain.User)
	for _, user := range ju.User {
		du, err := domain.ReNewUser(user.AccountID, user.Id, user.Name)
		if err != nil {
			return nil, WrapInfraError("JSONからUserオブジェクトの再生成に失敗しました", err)
		}
		dul[du.AccountID()] = *du
	}

	var jp jsonPass
	passJSONErr := share.JsonToStruct(&jp, s.passJSON)
	if passJSONErr != nil {
		return nil, WrapInfraError("JSONの取得に失敗しました", passJSONErr)
	}
	// JSON を インスタンス に保持させるため
	var dpl = make(map[domain.AccountID]domain.Pass)
	for _, pass := range jp.Pass {
		dp, err := domain.ReNewPass(pass.PassID, pass.AccountID, pass.Password, pass.IssuedAt)
		if err != nil {
			return nil, WrapInfraError("JSONからPassオブジェクトの再生成に失敗しました", err)
		}
		dpl[dp.AccountID()] = *dp
	}

	return &credentialRepoJSON{dul, dpl}, nil
}

func (crj *credentialRepoJSON) Insert(c domain.Credential) (*domain.Credential, domain.DomainError) {
	// なぜか1個ずつやらないと構文エラーらしいから
	u := c.User()
	aID := u.AccountID()

	crj.dataUser[aID] = c.User()
	crj.dataPass[aID] = c.Pass()

	// 保存する
	err := crj.save()
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (crj *credentialRepoJSON) Fetch() ([]domain.Credential, domain.DomainError) {

	var cl []domain.Credential
	for _, user := range crj.dataUser {
		c, err := domain.NewCredential(user, crj.dataPass[user.AccountID()])
		if err != nil {
			return nil, err
		}
		cl = append(cl, *c)
	}
	return cl, nil
}

func (crj *credentialRepoJSON) FindByAccountId(accountID domain.AccountID) (*domain.Credential, domain.DomainError) {
	for _, user := range crj.dataUser {
		if user.AccountID() == accountID {
			c, err := domain.NewCredential(user, crj.dataPass[user.AccountID()])
			if err != nil {
				return nil, err
			}
			return c, nil
		}
	}
	// DomainError interface が戻り値だけど 実体は InfraError 型で値を返却できる
	// TODO: ただ見つからなかっただけで エラーじゃない? &domain.Credential{} にすべき?
	return nil, NewInfraError("該当のユーザーが見つかりません")
}

func (crj *credentialRepoJSON) save() InfraError {
	var jul jsonUser
	for _, u := range crj.dataUser {
		// json で保存するために strct に プリミティブ型が要求されるため プリミティブ型にしないといけない
		// なぜか1個ずつやらないと構文エラーらしいから
		aID := u.AccountID()
		uID := u.ID()
		uName := u.Name()
		ju := user{
			aID.Val(),
			uID.Val(),
			uName.Val(),
		}
		jul.User = append(jul.User, ju)
	}

	// strct を ファイル保存用 の bytes に変換
	jub, err := json.Marshal(jul)
	if err != nil {
		return WrapInfraError("構造体をバイトに変換できませんでした", err)
	}

	var jpl jsonPass
	for _, p := range crj.dataPass {
		// json で保存するために strct に プリミティブ型が要求されるため プリミティブ型にしないといけない
		// なぜか1個ずつやらないと構文エラーらしいから
		pID := p.PassID()
		aID := p.AccountID()
		password := p.Password()
		d := p.IssuedAt()

		jp := pass{
			pID.Val(),
			aID.Val(),
			password.ToHash(),
			d.MyFormat(),
		}
		jpl.Pass = append(jpl.Pass, jp)
	}

	// strct を ファイル保存用 の bytes に変換
	jpb, err := json.Marshal(jpl)
	if err != nil {
		return WrapInfraError("構造体をバイトに変換できませんでした", err)
	}

	// 同期して保存したいファイル群を用意する
	s := NewStoragePath()
	saveJSON, err := share.NewSyncJSONs(jub, s.userJSON)
	if err != nil {
		return WrapInfraError("保存でエラーになりました", err)
	}
	err02 := saveJSON.Add(jpb, s.passJSON)
	if err02 != nil {
		return WrapInfraError("保存でエラーになりました", err02)
	}
	err03 := saveJSON.Save()
	if err03 != nil {
		return WrapInfraError("保存でエラーになりました", err03)
	}

	return nil
}
