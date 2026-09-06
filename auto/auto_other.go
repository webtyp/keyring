//go:build !linux && !darwin && !windows && !wasm

package auto

import "webtyp.com/keyring"

// Provider returns Fallback on unsupported platforms.
func Provider() keyring.Provider { return keyring.Fallback{} }
