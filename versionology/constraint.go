package versionology

import (
	"github.com/Masterminds/semver/v3"
	"github.com/paketo-buildpacks/packit/v2/cargo"
)

// Constraint largely mimics cargo.ConfigMetadataDependencyConstraint but has
// a semver.Constraints instead of a string
type Constraint struct {
	Constraint *semver.Constraints
	ID         string
	Patches    int
}

// NewConstraint will translate a cargo.ConfigMetadataDependencyConstraint into a Constraint
func NewConstraint(c cargo.ConfigMetadataDependencyConstraint) (Constraint, error) {
	semverConstraint, err := semver.NewConstraint(c.Constraint)

	if err != nil {
		return Constraint{}, err
	}

	return Constraint{
		ID:         c.ID,
		Patches:    c.Patches,
		Constraint: semverConstraint,
	}, nil
}

// Check tests if a version satisfies the constraints.
func (c Constraint) Check(versionFetcher VersionFetcher) bool {
	return c.Constraint.Check(versionFetcher.Version())
}
