---
title: "注释篇"
weight: 1
# bookFlatSection: false
# bookToc: true
# bookHidden: false
# bookCollapseSection: false
# bookComments: false
# bookSearchExclude: false
---
看 `Kubernetes` 代码的过程中，不断回忆起之前看的代码封装相关的书籍，比如 《重构》、《代码整洁之道》和《设计模式》等，发现在 `Kubernetes` 不断在践行这些书籍里面的理论和技巧。

也正因如此，在读代码的过程中能通过符合直觉的方式去推导想要了解的内容，也能快速了解到代码的意图。Kubernetes 源码里面，优秀的注释和变量命名也是帮助开发者更好了解代码设计的意图。那 Kubernetes 源码中注释和变量命名有哪些值得我们学习的呢？

# 变量篇

### 变量名不是越详细越好

**如果变量名要精确表达具体的意思，势必会遇到一个问题，就是长度会太长。**但是，当一个非常长的变量名反复出现在代码中的时候，这种感觉就像是 别洛焦尔斯基 和 特维尔斯基 这么多个“司机”，我们看起来也会十分的头疼。

为了避免这种过分精确的重复命名带来的困惑，我们可以通过上下文的语义，**帮助短小精悍的变量名表达更多的含义。**

```go
func (q *graceTerminateRSList) remove(rs *listItem) bool{
		//...
}
```

在 `Kubernetes` 的  `graceTerminateRSList`  结构体的定义，我们就不需要写`graceTerminateRealServerList` ，因为对应上下文的 `listItem` 的时候内部已经用了全称定义。

```go
type listItem struct {
	VirtualServer *utilipvs.VirtualServer
	RealServer    *utilipvs.RealServer
}
```

所以在这个语境下面， `rs` 对应的只会是 `realServer` 而不会是 `replicaSet` 或者是其他的，**如果有这种歧义的可能性，那么就不能进行这种缩写。**

还有`graceTerminateRSList`  的移除 `rs` 方法，我们不需要写 `removeRS` 或者是 `removeRealServer` ，在传入的参数签名就已经有了 `rs *listItem` ，所以这个方法移除的只能是 `rs` ，在方法名加上 `rs` 反而显得冗余了。

**在起名字的时候，尽可能短的命名承载更多的意思。**

```go
func CountNumber(nums []int, n int) (count int) {
	for i := 0; i < len(nums); i++ {
		// 如果要赋值 则 v := nums[i]
		if nums[i] == n {
			count++
		}
	}
	return
}

func CountNumberBad(nums []int, n int) (count int) {
	for index := 0; index < len(nums); index++ {
		value := nums[index]
		if value == n {
			count++
		}
	}
	return
}

```



`index`并不比 `i` 承载了更多信息， `value` 也不比 `v` 更好，所以，在这个例子中，是可以使用缩写来代替的。但是，**缩写也并不完全带来好处，需要看具体的场景是否会产生歧义而确定。**

### 变量名需要避免理解歧义

想要表达参加活动的用户数（  `int` 类型），那么用 `userCount` 比用 `user` 或者 `users` 更好。**因为** `user` **可以表示用户信息的对象，而** `users` **可能是用户信息的切片，如果使用这两个会给用户带来歧义。**

再看一个例子。 `min` 在某些情况下可以表示最小值(minimum)，也可以表示分钟(minutes)，如果我们在一些比较容易混淆的场景下，我们就用全拼来代替缩写。

```go
// 计算最小价格和促销活动的剩余时间
func main() {
	// 商品价格列表
	prices := []float64{12.99, 9.99, 15.99, 8.49}

	// 各个商品促销剩余时间（分钟）
	remainingMinutes := []int{30, 45, 10, 20}

	// min := findMinPrice(prices) // 变量名 "min"：表示最小价格
	minPrice := findMinPrice(prices) 
	fmt.Printf("商品的最低价格: $%.2f\n", min)

	// min = findMinTime(remainingMinutes) // 变量名 "min"：表示剩余的最短时间
	remainingMinute := findMinTime(remainingMinutes) 
	fmt.Printf("促销活动的最短剩余时间: %d minutes\n", min)
}
```

在这个例子中， `min` 不仅能表示商品最低价格，还能表示活动剩余的最小分钟数，所以在这种情况下我们就不要使用缩写。这样我们就能明确的区分出找到的是最小价格还是最小分钟数。

### 相同含义的变量需保持一致

**整个项目中代表相同含义的变量名字应该尽量保持一致**，如果项目里把用户 `id` 写成 `UserId` ，那么在其他地方进行复制的时候，就不要把他改成 `Uid` ，这样我们就会疑惑 `UserId` 和 `Uid` 是否是同个东西。

不要小看这种情况，有的时候因为多个系统都存在用户id，我们可能都需要进行存储，如果不加前缀进行区分，那么在需要使用的时候会无从下手。

比如用户A作为买家的角色，买了卖家的某个商品，并且由骑手进行派送。

这里就出现了三个用户id：买家、卖家和骑手。

这个时候我们可以通过增加模块前缀的方式来区分：**BuyerId 、SellerId和DriverId。**

**这个时候我们也尽量不要进行缩写，因为这三个已经足够简短，如果我们把函数入参的 SellerId 缩写成 Sid 的话，当后面的需求有一个商铺id（ShopId）的概念，这个时候我们就会疑惑，Sid 是对应的 `SellerId` 还是 `ShopId` 。** 如果有个人卖家就是通过 `SellerId` 填充了 `ShopId` ，那这个时候就会造成线上的 BUG。
