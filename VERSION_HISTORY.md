# Version History

Complete version history and release notes for mp4Label.

## [v0.2.5] - 2026-02-05

### Overview
Version 0.2.5 adds model annotation comparison functionality specifically designed for algorithm engineers to evaluate model performance.

### New Features

#### Model Annotation Comparison Panel
- **Purpose**: Allow algorithm engineers to compare model-generated annotations with human ground truth
- **4th Panel**: New optional panel appears on the right-most side when configured
- **Dynamic Layout**: Seamlessly switches between 3-panel (annotators) and 4-panel (algorithm engineers) modes
- **Read-only Display**: Model annotations shown in read-only mode with visual distinction (gray background)
- **Same Format**: Model annotations use identical `.txt` file format as human annotations

### Configuration Changes
- **New Field**: `model_annotation_dir` (optional)
- **Location**: Configuration modal and `~/.mp4label/config.json`
- **Validation**: Automatically validates directory existence if specified
- **Backward Compatible**: Existing users experience no changes unless they configure this field

### Technical Implementation
- **Backend**: New API endpoint `/api/model-annotation/:filename`
- **Backend**: Extended `Config` struct with `ModelAnnotationDir` field
- **Frontend**: Dynamic panel visibility based on configuration
- **Frontend**: Separate rendering logic for model annotations
- **Frontend**: Visual styling to differentiate model annotations from editable annotations

### User Workflows

#### For Annotators (Unchanged)
- Continue using 3-panel layout
- No impact on existing annotation workflow
- Model panel remains hidden

#### For Algorithm Engineers (New)
1. Configure model annotation directory in settings
2. Select video to compare
3. View human annotation in center-right panel (editable)
4. View model annotation in right-most panel (read-only)
5. Identify differences and areas for model improvement

### Benefits
- **Quality Evaluation**: Quickly assess model annotation quality
- **Error Analysis**: Identify systematic errors in model predictions
- **Iteration Speed**: Faster feedback loop for model development
- **Side-by-Side**: Easy visual comparison of model vs human annotations

---

## [v0.2.4] - 2026-02-04

### Overview
Version 0.2.4 adds task file support and statistics display for better annotation management.

### New Features

#### 📊 Task File Support
- Specify subset of videos to annotate via text file
- Useful for collaborative annotation workflows
- One video name per line (without .mp4 extension)
- Automatically filters video list

#### 📈 Statistics Display
- Real-time statistics in left sidebar
- Shows: Total, Annotated, Pre-annotated, Unannotated counts
- Updates automatically when annotations saved/deleted

---

## [v0.2.3] - 2026-02-04

### New Features

#### ⌨️ Keyboard Shortcut Optimization
- **Arrow keys timing**: Changed from 5s to 0.5s for more precise control
- **Left arrow (←)**: Rewind 0.5 seconds
- **Right arrow (→)**: Forward 0.5 seconds
- More accurate positioning for frame-by-frame annotation

#### 📋 One-Click Timestamp Copy
- **Click to copy**: Click on the timestamp display to copy to clipboard
- **Visual feedback**: Shows "Copied!" floating message after clicking
- Auto-dismiss after 2 seconds
- Convenient for sharing or documenting timestamps

### Bug Fixes

#### 🐛 Non-Tutorial Checkbox Click Area
- **Fixed**: Only clicking the checkbox itself marks as "non-tutorial"
- **Before**: Clicking anywhere in the right-bottom area triggered the checkbox
- **After**: Must click the actual checkbox input element

### Documentation

#### 📚 Documentation Consolidation
- **VERSION_HISTORY.md**: Comprehensive version history with detailed release notes
- **DOCS.md**: Unified documentation including user guide, build instructions, and testing
- **Cleaned up**: Removed redundant documentation files
- **Internationalization**: All documentation and code comments now in English

### Build System

#### 🔨 Multi-Platform Build
- **Default target**: `make` now builds all platforms by default
- **Single command**: One command to build Windows, Linux, and macOS binaries
- **Simplified workflow**: No need to remember multiple commands

---

## [v0.2.2] - 2026-02-03

### New Features

#### 📦 Web File Embedding
- **Single executable**: Uses Go embed to bundle web directory into executable
- **Portable**: No external files needed, double-click to run
- **Cross-platform**: Windows, macOS, Linux all support single-file distribution
- **Simple deployment**: Just copy one file

#### 🎯 Step Selection State
- **Visual feedback**: Selected step has blue border and background
- **Smart insertion**: Press I key to insert new step below selected one
- **Auto-select**: Newly inserted steps are automatically selected and focused
- **Default selection**: First step auto-selected when loading annotation
- **Delete management**: Selection state adjusts automatically when deleting steps

#### ⚙️ Configuration Dialog Improvements
- **Prevent accidental close**: Clicking outside dialog shows shake animation
- **Quote support**: Automatically recognizes and cleans single/double quotes in paths
- **Whitespace handling**: Auto-trims leading/trailing spaces
- **Better UX**: Users must click "Save" or "Cancel" to close

