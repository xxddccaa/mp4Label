# mp4Label - Video Annotation Tool

mp4Label is a cross-platform video annotation tool built with Go, designed for annotating video tutorial content. It provides a web-based interface for viewing videos, editing annotations, and saving results.

## ✨ Features

- 📦 **Single Executable**: Web files embedded, no external dependencies needed
- 📹 **Video Management**: Scans MP4 files, automatic numbering
- 📝 **Pre-annotation Support**: Load and edit existing pre-annotations
- 🤖 **Model Comparison**: Optional 4th panel to compare model vs human annotations (for algorithm engineers)
- ✏️ **Annotation Editor**: Edit tutorial titles, timestamps, and step descriptions
- 🎬 **Professional Player**: Video.js-based with zoom, fullscreen, speed control
- 💾 **Auto-Save**: Real-time auto-save with visual status indicator
- 👁️ **View Mode**: Read-only `view` subcommand for safe annotation review (CLI config, no editing)
- ⚡ **Quick Annotation**: One-click timestamp insertion with millisecond precision
- 🎯 **Smart Insertion**: Press I to insert below selected step
- 🎯 **Click to Seek**: Click timestamp to jump to video position
- 📋 **Copy Timestamp**: Click timestamp display to copy to clipboard
- 🔄 **Drag-and-Drop**: Reorder steps by dragging
- ⌨️ **Keyboard Shortcuts**: Space to play/pause, arrows to seek (0.5s), I to insert
- 🎛️ **Status Filtering**: Filter videos by annotation status
- 📊 **Task File Support**: Specify a subset of videos to annotate via text file
- 📂 **Task Directory**: Multiple task groups in a directory, switchable via sidebar dropdown
- 📈 **Statistics Display**: Real-time display of total/annotated/pre-annotated/unannotated counts
- ⚙️ **Smart Config**: Handles quoted paths, prevents accidental dialog close
- 💾 **Persistent Config**: Saves to `~/.mp4label/config.json`
- 🌐 **Modern Web UI**: Clean, intuitive interface
- 🔍 **Search Function**: Real-time video filtering
- ✅ **Format Validation**: Automatic validation of timestamps and format

## 🚀 Quick Start

### macOS/Linux

```bash
# 1. Navigate to project directory
cd /path/to/mp4Label

# 2. Build all platforms (default)
make

# 3. Run
./bin/mp4label web
```

### Windows

```bash
# 1. Navigate to project directory
cd /path/to/mp4Label

# 2. Build Windows version
make windows

# 3. Run mp4label.exe (or copy to Windows system)
bin/mp4label.exe web
```

### First-Time Setup

1. Open browser: `http://localhost:8080`
2. Click "Config" button
3. Set paths (type manually or click 📁📄 to open native dialog):
   - **Video Directory**: Folder with MP4 files (required)
   - **Pre-annotation Directory**: Folder with pre-annotations (optional)
   - **Output Directory**: Where to save annotations (required)
   - **Task File**: Text file with video names to annotate (optional, one name per line without .mp4 extension)
   - **Model Annotation Directory**: Folder with model-generated annotations (optional, for algorithm engineers)
4. Click "Save"

### View Mode (Read-Only)

```bash
# Review annotations without editing risk
mp4label view -video-dir /path/to/videos -output-dir /path/to/annotations

# With custom port and all options
mp4label view -port 9090 -video-dir /videos -output-dir /output -task-file /task.txt
```

In view mode: all editing is disabled, config is specified via CLI flags, no config file is read or modified.

## 📚 Documentation

- **[DOCS.md](DOCS.md)** - Complete documentation (user guide, build instructions, testing)
- **[VERSION_HISTORY.md](VERSION_HISTORY.md)** - Detailed version history and release notes
- **[CHANGELOG.md](CHANGELOG.md)** - Detailed changelog
- **[CLAUDE.md](CLAUDE.md)** - AI development guide (for AI-assisted development)

## 🔨 Building

### Build All Platforms (Default)

```bash
make
```

Generates:
- `bin/mp4label.exe` - Windows 64-bit
- `bin/mp4label-linux` - Linux 64-bit
- `bin/mp4label-darwin` - macOS 64-bit

### Build Specific Platform

