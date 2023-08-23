package main

import (
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"strings"
	"time"
)

func main() {
	stopCh := make(chan struct{})
	kubeClient, err := buildKubeClient()
	if err != nil {
		println(err.Error())
		panic(err.Error())
	}

	listOptions := informers.WithTweakListOptions(func(options *v12.ListOptions) {
	})
	//测试Resync能力
	informerFactory := informers.NewSharedInformerFactoryWithOptions(kubeClient, 5*time.Second, listOptions)
	informer := informerFactory.Core().V1().Pods().Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod := obj.(*v1.Pod)
			if strings.Contains(pod.Name, "custom-scheduler") {
				println("add ", pod.Name, "  ", time.Now().GoString())
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			pod := newObj.(*v1.Pod)
			if strings.Contains(pod.Name, "custom-scheduler") {
				println("update ", pod.Name, "  ", time.Now().GoString())
			}
		},
	})
	informerFactory.Start(stopCh)
	go func() {
		time.Sleep(5 * time.Minute)
		println("stop server....")
		stopCh <- struct{}{}
	}()
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
