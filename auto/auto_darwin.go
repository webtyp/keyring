//go:build darwin

package auto

import (
	"webtyp.com/keyring"
	"webtyp.com/keyring/darwin"
)

// Provider returns the Keychain backend.
func Provider() keyring.Provider { return darwin.New() }
