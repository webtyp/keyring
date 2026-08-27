package auto

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/tinywasm/keyring"
)

const (
	// ErrGoModNotFound is returned when no go.mod exists in dir or any parent.
	ErrGoModNotFound = keyring.Error("keyring/auto: no go.mod found walking up from the given directory")
	// ErrNoModuleDirective is returned when go.mod has no `module` line.
	ErrNoModuleDirective = keyring.Error("keyring/auto: go.mod has no module directive")
)

// findModulePath walks up from dir looking for go.mod and returns the value
// of its `module` directive (e.g. "github.com/veltylabs/misitio").
func findModulePath(dir string) (string, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	for {
		data, err := os.ReadFile(filepath.Join(dir, "go.mod"))
		if err == nil {
			for _, line := range strings.Split(string(data), "\n") {
				line = strings.TrimSpace(line)
				if after, ok := strings.CutPrefix(line, "module "); ok {
					return strings.TrimSpace(after), nil
				}
			}
			return "", ErrNoModuleDirective
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", ErrGoModNotFound
		}
		dir = parent
	}
}

// OpenForModule opens (auto-selecting the platform's provider, same as
// OpenKeyring) a Keyring scoped to the Go module found by walking up from
// dir — pass "." for the current working directory. The module's own path is
// the service namespace: it is already globally unique, so two different
// projects' identically named keys never collide, and no project has to
// remember to declare a service name of its own.
func OpenForModule(dir string) (*Keyring, error) {
	service, err := findModulePath(dir)
	if err != nil {
		return nil, err
	}
	return OpenKeyring(service), nil
}
