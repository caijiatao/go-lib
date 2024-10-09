package interface_demo

type Cart struct {
	Total float64
	Items []string
}

type Promotion interface {
	// Apply 在购物车上应用促销策略，返回的是新的购物车实体
	Apply(cart Cart) Cart
}

// DiscountByAmount 满减促销
type DiscountByAmount struct {
	Threshold float64
	Discount  float64
}

func (d DiscountByAmount) Apply(cart Cart) Cart {
	if cart.Total >= d.Threshold {
		cart.Total -= d.Discount
	}
	return cart
}

// DiscountByPercentage 打折促销
type DiscountByPercentage struct {
	Percentage float64
}

func (d DiscountByPercentage) Apply(cart Cart) Cart {
	cart.Total *= (1 - d.Percentage)
	return cart
}

// FreeGift 赠品促销
type FreeGift struct {
	Gift string
}

func (g FreeGift) Apply(cart Cart) Cart {
	cart.Items = append(cart.Items, g.Gift)
	return cart
}

// ApplyPromotion 应用在活动期限的促销策略
func ApplyPromotion(cart Cart) Cart {
	promotions := GetEffectPromotion()

	for _, promotion := range promotions {
		cart = promotion.Apply(cart)
	}

	return cart
}

func GetEffectPromotion() []Promotion {
	promotions := make([]Promotion, 0)

	// 满减促销
	discountByAmount := DiscountByAmount{Threshold: 100, Discount: 20}
	promotions = append(promotions, discountByAmount)

	// 打折促销
	discountByPercentage := DiscountByPercentage{Percentage: 0.1}
	promotions = append(promotions, discountByPercentage)

	// 赠品促销
	freeGift := FreeGift{Gift: "Free Mug"}
	promotions = append(promotions, freeGift)

	// 如果有其他策略可以在这里初始化

	return promotions
}
