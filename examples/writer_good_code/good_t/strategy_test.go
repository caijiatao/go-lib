package good_t

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCalculatePriceF(t *testing.T) {
	type args struct {
		user  User
		price float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "",
			args: args{
				user: User{
					CustomerType: SVIP,
				},
				price: 100,
			},
			want: 50,
		},
		{
			name: "",
			args: args{
				user: User{
					CustomerType: VIP,
				},
				price: 100,
			},
			want: 80,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, CalculatePriceF(tt.args.user, tt.args.price), "CalculatePriceF(%v, %v)", tt.args.user, tt.args.price)
		})
	}
}
