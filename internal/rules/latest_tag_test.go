package rules_test

import (
	"testing"

	"github.com/ThelonKarrde/Kube-Image-Guardian/internal/rules"
)

func TestIsUsingLatestTag(t *testing.T) {
	r, _ := rules.IsUsingLatestTag("nginx:latest")
	if r == false {
		t.Errorf("Incorrect result: %t, expected %t", r, true)
	}
}

func TestIsUsingLatestTagImplicitly(t *testing.T) {
	r, _ := rules.IsUsingLatestTag("nginx")
	if r == false {
		t.Errorf("Incorrect result: %t, expected %t", r, true)
	}
}

func TestIsUsingLatestTagNotLatest(t *testing.T) {
	r, _ := rules.IsUsingLatestTag("nginx:1.3.2")
	if r == true {
		t.Errorf("Incorrect result: %t, expected %t", r, true)
	}
}
