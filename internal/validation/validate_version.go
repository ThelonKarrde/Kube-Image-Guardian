package validation

import (
	"github.com/ThelonKarrde/Kube-Image-Guardian/internal/rules"
	v1 "k8s.io/api/core/v1"
	"strings"
)

func (v *Validator) versionValidation(containers []v1.Container) (bool, error) {
	var result bool = false
	for _, c := range containers {
		b, err := rules.IsUsingLatestTag(c.Image)
		if err != nil {
			return false, err
		}
		if !b {
			imgT := strings.Split(c.Image, ":")
			if val, ok := v.allowConfig.DesiredVersions[imgT[0]]; ok {
				result, err := rules.MinimalVersion(val, imgT[1])
				if err != nil {
					return false, err
				}
				if result {
					return false, nil
				}
			}
		}
	}
	return result, nil
}
