// +build mage

package main

import (
	"os"
	"path/filepath"

	"github.com/nocquidant/mageproj/mgp"
)

const (
	projectName = "go-hello"
	groupName   = "nocquidant"
	buildDir    = "build"
	ldFlags     = "-s -X main.version=$VERSION -X main.buildDate=$BUILD_DATE"
	dckRegistry = "" // default
	gitRepo     = "github.com"
)

var proj *mgp.MageProject

func init() {
	proj = &mgp.MageProject{
		ProjectName: projectName,
		GroupName:   groupName,
		BuildDir:    buildDir,
		PackageName: filepath.Join(gitRepo, groupName, projectName),
		LdFlags:     ldFlags,
		DckRegistry: dckRegistry,
		DckImage:    filepath.Join(dckRegistry, groupName, projectName),
	}
	proj = mgp.InitMageProject(currentDir(), proj)
}

func currentDir() string {
	workdir, err := os.Getwd()
	if err != nil {
		workdir = "."
	}
	return workdir
}

// Validate runs go format and linters
func Validate() error {
	return proj.Validate()
}

// Test runs tests with go test
func Test() error {
	return proj.Test()
}

// Build builds binary in build dir
func Build() error {
	return proj.Build()
}

// Clean removes the build directory
func Clean() error {
	return proj.Clean()
}

// DockerBuild builds Docker image
func DockerBuild() error {
	return proj.DockerBuild()
}

// DockerPush pushes Docker image to Artifactory
func DockerPush() error {
	if os.Getenv("DOCKER_USR") == "" {
		os.Setenv("DOCKER_USR", os.Getenv("DOCKER_USER"))
	}
	if os.Getenv("DOCKER_PWD") == "" {
		os.Setenv("DOCKER_PWD", os.Getenv("DOCKER_PASS"))
	}
	return proj.DockerPush()
}

// PrintInfo prints information used internally
func PrintInfo() {
	proj.PrintInfo()
}
