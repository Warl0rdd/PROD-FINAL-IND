package ads

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"testing"
)

func TestAdScore(t *testing.T) {
	type args struct {
		cpi float64
		cpc float64
		rel float64
		r0  float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "AdScore",
			args: args{
				cpi: 2.0,
				cpc: 5.0,
				rel: 1.0,
				r0:  0.5,
			},
			want: 6.966,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := cmp.Options{
				cmpopts.EquateApprox(0, 0.01),
			}

			if diff := cmp.Diff(AdScore(tt.args.cpi, tt.args.cpc, tt.args.rel, tt.args.r0), tt.want, opts); diff != "" {
				t.Errorf("actual = %v, want %v", diff, tt.want)
			}
		})
	}
}

func TestLogistic(t *testing.T) {
	type args struct {
		x  float64
		r0 float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "Logistic",
			args: args{
				x:  0.6,
				r0: 0.5,
			},
			want: 0.731,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := cmp.Options{
				cmpopts.EquateApprox(0, 0.01),
			}

			if diff := cmp.Diff(Logistic(tt.args.x, tt.args.r0), tt.want, opts); diff != "" {
				t.Errorf("actual = %v, want %v", diff, tt.want)
			}
		})
	}
}

func TestNormalization(t *testing.T) {
	type args struct {
		rel float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "Normalization",
			args: args{
				rel: 1073741823.5,
			},
			want: 0.5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := cmp.Options{
				cmpopts.EquateApprox(0, 0.01),
			}

			if diff := cmp.Diff(Normalization(tt.args.rel), tt.want, opts); diff != "" {
				t.Errorf("actual = %v, want %v", diff, tt.want)
			}
		})
	}
}
