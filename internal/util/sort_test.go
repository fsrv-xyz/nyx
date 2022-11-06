package util_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fsrv-xyz/nyx/internal/check"
	"github.com/fsrv-xyz/nyx/internal/util"
)

func TestSortByCheckName(t *testing.T) {
	type args struct {
		checks []check.GenericCheck
	}
	tests := []struct {
		name string
		args args
		want []check.GenericCheck
	}{
		{
			name: "sorts checks by name",
			args: args{
				checks: []check.GenericCheck{
					{Name: "b"},
					{Name: "a"},
					{Name: "d"},
					{Name: "g"},
				},
			},
			want: []check.GenericCheck{
				{Name: "a"},
				{Name: "b"},
				{Name: "d"},
				{Name: "g"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := util.SortByCheckName(tt.args.checks); !reflect.DeepEqual(got, tt.want) {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
