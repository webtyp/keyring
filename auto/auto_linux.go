//go:build linux

package auto

import (
	"webtyp.com/keyring"
	"webtyp.com/keyring/linux"
)

// Provider returns the Secret Service backend.
func Provider() keyring.Provider { return linux.New() }
