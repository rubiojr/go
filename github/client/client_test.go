package client

import (
	"os"
	"testing"
)

func TestClient(t *testing.T) {
	os.Setenv("GITHUB_TOKEN", "foobar")
	c1, _ := Singleton()
	c2, _ := Singleton()

	if c1 == nil {
		t.Error("Invalid client returned")
	}

	if c1 != c2 {
		t.Error("Singletons aren't that singletons")
	}
}
