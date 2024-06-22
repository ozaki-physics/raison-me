package domain_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/ozaki-physics/raison-me/info/authN/domain"
)

func TestNewAccountID(t *testing.T) {
	// テスト対象に渡す必要がある引数
	type args struct {
	}
	type want struct {
		first  string
		second string
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
			name: "AccountIDが生成できるか?",
			args: args{},
			want: want{
				first:  "a",
				second: "",
			},
			err: err{
				hasErr: false,
				msg:    "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := domain.NewAccountID()

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

			if strings.Split(got.Val(), "-")[0] != tt.want.first {
				t.Errorf("実際の値は %v, 想定した値は %v", got.Val(), tt.want)
			}
		})
	}
}

func TestReNewAccountID(t *testing.T) {
	// テスト対象に渡す必要がある引数
	type args struct {
		data string
	}
	type want struct {
		first  string
		second string
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
			name: "文字列を渡してAccountIDが生成できるか?",
			args: args{
				data: "a-asdf",
			},
			want: want{
				first:  "a",
				second: "asdf",
			},
			err: err{
				hasErr: false,
				msg:    "",
			},
		},
		{
			name: "AccountIDに適さないプレフィックスでエラーになるか?",
			args: args{
				data: "u-asdf",
			},
			want: want{
				first:  "u",
				second: "asdf",
			},
			err: err{
				hasErr: true,
				msg:    "AccountIDに設定できないプレフィックスです",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := domain.ReNewAccountID(tt.args.data)

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

			if strings.Split(got.Val(), "-")[0] != tt.want.first {
				t.Errorf("実際の値は %v, 想定した値は %v", got.Val(), tt.want)
			}
		})
	}
}
