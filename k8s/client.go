package k8s

import (
	"fmt"
	"log"
	"os"

	"github.com/Sho2010/cinderella-simple/config"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

//TODO: どうしようかな
var (
	ServiceAccountNamespace = "default"
)

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
		fmt.Printf("home: %s", home)
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
