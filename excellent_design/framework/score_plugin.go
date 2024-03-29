package framework

type ScorePlugin interface {
	Score(NodeName string, pod string) (score int)
}

type ImageLocality struct{}

func (pl *ImageLocality) Score(NodeName string, pod string) (score int) {
	return 0
}

type NodeAffinity struct{}

func (pl *NodeAffinity) Score(NodeName string, pod string) (score int) {
	return 0
}
