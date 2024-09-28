# Go 反射性能究竟如何？

Tags: Go
Last edited time: May 29, 2024 5:05 PM
状态: 已发布

# **一.为什么需要反射？**

首先要明白反射能有什么好处，如果它不能带来任何好处，那么实际上我们不用它也就不需要担心性能带来的影响。

> Go 语言反射的实现原理
> 

Go 语言的语法元素很少、设计简单，所以它没有特别强的表达能力，但是 Go 语言的 `reflect` 包能够弥补它在语法上的一些劣势。

反射能减少重复的编码工作，工具包通过反射来处理不同结构体入参。比如当我们不知道具体数据类型的时候，反射可以帮助我们处理这部分情况，通过 `reflect` 包能够确定数据的类型来做不同的分支处理。

序列化和反序列化库依赖反射来将结构体转换为其他格式，或者从这些格式转换回结构体。例如平时使用的很多的`json.Marshal` 和 `json.Unmarshal` ，反射可以遍历结构体的字段，并根据字段的类型和标签进行相应的转换。

# 二.应用场景

## **2.1 反射判断结构体是否为空**

**业务场景：通过反射的方式可以在传入的结构体为空直接进行返回，而不对SQL 拼接，从而避免了全表扫描有慢SQL**

如果不使用反射，那么当需要判断一个结构体是否为空时则需要一个字段一个字段进行判断，实现如下：

```go
type aStruct struct {
	Name stringMale string}

func (s *aStruct) IsEmpty() bool {
	return s.Male == "" && s.Name == ""
}

type complexSt struct {
	A        aStruct
	S        []stringIntValue int
}

func (c *complexSt) IsEmpty() bool {
	return c.A.IsEmpty() && len(c.S) == 0 && c.IntValue == 0
}
```

这时如果需要增加一个新结构体的判空，那么就需要实现对应的方法对每个字段进行判断

我们来具体看一下反射能够怎么实现。

首先通过反射来对每个字段类型先进行判空定义

```go
func intValueIsEmpty(v reflect.Value) bool {
	return v.Int() == 0
}

func uintValueIsEmpty(v reflect.Value) bool {
	return v.Uint() == 0
}

func completeValueIsEmpty(v reflect.Value) bool {
	return v.Complex() == 0
}

func stringValueIsEmpty(v reflect.Value) bool {
	return len(v.String()) == 0
}

func sliceValueIsEmpty(v reflect.Value) bool {
	return v.IsNil() || v.Len() == 0
}

// 如果有其他想增加的判断可以直接增加方法实现
```

将对应的空值映射成 `map` 方便使用

```go
type valueIsEmptyFuncfunc(v reflect.Value)boolvar valueIsEmptyFuncMap =map[reflect.Kind]valueIsEmptyFunc{
	reflect.Int:        intValueIsEmpty,
	reflect.Int8:       intValueIsEmpty,
	reflect.Int16:      intValueIsEmpty,
	reflect.Int32:      intValueIsEmpty,
	reflect.Int64:      intValueIsEmpty,
	reflect.Uint:       uintValueIsEmpty,
	reflect.Uint8:      uintValueIsEmpty,
	reflect.Uint16:     uintValueIsEmpty,
	reflect.Uint32:     uintValueIsEmpty,
	reflect.Uint64:     uintValueIsEmpty,
	reflect.Complex64:  completeValueIsEmpty,
	reflect.Complex128: completeValueIsEmpty,
	reflect.Array:      sliceValueIsEmpty,
	reflect.Slice:      sliceValueIsEmpty,
	reflect.String:     stringValueIsEmpty,
}
```

判断是否为空的结构体通过将入参定义成 `interface` 来实现通用

```go
func IsStructEmpty(v interface{}) bool {
	vType := reflect.TypeOf(v)
	value := reflect.ValueOf(v)
	for i := 0; i < value.NumField(); i++ {
		field := vType.Field(i)
		// 结构体变量递归判断
		if field.Type.Kind() == reflect.Struct && field.Type.NumField() > 0 {
			if isEmpty := IsStructEmpty(value.FieldByName(field.Name).Interface()); !isEmpty {
				return false
			}
		}
		if isEmptyFunc, ok := valueIsEmptyFuncMap[field.Type.Kind()]; ok {
			if isEmpty := isEmptyFunc(value.Field(i)); !isEmpty {
				return false
			}
		}
	}
	return true
}
```

