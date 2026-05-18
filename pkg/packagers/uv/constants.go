// SPDX-FileCopyrightText: © 2026 Idiap Research Institute <contact@idiap.ch>
// SPDX-FileContributor: Samuel Gaist <samuel.gaist@idiap.ch>
//
// SPDX-License-Identifier: Apache-2.0

package uvinstall

const (
	// UvEnvLayer is the name of the layer into which uv environment is installed.
	UvEnvLayer = "uv-env"

	// UvEnvCache is the name of the layer that is used as the uv package directory.
	UvEnvCache = "uv-env-cache"

	// UvEnvPlanEntry is the name of the Build Plan requirement that this buildpack provides.
	UvEnvPlanEntry = "uv-environment"

	// UvPlanEntry is the name of the Build Plan requirement for the uv
	// dependency that this buildpack requires.
	UvPlanEntry = "uv"

	// UvEnvLayerCacheSha is the key in the Layer Content Metadata used to determine if layer
	// can be reused.
	UvEnvLayerCacheSha = "layer-cache-sha"

	// LockfileName is the name of the export file from which the buildpack reinstalls packages
	LockfileName = "uv.lock"

	// Config environmental variables

	// Used to specify one or more directories to pass to `--find-links`
	UvFindLinks = "BP_UV_FIND_LINKS"

	// List of additional groups divided by comma which should be installed
	UvInstallGroups = "BP_UV_INSTALL_GROUPS"

	// UV will assert that the uv.lock remains unchanged
	UvLocked = "BP_UV_LOCKED"

	// UV will compile Python source files to bytecode after installation
	UvCompileByteCode = "BP_UV_COMPILE_BYTECODE"

	// Enables preview mode for uv
	UvPreview = "BP_UV_PREVIEW"
)
