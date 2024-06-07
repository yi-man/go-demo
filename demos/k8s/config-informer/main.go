package main

import (
	"fmt"
	"path/filepath"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	homeDir "k8s.io/client-go/util/homedir"
)

func main() {
	kubeconfig := filepath.Join(
		homeDir.HomeDir(), ".kube", "core-dev-poodle.cfg",
	)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	factory := informers.NewSharedInformerFactory(clientset, time.Minute*10)
	configMapInformer := factory.Core().V1().ConfigMaps().Informer()

	configMapInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			configMap := obj.(*v1.ConfigMap)
			fmt.Printf("ConfigMap added: %s/%s\n", configMap.Namespace, configMap.Name)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldConfigMap := oldObj.(*v1.ConfigMap)
			// newConfigMap := newObj.(*v1.ConfigMap)
			fmt.Printf("ConfigMap updated: %s/%s\n", oldConfigMap.Namespace, oldConfigMap.Name)
		},
		DeleteFunc: func(obj interface{}) {
			configMap := obj.(*v1.ConfigMap)
			fmt.Printf("ConfigMap deleted: %s/%s\n", configMap.Namespace, configMap.Name)
		},
	})

	stopCh := make(chan struct{})
	defer close(stopCh)

	factory.Start(stopCh)
	factory.WaitForCacheSync(stopCh)

	<-stopCh
}
