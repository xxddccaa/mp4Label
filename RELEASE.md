# mp4Label Release Process

Standard release procedure for mp4Label. Follow these steps in order.

## Prerequisites

- All code changes are complete and tested
- Go 1.22+ installed
- Git configured with push access to origin

## Steps

### 1. Determine Version Number

Follow semantic versioning (MAJOR.MINOR.PATCH):
- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)  
- **PATCH**: Bug fixes

### 2. Update Version in Code

Update the default version fallback in `Makefile`:
```makefile
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "vX.Y.Z")
```

Update the version constant in `cmd/mp4label/main.go`:
```go
var version = "vX.Y.Z"
```

### 3. Update Documentation

All of these files must reflect the new version:

| File | What to Update |
|------|---------------|
| `VERSION_HISTORY.md` | Add new version section at top with date, features, fixes |
| `CHANGELOG.md` | Add new version section at top with technical changes |
| `README.md` | Update version number, recent updates list |
| `DOCS.md` | Add docs for any new features, update API list if changed |
| `CLAUDE.md` | Update "Recent Changes" section, version history, API list, last updated date |
| `使用教程.md` | Update user-facing tutorial if UI or workflow changed |

### 4. Build All Platforms

```bash
make clean
make all-platforms
```

Verify:
```bash
./bin/mp4label-darwin version
# Should show: mp4Label vX.Y.Z
```

### 5. Git Commit

```bash
git add -A
git commit -m "Release vX.Y.Z: brief description"
```

### 6. Git Tag

```bash
git tag -a vX.Y.Z -m "Release vX.Y.Z: brief description"
```

### 7. Rebuild with Clean Tag

After tagging, rebuild so binaries embed the clean tag (without `-dirty`):
```bash
make clean
make all-platforms
./bin/mp4label-darwin version
# Should show: mp4Label vX.Y.Z (no -dirty suffix)
```

### 8. Git Push

```bash
git push origin main
git push origin vX.Y.Z
```

Or push all tags at once:
```bash
git push origin main --tags
```

## Quick Checklist

- [ ] Version updated in Makefile and main.go
- [ ] VERSION_HISTORY.md updated
- [ ] CHANGELOG.md updated
- [ ] README.md version updated
- [ ] DOCS.md updated (if features changed)
- [ ] CLAUDE.md updated
- [ ] 使用教程.md updated (if UI changed)
- [ ] `make clean && make all-platforms` succeeds
- [ ] `./bin/mp4label-darwin version` shows correct version
- [ ] Git commit created
- [ ] Git tag created
- [ ] Rebuild after tag (clean version string)
- [ ] Push code and tags to origin

## Notes

- The `VERSION` in Makefile uses `git describe --tags` which auto-derives version from tags
- After tagging, the binary version string will be clean (e.g., `v0.2.9` instead of `v0.2.7-dirty`)
- Always rebuild after tagging to embed the correct version
- Binary outputs go to `bin/` directory (gitignored)
