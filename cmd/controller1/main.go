package main

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"time"
)

func main() {
	stopCh := make(chan struct{})
	kubeClient, err := buildKubeClient()
	if err != nil {
		println(err.Error())
		panic(err.Error())
	}
	informerFactory := informers.NewSharedInformerFactoryWithOptions(kubeClient, 5*time.Second)
	//podLister := informerFactory.Core().V1().Pods().Lister()
	informer := informerFactory.Core().V1().Pods().Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod := obj.(*v1.Pod)
			println("add ", pod.Name)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			pod := newObj.(*v1.Pod)
			println("update  ", pod.Name)
		},
	})
	informerFactory.Start(stopCh)
	<-stopCh
}

func buildKubeClient() (*kubernetes.Clientset, error) {
	kubeConfigDir := homedir.HomeDir() + "/.kube/config"
	println("kubeconfig=", kubeConfigDir)

	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigDir)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}
