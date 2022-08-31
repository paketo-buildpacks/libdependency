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
	"golang.org/x/exp/slices"
)

type GenerateMetadataFunc func(version versionology.HasVersion) (cargo.ConfigMetadataDependency, error)

func NewMetadata(id string, getNewVersions libdependency.HasVersionsFunc, generateMetadata GenerateMetadataFunc, targets ...string) {
	var (
		buildpackTomlPath      string
		metadataFile           string
		targetsFile            string
		buildpackTomlPathUsage = "full path to the buildpack.toml file, using only one of camelCase, snake_case, or dash_case"
		metadataFileUsage      = "output filename into which to write the JSON metadata, using only one of camelCase, snake_case, or dash_case"
		targetsFileUsage       = `output filename into which to write a JSON array of targets (e.g. ["bionic","jammy"])`
	)

	flag.StringVar(&buildpackTomlPath, "buildpackTomlPath", "", buildpackTomlPathUsage)
	flag.StringVar(&buildpackTomlPath, "buildpack_toml_path", buildpackTomlPath, buildpackTomlPathUsage)
	flag.StringVar(&buildpackTomlPath, "buildpack-toml-path", buildpackTomlPath, buildpackTomlPathUsage)
	flag.StringVar(&metadataFile, "metadataFile", "", metadataFileUsage)
	flag.StringVar(&metadataFile, "metadata_file", metadataFile, metadataFileUsage)
	flag.StringVar(&metadataFile, "metadata-file", metadataFile, metadataFileUsage)
	flag.StringVar(&targetsFile, "targetsFile", "", targetsFileUsage)
	flag.StringVar(&targetsFile, "targets_file", targetsFile, targetsFileUsage)
	flag.StringVar(&targetsFile, "targets-file", targetsFile, targetsFileUsage)
	flag.Parse()

	validate(buildpackTomlPath, metadataFile, targetsFile)

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
			return generateMetadata(hasVersion)
		})
	if err != nil {
		panic(err)
	}

	metadataJson, err := workflows.ToWorkflowJson(dependencies)
	if err != nil {
		panic(fmt.Errorf("unable to marshall metadata json, with error=%w", err))
	}

	if err = os.WriteFile(metadataFile, []byte(metadataJson), os.ModePerm); err != nil {
		panic(fmt.Errorf("cannot write to %s: %w", metadataFile, err))
	} else {
		fmt.Printf("Wrote metadata to %s\n", metadataFile)
	}

	slices.Sort(targets)
	targetsJson, err := workflows.ToWorkflowJson(targets)
	if err != nil {
		panic(fmt.Errorf("unable to marshall targets json, with error=%w", err))
	}

	if err = os.WriteFile(targetsFile, []byte(targetsJson), os.ModePerm); err != nil {
		panic(fmt.Errorf("cannot write to %s: %w", targetsFile, err))
	} else {
		fmt.Printf("Wrote targets to %s\n", targetsFile)
	}
}

func validate(buildpackTomlPath, metadataFile, targetsFile string) {
	if exists, err := fs.Exists(buildpackTomlPath); err != nil {
		panic(err)
	} else if !exists {
		panic(fmt.Errorf("could not locate buildpack.toml at '%s'", buildpackTomlPath))
	}

	if metadataFile == "" {
		panic("metadataFile is required")
	}

	if targetsFile == "" {
		panic("targetsFile is required")
	}
}
