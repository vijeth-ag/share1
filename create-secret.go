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

	// Create a Secret object with base64-encoded data
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-secret",
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			"username": []byte("YWRtaW4="), // base64-encoded "admin"
			"password": []byte("cGFzc3dvcmQxMjM="), // base64-encoded "password123"
		},
	}

	// Create or update the Secret in the OpenShift cluster
	err = createOrUpdateSecret(openshiftClient, secret)
	if err != nil {
		fmt.Printf("Error creating or updating Secret: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Secret created or updated successfully.")
}

func getOpenShiftClientSet(kubeconfigPath string) (*occlient.Clientset, error) {
	// Load the Kubernetes configuration from the specified file path
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	// Create an OpenShift client using the configuration
	clientset, err := occlient.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func createOrUpdateSecret(clientset *occlient.Clientset, secret *corev1.Secret) error {
	_, err := clientset.CoreV1().Secrets("default").Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}
