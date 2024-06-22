package domain_test

import (
	"errors"
	"testing"

	"github.com/ozaki-physics/raison-me/info/authN/domain"
)

func TestReNewPass(t *testing.T) {
	// テスト対象に渡す必要がある引数
	type args struct {
		passID    string
		accountID string
		password  string
		iat       string
	}
	type want struct {
		passID    string
		accountID string
		password  string
		iat       string
	}
	type err struct {
		hasErr bool
		msg    string
	}
	// テスト用の値たち
	tests := []struct {
		name string
		args args
		want want
		err  err
	}{
		// テストケースたち
		{
			name: "Passオブジェクトをプリミティブ型から生成できるか?",
			args: args{
				passID:    "p-123",
				accountID: "a-456",
				password:  "$2a$10$UTmmO8T1nfe0vP28Hbl0.uUM/b00yVAY9Ck9QGv3ETqp1PAOtjhPO",
				iat:       "2023-05-27T22:21:00+09:00",
			},
			want: want{
				passID:    "p-123",
				accountID: "a-456",
				password:  "$2a$10$UTmmO8T1nfe0vP28Hbl0.uUM/b00yVAY9Ck9QGv3ETqp1PAOtjhPO",
				iat:       "2023-05-27T22:21:00+09:00",
			},
			err: err{
				hasErr: false,
				msg:    "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := domain.ReNewPass(
				tt.args.passID,
				tt.args.accountID,
				tt.args.password,
				tt.args.iat,
			)

			if (err != nil) != tt.err.hasErr {
				t.Errorf("実際の error = %v, hasErr %v", err, tt.err.hasErr)
				return
			}

			if (err != nil) && tt.err.hasErr {
				var de = domain.NewDomainError("")
				if errors.As(err, &de) {
					if err.Error() != tt.err.msg {
						t.Errorf("実際のエラーは %v, 想定されるエラーは %v", err.Error(), tt.err.msg)
					}
				} else {
					t.Errorf("実際の error = %v, hasErr %v", err, tt.err.hasErr)
				}
				return
			}

			// それぞれの値オブジェクトの生成は正常と仮定する
			a, _ := domain.ReNewPassID(tt.want.passID)
			b, _ := domain.ReNewAccountID(tt.want.accountID)
			c, _ := domain.ReNewPassword(tt.want.password)
			d, _ := domain.ReNewDate(tt.want.iat)

			if got.PassID() != a {
				t.Errorf("実際の値は %v, 想定した値は %v", got.PassID(), a)
			}
			if got.AccountID() != b {
				t.Errorf("実際の値は %v, 想定した値は %v", got.AccountID(), b)
			}
			if got.Password() != c {
				t.Errorf("実際の値は %v, 想定した値は %v", got.Password(), c)
			}
			dd := got.IssuedAt()
			if dd.MyFormat() != d.MyFormat() {
				t.Errorf("実際の値は %v, 想定した値は %v", got.IssuedAt(), d)
			}
		})
	}
}
