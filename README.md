# libdependency

Library for dependency-management related functionality

Package libdependency provides a library for buildpack authors to construct code to retrieve new versions and their associated metadata.

It also contains various helpers and actions to assist the overall dependency management workflow.

See: [Buildpack Dependency Management Improvement Overview RFC](https://github.com/paketo-buildpacks/rfcs/blob/main/text/dependencies/rfcs/0003-dependency-management-overview.md)

## Retrieval

The `retrieve` subpackage has an entrypoint func called `NewMetadata` that takes in a buildpack id.

In addition the `retrieve` subpackage has an entrypoint func called `NewMetadataWithPlatform` which supports multi-arch dependency updates that takes in a buildpack id.

See the `godoc` for that package for additional information.

## Tests

Run the unit tests with:

```
./scripts/unit.sh
```

The integration tests in `integration/` call the live GitHub API. They are behind the `integration` build tag, so `go test ./...` and `scripts/unit.sh` skip them. CI does not run them. Run them with:

```
go test -tags integration ./integration/... -run Integration
```

Set `GIT_TOKEN` to a GitHub token to avoid the unauthenticated rate limit of 60 requests per hour. One run uses about 8 requests.