#### 🔢 Video Numbering System
- Fixed numbering for each video (#0001 - #9999)
- Numbers assigned by list order
- Easy identification and reference

#### 🎛️ Status Filtering
- New filter buttons: All, Unannotated, Pre-annotated, Annotated
- Quick filtering of video list
- Combinable with search function

### Improvements

#### ⏱️ Video Time Display Optimization
- Display format changed from "remaining time" to "played/total duration"
- More intuitive time information

#### 📐 Layout Re-optimization
- Right editor width: 480px → **550px**
- More spacious editing area

#### 🐛 Bug Fixes
- Fixed timestamp seek auto-play issue
- Now maintains pause state after seeking
- Fixed delete button centering issue

---

## [v0.2.1] - 2026-02-03

### New Features

#### ⏱️ Millisecond Precision Timestamps
- Timestamp format upgraded from `mm:ss` to `mm:ss.SSS`
- Supports millisecond precision (e.g., 12:32.766)
- Old format annotations auto-converted (adds .000)
- Inserted timestamps automatically include milliseconds

#### 🔄 Step Drag-and-Drop Sorting
- Support dragging step cards to reorder
- Drag handle (⋮⋮) indicator
- Visual feedback during drag
- Automatic renumbering

#### 📐 Layout Optimization
- Right editor width increased from 400px to 480px
- Better suited for vertical phone video annotation
- Timestamp input width increased from 80px to 100px

#### 🐛 Bug Fixes
- Fixed incorrect timestamp seek position
- Improved timestamp parsing logic for millisecond precision

---

## [v0.2.0] - 2026-02-03

### New Features

#### 🎬 Professional Video Player
- Uses Video.js instead of native HTML5 video tag
- Supports zoom, fullscreen, picture-in-picture
- Better cross-browser compatibility and performance
- Custom control bar for annotation workflow

#### ⚡ Quick Annotation Feature
- **Insert timestamp button**: One-click insert current playback time to new step
- **Auto-pause**: Video auto-pauses after inserting timestamp
- **Auto-focus**: Focus moves to description input automatically
- **Keyboard shortcut**: Press I key to quickly insert timestamp

#### 🎯 Click-to-Seek Feature
- Click timestamp in step list to seek video to that position
- Auto-plays for 1 second then pauses for quick verification
- Timestamp has visual feedback on hover
- Convenient for validating and adjusting annotations

#### ⌨️ Keyboard Shortcut Support
- `Space`: Play/pause video
- `←` (Left arrow): Rewind 5 seconds
- `→` (Right arrow): Forward 5 seconds
- `I`: Insert current timestamp to new step
- `F`: Fullscreen (when player focused)

#### 🎚️ Playback Speed Control
- Added playback speed selector
- Supports 0.5x, 0.75x, 1x, 1.25x, 1.5x, 2x speeds
- Convenient for quick browsing or precise annotation

### Improvements

#### UI Optimization
- Redesigned video control bar layout
- Improved time display style, more prominent
- Optimized video container size for better space utilization
- Added click hint and hover effect to timestamp input

#### User Experience
- Smoother annotation workflow
- Reduced mouse operations, improved annotation efficiency
- Better video playback control precision
- Enhanced video and annotation content interaction

### Documentation

#### New Documentation
- **USER_GUIDE.md**: Detailed user guide
  - Quick start tutorial
  - Complete annotation workflow
  - Feature details and tips
  - FAQ

#### Updated Documentation
- **README.md**: Updated feature list and usage instructions
- **REQUIREMENTS.md**: Added new feature requirements and workflow design

### Technical Changes

- Introduced Video.js 8.10.0 (via CDN)
- Refactored video player initialization logic
- Added timestamp seek function `seekToTimestamp()`
- Added quick insert function `insertCurrentTimestamp()`
- Improved event listener organization

---

## [v0.1.0] - 2026-02-02

### Initial Features

- Video file scanning and list display
- Basic video playback
- Annotation editor (tutorial title, step list)
- Pre-annotation file reading
- Annotation save function
- Configuration management (video dir, pre-annotation dir, output dir)
- Basic search function
- Web interface implementation

### Tech Stack

- Backend: Go 1.22+
- Frontend: HTML5 + CSS + JavaScript (vanilla)
- No external dependencies (Go standard library only)

---

## Version Format

Version numbers follow Semantic Versioning (MAJOR.MINOR.PATCH):

- **MAJOR**: Incompatible API changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

## Build Information

Each version is embedded with version information:

```bash
# View version
./bin/mp4label version

# Output example:
# mp4Label v0.2.3
# Video annotation tool for tutorial videos
```

Version is automatically determined by git tags, defaults to v0.1 if no tags exist.

## Upgrade Guide

### General Upgrade Steps

1. **Backup your work** (optional but recommended)
   ```bash
   cp -r ~/.mp4label ~/.mp4label.backup
   ```

2. **Pull latest code**
   ```bash
   git pull origin main
   ```

3. **Clean and rebuild**
   ```bash
   make clean
   make
   ```

4. **Run new version**
   ```bash
   ./bin/mp4label web
   ```

5. **Force refresh browser**
   - Windows/Linux: `Ctrl+F5`
   - macOS: `Cmd+Shift+R`

### Breaking Changes

#### v0.2.0 → v0.2.1
- Timestamp format changed to millisecond precision
- Old annotations automatically converted
- No manual migration needed

#### v0.1.x → v0.2.x
- Command structure changed to subcommand style
- Old: `./mp4label -port 8080`
- New: `./mp4label web -port 8080`

## Roadmap

### Planned Features

#### v0.3.x
- [ ] Batch annotation export
- [ ] Annotation templates
- [ ] Custom keyboard shortcuts
- [ ] Dark mode support

#### v0.4.x
- [ ] Multi-language support (i18n)
- [ ] Video thumbnail preview
- [ ] Annotation statistics dashboard
- [ ] Team collaboration features

#### v1.0.x
- [ ] Plugin system
- [ ] Advanced video editing integration
- [ ] Cloud storage support
- [ ] API for external tools

## Support

For issues, feature requests, or questions:

1. Check this VERSION_HISTORY.md
2. Review DOCS.md for detailed documentation
3. Submit an issue to the project repository

---

**Thank you for using mp4Label!** 🎬✨
