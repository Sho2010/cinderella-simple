package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Sho2010/cinderella-simple/audit"
	"github.com/Sho2010/cinderella-simple/config"
	"github.com/Sho2010/cinderella-simple/controller"
	"github.com/Sho2010/cinderella-simple/k8s"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	fmt.Println("hello cinderella")

	client := initClient()

	// gen(client)

	controller.StoreSecret()

	c, err := k8s.NewCleaner(client)
	if err != nil {
		panic(err)
	}

	go c.Start(context.TODO())
	h := audit.LogHandler{
		LogWriter: os.Stdout,
	}

	h.Start(audit.AuditCh)

	select {} // Block all
}

func gen(client kubernetes.Interface) {
	gen := k8s.KubeconfigGenerator{
		Client: client,
	}
	gen.Generate(os.Stdout, "test", "glass")

}

func initClient() kubernetes.Interface {
	var kubeClient kubernetes.Interface
	var server string

	if _, err := rest.InClusterConfig(); err != nil {
		kubeClient, server = k8s.GetClientOutOfCluster()
	} else {
		kubeClient, server = k8s.GetClient()
	}

	c := config.GetConfig()
	c.KubeServer = server

	return kubeClient
}
