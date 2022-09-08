# Libdependency

This exists to assist buildpack authors construct code to retrieve new versions and their associated metadata.

It also contains various helpers and actions to assist the overall dependencies workflow.

## Integration

### Retrieval

The retrieve subpackage has an entrypoint func called `NewMetadata` that takes in various information and functions.

```go
type HasVersionsFunc func() ([]versionology.HasVersion, error)
type GenerateMetadataFunc func(version versionology.HasVersion) (cargo.ConfigMetadataDependency, error)

func NewMetadata(id string, getNewVersions HasVersionsFunc, generateMetadata GenerateMetadataFunc, targets ...string)
```

The role of `HasVersionsFunc` is to return all known versions from an online source.

- For `nginx` this could be obtained from https://github.com/nginx/nginx/tags
- For `bundler` this could be obtained from https://rubygems.org/api/v1/versions/bundler.json

The implementation of `HasVersionsFunc` returns an array of structs that implement `versionology.HasVersion`.
This allows buildpack authors to use whatever custom struct they need, since `libdependency` only needs the version.

The role of `GenerateMetadataFunc` is to take in a single version and generate all the associated metadata for it.
That way `NewMetadata` can assemble the `metadata.json` file containing all new metadata for all new versions.