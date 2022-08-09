package metadata

import (
	"flag"
	"fmt"
	"os"

	"github.com/joshuatcasey/libdependency"
	"github.com/joshuatcasey/libdependency/workflows"
	"github.com/paketo-buildpacks/packit/v2/cargo"
)

type GetNewMetadataFunc func(id, name string, config cargo.Config) ([]cargo.ConfigMetadataDependency, error)

func GetNewMetadata(f GetNewMetadataFunc) {
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

	config, err := libdependency.ParseBuildpackToml(buildpackTomlPath)
	if err != nil {
		panic(err)
	}

	versions, err := f(id, name, config)

	json, err := workflows.ToWorkflowJson(versions)
	if err != nil {
		panic(fmt.Errorf("unable to marshall json, with error=%w", err))
	}

	err = os.WriteFile(outputFile, []byte(json), os.ModePerm)
	if err != nil {
		panic(fmt.Errorf("cannot write to %s: %w", outputFile, err))
	}
}
