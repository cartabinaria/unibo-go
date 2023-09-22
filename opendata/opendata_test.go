package opendata

import (
	"reflect"
	"testing"
)

func TestResources_GetByAlias(t *testing.T) {

	singleAliasRes := Resource{Alias: "alias1"}
	multipleAliasRes := Resource{Alias: "alias1, alias2"}

	type args struct {
		alias string
	}
	tests := []struct {
		name  string
		r     Resources
		args  args
		want  *Resource
		want1 bool
	}{
		{name: "Single Alias", r: Resources{singleAliasRes}, args: args{alias: "alias1"}, want: &singleAliasRes, want1: true},
		{name: "Multiple Alias", r: Resources{multipleAliasRes}, args: args{alias: "alias1"}, want: &multipleAliasRes, want1: true},
		{name: "Multiple Alias in wrong order", r: Resources{multipleAliasRes}, args: args{alias: "alias2"}, want: &multipleAliasRes, want1: true},
		{name: "Multiple Alias not found", r: Resources{multipleAliasRes}, args: args{alias: "alias3"}, want: nil, want1: false},
		{name: "Multiple resources", r: Resources{singleAliasRes, multipleAliasRes}, args: args{alias: "alias2"}, want: &multipleAliasRes, want1: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.r.GetByAlias(tt.args.alias)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByAlias() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetByAlias() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
