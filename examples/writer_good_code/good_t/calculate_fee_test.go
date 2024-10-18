package good_t

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_calculateFee(t *testing.T) {
	type args struct {
		transactionNum int
	}
	tests := []struct {
		name         string
		args         args
		wantTotalFee int
	}{
		{
			name: "",
			args: args{
				transactionNum: 5,
			},
			wantTotalFee: 100,
		},
		{
			name: "",
			args: args{
				transactionNum: 10,
			},
			wantTotalFee: 100 + 5*10,
		},
		{
			name: "",
			args: args{
				transactionNum: 6,
			},
			wantTotalFee: 110,
		},
		{
			name: "",
			args: args{
				transactionNum: 30,
			},
			wantTotalFee: 100 + 150 + (30 - 20),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.wantTotalFee, calculateFee(tt.args.transactionNum), "calculateFee(%v)", tt.args.transactionNum)
		})
	}
}
