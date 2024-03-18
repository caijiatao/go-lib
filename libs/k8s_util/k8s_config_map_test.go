package k8s_util

import (
	"context"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
	"time"
)

func TestConfigWatcher(t *testing.T) {
	namespace := "test-namespace"
	configMapName := "test-config-map"
	client, err := GetClient()
	assert.Nil(t, err)

	exists, err := HasConfigMap(context.Background(), namespace, configMapName)
	assert.Nil(t, err)
	if exists {
		err = client.CoreV1().ConfigMaps(namespace).Delete(context.Background(), configMapName, metav1.DeleteOptions{})
		assert.Nil(t, err)
	}

	type testStruct struct {
		RecVersion string `json:"REC_VERSION"`
	}

	// 1.创建configMap
	testConfigMapData := testStruct{
		RecVersion: "1.0.0",
	}
	createData := map[string]string{
		"REC_VERSION": testConfigMapData.RecVersion,
	}
	configMap := corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: namespace,
		},
		Data: createData,
	}
	ctx := context.Background()
	_, err = client.CoreV1().ConfigMaps(namespace).Create(ctx, &configMap, metav1.CreateOptions{})
	assert.Nil(t, err)

	c := cache.New(cache.NoExpiration, cache.NoExpiration)
	cancel := ConfigWatcher(ctx, namespace, configMapName, ConfigMapOnChangeCache[testStruct](c))
	// 2.获取配置并判断
	time.Sleep(time.Second)
	value, ok := c.Get(GetConfigMapCacheKey(namespace, configMapName))
	assert.True(t, ok)
	configMapValue := value.(testStruct)
	assert.Equal(t, testConfigMapData.RecVersion, configMapValue.RecVersion)

	// 3.修改configMap
	testConfigMapData.RecVersion = "1.0.1"
	createData["REC_VERSION"] = testConfigMapData.RecVersion
	configMap.Data = createData
	_, err = client.CoreV1().ConfigMaps(namespace).Update(context.Background(), &configMap, metav1.UpdateOptions{})
	assert.Nil(t, err)
	time.Sleep(time.Second)
	value, ok = c.Get(GetConfigMapCacheKey(namespace, configMapName))
	assert.True(t, ok)
	configMapValue = value.(testStruct)
	assert.Equal(t, testConfigMapData.RecVersion, configMapValue.RecVersion)

	// 4.删除configMap
	err = client.CoreV1().ConfigMaps(namespace).Delete(context.Background(), configMapName, metav1.DeleteOptions{})
	assert.Nil(t, err)
	cancel()
}
