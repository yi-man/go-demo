package cluster

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	utilruntime "k8s.io/apimachinery/pkg/util/runtime"

	"github.com/zeromicro/go-zero/core/logx"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	deletePolicy       = metav1.DeletePropagationBackground
	gracePeriodSeconds = int64(30)
	secretTLSKey       = "secret-tls"
	ingressClass       = "nginx"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(appsv1.AddToScheme(scheme))
}

type ClusterClient struct {
	client     client.Client
	KubeClient *kubernetes.Clientset

	NameSpace string

	HarborHost      string
	HarborUserName  string
	HarborPwd       string
	HarborSecretKey string
}

func NewClusterClient(kubeConfigBytes []byte, namespace string, harborHost string, harborUserName string, harborPwd string) (*ClusterClient, error) {
	ClusterConfig, err := clientcmd.RESTConfigFromKubeConfig(kubeConfigBytes)
	if err != nil {
		logx.Error(err)
		return nil, err
	}

	k8sClientSet, err := kubernetes.NewForConfig(ClusterConfig)
	if err != nil {
		logx.Error(err)
		return nil, err
	}

	clientSet, err := client.New(ClusterConfig, client.Options{
		Scheme: scheme,
	})

	if err != nil {
		logx.Error(err)
		return nil, err
	}

	return &ClusterClient{
		client:     clientSet,
		KubeClient: k8sClientSet,

		NameSpace: namespace,

		HarborHost:      harborHost,
		HarborUserName:  harborUserName,
		HarborPwd:       harborPwd,
		HarborSecretKey: strings.ToLower(harborUserName) + "-private-secret-key",
	}, nil
}

func int32Ptr(i int32) *int32 { return &i }

func (cc *ClusterClient) CreateSecret(ctx context.Context) error {
	// 检查是否已经存在相同名称的 Secret
	_, err := cc.KubeClient.CoreV1().Secrets(cc.NameSpace).Get(
		ctx,
		cc.HarborSecretKey,
		metav1.GetOptions{},
	)
	if err == nil {
		// 如果已存在相同名称的 Secret，可以返回错误或者执行其他逻辑
		return err
	}

	// 创建一个 Docker Registry 认证的 Secret
	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cc.HarborSecretKey,
			Namespace: cc.NameSpace,
		},
		Type: v1.SecretTypeDockerConfigJson,
		Data: map[string][]byte{
			".dockerconfigjson": []byte(fmt.Sprintf(`{
			"auths": {
				"%s": {
					"username": "%s",
					"password": "%s",
					"auth": "%s"
				}
			}
		}`, cc.HarborHost, cc.HarborUserName, cc.HarborPwd, base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", cc.HarborUserName, cc.HarborPwd))))),
		},
	}

	_, err = cc.KubeClient.CoreV1().Secrets(cc.NameSpace).Create(
		ctx,
		secret,
		metav1.CreateOptions{},
	)
	if err != nil {
		return err
	}

	return nil
}
