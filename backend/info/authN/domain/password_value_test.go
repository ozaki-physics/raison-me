package domain_test

import (
	"errors"
	"testing"

	"github.com/ozaki-physics/raison-me/info/authN/domain"
)

func TestNewPassword(t *testing.T) {
	// テスト対象に渡す必要がある引数
	type args struct {
		data string
	}
	type want struct {
		first string
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
		// テストケース
		{
			name: "ハッシュ文字列を渡して生成できるか?",
			args: args{
				data: "$2a$10$UTmmO8T1nfe0vP28Hbl0.uUM/b00yVAY9Ck9QGv3ETqp1PAOtjhPO",
			},
			want: want{
				first: "$2a$10$UTmmO8T1nfe0vP28Hbl0.uUM/b00yVAY9Ck9QGv3ETqp1PAOtjhPO",
			},
			err: err{
				hasErr: false,
			},
		},
		{
			name: "ブランクでエラーになるか?",
			args: args{
				data: "",
			},
			want: want{
				first: "",
			},
			err: err{
				hasErr: true,
				msg:    "パスワードを入力してください",
			},
		},
		{
			name: "bcrypt 形式ではない文字列でエラーになるか?",
			args: args{
				data: "NotHashText",
			},
			want: want{
				first: "",
			},
			err: err{
				hasErr: true,
				msg:    "パスワードハッシュが不正です",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := domain.NewPassword(tt.args.data)

			if (err != nil) != tt.err.hasErr {
				t.Errorf("error の有無が想定と異なります, 実際のエラー: %v, 想定されるエラーの有無: %v", err, tt.err.hasErr)
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

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got.HashedText() != tt.want.first {
				t.Errorf("実際の値は %v, 想定した値は %v", got.HashedText(), tt.want.first)
			}
		})
	}
}

func TestReNewPassword(t *testing.T) {
	// テスト対象に渡す必要がある引数
	type args struct {
		data string
	}
	type want struct {
		first string
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
			name: "ハッシュ文字列を渡して生成できるか?",
			args: args{
				data: "$2a$10$UTmmO8T1nfe0vP28Hbl0.uUM/b00yVAY9Ck9QGv3ETqp1PAOtjhPO",
			},
			want: want{
				first: "$2a$10$UTmmO8T1nfe0vP28Hbl0.uUM/b00yVAY9Ck9QGv3ETqp1PAOtjhPO",
			},
			err: err{
				hasErr: false,
				msg:    "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := domain.ReNewPassword(tt.args.data)

			if (err != nil) != tt.err.hasErr {
				t.Errorf("error の有無が想定と異なります, 実際のエラー: %v, 想定されるエラーの有無: %v", err, tt.err.hasErr)
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got.HashedText() != tt.want.first {
				t.Errorf("実際の値は %v, 想定した値は %v", got.HashedText(), tt.want.first)
			}
		})
	}
}