```bash
make cmd-mp4label  # Current platform only
make windows       # Windows 64-bit
make linux         # Linux 64-bit
make darwin        # macOS 64-bit
```

### Other Commands

```bash
make clean         # Clean build artifacts
make run           # Build and run
make install       # Install to /usr/local/bin/ (macOS/Linux)
```

## ⌨️ Keyboard Shortcuts

| Key | Function |
|-----|----------|
| `Space` | Play/Pause video |
| `←` | Rewind 0.5 seconds |
| `→` | Forward 0.5 seconds |
| `I` | Insert timestamp at current position |

## 🎯 Annotation Workflow

1. **Select Video**: Click video in left sidebar
2. **Enter Title**: Type tutorial title in right editor
3. **Watch & Annotate**:
   - Press `Space` to play
   - Press `I` when reaching key step (video auto-pauses)
   - Type step description
   - Press `Space` to continue
4. **Verify**: Click timestamps to verify accuracy
5. **Save**: Click "Save" button

## 📝 Annotation Format

### Tutorial Video

```
Tutorial Title

1) 00:11.230 Step description
2) 00:25.560 Step description
3) 00:42.890 Step description
```

- Timestamp format: `mm:ss.SSS` (millisecond precision)
- Legacy format `mm:ss` auto-converts to `mm:ss.000`

### Non-Tutorial Video

```
[not tutorial]
```

## 🏗️ Project Structure

```
mp4Label/
├── cmd/mp4label/        # Main entry point
├── pkg/                 # Core packages
│   ├── server/          # Web server
│   ├── annotation/      # Annotation processing
│   ├── video/           # Video handling
│   └── config/          # Configuration
├── web/                 # Frontend resources
│   ├── index.html
│   └── static/
│       ├── css/
│       └── js/
├── bin/                 # Compiled binaries
├── Makefile             # Build automation
├── README.md            # This file
├── DOCS.md              # Complete documentation
├── VERSION_HISTORY.md   # Version history
├── CHANGELOG.md         # Detailed changelog
└── CLAUDE.md            # AI development guide
```

## 🔧 API Endpoints

- `GET /` - Main page
- `GET /api/videos` - Get video list
- `GET /api/video/:filename` - Stream video file
- `GET /api/annotation/:filename` - Get annotation
- `POST /api/annotation/:filename` - Save annotation
- `DELETE /api/annotation/:filename` - Delete annotation
- `GET /api/config` - Get configuration
- `POST /api/config` - Save configuration
- `GET /api/mode` - Get server mode (readonly/edit)
- `GET /api/task-groups` - List available task groups
- `GET /api/dialog` - Open native file/folder dialog

## 💻 Tech Stack

- **Backend**: Go 1.22+
- **Frontend**: HTML5 + CSS + JavaScript (vanilla)
- **Video Player**: Video.js 8.10.0
- **Dependencies**: None (Go standard library only)

## 📦 Single-File Distribution

✨ All builds have web files embedded

- ✅ Only the executable needed
- ✅ No external files required
- ✅ Direct user distribution
- ✅ Cross-platform compatible

## 🐛 Troubleshooting

See [DOCS.md](DOCS.md#troubleshooting) for detailed troubleshooting guide.

Common issues:
- **Video won't play**: Check video directory path and file format
- **Annotation won't save**: Verify output directory exists and is writable
- **Empty video list**: Ensure directory contains `.mp4` files
- **Keyboard shortcuts don't work**: Click outside input fields first

## 📄 License

MIT License

## 🤝 Contributing

Contributions welcome! Please:

1. Check existing issues and documentation
2. Fork the repository
3. Create a feature branch
4. Test thoroughly
5. Submit a pull request

## 📌 Version

Current version: **v0.2.9** (2026-02-11)

Recent updates:
- **v0.2.9** - Task directory with switchable task groups for team annotation
- **v0.2.8** - Read-only view mode with CLI config for safe annotation review
- **v0.2.7** - Auto-save annotations, fix 1ms timestamp precision issue
- **v0.2.6** - Native OS file/folder dialogs for configuration
- **v0.2.5** - Model annotation comparison panel for algorithm engineers

See [VERSION_HISTORY.md](VERSION_HISTORY.md) for complete history.

---

**Happy Annotating!** 🎬✨
