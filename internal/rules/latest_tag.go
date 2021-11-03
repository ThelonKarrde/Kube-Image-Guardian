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
	if !strings.Contains(normalized.String(), ":") {
		return true, nil
	}
	if !strings.HasSuffix(normalized.String(), ":latest") {
		return false, nil
	}
	return true, nil
}
