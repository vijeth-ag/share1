package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ServiceNowResponse struct {
	Result []struct {
		SysID  string `json:"sys_id"`
		Number string `json:"number"`
		// Add more fields here based on your requirements
	} `json:"result"`
}

func getSysIDByREQID(instanceURL, username, password, tableName, reqID string) (string, error) {
	// ServiceNow API endpoint URL for the specified table
	apiURL := fmt.Sprintf("%s/api/now/table/%s", instanceURL, tableName)

	// Set up basic authentication with username and password
	client := &http.Client{}
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(username, password)

	// Set the filter to retrieve the record with the specified REQID
	req.URL.RawQuery = "sysparm_query=number=" + reqID

	// Send GET request to ServiceNow API
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch data. Status code: %d", resp.StatusCode)
	}

	var snResp ServiceNowResponse
	if err := json.NewDecoder(resp.Body).Decode(&snResp); err != nil {
		return "", err
	}

	// Check if a record is found with the provided REQID
	if len(snResp.Result) > 0 {
		return snResp.Result[0].SysID, nil
	}

	return "", nil // No record found for the given REQID
}

func main() {
	// Replace these variables with your ServiceNow instance URL, username, password, table name, and REQID
	instanceURL := "https://dev168296.service-now.com"
	username := "admin"
	password := "WPmfr7sMB*%1"
	tableName := "sc_request" // Change to the appropriate table name (e.g., "incident" or "sc_request")
	reqID := "REQ00w10001"

	// Fetch the sys_id of the specified REQID from the specified table
	sysID, err := getSysIDByREQID(instanceURL, username, password, tableName, reqID)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if sysID != "" {
		fmt.Printf("Sys_id of REQID %s: %s\n", reqID, sysID)
	} else {
		fmt.Println("REQID not found.")
	}
}
