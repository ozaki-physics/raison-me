package domain_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/ozaki-physics/raison-me/info/authN/domain"
)

func TestNewID(t *testing.T) {
	// テスト対象に渡す必要がある引数
	type args struct{}
	// テスト用の値たち
	tests := []struct {
		name   string
		args   args
		want   string
		hasErr bool
	}{
		{
			name:   "IDが生成できるか?",
			args:   args{},
			want:   "", // 生成されるIDは毎回変わるため、特定の値を期待することはできない
			hasErr: false,
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

			got, err := domain.NewID()

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

			if err != nil {
				t.Errorf("NewID() error = %v", err)
			}

			parsed, err02 := uuid.Parse(got.Val())
			if err02 != nil {
				t.Errorf("uuid.Parse() error = %v", err02)
			}

			if parsed.Version() != 7 {
				t.Errorf("実際のバージョンは %v, 想定した値は 7", parsed.Version())
			}
		})
	}
}

func TestNilID(t *testing.T) {
	got := domain.NilID()
	if got.Val() != uuid.Nil.String() {
		t.Errorf("実際の値は %v, 想定した値は %v", got.Val(), uuid.Nil.String())
	}
}

func TestIsNilID(t *testing.T) {
	a := domain.NilID()
	if a.IsNilID() != true {
		t.Errorf("NilIDのはずなのに 判定メソッドが NilIDじゃない となっています")
	}
	b, _ := domain.NewID()
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
			name: "UUID v7文字列を渡して生成できるか?",
			args: args{data: "018f2f4e-8c1d-7b33-a1aa-4c0b0f0e1d01"},
			want: want{first: "018f2f4e-8c1d-7b33-a1aa-4c0b0f0e1d01"},
			err:  err{hasErr: false, msg: ""},
		},
		{
			name: "空文字でエラーになるか?",
			args: args{data: ""},
			want: want{first: "", second: ""},
			err:  err{hasErr: true, msg: "ID 生成時の値が存在しません"},
		},
		{
			name: "UUIDでない文字列でエラーになるか?",
			args: args{data: "not-uuid"},
			want: want{first: "", second: ""},
			err:  err{hasErr: true, msg: "ストレージの ID フォーマットが不正です"},
		},
		{
			name: "UUID v4ではエラーになるか?",
			args: args{data: "550e8400-e29b-41d4-a716-446655440000"},
			want: want{first: "", second: ""},
			err:  err{hasErr: true, msg: "ID は UUID v7 の必要があります"},
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

			if (err != nil) != tt.err.hasErr {
				t.Errorf("実際の error = %v, hasErr %v", err, tt.err.hasErr)
				return
			}

			if got.Val() != tt.want.first {
				t.Errorf("実際の値は %v, 想定した値は %v", got.Val(), tt.want.first)
			}
		})
	}
}
