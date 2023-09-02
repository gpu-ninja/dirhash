# dirhash

Cryptographically checksums a directory and its contents. Includes support for signing checksums with Ed25519.

## Usage

### Hash and Sign a Directory

```
dirhash -k ~/.ssh/id_ed25519 /path/to/dir
```

### Verify a Signed Directory

```
dirhash verify -k ~/.ssh/id_ed25519.pub <hash> /path/to/dir
```