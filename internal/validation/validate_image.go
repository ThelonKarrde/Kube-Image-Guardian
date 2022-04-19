package validation

import (
	"fmt"

	"github.com/ThelonKarrde/Kube-Image-Guardian/internal/rules"
	v1 "k8s.io/api/core/v1"
)

func (v *Validator) imageValidation(containers []v1.Container) (string, error) {
	for _, container := range containers {
		isLatest, err := rules.IsUsingLatestTag(container.Image)
		if err != nil {
			v.errorLog.Printf("Error normalizing image name: %s", container.Image)
			return "processing image name", err
		}
		if !v.allowConfig.AllowLatest && isLatest {
			if v.allowConfig.LogOnly {
				v.infoLog.Println(fmt.Sprintf("[LogOnly mode] Check for container %s; image %s is not allowed because latest tag prohibited.", container.Name, container.Image))
			} else {
				return "Latest tag is not allowed", nil
			}
		}
		rep, err := rules.FromAllowedRepository(v.allowConfig.AllowedRepositories, container.Image)
		if err != nil {
			v.errorLog.Printf("Error normalizing image name: %s", container.Image)
			return "processing image name", err
		}
		if !rep {
			reg, err := rules.FromAllowedRegistries(v.allowConfig.AllowedRegistries, container.Image)
			if err != nil {
				v.errorLog.Printf("Error normalizing image name: %s", container.Image)
				return "processing image name", err
			}
			if !reg {
				if v.allowConfig.LogOnly {
					v.infoLog.Println(fmt.Sprintf("[LogOnly mode] Check for container %s; image %s is not allowed because image from untrusted registry.", container.Name, container.Image))
				} else {
					return "Image is from untrusted registry", nil
				}
			}
		}
	}
	return "", nil
}
