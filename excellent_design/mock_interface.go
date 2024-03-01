package excellent_design

type OrderAPI interface {
	GetOrderId() string
}

type realOrderImpl struct{}

func (r *realOrderImpl) GetOrderId() string {
	return ""
}

type mockOrderImpl struct{}

func (m *mockOrderImpl) GetOrderId() string {
	return "mock"
}
