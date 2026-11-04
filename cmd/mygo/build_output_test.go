package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildRejectsOutputInResources(t *testing.T) {
	for _, web := range []bool{false, true} {
		c := resourceProject(t, `{"out":"resources"}`)
		var err error
		if web {
			_, err = buildWeb(c, buildOptions{})
		} else {
			_, err = buildPlatform(c, "linux", "amd64", buildOptions{})
		}
		if err == nil || !strings.Contains(err.Error(), "build output") || !strings.Contains(err.Error(), "resource directory") {
			t.Fatalf("web=%v: %v; want output/resource overlap rejected before compilation", web, err)
		}
	}
}

func TestBuildOutputDir(t *testing.T) {
	for _, tc := range []struct {
		out    string
		reject bool
	}{
		{".", false}, {"dist", false}, {"legal", false},
		{"resources", true}, {"resources/generated", true},
		{"vendor/models", true}, {"vendor/models/build", true},
	} {
		t.Run(tc.out, func(t *testing.T) {
			c := resourceProject(t, `{"resources":["legal/NOTICE.txt","vendor/models"]}`)
			c.Out = tc.out
			_, err := c.buildOutputDir()
			if (err != nil) != tc.reject {
				t.Fatalf("output %q: %v", tc.out, err)
			}
		})
	}
	c := resourceProject(t, `{}`)
	c.Out = t.TempDir()
	if _, err := c.buildOutputDir(); err != nil {
		t.Fatal(err)
	}
}

func TestBuildOutputDirSymlink(t *testing.T) {
	c := resourceProject(t, `{}`)
	link := filepath.Join(c.root, "output-link")
	if err := os.Symlink(c.path("resources"), link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	// The output does not exist yet, but its parent resolves inside resources.
	c.Out = filepath.Join(link, "generated")
	if _, err := c.buildOutputDir(); err == nil {
		t.Fatal("accepted output through a symlink to resources")
	}
}
