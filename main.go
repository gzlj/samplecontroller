package main

import (
	"flag"
	"fmt"
	samplecontrollerv1 "github.com/gzlj/samplecontroller/pkg/apis/samplecontroller/v1"
	clientset "github.com/gzlj/samplecontroller/pkg/generated/clientset/versioned"
	informers "github.com/gzlj/samplecontroller/pkg/generated/informers/externalversions"
	//"k8s.io/sample-controller/pkg/generated/informers/externalversions/samplecontroller/v1alpha1"
	samplecontrollerinformerv1 "github.com/gzlj/samplecontroller/pkg/generated/informers/externalversions/samplecontroller/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"

	"time"

	/*
		clientset "k8s.io/sample-controller/pkg/generated/clientset/versioned"
		informers "k8s.io/sample-controller/pkg/generated/informers/externalversions"

		 */
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/apimachinery/pkg/labels"
)
func main() {


	var (
		kubeConfigFile *string
		err            error
		config         *rest.Config
		databaseClient *clientset.Clientset
		//informers.
		//informers.DatabaseInformer
		databaseInformer samplecontrollerinformerv1.DatabaseInformer
	)
	//samplecontrollerv1.Database{}

	//kubeConfigFile = flag.String("kubeConfigFile", "/root/.kube/config", "kubernetes config file path")
	kubeConfigFile = flag.String("kubeConfigFile", "", "kubernetes config file path")
	flag.Parse()
	config, err = clientcmd.BuildConfigFromFlags("127.0.0.1:8080", *kubeConfigFile)
	log.Println("config:", config)
	databaseClient, err = clientset.NewForConfig(config)
	myDatabase, _  := databaseClient.SamplecontrollerV1().Databases("default").Get("my-database", metav1.GetOptions{})
	log.Println("data base :", myDatabase)
	fmt.Printf("===>Database Name:%v(%v,%v,%v)\n", myDatabase.Name, myDatabase.Spec.User, myDatabase.Spec.Password, myDatabase.Spec.Encoding)
	log.Println(err)
	//informers.n
	informerFactory := informers.NewSharedInformerFactory(databaseClient, time.Second*30)
	databaseInformer = informerFactory.Samplecontroller().V1().Databases()
	databaseInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {log.Println("add :", obj.(*samplecontrollerv1.Database).Name)},

	})


	stopCh := make(chan struct{})
	informerFactory.Start(stopCh)

	time.Sleep(time.Duration(2)*time.Second)
	databaseLister := databaseInformer.Lister()
	allDatabases, _ := databaseLister.List(labels.Everything())
	log.Println("allDatabases", allDatabases)
	for _, p := range allDatabases {
		fmt.Printf("list database: %v\n", p.Name)
	}
	<- stopCh
}
