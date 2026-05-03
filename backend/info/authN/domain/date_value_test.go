package domain_test

import (
	"errors"
	"testing"
	"time"

	"github.com/ozaki-physics/raison-me/info/authN/domain"
)

func TestReNewDate(t *testing.T) {
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
				data: "2023-05-21T21:00:00+09:00",
			},
			want: want{
				first: "2023-05-21T21:00:00+09:00",
			},
			err: err{
				hasErr: false,
				msg:    "",
			},
		},
		{
			name: "空文字列を渡して生成できないか?",
			args: args{
				data: "u-asdf",
			},
			want: want{
				first: "",
			},
			err: err{
				hasErr: true,
				msg:    "日付を入力してください",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := time.Parse(time.RFC3339, tt.args.data)
			if (err != nil) != tt.err.hasErr {
				t.Errorf("実際の error = %v, hasErr %v", err, tt.err.hasErr)
				return
			}
			got, err := domain.ReNewDate(d)

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

			if got.MyFormat() != tt.want.first {
				t.Errorf("実際の値は %v, 想定した値は %v", got.MyFormat(), tt.want)
			}
		})
	}
}

func TestNewDateFromString(t *testing.T) {
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
	tests := []struct {
		name string
		args args
		want want
		err  err
	}{
		{
			name: "文字列を渡して生成できるか?",
			args: args{
				data: "2023-05-21T21:00:00+09:00",
			},
			want: want{
				first: "2023-05-21T21:00:00+09:00",
			},
			err: err{
				hasErr: false,
				msg:    "",
			},
		},
		{
			name: "空文字列を渡して生成できないか?",
			args: args{
				data: "",
			},
			want: want{
				first: "",
			},
			err: err{
				hasErr: true,
				msg:    "日付のパースに失敗しました",
			},
		},
		{
			name: "日付のフォーマットに合わない文字列を渡して生成できないか?",
			args: args{
				data: "u-asdf",
			},
			want: want{
				first: "",
			},
			err: err{
				hasErr: true,
				msg:    "日付のパースに失敗しました",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := domain.NewDateFromString(tt.args.data)

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

			if got.MyFormat() != tt.want.first {
				t.Errorf("実際の値は %v, 想定した値は %v", got.MyFormat(), tt.want)
			}
		})
	}
}
