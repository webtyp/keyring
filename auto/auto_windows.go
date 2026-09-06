//go:build windows

package auto

import (
	"webtyp.com/keyring"
	"webtyp.com/keyring/windows"
)

// Provider returns the Credential Manager backend.
func Provider() keyring.Provider { return windows.New() }
