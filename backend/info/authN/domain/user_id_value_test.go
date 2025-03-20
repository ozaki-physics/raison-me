package domain_test

import (
	"errors"
	"testing"

	"github.com/ozaki-physics/raison-me/info/authN/domain"
)

func TestNewUserID(t *testing.T) {
	// テスト対象に渡す必要がある引数
	type args struct {
		id string
	}
	// テスト用の値たち
	tests := []struct {
		name   string
		args   args
		want   string
		hasErr bool
	}{
		// テストケース
		{
			name: "文字列を渡して生成できるか?",
			args: args{
				id: "sena",
			},
			want:   "sena",
			hasErr: false,
		},
		{
			name: "ブランクでドメインエラーになるか?",
			args: args{
				id: "",
			},
			want:   "ユーザーIDを入力してください",
			hasErr: true,
		},
	}
	// テスト関数の実行前
	// 処理
	defer func() {
		// テスト関数の実行後
		// 処理
	}()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// テストケースの実行前
			// 処理
			defer func() {
				// テストケースの実行後
				// 処理
			}()

			got, err := domain.NewUserID(tt.args.id)

			if (err != nil) != tt.hasErr {
				t.Errorf("実際の error = %v, hasErr %v", err, tt.hasErr)
				return
			}

			if (err != nil) && tt.hasErr {
				var de = domain.NewDomainError("")
				if errors.As(err, &de) {
					if err.Error() != tt.want {
						t.Errorf("実際のエラーは %v, 想定されるエラーは %v", err.Error(), tt.want)
					}
				} else {
					t.Errorf("実際の error = %v, hasErr %v", err, tt.hasErr)
				}
				return
			}

			if got.Val() != tt.want {
				t.Errorf("実際の値は %v, 想定した値は %v", got.Val(), tt.want)
			}
		})
	}
}
