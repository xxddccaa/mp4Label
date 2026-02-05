# Changelog

All notable changes to this project will be documented in this file.

## [v0.2.6] - 2026-02-05

### New Features

#### 🗂️ Native OS File Dialog
- Added native file/folder picker accessed via 📁📄 buttons
- Backend launches OS dialogs and returns full absolute paths
- Supports both directory selection (4 paths) and file selection (1 path)
- One-click selection fills the input field immediately
- Removes browser security limitations for path selection

### Technical Changes

#### Backend (Go)
- Added `GET /api/dialog` endpoint in `server.go`
- macOS: uses `osascript` native dialogs
- Windows: uses PowerShell dialog APIs
- Linux: uses `zenity` for file/folder selection

#### Frontend
- Reused existing 📁📄 buttons to trigger native dialogs
- `openBrowser(inputId, mode)` calls backend dialog API and sets the input path
- Removed custom browser modal UI

### User Experience
- Click 📁/📄 buttons to open native OS picker
- Select folder/file and get full path automatically
- Reduces configuration errors from manual path entry
- Faster setup and clearer workflow

---

## [v0.2.5] - 2026-02-05

### New Features

#### 🤖 Model Annotation Comparison Panel
- Added optional "Model Annotation Directory" configuration for algorithm engineers
- New 4th panel (right-most) displays model-generated annotations in read-only mode
- Dynamic layout: panel only appears when model annotation directory is configured
- Allows side-by-side comparison between human annotations and model predictions
- Helps algorithm teams evaluate model performance and identify improvement areas

### Technical Changes
- Backend: Added `ModelAnnotationDir` field to `Config` struct
- Backend: New API endpoint `GET /api/model-annotation/:filename` for fetching model annotations
- Frontend: Added model panel with read-only display of model annotations
- Frontend: Dynamic 4-panel layout that adapts based on configuration
- Frontend: Visual distinction (gray background) indicates read-only model annotations

### Configuration
- New optional field: `model_annotation_dir` in config
- Backward compatible: existing users see no changes unless they configure this field
- Model annotation files use same format as human annotations (`.txt` files)

### User Experience
- Algorithm engineers can now compare model output with ground truth annotations
- Annotators continue using 3-panel layout (unchanged experience)
- Seamless switching between 3-panel and 4-panel modes based on config

---

## [v0.2.4] - 2026-02-04

### New Features

#### 📊 Task File Support
- Added optional "Task File" configuration to specify a subset of videos to annotate
- Useful for collaborative annotation where different users work on different video sets
- Task file format: one video name per line (without .mp4 extension)
- Automatically filters video list to show only videos listed in task file
- Non-existent video names in task file are silently ignored

#### 📈 Statistics Display
- Added real-time statistics panel in left sidebar showing:
  - **Total**: Total number of videos (filtered by task file if set)
  - **Annotated**: Number of videos with completed annotations
  - **Pre-annotated**: Number of videos with pre-annotations only
  - **Unannotated**: Number of videos without any annotation
- Statistics update automatically when annotations are saved or deleted

### Technical Changes
- Backend: Modified `ScanVideos` function to accept optional task file parameter
- Backend: Added `loadTaskFile` function to parse task file content
- Backend: Enhanced video API response to include statistics
- Frontend: Updated config dialog to include task file input
- Frontend: Added statistics display panel with grid layout
- Frontend: Enhanced loadVideos function to handle new response format

### Documentation
- Added `TASK_FILE_GUIDE.md` with detailed usage instructions and examples
- Updated README.md to mention new features
- Added task file configuration in setup instructions

## [v0.2.3] - 2026-02-04

### New Features

#### ⌨️ Keyboard Shortcut Optimization
- Changed arrow key timing from 5s to 0.5s for more precise control
- Left arrow (←): Rewind 0.5 seconds (was 5s)
- Right arrow (→): Forward 0.5 seconds (was 5s)
- Enables frame-by-frame annotation precision

#### 📋 One-Click Timestamp Copy
- Click timestamp display to copy to clipboard
- Shows "Copied!" floating message after clicking
- Auto-dismiss after 2 seconds
- Works on modern and legacy browsers (with fallback)

