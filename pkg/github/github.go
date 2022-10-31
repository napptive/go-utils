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
	"errors"
	"fmt"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/napptive/nerrors/pkg/nerrors"
	"github.com/rs/zerolog/log"
	"os"
	"os/exec"
	"path/filepath"
)

// GHProvider with a struct to manage GitHub repositories
type GHProvider struct {
	// PersonalAccessToken with a Personal API token
	PersonalAccessToken string
	// GitHubUser with the name of the gitHub user
	GitHubUser string
}

// NewGitHubProvider returns a GHProvider to manage GitHub operations
func NewGitHubProvider(pat string, user string) *GHProvider {

	return &GHProvider{
		PersonalAccessToken: pat,
		GitHubUser:          user,
	}
}

// Clone Downloads a repository in a path
func (gp *GHProvider) Clone(repoURL string, outputPath string) error {

	log.Debug().Str("repoURL", repoURL).Str("outputPath", outputPath).Msg("cloning repository")

	if err := gp.createIfNotExists(outputPath); err != nil {
		return err
	}

	var cloneOptions git.CloneOptions

	if gp.PersonalAccessToken != "" {
		cloneOptions = git.CloneOptions{
			Auth: &http.BasicAuth{
				Username: "git",
				Password: gp.PersonalAccessToken,
			},
			URL:      repoURL,
			Progress: os.Stdout,
		}
	} else {
		cloneOptions = git.CloneOptions{
			URL:      repoURL,
			Progress: os.Stdout,
		}
	}

	if _, err := git.PlainClone(outputPath, false, &cloneOptions); err != nil {
		log.Error().Err(err).Msg("error cloning repo")
		return nerrors.NewInternalErrorFrom(err, "error cloning repo")
	}

	log.Debug().Str("repo url", repoURL).Msg("repo successfully cloned")

	return nil
}

// ConfigurePrivateRepositories configures git options to download private repositories if required
func (gp *GHProvider) ConfigurePrivateRepositories() {

	if gp.PersonalAccessToken != "" && gp.GitHubUser != "" {
		// 	git config --global  url."https://${user}:${personal_access_token}@github.com".insteadOf  "https://github.com"
		args := []string{"config", "--global", fmt.Sprintf("url.https://%s:%s@github.com.insteadOf", gp.GitHubUser, gp.PersonalAccessToken), "https://github.com"}
		toExecute := exec.Command("git", args...)
		stdOut, err := toExecute.CombinedOutput()
		if err != nil {
			log.Error().Err(err).Str("stdOut", string(stdOut)).Msg("error executing git config command.")
		}
	} else {
		log.Warn().Msg("If there is any private repo the launcher could failed. Fill RepositoryAccessToken and RepositoryUser")
	}
}

// DownloadRepoFile downloads a file from a repository
func (gp *GHProvider) DownloadRepoFile(repoURL string, filePath string) ([]byte, error) {
	log.Debug().Str("repoURL", repoURL).Msg("getting file from repository")

	//if err := gp.createIfNotExists(outputPath); err != nil {
	//	return nil, err
	//}

	// TODO: tmp_dir
	dir, err := os.MkdirTemp("", "clone-repo")
	if err != nil {
		log.Error().Err(err).Str("url repo", repoURL).Str("file path", filePath).Msg("error creating temp directory")
		return nil, nerrors.NewInternalErrorFrom(err, "error downloading file from a repo")
	}
	defer os.RemoveAll(dir)

	if err := gp.Clone(repoURL, dir); err != nil {
		log.Error().Err(err).Str("url repo", repoURL).Str("file path", filePath).Msg("error in clone repository")
		return nil, nerrors.NewInternalErrorFrom(err, "error downloading file from a repo")
	}

	//file, err := os.Open(filepath.Join(dir.Name(), filePath))
	//if err != nil {
	//	log.Error().Err(err).Str("url repo", repoURL).Str("file path", filePath).Msg("error opening file")
	//	return nerrors.NewInternalErrorFrom(err, "error downloading file from a repo")
	//}

	content, err := os.ReadFile(filepath.Join(dir, filePath))
	if err != nil {
		log.Error().Err(err).Str("url repo", repoURL).Str("file path", filePath).Msg("error creating output file")
		return nil, nerrors.NewInternalErrorFrom(err, "error downloading file from a repo")
	}

	//outputFile, err := os.Create(outputPath)
	//if err != nil {
	//	log.Error().Err(err).Str("url repo", repoURL).Str("file path", filePath).Msg("error creating output file")
	//	return nerrors.NewInternalErrorFrom(err, "error downloading file from a repo")
	//}
	//
	//if _, err := io.Copy(outputFile, file); err != nil {
	//	log.Error().Err(err).Str("url repo", repoURL).Str("file path", filePath).Msg("error copying file")
	//	return nerrors.NewInternalErrorFrom(err, "error downloading file from a repo")
	//}

	log.Debug().Str("repo url", repoURL).Str("file path", filePath).Msg("File successfully downloaded")

	return content, nil
}

// createIfNotExists check if a directory exists and create it if it does not
func (gp *GHProvider) createIfNotExists(outputPath string) error {

	_, err := os.Stat(outputPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Warn().Str("outputPath", outputPath).Msg("directory does not exist, creating it!")
			if err := os.Mkdir(outputPath, os.ModePerm); err != nil {
				log.Error().Err(err).Msg("error creating the out path.")
				return nerrors.NewFailedPreconditionErrorFrom(err, "error creating the out path")
			}
		} else {
			log.Error().Err(err).Msg("error asking for the out path.")
			return nerrors.NewFailedPreconditionErrorFrom(err, "error asking for the output path")
		}
	}
	return nil
}
