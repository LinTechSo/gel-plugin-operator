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

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	http "github.com/LinTechSo/gel-plugin-operator/internal/http"

	lokiv1alpha1 "github.com/LinTechSo/gel-plugin-operator/api/v1alpha1"
)

// GrafanaEnterpriseLogsAccessPolicyReconciler reconciles a GrafanaEnterpriseLogsAccessPolicy object
type GrafanaEnterpriseLogsAccessPolicyReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=loki.hamravesh.com,resources=grafanaenterpriselogsaccesspolicies,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=loki.hamravesh.com,resources=grafanaenterpriselogsaccesspolicies/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=loki.hamravesh.com,resources=grafanaenterpriselogsaccesspolicies/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the GrafanaEnterpriseLogsAccessPolicy object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (r *GrafanaEnterpriseLogsAccessPolicyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// Fetch the GrafanaEnterpriseLogsTenant instance
	var instance = &lokiv1alpha1.GrafanaEnterpriseLogsAccessPolicy{}
	err := r.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// object not found, could have been deleted after
			// reconcile request, hence don't requeue
			log.Log.Info("GrafanaEnterpriseLogsAccessPolicy resource not found. Ignoring since object might be deleted")
			return ctrl.Result{}, nil
		}
		// error reading the object, requeue the request
		return ctrl.Result{}, err
	}
	log.Log.Info("Reconciling GrafanaEnterpriseLogsAccessPolicy", "instance", instance.Name)

	// initialize realm
	realm, scopes, err := r.jsondataImplementation(ctx, instance)
	if err != nil {
		log.Log.Error(err, "Failed to prepare the desired state of GrafanaEnterpriseLogs AccessPolicy realm and scopes")
		return ctrl.Result{}, err
	}
	log.Log.Info("Prepare the desired state of GrafanaEnterpriseLogs AccessPolicy realm and scopes")

	// Add your finalizer when creating a new object.
	if instance.ObjectMeta.DeletionTimestamp.IsZero() {
		if !containsString(instance.ObjectMeta.Finalizers, finalizerName) {
			instance.ObjectMeta.Finalizers = append(instance.ObjectMeta.Finalizers, finalizerName)
			if err := r.Update(context.Background(), instance); err != nil {
				return ctrl.Result{}, err
			}
		}
		// Create associated request or perform other cleanup
		err := r.createAssociatedRequestForAccessPolicy(ctx, instance, realm, scopes)
		if err != nil {
			log.Log.Error(err, "Failed to create associated request", "instance", instance.Name)
			return ctrl.Result{}, err
		}
		log.Log.Info("Create associated request", "instance", instance.Name)

	} else {
		// Object is being deleted.
		if containsString(instance.ObjectMeta.Finalizers, finalizerName) {

			// Delete associated resources or perform other cleanup
			err := r.deleteAssociatedResources(ctx, instance, realm, scopes)
			if err != nil {
				log.Log.Error(err, "Failed to delete associated resources")
				return ctrl.Result{}, err
			}
			log.Log.Info("Delete associated resources", "instance", instance.Name)

			// Remove the finalizer.
			instance.ObjectMeta.Finalizers = removeString(instance.ObjectMeta.Finalizers, finalizerName)
			if err := r.Update(context.Background(), instance); err != nil {
				log.Log.Error(err, "Failed to remove the finalizer")
				return ctrl.Result{}, err
			}
			log.Log.Info("Remove finalizer", "instance", instance.Name)
		}
		// Return to allow Kubernetes to delete the object.
		return ctrl.Result{}, nil
	}
	return ctrl.Result{}, nil
}

// jsondataImplementation performs the creation of associated json data
func (r *GrafanaEnterpriseLogsAccessPolicyReconciler) jsondataImplementation(ctx context.Context, instance *lokiv1alpha1.GrafanaEnterpriseLogsAccessPolicy) (realm map[string]interface{}, scopes []string, err error) {
	_ = log.FromContext(ctx)

	// Define the accessScope
	var ac = instance.Spec.TenantInfoRef.AccessPolicies
	for _, acpo := range ac {
		scopes = append(scopes, acpo)
	}
	log.Log.Info("Define the accessScope list", "instance", instance.Name)

	// Define the labelSelector
	var ls = instance.Spec.TenantInfoRef.LabelSelectors
	labelPolicies := make([]map[string]string, 0)
	for _, selector := range ls {
		policy := map[string]string{
			"selector": selector,
		}
		labelPolicies = append(labelPolicies, policy)
	}
	log.Log.Info("Define the labelSelector map", "labelSelector", labelPolicies, "instance", instance.Name)

	// Create a json list request from the GrafanaEnterpriseLogsAccessPolicy object
	var tenant = instance.Spec.TenantInfoRef
	realm = map[string]interface{}{
		"cluster":        tenant.ClusterName,
		"label_policies": labelPolicies,
		"tenant":         tenant.TenantName,
	}

	log.Log.Info("Define the realm", "realm", realm, "scopes", scopes, "instance", instance.Name)
	return realm, scopes, nil
}

// createAssociatedRequest performs the creation of associated resources
func (r *GrafanaEnterpriseLogsAccessPolicyReconciler) createAssociatedRequestForAccessPolicy(ctx context.Context, instance *lokiv1alpha1.GrafanaEnterpriseLogsAccessPolicy, realm map[string]interface{}, scopes []string) error {
	_ = log.FromContext(ctx)

	var metadataName = instance.ObjectMeta.Name
	data := map[string]interface{}{
		"name":         metadataName,
		"display_name": metadataName,
		"status":       "active",
		"realms":       []map[string]interface{}{realm},
		"scopes":       scopes,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Log.Error(err, "Error marshaling", "data", data)
		return err
	}
	log.Log.Info("Create Associated Resources", "instance", instance.Name)

	_, err = http.CreateAccessPolicyApiRequest(ctx, jsonData, err)
	if err != nil {
		return err
	}

	log.Log.Info("Requesting AccessPolicy creation", "instance", instance.Name)
	return nil
}

// deleteAssociatedResources performs the cleanup of associated resources
func (r *GrafanaEnterpriseLogsAccessPolicyReconciler) deleteAssociatedResources(ctx context.Context, instance *lokiv1alpha1.GrafanaEnterpriseLogsAccessPolicy, realm map[string]interface{}, scopes []string) error {
	// Your logic for deleting associated resources here...
	var metadataName = instance.ObjectMeta.Name
	data := map[string]interface{}{
		"display_name": metadataName,
		"status":       "inactive",
		"realms":       []map[string]interface{}{realm},
		"scopes":       scopes,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Log.Error(err, "Error marshaling", "data", data)
		return err
	}
	log.Log.Info("Delete Associated Resources", "instance", instance.Name)

	_, err = http.DeleteAccessPolicy(ctx, jsonData, metadataName, err)
	if err != nil {
		return err
	}

	log.Log.Info("Requesting AccessPolicy deletion", "instance", instance.Name)
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GrafanaEnterpriseLogsAccessPolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&lokiv1alpha1.GrafanaEnterpriseLogsAccessPolicy{}).
		Complete(r)
}
