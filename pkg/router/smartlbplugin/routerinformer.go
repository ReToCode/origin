package smartlbplugin

import (
	"k8s.io/kubernetes/pkg/client/clientset_generated/clientset"

	kinformers "k8s.io/kubernetes/pkg/client/informers/informers_generated/externalversions"
	"k8s.io/kubernetes/pkg/api/v1"

	"k8s.io/client-go/tools/cache"
	"time"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
)

// CreateRouterInformer creates an informer for the pod object and calls the plugin on changes
func CreateAndRunRouterInformer(kClient clientset.Interface, plugin *SmartLBPlugin) {
	informerFactory := kinformers.NewSharedInformerFactory(kClient, 5 * time.Minute)
	podInformer := informerFactory.Core().V1().Pods()

	podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		UpdateFunc: func(oldPod, newPod interface{}) {
			plugin.HandlePod(watch.Modified, newPod.(*v1.Pod))
		},
		DeleteFunc: func(p interface{}) {
			plugin.HandlePod(watch.Deleted, p.(*v1.Pod))
		},
		AddFunc: func(p interface{}) {
			plugin.HandlePod(watch.Added, p.(*v1.Pod))
		},
	})

	informerFactory.Start(wait.NeverStop)
}
