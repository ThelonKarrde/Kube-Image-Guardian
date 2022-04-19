package rules

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
)

func LimitsDefined(enforced bool, containers []v1.Container) (bool, []string) {
	var result bool = true
	var clList []string
	for _, c := range containers {
		if c.Resources.Limits == nil && enforced {
			result = false
			clList = append(clList, fmt.Sprintf("Limits are undefined for container: %s", c.Name))
		}
	}
	return result, clList
}
