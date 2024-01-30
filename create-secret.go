package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/openshift/api/config/v1"
	occlient "github.com/openshift/client-go/config/clientset/versioned"
	ocv1 "github.com/openshift/api/config/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	// Set the kubeconfig path
	kubeconfigPath := filepath.Join(os.Getenv("HOME"), ".kube", "config")

	// Initialize the OpenShift client
	openshiftClient, err := getOpenShiftClientSet(kubeconfigPath)
	if err != nil {
		fmt.Printf("Error creating OpenShift client: %v\n", err)
		os.Exit(1)
	}

	// Create a Secret object
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-secret",
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			"username": []byte("admin"),
			"password": []byte("password123"),
		},
	}

	// Create or update the Secret in the OpenShift cluster
	err = createOrUpdateSecret(openshiftClient
