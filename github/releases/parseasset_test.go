package releases

import (
	"testing"
)

func testAsset(asset *ReleaseAsset, t *testing.T) {
	if asset.Arch != "386" {
		t.Errorf("Error parsing architecture %s", asset.Arch)
	}

	if asset.OS != "freebsd" {
		t.Errorf("Error parsing OS %s", asset.OS)
	}

	if asset.Version != "v1.51.0" {
		t.Error("Error parsing version")
	}

	if asset.Name != "foobar" {
		t.Errorf("Error parsing name %s", asset.Name)
	}
}

func TestParseAsset(t *testing.T) {
	asset := ParseAsset("foobar-v1.51.0-freebsd-386.zip")
	testAsset(asset, t)

	if asset.Filename != "foobar-v1.51.0-freebsd-386.zip" {
		t.Error("Error parsing filename")
	}

	asset = ParseAsset("foobar-v1.51.0-FreeBSD-i386.zip")
	testAsset(asset, t)
}

func TestBinNoArchive(t *testing.T) {
	asset := ParseAsset("foobar-x86_64")

	if asset.Name != "foobar" {
		t.Errorf("Error parsing name %s", asset.Name)
	}

	if asset.Arch != "amd64" {
		t.Errorf("Error parsing arch %s", asset.Arch)
	}
}

func TestArchNormalization(t *testing.T) {
	asset := ParseAsset("foobar-v1.51.0-macos-x86_64.zip")
	if asset.Arch != "amd64" {
		t.Errorf("Error normalizing arch %s", asset.Arch)
	}
}

func TestOSNormalization(t *testing.T) {
	asset := ParseAsset("foobar-v1.51.0-macos-386.zip")
	if asset.OS != "darwin" {
		t.Error("Error normalizing OS")
	}

	asset = ParseAsset("foobar-v1.51.0-mac-386.zip")
	if asset.OS != "darwin" {
		t.Error("Error normalizing OS")
	}

	asset = ParseAsset("foobar-v1.51.0-darwin-386.zip")
	if asset.OS != "darwin" {
		t.Error("Error normalizing OS")
	}

	asset = ParseAsset("foobar-v1.51.0-osx-386.zip")
	if asset.OS != "darwin" {
		t.Error("Error normalizing OS")
	}

	asset = ParseAsset("foobar-v1.51.0-win-386.zip")
	if asset.OS != "windows" {
		t.Error("Error normalizing OS")
	}
}

func TestParseVersionWithoutV(t *testing.T) {
	asset := ParseAsset("foobar-1.2.3-freebsd-386.zip")

	if asset.Version != "1.2.3" {
		t.Error("Error parsing version")
	}
}

func TestParseMissingOS(t *testing.T) {
	asset := ParseAsset("foobar-v1.51.0-386.zip")

	if asset.OS != "unknown" {
		t.Error("Error parsing OS")
	}
}

func TestParseMissingVersion(t *testing.T) {
	asset := ParseAsset("foobar-freebsd-386.zip")

	if asset.Version != "unknown" {
		t.Error("Error parsing version")
	}
}

func TestParseWeirdNames(t *testing.T) {
	asset := ParseAsset("foobarWE-sdfsd--we-1.2.3.4-freebsd-386.zip")

	if asset.Name != "foobarWE-sdfsd--we" {
		t.Errorf("Error parsing name, found %s", asset.Name)
	}

	asset = ParseAsset("foobarWE1.2.3.4-freebsd-386.zip")
	if asset.Name != "foobarWE1.2.3.4" {
		t.Errorf("Error parsing name, found %s", asset.Name)
	}
	if asset.Version != "unknown" {
		t.Errorf("Error parsing unknown, found %s", asset.Version)
	}
}
