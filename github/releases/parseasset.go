package releases

import (
	"fmt"
	"regexp"
	"strings"
)

// ReleaseAsset represents a GitHub release asset
//
// Fields are normalized and lowercase.
type ReleaseAsset struct {
	Filename string
	// Normalized architecture
	Arch string
	// Normalized operating system
	OS        string
	Version   string
	IsArchive bool
	IsBinary  bool
	Name      string
}

var archNormMap = map[string]string{
	"386":    "386",
	"i386":   "386",
	"32bit":  "386",
	"i686":   "386",
	"x86":    "386",
	"amd64":  "amd64",
	"64bit":  "amd64",
	"x86_64": "amd64",
}

var osNormMap = map[string]string{
	"linux":   "linux",
	"lin":     "linux",
	"darwin":  "darwin",
	"osx":     "darwin",
	"macos":   "darwin",
	"mac":     "darwin",
	"windows": "windows",
	"win":     "windows",
	"mswin":   "windows",
}

var versionReg = regexp.MustCompile(`(?i)^v?\d+(\d|\.)+$`)
var archReg = regexp.MustCompile(`(?i)^arm|mips|mipsle|amd64|x86_64|x64|i?(3|6)86|(32|64)bit|x86$`)
var osReg = regexp.MustCompile(`(?i)^dragonflybsd|freebsd|linux|darwin|mac|macos|osx|windows|win|netbsd|openbsd|solaris|plan9$`)
var extReg = regexp.MustCompile(`\.(tgz|tar\.(gz|bz|bz2)|zip|7zip|tbz2)$`)

// ParseAsset parses the asset file name to figure out the
// architecture, version, operating system, etc from a GitHub
// release asset.
//
func ParseAsset(name string) *ReleaseAsset {
	nameWithoutExt := extReg.ReplaceAllString(name, "")

	// We don't want to tokenize the x86_64 arch
	nameWithoutExt = strings.Replace(nameWithoutExt, "x86_64", "amd64", 1)
	var dashUnder = regexp.MustCompile(`[-_]`)
	nameWithoutExt = dashUnder.ReplaceAllString(nameWithoutExt, " ")

	tokens := strings.Split(nameWithoutExt, " ")

	asset := &ReleaseAsset{
		Arch:    "unknown",
		OS:      "unknown",
		Version: "unknown",
		Name:    "",
	}
	asset.Name = tokens[0]
	asset.Filename = name
	asset.IsArchive = isArchive(name)

	for _, t := range tokens[1:] {
		switch {
		case archReg.MatchString(t):
			asset.Arch = normalizeArch(strings.ToLower(t))
			asset.IsBinary = true
		case versionReg.MatchString(t) && !archReg.MatchString(t):
			asset.Version = t
		case osReg.MatchString(t):
			asset.OS = normalizeOS(strings.ToLower(t))
		default:
			if asset.Name == "" {
				asset.Name = t
			} else {
				asset.Name = fmt.Sprintf("%s-%s", asset.Name, t)
			}
		}
	}

	return asset
}

func normalizeArch(arch string) string {
	narch := archNormMap[arch]
	if narch != "" {
		return narch
	}

	return arch
}

func normalizeOS(os string) string {
	nos := osNormMap[os]
	if nos != "" {
		return nos
	}

	return os
}

func isArchive(name string) bool {
	return strings.HasSuffix(name, ".zip") ||
		strings.HasSuffix(name, ".tar.gz") ||
		strings.HasSuffix(name, ".tgz") ||
		strings.HasSuffix(name, ".tar")
}
