package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/openshift/client-go/config"
	"github.com/openshift/client-go/user/clientset/versioned"
	v1 "github.com/openshift/api/user/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/homedir"
)

func main() {
	// Set the path to your OpenShift kubeconfig file.
	kubeconfigPath := filepath.Join(homedir.HomeDir(), ".kube", "config")

	// Load the OpenShift configuration from the specified kubeconfig file.
	ocConfig, err := config.LoadKubeConfig(kubeconfigPath)
	if err != nil {
		fmt.Printf("Error loading kubeconfig: %v\n", err)
		os.Exit(1)
	}

	// Create an OpenShift clientset.
	clientset, err := versioned.NewForConfig(ocConfig)
	if err != nil {
		fmt.Printf("Error creating OpenShift clientset: %v\n", err)
		os.Exit(1)
	}

	// Define the Group object.
	group := &v1.Group{
		ObjectMeta: metav1.ObjectMeta{
			Name: "example-group",
		},
		Users: []string{
			"user1",
			"user2",
		},
	}

	// Create the Group resource.
	createdGroup, err := clientset.UserV1().Groups().Create(group)
	if err != nil {
		fmt.Printf("Error creating Group: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Group created successfully with name: %s\n", createdGroup.Name)
}
