package main

import (
	"os"
	"os/user"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/gaia-docker/tugbot-kubernetes/action"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/cache"
	"k8s.io/kubernetes/pkg/client/restclient"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/controller/framework"
	"k8s.io/kubernetes/pkg/fields"
	"os/signal"
	"syscall"
)

var (
	namespace string
	kube      *client.Client
	stop      chan struct{}
)

func main() {
	initialize()
	watchEvents()
	waitForInterrupt()
}

func initialize() {
	setLogLevel()
	namespace = getNamespace()
	kube = getKubernetesClient()
	stop = make(chan struct{})
}

func setLogLevel() {
	if os.Getenv("TUGBOT_KUBERNETES_LOG_LEVEL") == "debug" {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}

func getKubernetesClient() *client.Client {
	kubeCertPath := os.Getenv("KUBERNETES_CERT_PATH")
	if kubeCertPath == "" {
		log.Info("KUBERNETES_CERT_PATH is not defined, looking for certificate files in .minikube folder under home directory.")
		usr, err := user.Current()
		if err != nil {
			log.Fatalln("Failed to get home directory.", err)
		}
		kubeCertPath = filepath.Join(usr.HomeDir, ".minikube")
		log.Infof("Looking for ca.crt, apiserver.crt, apiserver.key in %s", kubeCertPath)
	}
	host := os.Getenv("KUBERNETES_HOST")
	if host == "" {
		host = "https://192.168.99.100:8443"
		log.Infof("KUBERNETES_HOST is not defined, using %s", host)
	}
	ret, err := client.New(&restclient.Config{
		Host: host,
		TLSClientConfig: restclient.TLSClientConfig{
			CAFile:   filepath.Join(kubeCertPath, "ca.crt"),
			CertFile: filepath.Join(kubeCertPath, "apiserver.crt"),
			KeyFile:  filepath.Join(kubeCertPath, "apiserver.key")},
	})
	if err != nil {
		log.Fatalln("Failed to create Kubernetes client. ", err)
	}

	return ret
}

func getNamespace() string {
	ret := os.Getenv("TUGBOT_KUBERNETES_NAMESPACE")
	if ret == "" {
		ret = api.NamespaceDefault
		log.Infof("TUGBOT_KUBERNETES_NAMESPACE is not defined, using %s", ret)
	}

	return ret
}

func watchEvents() {
	go func() {
		watchList := cache.NewListWatchFromClient(kube, "events", api.NamespaceAll, fields.Everything())
		_, eventController := framework.NewInformer(watchList,
			&api.Event{},
			0,
			framework.ResourceEventHandlerFuncs{
				AddFunc: func(obj interface{}) {
					updateJobs(obj)
				},
				DeleteFunc: func(obj interface{}) {
					updateJobs(obj)
				},
			},
		)
		log.Info("Start watching for Kubernetes Events...")
		eventController.Run(stop)
	}()
}

func updateJobs(event interface{}) {
	action.UpdateJobs(kube.BatchClient.Jobs(namespace), event.(*api.Event))
}

func waitForInterrupt() {
	// Graceful shut-down on SIGINT/SIGTERM/SIGQUIT
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	<-c
	log.Info("Stopping Tugbot...")
	close(stop)
	os.Exit(1)
}
