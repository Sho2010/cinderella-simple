package k8s

import (
	"log"
	"os"

	"github.com/Sho2010/cinderella-simple/config"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	ManagedLabel      = "app.kubernetes.io/managed-by"
	ManagedLabelValue = "cinderella"
)

// TODO: configもしくは自分自身(POD)から取得する
var _cinderellaNamespace = ""

var _managedResourceLabels = map[string]string{
	ManagedLabel: ManagedLabelValue,
}

//TODO: client取得するときにHost一緒に返すのあまりにも使いづらいのでなんとかする
func GetDefaultClient() (kubernetes.Interface, string) {
	var kubeClient kubernetes.Interface
	var server string

	if _, err := rest.InClusterConfig(); err != nil {
		kubeClient, server = GetClientOutOfCluster()
	} else {
		kubeClient, server = GetClient()
	}

	c := config.GetConfig()
	c.KubeServer = server

	return kubeClient, server
}

func GetClient() (kubernetes.Interface, string) {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Can not get kubernetes config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Can not create kubernetes client: %v", err)
	}

	return clientset, config.Host
}

func buildOutOfClusterConfig() (*rest.Config, error) {
	kubeconfigPath := os.Getenv("KUBECONFIG")
	if kubeconfigPath == "" {
		home, _ := os.UserHomeDir()
		kubeconfigPath = home + "/.kube/config"
	}
	return clientcmd.BuildConfigFromFlags("", kubeconfigPath)
}

// GetClientOutOfCluster returns a k8s clientset to the request from outside of cluster
func GetClientOutOfCluster() (kubernetes.Interface, string) {
	config, err := buildOutOfClusterConfig()
	if err != nil {
		log.Fatalf("Can not get kubernetes config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Can not get kubernetes config: %v", err)
	}

	return clientset, config.Host
}

// GetCinderellaNamespace returns the namespace of cinderella
// it returns the get from CINDERELLA_POD_NAMESPACE environment variables
// if it is empty, it returns the get from configfile
func GetCinderellaNamespace() string {
	// See: https://kubernetes.io/ja/docs/tasks/inject-data-application/environment-variable-expose-pod-information/
	if _cinderellaNamespace == "" {
		log.Println("Get cinderella Namespace from CINDERELLA_POD_NAMESPACE environment variable")
		_cinderellaNamespace = os.Getenv("CINDERELLA_POD_NAMESPACE")
	}

	if _cinderellaNamespace == "" {
		log.Println("Get cinderella Namespace from configfile")
		_cinderellaNamespace = config.GetConfig().Namespace
	}

	return _cinderellaNamespace
}
