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

	"github.com/paketo-buildpacks/python-packagers/integration"
)

var (
	buildpack                 string
	minicondaBuildpack        string
	minicondaBuildpackOffline string
	buildPlanBuildpack        string

	buildpackInfo struct {
		Buildpack struct {
			ID   string
			Name string
		}
	}

	config struct {
		MinicondaBuildpack string `json:"miniconda"`
		BuildPlanBuildpack string `json:"build-plan"`
	}
)

func TestCondaIntegration(t *testing.T) {
	// Do not truncate Gomega matcher output
	// The buildpack output text can be large and we often want to see all of it.
	format.MaxLength = 0
	SetDefaultEventuallyTimeout(30 * time.Second)

	Expect := NewWithT(t).Expect

	root, err := filepath.Abs("./../..")
	Expect(err).ToNot(HaveOccurred())

	file, err := os.Open("./../../buildpack.toml")
	Expect(err).NotTo(HaveOccurred())

	_, err = toml.NewDecoder(file).Decode(&buildpackInfo)
	Expect(err).NotTo(HaveOccurred())
	Expect(file.Close()).To(Succeed())

	file, err = os.Open("integration.json")
	Expect(err).NotTo(HaveOccurred())

	Expect(json.NewDecoder(file).Decode(&config)).To(Succeed())
	Expect(file.Close()).To(Succeed())

	buildpackStore := integration_helpers.NewBuildpackStore("conda")

	buildpack, err = buildpackStore.Get.
		WithVersion("1.2.3").
		Execute(root)
	Expect(err).NotTo(HaveOccurred())

	minicondaBuildpack, err = buildpackStore.Get.
		WithVersion("1.2.3").
		Execute(config.MinicondaBuildpack)
	Expect(err).NotTo(HaveOccurred())

	minicondaBuildpackOffline, err = buildpackStore.Get.
		WithVersion("1.2.3").
		WithOfflineDependencies().
		Execute(config.MinicondaBuildpack)
	Expect(err).NotTo(HaveOccurred())

	buildPlanBuildpack, err = buildpackStore.Get.
		Execute(config.BuildPlanBuildpack)
	Expect(err).NotTo(HaveOccurred())

	suite := spec.New("Integration", spec.Report(report.Terminal{}), spec.Parallel())
	suite("Default", testDefault)
	suite("LayerReuse", testLayerReuse)
	suite("LockFile", testLockFile)
	suite("Logging", testLogging)
	suite("Offline", testOffline)
	suite.Run(t)
}
