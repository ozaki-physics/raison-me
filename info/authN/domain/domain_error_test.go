package domain_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ozaki-physics/raison-me/info/authN/domain"
)

func TestNewDomainError(t *testing.T) {
	// テスト対象に渡す必要がある引数
	type args struct {
		msg string
	}
	type want struct {
		sameType bool
		msg      string
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
			name: "ドメイン層の独自エラーを生成できるか?",
			args: args{
				msg: "hello",
			},
			want: want{
				sameType: true,
				msg:      "hello",
			},
			err: err{
				hasErr: false,
				msg:    "",
			},
		},
		{
			name: "生成時にメッセージを格納しなかったとき エラー構造体にエラーメッセージが格納されるか?",
			args: args{
				msg: "",
			},
			want: want{
				sameType: true,
				msg:      "ドメイン層のエラーを生成しましたが エラーメッセージが格納されていません",
			},
			err: err{
				hasErr: false,
				msg:    "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := domain.NewDomainError(tt.args.msg)
			var de domain.DomainError

			if errors.As(got, &de) != tt.want.sameType {
				t.Errorf("DomainError じゃない型のエラーになっています")
			}

			if got.Error() != tt.want.msg {
				t.Error("エラーメッセージの内容が異なって生成されています")
			}
		})
	}
}

func TestWrapDomainError(t *testing.T) {
	// テスト対象に渡す必要がある引数
	type args struct {
		msg01 string
		msg02 string
	}
	type want struct {
		sameType bool
		msg01    string
		msg02    string
		fullMsg  string
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
			name: "ドメイン層の独自エラーでラップできるか?",
			args: args{
				msg01: "hello",
				msg02: "world",
			},
			want: want{
				sameType: true,
				msg01:    "hello",
				msg02:    "world",
				fullMsg:  "world: hello",
			},
			err: err{
				hasErr: false,
				msg:    "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := fmt.Errorf(tt.args.msg01)
			got := domain.WrapDomainError(tt.args.msg02, e)
			got2 := fmt.Errorf("ラップ: %w", got)
			got3 := fmt.Errorf("ラップ: %w", got2)
			var de domain.DomainError

			if errors.As(got, &de) != tt.want.sameType {
				t.Errorf("DomainError じゃない型のエラーになっています")
			}

			if got.Error() != tt.want.msg02 {
				t.Error("エラーメッセージの内容が異なって生成されています")
			}

			if got2.Error() != fmt.Sprintf("ラップ: %v", tt.want.msg02) {
				t.Error("ドメイン層の独自エラーをラップしたのに エラーメッセージの内容が異なって生成されています")
			}

			if uwe := errors.Unwrap(got); uwe != nil {
				if errors.Is(uwe, e) != true {
					t.Error("ドメイン層の独自エラーの Unwrap() が上手に動いてないかも")
				}
			}

			// 任意のエラー型が存在することを確認したのち 型アサーションして 該当エラー型のメソッドを使う
			if errors.As(got3, &de) {
				var wrap error
				wrap = got3
				for {
					// 無限に wrap 変数 に再代入して 詳細情報を取り出し続ける
					if isE, ok := wrap.(domain.DomainError); ok {
						if isE.FullError() != tt.want.fullMsg {
							t.Error("ラップしたメッセージが異なります")
						}
						break
					}

					if unwrappedErr := errors.Unwrap(wrap); unwrappedErr != nil {
						// 再度 Unwrap するため
						wrap = unwrappedErr
					} else {
						break
					}
				}
			}
		})
	}
}
