package main

import (
	"context"
	"fmt"

	"github.com/openshift/client-go/config/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

func main() {
	// Create an OpenShift client
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(fmt.Errorf("failed to create OpenShift config: %v", err))
	}

	client, err := versioned.NewForConfig(config)
	if err != nil {
		panic(fmt.Errorf("failed to create OpenShift client: %v", err))
	}

	// Fetch the OAuth resource
	oauth, err := client.ConfigV1().OAuths().Get(context.Background(), "cluster", metav1.GetOptions{})
	if err != nil {
		panic(fmt.Errorf("failed to get OAuth resource: %v", err))
	}

	// Modify the OAuth settings
	oauth.Spec.IdentityProviders[0].Name = "mywebhook" // Assuming the first identity provider is the webhook

	// Update the OAuth resource
	_, err = client.ConfigV1().OAuths().Update(context.Background(), oauth, metav1.UpdateOptions{})
	if err != nil {
		panic(fmt.Errorf("failed to update OAuth resource: %v", err))
	}

	fmt.Println("OAuth resource updated successfully")
}
