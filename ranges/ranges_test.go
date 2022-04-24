package ranges

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResolveIndexedRanges(t *testing.T) {
	alphabet := []string{"A", "B", "C", "D", "E", "F", "a", "b", "c", "d", "e", "f"}

	type args struct {
		inputs []string
		values []string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "no ranges",
			args: args{
				inputs: []string{"a", "c"},
				values: alphabet,
			},
			want:    []string{"a", "c"},
			wantErr: assert.NoError,
		},
		{
			name: "single range only",
			args: args{
				inputs: []string{"b-c"},
				values: alphabet,
			},
			want:    []string{"b", "c"},
			wantErr: assert.NoError,
		},
		{
			name: "double range",
			args: args{
				inputs: []string{"A-B", "b-d"},
				values: alphabet,
			},
			want:    []string{"A", "B", "b", "c", "d"},
			wantErr: assert.NoError,
		},
		{
			name: "range/literal mix",
			args: args{
				inputs: []string{"b", "A-B", "d"},
				values: alphabet,
			},
			want:    []string{"b", "A", "B", "d"},
			wantErr: assert.NoError,
		},
		{
			name: "range with whitespace",
			args: args{
				inputs: []string{"A - B"},
				values: alphabet,
			},
			want:    []string{"A", "B"},
			wantErr: assert.NoError,
		},
		{
			name: "zero length range",
			args: args{
				inputs: []string{"A- A"},
				values: alphabet,
			},
			want:    []string{"A- A"},
			wantErr: assert.NoError,
		},
		{
			name: "reverse range (wrap)",
			args: args{
				inputs: []string{"e-A"},
				values: alphabet,
			},
			want:    []string{"e", "f", "A"},
			wantErr: assert.NoError,
		},

		// invalid configurations
		{
			name: "open range",
			args: args{
				inputs: []string{"B-"},
				values: alphabet,
			},
			want:    []string{"B-"},
			wantErr: assert.NoError,
		},
		{
			name: "range with invalid values",
			args: args{
				inputs: []string{"My-Car"},
				values: alphabet,
			},
			want:    []string{"My-Car"},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveIndexedRanges(tt.args.inputs, tt.args.values)
			if !tt.wantErr(t, err, fmt.Sprintf("ResolveIndexedRanges(%v, %v)", tt.args.inputs, tt.args.values)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ResolveIndexedRanges(%v, %v)", tt.args.inputs, tt.args.values)
		})
	}
}

func TestResolveIntRanges(t *testing.T) {
	type args struct {
		inputs []string
	}
	tests := []struct {
		name    string
		args    args
		want    []int
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "no range",
			args: args{
				inputs: []string{"3"},
			},
			want:    []int{3},
			wantErr: assert.NoError,
		},
		{
			name: "simple range",
			args: args{
				inputs: []string{"3-5"},
			},
			want:    []int{3, 4, 5},
			wantErr: assert.NoError,
		},
		{
			name: "multiple ranges",
			args: args{
				inputs: []string{"3-5", "4-6"},
			},
			want:    []int{3, 4, 5, 4, 5, 6},
			wantErr: assert.NoError,
		},
		{
			name: "mixed",
			args: args{
				inputs: []string{"3-5", "6"},
			},
			want:    []int{3, 4, 5, 6},
			wantErr: assert.NoError,
		},
		// invalid inputs
		{
			name: "reverse range",
			args: args{
				inputs: []string{"5-3"},
			},
			want:    []int{},
			wantErr: assert.Error,
		},
		{
			name: "null range",
			args: args{
				inputs: []string{"5-5"},
			},
			want:    []int{},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveIntRanges(tt.args.inputs)
			if !tt.wantErr(t, err, fmt.Sprintf("ResolveIntRanges(%v)", tt.args.inputs)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ResolveIntRanges(%v)", tt.args.inputs)
		})
	}
}
