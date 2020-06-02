package client

import (
	"os"
	"sync"
	"testing"
)

func TestSingleton(t *testing.T) {
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

func TestCachingSingleton(t *testing.T) {
	ghcInstance = nil
	os.Setenv("GITHUB_TOKEN", "foobar")
	c1, err := CachingSingleton("mem:")
	if err != nil {
		t.Fatal("Error creating the cache")
	}

	if c1 == nil {
		t.Error("should return a client")
	}

	c2, err := CachingSingleton("mem:")
	if err != nil {
		t.Fatal("Error creating the cache")
	}

	if c1 != c2 {
		t.Error("Singletons aren't that singletons")
	}
}

func TestMemClient(t *testing.T) {
	os.Setenv("GITHUB_TOKEN", "foobar")
	_, err := CachingSingleton("mem:")
	if err != nil {
		t.Errorf("Error creating a mem cache")
	}
}

func TestDiskClient(t *testing.T) {
	ghcCachingOnce = sync.Once{}
	_, err := CachingSingleton("file:///dev/null")
	if err != nil {
		t.Errorf("Error creating a disk cache")
	}
}

func TestLRUClient(t *testing.T) {
	ghcCachingOnce = sync.Once{}
	_, err := CachingSingleton("lru:")
	if err != nil {
		t.Errorf("Error creating a LRU cache")
	}
}
func TestInvalidCacheURL(t *testing.T) {
	ghcCachingOnce = sync.Once{}
	_, err := CachingSingleton("foo:")
	if err == nil {
		t.Errorf("Should fail with an unsupported cache")
	}
}
