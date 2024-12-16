package consistent_hash

import (
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

type HashFunc func(data []byte) uint32

type ConsistentHashOpt func(*ConsistentHash)

type ConsistentHash struct {
	hashFunc HashFunc // 自定义哈希函数

	replicas int               // 每个真实节点的虚拟节点数
	ring     []uint32          // 哈希环，存储虚拟节点的哈希值
	nodes    map[uint32]string // 哈希值到真实节点的映射

	mu sync.RWMutex // 保护数据的一致性
}

func WithHashFunc(f HashFunc) ConsistentHashOpt {
	return func(ch *ConsistentHash) {
		ch.hashFunc = f
	}
}

func WithReplica(replicas int) ConsistentHashOpt {
	return func(ch *ConsistentHash) {
		ch.replicas = replicas
	}
}

func (ch *ConsistentHash) getVirtualNodeKey(node string, replicaIndex int) []byte {
	return []byte(node + "#" + strconv.Itoa(replicaIndex))
}

// NewConsistentHash 创建一致性哈希对象
func NewConsistentHash(opts ...ConsistentHashOpt) *ConsistentHash {
	consistentHash := &ConsistentHash{
		replicas: 3,
		hashFunc: crc32.ChecksumIEEE,
		ring:     []uint32{},
		nodes:    make(map[uint32]string),
	}

	for _, opt := range opts {
		opt(consistentHash)
	}

	return consistentHash
}

// AddNode 添加真实节点并生成虚拟节点
func (ch *ConsistentHash) AddNode(node string) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	// 为每个真实节点生成虚拟节点
	for i := 0; i < ch.replicas; i++ {
		virtualNodeKey := ch.getVirtualNodeKey(node, i)
		hashValue := ch.hashFunc(virtualNodeKey)
		ch.ring = append(ch.ring, hashValue)
		ch.nodes[hashValue] = node
	}

	// 方便后面查找 key 对应的节点索引
	sort.Slice(ch.ring, func(i, j int) bool {
		return ch.ring[i] < ch.ring[j]
	})
}

// RemoveNode 移除真实节点及其虚拟节点
func (ch *ConsistentHash) RemoveNode(node string) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	// 移除该真实节点对应的虚拟节点
	for i := 0; i < ch.replicas; i++ {
		virtualNodeKey := ch.getVirtualNodeKey(node, i)
		hashValue := ch.hashFunc(virtualNodeKey)
		delete(ch.nodes, hashValue)

		// 从环中删除
		for j, v := range ch.ring {
			if v == hashValue {
				ch.ring = append(ch.ring[:j], ch.ring[j+1:]...)
				break
			}
		}
	}
}

// GetNode 根据键找到对应的真实节点
func (ch *ConsistentHash) GetNode(key string) string {
	ch.mu.RLock()
	defer ch.mu.RUnlock()

	if len(ch.ring) == 0 {
		return ""
	}

	hashValue := ch.hashFunc([]byte(key))

	// 使用二分查找找到第一个大于或等于 hashValue 的虚拟节点
	idx := sort.Search(len(ch.ring), func(i int) bool {
		return ch.ring[i] >= hashValue
	})

	// 如果超出最大值，则环回到第一个节点
	if idx == len(ch.ring) {
		idx = 0
	}

	return ch.nodes[ch.ring[idx]]
}
