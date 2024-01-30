package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
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
	configMap := ConfigMap{
		APIVersion: "v1",
		Kind:       "Config",
		Clusters: []Cluster{
			{
				Name: "my-cluster",
				Cluster: struct {
					Server string `yaml:"server"`
				}{
					Server: "https://cluster-api-server",
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

	configMap.Clusters[0].Cluster.Server = "https://myauth.com/v1/authentication"

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
