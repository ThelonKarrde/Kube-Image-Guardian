package rules_test

import (
	"testing"

	"github.com/ThelonKarrde/Kube-Image-Guardian/internal/rules"
)

func TestFromAllowedRegistries(t *testing.T) {
	registries := []string{"docker.io", "gcr.io"}
	r, _ := rules.FromAllowedRegistries(registries, "hub.docker.com/r/nginx")
	if r == true {
		t.Errorf("Incorrect result: %t, expected %t", r, false)
	}
	r, _ = rules.FromAllowedRegistries(registries, "nginx")
	if r == false {
		t.Errorf("Incorrect result: %t, expected %t", r, true)
	}
}
