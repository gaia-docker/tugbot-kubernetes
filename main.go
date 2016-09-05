package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	"k8s.io/kubernetes/pkg/apis/batch"
	"k8s.io/kubernetes/pkg/client/restclient"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"os/user"
)

func main() {
	createdJob, err := getKubernetesClient().BatchClient.Jobs("default").Create(getJob())
	if err != nil {
		log.Fatalln("Failed to unmarshal Job:", err)
	}
	log.Println("Job Created :)", createdJob)
}

func getJob() *batch.Job {
	file, err := ioutil.ReadFile("job.yaml")
	if err != nil {
		log.Fatal("Failed to read file.", err)
	}
	var ret batch.Job
	if err := yaml.Unmarshal(file, &ret); err != nil {
		log.Fatalln("Failed to unmarshal Job.", err)
	}

	return &ret
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
		log.Println(kubeCertPath)
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
