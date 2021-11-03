package rules

import (
	"log"

	"github.com/containers/image/docker/reference"
)

func FromAllowedRepository(repositories []string, image string) (bool, error) {
	normalized, err := reference.ParseNormalizedNamed(image)
	if err != nil {
		log.Printf("Image normalizing error: %s\n", err)
		return false, err
	}
	repository := normalized.Name()
	if ifIn(repositories, repository) {
		return true, nil
	}
	return false, nil
}
