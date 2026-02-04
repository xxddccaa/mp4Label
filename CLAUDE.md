# mp4Label - AI Development Guide

This document provides context for AI assistants (like Claude) to help develop this project effectively.

## Project Overview

**Project Name:** mp4Label  
**Purpose:** Video annotation tool for tutorial videos  
**Language:** Go (backend) + Vanilla JavaScript (frontend)  
**Architecture:** Single-page web application with embedded resources  
**Target Platforms:** Windows, Linux, macOS (cross-platform)  

## Key Technical Characteristics

### 1. Single-File Distribution
- Web resources are **embedded** into the Go binary using `go:embed`
- Two `web/` directories exist:
  - `web/` - Source directory
  - `cmd/mp4label/web/` - Temporary copy for embedding (auto-generated during build)
- Build process: `prepare-web` copies files, then embeds them

### 2. Version Management

**Version Format:** Semantic Versioning (MAJOR.MINOR.PATCH)
- Example: v0.2.3
- Automatically derived from git tags
- Default: v0.1 if no tags exist

**Version is embedded in binary:**
```bash
./bin/mp4label-darwin version
# Output: mp4Label v0.2.3
```

**Update version:**
1. Update in `Makefile` (default VERSION)
2. Create git tag if needed
3. Update `VERSION_HISTORY.md`
4. Update `CHANGELOG.md`
5. Update `README.md` version section

### 3. Build System

**Makefile Targets:**
```makefile
make                # Build all platforms (default)
make cmd-mp4label   # Build current platform only
make windows        # Build Windows binary
make linux          # Build Linux binary
make darwin         # Build macOS binary
make clean          # Remove all build artifacts
make run            # Build and run
```

**Build Process:**
1. `prepare-web` - Copy web/ to cmd/mp4label/web/
2. Build with embedded files
3. `clean-web` - Remove temporary web copy

**Output Location:** `bin/`
- `mp4label.exe` (Windows)
- `mp4label-linux` (Linux)
- `mp4label-darwin` (macOS)

## Release Process

### When Creating a New Release

1. **Code Changes**
   - Implement features/fixes
   - Test thoroughly on all platforms

2. **Version Number**
   - Increment version in `Makefile` (VERSION variable)
   - Follow semantic versioning rules:
     - MAJOR: Breaking changes
     - MINOR: New features (backward compatible)
     - PATCH: Bug fixes

3. **Documentation Updates (CRITICAL)**
   - Update `VERSION_HISTORY.md`:
     - Add new version section at top
     - Document all changes with categories
     - Include date (YYYY-MM-DD format)
   - Update `CHANGELOG.md`:
     - Same content as VERSION_HISTORY.md
     - Keep detailed technical notes
   - Update `README.md`:
     - Update version number in "Version" section
     - Add to recent updates list

4. **Build All Platforms**
   ```bash
   make clean
   make  # Builds all platforms by default
   ```

5. **Verify Build**
   ```bash
   ./bin/mp4label-darwin version
   # Should show new version number
   ```

6. **Git Tag (Optional but Recommended)**
   ```bash
   git tag -a v0.2.4 -m "Release v0.2.4: Description"
   git push origin v0.2.4
   ```

## Documentation Standards

### Core Documentation Files

1. **README.md** - Project overview and quick start (internationalized to English)
2. **DOCS.md** - Complete documentation (user guide, build instructions, testing)
3. **VERSION_HISTORY.md** - Comprehensive version history
4. **CHANGELOG.md** - Detailed technical changelog
5. **CLAUDE.md** - This file (AI development guide)

### Documentation Rules

- **Language:** All documentation in English (as of v0.2.3)
- **Format:** Markdown with clear structure
- **Updates:** Must update docs when releasing new version
- **No Redundancy:** Avoid duplicate documentation files

### What NOT to Create

❌ Don't create these types of files (we consolidated them):
- BUILD.md, BUILD_WINDOWS.md
- QUICKSTART.md, USER_GUIDE.md
- IMPROVEMENTS.md, OPTIMIZATION_*.md
- UPDATE_v*.md, VERSION_UPDATE.md
- REQUIREMENTS.md, SUMMARY.md
- TESTING.md, WINDOWS_USAGE.md

✅ Instead, add content to existing docs:
- Build info → DOCS.md
- User guides → DOCS.md
- Version updates → VERSION_HISTORY.md
- Requirements → DOCS.md

## Code Standards

### Go Code

**Style:**
- Use `gofmt` formatting
- Clear function names
- Comments in English
- Error handling with proper messages

**Packages:**
- `cmd/mp4label/` - Main entry point
- `pkg/server/` - HTTP server
- `pkg/annotation/` - Annotation parsing and validation
- `pkg/video/` - Video file handling
- `pkg/config/` - Configuration management

**No External Dependencies:**
- Use only Go standard library
- Keep binary size small
- Maximize portability

### Frontend Code

**Structure:**
- `web/index.html` - HTML structure
- `web/static/css/style.css` - Styles
- `web/static/js/app.js` - JavaScript logic

**Standards:**
- Vanilla JavaScript (no frameworks)
- English for all UI text and comments
- Video.js for video player
- Clear function names with JSDoc-style comments

**Key JavaScript Patterns:**
```javascript
// Global state management
let currentVideo = null;
let currentAnnotation = null;
let player = null; // Video.js instance

// Async/await for API calls
async function loadVideos() {
    const response = await fetch('/api/videos');
    videos = await response.json();
}

// Event-driven architecture
setupEventListeners();
```

## Common Development Tasks

### Adding a New Feature

