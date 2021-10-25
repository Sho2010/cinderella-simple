package controller

import (
	"bytes"
	"io"

	"github.com/Sho2010/cinderella-simple/config"
	"github.com/Sho2010/cinderella-simple/encrypt"
	"github.com/Sho2010/cinderella-simple/k8s"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func StoreSecret() {
	r, w := io.Pipe()
	defer r.Close()

	go func() {
		defer w.Close()
		gen := k8s.KubeconfigGenerator{
			Client: initClient(),
		}
		gen.Generate(w, "test", "glass")
		w.Close()
	}()

	enc := encrypt.ZipEncrypter{
		Password: "password",
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	enc.Encrypt("/tmp/zip/abc.zip", buf.Bytes())

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
