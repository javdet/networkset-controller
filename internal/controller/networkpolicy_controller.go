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

	"net"
	"regexp"

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
type NetworkPolicyReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

var destinationSelector = regexp.MustCompile(`^DNS_RESOLVER\s==\s'(?P<domain>.*)'$`)

var controllerNetworksetsLog = ctrl.Log.WithName("controller").WithName("Networkpolicy")

//+kubebuilder:rbac:groups=crd.projectcalico.org,resources=networkpolicies,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=crd.projectcalico.org,resources=networkpolicies/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=crd.projectcalico.org,resources=networkpolicies/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Memcached object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *NetworkPolicyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	controllerNetworksetsLog.Info("start reconcile", "request", req.NamespacedName)
	instance := &calicov3.NetworkPolicy{}
	networkSetList := &calicov3.NetworkSetList{}

	opts := []client.ListOption{
		client.InNamespace(req.NamespacedName.Namespace),
		client.MatchingLabels{
			"control-plane":        "networksets-operator",
			"parent-networkPolicy": req.NamespacedName.Name,
		},
	}
	err := r.List(ctx, networkSetList, opts...)
	if err != nil {
		controllerNetworksetsLog.Error(err, "cannot get object NetworkSet list")
		return ctrl.Result{}, err
	}

	err = r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		controllerNetworksetsLog.Error(err, "cannot get object NetworkPolicy")
		if apierrors.IsNotFound(err) {
			deleteNetworkSetList := r.getNetworkSetList(req.NamespacedName.Name, req.NamespacedName.Namespace, networkSetList)
			for _, networkSet := range deleteNetworkSetList.Items {
				if networkSet.GetName() != "" {
					controllerNetworksetsLog.Info("Remove networkset", "name", networkSet.GetName())
					err = r.Delete(
						ctx,
						&networkSet,
					)
					if err != nil {
						controllerNetworksetsLog.Error(err, "cannot delete NetworkSet")
						return ctrl.Result{}, err
					}
					monitoring.NetworksetControllerNetworksetDeleted.Inc()
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
			controllerNetworksetsLog.Info("Found domain", "request", req.NamespacedName, "domain", domain)
			ipAddress, err := resolveDomain(domain)
			if err != nil {
				return ctrl.Result{}, err
			}

			networkSet := r.getNetworkSet(req.NamespacedName.Name, req.NamespacedName.Namespace, domain, networkSetList)
			if networkSet.GetName() != "" {
				controllerNetworksetsLog.Info("Update existing networkset", "request", req.NamespacedName, "name", fmt.Sprint(req.NamespacedName.Name, "-", ruleNumber))
				err = r.Update(
					ctx,
					updateNetworkset(instance, networkSet, domain, ipAddress),
				)
				if err != nil {
					controllerNetworksetsLog.Error(err, "cannot update NetworkSet", "name", fmt.Sprint(req.NamespacedName.Name, "-", transformDomain(domain)))
					monitoring.NetworksetControllerNetworksetUpdateFailed.Inc()
					return ctrl.Result{}, err
				}
				monitoring.NetworksetControllerNetworksetUpdated.Inc()
			} else {
				controllerNetworksetsLog.Info("Create networkset", "request", req.NamespacedName, "name", fmt.Sprint(req.NamespacedName.Name, "-", transformDomain(domain)))
				err = r.Create(
					ctx,
					createNetworkset(instance, domain, ipAddress),
				)
				if err != nil {
					controllerNetworksetsLog.Error(err, "cannot create NetworkSet", "name", fmt.Sprint(req.NamespacedName.Name, "-", transformDomain(domain)))
					monitoring.NetworksetControllerNetworksetCreationFailed.Inc()
					return ctrl.Result{}, err
				}
				monitoring.NetworksetControllerNetworksetCreated.Inc()
			}
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NetworkPolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&calicov3.NetworkPolicy{}).
		WithOptions(controller.Options{MaxConcurrentReconciles: 2}).
		Complete(r)
}

func resolveDomain(domain string) ([]string, error) {
	ips, err := net.LookupIP(domain)
	if err != nil {
		controllerNetworksetsLog.Error(err, "Error resolving domain")
		monitoring.NetworksetControllerResolveFailed.Inc()
		return []string{}, err
	}
	var ipStrings []string
	for _, ip := range ips {
		ipStrings = append(ipStrings, fmt.Sprint(ip.String(), "/32"))
	}
	monitoring.NetworksetControllerResolveSuccesful.Inc()

	return ipStrings, nil
}
