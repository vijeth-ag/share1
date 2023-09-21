package main

import (
    "fmt"
    "os"

    "github.com/openshift/api/user/v1"
    "github.com/openshift/client-go/user/clientset/versioned"
    "k8s.io/client-go/rest"
)

func main() {
    // Load the Kubernetes/OpenShift configuration
    config, err := rest.InClusterConfig()
    if err != nil {
        fmt.Printf("Error creating in-cluster config: %v\n", err)
        os.Exit(1)
    }

    // Create a clientset for the OpenShift API
    userClient, err := versioned.NewForConfig(config)
    if err != nil {
        fmt.Printf("Error creating OpenShift client: %v\n", err)
        os.Exit(1)
    }

    // Define the group resource you want to create
    group := &v1.Group{
        // Set the metadata for the group
        ObjectMeta: metav1.ObjectMeta{
            Name: "example-group",
        },
        Users: []string{
            "user1",
            "user2",
        },
    }

    // Create the group resource
    createdGroup, err := userClient.UserV1().Groups().Create(group)
    if err != nil {
        fmt.Printf("Error creating group: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Created group: %s\n", createdGroup.Name)
}
