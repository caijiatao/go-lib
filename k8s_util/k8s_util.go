package k8s_util

import (
	"bytes"
	"context"
	"io"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
	"sync"
)

var (
	kubeConfig         *restclient.Config
	kubeConfigInitOnce sync.Once
)

func getKubeConfig() *restclient.Config {
	kubeConfigInitOnce.Do(func() {
		kubeconfig := filepath.Join("./", ".kube", "config")
		var err error
		kubeConfig, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err)
		}
	})
	return kubeConfig
}

func GetClient() (*kubernetes.Clientset, error) {
	client, err := kubernetes.NewForConfig(getKubeConfig())
	if err != nil {
		return nil, err
	}
	return client, nil
}

func HasNamespace(ctx context.Context, name string) (bool, error) {
	client, err := GetClient()
	if err != nil {
		return false, err
	}
	_, err = client.CoreV1().Namespaces().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func HasConfigMap(ctx context.Context, namespace, name string) (bool, error) {
	client, err := GetClient()
	if err != nil {
		return false, err
	}
	_, err = client.CoreV1().ConfigMaps(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func GetYamlStructByFileName(fileName string) (unstructuredObjs []*unstructured.Unstructured, err error) {
	unstructuredObjs = make([]*unstructured.Unstructured, 0)
	b, err := os.ReadFile(fileName)
	decoder := yamlutil.NewYAMLOrJSONDecoder(bytes.NewReader(b), 100)
	for {
		var rawObj runtime.RawExtension
		if err = decoder.Decode(&rawObj); err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}

		obj, _, err := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
		if err != nil {
			return nil, err
		}
		unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
		if err != nil {
			return nil, err
		}
		unstructuredObj := &unstructured.Unstructured{Object: unstructuredMap}
		unstructuredObjs = append(unstructuredObjs, unstructuredObj)
	}
	return unstructuredObjs, nil
}

func CleanResource(ctx context.Context, unstructuredObj *unstructured.Unstructured) (err error) {
	dri, err := GetDynamicResourceInterface(unstructuredObj)
	if err != nil {
		return err
	}
	// 没有资源需要清理，直接返回
	_, err = dri.Get(ctx, unstructuredObj.GetName(), metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	err = dri.Delete(ctx, unstructuredObj.GetName(), metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

func CreateResource(ctx context.Context, unstructuredObj *unstructured.Unstructured) (uncastObj *unstructured.Unstructured, err error) {
	dri, err := GetDynamicResourceInterface(unstructuredObj)
	if err != nil {
		return nil, err
	}
	if uncastObj, err = dri.Create(ctx, unstructuredObj, metav1.CreateOptions{}); err != nil {
		return nil, err
	}
	return uncastObj, nil
}

func GetDynamicResourceInterface(unstructuredObj *unstructured.Unstructured) (dri dynamic.ResourceInterface, err error) {
	clientset, err := GetClient()
	if err != nil {
		return nil, err
	}
	gr, err := restmapper.GetAPIGroupResources(clientset.Discovery())
	if err != nil {
		return nil, err
	}

	mapper := restmapper.NewDiscoveryRESTMapper(gr)
	mapping, err := mapper.RESTMapping(unstructuredObj.GroupVersionKind().GroupKind(), unstructuredObj.GroupVersionKind().Version)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	dynamicClient, err := dynamic.NewForConfig(getKubeConfig())
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		if unstructuredObj.GetNamespace() == "" {
			unstructuredObj.SetNamespace("default")
		}
		dri = dynamicClient.Resource(mapping.Resource).Namespace(unstructuredObj.GetNamespace())
	} else {
		dri = dynamicClient.Resource(mapping.Resource)
	}
	return dri, nil
}
