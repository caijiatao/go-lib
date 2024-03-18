package excellent_design

import (
	"fmt"
	"testing"
)

func TestGetOrderId(t *testing.T) {
	orderAPI := &mockOrderImpl{} // 如果要获取订单id，且不是测试的重点，这里直接初始化成mock的结构体
	fmt.Println(orderAPI.GetOrderId())
}

type testResp struct {
	code string
}

func s() (resp *testResp) {
	resp.code = "200"
	return resp
}

func TestResp(t *testing.T) {
	s()
}
