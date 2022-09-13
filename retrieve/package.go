// Package retrieve assists buildpack authors to create the retrieve step as outlined in the [Dependency RFCs],
// especially `Dependency Management Phase 2: Workflow and Github Action Generalization RFC`.
//
// Buildpacks should call retrieve.NewMetadata from the `main` func of their retrieval step and pass in the necessary
// function implementations:
// - libdependency.VersionFetcherFunc: to retrieve all versions of a dependency from the internet
// - retrieve.GenerateMetadataFunc: to generate metadata for a specific version
//
// [Dependency RFCs]: https://github.com/paketo-buildpacks/rfcs/tree/main/text/dependencies/rfcs
package retrieve
