// TODO: 作らないといけないけど 面倒くさい笑
package share_test

import (
	"testing"

	"github.com/ozaki-physics/raison-me/info/authN/infra"
	"github.com/ozaki-physics/raison-me/share"
)

func TestJsonToStruct(t *testing.T) {
	type user struct {
		AccountID string `json:"account_id"`
		Id        string `json:"id"`
		Name      string `json:"name"`
	}
	type jsonUser struct {
		User []user `json:"user"`
	}
	var ju jsonUser

	type args struct {
		goStruct interface{}
		path     string
	}

	path := infra.NewStoragePath()

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "JSONを構造体に変換できるか",
			args: args{
				ju,
				path.GetUserPath(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := share.JsonToStruct(tt.args.goStruct, tt.args.path)
			t.Errorf("%+v\n", tt.args.goStruct)
			if (err != nil) != tt.wantErr {
				t.Errorf("JsonToStruct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
