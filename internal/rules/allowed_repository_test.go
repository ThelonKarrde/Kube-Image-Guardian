package rules_test

import (
	"testing"

	"github.com/ThelonKarrde/Kube-Image-Guardian/internal/rules"
)

func TestFromAllowedRepositoryTrue(t *testing.T) {
	registries := []string{"docker.io/library/nginx"}
	r, _ := rules.FromAllowedRepository(registries, "nginx")
	if r == false {
		t.Errorf("Incorrect result: %t, expected %t", r, false)
	}
}

func TestFromAllowedRepositoryFalse(t *testing.T) {
	registries := []string{"docker.io/library/nginx"}
	r, _ := rules.FromAllowedRepository(registries, "alpine")
	if r == true {
		t.Errorf("Incorrect result: %t, expected %t", r, false)
	}
}
