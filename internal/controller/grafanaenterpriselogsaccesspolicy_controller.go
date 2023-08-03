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
	"fmt"

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

	rlog := log.Log.WithValues("loki", req.NamespacedName)
	rlog.Info("Reconciling AccessPolicy creation request")

	// Fetch the GrafanaEnterpriseLogsAccessPolicy instance
	instance := &lokiv1alpha1.GrafanaEnterpriseLogsAccessPolicy{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if client.IgnoreNotFound(err) != nil {
			rlog.Error(err, "unable to list GrafanaEnterpriseLogsAccessPolicy")
			return ctrl.Result{}, err
		}
		// Request object not found, could have been deleted after reconcile request.
		// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
		// Return and don't requeue
		return ctrl.Result{}, nil
	}
	// Define the accessScope
	var ac = instance.Spec.TenantInfoRef.AccessPolicies
	var scopes []string
	for _, acpo := range ac {
		scopes = append(scopes, acpo)
	}

	// Define the labelSelector
	var ls = instance.Spec.TenantInfoRef.LabelSelectors
	labelPolicies := make([]map[string]string, 0)
	for _, selector := range ls {
		policy := map[string]string{
			"selector": selector,
		}
		labelPolicies = append(labelPolicies, policy)
	}

	// Create a json list request from the GrafanaEnterpriseLogsAccessPolicy object
	var tenant = instance.Spec.TenantInfoRef
	realm := map[string]interface{}{
		"cluster":        tenant.ClusterName,
		"label_policies": labelPolicies,
		"tenant":         tenant.TenantName,
	}

	data := map[string]interface{}{
		"name":         tenant.TenantName,
		"display_name": tenant.TenantName,
		"status":       "active",
		"realms":       []map[string]interface{}{realm},
		"scopes":       scopes,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling:", err)
		return ctrl.Result{}, err
	}

	rlog.Info("Requesting AccessPolicy creation")
	_, err = http.CreateAccessPolicyApiRequest(jsonData, err)
	if err != nil {
		return ctrl.Result{}, err
	}

	if instance.ObjectMeta.DeletionTimestamp.IsZero() {
		// Object is being deleted, perform cleanup or other actions here

		// Example: Delete associated resources or perform other cleanup
		err := r.deleteAssociatedResources(ctx, instance)
		if err != nil {
			return ctrl.Result{}, err
		}

		// Object has been deleted, no need to requeue
		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, nil
}

// deleteAssociatedResources performs the cleanup of associated resources
func (r *GrafanaEnterpriseLogsAccessPolicyReconciler) deleteAssociatedResources(ctx context.Context, obj *lokiv1alpha1.GrafanaEnterpriseLogsAccessPolicy) error {
	// Your logic for deleting associated resources here...
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GrafanaEnterpriseLogsAccessPolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&lokiv1alpha1.GrafanaEnterpriseLogsAccessPolicy{}).
		Complete(r)
}