### Bug Fixes

#### 🐛 Non-Tutorial Checkbox Click Area
- Fixed: Only clicking checkbox itself marks as "non-tutorial"
- Before: Clicking label area also triggered checkbox
- After: Must click the actual checkbox input element
- Improved user control and prevented accidental clicks

### Documentation

#### 📚 Documentation Consolidation
- Created VERSION_HISTORY.md with comprehensive version history
- Created DOCS.md with unified documentation
- Internationalized all documentation to English
- Removed redundant documentation files
- Cleaned up project structure

### Build System

#### 🔨 Multi-Platform Build Enhancement
- Changed default `make` target to build all platforms
- Single command now generates Windows, Linux, and macOS binaries
- Simplified build workflow for distribution
- Added build summary output

### Internationalization

#### 🌐 UI and Documentation
- All HTML interface text now in English
- All code comments translated to English
- All documentation translated to English
- Better for international collaboration

---

## [v0.2.2] - 2026-02-03 (Previous Version)

### 重要变更

#### 🔧 命令行接口重构
- **子命令方式**：采用 `mp4label <子命令>` 结构
- **启动服务器**：`mp4label web` 替代原来的 `mp4label`
- **版本显示**：新增 `mp4label version` 命令
- **帮助信息**：新增 `mp4label help` 命令
- **更好的用户体验**：命令意图更清晰，符合现代工具规范

**迁移说明：**
```bash
# 旧方式（已废弃）
./mp4label -port 8080

# 新方式
./mp4label web -port 8080
```

### 用户体验优化

#### ⌨️ 全局键盘快捷键
- **改进**：键盘快捷键现在在整个页面都能工作
- 无需先点击视频播放器即可使用左右箭头键
- 智能判断：在输入框中不会触发快捷键
- 支持快捷键：空格、←、→、I

#### 🎨 视频控制栏布局优化
- **改进**：控制栏改为垂直布局，避免元素重叠
- 第一行：视频名称 + 当前时间
- 第二行：操作按钮（插入时间戳、播放速度）
- 长视频名称自动截断显示（...省略号）
- 时间显示固定宽度，不会被挤压

#### ⏱️ 视频时长显示优化
- **改进**：优化 Video.js 时间显示配置
- 进度条后清晰显示：当前时间 / 总时长
- 增大时间显示字体，更清晰可读
- 优化元素间距，布局更合理

---

## [v0.2.2] - 2026-02-03

### 新增功能

#### 📦 Web 文件嵌入
- **单文件分发**：使用 Go embed 将 web 目录完整嵌入到可执行文件
- **即开即用**：无需额外文件，双击 exe 即可运行
- **跨平台支持**：Windows、macOS、Linux 都支持单文件运行
- **简化部署**：只需复制一个文件即可

#### 🎯 步骤选中状态
- **视觉反馈**：选中的步骤有明显的蓝色边框和背景
- **智能插入**：按 I 键在选中步骤下方插入新步骤
- **自动选中**：新插入的步骤自动选中并聚焦
- **默认选中**：加载标注时自动选中第一个步骤
- **删除管理**：删除步骤时自动调整选中状态

#### ⚙️ 配置交互改进
- **防止误关闭**：点击配置对话框外部会晃动提示，不会关闭
- **引号支持**：自动识别和清理路径中的单引号和双引号
- **空格处理**：自动去除路径首尾的空格
- **更好的 UX**：用户必须点击"保存"或"取消"才能关闭配置

#### 🔢 视频编号系统
- 为每个视频分配固定编号（#0001 - #9999）
- 编号按视频列表顺序分配
- 方便识别和引用视频

#### 🎛️ 状态筛选功能
- 新增筛选按钮：全部、未标注、预标注、已标注
- 快速过滤视频列表
- 支持与搜索功能组合使用

### 改进

#### ⏱️ 视频时间显示优化
- 显示格式从"剩余时间"改为"已播放/总时长"
- 更直观的时间信息

