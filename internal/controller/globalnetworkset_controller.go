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
type GlobalNetworkSetReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

var controllerGlobalNetworksetLog = ctrl.Log.WithName("controller").WithName("GlobalNetworksets")

func (r *GlobalNetworkSetReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	controllerNetworksetLog.Info("start reconcile", "request", req.NamespacedName)

	globalNetworkSetList := &calicov3.GlobalNetworkSetList{}
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
			err := r.List(ctx, globalNetworkSetList, opts...)
			if err != nil {
				controllerGlobalNetworksetLog.Error(err, "cannot get object NetworkSet list")
				return ctrl.Result{}, err
			}

			for _, globalNetworkSet := range globalNetworkSetList.Items {
				domain := globalNetworkSet.GetLabels()[selectorName]
				newIpAddress, err := resolveDomain(domain)
				if err != nil {
					return ctrl.Result{}, err
				}
				oldIpAddress := globalNetworkSet.Spec.Nets
				match, err := arraysMatch(newIpAddress, oldIpAddress)
				if err != nil {
					return ctrl.Result{}, err
				}

				if !match {
					controllerGlobalNetworksetLog.Info("Update dns networkset", "Networkset", globalNetworkSet.GetName())
					err = r.Update(
						ctx,
						updateGlobalNetworkset(&calicov3.GlobalNetworkPolicy{
							ObjectMeta: metav1.ObjectMeta{
								Name: globalNetworkSet.GetLabels()["parent-networkPolicy"],
							},
						},
							&globalNetworkSet,
							domain,
							newIpAddress),
					)
					if err != nil {
						controllerNetworksetLog.Error(err, "cannot update NetworkSet", "name", fmt.Sprint(req.NamespacedName.Name, "-", transformDomain(domain)))
						monitoring.NetworksetControllerGlobalNetworksetUpdateFailed.Inc()
						return ctrl.Result{}, err
					}
					monitoring.NetworksetControllerGlobalNetworksetUpdated.Inc()
				}
			}
		}
	}
}

func (r *GlobalNetworkSetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&calicov3.GlobalNetworkSet{}).
		WithOptions(controller.Options{MaxConcurrentReconciles: 1}).
		Complete(r)
}
