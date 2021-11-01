package rules

import (
	"log"
	"strings"

	"github.com/containers/image/docker/reference"
)

func IsUsingLatestTag(imageName string) (bool, error) {
	normalized, err := reference.ParseNormalizedNamed(imageName)
	if err != nil {
		log.Printf("Image normalizing error: %s\n", err)
		return false, err
	}
	log.Printf("Normalized name: %s\n", normalized)
	if !strings.HasSuffix(normalized.Name(), ":latest") {
		return false, nil
	}
	return true, nil
}
