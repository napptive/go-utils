/*
 *   Copyright 2022 Napptive
 *
 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at
 *
 *        https://www.apache.org/licenses/LICENSE-2.0
 *
 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 */

package github

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"os"
)

var _ = ginkgo.Describe("GitHub tests", func() {

	if run := os.Getenv("RUN_GITHUB_TEST"); run == "" {
		log.Warn().Msg("Github tests are skipped")
		return
	}

	pat := os.Getenv("PAT")
	user := os.Getenv("GITHUB_USER")
	repoUrl := "https://github.com/napptive/event-exporter"
	outputPath := "/tmp_test/"

	ginkgo.It("should be able to clone a repo", func() {
		prov := NewGitHubProvider(pat, user)
		err := prov.Clone(repoUrl, outputPath)
		gomega.Expect(err).Should(gomega.Succeed())
	})

	ginkgo.It("Should be able to get a file from a repo", func() {
		prov := NewGitHubProvider(pat, user)
		content, err := prov.DownloadRepoFile(repoUrl, "./contributing.md")
		gomega.Expect(err).Should(gomega.Succeed())
		gomega.Expect(content).ShouldNot(gomega.BeNil())
		gomega.Expect(content).ShouldNot(gomega.BeEmpty())
		log.Info().Str("content", string(content)).Msg("content file")
	})

	ginkgo.It("Should not be able to get a non existing file from a repo", func() {
		prov := NewGitHubProvider(pat, user)
		_, err := prov.DownloadRepoFile(repoUrl, "./error_file.md")
		gomega.Expect(err).ShouldNot(gomega.Succeed())
	})

})
