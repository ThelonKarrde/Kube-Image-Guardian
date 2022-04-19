package validation

import (
	"github.com/ThelonKarrde/Kube-Image-Guardian/internal/rules"
	v1 "k8s.io/api/core/v1"
)

func (v *Validator) resourceValidation(containers []v1.Container) (bool, []string) {
	reqRes, reqL := rules.RequestsDefined(v.allowConfig.RequestsDefined, containers)
	limRes, limL := rules.LimitsDefined(v.allowConfig.LimitsDefined, containers)
	if reqRes && limRes {
		return true, nil
	}
	return false, append(reqL, limL...)
}
