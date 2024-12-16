package consistent_hash

import (
	"hash/crc32"
	"testing"
)

// 测试创建一致性哈希对象
func TestNewConsistentHash(t *testing.T) {
	ch := NewConsistentHash(
		WithReplica(5),
		WithHashFunc(crc32.ChecksumIEEE),
	)

	if ch.replicas != 5 {
		t.Errorf("expected replicas = 5, got %d", ch.replicas)
	}

	if ch.hashFunc == nil {
		t.Error("hash function is not initialized")
	}
}

// 测试添加节点和获取节点
func TestAddNodeAndGetNode(t *testing.T) {
	ch := NewConsistentHash(WithReplica(3))

	ch.AddNode("NodeA")
	ch.AddNode("NodeB")
	ch.AddNode("NodeC")

	testKeys := []string{"Key1", "Key2", "Key3", "Key4", "Key5"}
	for _, key := range testKeys {
		node := ch.GetNode(key)
		if node == "" {
			t.Errorf("failed to get node for key %s", key)
		} else {
			t.Logf("Key %s -> Node %s", key, node)
		}
	}
}

// 测试删除节点
func TestRemoveNode(t *testing.T) {
	ch := NewConsistentHash(WithReplica(3))

	ch.AddNode("NodeA")
	ch.AddNode("NodeB")
	ch.AddNode("NodeC")

	// 移除 NodeB
	ch.RemoveNode("NodeB")

	testKeys := []string{"Key1", "Key2", "Key3", "Key4", "Key5"}
	for _, key := range testKeys {
		node := ch.GetNode(key)
		if node == "" {
			t.Errorf("failed to get node for key %s after removing NodeB", key)
		} else {
			t.Logf("Key %s -> Node %s", key, node)
		}
	}
}

// 测试虚拟节点数量是否正确
func TestVirtualNodeCount(t *testing.T) {
	ch := NewConsistentHash(WithReplica(3))

	ch.AddNode("NodeA")

	expectedVirtualNodes := 3
	actualVirtualNodes := 0
	for _, node := range ch.nodes {
		if node == "NodeA" {
			actualVirtualNodes++
		}
	}

	if actualVirtualNodes != expectedVirtualNodes {
		t.Errorf("expected %d virtual nodes, got %d", expectedVirtualNodes, actualVirtualNodes)
	}
}

// 测试边界条件：空哈希环
func TestEmptyRing(t *testing.T) {
	ch := NewConsistentHash()

	key := "Key1"
	node := ch.GetNode(key)

	if node != "" {
		t.Errorf("expected no node for key %s on empty ring, got %s", key, node)
	}
}

// 测试高并发访问一致性
func TestConcurrentAccess(t *testing.T) {
	ch := NewConsistentHash(WithReplica(3))

	ch.AddNode("NodeA")
	ch.AddNode("NodeB")
	ch.AddNode("NodeC")

	testKeys := []string{"Key1", "Key2", "Key3", "Key4", "Key5"}

	// 并发测试
	t.Run("concurrent-get", func(t *testing.T) {
		t.Parallel()
		for _, key := range testKeys {
			go func(k string) {
				node := ch.GetNode(k)
				if node == "" {
					t.Errorf("failed to get node for key %s", k)
				}
			}(key)
		}
	})
}
