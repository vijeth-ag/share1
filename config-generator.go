package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type ConfigMap struct {
	APIVersion     string    `yaml:"apiVersion"`
	Kind           string    `yaml:"kind"`
	Clusters       []Cluster `yaml:"clusters"`
	Users          []User    `yaml:"users"`
	Contexts       []Context `yaml:"contexts"`
	CurrentContext string    `yaml:"current-context"`
}

type Cluster struct {
	Name    string `yaml:"name"`
	Cluster struct {
		Server string `yaml:"server"`
		IP     string `yaml:"ip"`
	} `yaml:"cluster"`
}

type User struct {
	Name string `yaml:"name"`
	User struct {
		ClientCertificateData string `yaml:"client-certificate-data"`
		ClientKeyData         string `yaml:"client-key-data"`
	} `yaml:"user"`
}

type Context struct {
	Name    string `yaml:"name"`
	Context struct {
		Cluster string `yaml:"cluster"`
		User    string `yaml:"user"`
	} `yaml:"context"`
}

func main() {
	// Fetch the cluster IP dynamically
	clusterIP, err := getClusterIP()
	if err != nil {
		fmt.Printf("Error fetching cluster IP: %v\n", err)
		os.Exit(1)
	}

	configMap := ConfigMap{
		APIVersion: "v1",
		Kind:       "Config",
		Clusters: []Cluster{
			{
				Name: "my-cluster",
				Cluster: struct {
					Server string `yaml:"server"`
					IP     string `yaml:"ip"`
				}{
					Server: fmt.Sprintf("https://%s", clusterIP),
					IP:     clusterIP,
				},
			},
		},
		Users: []User{
			{
				Name: "my-user",
				User: struct {
					ClientCertificateData string `yaml:"client-certificate-data"`
					ClientKeyData         string `yaml:"client-key-data"`
				}{
					ClientCertificateData: "...", // Base64-encoded certificate data
					ClientKeyData:         "...", // Base64-encoded key data
				},
			},
		},
		Contexts: []Context{
			{
				Name: "my-context",
				Context: struct {
					Cluster string `yaml:"cluster"`
					User    string `yaml:"user"`
				}{
					Cluster: "my-cluster",
					User:    "my-user",
				},
			},
		},
		CurrentContext: "my-context",
	}

	ip, err := getClusterIP()
	if err != nil {
		log.Println("getclusterip err", err)
	}
	log.Println("CLUSTER IP", ip)

	yamlData, err := yaml.Marshal(configMap)
	if err != nil {
		fmt.Printf("Error marshalling YAML: %v\n", err)
		os.Exit(1)
	}

	err = ioutil.WriteFile("output.yaml", yamlData, 0644)
	if err != nil {
		fmt.Printf("Error writing YAML file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("YAML file generated successfully.")
}

func createKubeClient(configPath string) (*kubernetes.Clientset, error) {
	// Load the Kubernetes configuration from the specified file path
	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		return nil, err
	}

	// Create a Kubernetes client using the configuration
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func getClusterIP() (string, error) {
	// config, err := rest.InClusterConfig()

	configPath := "/Users/vijeth.ag/.kube/config"

	// Create a Kubernetes client using the configuration file
	client, err := createKubeClient(configPath)
	if err != nil {
		fmt.Printf("Error creating Kubernetes client: %v\n", err)
		os.Exit(1)
	}

	service, err := client.CoreV1().Services("default").Get(context.TODO(), "kubernetes", metav1.GetOptions{})
	if err != nil {
		log.Println("err3", err)
		return "", err
	}

	return service.Spec.ClusterIP, nil
}
