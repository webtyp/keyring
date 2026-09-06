# Gestión de secretos (keyring)

`webtyp/keyring` proporciona almacenamiento seguro de credenciales sobre cuatro backends:
Windows Credential Manager, macOS Keychain, Linux Secret Service y Navegador (WebCrypto + IndexedDB).

## API

### Keyring genérico (por servicio con inyección de proveedor)

```go
import (
    "webtyp.com/keyring"
    "webtyp.com/keyring/linux"
)

kr, err := keyring.NewKeyring("my-app", linux.New())
if err != nil {
    // keyring no disponible
}
kr.SetLog(log.Printf) // opcional

kr.Set("github_token", "ghp_...")   // guardar
token, err := kr.Get("github_token") // leer
kr.Delete("github_token")            // borrar
```

### Selección automática (shortcut `auto`)

```go
import keyring "webtyp.com/keyring/auto"

kr := keyring.OpenKeyring("my-app")
```

### KeyManager (servicio `updater-cicd`)

```go
import (
    "webtyp.com/keyring"
    "webtyp.com/keyring/auto"
)

km := keyring.New(auto.Provider())
km.Setup("hmac-secret-value", "github-pat-value")
```

## Backend intercambiable (tests)

Las operaciones pasan por la interfaz `Provider`. Los tests inyectan un proveedor en memoria (`tests.NewMemProvider()`) a través del constructor de `Keyring`, permitiendo ejecución paralela segura.

## Requisitos del sistema

- **Linux**: Secret Service sobre `webtyp/dbus` (`gnome-keyring` auto-instalable con `Ensure`).
- **Windows**: Credential Manager vía `advapi32.dll`.
- **macOS**: Keychain vía `/usr/bin/security`.
- **Navegador**: WebCrypto API e IndexedDB vía `syscall/js`.
