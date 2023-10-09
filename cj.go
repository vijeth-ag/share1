package main

import (
    "context"
    "fmt"
    "os"
    "time"

    batchv1 "k8s.io/api/batch/v1beta1"
    "k8s.io/apimachinery/pkg/util/intstr"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/util/homedir"
    "k8s.io/client-go/util/kubeconfig"
)

func main() {
    // Load Kubernetes configuration from the default location or use the one specified in KUBECONFIG.
    kubeconfigPath := filepath.Join(homedir.HomeDir(), ".kube", "config")
    config, err := kubeconfig.NewForConfig(kubeconfigPath)
    if err != nil {
        fmt.Printf("Error loading kubeconfig: %v\n", err)
        os.Exit(1)
    }

    // Create a Kubernetes clientset.
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        fmt.Printf("Error creating clientset: %v\n", err)
        os.Exit(1)
    }

    // Define the CronJob object.
    cronJob := &batchv1.CronJob{
        ObjectMeta: metav1.ObjectMeta{
            Name: "my-cronjob",
        },
        Spec: batchv1.CronJobSpec{
            Schedule: "*/1 * * * *", // Replace with your desired schedule.
            JobTemplate: batchv1.JobTemplateSpec{
                Spec: batchv1.JobSpec{
                    Template: corev1.PodTemplateSpec{
                        Spec: corev1.PodSpec{
                            Containers: []corev1.Container{
                                {
                                    Name:  "my-container",
                                    Image: "my-image:latest", // Replace with your container image.
                                },
                            },
                        },
                    },
                },
            },
        },
    }

    // Create the CronJob in the desired namespace.
    namespace := "default" // Replace with your desired namespace.
    createdCronJob, err := clientset.BatchV1beta1().CronJobs(namespace).Create(context.TODO(), cronJob, metav1.CreateOptions{})
    if err != nil {
        fmt.Printf("Error creating CronJob: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("CronJob created with name: %s\n", createdCronJob.Name)
}
