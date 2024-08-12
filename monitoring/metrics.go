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

package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

// MetricDescription is an exported struct that defines the metric description (Name, Help)
// as a new type named MetricDescription.
type MetricDescription struct {
	Name string
	Help string
	Type string
}

// metricsDescription is a map of string keys (metrics) to MetricDescription values (Name, Help).
var metricDescription = map[string]MetricDescription{
	"NetworksetControllerResolveFailed": {
		Name: "networkset_controller_resolve_failed",
		Help: "Total number of failed resolve attempts.",
		Type: "Counter",
	},
	"NetworksetControllerResolveSuccesful": {
		Name: "networkset_controller_resolve_succesful",
		Help: "Total number of succesful resolve attempts.",
		Type: "Counter",
	},
	"NetworksetControllerNetworksetCreated": {
		Name: "networkset_controller_networkset_created",
		Help: "Total number of successful created networksets.",
		Type: "Counter",
	},
	"NetworksetControllerNetworksetCreationFailed": {
		Name: "networkset_controller_networkset_failed",
		Help: "Total number of failed creating networksets.",
		Type: "Counter",
	},
	"NetworksetControllerGlobalNetworksetCreated": {
		Name: "networkset_controller_globalnetworkset_created",
		Help: "Total number of successful created globalnetworksets.",
		Type: "Counter",
	},
	"NetworksetControllerGlobalNetworksetFailed": {
		Name: "networkset_controller_globalnetworkset_failed",
		Help: "Total number of failed creating globalnetworksets.",
		Type: "Counter",
	},
	"NetworksetControllerNetworksetUpdated": {
		Name: "networkset_controller_networkset_updated",
		Help: "Total number of successful updated networksets.",
		Type: "Counter",
	},
	"NetworksetControllerNetworksetUpdateFailed": {
		Name: "networkset_controller_networkset_update_failed",
		Help: "Total number of failed updating networksets.",
		Type: "Counter",
	},
	"NetworksetControllerGlobalNetworksetUpdated": {
		Name: "networkset_controller_globalnetworkset_updated",
		Help: "Total number of successful updated globalnetworksets.",
		Type: "Counter",
	},
	"NetworksetControllerGlobalNetworksetUpdateFailed": {
		Name: "networkset_controller_globalnetworkset_update_failed",
		Help: "Total number of failed updating globalnetworksets.",
		Type: "Counter",
	},
	"NetworksetControllerNetworksetDeleted": {
		Name: "networkset_controller_networkset_deleted",
		Help: "Total number of successful deleted networksets.",
		Type: "Counter",
	},
	"NetworksetControllerNetworksetDeletionFailed": {
		Name: "networkset_controller_networkset_deletion_failed",
		Help: "Total number of failed deletion networksets.",
		Type: "Counter",
	},
	"NetworksetControllerGlobalNetworksetDeleted": {
		Name: "networkset_controller_globalnetworkset_deleted",
		Help: "Total number of successful deleted globalnetworksets.",
		Type: "Counter",
	},
	"NetworksetControllerGlobalNetworksetDeletionFailed": {
		Name: "networkset_controller_globalnetworkset_deletion_failed",
		Help: "Total number of failed deletion globalnetworksets.",
		Type: "Counter",
	},
}

