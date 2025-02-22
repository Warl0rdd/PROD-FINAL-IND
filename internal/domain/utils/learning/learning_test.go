package learning

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"solution/internal/adapters/database/postgres"
	"testing"
)

func TestGenNewR0(t *testing.T) {
	type args struct {
		oldR0 float64
		data  []postgres.GetImpressionsForLearningRow
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "New R0 Clicked",
			args: args{
				oldR0: 0.5,
				data: []postgres.GetImpressionsForLearningRow{
					{
						Score:        0.639,
						ClickedAfter: true,
					},
				},
			},
			want: 0.48,
		},
		{
			name: "New R0 didn't click",
			args: args{
				oldR0: 0.5,
				data: []postgres.GetImpressionsForLearningRow{
					{
						Score:        0.639,
						ClickedAfter: false,
					},
				},
			},
			want: 0.58,
		},
		{
			name: "Old R0",
			args: args{
				oldR0: 0.5,
				data:  []postgres.GetImpressionsForLearningRow{},
			},
			want: 0.5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := cmp.Options{
				cmpopts.EquateApprox(0, 0.01),
			}

			if diff := cmp.Diff(GenNewR0(tt.args.oldR0, tt.args.data), tt.want, opts); diff != "" {
				t.Errorf("actual = %v, want %v", diff, tt.want)
			}
		})
	}
}

func BenchmarkGenNewR0(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenNewR0(0.5, []postgres.GetImpressionsForLearningRow{
			{
				Score:        0.639,
				ClickedAfter: true,
			},
			{
				Score:        0.638,
				ClickedAfter: false,
			},
			{
				Score:        0.7,
				ClickedAfter: true,
			},
			{
				Score:        0.1,
				ClickedAfter: false,
			},
			{
				Score:        0.2,
				ClickedAfter: true,
			},
			{
				Score:        0.87,
				ClickedAfter: false,
			},
			{
				Score:        0.5,
				ClickedAfter: true,
			},
			{
				Score:        0.9,
				ClickedAfter: false,
			},
			{
				Score:        0.33,
				ClickedAfter: true,
			},
			{
				Score:        0.1,
				ClickedAfter: false,
			},
			{
				Score:        0.123,
				ClickedAfter: true,
			},
			{
				Score:        0.543,
				ClickedAfter: false,
			},
			{
				Score:        0.123123,
				ClickedAfter: true,
			},
			{
				Score:        0.76,
				ClickedAfter: false,
			},
			{
				Score:        0.642,
				ClickedAfter: true,
			},
		})
	}
}
