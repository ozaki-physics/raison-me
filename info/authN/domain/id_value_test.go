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
