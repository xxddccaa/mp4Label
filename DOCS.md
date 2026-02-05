# mp4Label Documentation

Complete documentation for mp4Label - a video annotation tool for tutorial videos.

## Table of Contents

- [Quick Start](#quick-start)
- [User Guide](#user-guide)
- [Build Instructions](#build-instructions)
- [Testing Guide](#testing-guide)
- [Keyboard Shortcuts](#keyboard-shortcuts)
- [Configuration](#configuration)
- [Troubleshooting](#troubleshooting)

---

## Quick Start

### 1. Get Started in 5 Minutes

#### macOS/Linux Users

```bash
# 1. Navigate to project directory
cd /path/to/mp4Label

# 2. Build for current platform
make

# 3. Run
./bin/mp4label web
```

#### Windows Users

```bash
# 1. Navigate to project directory
cd /path/to/mp4Label

# 2. Build Windows version
make windows

# 3. Run (or copy bin/mp4label.exe to Windows system)
bin/mp4label.exe web
```

### 2. Open in Browser

Visit: `http://localhost:8080`

### 3. Configure Paths

1. Click "Config" button in top-right corner
2. Set paths:
   - **Video Directory**: Directory containing MP4 files (required)
   - **Pre-annotation Directory**: Directory with pre-annotation files (optional)
   - **Output Directory**: Where to save annotation files (required)
   - **Task File**: Text file with video names to annotate (optional)
   - **Model Annotation Directory**: Directory with model-generated annotations (optional, for algorithm engineers)
3. Click "Save"

Configuration is automatically saved to `~/.mp4label/config.json`

---

## User Guide

### Annotation Workflow

#### Quick Annotation (Recommended)

This is the most efficient way to annotate videos:

**Step 1: Select Video**
- Find the video in the left sidebar
- Click to load it

**Step 2: Enter Tutorial Title**
- Type the tutorial title in the right editor
- Example: "Motion Blur Effect Tutorial"

**Step 3: Watch and Annotate**
1. Press `Space` or click play to start video
2. When video reaches a key step:
   - **Press `I` key** (or click "Insert Timestamp" button)
   - Video auto-pauses
   - New step created with current timestamp
   - Cursor auto-focuses to description field
3. Type step description (e.g., "Import footage to editor")
4. Press `Space` to continue playing
5. Repeat steps 2-4 for all key steps

**Step 4: Verify Annotation**
- Review step list on right
- **Click any timestamp** to seek video to that position
- Verify timestamp accuracy
- Edit if needed

**Step 5: Save**
- Click "Save" button in top-right
- Annotation saved to output directory

#### Editing Pre-annotations

If video has AI pre-annotation:

1. Select video - pre-annotation loads automatically
2. Video list item shows "Pre-annotated" badge
3. Verify and edit:
   - **Click step timestamps** to view corresponding video frames
   - Fix inaccurate timestamps or descriptions
   - Delete extra steps (click × button)
   - Add missing steps (use `I` key)
4. Save corrected annotation

### Interface Layout

```
┌─────────────────────────────────────────────────────────────────┐
│  mp4Label - Video Annotation Tool                    [Config]   │
├──────────┬──────────────────────────────┬─────────────────────┤
│          │                              │  Annotation Editor   │
│  Video   │       Video Player           │             [Delete] │
│  List    │                              │             [Save]   │
│          │                              │                      │
│ [Search] │                              │  Tutorial Title:     │
│ ─────────│                              │  ___________         │
│ [All]    │                              │                      │
│ [None]   │                              │  Steps:              │
│ [Pre]    │                              │  1) 00:11.230 desc   │
│ [Done]   │                              │  2) 00:25.560 desc   │
│ ─────────│                              │  3) 00:42.890 desc   │
│ #0001    │                              │                      │
│ video1   │                              │  [Add Step]          │
│ #0002    │                              │                      │
│ video2   │                              │  [ ] Non-tutorial    │
│   ...    ├──────────────────────────────┤                      │
│          │ video_name.mp4    00:15.230  │                      │
│          │ [⏱Insert] [Speed: 1x]        │                      │
└──────────┴──────────────────────────────┴─────────────────────┘
```

### Features

#### Video Playback Control

**Player Buttons**
- **Play/Pause**: Click player center or press Space
- **Progress bar**: Drag to seek to any position
- **Volume**: Click volume icon to adjust
- **Fullscreen**: Click fullscreen button or double-click video

**Playback Speed**
- Select speed below player: 0.5x, 0.75x, 1x, 1.25x, 1.5x, 2x
- Recommendations:
  - First viewing: 1x to understand content
  - Annotating: 0.75x or 0.5x for easier key moment capture
  - Quick browsing: 1.5x or 2x

#### Annotation Editing

**Timestamp Format**
- Format: `mm:ss.SSS` (minutes:seconds.milliseconds)
- Examples:
  - `00:11.000` (11 seconds)
  - `01:23.456` (1 minute 23 seconds 456 milliseconds)
  - `10:05.789` (10 minutes 5 seconds 789 milliseconds)
- Legacy format: `mm:ss` auto-converts to `mm:ss.000`

**Three Ways to Insert Timestamp**
1. **Press I key** (fastest)
   - Press I when video reaches key position
   - Auto-creates new step and pauses video
   
2. **Click button**
   - Click "Insert Timestamp" button
   - Same effect as I key

3. **Manual add**
   - Click "Add Step" button
   - Manually enter timestamp and description

**Click to Seek**
- Click **timestamp** in step list (not description)
- Video seeks to that time point
- Auto-pauses for verification

**Delete Step**
- Click `×` button in step top-right corner
- Step numbers automatically renumber

**Drag-and-Drop Sorting**
- Click drag handle (⋮⋮) on left of step
- Drag to target position
- Release mouse to complete
- Step numbers automatically renumber

**Copy Timestamp**
- Click the timestamp display below video player
- Timestamp copied to clipboard
- Shows "Copied!" message

#### Video List Status

Status badges in video list:

- 🟡 **Pre-annotated**: AI pre-annotation file exists
- 🟢 **Annotated**: Annotation file saved
- 🔴 **None**: No annotation

#### Status Filtering

Four filter buttons above video list:

| Button | Description | Filter Condition |
|--------|-------------|------------------|
| **All** | Show all videos | No filter |
| **None** | Unannotated videos | No pre-annotation & no annotation |
| **Pre** | Pre-annotated videos | Has pre-annotation & no annotation |
| **Done** | Annotated videos | Has annotation |

Combinable with search function.

#### Search Function

Enter keywords in search box at top of video list for real-time filtering.

---

## Build Instructions

### Prerequisites

- Go 1.22 or higher
- Git (for version tagging)

### Build Commands

#### Build All Platforms (Default)

```bash
make
```

This generates:
- `bin/mp4label.exe` - Windows 64-bit
- `bin/mp4label-linux` - Linux 64-bit
- `bin/mp4label-darwin` - macOS 64-bit

#### Build Specific Platform

```bash
# Current platform only
make cmd-mp4label

# Windows 64-bit
make windows

# Linux 64-bit
make linux

# macOS 64-bit
make darwin
```

#### Clean Build Artifacts

```bash
make clean
```

#### Build and Run

```bash
make run
```

#### Install to System (macOS/Linux)

```bash
make install  # Installs to /usr/local/bin/
```

### Build Output

All compiled files are in `bin/` directory:

```
bin/
├── mp4label          # Current platform (macOS in this example)
├── mp4label.exe      # Windows version
├── mp4label-linux    # Linux version
└── mp4label-darwin   # macOS version (cross-compiled)
```

### Version Information

Version is automatically determined by git tags:

```bash
# View version
./bin/mp4label version

# Output:
# mp4Label v0.2.3
# Video annotation tool for tutorial videos
```

Default version is v0.1 if no git tags exist.

### Single-File Distribution

✨ **Important**: All builds have web files embedded

- ✅ Only need the executable file
- ✅ No external files needed
- ✅ Can be directly distributed to users
- ✅ Cross-platform compatible

**Windows users**: Just `mp4label.exe` file  
**Linux users**: Just `mp4label-linux` file  
**macOS users**: Just `mp4label` or `mp4label-darwin` file

### Cross-Compilation

Supports building Windows version on macOS/Linux:

```bash
# Build Windows version
make windows

# Build all platforms
make all-platforms
```

### Build Process Details

The build uses Go's `embed` feature to bundle web resources:

1. **prepare-web**: Copies `web/` directory to `cmd/mp4label/web/`
2. **build**: Compiles Go code with embedded web files
3. **clean-web**: Removes temporary web copy

This is why there are two `web/` directories - see BUILD.md for details.

---

## Testing Guide

### Manual Testing Checklist

#### Basic Functionality
- [ ] Launch application successfully
- [ ] Configuration dialog opens and saves
- [ ] Video list loads and displays
- [ ] Video plays in player
- [ ] Can add/edit/delete steps
- [ ] Can save annotation
- [ ] Annotation file created in output directory

#### Video Playback
- [ ] Play/pause works (button and Space key)
- [ ] Progress bar seeking works
- [ ] Volume control works
- [ ] Fullscreen works
- [ ] Playback speed adjustment works

#### Keyboard Shortcuts
- [ ] Space: Play/pause
- [ ] ← (Left): Rewind 0.5 seconds
- [ ] → (Right): Forward 0.5 seconds
- [ ] I: Insert timestamp
- [ ] Shortcuts don't trigger in input fields

#### Annotation Features
- [ ] Insert timestamp button works
- [ ] I key inserts timestamp
- [ ] Timestamp auto-pauses video
- [ ] Focus moves to description field
- [ ] Click timestamp seeks video
- [ ] Timestamp format validation works
- [ ] Can drag-and-drop steps to reorder
- [ ] Steps renumber after reordering
- [ ] Copy timestamp by clicking time display
- [ ] Shows "Copied!" message

#### UI/UX
- [ ] Video list search works
- [ ] Status filter buttons work
- [ ] Video numbering displays correctly
- [ ] Selected step has visual feedback
- [ ] Insert timestamp adds below selected step
- [ ] Config dialog prevents accidental close
- [ ] Delete button centered and clickable
- [ ] Non-tutorial checkbox works correctly (only checkbox area)

#### Pre-annotation
- [ ] Pre-annotation loads if exists
- [ ] Status badge shows "Pre-annotated"
- [ ] Can edit pre-annotation
- [ ] Saving overwrites pre-annotation

#### Edge Cases
- [ ] Empty video directory handled gracefully
- [ ] Invalid paths show error message
- [ ] Very long video names display correctly
- [ ] Videos with no duration handle properly
- [ ] Can annotate 0:00 timestamp
- [ ] Large step count (100+) renders OK

### Automated Testing

Currently manual testing only. Future versions may include:

- Unit tests for Go backend
- Integration tests for API endpoints
- E2E tests for web interface

### Performance Testing

- [ ] 1000+ videos load in reasonable time (< 5s)
- [ ] Smooth playback with 100+ steps
- [ ] No memory leaks during extended use
- [ ] Search filters 1000+ videos instantly

### Cross-Platform Testing

Test on all target platforms:

- [ ] Windows 10/11 (64-bit)
- [ ] macOS (Intel and Apple Silicon)
- [ ] Linux (Ubuntu, Fedora, etc.)

### Browser Compatibility

Test on multiple browsers:

- [ ] Chrome/Chromium
- [ ] Firefox
- [ ] Safari (macOS)
- [ ] Edge

---

## Keyboard Shortcuts

### Global Shortcuts

| Key | Function | Notes |
|-----|----------|-------|
| `Space` | Play/Pause video | Works anywhere except input fields |
| `←` | Rewind 0.5 seconds | Precise frame-by-frame control |
| `→` | Forward 0.5 seconds | Precise frame-by-frame control |
| `I` | Insert timestamp | Creates new step with current time |

### Player Shortcuts

When video player is focused:

| Key | Function |
|-----|----------|
| `F` | Toggle fullscreen |
| `M` | Mute/unmute |
| `↑` | Volume up |
| `↓` | Volume down |

### Tips

- **Shortcuts disabled in input fields**: Type normally without triggering actions
- **Precise control**: 0.5s arrow keys ideal for frame-accurate annotation
- **Workflow optimization**: Space → I → type → Space for fastest annotation

---

## Configuration

### Configuration File

Location: `~/.mp4label/config.json`

Format:
```json
{
  "video_dir": "/path/to/videos",
  "pre_annotation_dir": "/path/to/pre-annotations",
  "output_dir": "/path/to/output",
  "task_file": "/path/to/task.txt",
  "model_annotation_dir": "/path/to/model-annotations"
}
```

### Path Requirements

- **video_dir**: Must exist and contain `.mp4` files
- **pre_annotation_dir**: Optional, will be ignored if not set
- **output_dir**: Must exist and be writable
- **task_file**: Optional, text file specifying which videos to annotate
- **model_annotation_dir**: Optional, for algorithm engineers to compare model annotations

### Task File Feature (v0.2.4+)

The task file feature allows you to specify a subset of videos to annotate, which is useful for:
- **Collaborative annotation**: Different team members work on different video sets
- **Batch processing**: Annotate videos in multiple rounds
- **Selective annotation**: Only annotate specific videos from a large collection

#### Task File Format

Create a plain text file with one video name per line (without `.mp4` extension):

```
video_name_1
video_name_2
video_name_3
```

**Example** (`task.txt`):
```
tutorial_intro_2024
advanced_effects_part1
beginner_guide_chapter2
```

#### How It Works

1. Set the task file path in configuration
2. Only videos listed in the task file will appear in the video list
3. Non-existent video names are silently ignored
4. Statistics panel shows: Total / Annotated / Pre-annotated / Unannotated counts
5. Leave task file empty to show all videos in the directory

#### Notes

- Empty lines are automatically ignored
- Video names are matched exactly (case-sensitive)
- The `.mp4` extension is automatically removed if present in the file
- Videos not in the task file are completely hidden from the interface

### Model Annotation Comparison (v0.2.5+)

The model annotation feature is designed for **algorithm engineers** to compare model-generated annotations with human ground truth annotations.

#### Purpose

- **Quality Evaluation**: Assess model annotation quality against human annotations
- **Error Analysis**: Identify systematic errors in model predictions
- **Iteration Speed**: Faster feedback loop for model development
- **Visual Comparison**: Side-by-side comparison of model vs human annotations

#### Setup

1. Train your model to output annotations in the same format as human annotations (`.txt` files)
2. Generate model annotations for your video dataset
3. Place model annotation files in a dedicated directory
4. Configure the model annotation directory path in settings

#### Configuration

```json
{
  "video_dir": "/path/to/videos",
  "pre_annotation_dir": "/path/to/pre-annotations",
  "output_dir": "/path/to/output",
  "task_file": "/path/to/task.txt",
  "model_annotation_dir": "/path/to/model-annotations"
}
```

#### How It Works

1. **Without Model Directory**: 3-panel layout (video list, player, annotation editor) - standard annotator view
2. **With Model Directory**: 4-panel layout - adds read-only model annotation panel on the right

#### Model Annotation Format

Model annotations must use the **same format** as human annotations:

**Tutorial Video**:
```
Tutorial Title

1) 00:11.230 Step description
2) 00:25.560 Step description
3) 00:42.890 Step description
```

**Non-Tutorial Video**:
```
[not tutorial]
```

#### Usage for Algorithm Engineers

1. Configure model annotation directory in settings
2. Select a video from the list
3. Compare:
   - **Center-right panel**: Human annotation (editable)
   - **Right-most panel**: Model annotation (read-only, gray background)
4. Identify differences and improvement areas
5. Iterate on model training

#### Notes

- Model panel only appears when `model_annotation_dir` is configured
- Model annotations are **read-only** - cannot be edited
- If model annotation file doesn't exist for a video, panel shows "No model annotation available"
- Annotators don't need to configure this - their workflow remains unchanged
- Model annotations are never overwritten by the application

### Quote Handling

The application automatically handles:
- Single quotes (`'`) in paths
- Double quotes (`"`) in paths
- Leading/trailing whitespace

You can paste paths directly from terminal without manual cleanup.

---

## Troubleshooting

### Video Won't Play

**Symptoms**: Video area blank or shows error

**Solutions**:
- Check video directory path in config
- Verify video file exists and is readable
- Ensure file is MP4 format
- Try refreshing browser (Ctrl+F5 / Cmd+Shift+R)
- Check browser console for errors (F12)

### Annotation Won't Save

**Symptoms**: "Save failed" error

**Solutions**:
- Check output directory path in config
- Verify output directory exists
- Ensure directory has write permissions
- Check annotation format (timestamp format, non-empty descriptions)
- Review browser console for detailed error

### Video List Empty

**Symptoms**: "No videos found" message

**Solutions**:
- Verify video directory contains `.mp4` files
- Check directory path is correct (no typos)
- Ensure directory has read permissions
- Try setting config again

### Keyboard Shortcuts Don't Work

**Symptoms**: Keys have no effect

**Solutions**:
- Click outside input fields first
- Refresh browser page
- Ensure video is loaded
- Check if shortcuts work in other browsers

### Config Dialog Won't Close

**Symptoms**: Clicking outside doesn't close dialog

**This is intentional**: Prevents accidental data loss. Click "Save" or "Cancel" button.

### Timestamp Jump Inaccurate

**Symptoms**: Clicking timestamp seeks to wrong position

**Solutions**:
- Verify timestamp format is correct (`mm:ss.SSS`)
- Check video duration is correct
- Try refreshing page
- Ensure video is fully loaded

### Copy Timestamp Not Working

**Symptoms**: No "Copied!" message appears

**Solutions**:
- Ensure clicking the actual timestamp display (below player)
- Check browser clipboard permissions
- Try a different browser
- Check browser console for errors

### Non-Tutorial Checkbox Area Issue

**Symptoms**: Clicking near checkbox triggers it

**Solutions**:
- Click directly on the checkbox itself
- This has been fixed in v0.2.3

---

## API Endpoints

For developers integrating with mp4Label:

### Video Management

- `GET /` - Main page
- `GET /api/videos` - Get video list
- `GET /api/video/:filename` - Stream video file

### Annotation Management

- `GET /api/annotation/:filename` - Get annotation content
- `POST /api/annotation/:filename` - Save annotation
- `DELETE /api/annotation/:filename` - Delete annotation

### Configuration

- `GET /api/config` - Get current configuration
- `POST /api/config` - Save configuration

### Request/Response Formats

**GET /api/videos**
```json
[
  {
    "filename": "video1.mp4",
    "has_pre_annotation": true,
    "has_annotation": false,
    "index": 0
  }
]
```

**GET /api/annotation/:filename**
```json
{
  "title": "Tutorial Title",
  "is_tutorial": true,
  "steps": [
    {
      "number": 1,
      "timestamp": "00:11.230",
      "description": "Import footage"
    }
  ]
}
```

**POST /api/annotation/:filename**
```json
{
  "title": "Tutorial Title",
  "is_tutorial": true,
  "steps": [...]
}
```

---

## Project Structure

```
mp4Label/
├── cmd/
│   └── mp4label/
│       ├── main.go              # Entry point
│       └── web/                 # Embedded web files (auto-copied)
├── pkg/
│   ├── server/                  # Web server
│   │   └── server.go
│   ├── annotation/              # Annotation processing
│   │   ├── parser.go
│   │   └── validator.go
│   ├── video/                   # Video handling
│   │   └── scanner.go
│   └── config/                  # Configuration management
│       └── config.go
├── web/                         # Source web resources
│   ├── index.html
│   └── static/
│       ├── css/
│       │   └── style.css
│       └── js/
│           └── app.js
├── bin/                         # Compiled binaries
├── Makefile                     # Build automation
├── go.mod                       # Go module definition
├── README.md                    # Project overview
├── DOCS.md                      # This file
├── VERSION_HISTORY.md           # Version history
└── CHANGELOG.md                 # Detailed changelog
```

---

## Tech Stack

- **Backend**: Go 1.22+
- **Frontend**: HTML5 + CSS + JavaScript (vanilla)
- **Video Player**: Video.js 8.10.0
- **Dependencies**: None (Go standard library only)

---

## Contributing

Contributions are welcome! Please:

1. Check existing issues and documentation
2. Fork the repository
3. Create a feature branch
4. Make your changes
5. Test thoroughly
6. Submit a pull request

---

## License

MIT License - see LICENSE file for details

---

## Support

For help, please:

1. Review this documentation
2. Check VERSION_HISTORY.md for changes
3. Search existing issues
4. Submit a new issue with:
   - Your OS and version
   - mp4Label version (`./bin/mp4label version`)
   - Steps to reproduce
   - Expected vs actual behavior

---

**Happy Annotating!** 🎬✨
