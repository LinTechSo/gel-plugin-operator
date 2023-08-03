package controller

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// test Sink
// export Loki_Endpoint_Address="http://45.129.36.204:31178"
// export Loki_Admin_Api_Token="X19hZG1pbl9fLWUyZGVmMWJkZmUwZTc2MmU6Yj5eMHUwMDE4bCI3LnsyKzBrPTQzIi0+"
func CreateAccessPolicyApiRequest(AccessPolicyData []byte, err error) (string, error) {
	var AccessPolicyUrlPrefix = "/admin/api/v3/accesspolicies"
	var LokiEndpointAddress = os.Getenv("Loki_Endpoint_Address")
	var LokiAdminApiToken = os.Getenv("Loki_Admin_Api_Token")
	var Address = LokiEndpointAddress + AccessPolicyUrlPrefix

	// Check if the environment variable is set
	if LokiEndpointAddress == "" {
		fmt.Println("Loki_Endpoint_Address is not set")
	}
	if LokiAdminApiToken == "" {
		fmt.Println("Loki_Admin_Api_Token is not set")
	}

	username := ""
	password := string(LokiAdminApiToken)

	auth := username + ":" + password
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))

	data := strings.NewReader(string(AccessPolicyData))

	fmt.Println(data)
	resp, err := http.NewRequest("POST", Address, data)
	if err != nil {
		// Handle the error
		fmt.Println(err, "unable to request")
		return "", err
	}
	resp.Header.Set("Authorization", basicAuth)
	resp.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	result, err := client.Do(resp)
	if err != nil {
		// Handle the error
		fmt.Println(err, "unable to get client result")
		return "", err
	}
	defer resp.Body.Close()

	// Check the response status code and handle accordingly
	if result.StatusCode != http.StatusOK {
		// Handle non-OK status codes
	}

	fmt.Println(result.Status)

	return "", err
}

func CreateTenantApiRequest(TenantName string, DisplayName string, ClusterName string, err error) (string, error) {
	var tenantUrlPrefix = "/admin/api/v3/tenants"
	var LokiEndpointAddress = os.Getenv("Loki_Endpoint_Address")
	var LokiAdminApiToken = os.Getenv("Loki_Admin_Api_Token")
	var Address = LokiEndpointAddress + tenantUrlPrefix
	// Check if the environment variable is set
	if LokiEndpointAddress == "" {
		fmt.Println("Loki_Endpoint_Address is not set")
	}
	if LokiAdminApiToken == "" {
		fmt.Println("Loki_Admin_Api_Token is not set")
	}

	// fmt.Println(Address)
	// fmt.Println(LokiAdminApiToken)

	username := ""
	password := string(LokiAdminApiToken)

	auth := username + ":" + password
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))

	postData := strings.NewReader(
		fmt.Sprintf(
			`{"name":"%s","display_name":"%s","cluster":"%s"}`, TenantName, DisplayName, ClusterName),
	)
	resp, err := http.NewRequest("POST", Address, postData)
	if err != nil {
		// Handle the error
		fmt.Println(err, "unable to request")
		return "", err
	}
	resp.Header.Set("Authorization", basicAuth)
	resp.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	result, err := client.Do(resp)
	if err != nil {
		// Handle the error
		fmt.Println(err, "unable to get client result")
		return "", err
	}
	defer resp.Body.Close()

	// Check the response status code and handle accordingly
	if result.StatusCode != http.StatusOK {
		// Handle non-OK status codes
	}

	fmt.Println(result.Status)

	return "", err
}

func CreateTokenApiRequest() {
	fmt.Println("Hello from CreateTokenApiRequest")
}
