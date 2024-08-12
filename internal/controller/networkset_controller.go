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
	"sort"
	"time"

	"github.com/go-logr/logr"
	"github.com/javdet/networksets-controller/monitoring"
	calicov3 "github.com/projectcalico/api/pkg/apis/projectcalico/v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// NetworkPolicyReconciler reconciles a Networkset object
type NetworkSetReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

var controllerNetworksetLog = ctrl.Log.WithName("controller").WithName("Networksets")

func (r *NetworkSetReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	controllerNetworksetLog.Info("start reconcile", "request", req.NamespacedName)

	networkSetList := &calicov3.NetworkSetList{}
	opts := []client.ListOption{
		client.MatchingLabels{
			"control-plane": "networksets-operator",
		},
	}
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctrl.Result{}, nil
		case <-ticker.C:
			// Resolve the domain names and update the NetworkSet
			err := r.List(ctx, networkSetList, opts...)
			if err != nil {
				controllerNetworksetLog.Error(err, "cannot get object NetworkSet list")
				return ctrl.Result{}, err
			}

			for _, networkSet := range networkSetList.Items {
				domain := networkSet.GetLabels()[selectorName]
				newIpAddress, err := resolveDomain(domain)
				if err != nil {
					return ctrl.Result{}, err
				}
				oldIpAddress := networkSet.Spec.Nets
				match, err := arraysMatch(newIpAddress, oldIpAddress)
				if err != nil {
					return ctrl.Result{}, err
				}

				if !match {
					controllerNetworksetLog.Info("Update dns networkset", "Networkset", networkSet.GetName())
					err = r.Update(
						ctx,
						updateNetworkset(&calicov3.NetworkPolicy{
							ObjectMeta: metav1.ObjectMeta{
								Name: networkSet.GetLabels()["parent-networkPolicy"],
							},
						},
							&networkSet,
							domain,
							newIpAddress),
					)
					if err != nil {
						controllerNetworksetLog.Error(err, "cannot update NetworkSet", "name", fmt.Sprint(req.NamespacedName.Name, "-", transformDomain(domain)))
						monitoring.NetworksetControllerNetworksetUpdateFailed.Inc()
						return ctrl.Result{}, err
					}
					monitoring.NetworksetControllerNetworksetUpdated.Inc()
				}
			}
		}
	}
}

func (r *NetworkSetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&calicov3.NetworkSet{}).
		WithOptions(controller.Options{MaxConcurrentReconciles: 1}).
		Complete(r)
}

func subnetsMatch(subnet1, subnet2 string) (bool, error) {
	_, ipNet1, err1 := net.ParseCIDR(subnet1)
	_, ipNet2, err2 := net.ParseCIDR(subnet2)

	if err1 != nil || err2 != nil {
		return false, fmt.Errorf("invalid CIDR notation")
	}

	return ipNet1.String() == ipNet2.String(), nil
}

func arraysMatch(array1, array2 []string) (bool, error) {
	if len(array1) != len(array2) {
		return false, nil
	}

	// Sort the arrays
	sort.Strings(array1)
	sort.Strings(array2)

	// Compare sorted arrays
	for i := range array1 {
		match, err := subnetsMatch(array1[i], array2[i])
		if err != nil {
			return false, err
		}
		if !match {
			return false, nil
		}
	}

	return true, nil
}
