package config

import (
	"github.com/rs/zerolog/log"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

func ConfigWatcher(clientSet *kubernetes.Clientset) {
	namespace := "namespace"
	configMapName := "configMap"

	listWatcher := cache.NewListWatchFromClient(
		clientSet.CoreV1().RESTClient(),
		"configmaps",
		namespace,
		fields.Everything(),
	)

	informer := cache.NewSharedInformer(
		listWatcher,
		&corev1.ConfigMap{},
		0, // No resync
	)

	_, err := informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		UpdateFunc: func(oldObj, newObj interface{}) {
			newConfigMap, ok := newObj.(*corev1.ConfigMap)
			if ok && newConfigMap.Name == configMapName {
				// Handle the update logic here
				log.Info().Msg("ConfigMap updated")
			}
		},
		// Handle AddFunc and DeleteFunc if needed
	})
	if err != nil {
		return
	}
	stopCh := make(chan struct{})
	defer close(stopCh)
	go informer.Run(stopCh)
	select {}
}
