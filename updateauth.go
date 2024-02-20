package main

import (
	"context"
	"fmt"

	userv1 "github.com/openshift/client-go/user/clientset/versioned/typed/user/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

func main() {
	// Create an OpenShift client
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(fmt.Errorf("failed to create OpenShift config: %v", err))
	}

	userClient, err := userv1.NewForConfig(config)
	if err != nil {
		panic(fmt.Errorf("failed to create OpenShift user client: %v", err))
	}

	// Fetch the existing Authentication object
	auth, err := userClient.Authentications().Get(context.Background(), "authentication-name", metav1.GetOptions{})
	if err != nil {
		panic(fmt.Errorf("failed to get Authentication object: %v", err))
	}

	// Modify the Authentication object
	auth.Spec.Type = "None"
	auth.Spec.WebhookAuthenticator.Kubeconfig.Name = "mywebhook"

	// Update the Authentication object
	_, err = userClient.Authentications().Update(context.Background(), auth, metav1.UpdateOptions{})
	if err != nil {
		panic(fmt.Errorf("failed to update Authentication object: %v", err))
	}

	fmt.Println("Authentication object updated successfully")
}
