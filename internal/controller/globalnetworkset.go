package controller

import (
	"fmt"
	"strings"

	calicov3 "github.com/projectcalico/api/pkg/apis/projectcalico/v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func createGlobalNetworkset(instance *calicov3.GlobalNetworkPolicy, ruleNumber int, domain string, ipAddress []string) *calicov3.GlobalNetworkSet {
	return &calicov3.GlobalNetworkSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "GlobalNetworkSet",
			APIVersion: "projectcalico.org/v3",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        fmt.Sprint(instance.GetName(), "-", transformDomain(domain)),
			Labels:      getLabels(instance.GetName(), domain),
			Annotations: getAnnotations(),
		},
		Spec: calicov3.GlobalNetworkSetSpec{
			Nets: ipAddress,
		},
	}
}

func updateGlobalNetworkset(instance *calicov3.GlobalNetworkPolicy, globalNetworkSet *calicov3.GlobalNetworkSet, domain string, ipAddress []string) *calicov3.GlobalNetworkSet {
	globalNetworkSet.SetLabels(getLabels(instance.GetName(), domain))
	globalNetworkSet.Spec.Nets = ipAddress
	return globalNetworkSet
}

func getGlobalNetworkSet(policyName string, PolicyNamespace string, domain string, globalNetworkSetList *calicov3.GlobalNetworkSetList) *calicov3.GlobalNetworkSet {
	result := &calicov3.GlobalNetworkSet{}
	for _, gloablNetworkSet := range globalNetworkSetList.Items {
		if strings.Contains(gloablNetworkSet.GetName(), policyName) {
			if gloablNetworkSet.GetLabels()[selectorName] == domain {
				result = &gloablNetworkSet
				break
			}
		}
	}

	return result
}

func getGlobalNetworkSetList(policyName string, PolicyNamespace string, globalNetworkSetList *calicov3.GlobalNetworkSetList) *calicov3.GlobalNetworkSetList {
	result := &calicov3.GlobalNetworkSetList{}
	for _, gloablNetworkSet := range globalNetworkSetList.Items {
		if strings.Contains(gloablNetworkSet.GetName(), policyName) {
			result.Items = append(result.Items, gloablNetworkSet)
		}
	}

	return result
}
