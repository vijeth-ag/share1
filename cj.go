package main

import (
    "context"
    "fmt"
    "os"
    "time"

    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/util/homedir"
    "k8s.io/client-go/util/kubeconfig"
    "github.com/robfig/cron/v3"
)

func main() {
    // Initialize cron scheduler.
    c := cron.New()

    // Add a cron job to delete the ServiceAccount every 30 minutes.
    _, err := c.AddFunc("*/30 * * * *", func() {
        deleteServiceAccount()
    })
    if err != nil {
        fmt.Printf("Error adding cron job: %v\n", err)
        os.Exit(1)
    }

    // Start the cron scheduler.
    c.Start()

    // Run the program indefinitely.
    select {}
}

func deleteServiceAccount() {
    // Load Kubernetes configuration from the default location or use the one specified in KUBECONFIG.
    kubeconfigPath := filepath.Join(homedir.HomeDir(), ".kube", "config")
    config, err := kubeconfig.NewForConfig(kubeconfigPath)
    if err != nil {
        fmt.Printf("Error loading kubeconfig: %v\n", err)
        return
    }

    // Create a Kubernetes clientset.
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        fmt.Printf("Error creating clientset: %v\n", err)
        return
    }

    // Define the ServiceAccount name and namespace to delete.
    namespace := "your-namespace" // Replace with the desired namespace.
    serviceAccountName := "your-serviceaccount" // Replace with the desired ServiceAccount name.

    // Delete the ServiceAccount.
    err = clientset.CoreV1().ServiceAccounts(namespace).Delete(context.TODO(), serviceAccountName, metav1.DeleteOptions{})
    if err != nil {
        fmt.Printf("Error deleting ServiceAccount: %v\n", err)
        return
    }

    fmt.Printf("ServiceAccount %s in namespace %s deleted\n", serviceAccountName, namespace)
}
