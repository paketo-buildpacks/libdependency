package versionology

import (
	"fmt"
	"sort"

	"github.com/Masterminds/semver/v3"
	"github.com/joshuatcasey/collections"
)

func HasVersionToString(semverVersions []HasVersion) []string {
	return collections.TransformFunc(semverVersions, func(version HasVersion) string {
		return version.GetVersion().String()
	})
}

func SemverToString(semverVersions []*semver.Version) []string {
	return collections.TransformFunc(semverVersions, func(v *semver.Version) string {
		return v.String()
	})
}

func ConstraintsToString(semverVersions []Constraint) []string {
	return collections.TransformFunc(semverVersions, func(c Constraint) string {
		return c.Constraint.String()
	})
}

func LogAllVersions(id, description string, versions []HasVersion) {
	fmt.Printf("Found %d versions of %s %s\n", len(versions), id, description)

	sort.Slice(versions, func(i, j int) bool {
		return versions[i].GetVersion().GreaterThan(versions[j].GetVersion())
	})

	fmt.Printf("[\n  ")
	strings := HasVersionToString(versions)

	maxWidth := make([]int, 5)
	for i, s := range strings {
		length := len(s)
		if length > maxWidth[i%5] {
			maxWidth[i%5] = length
		}
	}

	for i, s := range strings {
		fmt.Printf(`"%s"`, s)

		if i != len(strings)-1 {
			fmt.Print(",")
		}

		if i != len(strings)-1 {
			if i > 0 && (i+1)%5 == 0 {
				fmt.Printf("\n  ")
			} else {
				fmt.Printf("%*s", 1+maxWidth[i%5]-len(s), "")
			}
		}
	}
	fmt.Printf("\n]\n")
}

func FilterVersionsByConstraints(inputVersions []*semver.Version, constraints []*semver.Constraints) []*semver.Version {
	if len(constraints) == 0 {
		return inputVersions
	}

	var outputVersions []*semver.Version

	for _, version := range inputVersions {
		for _, constraint := range constraints {
			if constraint.Check(version) {
				outputVersions = append(outputVersions, version)
			}
		}
	}

	return outputVersions
}

func FilterUpstreamVersionsByConstraints(id string, upstreamVersions HasVersionArray, constraints []Constraint, dependencies HasVersionArray) HasVersionArray {
	constraintsToDependencies := make(map[Constraint][]HasVersion)

	for _, dependency := range dependencies {
		for _, constraint := range constraints {
			if constraint.Check(dependency) {
				constraintsToDependencies[constraint] = append(constraintsToDependencies[constraint], dependency)
			}
		}
	}

	constraintsToInputVersion := make(map[Constraint][]HasVersion)

	for _, version := range upstreamVersions {
		for _, constraint := range constraints {
			if constraint.Check(version) {
				constraintsToInputVersion[constraint] = append(constraintsToInputVersion[constraint], version)
			}
		}
	}

	for constraint, versions := range constraintsToInputVersion {
		constraintDescription := fmt.Sprintf("for constraint %s", constraint.Constraint.String())
		LogAllVersions(id, constraintDescription, versions)
	}

	constraintsToOutputVersions := make(map[Constraint][]HasVersion)

	for constraint, upstreamVersionsForConstraint := range constraintsToInputVersion {
		existingDependencies := constraintsToDependencies[constraint]

	ConstraintsToInputVersionLoop:
		for _, upstreamVersionForConstraint := range upstreamVersionsForConstraint {
			for _, existingDependency := range existingDependencies {
				if upstreamVersionForConstraint.GetVersion().LessThan(existingDependency.GetVersion()) || upstreamVersionForConstraint.GetVersion().Equal(existingDependency.GetVersion()) {
					continue ConstraintsToInputVersionLoop
				}
			}
			constraintsToOutputVersions[constraint] = append(constraintsToOutputVersions[constraint], upstreamVersionForConstraint)
		}
	}

	var outputVersions []HasVersion

	for constraint, constraintsToOutputVersion := range constraintsToOutputVersions {
		sort.Slice(constraintsToOutputVersion, func(i, j int) bool {
			return constraintsToOutputVersion[i].GetVersion().LessThan(constraintsToOutputVersion[j].GetVersion())
		})

		if constraint.Patches < len(constraintsToOutputVersion) {
			constraintsToOutputVersion = constraintsToOutputVersion[len(constraintsToOutputVersion)-constraint.Patches:]
		}

		constraintDescription := fmt.Sprintf("for constraint %s, after limiting for %d patches", constraint.Constraint.String(), constraint.Patches)
		LogAllVersions(id, constraintDescription, constraintsToOutputVersion)

		outputVersions = append(outputVersions, constraintsToOutputVersion...)
	}

	if len(constraints) < 1 {
	ZeroConstraintsLoop:
		for _, upstreamVersion := range upstreamVersions {
			for _, dependency := range dependencies {
				if upstreamVersion.GetVersion().LessThan(dependency.GetVersion()) || upstreamVersion.GetVersion().Equal(dependency.GetVersion()) {
					continue ZeroConstraintsLoop
				}
			}
			outputVersions = append(outputVersions, upstreamVersion)
		}
	}

	LogAllVersions(id, "as new versions", outputVersions)
	return outputVersions
}
