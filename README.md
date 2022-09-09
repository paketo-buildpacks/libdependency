# Libdependency

Package libdependency provides a library for buildpack authors to construct code to retrieve new versions and their associated metadata.

It also contains various helpers and actions to assist the overall dependency management workflow.

## Integration

### Retrieval

The retrieve subpackage has an entrypoint func called `NewMetadata` that takes in various information and functions.

```go
type HasVersionsFunc func() ([]versionology.HasVersion, error)
type GenerateMetadataFunc func(version versionology.HasVersion) (cargo.ConfigMetadataDependency, error)

func NewMetadata(id string, getNewVersions HasVersionsFunc, generateMetadata GenerateMetadataFunc, targets ...string)
```

The role of `HasVersionsFunc` is to return all known versions from an online source as an array of `versionology.HasVersion`.

Buildpacks authors can choose the source of these versions. Some examples include:

- `nginx` versions from https://github.com/nginx/nginx/tags
- `bundler` versions from https://rubygems.org/api/v1/versions/bundler.json

The role of `GenerateMetadataFunc` is to take in a single version and generate all the associated metadata for it.
That way `NewMetadata` can assemble the `metadata.json` file containing all new metadata for all new versions.