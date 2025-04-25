// SPDX-FileCopyrightText: © 2025 Idiap Research Institute <contact@idiap.ch>
// SPDX-FileContributor: Samuel Gaist <samuel.gaist@idiap.ch>
//
// SPDX-License-Identifier: Apache-2.0

package integration_helpers

import (
	"os"
	"path/filepath"

	"github.com/ForestEckhardt/freezer"
	"github.com/ForestEckhardt/freezer/github"
	"github.com/paketo-buildpacks/occam"
	"github.com/paketo-buildpacks/occam/packagers"
)

func NewBuildpackStore(suffix string) occam.BuildpackStore {
	gitToken := os.Getenv("GIT_TOKEN")
	cacheManager := freezer.NewCacheManager(filepath.Join(os.Getenv("HOME"), ".freezer-cache", suffix))
	releaseService := github.NewReleaseService(github.NewConfig("https://api.github.com", gitToken))
	packager := packagers.NewJam()
	namer := freezer.NewNameGenerator()

	return occam.NewBuildpackStore().WithLocalFetcher(
		freezer.NewLocalFetcher(
			&cacheManager,
			packager,
			namer,
		)).WithRemoteFetcher(
		freezer.NewRemoteFetcher(
			&cacheManager,
			releaseService, packager,
		)).WithCacheManager(
		&cacheManager,
	)
}