1. **Code Changes**
   - Backend: Add to appropriate package in `pkg/`
   - Frontend: Update `web/` files
   - Test locally

2. **Documentation**
   - Add feature description to VERSION_HISTORY.md
   - Update CHANGELOG.md
   - Update DOCS.md if user-facing

3. **Version Bump**
   - MINOR version for new features
   - Update Makefile VERSION

4. **Build and Test**
   ```bash
   make clean
   make
   ./bin/mp4label-darwin web
   # Test in browser
   ```

### Fixing a Bug

1. **Code Changes**
   - Identify and fix issue
   - Test thoroughly

2. **Documentation**
   - Add to VERSION_HISTORY.md under "Bug Fixes"
   - Update CHANGELOG.md
   - Note the bug and fix

3. **Version Bump**
   - PATCH version for bug fixes
   - Update Makefile VERSION

### Updating UI/UX

1. **HTML Changes:** `web/index.html`
2. **Styling:** `web/static/css/style.css`
3. **Behavior:** `web/static/js/app.js`
4. **Rebuild:** `make clean && make`

### Keyboard Shortcuts

Current shortcuts (as of v0.2.3):
- `Space` - Play/Pause
- `←` - Rewind 0.5 seconds
- `→` - Forward 0.5 seconds
- `I` - Insert timestamp

**When modifying:**
- Update both player shortcuts and global shortcuts in app.js
- Update README.md keyboard shortcuts table
- Update DOCS.md keyboard shortcuts section

## Testing Checklist

Before releasing a new version, test:

### Build Testing
- [ ] `make clean` works
- [ ] `make` builds all 3 platforms
- [ ] All binaries created in `bin/`
- [ ] Version number correct in binary

### Functionality Testing
- [ ] Video list loads
- [ ] Video plays
- [ ] Timestamp insertion works
- [ ] Save annotation works
- [ ] Load pre-annotation works
- [ ] Search/filter works
- [ ] Keyboard shortcuts work
- [ ] Drag-and-drop step reordering works
- [ ] Click timestamp to copy works
- [ ] Click timestamp in step list seeks video

### UI/UX Testing
- [ ] All text in English
- [ ] Responsive layout
- [ ] No console errors
- [ ] Modal dialogs work correctly
- [ ] Non-tutorial checkbox only responds to checkbox click

### Cross-Platform Testing
- [ ] Windows binary works
- [ ] Linux binary works
- [ ] macOS binary works

## API Endpoints Reference

```
GET  /                           Main page
GET  /api/videos                 List videos
GET  /api/video/:filename        Stream video
GET  /api/annotation/:filename   Get annotation
POST /api/annotation/:filename   Save annotation
DELETE /api/annotation/:filename Delete annotation
GET  /api/config                 Get config
POST /api/config                 Save config
```

## Project History Context

### Major Milestones

- **v0.1.0** (2026-02-02): Initial release
- **v0.2.0** (2026-02-03): Video.js player, quick annotation features
- **v0.2.1** (2026-02-03): Millisecond precision, drag-and-drop
- **v0.2.2** (2026-02-03): Video numbering, status filtering
- **v0.2.3** (2026-02-04): Keyboard shortcuts (0.5s), copy timestamp, bug fixes, internationalization

### Recent Changes (v0.2.3)

1. **Keyboard shortcuts:** Arrow keys now 0.5s instead of 5s
2. **Copy timestamp:** Click timestamp display to copy
3. **Bug fix:** Non-tutorial checkbox click area fixed
4. **Documentation:** Consolidated from 20+ files to 5 files
5. **Internationalization:** All UI and docs in English
6. **Build system:** `make` now builds all platforms by default

## Important Notes for AI Assistants

### When Making Changes

1. **Always update documentation** when changing features
2. **Test the build** after code changes
3. **Follow existing patterns** in the codebase
4. **Keep it simple** - no unnecessary dependencies
5. **English only** - all new text should be in English

### Version History Format

Use this structure in VERSION_HISTORY.md:

```markdown
## [vX.Y.Z] - YYYY-MM-DD

### New Features
- Feature description with emoji

### Improvements
- Improvement description

### Bug Fixes
- Bug fix description

### Documentation
- Documentation changes
```

### When User Requests New Features

1. Confirm understanding of the feature
2. Implement the feature
3. Update all relevant documentation
4. Update version number
5. Test the build
6. Provide summary of changes

### File Locations Quick Reference

```
Source code:
- Go: cmd/mp4label/*.go, pkg/**/*.go
- HTML: web/index.html
- CSS: web/static/css/style.css
- JS: web/static/js/app.js

Documentation:
- Overview: README.md
- Complete: DOCS.md
- Versions: VERSION_HISTORY.md
- Changelog: CHANGELOG.md
- AI Guide: CLAUDE.md (this file)

Build:
- Makefile
- go.mod
- bin/ (output)
```

## Troubleshooting Common Issues

### Build Fails
- Check if `web/` directory exists
- Ensure Go 1.22+ installed
- Run `make clean` first

### Version Number Wrong
- Update VERSION in Makefile
- Rebuild: `make clean && make`

### Web Files Not Updating
- Run `make clean` to remove embedded files
- Rebuild to re-embed

### Documentation Out of Sync
- Check VERSION_HISTORY.md matches CHANGELOG.md
- Ensure README.md version number updated

## Contact & Support

For questions or issues:
1. Check DOCS.md
2. Check VERSION_HISTORY.md
3. Check this file (CLAUDE.md)
4. Review recent git commits for context

---

**Last Updated:** 2026-02-04 (v0.2.3)  
**Maintained for:** AI-assisted development with Claude and similar tools
