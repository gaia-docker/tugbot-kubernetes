package main

import (
	"os"
	"os/user"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"k8s.io/kubernetes/pkg/client/restclient"
	client "k8s.io/kubernetes/pkg/client/unversioned"

	"github.com/gaia-docker/tugbot-kubernetes/action"
)

func main() {
	kube := getKubernetesClient()
	action.UpdateJobs(kube.Jobs(getNamespace()), "deployment")
}

func getKubernetesClient() *client.Client {
	kubeCertPath := os.Getenv("KUBERNETES_CERT_PATH")
	if kubeCertPath == "" {
		log.Println("KUBERNETES_CERT_PATH is not defined, looking for certificate files in .minikube folder under home directory.")
		usr, err := user.Current()
		if err != nil {
			log.Fatalln("Failed to get home directory.", err)
		}
		kubeCertPath = filepath.Join(usr.HomeDir, ".minikube")
	}
	host := os.Getenv("KUBERNETES_HOST")
	if host == "" {
		host = "https://192.168.99.100:8443"
		log.Println("KUBERNETES_HOST is not defined, using", host)
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
	return "default"
}
