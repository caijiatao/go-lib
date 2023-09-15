package k8s_util

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog/log"
	"golib/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	k8sCache "k8s.io/client-go/tools/cache"
)

type configMapOperateType int

const (
	updateConfigMap configMapOperateType = iota + 1
	deleteConfigMap
)

type ConfigMapOnChangeFunc func(namespace, configMapName string, oldDataMap, newDataMap *corev1.ConfigMap, operate configMapOperateType) error

func ConfigWatcher(ctx context.Context, namespace, configMapName string, onChangeFunc ConfigMapOnChangeFunc) (cancel func()) {
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

	err = initConfigMapCache(ctx, namespace, configMapName, onChangeFunc)
	if err != nil {
		panic(err)
	}

	_, err = informer.AddEventHandler(k8sCache.ResourceEventHandlerFuncs{
		UpdateFunc: func(oldObj, newObj interface{}) {
			log.Info().Msg("ConfigMap updated")
			oldConfigMap, ok := oldObj.(*corev1.ConfigMap)
			if !ok {
				logger.Errorf("Old object is not a ConfigMap: %v", oldObj)
				return
			}

			newConfigMap, ok := newObj.(*corev1.ConfigMap)
			if !ok {
				logger.Errorf("New object is not a ConfigMap: %v", newObj)
				return
			}
			if newConfigMap.Name != configMapName {
				logger.Errorf("New object config name error: %s , configMapName :%s", newConfigMap.Name, configMapName)
				return
			}
			err = onChangeFunc(namespace, configMapName, oldConfigMap, newConfigMap, updateConfigMap)
			if err != nil {
				logger.Errorf("update config map err:%s, configMapName:%s, namespace:%s", err.Error(), configMapName, namespace)
			}
		},
		// Handle AddFunc and DeleteFunc if needed
		DeleteFunc: func(obj interface{}) {
			log.Info().Msg("ConfigMap deleted")
			err = onChangeFunc(namespace, configMapName, nil, nil, deleteConfigMap)
			if err != nil {
				logger.Errorf("delete config map err:%s, configMapName:%s, namespace:%s", err.Error(), configMapName, namespace)
			}
		},
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

func initConfigMapCache(ctx context.Context, namespace, configMapName string, onChangeFunc ConfigMapOnChangeFunc) (err error) {
	clientSet, err := GetClient()
	if err != nil {
		panic(err)
	}
	configMap, err := clientSet.CoreV1().ConfigMaps(namespace).Get(ctx, configMapName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	return onChangeFunc(namespace, configMapName, nil, configMap, updateConfigMap)
}

func ConfigMapOnChangeCache[T any](c *cache.Cache) ConfigMapOnChangeFunc {
	return func(namespace, configMapName string, oldDataMap, newDataMap *corev1.ConfigMap, operate configMapOperateType) (err error) {
		switch operate {
		case updateConfigMap:
			err = updateConfigMapOnChange[T](c, namespace, configMapName, oldDataMap, newDataMap)
		case deleteConfigMap:
			err = deleteConfigMapOnChange(c, namespace, configMapName)
		}
		return err
	}
}

func updateConfigMapOnChange[T any](c *cache.Cache, namespace, configMapName string, oldDataMap, newDataMap *corev1.ConfigMap) (err error) {
	dataByte, err := json.Marshal(newDataMap.Data)
	if err != nil {
		return err
	}
	data := new(T)
	err = json.Unmarshal(dataByte, data)
	if err != nil {
		return err
	}
	c.SetDefault(GetConfigMapCacheKey(namespace, configMapName), *data)
	return nil
}

func deleteConfigMapOnChange(c *cache.Cache, namespace, configMapName string) (err error) {
	c.Delete(GetConfigMapCacheKey(namespace, configMapName))
	return nil
}