var (
	// MemcachedDeploymentSizeUndesiredCountTotal will count how many times was required
	// to perform the operation to ensure that the number of replicas on the cluster
	// is the same as the quantity desired and specified via the custom resource size spec.
	NetworksetControllerResolveFailed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: metricDescription["NetworksetControllerResolveFailed"].Name,
			Help: metricDescription["NetworksetControllerResolveFailed"].Help,
		},
	)
	NetworksetControllerResolveSuccesful = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: metricDescription["NetworksetControllerResolveSuccesful"].Name,
			Help: metricDescription["NetworksetControllerResolveSucessful"].Help,
		},
	)
	NetworksetControllerNetworksetCreated = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: metricDescription["NetworksetControllerNetworksetCreated"].Name,
			Help: metricDescription["NetworksetControllerNetworksetCreated"].Help,
		},
	)
	NetworksetControllerNetworksetCreationFailed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: metricDescription["NetworksetControllerNetworksetCreationFailed"].Name,
			Help: metricDescription["NetworksetControllerNetworksetCreationFailed"].Help,
		},
	)
	NetworksetControllerGlobalNetworksetCreated = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: metricDescription["NetworksetControllerGlobalNetworksetCreated"].Name,
			Help: metricDescription["NetworksetControllerGlobalNetworksetCreated"].Help,
		},
	)
	NetworksetControllerGlobalNetworksetCreationFailed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: metricDescription["NetworksetControllerGlobalNetworksetCreationFailed"].Name,
			Help: metricDescription["NetworksetControllerGlobalNetworksetCreationFailed"].Help,
		},
	)
	NetworksetControllerNetworksetUpdated = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: metricDescription["NetworksetControllerNetworksetUpdated"].Name,
			Help: metricDescription["NetworksetControllerNetworksetUpdated"].Help,
		},
	)
	NetworksetControllerNetworksetUpdateFailed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: metricDescription["NetworksetControllerNetworksetUpdateFailed"].Name,
			Help: metricDescription["NetworksetControllerNetworksetUpdateFailed"].Help,
		},
	)
	NetworksetControllerGlobalNetworksetUpdated = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: metricDescription["NetworksetControllerGlobalNetworksetUpdated"].Name,
			Help: metricDescription["NetworksetControllerGlobalNetworksetUpdated"].Help,
		},
	)
	NetworksetControllerGlobalNetworksetUpdateFailed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: metricDescription["NetworksetControllerGlobalNetworksetUpdateFailed"].Name,
			Help: metricDescription["NetworksetControllerGlobalNetworksetUpdateFailed"].Help,
		},
	)
	NetworksetControllerNetworksetDeleted = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: metricDescription["NetworksetControllerNetworksetDeleted"].Name,
			Help: metricDescription["NetworksetControllerNetworksetDeleted"].Help,
		},
	)
	NetworksetControllerNetworksetDeletionFailed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: metricDescription["NetworksetControllerNetworksetDeletionFailed"].Name,
			Help: metricDescription["NetworksetControllerNetworksetDeletionFailed"].Help,
		},
	)
	NetworksetControllerGlobalNetworksetDeleted = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: metricDescription["NetworksetControllerGlobalNetworksetDeleted"].Name,
			Help: metricDescription["NetworksetControllerGlobalNetworksetDeleted"].Help,
		},
	)
	NetworksetControllerGlobalNetworksetDeletionFailed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: metricDescription["NetworksetControllerGlobalNetworksetDeletionFailed"].Name,
			Help: metricDescription["NetworksetControllerGlobalNetworksetDeletionFailed"].Help,
		},
	)
)

// RegisterMetrics will register metrics with the global prometheus registry
func RegisterMetrics() {
	metrics.Registry.MustRegister(NetworksetControllerResolveFailed)
	metrics.Registry.MustRegister(NetworksetControllerResolveSuccesful)
	metrics.Registry.MustRegister(NetworksetControllerNetworksetCreated)
	metrics.Registry.MustRegister(NetworksetControllerNetworksetCreationFailed)
	metrics.Registry.MustRegister(NetworksetControllerGlobalNetworksetCreated)
	//metrics.Registry.MustRegister(NetworksetControllerGlobalNetworksetCreationFailed)
	metrics.Registry.MustRegister(NetworksetControllerNetworksetUpdated)
	metrics.Registry.MustRegister(NetworksetControllerNetworksetUpdateFailed)
	metrics.Registry.MustRegister(NetworksetControllerGlobalNetworksetUpdated)
	metrics.Registry.MustRegister(NetworksetControllerGlobalNetworksetUpdateFailed)
	metrics.Registry.MustRegister(NetworksetControllerNetworksetDeleted)
	metrics.Registry.MustRegister(NetworksetControllerNetworksetDeletionFailed)
	metrics.Registry.MustRegister(NetworksetControllerGlobalNetworksetDeleted)
	metrics.Registry.MustRegister(NetworksetControllerGlobalNetworksetDeletionFailed)
}

// ListMetrics will create a slice with the metrics available in metricDescription
func ListMetrics() []MetricDescription {
	v := make([]MetricDescription, 0, len(metricDescription))
	// Insert value (Name, Help) for each metric
	for _, value := range metricDescription {
		v = append(v, value)
	}

	return v
}
