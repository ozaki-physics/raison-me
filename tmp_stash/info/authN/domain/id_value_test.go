package domain_test

import (
	"errors"
	"testing"

	"github.com/ozaki-physics/raison-me/info/authN/domain"
)

func TestNewID(t *testing.T) {
	// テスト対象に渡す必要がある引数
	type args struct {
		prefixID string
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
			name: "正しいプレフィックスで生成できるか?",
			args: args{
				prefixID: "u",
			},
			want:   "u",
			hasErr: false,
		},
		{
			name: "異なるプレフィックスで生成できるか?",
			args: args{
				prefixID: "a",
			},
			want:   "a",
			hasErr: false,
		},
		{
			name: "プレフィックスがブランクでドメインエラーになるか?",
			args: args{
				prefixID: "",
			},
			want:   "ID生成時のプレフィックスが存在しません",
			hasErr: true,
		},
		{
			name: "プレフィックスの最大値で生成できるか?",
			args: args{
				prefixID: "123456789",
			},
			want:   "123456789",
			hasErr: false,
		},
		{
			name: "プレフィックスが長すぎでドメインエラーになるか?",
			args: args{
				prefixID: "1234567890",
			},
			want:   "ID生成時のプレフィックスが長すぎます",
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

			got, err := domain.NewID(tt.args.prefixID)

			if (err != nil) != tt.hasErr {
				t.Errorf("NewID() error = %v, hasErr %v", err, tt.hasErr)
				return
			}

			if (err != nil) && tt.hasErr {
				var de = domain.NewDomainError("")
				if errors.As(err, &de) {
					if err.Error() != tt.want {
						t.Errorf("実際のエラーは %v, 想定されるエラーは %v", err.Error(), tt.want)
					}
				} else {
					t.Errorf("NewID() error = %v, hasErr %v", err, tt.hasErr)
				}
				return
			}

			if got.PrefixID() != tt.want {
				t.Errorf("実際のプレフィックスは %v, 想定したプレフィックスは %v", got.PrefixID(), tt.want)
			}
		})
	}
}

func TestNilID(t *testing.T) {
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
			name: "nilなIDが生成できるか?",
			args: args{},
			want: want{
				first:  "nil",
				second: "nil",
			},
			err: err{
				hasErr: false,
				msg:    "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := domain.NilID()

			if got.Val() != tt.want.first+"-"+tt.want.second {
				t.Errorf("実際の値は %v, 想定した値は %v", got.Val(), tt.want)
			}
		})
	}
}

func TestIsNilID(t *testing.T) {
	a := domain.NilID()
	if a.IsNilID() != true {
		t.Errorf("NilIDのはずなのに 判定メソッドが NilIDじゃない となっています")
	}
	b, _ := domain.NewID("a")
	if b.IsNilID() != false {
		t.Errorf("NilIDじゃないはずなのに 判定メソッドが NilID となっています")
	}
}

func TestReNewID(t *testing.T) {
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
			name: "文字列を渡して生成できるか?",
			args: args{
				data: "u-asdf",
			},
			want: want{
				first:  "u",
				second: "asdf",
			},
			err: err{
				hasErr: false,
				msg:    "",
			},
		},
		{
			name: "文字列にハイフンがあっても生成できるか?",
			args: args{
				data: "u-asdf-ghjk",
			},
			want: want{
				first:  "u",
				second: "asdf-ghjk",
			},
			err: err{
				hasErr: false,
				msg:    "",
			},
		},
		{
			name: "文字列の後半が無いとエラーになるか?",
			args: args{
				data: "u-",
			},
			want: want{
				first:  "u",
				second: "",
			},
			err: err{
				hasErr: true,
				msg:    "ID生成時の値が存在しません",
			},
		},
		{
			name: "文字列の前半が無いとエラーになるか?",
			args: args{
				data: "-asdf",
			},
			want: want{
				first:  "",
				second: "asdf",
			},
			err: err{
				hasErr: true,
				msg:    "ID生成時のプレフィックスが存在しません",
			},
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

			got, err := domain.ReNewID(tt.args.data)

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

			if got.Val() != tt.want.first+"-"+tt.want.second {
				t.Errorf("実際の値は %v, 想定した値は %v", got.Val(), tt.want)
			}
		})
	}
}
