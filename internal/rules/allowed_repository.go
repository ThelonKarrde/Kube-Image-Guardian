package rules

import (
	"log"
	"strings"

	"github.com/containers/image/docker/reference"
)

func FromAllowedRepository(repositories []string, image string) (bool, error) {
	normalized, err := reference.ParseNormalizedNamed(image)
	if err != nil {
		log.Printf("Image normalizing error: %s\n", err)
		return false, err
	}
	repository := normalized.Name()
	i := strings.Index(repository, ":")
	if i > 0 {
		repository = repository[0:i]
	}
	log.Print(repositories)
	log.Print(repository)
	if ifIn(repositories, repository) {
		return true, nil
	}
	return false, nil
}
