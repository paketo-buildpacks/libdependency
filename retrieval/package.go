// Package retrieval assists buildpack authors to create the retrieve step as outlined in the [Dependency RFCs],
// especially `Dependency Management Phase 2: Workflow and Github Action Generalization RFC`.
//
// Buildpacks should call retrieval.NewMetadata from the `main` func of their retrieval step.
// See the godoc for that function for additional information.
//
// [Dependency RFCs]: https://github.com/paketo-buildpacks/rfcs/tree/main/text/dependencies/rfcs
package retrieval
