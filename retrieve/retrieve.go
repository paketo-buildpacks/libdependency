package retrieve

import (
	"flag"
	"fmt"
	"os"

	"github.com/joshuatcasey/collections"
	"github.com/joshuatcasey/libdependency"
	"github.com/joshuatcasey/libdependency/versionology"
	"github.com/joshuatcasey/libdependency/workflows"
	"github.com/paketo-buildpacks/packit/v2/cargo"
	"github.com/paketo-buildpacks/packit/v2/fs"
)

type GenerateMetadataFunc func(id, name string, version versionology.HasVersion) (cargo.ConfigMetadataDependency, error)

func NewMetadata(getNewVersions libdependency.HasVersionsFunc, generateMetadata GenerateMetadataFunc) {
	var (
		buildpackTomlPath      string
		outputFile             string
		id                     string
		name                   string
		buildpackTomlPathUsage = "full path to the buildpack.toml file, using only one of camelCase, snake_case, or dash_case"
		outputFileUsage        = "output filename into which to write the JSON metadata, using only one of camelCase, snake_case, or dash_case"
	)

	flag.StringVar(&buildpackTomlPath, "buildpackTomlPath", "", buildpackTomlPathUsage)
	flag.StringVar(&buildpackTomlPath, "buildpack_toml_path", buildpackTomlPath, buildpackTomlPathUsage)
	flag.StringVar(&buildpackTomlPath, "buildpack-toml-path", buildpackTomlPath, buildpackTomlPathUsage)
	flag.StringVar(&outputFile, "outputFile", "", outputFileUsage)
	flag.StringVar(&outputFile, "output_file", outputFile, outputFileUsage)
	flag.StringVar(&outputFile, "output-file", outputFile, outputFileUsage)
	flag.StringVar(&id, "id", "", "id of the dependency")
	flag.StringVar(&name, "name", "", "name of the dependency")
	flag.Parse()

	validate(buildpackTomlPath, outputFile, id, name)

	config, err := libdependency.ParseBuildpackToml(buildpackTomlPath)
	if err != nil {
		panic(err)
	}

	newVersions, err := libdependency.GetNewVersionsForId(id, config, getNewVersions)
	if err != nil {
		panic(err)
	}

	dependencies, err := collections.TransformFuncWithError(newVersions,
		func(hasVersion versionology.HasVersion) (cargo.ConfigMetadataDependency, error) {
			fmt.Printf("Generating metadata for %s\n", hasVersion.GetVersion().String())
			return generateMetadata(id, name, hasVersion)
		})
	if err != nil {
		panic(err)
	}

	json, err := workflows.ToWorkflowJson(dependencies)
	if err != nil {
		panic(fmt.Errorf("unable to marshall json, with error=%w", err))
	}

	if err = os.WriteFile(outputFile, []byte(json), os.ModePerm); err != nil {
		panic(fmt.Errorf("cannot write to %s: %w", outputFile, err))
	} else {
		fmt.Printf("Wrote output to %s\n", outputFile)
	}
}

func validate(buildpackTomlPath, outputFile, id, name string) {
	if exists, err := fs.Exists(buildpackTomlPath); err != nil {
		panic(err)
	} else if !exists {
		panic(fmt.Errorf("could not locate buildpack.toml at '%s'", buildpackTomlPath))
	}

	if outputFile == "" {
		panic("outputFile is required")
	}

	if id == "" {
		panic("id is required")
	}

	if name == "" {
		panic("name is required")
	}
}
