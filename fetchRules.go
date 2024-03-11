package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/jsonpath"
)

func main() {
	// Create a Kubernetes client
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// Name of the ValidatingWebhookConfiguration
	webhookName := "your-webhook-name"

	// Get the ValidatingWebhookConfiguration
	webhookConfig, err := clientset.AdmissionregistrationV1().ValidatingWebhookConfigurations().Get(context.Background(), webhookName, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}

	// Extract and print the rules
	var rules []string
	for _, webhook := range webhookConfig.Webhooks {
		for _, rule := range webhook.Rules {
			rules = append(rules, rule.String())
		}
	}

	// Create JSON representation of the rules
	rulesJSON, err := json.Marshal(rules)
	if err != nil {
		panic(err.Error())
	}

	// Print JSON
	fmt.Println(string(rulesJSON))
}
