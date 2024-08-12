/*
Copyright 2024.

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
	"fmt"

	"github.com/go-logr/logr"
	"github.com/javdet/networksets-controller/monitoring"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/log"

	calicov3 "github.com/projectcalico/api/pkg/apis/projectcalico/v3"
)

// NetworkPolicyReconciler reconciles a Networkset object
type GlobalNetworkPolicyReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

var controllerGlobalNetworksetsLog = ctrl.Log.WithName("controller").WithName("GlobalNetworkpolicy")

//+kubebuilder:rbac:groups=crd.projectcalico.org,resources=globalnetworkpolicies,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=crd.projectcalico.org,resources=globalnetworkpolicies/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=crd.projectcalico.org,resources=globalnetworkpolicies/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Memcached object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *GlobalNetworkPolicyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	controllerGlobalNetworksetsLog.Info("start reconcile", "request", req.NamespacedName)

	instance := &calicov3.GlobalNetworkPolicy{}
	globalNetworkSetList := &calicov3.GlobalNetworkSetList{}

	opts := []client.ListOption{
		client.MatchingLabels{
			"control-plane":        "networksets-operator",
			"parent-networkPolicy": req.NamespacedName.Name,
		},
	}
	err := r.List(ctx, globalNetworkSetList, opts...)
	if err != nil {
		controllerGlobalNetworksetsLog.Error(err, "cannot get object GlobalNetworkSet list")
		return ctrl.Result{}, err
	}
	err = r.Get(ctx, req.NamespacedName, instance)

	if err != nil {
		controllerGlobalNetworksetsLog.Error(err, "cannot get object GlobalNetworkPolicy")
		if apierrors.IsNotFound(err) {
			deleteGlobalNetworkSetList := getGlobalNetworkSetList(req.NamespacedName.Name, req.NamespacedName.Namespace, globalNetworkSetList)
			for _, globalNetworkSet := range deleteGlobalNetworkSetList.Items {
				if globalNetworkSet.GetName() != "" {
					controllerGlobalNetworksetsLog.Info("Remove globalnetworkset", "name", globalNetworkSet.GetName())
					err = r.Delete(
						ctx,
						&globalNetworkSet,
					)
					if err != nil {
						controllerGlobalNetworksetsLog.Error(err, "cannot delete GlobalNetworkSet")
						monitoring.NetworksetControllerGlobalNetworksetDeletionFailed.Inc()
						return ctrl.Result{}, err
					}
					monitoring.NetworksetControllerGlobalNetworksetDeleted.Inc()
				}
			}
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	var domain string
	for ruleNumber, rule := range instance.Spec.Egress {
		matches := destinationSelector.FindStringSubmatch(rule.Destination.Selector)
		if len(matches) > 1 {
			domain = matches[1]
			controllerGlobalNetworksetsLog.Info("Found domain", "request", req.NamespacedName, "domain", domain)
			ipAddress, err := resolveDomain(domain)
			if err != nil {
				return ctrl.Result{}, err
			}

			networkSet := getGlobalNetworkSet(req.NamespacedName.Name, req.NamespacedName.Namespace, domain, globalNetworkSetList)
			if networkSet.GetName() != "" {
				controllerGlobalNetworksetsLog.Info("Update existing networkset", "request", req.NamespacedName, "name", fmt.Sprint(req.NamespacedName.Name, "-", ruleNumber))
				err = r.Update(
					ctx,
					updateGlobalNetworkset(instance, networkSet, domain, ipAddress),
				)
				if err != nil {
					controllerGlobalNetworksetsLog.Error(err, "cannot update NetworkSet", "name", fmt.Sprint(req.NamespacedName.Name, "-", ruleNumber))
					monitoring.NetworksetControllerGlobalNetworksetUpdateFailed.Inc()
					return ctrl.Result{}, err
				}
				monitoring.NetworksetControllerGlobalNetworksetUpdated.Inc()
			} else {
				controllerGlobalNetworksetsLog.Info("Create networkset", "request", req.NamespacedName, "name", fmt.Sprint(req.NamespacedName.Name, "-", ruleNumber))
				err = r.Create(
					ctx,
					createGlobalNetworkset(instance, ruleNumber, domain, ipAddress),
				)
				if err != nil {
					controllerGlobalNetworksetsLog.Error(err, "cannot create NetworkSet", "name", fmt.Sprint(req.NamespacedName.Name, "-", ruleNumber))
					monitoring.NetworksetControllerGlobalNetworksetCreationFailed.Inc()
					return ctrl.Result{}, err
				}
				monitoring.NetworksetControllerGlobalNetworksetCreated.Inc()
			}
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GlobalNetworkPolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&calicov3.GlobalNetworkPolicy{}).
		WithOptions(controller.Options{MaxConcurrentReconciles: 2}).
		Complete(r)
}
