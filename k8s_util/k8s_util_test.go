package k8s_util

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"testing"
	"time"
)

var (
	testNamespace = "exploit-ghrec"
	testClient    *kubernetes.Clientset
)

func TestMain(m *testing.M) {
	var err error
	testClient, err = GetClient()
	if err != nil {
		panic(err)
	}
	m.Run()
}

func TestGetClient(t *testing.T) {
	client, err := GetClient()
	assert.Nil(t, err)
	list, err := client.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	assert.Greater(t, len(list.Items), 0)
	assert.Nil(t, err)
}

func TestCreateNamespace(t *testing.T) {
	hasNamespace, err := HasNamespace(context.Background(), testNamespace)
	assert.Nil(t, err)
	if hasNamespace {
		err = testClient.CoreV1().Namespaces().Delete(context.Background(), testNamespace, metav1.DeleteOptions{})
		assert.Nil(t, err)
	}
	nsName := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: testNamespace,
		},
	}
	resp, err := testClient.CoreV1().Namespaces().Create(context.Background(), nsName, metav1.CreateOptions{})
	assert.Nil(t, err)
	fmt.Println(resp)
}

func TestCreateConfigMap(t *testing.T) {
	configMapName := "airec-config-map"
	hasConfigMap, err := HasConfigMap(context.Background(), testNamespace, configMapName)
	assert.Nil(t, err)
	if hasConfigMap {
		err = testClient.CoreV1().ConfigMaps(testNamespace).Delete(context.Background(), configMapName, metav1.DeleteOptions{})
		assert.Nil(t, err)
	}
	createData := map[string]string{
		"REC_VERSION": fmt.Sprintf("%d", time.Now().Unix()),
	}
	configMap := corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: testNamespace,
		},
		Data: createData,
	}
	resp, err := testClient.CoreV1().ConfigMaps(testNamespace).Create(context.Background(), &configMap, metav1.CreateOptions{})
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, createData, resp.Data)

	// 更新Config map
	updateData := map[string]string{
		"REC_VERSION": fmt.Sprintf("%d", time.Now().Unix()),
	}
	configMap.Data = updateData
	resp, err = testClient.CoreV1().ConfigMaps(testNamespace).Update(context.Background(), &configMap, metav1.UpdateOptions{})
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, updateData, resp.Data)
}

func TestHasNamespace(t *testing.T) {
	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "",
			args: args{
				ctx:  context.TODO(),
				name: "not-exists",
			},
			want:    false,
			wantErr: assert.NoError,
		},
		{
			name: "",
			args: args{
				ctx:  context.TODO(),
				name: "test-namespace",
			},
			want:    true,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HasNamespace(tt.args.ctx, tt.args.name)
			if !tt.wantErr(t, err, fmt.Sprintf("HasNamespace(%v, %v)", tt.args.ctx, tt.args.name)) {
				return
			}
			assert.Equalf(t, tt.want, got, "HasNamespace(%v, %v)", tt.args.ctx, tt.args.name)
		})
	}
}

func TestCreateResource(t *testing.T) {
	// 1.读取yaml文件
	objs, err := GetYamlStructByFileName("./airec-server.yaml")
	assert.Nil(t, err)
	assert.Len(t, objs, 2)
	for _, obj := range objs {
		assert.Equal(t, "", obj.GetNamespace())
		assert.Equal(t, "airec-server-test", obj.GetName())
		assert.NotEqual(t, "", obj.GetKind())
		obj.SetNamespace(testNamespace)
	}

	ctx := context.Background()
	for _, obj := range objs {
		err = CleanResource(ctx, obj)
		assert.Nil(t, err)
	}

	// 设置namespace 并创建资源
	for _, obj := range objs {
		assert.Equal(t, testNamespace, obj.GetNamespace())
		uncastObj, err := CreateResource(ctx, obj)
		assert.Nil(t, err)
		assert.NotNil(t, uncastObj)
	}
	// 清理资源
	for _, obj := range objs {
		err = CleanResource(ctx, obj)
		assert.Nil(t, err)
	}
}
