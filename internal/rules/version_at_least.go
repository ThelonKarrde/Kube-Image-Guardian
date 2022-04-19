package rules

import (
	"github.com/Masterminds/semver"
)

func MinimalVersion(constraint string, actual string) (bool, error) {
	var err error
	c, err := semver.NewConstraint(constraint)
	if err != nil {
		return false, err
	}
	v, err := semver.NewVersion(actual)
	if err != nil {
		return false, err
	}
	if c.Check(v) {
		return true, nil
	}
	return false, nil
}
