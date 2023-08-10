package controller

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"

	"context"

	"sigs.k8s.io/controller-runtime/pkg/log"
)

// export Loki_Endpoint_Api_Version=v3
// export Loki_Admin_Api_Token=X
// export Loki_Endpoint_Address=http://localhost:3100
func ReadEnvironmentVariables(ctx context.Context, err error) (string, string, string, error) {
	_ = log.FromContext(ctx)

	var LokiEndpointApiVersion = os.Getenv("Loki_Endpoint_Api_Version")
	if LokiEndpointApiVersion == "v3" {
		var UrlPrefix = "/admin/api/v3/"
		var LokiEndpointAddress = os.Getenv("Loki_Endpoint_Address")
		var LokiAdminApiToken = os.Getenv("Loki_Admin_Api_Token")
		// Check if the environment variable is set
		if LokiEndpointAddress == "" && LokiAdminApiToken == "" {
			log.Log.Error(err, "Loki Endpoint Address or Loki Admin Api Token is not set")
			return "", "", "", err
		}

		return LokiAdminApiToken, LokiEndpointAddress, UrlPrefix, nil
	} else {
		log.Log.Info("Loki Endpoint Api Version v3 just supported at this time")
	}
	log.Log.Error(err, "Loki Admin API Version is not set")
	return "", "", "", err
}

// Creation methods
func CreateAccessPolicyApiRequest(ctx context.Context, AccessPolicyData []byte, err error) (string, error) {
	_ = log.FromContext(ctx)

	LokiAdminApiToken, LokiEndpointAddress, UrlPrefix, err := ReadEnvironmentVariables(ctx, err)
	var Address = LokiEndpointAddress + UrlPrefix + "accesspolicies/"
	log.Log.Info("Endpoint Address", "URL", Address)

	username := ""
	password := string(LokiAdminApiToken)

	auth := username + ":" + password
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))

	data := strings.NewReader(string(AccessPolicyData))
	fmt.Println("Create AccessPolicy API request", data)

	resp, err := http.NewRequest("POST", Address, data)
	if err != nil {
		log.Log.Error(err, "Failed to create HTTP request", "URL", Address)
		return "", err
	}
	resp.Header.Set("Authorization", basicAuth)
	resp.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	result, err := client.Do(resp)
	if err != nil {
		log.Log.Error(err, "Failed Response", "result", result)
		return "", err
	}
	defer resp.Body.Close()

	// Check the response status code and handle accordingly
	if result.StatusCode != http.StatusOK {
		log.Log.Info("Failed status code")
	}

	log.Log.Info("Access Policy request code", " status", result.Status)
	return "", err
}

func CreateTenantApiRequest(ctx context.Context, TenantName string, DisplayName string, ClusterName string, status string, err error) (string, error) {
	_ = log.FromContext(ctx)

	LokiAdminApiToken, LokiEndpointAddress, UrlPrefix, err := ReadEnvironmentVariables(ctx, err)
	var Address = LokiEndpointAddress + UrlPrefix + "tenants/"
	log.Log.Info("Endpoint Address", "URL", Address)

	username := ""
	password := string(LokiAdminApiToken)

	auth := username + ":" + password
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))

	postData := strings.NewReader(
		fmt.Sprintf(
			`{"name":"%s","display_name":"%s","cluster":"%s", "status":"%s"}`, TenantName, DisplayName, ClusterName, status),
	)
	fmt.Println("Create tenant API request", postData)

	resp, err := http.NewRequest("POST", Address, postData)
	if err != nil {
		log.Log.Error(err, "unable to request", "result", resp)
		return "", err
	}
	resp.Header.Set("Authorization", basicAuth)
	resp.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	result, err := client.Do(resp)
	if err != nil {
		log.Log.Error(err, "unable to get client result", "result", result)
		return "", err
	}
	defer resp.Body.Close()

	// Check the response status code and handle accordingly
	if result.StatusCode != http.StatusOK {
		log.Log.Info("Failed status ")
	}

	log.Log.Info("tenant request ", "code status", result.Status)
	return "", err
}

