package excellent_design

import (
	"golib/libs/concurrency"
	"testing"
	"time"
)

func TestKubelet_Run(t *testing.T) {
	kl := Kubelet{}

	updates := make(chan Pod, 5)
	wg := concurrency.Group{}
	wg.Start(func() {
		for i := 0; i < 10; i++ {
			updates <- Pod{Status: "running"}
		}
	})

	wg.Start(func() {
		kl.Run(updates)
		time.Sleep(2 * time.Second)
	})
	wg.Wait()

}
