package rules

import (
	"github.com/containers/image/docker/reference"
)

func FromAllowedRepository(repositories []string, image string) (bool, error) {
	normalized, err := reference.ParseNormalizedNamed(image)
	if err != nil {
		return false, err
	}
	repository := normalized.Name()
	if ifIn(repositories, repository) {
		return true, nil
	}
	return false, nil
}
