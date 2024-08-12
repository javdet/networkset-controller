package controller

import (
	"fmt"
	"strings"

	calicov3 "github.com/projectcalico/api/pkg/apis/projectcalico/v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const selectorName = "DNS_RESOLVER"

func createNetworkset(instance *calicov3.NetworkPolicy, domain string, ipAddress []string) *calicov3.NetworkSet {
	return &calicov3.NetworkSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "NetworkSet",
			APIVersion: "projectcalico.org/v3",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        fmt.Sprint(instance.GetName(), "-", transformDomain(domain)),
			Namespace:   instance.GetNamespace(),
			Labels:      getLabels(instance.GetName(), domain),
			Annotations: getAnnotations(),
		},
		Spec: calicov3.NetworkSetSpec{
			Nets: ipAddress,
		},
	}
}

func updateNetworkset(instance *calicov3.NetworkPolicy, networkSet *calicov3.NetworkSet, domain string, ipAddress []string) *calicov3.NetworkSet {
	networkSet.SetLabels(getLabels(instance.GetName(), domain))
	networkSet.Spec.Nets = ipAddress
	return networkSet
}

// getLabels get common labels
func getLabels(name string, domain string) map[string]string {
	return map[string]string{
		selectorName:           domain,
		"parent-networkPolicy": name,
		"control-plane":        "networksets-operator",
	}
}

// getAnnotations get common annotations
func getAnnotations() map[string]string {
	return map[string]string{
		"operator":      "networksets",
		"control-plane": "networksets-operator",
	}
}

func (r *NetworkPolicyReconciler) getNetworkSet(policyName string, PolicyNamespace string, domain string, networkSetList *calicov3.NetworkSetList) *calicov3.NetworkSet {
	result := &calicov3.NetworkSet{}
	for _, networkSet := range networkSetList.Items {
		if strings.Contains(networkSet.GetName(), policyName) {
			if networkSet.GetLabels()[selectorName] == domain {
				result = &networkSet
				break
			}
		}
	}

	return result
}

func (r *NetworkPolicyReconciler) getNetworkSetList(policyName string, PolicyNamespace string, networkSetList *calicov3.NetworkSetList) *calicov3.NetworkSetList {
	result := &calicov3.NetworkSetList{}
	for _, networkSet := range networkSetList.Items {
		if strings.Contains(networkSet.GetName(), policyName) {
			result.Items = append(result.Items, networkSet)
		}
	}

	return result
}

func transformDomain(domain string) string {
	domain = strings.ReplaceAll(domain, ".", "-")
	domain = strings.ReplaceAll(domain, ":", "-")
	return domain
}
