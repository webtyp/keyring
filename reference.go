package keyring

// Scheme is the .env marker meaning "this value lives in the keyring, not on
// disk" — e.g. IAM_CLIENT_SECRET=keyring://. tinywasm/env resolves it via
// auto.OpenForModule; tinywasm/app's SECRETS section is how a developer fills
// it in. Full convention: KEYRING_DOTENV_MASTER_PLAN.md.
const Scheme = "keyring://"

// IsReference reports whether raw is exactly the Scheme marker rather than a
// literal value already stored inline in the .env file.
func IsReference(raw string) bool { return raw == Scheme }
