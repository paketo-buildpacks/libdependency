package versionology

import (
	"fmt"
	"sort"

	"github.com/joshuatcasey/collections"
)

// VersionFetcherToString translates from an array of VersionFetcher to an array of strings.
// Primarily intended as a test helper.
func VersionFetcherToString(semverVersions []VersionFetcher) []string {
	return collections.TransformFunc(semverVersions, func(version VersionFetcher) string {
		return version.Version().String()
	})
}

// ConstraintsToString translates from an array of Constraints to an array of strings
// Primarily intended as a test helper.
func ConstraintsToString(semverVersions []Constraint) []string {
	return collections.TransformFunc(semverVersions, func(c Constraint) string {
		return c.Constraint.String()
	})
}

// LogAllVersions will print out a JSON array of the versions arranged as a block table.
// See Example tests for demonstration.
func LogAllVersions(id, description string, versions []VersionFetcher) {
	fmtString := "Found %d versions of %s %s\n"
	if len(versions) == 1 {
		fmtString = "Found %d version of %s %s\n"
	}
	fmt.Printf(fmtString, len(versions), id, description)

	sort.Slice(versions, func(i, j int) bool {
		return versions[i].Version().GreaterThan(versions[j].Version())
	})

	fmt.Printf("[\n  ")
	strings := VersionFetcherToString(versions)

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

// FilterUpstreamVersionsByConstraints will return only those versions with the following properties:
// - contained in upstreamVersions
// - satisfy at least one constraint
// - newer than all existing dependencies
func FilterUpstreamVersionsByConstraints(
	id string,
	upstreamVersions VersionFetcherArray,
	constraints []Constraint,
	existingVersion VersionFetcherArray) VersionFetcherArray {

	constraintsToDependencies := make(map[Constraint]VersionFetcherArray)

	for _, dependency := range existingVersion {
		for _, constraint := range constraints {
			if constraint.Check(dependency) {
				constraintsToDependencies[constraint] = append(constraintsToDependencies[constraint], dependency)
			}
		}
	}

	constraintsToInputVersion := make(map[Constraint][]VersionFetcher)

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

	constraintsToOutputVersions := make(map[Constraint][]VersionFetcher)

	for constraint, upstreamVersionsForConstraint := range constraintsToInputVersion {
		existingDependencies := constraintsToDependencies[constraint]

	ConstraintsToInputVersionLoop:
		for _, upstreamVersionForConstraint := range upstreamVersionsForConstraint {
			for _, existingDependency := range existingDependencies {
				if upstreamVersionForConstraint.Version().LessThan(existingDependency.Version()) || upstreamVersionForConstraint.Version().Equal(existingDependency.Version()) {
					continue ConstraintsToInputVersionLoop
				}
			}
			constraintsToOutputVersions[constraint] = append(constraintsToOutputVersions[constraint], upstreamVersionForConstraint)
		}
	}

	var outputVersions []VersionFetcher

	for constraint, constraintsToOutputVersion := range constraintsToOutputVersions {
		sort.Slice(constraintsToOutputVersion, func(i, j int) bool {
			return constraintsToOutputVersion[i].Version().LessThan(constraintsToOutputVersion[j].Version())
		})

		if constraint.Patches < len(constraintsToOutputVersion) {
			constraintsToOutputVersion = constraintsToOutputVersion[len(constraintsToOutputVersion)-constraint.Patches:]
		}

		constraintDescription := fmt.Sprintf("newer than '%s' for constraint %s, after limiting for %d patches",
			constraintsToDependencies[constraint].GetNewestVersion(),
			constraint.Constraint.String(),
			constraint.Patches)
		LogAllVersions(id, constraintDescription, constraintsToOutputVersion)

		outputVersions = append(outputVersions, constraintsToOutputVersion...)
	}

	if len(constraints) < 1 {
	ZeroConstraintsLoop:
		for _, upstreamVersion := range upstreamVersions {
			for _, dependency := range existingVersion {
				if upstreamVersion.Version().LessThan(dependency.Version()) || upstreamVersion.Version().Equal(dependency.Version()) {
					continue ZeroConstraintsLoop
				}
			}
			outputVersions = append(outputVersions, upstreamVersion)
		}
	}

	LogAllVersions(id, "as new versions", outputVersions)
	return outputVersions
}