这个时候只需要传入对应的结构体就可以得到对应的数据是否为空，不需要重复进行实现。

### **性能对比**

```go
func BenchmarkReflectIsStructEmpty(b *testing.B) {
	s := complexSt{
		A:        aStruct{},
		S:        make([]string, 0),
		IntValue: 0,
	}
	for i := 0; i < b.N; i++ {
		IsStructEmpty(s)
	}
}

func BenchmarkNormalIsStructEmpty(b *testing.B) {
	s := complexSt{
		A:        aStruct{},
		S:        make([]string, 0),
		IntValue: 0,
	}
	for i := 0; i < b.N; i++ {
		s.IsEmpty()
	}
}
```

执行性能测试

```go
# -benchmem 查看每次分配内存的次数
# -benchtime=3s 执行的时间指定为3s，一般1s、3s、5s得到的结果差不多的，如果性能较差，执行时间越长得到的性能平均值越准确
# -count=3 执行次数，多次执行能保证准确性
# -cpu n 指定cpu的核数，一般情况下CPU核数增加会提升性能，但也不是正相关的关系，因为核数多了之后上下文切换会带来影响，需要看是IO密集型还是CPU密集型的应用，多协程的测试中可以进行对比
go test -bench="." -benchmem -cpuprofile=cpu_profile.out -memprofile=mem_profile.out -benchtime=3s -count=3 .
```

执行结果：

```go
BenchmarkReflectIsStructEmpty-16                 8126697               493 ns/op             112 B/op          7 allocs/op
BenchmarkReflectIsStructEmpty-16                 6139268               540 ns/op             112 B/op          7 allocs/op
BenchmarkReflectIsStructEmpty-16                 7222296               465 ns/op             112 B/op          7 allocs/op
BenchmarkNormalIsStructEmpty-16                 1000000000               0.272 ns/op           0 B/op          0 allocs/op
BenchmarkNormalIsStructEmpty-16                 1000000000               0.285 ns/op           0 B/op          0 allocs/op
BenchmarkNormalIsStructEmpty-16                 1000000000               0.260 ns/op           0 B/op          0 allocs/op
```

### **结果分析**

结果字段的含义：

| **结果项** | **含义** |
| --- | --- |
| BenchmarkReflectIsStructEmpty-16 | BenchmarkReflectIsStructEmpty 是测试的函数名-16 表示GOMAXPROCS （线程数）的值为16 |
| 2899022 | 一共执行了2899022次 |
| 401 ns/op | 表示平均每次操作花费了401纳秒 |
| 112 B/op | 表示每次操作申请了112 Byte的内存 |
| 7 allocs/op | 表示申请了七次内存 |

反射判断每次操作的耗时大约是直接判断的1000倍，且带来了额外七次的内存分配，每次会增加112Byte，这样看下来性能比直接操作还是会下降不少的。

## **2.2 反射复制结构体同名字段**

在实际业务接口中我们经常需要对数据进行 `DTO` 和 `VO` 的转换，并且大部分时候是同名字段的复制，这个时候如果不使用反射则需要对每个字段进行复制，并且在新产生一个结构体需要复制时，则需要再重复进行 如下`new` 方法的编写，会带来大量的重复工作：

```go
type aStruct struct {
	Name stringMale string
}

type aStructCopy struct {
	Name stringMale string
}

func newAStructCopyFromAStruct(a *aStruct) *aStructCopy {
	return &aStructCopy{
		Name: a.Name,
		Male: a.Male,
	}
}
```

使用反射来对结构体进行复制，在有需要复制的新结构体时我们只需要将结构体指针传入即可进行同名字段的复制，实现如下：

```go
func CopyIntersectionStruct(src, dst interface{}) {
	sElement := reflect.ValueOf(src).Elem()
	dElement := reflect.ValueOf(dst).Elem()
	for i := 0; i < dElement.NumField(); i++ {
		dField := dElement.Type().Field(i)
		sValue := sElement.FieldByName(dField.Name)
		if !sValue.IsValid() {
			continue}
		value := dElement.Field(i)
		value.Set(sValue)
	}
}
```

### **性能对比**

Benchmark Test的代码如下：

