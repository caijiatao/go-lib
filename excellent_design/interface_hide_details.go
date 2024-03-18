package excellent_design

import "fmt"

type Pod struct {
	Status string
}

type Kubelet struct{}

func (kl *Kubelet) HandlePodAdditions(pods []*Pod) {
	for _, pod := range pods {
		fmt.Printf("create pods : %s\n", pod.Status)
	}
}

func (kl *Kubelet) Run(updates <-chan Pod) {
	fmt.Println(" run kubelet")
	go kl.syncLoop(updates, kl)
}

func (kl *Kubelet) syncLoop(updates <-chan Pod, handler SyncHandler) {
	for {
		select {
		case pod := <-updates:
			handler.HandlePodAdditions([]*Pod{&pod})
		}
	}
}

type SyncHandler interface {
	HandlePodAdditions(pods []*Pod)
}
