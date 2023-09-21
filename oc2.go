package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/client-go/util/retry"
	"k8s.io/client-go/util/wait"
	"k8s.io/client-go/util/yaml"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
)

func main() {
	// Set the path to your OpenShift kubeconfig file.
	kubeconfigPath := filepath.Join(homedir.HomeDir(), ".kube", "config")

	// Load the OpenShift configuration from the specified kubeconfig file.
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		fmt.Printf("Error loading kubeconfig: %v\n", err)
		os.Exit(1)
	}

	// Create an OpenShift clientset.
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error creating OpenShift clientset: %v\n", err)
		os.Exit(1)
	}

	// Define the Group YAML manifest.
	groupYAML := `
apiVersion: user.openshift.io/v1
kind: Group
metadata:
  name: example-group
`

	// Unmarshal the YAML into a Group object.
	var group metav1.Group
	if err := yaml.Unmarshal([]byte(groupYAML), &group); err != nil {
		fmt.Printf("Error unmarshalling Group YAML: %v\n", err)
		os.Exit(1)
	}

	// Create the Group resource.
	createdGroup, err := clientset.UserV1().Groups().Create(context.TODO(), &group, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("Error creating Group: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Group created successfully with name: %s\n", createdGroup.Name)
}
