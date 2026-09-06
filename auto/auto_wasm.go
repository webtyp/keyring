//go:build wasm

package auto

import (
	"webtyp.com/keyring"
	"webtyp.com/keyring/browser"
)

// Provider returns the WebCrypto/IndexedDB backend.
func Provider() keyring.Provider { return browser.New() }
