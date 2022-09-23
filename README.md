# Libdependency

Package libdependency provides a library for buildpack authors to construct code to retrieve new versions and their associated metadata.

It also contains various helpers and actions to assist the overall dependency management workflow.

## Retrieval

The `retrieve` subpackage has an entrypoint func called `NewMetadata` that takes in a buildpack id.
See the `godoc` for that package for additional information.