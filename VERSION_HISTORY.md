# Version History

Complete version history and release notes for mp4Label.

## [v0.2.9] - 2026-02-11

### Overview
Version 0.2.9 adds **Task Directory** support — a directory containing multiple `.txt` task files, each representing a task group. Users can switch between task groups via a dropdown in the sidebar, making it easy for teams to review their own assigned video sets.

### New Features

#### 📂 Task Directory (Task Groups)
- **New config field**: `task_dir` — a directory containing multiple `.txt` task files
- **Sidebar dropdown**: When configured, a task group selector appears above the search box
- **Dynamic switching**: Select a task group to filter the video list instantly (no reload)
- **"All Videos" option**: First option in dropdown shows ALL videos unfiltered (ignores both task_dir and task_file)
- **Live video count**: Label shows current video count (e.g., "Task Group (42 videos):")
- **Visual prominence**: Dropdown styled with distinct blue border and background for easy identification
- **Both modes**: Works in both `web` (edit) and `view` (read-only) modes
- **Config modal**: New "Task Directory" input with browse button in configuration dialog
- **Persisted**: Saved to `~/.mp4label/config.json` as `task_dir` field

#### 🔌 New API Endpoint
- `GET /api/task-groups` — returns `{"available": true, "groups": ["task_001.txt", "task_002.txt", ...]}`
- `GET /api/videos?task=filename.txt` — filters videos by a specific task file from TaskDir

#### ⚙️ CLI Changes
- `mp4label web -task-dir /path/to/tasks` — specify task directory in edit mode
- `mp4label view ... -task-dir /path/to/tasks` — specify task directory in view mode
- **Conflict handling**: When both `-task-dir` and `-task-file` are specified, `-task-dir` takes priority with a clear `WARNING` message

### Priority Rules
- `-task-dir` > `-task-file` (when both specified)
- CLI flags > config file values
- "All Videos" selection ignores all task filtering

### Security
- Path traversal protection: `?task=` parameter rejects `..`, `/`, `\`
- File must end with `.txt`
- Absolute path verification ensures file is within TaskDir

### Use Cases
- **Team annotation**: Each annotator has their own task .txt file, selects their group
- **Batch review**: QA reviewer switches between groups to check different annotator's work
- **Large projects**: Split 1000s of videos into manageable groups of ~20-50

### Example
```bash
# Directory structure
/tasks/
  ├── annotator_alice.txt    # 50 video names
  ├── annotator_bob.txt      # 50 video names
  └── annotator_carol.txt    # 53 video names

# Start in view mode
mp4label view -video-dir /videos -output-dir /output -task-dir /tasks

# Or in edit mode
mp4label web -task-dir /tasks
```

---

## [v0.2.8] - 2026-02-11

### Overview
Version 0.2.8 adds a **read-only view mode** via the new `view` subcommand. This allows reviewing completed annotations without risk of accidental modification, with all config parameters specified via CLI flags.

### New Features

#### 👁️ Read-Only View Mode (`view` subcommand)
- **New subcommand**: `mp4label view` starts the server in read-only mode
- **CLI config**: All config parameters specified via command-line flags (no config file read/write)
  - `-video-dir` (required)
  - `-output-dir` (required)
  - `-pre-annotation-dir` (optional)
  - `-task-file` (optional)
  - `-model-annotation-dir` (optional)
  - `-port` (default: 8080)
- **Full lockdown**: All write operations return HTTP 403 Forbidden:
  - Save annotation (POST `/api/annotation/`)
  - Delete annotation (DELETE `/api/annotation/`)
  - Save config (POST `/api/config`)
  - Open file dialog (GET `/api/dialog`)
- **Frontend enforcement**: UI disabled in addition to backend protection
  - Save, Delete, Add Step, Insert Timestamp buttons disabled
  - All inputs/textareas set to read-only
  - Drag-and-drop step reordering disabled
  - Step remove buttons hidden
  - Auto-save system completely disabled
  - `I` key shortcut disabled
  - Config modal: inputs read-only, browse buttons hidden, save button hidden
- **Visual indicators**:
  - Red "View Only" badge in header
  - "🔒 View Only" status in editor
  - Disabled buttons have reduced opacity

#### 🔌 New API Endpoint
- `GET /api/mode` — returns `{"readonly": true/false}` for frontend to detect server mode

### What Works in View Mode
- Video playback (play, pause, seek, speed control, fullscreen)
- Video list browsing and filtering
- Viewing annotations (all content displayed)
- Model annotation comparison panel
- Click timestamp to seek video
- Copy timestamp to clipboard
- Keyboard shortcuts for playback (Space, ←, →)
- Config modal viewing (read-only)

### Use Cases
- **Annotation review**: Review completed annotations without accidental edits
- **Demonstration**: Share annotated results with stakeholders
- **Quality assurance**: QA team can review without modification access
- **Portable deployment**: Specify exact paths via CLI without touching config file

### Examples
```bash
# Basic view mode
mp4label view -video-dir /path/to/videos -output-dir /path/to/annotations

