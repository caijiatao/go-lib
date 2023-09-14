package k8s_util

import (
	"encoding/json"
	"fmt"
	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog/log"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	k8sCache "k8s.io/client-go/tools/cache"
)

type ConfigMapOnChangeFunc func(namespace, configMapName string, oldDataMap, newDataMap *corev1.ConfigMap)

func ConfigWatcher(namespace, configMapName string, onChangeFunc ConfigMapOnChangeFunc) (cancel func()) {
	clientSet, err := GetClient()
	if err != nil {
		panic(err)
	}

	listWatcher := k8sCache.NewListWatchFromClient(
		clientSet.CoreV1().RESTClient(),
		"configmaps",
		namespace,
		fields.Everything(),
	)

	informer := k8sCache.NewSharedInformer(
		listWatcher,
		&corev1.ConfigMap{},
		0, // No resync
	)

	_, err = informer.AddEventHandler(k8sCache.ResourceEventHandlerFuncs{
		UpdateFunc: func(oldObj, newObj interface{}) {
			log.Info().Msg("ConfigMap updated")
			newConfigMap, ok := newObj.(*corev1.ConfigMap)

			oldConfigMap, ok := oldObj.(*corev1.ConfigMap)
			if ok && newConfigMap.Name == configMapName {
				// Handle update logic here
				onChangeFunc(namespace, configMapName, oldConfigMap, newConfigMap)
			}

		},
		// Handle AddFunc and DeleteFunc if needed
	})
	if err != nil {
		panic(err)
	}

	stopCh := make(chan struct{})
	cancel = func() {
		close(stopCh)
	}
	go informer.Run(stopCh)
	return cancel
}

func GetConfigMapCacheKey(namespace, configMapName string) string {
	return fmt.Sprintf("%s@@%s", namespace, configMapName)
}

func ConfigMapOnChangeCache[T any](c *cache.Cache) ConfigMapOnChangeFunc {
	return func(namespace, configMapName string, oldDataMap, newDataMap *corev1.ConfigMap) {
		dataByte, err := json.Marshal(newDataMap.Data)
		if err != nil {
			return
		}
		data := new(T)
		err = json.Unmarshal(dataByte, data)
		if err != nil {
			return
		}
		c.SetDefault(GetConfigMapCacheKey(namespace, configMapName), *data)
	}
}
