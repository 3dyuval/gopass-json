package main

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/spf13/cobra"
)

const (
	repoOwner = "3dyuval"
	repoName  = "gopass-json"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Check for a newer release and update the binary",
	RunE: func(cmd *cobra.Command, args []string) error {
		current, err := semver.ParseTolerant(version)
		if err != nil {
			return fmt.Errorf("could not parse current version %q: %w", version, err)
		}

		latest, latestTag, err := fetchLatestVersion()
		if err != nil {
			return fmt.Errorf("could not fetch latest release: %w", err)
		}

		fmt.Printf("current: v%s\n", current)
		fmt.Printf("latest:  %s\n", latestTag)

		if !latest.GT(current) {
			fmt.Println("already up to date")
			return nil
		}

		fmt.Printf("updating to %s...\n", latestTag)

		binPath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("could not determine binary location: %w", err)
		}
		binPath, err = filepath.EvalSymlinks(binPath)
		if err != nil {
			return fmt.Errorf("could not resolve binary path: %w", err)
		}

		if err := downloadAndReplace(latestTag, binPath); err != nil {
			return fmt.Errorf("update failed: %w", err)
		}

		fmt.Printf("updated to %s at %s\n", latestTag, binPath)
		return nil
	},
}

func init() {
	root.AddCommand(updateCmd)
}

type githubRelease struct {
	TagName string `json:"tag_name"`
}

func fetchLatestVersion() (semver.Version, string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", repoOwner, repoName)
	resp, err := http.Get(url) //nolint:noctx
	if err != nil {
		return semver.Version{}, "", err
	}
	defer resp.Body.Close()

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return semver.Version{}, "", err
	}

	v, err := semver.ParseTolerant(strings.TrimPrefix(release.TagName, "v"))
	return v, release.TagName, err
}

func downloadAndReplace(tag, binPath string) error {
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	archiveName := fmt.Sprintf("%s-%s-%s-%s.tar.gz", repoName, tag, goos, goarch)
	url := fmt.Sprintf("https://github.com/%s/%s/releases/download/%s/%s", repoOwner, repoName, tag, archiveName)

	resp, err := http.Get(url) //nolint:noctx
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: HTTP %d for %s", resp.StatusCode, url)
	}

	newBin, err := extractBinary(resp.Body, repoName)
	if err != nil {
		return err
	}
	defer os.Remove(newBin)

	if err := os.Chmod(newBin, 0o755); err != nil {
		return err
	}

	return os.Rename(newBin, binPath)
}

func extractBinary(r io.Reader, binaryName string) (string, error) {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return "", err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		if filepath.Base(hdr.Name) != binaryName {
			continue
		}

		tmp, err := os.CreateTemp("", binaryName+"-*")
		if err != nil {
			return "", err
		}
		if _, err := io.Copy(tmp, tr); err != nil {
			tmp.Close()
			os.Remove(tmp.Name())
			return "", err
		}
		tmp.Close()
		return tmp.Name(), nil
	}

	return "", fmt.Errorf("binary %q not found in archive", binaryName)
}