#### 📐 布局再次优化
- 右侧编辑器宽度：480px → **550px**
- 提供更宽敞的编辑空间

#### 🐛 Bug 修复
- 修复点击时间戳跳转后自动播放的问题
- 现在跳转后保持暂停状态
- 修复删除按钮位置不居中的问题

---

## [v0.2.1] - 2026-02-03

### 新增功能

#### ⏱️ 毫秒精度时间戳
- 时间戳格式从 `mm:ss` 升级到 `mm:ss.SSS`
- 支持毫秒级精度（例如：12:32.766）
- 旧格式标注文件自动转换（添加 .000）
- 插入时间戳自动包含毫秒

#### 🔄 步骤拖拽排序
- 支持拖拽步骤卡片重新排序
- 拖拽手柄（⋮⋮）提示
- 拖拽时视觉反馈
- 自动重新编号

#### 📐 布局优化
- 右侧编辑器宽度从 400px 增加到 480px
- 更适合手机竖屏视频标注
- 时间戳输入框宽度从 80px 增加到 100px

#### 🐛 Bug 修复
- 修复时间戳跳转位置不正确的问题
- 改进时间戳解析逻辑，支持毫秒精度

---

## [v0.2.0] - 2026-02-03

### 新增功能

#### 🎬 专业视频播放器
- 使用 Video.js 替换原生 HTML5 video 标签
- 支持视频缩放、全屏、画中画等高级功能
- 更好的跨浏览器兼容性和性能
- 自定义控制栏，更符合标注工作流

#### ⚡ 快速标注功能
- **插入时间戳按钮**：一键插入当前播放时间到新步骤
- **自动暂停**：插入时间戳后视频自动暂停，方便输入描述
- **自动聚焦**：焦点自动移到描述输入框，提高输入效率
- **键盘快捷键**：按 I 键快速插入时间戳

#### 🎯 点击跳转功能
- 点击步骤列表中的时间戳，视频自动跳转到对应位置
- 跳转后自动播放 1 秒并暂停，方便快速查看
- 时间戳添加视觉反馈（悬停高亮）
- 便于验证和调整标注内容

#### ⌨️ 快捷键支持
- `空格键`：播放/暂停视频
- `←`（左箭头）：后退 5 秒
- `→`（右箭头）：前进 5 秒
- `I`：插入当前时间戳到新步骤
- `F`：全屏（在播放器焦点时）

#### 🎚️ 播放速度控制
- 新增播放速度选择器
- 支持 0.5x、0.75x、1x、1.25x、1.5x、2x 倍速
- 方便快速浏览或精确标注

### 改进

#### 界面优化
- 重新设计视频控制栏布局
- 改进时间显示样式，更加醒目
- 优化视频容器尺寸，更好利用空间
- 时间戳输入框添加点击提示和悬停效果

#### 用户体验
- 标注工作流更加流畅
- 减少鼠标操作，提高标注效率
- 改善视频播放控制的精确度
- 增强视频和标注内容的联动

### 文档

#### 新增文档
- **USER_GUIDE.md**：详细的用户使用指南
  - 快速开始教程
  - 完整的标注工作流说明
  - 功能详解和使用技巧
  - 常见问题解答

#### 更新文档
- **README.md**：更新功能列表和使用说明
- **REQUIREMENTS.md**：补充新功能需求和工作流设计

### 技术变更

- 引入 Video.js 8.10.0（通过 CDN）
- 重构视频播放器初始化逻辑
- 新增时间戳跳转函数 `seekToTimestamp()`
- 新增快速插入函数 `insertCurrentTimestamp()`
- 改进事件监听器组织结构

---

## [v0.1.0] - 2026-02-02

### 初始功能

- 视频文件扫描和列表显示
- 基础视频播放功能
- 标注编辑器（教程题目、步骤列表）
- 预标注文件读取
- 标注保存功能
- 配置管理（视频目录、预标注目录、输出目录）
- 基础搜索功能
- Web 界面实现

### 技术栈

- 后端：Go 1.22+
- 前端：HTML5 + CSS + JavaScript（原生）
- 无外部依赖（仅 Go 标准库）
