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