func CreateTokenApiRequest(ctx context.Context, TokenName string, DisplayName string, TokenExpiration string, AccessPolicyName string, status string, err error) (*http.Response, error) {
	_ = log.FromContext(ctx)

	LokiAdminApiToken, LokiEndpointAddress, UrlPrefix, err := ReadEnvironmentVariables(ctx, err)
	var Address = LokiEndpointAddress + UrlPrefix + "tokens/"
	log.Log.Info("Endpoint Address", "URL", Address)

	username := ""
	password := string(LokiAdminApiToken)

	auth := username + ":" + password
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))

	postData := strings.NewReader(
		fmt.Sprintf(
			`{"name":"%s","display_name":"%s","expiration":"%s", "access_policy":"%s", "status":"%s"}`, TokenName, DisplayName, TokenExpiration, AccessPolicyName, status),
	)
	fmt.Println("Create tenant API request", postData)

	resp, err := http.NewRequest("POST", Address, postData)
	if err != nil {
		log.Log.Error(err, "unable to request", "result", resp)
		return nil, err
	}
	resp.Header.Set("Authorization", basicAuth)
	resp.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	result, err := client.Do(resp)
	if err != nil {
		log.Log.Error(err, "unable to get client result", "result", result)
		return nil, err
	}
	defer resp.Body.Close()

	// Check the response status code and handle accordingly
	if result.StatusCode != http.StatusOK {
		log.Log.Info("Failed status ")
	}
	log.Log.Info("token request ", "result", result)
	log.Log.Info("token request ", "code status", result.Status)
	return result, err
}

// Deletion methods
func DeleteTenant(ctx context.Context, TenantName string, ClusterName string, status string, err error) (string, error) {
	_ = log.FromContext(ctx)

	LokiAdminApiToken, LokiEndpointAddress, UrlPrefix, err := ReadEnvironmentVariables(ctx, err)
	var Address = LokiEndpointAddress + UrlPrefix + "tenants/" + TenantName
	log.Log.Info("Endpoint Address", "URL", Address)

	username := ""
	password := string(LokiAdminApiToken)

	auth := username + ":" + password
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))

	postData := strings.NewReader(
		fmt.Sprintf(
			`{"cluster":"%s", "status":"%s"}`, ClusterName, status),
	)
	fmt.Println("DELETE tenant API request", postData)

	resp, err := http.NewRequest("PUT", Address, postData)
	if err != nil {
		// Handle the error
		log.Log.Error(err, "unable to request", "URL", Address)
		return "", err
	}
	resp.Header.Set("Authorization", basicAuth)
	resp.Header.Add("If-Match", "1")
	resp.Header.Add("Content-Type", "text/plain")

	client := &http.Client{}
	result, err := client.Do(resp)
	if err != nil {
		log.Log.Error(err, "unable to get client result", "URL", Address)
		return "", err
	}
	defer resp.Body.Close()

	// Check the response status code and handle accordingly
	if result.StatusCode != http.StatusOK {
		// Handle non-OK status codes
		log.Log.Info("Failed status ")
	}

	log.Log.Info("tenant request ", "code status", result.Status)
	return "", err
}

func DeleteAccessPolicy(ctx context.Context, jsonData []byte, tenant string, err error) (string, error) {
	_ = log.FromContext(ctx)

	LokiAdminApiToken, LokiEndpointAddress, UrlPrefix, err := ReadEnvironmentVariables(ctx, err)
	var Address = LokiEndpointAddress + UrlPrefix + "accesspolicies/" + tenant
	log.Log.Info("Endpoint Address", "URL", Address)

	username := ""
	password := string(LokiAdminApiToken)

	auth := username + ":" + password
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))

	data := strings.NewReader(string(jsonData))
	fmt.Println("DELETE Access Policy API request", data)

	resp, err := http.NewRequest("PUT", Address, data)
	if err != nil {
		log.Log.Error(err, "unable to request", "URL", Address)
		return "", err
	}
	resp.Header.Set("Authorization", basicAuth)
	resp.Header.Add("If-Match", "1")
	resp.Header.Add("Content-Type", "text/plain")

	client := &http.Client{}
	result, err := client.Do(resp)
	if err != nil {
		log.Log.Error(err, "unable to get client result", "URL", Address)
		return "", err
	}
	defer resp.Body.Close()

	// Check the response status code and handle accordingly
	if result.StatusCode != http.StatusOK {
		log.Log.Info("Failed status code")
	}

	log.Log.Info("Access Policy deletion request code", "status", result.Status)
	return "", err
}

func DeleteToken() {
	fmt.Println("Hello from DeleteToken")
}
