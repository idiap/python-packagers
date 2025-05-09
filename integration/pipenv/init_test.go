// SPDX-FileCopyrightText: Copyright (c) 2013-Present CloudFoundry.org Foundation, Inc. All Rights Reserved.
//
// SPDX-License-Identifier: Apache-2.0

package integration_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"

	integration_helpers "github.com/paketo-buildpacks/python-packagers/integration"
)

var buildpackInfo struct {
	Buildpack struct {
		ID   string
		Name string
	}
	Metadata struct {
		Dependencies []struct {
			Version string
		}
	}
}

var settings struct {
	Buildpacks struct {
		CPython struct {
			Online  string
			Offline string
		}
		Pip struct {
			Online  string
			Offline string
		}
		Pipenv struct {
			Online  string
			Offline string
		}
		PipenvInstall struct {
			Online  string
			Offline string
		}
		BuildPlan struct {
			Online string
		}
	}

	Config struct {
		CPython   string `json:"cpython"`
		Pip       string `json:"pip"`
		Pipenv    string `json:"pipenv"`
		BuildPlan string `json:"build-plan"`
	}
}

func TestPipenvIntegration(t *testing.T) {
	// Do not truncate Gomega matcher output
	// The buildpack output text can be large and we often want to see all of it.
	format.MaxLength = 0
	Expect := NewWithT(t).Expect

	file, err := os.Open("integration.json")
	Expect(err).NotTo(HaveOccurred())

	Expect(json.NewDecoder(file).Decode(&settings.Config)).To(Succeed())
	Expect(file.Close()).To(Succeed())

	file, err = os.Open("./../../buildpack.toml")
	Expect(err).NotTo(HaveOccurred())

	_, err = toml.NewDecoder(file).Decode(&buildpackInfo)
	Expect(err).NotTo(HaveOccurred())
	Expect(file.Close()).To(Succeed())

	root, err := filepath.Abs("./../..")
	Expect(err).ToNot(HaveOccurred())

	buildpackStore := integration_helpers.NewBuildpackStore("pipenv")

	settings.Buildpacks.PipenvInstall.Online, err = buildpackStore.Get.
		WithVersion("1.2.3").
		Execute(root)
	Expect(err).NotTo(HaveOccurred())

	settings.Buildpacks.PipenvInstall.Offline, err = buildpackStore.Get.
		WithVersion("1.2.3").
		WithOfflineDependencies().
		Execute(root)
	Expect(err).NotTo(HaveOccurred())

	settings.Buildpacks.Pip.Online, err = buildpackStore.Get.
		Execute(settings.Config.Pip)
	Expect(err).NotTo(HaveOccurred())

	settings.Buildpacks.Pip.Offline, err = buildpackStore.Get.
		WithOfflineDependencies().
		Execute(settings.Config.Pip)
	Expect(err).NotTo(HaveOccurred())

	settings.Buildpacks.Pipenv.Online, err = buildpackStore.Get.
		Execute(settings.Config.Pipenv)
	Expect(err).NotTo(HaveOccurred())

	settings.Buildpacks.Pipenv.Offline, err = buildpackStore.Get.
		WithOfflineDependencies().
		Execute(settings.Config.Pipenv)
	Expect(err).NotTo(HaveOccurred())

	settings.Buildpacks.CPython.Online, err = buildpackStore.Get.
		Execute(settings.Config.CPython)
	Expect(err).NotTo(HaveOccurred())

	settings.Buildpacks.CPython.Offline, err = buildpackStore.Get.
		WithOfflineDependencies().
		Execute(settings.Config.CPython)
	Expect(err).NotTo(HaveOccurred())

	settings.Buildpacks.BuildPlan.Online, err = buildpackStore.Get.
		Execute(settings.Config.BuildPlan)
	Expect(err).NotTo(HaveOccurred())

	SetDefaultEventuallyTimeout(30 * time.Second)

	suite := spec.New("Integration", spec.Report(report.Terminal{}))
	suite("Default", testDefault, spec.Parallel())
	suite("Offline", testOffline, spec.Parallel())
	suite.Run(t)
}
