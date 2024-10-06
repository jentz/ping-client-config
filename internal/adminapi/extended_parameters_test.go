package adminapi

import (
	pfclient "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"testing"
)

func TestGetSingleExtendedParameterValue(t *testing.T) {
	type args struct {
		parameters   *map[string]pfclient.ParameterValues
		key          string
		defaultValue string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "single value",
			args: args{&map[string]pfclient.ParameterValues{
				"key": {Values: []string{"value"}},
			},
				"key",
				"default",
			},
			want: "value"},
		{name: "no value",
			args: args{&map[string]pfclient.ParameterValues{},
				"key",
				"default",
			},
			want: "default",
		},
		{name: "nil map",
			args: args{nil,
				"key",
				"default",
			},
			want: "default",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSingleExtendedParameterValue(tt.args.parameters, tt.args.key, tt.args.defaultValue); got != tt.want {
				t.Errorf("GetSingleExtendedParameterValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
