package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	configv1 "github.com/openshift/api/config/v1"
	ocpclientv1 "github.com/openshift/client-go/config/clientset/versioned/typed/config/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	kubeconfig := flag.String("kubeconfig", filepath.Join(homedir.HomeDir(), ".kube", "config"), "absolute path to the kubeconfig file")
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		fmt.Printf("Error building kubeconfig: %v\n", err)
		os.Exit(1)
	}

	// Create OpenShift client
	openshiftClient, err := ocpclientv1.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error building OpenShift client: %v\n", err)
		os.Exit(1)
	}

	// Create Kubernetes client using the OpenShift client's RESTConfig
	kubeClient, err := kubernetes.NewForConfig(openshiftClient.Config)
	if err != nil {
		fmt.Printf("Error building Kubernetes client: %v\n", err)
		os.Exit(1)
	}

	tokenReview := &configv1.TokenReview{
		Spec: configv1.TokenReviewSpec{
			Token: "YOUR_ACCESS_TOKEN_HERE", // Replace with the actual token you want to review
		},
	}

	result, err := kubeClient.ConfigV1().TokenReviews().Create(context.TODO(), tokenReview, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("Error performing TokenReview: %v\n", err)
		os.Exit(1)
	}

	if result.Status.Authenticated {
		fmt.Println("Token is valid.")
	} else {
		fmt.Println("Token is not valid.")
	}
}
