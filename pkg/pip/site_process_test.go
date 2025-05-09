// SPDX-FileCopyrightText: Copyright (c) 2013-Present CloudFoundry.org Foundation, Inc. All Rights Reserved.
//
// SPDX-License-Identifier: Apache-2.0

package pipinstall_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/paketo-buildpacks/packit/v2/pexec"
	pipinstall "github.com/paketo-buildpacks/python-packagers/pkg/pip"
	"github.com/paketo-buildpacks/python-packagers/pkg/pip/fakes"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testSiteProcess(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		layerPath  string
		executable *fakes.Executable

		process pipinstall.SiteProcess
	)

	it.Before(func() {
		layerPath = t.TempDir()

		executable = &fakes.Executable{}
		executable.ExecuteCall.Stub = func(execution pexec.Execution) error {
			if execution.Stdout != nil {
				_, err := fmt.Fprintln(execution.Stdout, filepath.Join(layerPath, "/pip/lib/python/site-packages"))
				Expect(err).NotTo(HaveOccurred())
			}
			return nil
		}

		process = pipinstall.NewSiteProcess(executable)
	})

	context("Execute", func() {
		context("there are site packages in the pip layer", func() {
			it("returns the full path to the packages", func() {
				sitePackagesPath, err := process.Execute(layerPath)
				Expect(err).NotTo(HaveOccurred())

				Expect(executable.ExecuteCall.Receives.Execution.Env).To(Equal(append(os.Environ(), fmt.Sprintf("PYTHONUSERBASE=%s", layerPath))))
				Expect(executable.ExecuteCall.Receives.Execution.Args).To(Equal([]string{"-m", "site", "--user-site"}))

				Expect(sitePackagesPath).To(Equal(filepath.Join(layerPath, "pip", "lib", "python", "site-packages")))
			})
		})

		context("failure cases", func() {
			context("site package lookup fails", func() {
				it.Before(func() {
					executable.ExecuteCall.Stub = func(execution pexec.Execution) error {
						_, err := fmt.Fprintln(execution.Stdout, "stdout output")
						Expect(err).NotTo(HaveOccurred())
						_, err = fmt.Fprintln(execution.Stderr, "stderr output")
						Expect(err).NotTo(HaveOccurred())
						return errors.New("locating site packages failed")
					}
				})

				it("returns an error", func() {
					_, err := process.Execute(layerPath)
					Expect(err).To(MatchError(ContainSubstring("failed to locate site packages:")))
					Expect(err).To(MatchError(ContainSubstring("stderr output")))
					Expect(err).To(MatchError(ContainSubstring("error: locating site packages failed")))
				})
			})

			context("when the site process returns nothing", func() {
				it.Before(func() {
					executable.ExecuteCall.Stub = nil
				})

				it("returns an error", func() {
					_, err := process.Execute(layerPath)
					Expect(err).To(MatchError("failed to locate site packages: output is empty"))
				})
			})
		})
	})
}
