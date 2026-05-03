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
		pepper  string
		plain   string
		wantErr string
	}{
		// テストケース
		{
			name:   "平文パスワードをハッシュ化できるか?",
			pepper: "server-side-pepper",
			plain:  "Password123!",
		},
		{
			name:    "ブランクでエラーになるか?",
			pepper:  "server-side-pepper",
			plain:   "",
			wantErr: "パスワードを入力してください",
		},
		{
			name:    "pepper を加えた結果が 72 バイト以上でエラーになるか?",
			pepper:  "0123456789",
			plain:   "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijkl",
			wantErr: "パスワードを短くしてください",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasher, err := infra.NewPasswordBcrypt(tt.pepper)
			if err != nil {
				t.Fatalf("unexpected constructor error: %v", err)
			}
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
		name        string
		pepper      string
		wrongPepper string
		plain       string
		wrong       string // 一致しないパスワード
		wantErr     string
	}{
		// テストケース
		{
			name:   "ハッシュ化したパスワードと平文が一致するか?",
			pepper: "server-side-pepper",
			plain:  "Password123!",
		},
		{
			name:        "異なる pepper でハッシュ化したパスワードと平文が一致しないか?",
			pepper:      "server-side-pepper",
			wrongPepper: "other-pepper",
			plain:       "Password123!",
			wantErr:     "パスワードが一致しません",
		},
		{
			name:        "パスワード不一致のときエラーになるか?",
			pepper:      "server-side-pepper",
			wrongPepper: "other-pepper",
			plain:       "Password234!",
			wrong:       "WrongPassword",
			wantErr:     "パスワードが一致しません",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasher, err := infra.NewPasswordBcrypt(tt.pepper)
			if err != nil {
				t.Fatalf("unexpected constructor error: %v", err)
			}
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

			if tt.wrongPepper != "" {
				wrongHasher, _ := infra.NewPasswordBcrypt(tt.wrongPepper)
				err = wrongHasher.Verify(password, tt.plain)
				if err == nil {
					t.Fatal("pepper が違うのに失敗しませんでした")
				}
				if err.Error() != tt.wantErr {
					t.Fatalf("実際のエラーは %v, 想定されるエラーは %v", err.Error(), tt.wantErr)
				}
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

func TestNewPasswordBcrypt(t *testing.T) {
	// テスト用の値たち
	tests := []struct {
		name    string
		pepper  string
		wantErr string
	}{
		// テストケース
		{
			name:   "pepper があれば hasher を生成できるか?",
			pepper: "server-side-pepper",
		},
		{
			name:    "空文字の pepper でエラーになるか?",
			pepper:  "",
			wantErr: "AUTHN_PASSWORD_PEPPER が設定されていません",
		},
		{
			name:    "空白だけの pepper でエラーになるか?",
			pepper:  "   ",
			wantErr: "AUTHN_PASSWORD_PEPPER が設定されていません",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasher, err := infra.NewPasswordBcrypt(tt.pepper)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatal("error が必要ですが nil でした")
				}
				if err.Error() != tt.wantErr {
					t.Fatalf("実際のエラーは %v, 想定されるエラーは %v", err.Error(), tt.wantErr)
				}
				if hasher != nil {
					t.Fatal("error 時は nil を返す想定です")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if hasher == nil {
				t.Fatal("hasher が nil です")
			}
		})
	}
}
