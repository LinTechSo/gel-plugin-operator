package controller

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func CreateAccessPolicyApiRequest(AccessPolicyData []byte, err error) (string, error) {
	var AccessPolicyUrlPrefix = "/admin/api/v3/accesspolicies"
	var LokiEndpointAddress = os.Getenv("Loki_Endpoint_Address")
	var LokiAdminApiToken = os.Getenv("Loki_Admin_Api_Token")
	// Check if the environment variable is set
	if LokiEndpointAddress == "" && LokiAdminApiToken == "" {
		fmt.Println("ENVs are not set")
		return "", err
	}
	var Address = LokiEndpointAddress + AccessPolicyUrlPrefix
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

func CreateTenantApiRequest(TenantName string, DisplayName string, ClusterName string, status string, err error) (string, error) {
	var tenantUrlPrefix = "/admin/api/v3/tenants"
	var LokiEndpointAddress = os.Getenv("Loki_Endpoint_Address")
	var LokiAdminApiToken = os.Getenv("Loki_Admin_Api_Token")
	// Check if the environment variable is set
	if LokiEndpointAddress == "" && LokiAdminApiToken == "" {
		fmt.Println("ENVs are not set")
		return "", err
	}

	var Address = LokiEndpointAddress + tenantUrlPrefix
	// fmt.Println(Address)
	// fmt.Println(LokiAdminApiToken)

	username := ""
	password := string(LokiAdminApiToken)

	auth := username + ":" + password
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))

	postData := strings.NewReader(
		fmt.Sprintf(
			`{"name":"%s","display_name":"%s","cluster":"%s", "status":"%s"}`, TenantName, DisplayName, ClusterName, status),
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

func DeleteTenant(TenantName string, DisplayName string, ClusterName string, status string, err error) (string, error) {
	var tenantUrlPrefix = "/admin/api/v3/tenants/"
	var LokiEndpointAddress = os.Getenv("Loki_Endpoint_Address")
	var LokiAdminApiToken = os.Getenv("Loki_Admin_Api_Token")
	// Check if the environment variable is set
	if LokiEndpointAddress == "" && LokiAdminApiToken == "" {
		fmt.Println("ENVs are not set")
		return "", err
	}

	var Address = LokiEndpointAddress + tenantUrlPrefix + TenantName
	fmt.Println(Address)
	username := ""
	password := string(LokiAdminApiToken)

	auth := username + ":" + password
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))

	postData := strings.NewReader(
		fmt.Sprintf(
			`{"cluster":"%s", "status":"%s"}`, ClusterName, status),
	)
	fmt.Println(postData)
	resp, err := http.NewRequest("PUT", Address, postData)
	if err != nil {
		// Handle the error
		fmt.Println(err, "unable to request")
		return "", err
	}
	resp.Header.Set("Authorization", basicAuth)
	resp.Header.Add("If-Match", "1")
	resp.Header.Add("Content-Type", "text/plain")

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

func DeleteAccessPolicy(jsonData []byte, tenant string, err error) (string, error) {
	var AccessPolicyUrlPrefix = "/admin/api/v3/accesspolicies/"
	var LokiEndpointAddress = os.Getenv("Loki_Endpoint_Address")
	var LokiAdminApiToken = os.Getenv("Loki_Admin_Api_Token")
	// Check if the environment variable is set
	if LokiEndpointAddress == "" && LokiAdminApiToken == "" {
		fmt.Println("ENVs are not set")
		return "", err
	}
	var Address = LokiEndpointAddress + AccessPolicyUrlPrefix + tenant
	username := ""
	password := string(LokiAdminApiToken)

	auth := username + ":" + password
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))

	data := strings.NewReader(string(jsonData))

	fmt.Println(data)
	resp, err := http.NewRequest("PUT", Address, data)
	if err != nil {
		// Handle the error
		fmt.Println(err, "unable to request")
		return "", err
	}
	resp.Header.Set("Authorization", basicAuth)
	resp.Header.Add("If-Match", "1")
	resp.Header.Add("Content-Type", "text/plain")

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
