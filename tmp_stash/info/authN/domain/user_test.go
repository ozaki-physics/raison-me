package domain_test

import (
	"errors"
	"testing"

	"github.com/ozaki-physics/raison-me/info/authN/domain"
)

func TestReNewUser(t *testing.T) {
	// テスト対象に渡す必要がある引数
	type args struct {
		accountID string
		userID    string
		name      string
	}
	type want struct {
		accountID string
		userID    string
		name      string
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
			name: "Userオブジェクトをプリミティブ型から生成できるか?",
			args: args{
				accountID: "a-123",
				userID:    "sena_us",
				name:      "sena",
			},
			want: want{
				accountID: "a-123",
				userID:    "sena_us",
				name:      "sena",
			},
			err: err{
				hasErr: false,
				msg:    "",
			},
		},
		{
			name: "AccountIDのエラーを戻せるか?",
			args: args{
				accountID: "p-123",
				userID:    "sena_us",
				name:      "sena",
			},
			want: want{
				accountID: "a-123",
				userID:    "sena_us",
				name:      "sena",
			},
			err: err{
				hasErr: true,
				msg:    "AccountIDに設定できないプレフィックスです",
			},
		},
		{
			name: "UserIDのエラーを戻せるか?",
			args: args{
				accountID: "a-123",
				userID:    "",
				name:      "sena",
			},
			want: want{
				accountID: "a-123",
				userID:    "sena_us",
				name:      "sena",
			},
			err: err{
				hasErr: true,
				msg:    "ユーザーIDを入力してください",
			},
		},
		{
			name: "UserNameのエラーを戻せるか?",
			args: args{
				accountID: "a-123",
				userID:    "sena_us",
				name:      "",
			},
			want: want{
				accountID: "a-123",
				userID:    "sena_us",
				name:      "sena",
			},
			err: err{
				hasErr: true,
				msg:    "ユーザー名を入力してください",
			},
		},
		{
			name: "エラーが複数あっても最初のエラーを戻せるか?",
			args: args{
				accountID: "a-123",
				userID:    "",
				name:      "",
			},
			want: want{
				accountID: "a-123",
				userID:    "sena_us",
				name:      "sena",
			},
			err: err{
				hasErr: true,
				msg:    "ユーザーIDを入力してください",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := domain.ReNewUser(tt.args.accountID, tt.args.userID, tt.args.name)

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
			a, _ := domain.ReNewAccountID(tt.want.accountID)
			b, _ := domain.NewUserID(tt.want.userID)
			c, _ := domain.NewUserName(tt.want.name)

			if got.AccountID() != a {
				t.Errorf("実際の値は %v, 想定した値は %v", got.AccountID(), a)
			}
			if got.ID() != b {
				t.Errorf("実際の値は %v, 想定した値は %v", got.ID(), b)
			}
			if got.Name() != c {
				t.Errorf("実際の値は %v, 想定した値は %v", got.Name(), c)
			}
		})
	}
}
