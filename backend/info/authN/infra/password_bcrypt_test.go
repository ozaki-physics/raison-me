package infra_test

import (
	"testing"

	"github.com/ozaki-physics/raison-me/info/authN/domain"
	"github.com/ozaki-physics/raison-me/info/authN/infra"
)

func TestPasswordBcryptHash(t *testing.T) {
	// テスト用の値たち
	tests := []struct {
		name    string
		plain   string
		wantErr string
	}{
		// テストケース
		{
			name:  "平文パスワードをハッシュ化できるか?",
			plain: "Password123!",
		},
		{
			name:    "ブランクでエラーになるか?",
			plain:   "",
			wantErr: "パスワードを入力してください",
		},
		{
			name:    "72バイト以上でエラーになるか?",
			plain:   "あいうえおあいうえおあいうえおあいうえおあいうえ",
			wantErr: "パスワードを短くしてください",
		},
	}

	hasher := infra.NewPasswordBcrypt()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := hasher.Hash(tt.plain)

			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("error が必要ですが nil でした")
				}
				if err.Error() != tt.wantErr {
					t.Fatalf("実際のエラーは %v, 想定されるエラーは %v", err.Error(), tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got == "" {
				t.Fatal("hash が空です")
			}
			if got == tt.plain {
				t.Fatal("平文と同じ文字列が返っています")
			}
		})
	}
}

func TestPasswordBcryptVerify(t *testing.T) {
	// テスト用の値たち
	tests := []struct {
		name    string
		plain   string
		wrong   string // 一致しないパスワード
		wantErr string
	}{
		// テストケース
		{
			name:  "ハッシュ化したパスワードと平文が一致するか?",
			plain: "Password123!",
		},
		{
			name:    "不一致のときエラーになるか?",
			plain:   "Password234!",
			wrong:   "WrongPassword",
			wantErr: "パスワードが一致しません",
		},
	}

	hasher := infra.NewPasswordBcrypt()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got01, err := hasher.Hash(tt.plain)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			got02, err := hasher.Hash(tt.plain)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got01 == got02 {
				t.Fatal("bcrypt の hash が同じになっています。bcrypt は同じ平文でも毎回違うハッシュを生成するはずです")
			}

			password, err := domain.NewPassword(got01)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			// ハッシュ化したパスワードと平文が一致するか
			err = hasher.Verify(password, tt.plain)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// 不一致のときエラーになるか
			if tt.wrong != "" {
				err = hasher.Verify(password, tt.wrong)
				if err == nil {
					t.Fatal("パスワード不一致なのに失敗しませんでした")
				}
				if err.Error() != tt.wantErr {
					t.Fatalf("実際のエラーは %v, 想定されるエラーは %v", err.Error(), tt.wantErr)
				}
			}
		})
	}
}
