package registry

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

// JavaVersion represents a Java version available for download
type JavaVersion struct {
	Version      string
	MajorVersion int
	Distribution string
	OS           string
	Arch         string
	DownloadURL  string
	Checksum     string
	FileName     string
}

// Registry manages available Java versions
type Registry struct {
	versions []JavaVersion
}

// NewRegistry creates a new registry instance
func NewRegistry() *Registry {
	return &Registry{
		versions: []JavaVersion{},
	}
}

// AvailableReleasesResponse represents the response from Adoptium's available_releases endpoint
type AvailableReleasesResponse struct {
	AvailableReleases []int `json:"available_releases"`
}

// AdoptiumRelease represents a release from Adoptium API
type AdoptiumRelease struct {
	Binary struct {
		OS           string `json:"os"`
		Architecture string `json:"architecture"`
		ImageType    string `json:"image_type"`
		Package      struct {
			Name     string `json:"name"`
			Link     string `json:"link"`
			Checksum string `json:"checksum"`
		} `json:"package"`
	} `json:"binary"`
	Version struct {
		Major    int `json:"major"`
		Minor    int `json:"minor"`
		Security int `json:"security"`
		Build    int `json:"build"`
	} `json:"version"`
}

// FetchAvailableVersions fetches available Java versions from Adoptium API
func (r *Registry) FetchAvailableVersions() error {
	// First, get the list of all available versions from Adoptium
	availableVersions, err := r.fetchAvailableVersionsList()
	if err != nil {
		return fmt.Errorf("failed to fetch available versions list: %w", err)
	}

	// Fetch each available version
	for _, majorVersion := range availableVersions {
		if err := r.fetchVersionFromAdoptium(majorVersion); err != nil {
			// Log error but continue with other versions
			fmt.Printf("Warning: Failed to fetch Java %d: %v\n", majorVersion, err)
		}
	}

	// Sort versions by major version (descending)
	sort.Slice(r.versions, func(i, j int) bool {
		return r.versions[i].MajorVersion > r.versions[j].MajorVersion
	})

	return nil
}

// fetchAvailableVersionsList fetches the list of available Java versions from Adoptium
func (r *Registry) fetchAvailableVersionsList() ([]int, error) {
	url := "https://api.adoptium.net/v3/info/available_releases"

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch available releases: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var releasesInfo AvailableReleasesResponse
	if err := json.Unmarshal(body, &releasesInfo); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return releasesInfo.AvailableReleases, nil
}

// fetchVersionFromAdoptium fetches a specific version from Adoptium API
func (r *Registry) fetchVersionFromAdoptium(majorVersion int) error {
	url := fmt.Sprintf("https://api.adoptium.net/v3/assets/latest/%d/hotspot", majorVersion)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch from Adoptium API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var releases []AdoptiumRelease
	if err := json.Unmarshal(body, &releases); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Determine current OS and Arch
	targetOS := runtime.GOOS
	if targetOS == "darwin" {
		targetOS = "mac"
	}

	targetArch := runtime.GOARCH
	if targetArch == "amd64" {
		targetArch = "x64"
	}
	// arm64 is usually aarch64 in Adoptium, but sometimes just arm64.
	// Adoptium API uses: x64, x32, ppc64, s390x, ppc64le, aarch64, arm
	if targetArch == "arm64" {
		targetArch = "aarch64"
	}

	// Extract binaries matching current OS/Arch
	for _, release := range releases {
		// Strict matching for OS and Architecture
		if release.Binary.OS == targetOS &&
			release.Binary.Architecture == targetArch &&
			release.Binary.ImageType == "jdk" {

			version := fmt.Sprintf("%d.%d.%d+%d",
				release.Version.Major,
				release.Version.Minor,
				release.Version.Security,
				release.Version.Build)

			r.versions = append(r.versions, JavaVersion{
				Version:      version,
				MajorVersion: release.Version.Major,
				Distribution: "Temurin",
				OS:           release.Binary.OS,
				Arch:         release.Binary.Architecture,
				DownloadURL:  release.Binary.Package.Link,
				Checksum:     release.Binary.Package.Checksum,
				FileName:     release.Binary.Package.Name,
			})
		}
	}

	return nil
}

// GetVersions returns all available versions
func (r *Registry) GetVersions() []JavaVersion {
	return r.versions
}

// FindVersion finds a version by major version number or full version string
func (r *Registry) FindVersion(versionStr string) (*JavaVersion, error) {
	// Try to parse as major version number
	if majorVersion, err := strconv.Atoi(versionStr); err == nil {
		for _, v := range r.versions {
			if v.MajorVersion == majorVersion {
				return &v, nil
			}
		}
		return nil, fmt.Errorf("version %d not found", majorVersion)
	}

	// Try to find exact match
	for _, v := range r.versions {
		if v.Version == versionStr {
			return &v, nil
		}
	}

	// Try partial match
	for _, v := range r.versions {
		if strings.HasPrefix(v.Version, versionStr) {
			return &v, nil
		}
	}

	return nil, fmt.Errorf("version %s not found", versionStr)
}

// GetMajorVersions returns unique major versions
func (r *Registry) GetMajorVersions() []int {
	seen := make(map[int]bool)
	var majors []int

	for _, v := range r.versions {
		if !seen[v.MajorVersion] {
			seen[v.MajorVersion] = true
			majors = append(majors, v.MajorVersion)
		}
	}

	sort.Sort(sort.Reverse(sort.IntSlice(majors)))
	return majors
}

// FindLatestForMajor finds the latest version for a specific major version
func (r *Registry) FindLatestForMajor(majorVersion int) (*JavaVersion, error) {
	// If versions are not loaded, fetch them
	if len(r.versions) == 0 {
		if err := r.FetchAvailableVersions(); err != nil {
			return nil, fmt.Errorf("failed to fetch versions: %w", err)
		}
	}

	// Find all versions matching the major version
	var candidates []JavaVersion
	for _, v := range r.versions {
		if v.MajorVersion == majorVersion {
			candidates = append(candidates, v)
		}
	}

	if len(candidates) == 0 {
		return nil, fmt.Errorf("no versions found for Java %d", majorVersion)
	}

	// Return the first one (versions are already sorted with latest first)
	return &candidates[0], nil
}
