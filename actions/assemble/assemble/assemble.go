package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/paketo-buildpacks/libdependency"
	"github.com/paketo-buildpacks/libdependency/workflows"
	"github.com/paketo-buildpacks/packit/v2/cargo"
	"github.com/paketo-buildpacks/packit/v2/fs"
)

func main() {
	var (
		id                     string
		buildpackTomlPath      string
		artifactPath           string
		buildpackTomlPathUsage = "full path to the buildpack.toml file, using only one of camelCase, snake_case, or dash_case"
		artifactPathUsage      = "full path to the directory containing metadata.json and compiled tarballs"
	)

	flag.StringVar(&id, "id", "", "id of the dependency")
	flag.StringVar(&buildpackTomlPath, "buildpackTomlPath", "", buildpackTomlPathUsage)
	flag.StringVar(&buildpackTomlPath, "buildpack_toml_path", buildpackTomlPath, buildpackTomlPathUsage)
	flag.StringVar(&buildpackTomlPath, "buildpack-toml-path", buildpackTomlPath, buildpackTomlPathUsage)
	flag.StringVar(&artifactPath, "artifactPath", "", artifactPathUsage)
	flag.StringVar(&artifactPath, "artifact_path", artifactPath, artifactPathUsage)
	flag.StringVar(&artifactPath, "artifact-path", artifactPath, artifactPathUsage)
	flag.Parse()

	fmt.Printf("id=%s\n", id)
	fmt.Printf("artifactPath=%s\n", artifactPath)
	fmt.Printf("buildpackTomlPath=%s\n", buildpackTomlPath)

	validate(id, artifactPath)

	prepareCommit(id, artifactPath, buildpackTomlPath)
}

func validate(id, artifactPath string) {
	if id == "" {
		panic("id is required")
	}

	if exists, err := fs.Exists(artifactPath); err != nil {
		panic(err)
	} else if !exists {
		panic(fmt.Errorf("directory %s does not exist", artifactPath))
	} else if fs.IsEmptyDir(artifactPath) {
		panic(fmt.Errorf("directory %s is empty", artifactPath))
	}
}

func prepareCommit(id, artifactPath, buildpackTomlPath string) {
	config, err := libdependency.ParseBuildpackToml(buildpackTomlPath)
	if err != nil {
		panic(fmt.Errorf("failed to parse buildpack.toml: %w", err))
	}

	artifacts := findArtifacts(artifactPath, id)

	fmt.Println("Found artifacts:")
	var item interface{} = artifacts
	workflowJson, err := workflows.ToWorkflowJson(item)
	if err != nil {
		panic(err)
	}
	fmt.Println(workflowJson)

	config.Metadata.Dependencies = append(config.Metadata.Dependencies, artifacts...)

	config = libdependency.PruneConfig(config)

	file, err := os.OpenFile(buildpackTomlPath, os.O_RDWR|os.O_TRUNC, 0600)
	if err != nil {
		panic(fmt.Errorf("failed to open buildpack config file: %w", err))
	}
	defer file.Close()

	err = cargo.EncodeConfig(file, config)
	if err != nil {
		panic(fmt.Errorf("failed to write buildpack config: %w", err))
	}
}

func findArtifacts(artifactDir string, id string) []cargo.ConfigMetadataDependency {
	var artifacts []cargo.ConfigMetadataDependency

	tarballGlob := filepath.Join(artifactDir, fmt.Sprintf("%s-*", id))
	if allDirsForArtifacts, err := filepath.Glob(tarballGlob); err != nil {
		panic(err)
	} else if len(allDirsForArtifacts) < 1 {
		fmt.Printf("no compiled artifact folders found: %s\n", tarballGlob)
		os.Exit(0)
	} else {
		fmt.Printf("Found compiled artifact folders:\n")
		for _, singleDirForArtifact := range allDirsForArtifacts {
			artifactBaseName := filepath.Base(singleDirForArtifact)
			fmt.Printf("- %s\n", artifactBaseName)

			metadata := getMetadata(singleDirForArtifact)
			target := getTarget(singleDirForArtifact)

			metadata.SHA256 = getSHA256AndValidate(singleDirForArtifact, artifactBaseName)
			metadata.URI = "TBD - needs to be uploaded to S3"
			metadata.Stacks = []string{target}
			artifacts = append(artifacts, metadata)
		}
	}
	return artifacts
}

func getSHA256AndValidate(artifactDir, artifactBaseName string) string {
	calculatedSHA256, err := fs.NewChecksumCalculator().Sum(filepath.Join(artifactDir, fmt.Sprintf("%s.tgz", artifactBaseName)))
	if err != nil {
		panic(err)
	}

	tarballSHA256 := getSHA256(artifactDir, artifactBaseName)

	if !strings.HasPrefix(tarballSHA256, calculatedSHA256) {
		fmt.Printf("SHA256 does not match! Expected=%s, Calculated=%s\n", tarballSHA256, calculatedSHA256)
		panic("SHA256 does not match!")
	}

	return calculatedSHA256
}

func getTarget(artifactDir string) string {
	targetBytes, err := os.ReadFile(filepath.Join(artifactDir, "target"))
	if err != nil {
		panic(err)
	}

	return strings.TrimSpace(string(targetBytes))
}

func getMetadata(artifactDir string) cargo.ConfigMetadataDependency {
	metadataBytes, err := os.ReadFile(filepath.Join(artifactDir, "metadata.json"))
	if err != nil {
		panic(err)
	}

	var metadata cargo.ConfigMetadataDependency
	err = json.Unmarshal(metadataBytes, &metadata)
	if err != nil {
		panic(err)
	}

	return metadata
}

func getSHA256(artifactDir string, artifactBaseName string) string {
	sha256Bytes, err := os.ReadFile(filepath.Join(artifactDir, fmt.Sprintf("%s.tgz.sha256", artifactBaseName)))
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(sha256Bytes))
}
