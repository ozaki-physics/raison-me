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
			name: "文字列を渡して生成できるか?",
			args: args{
				data: "Password123!",
			},
			want: want{
				// $2a$10$UTmmO8T1nfe0vP28Hbl0.uUM/b00yVAY9Ck9QGv3ETqp1PAOtjhPO など
				first: "",
			},
			err: err{
				hasErr: false,
				msg:    "",
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
			name: "72バイト以上でエラーになるか?",
			args: args{
				data: "あいうえおあいうえおあいうえおあいうえおあいうえ",
			},
			want: want{
				first: "",
			},
			err: err{
				hasErr: true,
				msg:    "パスワードを短くしてください",
			},
		},
		// {
		// 	// bcrypt.DefaultCost を 31以上にしたら発生したりする
		// 	name: "ハッシュ化できないエラーになるか?",
		// 	args: args{
		// 		data: "あいうえおあいうえおあいうえおあいうえおあいう",
		// 	},
		// 	want: want{
		// 		first: "",
		// 	},
		// 	err: err{
		// 		hasErr: true,
		// 		msg:    "",
		// 	},
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := domain.NewPassword(tt.args.data)

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

			// if got.ToHash() != tt.want.first {}
			// t.Errorf(got.ToHash())
			if got.ToHash() == "" {
				t.Errorf("実際の値は %v, 想定した値は %v", got.ToHash(), tt.want)
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
			name: "ハッシュ化された文字列を渡して生成できるか?",
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

			if got.ToHash() != tt.want.first {
				t.Errorf("実際の値は %v, 想定した値は %v", got.ToHash(), tt.want)
			}
		})
	}
}