```go
func BenchmarkCopyIntersectionStruct(b *testing.B) {
	a := &aStruct{
		Name: "test",
		Male: "test",
	}
	for i := 0; i < b.N; i++ {
		var ac aStructCopy
		CopyIntersectionStruct(a, &ac)
	}
}

func BenchmarkNormalCopyIntersectionStruct(b *testing.B) {
	a := &aStruct{
		Name: "test",
		Male: "test",
	}
	for i := 0; i < b.N; i++ {
		newAStructCopyFromAStruct(a)
	}
}
```

**运行性能测试**

```go
go test -bench="." -benchmem -cpuprofile=cpu_profile.out -memprofile=mem_profile.out -benchtime=3s -count=3 .
```

### 结果分析

```go
BenchmarkCopyIntersectionStruct-16              10787202               352 ns/op              64 B/op          5 allocs/op
BenchmarkCopyIntersectionStruct-16              10886558               304 ns/op              64 B/op          5 allocs/op
BenchmarkCopyIntersectionStruct-16              10147404               322 ns/op              64 B/op          5 allocs/op
BenchmarkNormalCopyIntersectionStruct-16        1000000000               0.277 ns/op           0 B/op          0 allocs/op
BenchmarkNormalCopyIntersectionStruct-16        1000000000               0.270 ns/op           0 B/op          0 allocs/op
BenchmarkNormalCopyIntersectionStruct-16        1000000000               0.259 ns/op           0 B/op          0 allocs/op
```

**与上面第一个运行结果相差无几**，反射的耗时**仍然是不使用反射的1000倍，内存分配也在每次多增加了64Byte**

在实际的业务场景中可能多次反射的组合使用，如果是需要对实际性能可以自行编写 `BenchmarkTest` 进行测试

火焰图对比可以更加明确的看出运行时间的占比

![Untitled](../../../Go%20%E5%8F%8D%E5%B0%84%E6%80%A7%E8%83%BD%E7%A9%B6%E7%AB%9F%E5%A6%82%E4%BD%95%EF%BC%9F%20d399f6992d11492fae0f4f5515d6669f/Untitled.png)

![Untitled](../../../Go%20%E5%8F%8D%E5%B0%84%E6%80%A7%E8%83%BD%E7%A9%B6%E7%AB%9F%E5%A6%82%E4%BD%95%EF%BC%9F%20d399f6992d11492fae0f4f5515d6669f/Untitled%201.png)

# **三.结论**

业务接口中我们假设接口的相应是10ms，一个反射方法的操作平均是 **400纳秒，**会带来的额外内存分配大概是 **64Byte~112Byte。**

> 1ms【毫秒】 = 1000μs【微秒】=1000 * 1000ns【纳秒】
> 

如果一个接口在链路上做了1000次反射的操作，单次操作大约会增加0.4ms的接口延时，一般单次请求中经过的中间件和业务操作也很少会达到这个次数，**所以这种响应时长的影响基本可以忽略不计。**实际业务中更多的损耗则是会在内存的复制和网络IO上。

但是反射在编码上也存在实实在在的问题就是维护起来会比普通业务代码更加困难，理解上会更加费劲，所以在使用时需要进行斟酌，避免过度使用导致代码的复杂度不断提高

# 四.推荐阅读

通过反射来实现具体类型的相互转换，也能通过反射来为变量设置具体的参数值，还有其他更多的用法可以参考《Go语言设计与实现》，里面也有关于Go `reflect` 的分析。

比如 `reflect.Type` 是一个类型的接口

```go
type Type interface {
        Kind() string
        ...
        Implements(u Type) bool
        ...
}
```

提供了类型操作有可能遇到的方法，比如 `Implements` 就可以用来判断一个结构体是否有实现某个结构。在使用的时候我们可以通过 `reflect.TypeOf` 获取任意变量的类型

```go
type test interface {
	test()
}

type A struct{}

func (receiver *A) test() {
	println("test")
}

func main() {
	a := &A{}
	aType := reflect.TypeOf(a)

	fmt.Println(aType.Kind())
	fmt.Println(aType.Implements(reflect.TypeOf((*test)(nil)).Elem()))
}

```

感谢你读到这里，如果喜欢云原生、Go、个人成长的内容可以关注我，我也会不定期推荐我自己读过的好书，让我们一起进步。