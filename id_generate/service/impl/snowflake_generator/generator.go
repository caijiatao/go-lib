package snowflake_generator

import (
	"math/rand"
	"sync"
	"time"
)

const (
	// NodeBits holds the number of bits to use for Node
	// 机器编号，初始化时随机在[0,2^10)中选择
	NodeBits uint8 = 10
	// StepBits holds the number of bits to use for Step
	// 每毫秒最多能生成2^12=4096个id
	StepBits uint8 = 12

	nodeMax   = -1 ^ (-1 << NodeBits)
	stepMask  = -1 ^ (-1 << StepBits)
	timeShift = NodeBits + StepBits
	nodeShift = StepBits
)

// A Node struct holds the basic information needed for a snowflake generator
// snowflakeGenerator
type Node struct {
	mu   sync.Mutex
	time int64
	node int64
	step int64
}

type ID int64

var (
	snowflakeGenerator *Node
)

// NewSnowflakeGenerateService returns a new snowflake Node that can be used to generate snowflake IDs
func NewSnowflakeGenerateService() *Node {
	snowflakeGenerator = &Node{
		time: 0,
		node: int64(rand.Intn(nodeMax + 1)),
		step: 0,
	}
	return snowflakeGenerator
}

// Generate creates and returns a unique snowflake ID
// To help guarantee uniqueness
// - Make sure your system is keeping accurate system time
// - Make sure you never have multiple nodes running with the same snowflakeGenerator ID
func (n *Node) Generate() ID {

	n.mu.Lock()
	defer n.mu.Unlock()
	now := time.Now().UnixNano() / 1000000

	if now == n.time {
		n.step = (n.step + 1) & stepMask

		if n.step == 0 {
			for now <= n.time {
				now = time.Now().UnixNano() / 1000000
			}
		}
	} else {
		// 机器编号Node随机生成，有极小概率重复
		// QPS低，随机step减小Node号相同造成id重复的概率
		// 参考QPS: 2022 99大促id地区QPS约为20，每毫秒内基本不会多次生成id
		n.step = int64(rand.Intn(stepMask + 1))
	}

	n.time = now

	r := ID((now)<<timeShift |
		(n.node << nodeShift) |
		(n.step),
	)

	return r
}

// Int64 returns an int64 of the snowflake ID
func (f ID) Int64() int64 {
	return int64(f)
}

// Uint64 returns an int64 of the snowflake ID
func (f ID) Uint64() uint64 {
	return uint64(f)
}