# Custom port with all options
mp4label view -port 9090 \
  -video-dir /data/videos \
  -output-dir /data/annotations \
  -pre-annotation-dir /data/pre-annotations \
  -task-file /data/task.txt \
  -model-annotation-dir /data/model-annotations
```

---

## [v0.2.7] - 2026-02-06

### Overview
Version 0.2.7 adds real-time auto-save for annotations and fixes a 1ms timestamp precision issue.

### New Features

#### 💾 Auto-Save Annotations
- **Real-time saving**: Annotations are automatically saved 1.5 seconds after any edit
- **Debounced writes**: Prevents excessive API calls; groups rapid edits into a single save
- **Visual status indicator**: Shows "Unsaved", "Saving...", "Saved", or "Save failed" next to editor header
- **Flush before navigation**: Switching videos immediately saves pending changes
- **Silent validation**: Only auto-saves when annotation data is complete and valid
- **Change detection**: Compares current data with last saved state to skip unnecessary saves
- **Manual save still works**: "Save" button remains available for explicit saves (cancels pending auto-save)

### Bug Fixes

#### 🐛 1ms Timestamp Precision Fix
- **Problem**: Inserting or copying timestamps showed a 1ms discrepancy from actual video time
- **Root cause**: `Math.floor((seconds - wholeSecs) * 1000)` suffered from IEEE 754 floating-point precision loss
- **Fix**: Now converts to total integer milliseconds via `Math.round(timeInSeconds * 1000)` first, then extracts minutes, seconds, and millis from that integer
- **Result**: Timestamps now match video time exactly with zero offset

### Technical Implementation

#### Auto-Save System
- `scheduleAutoSave()`: Debounced scheduler (1.5s delay)
- `performAutoSave()`: Validates, checks for changes, then sends POST request
- `flushAutoSave()`: Immediate save for video switching
- `syncAnnotationFromForm()`: Syncs form inputs to data model before saving
- `validateAnnotationSilently()`: Checks data completeness without alert dialogs
- `updateAutoSaveStatus(status)`: Updates the visual status indicator

#### Trigger Points
- Tutorial title input
- Step timestamp change/input (including paste)
- Step description input
- Add step / Remove step
- Drag-and-drop step reorder
- Insert timestamp (I key or button)
- Non-tutorial checkbox toggle

#### Timestamp Fix
- `formatTimestamp()`: Uses integer arithmetic to avoid floating-point rounding errors
- Before: `Math.floor((seconds - wholeSecs) * 1000)` → could lose 1ms
- After: `Math.round(timeInSeconds * 1000)` → exact millisecond value

### User Workflow Improvement
- **Before**: Annotators lost all work when navigating to another video without saving
- **After**: All edits are automatically saved in real-time; no data loss
- **Before**: Timestamps off by 1ms required manual correction
- **After**: Timestamps match video time exactly

---

## [v0.2.6] - 2026-02-05

### Overview
Version 0.2.6 adds native OS file/folder dialogs so users can pick paths directly.

### New Features

#### 🗂️ Native OS Dialogs
- **Browse Button (📁📄)**: Click to open native system dialog
- **No Browser Limitation**: Uses backend to launch OS dialogs, returns full absolute path
- **Mode Aware**:
  - 📁 for directories (Video, Pre-annotation, Output, Model directories)
  - 📄 for files (Task file)
- **One Click**: Pick a path and it is filled into the input field

### Technical Implementation

#### Backend (Go)
- New API endpoint: `GET /api/dialog?mode=file|directory`
- Launches native OS dialog and returns selected absolute path
- macOS: `osascript` `choose folder` / `choose file`
- Windows: PowerShell `FolderBrowserDialog` / `OpenFileDialog`
- Linux: `zenity` with file/folder selection

#### Frontend
- Reused 📁📄 buttons in config dialog
- `openBrowser(inputId, mode)` calls backend dialog API and fills input
- Fallback alert if dialog cannot be opened

### User Workflow
1. Click Config button
2. Click 📁 or 📄 button next to any path input
3. Use native OS dialog to pick a folder or file
4. Path automatically fills into input field

### Benefits
- **No Manual Typing**: Browse visually instead of typing paths
- **Fewer Errors**: Select existing paths, avoiding typos
- **Faster Setup**: Quick navigation through filesystem
- **User Friendly**: No need to remember exact paths
- **Cross-platform**: Works on macOS, Linux, Windows

---

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
