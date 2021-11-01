package rules

import (
	"log"
	"strings"

	"github.com/containers/image/docker/reference"
)

func ifIn(lStr []string, str string) bool {
	for _, s := range lStr {
		if s == str {
			return true
		}
	}
	return false
}

func FromAllowedRegistries(registries []string, image string) (bool, error) {
	normalized, err := reference.ParseNormalizedNamed(image)
	if err != nil {
		log.Printf("Image normalizing error: %s\n", err)
		return false, err
	}
	registry := strings.Split(normalized.Name(), "/")
	if ifIn(registries, registry[0]) {
		return true, nil
	}
	return false, nil
}
