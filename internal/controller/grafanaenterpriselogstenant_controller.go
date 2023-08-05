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

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	http "github.com/LinTechSo/gel-plugin-operator/internal/http"

	lokiv1alpha1 "github.com/LinTechSo/gel-plugin-operator/api/v1alpha1"
)

// GrafanaEnterpriseLogsTenantReconciler reconciles a GrafanaEnterpriseLogsTenant object
type GrafanaEnterpriseLogsTenantReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=loki.hamravesh.com,resources=grafanaenterpriselogstenants,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=loki.hamravesh.com,resources=grafanaenterpriselogstenants/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=loki.hamravesh.com,resources=grafanaenterpriselogstenants/finalizers,verbs=update
const finalizerName = "hamravesh.com/finalizer"

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the GrafanaEnterpriseLogsTenant object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (r *GrafanaEnterpriseLogsTenantReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// Fetch the GrafanaEnterpriseLogsTenant instance
	var instance = &lokiv1alpha1.GrafanaEnterpriseLogsTenant{}
	err := r.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// object not found, could have been deleted after
			// reconcile request, hence don't requeue
			return ctrl.Result{}, nil
		}
		// error reading the object, requeue the request
		return ctrl.Result{}, err
	}
	log.Log.Info("Reconciling GrafanaEnterpriseLogsTenant", "instance", instance.Name)

	// Add your finalizer when creating a new object.
	if instance.ObjectMeta.DeletionTimestamp.IsZero() {
		if !containsString(instance.ObjectMeta.Finalizers, finalizerName) {
			instance.ObjectMeta.Finalizers = append(instance.ObjectMeta.Finalizers, finalizerName)
			if err := r.Update(context.Background(), instance); err != nil {
				return ctrl.Result{}, err
			}
		}
		log.Log.Info("Finalizing GrafanaEnterpriseLogsTenant", "instance", instance.Name)

		// Create associated request or perform other cleanup
		err := r.createAssociatedRequesForTenant(ctx, instance, err)
		if err != nil {
			log.Log.Error(err, "Failed to create associated request")
			return ctrl.Result{}, err
		}
		log.Log.Info("Created associated request", "instance", instance.Name)
	} else {
		// Object is being deleted.
		if containsString(instance.ObjectMeta.Finalizers, finalizerName) {

			// Example: Delete associated resources or perform other cleanup
			err := r.deleteAssociatedResources(ctx, instance, err)
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

// createAssociatedRequesForTenant performs the creation of associated tenant
func (r *GrafanaEnterpriseLogsTenantReconciler) createAssociatedRequesForTenant(ctx context.Context, instance *lokiv1alpha1.GrafanaEnterpriseLogsTenant, err error) error {
	_ = log.FromContext(ctx)
	var tenant = instance.Spec.TenantInfo
	var status = "active"
	var metadataName = instance.ObjectMeta.Name

	_, err = http.CreateTenantApiRequest(ctx, tenant.Name, metadataName, tenant.ClusterName, status, err)
	if err != nil {
		return err
	}

	log.Log.Info("Create tenant API request", "tenant", tenant.Name, "instance", instance.Name)
	return nil
}

// deleteAssociatedResources performs the cleanup of associated resources
func (r *GrafanaEnterpriseLogsTenantReconciler) deleteAssociatedResources(ctx context.Context, instance *lokiv1alpha1.GrafanaEnterpriseLogsTenant, err error) error {
	_ = log.FromContext(ctx)

	// Your logic for deleting associated resources here...
	var status = "inactive"
	var tenantInfo = instance.Spec.TenantInfo
	_, err = http.DeleteTenant(ctx, tenantInfo.Name, tenantInfo.ClusterName, status, err)
	if err != nil {
		return err
	}

	log.Log.Info("Delete tenant API request", "tenant", tenantInfo.Name, "instance", instance.Name)
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GrafanaEnterpriseLogsTenantReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&lokiv1alpha1.GrafanaEnterpriseLogsTenant{}).
		Complete(r)
}
