package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	v1 "github.com/openshift/api/user/v1"
	userv1 "github.com/openshift/client-go/user/clientset/versioned/typed/user/v1"
	"github.com/robfig/cron"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

// User represents the structure of user data fetched from the API
type User struct {
	Username string `json:"username"`
}

// fetchUsersFromAPI fetches the list of users from the API endpoint
func fetchUsersFromAPI() ([]User, error) {
	// Perform HTTP request to fetch users from API
	// resp, err := http.Get("https://example.com/api/users")
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to fetch users from API: %v", err)
	// }
	// defer resp.Body.Close()

	// // Read response body
	// var users []User
	// if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
	// 	return nil, fmt.Errorf("failed to decode JSON response: %v", err)
	// }

	users := []User{}

	vj := User{
		Username: "vj",
	}

	vi := User{
		Username: "vi",
	}

	vh := User{
		Username: "vh",
	}

	users = append(users, vj)
	users = append(users, vi)
	users = append(users, vh)

	return users, nil
}

// addUserToGroup adds the user to the OpenShift group
func addUserToGroup(userClient userv1.UserV1Interface, username string, groupName string) error {
	// Check if the user is already a member of the group
	group, err := userClient.Groups().Get(context.Background(), groupName, metav1.GetOptions{})
	if err != nil {
		// return fmt.Errorf("failed to get group %s: %v", groupName, err)
		log.Println("creatgin grp")

		group = &v1.Group{
			ObjectMeta: metav1.ObjectMeta{Name: "grpname"},
		}

		_, err = userClient.Groups().Create(context.Background(), group, metav1.CreateOptions{})
		if err != nil {
			log.Println("err creating grp", err)
		}
	}

	// Check if the user is already a member
	for _, member := range group.Users {
		if member == username {
			return nil // User is already a member, no need to add again
		}
	}

	// Add the user to the group
	group.Users = append(group.Users, username)

	// Update the group
	_, err = userClient.Groups().Update(context.Background(), group, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update group %s: %v", groupName, err)
	}

	return nil
}

func main() {
	// Set up OpenShift client
	// config, err := rest.InClusterConfig()
	// if err != nil {
	// 	log.Fatalf("Error creating OpenShift config: %v", err)
	// }

	kubeconfig := filepath.Join("/Users/vijeth.ag/.kube/config")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)

	if err != nil {
		log.Println("err getting config", err)
	}

	userClient, err := userv1.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating OpenShift user client: %v", err)
	}

	// Set up cron job
	c := cron.New()

	// Schedule the job to run every 24 hours
	c.AddFunc("@every 3s", func() {
		users, err := fetchUsersFromAPI()
		if err != nil {
			log.Printf("Error fetching users from API: %v", err)
			return
		}

		log.Println(len(users))
		for _, user := range users {
			err := addUserToGroup(userClient, user.Username, "grpname")
			if err != nil {
				log.Printf("Error adding user %s to group: %v", user.Username, err)
			}
		}

		log.Println("Users synced successfully")
	})

	c.Start()

	// Keep the application running
	select {}
}
