/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"encoding/json"
	"time"

	http "github.com/LinTechSo/gel-plugin-operator/internal/http"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	lokiv1alpha1 "github.com/LinTechSo/gel-plugin-operator/api/v1alpha1"
)

// GrafanaEnterpriseLogsTokenReconciler reconciles a GrafanaEnterpriseLogsToken object
type GrafanaEnterpriseLogsTokenReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

var TokenStructureResponse struct {
	Name         string    `json:"name"`
	DisplayName  string    `json:"display_name"`
	CreatedAt    time.Time `json:"created_at"`
	Status       string    `json:"status"`
	AccessPolicy string    `json:"access_policy"`
	Expiration   time.Time `json:"expiration"`
	Token        string    `json:"token"`
}

//+kubebuilder:rbac:groups=loki.gelops.com,resources=grafanaenterpriselogstokens,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=loki.gelops.com,resources=grafanaenterpriselogstokens/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=loki.gelops.com,resources=grafanaenterpriselogstokens/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the GrafanaEnterpriseLogsToken object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (r *GrafanaEnterpriseLogsTokenReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// Fetch the GrafanaEnterpriseLogsTenant instance
	var instance = &lokiv1alpha1.GrafanaEnterpriseLogsToken{}
	err := r.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// object not found, could have been deleted after
			// reconcile request, hence don't requeue
			log.Log.Info("GrafanaEnterpriseLogsToken resource not found. Ignoring since object might be deleted")
			return ctrl.Result{}, nil
		}
		// error reading the object, requeue the request
		return ctrl.Result{}, err
	}
	log.Log.Info("Reconciling GrafanaEnterpriseLogsToken", "instance", instance.Name)

	// Add finalizer when creating a new object.
	if instance.ObjectMeta.DeletionTimestamp.IsZero() {
		if !containsString(instance.ObjectMeta.Finalizers, finalizerName) {
			instance.ObjectMeta.Finalizers = append(instance.ObjectMeta.Finalizers, finalizerName)
			if err := r.Update(context.Background(), instance); err != nil {
				return ctrl.Result{}, err
			}
		}
		log.Log.Info("Finalizing GrafanaEnterpriseLogsToken", "instance", instance.Name)

		// Create associated request or perform other cleanup
		err := r.createAssociatedRequesForToken(ctx, instance, err)
		if err != nil {
			log.Log.Error(err, "Failed to create associated request")
			return ctrl.Result{}, err
		}
		log.Log.Info("Created associated request", "instance", instance.Name)
	} else {
		// Object is being deleted.
		if containsString(instance.ObjectMeta.Finalizers, finalizerName) {

			// Example: Delete associated resources or perform other cleanup
			err := r.deleteAssociatedResourcesForToken(ctx, instance, err)
			if err != nil {
				return ctrl.Result{}, err
			}
			log.Log.Info("Deleted associated resources", "instance", instance.Name)

			// Remove the finalizer.
			instance.ObjectMeta.Finalizers = removeString(instance.ObjectMeta.Finalizers, finalizerName)
			if err := r.Update(context.Background(), instance); err != nil {
				return ctrl.Result{}, err
			}
			log.Log.Info("Removed finalizer", "instance", instance.Name)
		}
		// Return to allow Kubernetes to delete the object.
		return ctrl.Result{}, nil
	}

	if err := r.Status().Update(ctx, instance); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// createAssociatedRequesForToken performs the creation of associated token
func (r *GrafanaEnterpriseLogsTokenReconciler) createAssociatedRequesForToken(ctx context.Context, instance *lokiv1alpha1.GrafanaEnterpriseLogsToken, err error) error {
	_ = log.FromContext(ctx)

	var status = "active"
	var tokenMetadataName = instance.ObjectMeta.Name
	var tokenMetadataNamespace = instance.ObjectMeta.Namespace
	var AccessPolicyInfo = instance.Spec.AccessPolicyRef.Name
	var AccessPolicyExpirationInfo = instance.Spec.ExpirationTime
	result, err := http.CreateTokenApiRequest(ctx, tokenMetadataName, tokenMetadataName, AccessPolicyExpirationInfo, AccessPolicyInfo, status, err)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(result.Body)
	err = decoder.Decode(&TokenStructureResponse)
	if err != nil {
		panic(err)
	}

	data := map[string]string{
		"token": TokenStructureResponse.Token,
	}
	// Convert map values to Base64-encoded byte slices
	encodedData := make(map[string][]byte)
	for key, value := range data {
		encodedData[key] = []byte(value)
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      tokenMetadataName,
			Namespace: tokenMetadataNamespace,
		},
		Data: encodedData,
		Type: corev1.SecretTypeOpaque,
	}

	_, err = ctrl.CreateOrUpdate(ctx, r.Client, secret, func() error {
		return controllerutil.SetControllerReference(instance, secret, r.Scheme)
	})
	if err != nil {
		return err
	}

	return nil
}

// deleteAssociatedResourcesForToken performs the cleanup of associated resources
func (r *GrafanaEnterpriseLogsTokenReconciler) deleteAssociatedResourcesForToken(ctx context.Context, instance *lokiv1alpha1.GrafanaEnterpriseLogsToken, err error) error {
	_ = log.FromContext(ctx)

	var status = "inactive"
	var tokenMetadataName = instance.ObjectMeta.Name
	_, err = http.DeleteToken(ctx, tokenMetadataName, status, err)
	if err != nil {
		return err
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GrafanaEnterpriseLogsTokenReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&lokiv1alpha1.GrafanaEnterpriseLogsToken{}).
		Complete(r)
}
