package parsing

import "testing"

func TestFloat64MustParse(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "Success",
			args: args{
				s: "123",
			},
			want: 123.0,
		},
		{
			name: "Fail",
			args: args{
				s: "abc",
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Float64MustParse(tt.args.s); got != tt.want {
				t.Errorf("Float64MustParse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntMustParse(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Success",
			args: args{
				s: "123",
			},
			want: 123,
		},
		{
			name: "Fail",
			args: args{
				s: "abc",
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IntMustParse(tt.args.s); got != tt.want {
				t.Errorf("IntMustParse() = %v, want %v", got, tt.want)
			}
		})
	}
}
