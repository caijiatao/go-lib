package good_t

type DiscountStrategy interface {
	Calculate(price float64) float64
}

type RegularUserDiscount struct{}

func (r *RegularUserDiscount) Calculate(price float64) float64 {
	return price // 无折扣
}

type VIPUserDiscount struct{}

func (v *VIPUserDiscount) Calculate(price float64) float64 {
	return price * 0.8 // 8折
}

type SVIPUserDiscount struct{}

func (s *SVIPUserDiscount) Calculate(price float64) float64 {
	return price * 0.5 // 5折
}

type CustomerType int

const (
	Regular CustomerType = iota + 1
	VIP
	SVIP
)

var (
	DiscountStrategyMap = map[CustomerType]DiscountStrategy{
		Regular: &RegularUserDiscount{},
		VIP:     &VIPUserDiscount{},
		SVIP:    &SVIPUserDiscount{},
	}
)

type User struct {
	CustomerType CustomerType
}

func CalculatePrice(user User, price float64) float64 {
	strategy := DiscountStrategyMap[user.CustomerType]
	return strategy.Calculate(price)
}

type fullDiscount struct {
	targetPrice float64
	discount    float64
}

func NormalCalculatePrice(user User, price float64) float64 {
	// 8折
	if user.CustomerType == VIP {
		return price * 0.8
	}
	// 5折
	if user.CustomerType == SVIP {
		discountPrice := price * 0.5
		if discountPrice > 300 {
			return discountPrice - 30
		}
		return discountPrice
	}
	return price
}

func VIPDiscountCalculate(price float64) float64 {
	return price * 0.8
}

func RegularDiscountCalculate(price float64) float64 {
	return price * 0.9
}

func SVIPDiscountCalculate(price float64) float64 {
	return price * 0.5
}

type CalculateHandle func(float64) float64

func (c CalculateHandle) Calculate(price float64) float64 {
	return c(price)
}

var (
	calculateFunc = map[CustomerType]CalculateHandle{
		VIP:     VIPDiscountCalculate,
		SVIP:    SVIPDiscountCalculate,
		Regular: RegularDiscountCalculate,
	}
)

func CalculatePriceF(user User, price float64) float64 {
	handle, ok := calculateFunc[user.CustomerType]
	if !ok {
		return price
	}
	return handle.Calculate(price)
}

type Discount struct {
	// 折扣率
	DiscountRate float64
}

func (d *Discount) Calculate(price float64) float64 {
	return price * d.DiscountRate
}

type FullDiscount struct {
	// 满多少钱可以减
	TargetPrice float64
	// 减多少钱
	Discount float64
}

func (f *FullDiscount) Calculate(price float64) float64 {
	if price >= f.TargetPrice {
		return price - f.Discount
	}
	return price
}

func getDiscounts(customerType CustomerType) []DiscountStrategy {
	if customerType == VIP {
		return []DiscountStrategy{
			&Discount{DiscountRate: 0.8},
			&FullDiscount{TargetPrice: 500, Discount: 30},
		}
	}

	if customerType == SVIP {
		return []DiscountStrategy{
			&Discount{DiscountRate: 0.5},
			&FullDiscount{TargetPrice: 300, Discount: 30},
		}
	}

	return []DiscountStrategy{
		&Discount{DiscountRate: 1},
	}
}

func Calculate(user User, price float64, discount Discount) float64 {
	// 获取折扣策略，这个可以支持配置
	strategies := getDiscounts(user.CustomerType)

	// 给价格应用折扣策略
	for _, strategy := range strategies {
		price = strategy.Calculate(price)
	}

	return price
}
