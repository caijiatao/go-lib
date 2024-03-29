package framework

import (
	"golib/libs/concurrency"
	"sync"
)

type Framework struct {
	sync.Mutex
	scorePlugins []ScorePlugin
}

func (f *Framework) RegisterScorePlugin(plugin ScorePlugin) {
	f.Lock()
	defer f.Unlock()
	f.scorePlugins = append(f.scorePlugins, plugin)
}

func (f *Framework) runScorePlugins(node string, pod string) int {
	var score int
	for _, plugin := range f.scorePlugins {
		score += plugin.Score(node, pod)
	}
	return score
}

type Pod struct {
	Name     string
	NodeName string // 绑定的节点
}

func (f *Framework) RunScorePlugins(nodes []string, pod *Pod) map[string]int {
	scores := make(map[string]int)
	p := concurrency.NewParallelizer(16)
	p.Until(len(nodes), func(i int) {
		scores[nodes[i]] = f.runScorePlugins(nodes[i], pod.Name)
	})
	// 省略绑定节点的逻辑
	return scores
}
