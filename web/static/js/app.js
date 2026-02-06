// 全局状态
let currentVideo = null;
let currentAnnotation = null;
let currentModelAnnotation = null; // 模型标注数据
let videos = [];
let videoStats = null; // 视频统计信息
let config = null;
let player = null; // Video.js player 实例
let selectedStepIndex = -1; // 当前选中的步骤索引，-1 表示没有选中
let currentFilter = 'all'; // 当前筛选状态
let autoSaveTimer = null; // 自动保存定时器
let lastSavedAnnotationJSON = null; // 上次保存的标注JSON，用于检测变化

const AUTO_SAVE_DELAY = 1500; // 自动保存延迟（毫秒）

// DOM 元素
const videoList = document.getElementById('videoList');
const searchInput = document.getElementById('searchInput');
const videoPlayer = document.getElementById('videoPlayer');
const currentVideoName = document.getElementById('currentVideoName');
const currentTime = document.getElementById('currentTime');
const tutorialTitle = document.getElementById('tutorialTitle');
const stepsContainer = document.getElementById('stepsContainer');
const addStepBtn = document.getElementById('addStepBtn');
const saveBtn = document.getElementById('saveBtn');
const deleteAnnotationBtn = document.getElementById('deleteAnnotationBtn');
const notTutorialCheck = document.getElementById('notTutorialCheck');
const notTutorialGroup = document.getElementById('notTutorialGroup');
const configBtn = document.getElementById('configBtn');
const configModal = document.getElementById('configModal');
const closeModal = document.querySelector('.close');
const saveConfigBtn = document.getElementById('saveConfigBtn');
const cancelConfigBtn = document.getElementById('cancelConfigBtn');
const insertTimestampBtn = document.getElementById('insertTimestampBtn');
const playbackRate = document.getElementById('playbackRate');


// 初始化
document.addEventListener('DOMContentLoaded', () => {
    initVideoPlayer();
    loadConfig();
    loadVideos();
    setupEventListeners();
    setupResizableHandles();
});

// 初始化 Video.js 播放器
function initVideoPlayer() {
    player = videojs('videoPlayer', {
        controls: true,
        fluid: false,
        aspectRatio: '16:9',
        playbackRates: [0.5, 0.75, 1, 1.25, 1.5, 2],
        controlBar: {
            children: [
                'playToggle',
                'volumePanel',
                'progressControl',
                'currentTimeDisplay',
                'timeDivider',
                'durationDisplay',
                'playbackRateMenuButton',
                'fullscreenToggle'
            ]
        }
    });

    // 监听播放器时间更新
    player.on('timeupdate', () => {
        updateTimeDisplay();
    });

    // 添加时间显示点击复制功能
    currentTime.style.cursor = 'pointer';
    currentTime.title = 'Click to copy timestamp';
    currentTime.addEventListener('click', copyCurrentTimestamp);

    // 添加键盘快捷键
    player.on('keydown', (e) => {
        // 空格键：播放/暂停
        if (e.keyCode === 32) {
            e.preventDefault();
            if (player.paused()) {
                player.play();
            } else {
                player.pause();
            }
        }
        // 左箭头：后退0.5秒
        else if (e.keyCode === 37) {
            e.preventDefault();
            player.currentTime(Math.max(0, player.currentTime() - 0.5));
        }
        // 右箭头：前进0.5秒
        else if (e.keyCode === 39) {
            e.preventDefault();
            player.currentTime(Math.min(player.duration(), player.currentTime() + 0.5));
        }
        // I 键：插入时间戳
        else if (e.keyCode === 73) {
            e.preventDefault();
            insertCurrentTimestamp();
        }
    });
}

// 设置事件监听
function setupEventListeners() {
    // 搜索
    searchInput.addEventListener('input', filterVideos);

    // 添加步骤
    addStepBtn.addEventListener('click', addStep);

    // 保存标注
    saveBtn.addEventListener('click', saveAnnotation);

    // 标题输入触发自动保存
    tutorialTitle.addEventListener('input', () => {
        scheduleAutoSave();
    });

    // 删除标注
    deleteAnnotationBtn.addEventListener('click', deleteAnnotation);

    // 插入时间戳按钮
    insertTimestampBtn.addEventListener('click', insertCurrentTimestamp);

    // 播放速度控制
    playbackRate.addEventListener('change', (e) => {
        if (player) {
            player.playbackRate(parseFloat(e.target.value));
        }
    });

    // 筛选按钮
    document.querySelectorAll('.filter-btn').forEach(btn => {
        btn.addEventListener('click', (e) => {
            // 更新按钮状态
            document.querySelectorAll('.filter-btn').forEach(b => b.classList.remove('active'));
            e.target.classList.add('active');
            
            // 更新筛选状态
            currentFilter = e.target.dataset.filter;
            
            // 重新渲染列表
            renderVideoList();
        });
    });

    // 全局键盘快捷键
    document.addEventListener('keydown', (e) => {
        // 如果焦点在输入框、文本域或可编辑元素，不处理快捷键
        const target = e.target;
        if (target.tagName === 'INPUT' || target.tagName === 'TEXTAREA' || target.isContentEditable) {
            return;
        }

        if (!player) return;

        switch(e.keyCode) {
            case 32: // 空格键：播放/暂停
                e.preventDefault();
                if (player.paused()) {
                    player.play();
                } else {
                    player.pause();
                }
                break;
            case 37: // 左箭头：后退0.5秒
                e.preventDefault();
                player.currentTime(Math.max(0, player.currentTime() - 0.5));
                break;
            case 39: // 右箭头：前进0.5秒
                e.preventDefault();
                player.currentTime(Math.min(player.duration(), player.currentTime() + 0.5));
                break;
            case 73: // I键：插入时间戳
                e.preventDefault();
                insertCurrentTimestamp();
                break;
        }
    });

    // 非教学视频复选框 - 只响应checkbox自身的change事件
    notTutorialCheck.addEventListener('change', (e) => {
        if (e.target.checked) {
            currentAnnotation = {
                title: '',
                is_tutorial: false,
                steps: []
            };
            renderEditor();
        } else {
            currentAnnotation = {
                title: '',
                is_tutorial: true,
                steps: []
            };
            renderEditor();
        }
        scheduleAutoSave();
    });

    // 阻止点击label区域触发checkbox
    notTutorialGroup.addEventListener('click', (e) => {
        // 只允许点击checkbox本身才触发
        if (e.target !== notTutorialCheck) {
            e.preventDefault();
            e.stopPropagation();
        }
    });

    // 配置对话框
    configBtn.addEventListener('click', () => {
        configModal.style.display = 'block';
        loadConfigToForm();
    });

    closeModal.addEventListener('click', () => {
        configModal.style.display = 'none';
    });

    cancelConfigBtn.addEventListener('click', () => {
        configModal.style.display = 'none';
    });

    saveConfigBtn.addEventListener('click', saveConfig);

    // 禁止点击模态框外部关闭（用户必须点击保存或取消）
    window.addEventListener('click', (e) => {
        if (e.target === configModal) {
            // 添加一个晃动动画提示用户必须点击保存或取消
            configModal.querySelector('.modal-content').classList.add('shake');
            setTimeout(() => {
                configModal.querySelector('.modal-content').classList.remove('shake');
            }, 500);
        }
    });
}

// 加载配置
async function loadConfig() {
    try {
        const response = await fetch('/api/config');
        config = await response.json();
        updateModelPanelVisibility();
    } catch (error) {
        console.error('加载配置失败:', error);
    }
}

// 更新模型面板的显示状态
function updateModelPanelVisibility() {
    const modelPanel = document.getElementById('modelPanel');
    if (config && config.model_annotation_dir && config.model_annotation_dir !== '') {
        modelPanel.classList.add('visible');
    } else {
        modelPanel.classList.remove('visible');
    }
}

// 加载配置到表单
function loadConfigToForm() {
    document.getElementById('videoDir').value = config?.video_dir || '';
    document.getElementById('preAnnotationDir').value = config?.pre_annotation_dir || '';
    document.getElementById('outputDir').value = config?.output_dir || '';
    document.getElementById('taskFile').value = config?.task_file || '';
    document.getElementById('modelAnnotationDir').value = config?.model_annotation_dir || '';
}

// 保存配置
async function saveConfig() {
    const videoDir = document.getElementById('videoDir').value.trim();
    const preAnnotationDir = document.getElementById('preAnnotationDir').value.trim();
    const outputDir = document.getElementById('outputDir').value.trim();
    const taskFile = document.getElementById('taskFile').value.trim();
    const modelAnnotationDir = document.getElementById('modelAnnotationDir').value.trim();

    if (!videoDir || !outputDir) {
        alert('Video directory and output directory cannot be empty');
        return;
    }

    try {
        const response = await fetch('/api/config', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                video_dir: videoDir,
                pre_annotation_dir: preAnnotationDir || '',
                output_dir: outputDir,
                task_file: taskFile || '',
                model_annotation_dir: modelAnnotationDir || ''
            })
        });

        if (response.ok) {
            configModal.style.display = 'none';
            await loadConfig();
            updateModelPanelVisibility();
            loadVideos();
            alert('Configuration saved successfully');
        } else {
            const error = await response.text();
            alert('Failed to save configuration: ' + error);
        }
    } catch (error) {
        console.error('Failed to save configuration:', error);
        alert('Failed to save configuration: ' + error.message);
    }
}

// 加载视频列表
async function loadVideos() {
    try {
        const response = await fetch('/api/videos');
        const data = await response.json();
        
        // 处理新的响应格式
        if (data.videos) {
            videos = data.videos;
            videoStats = data.stats;
        } else {
            // 兼容旧格式
            videos = data;
            videoStats = null;
        }
        
        // 为每个视频添加索引（用于编号）
        videos.forEach((video, index) => {
            video.index = index;
        });
        
        renderVideoList();
        updateStatsDisplay();
    } catch (error) {
        console.error('Failed to load video list:', error);
        videoList.innerHTML = '<div class="loading">Failed to load videos, please check configuration</div>';
    }
}

// 更新统计信息显示
function updateStatsDisplay() {
    const statsPanel = document.getElementById('statsPanel');
    if (!videoStats) {
        statsPanel.style.display = 'none';
        return;
    }
    
    document.getElementById('statsTotal').textContent = videoStats.total || 0;
    document.getElementById('statsAnnotated').textContent = videoStats.annotated || 0;
    document.getElementById('statsPreAnnotated').textContent = videoStats.pre_annotated || 0;
    document.getElementById('statsUnannotated').textContent = videoStats.unannotated || 0;
    
    statsPanel.style.display = 'block';
}

// 渲染视频列表
function renderVideoList() {
    if (videos.length === 0) {
        videoList.innerHTML = '<div class="loading">No video files found</div>';
        return;
    }

    const filteredVideos = filterVideos();
    
    if (filteredVideos.length === 0) {
        videoList.innerHTML = '<div class="loading">No videos match the filter</div>';
        return;
    }
    
    videoList.innerHTML = filteredVideos.map((video, index) => {
        const statusBadges = [];
        if (video.has_pre_annotation) {
            statusBadges.push('<span class="status-badge pre-annotation">预标注</span>');
        }
        if (video.has_annotation) {
            statusBadges.push('<span class="status-badge annotation">已标注</span>');
        }
        if (!video.has_pre_annotation && !video.has_annotation) {
            statusBadges.push('<span class="status-badge none">未标注</span>');
        }

        // 使用视频在原始列表中的索引作为编号
        const videoNumber = String(video.index + 1).padStart(4, '0');
        // 保持当前选中视频的高亮状态
        const isActive = video.filename === currentVideo;

        return `
            <div class="video-item${isActive ? ' active' : ''}" data-filename="${video.filename}">
                <div class="video-item-name">
                    <span class="video-number">#${videoNumber}</span>
                    <span>${video.filename}</span>
                </div>
                <div class="video-item-status">${statusBadges.join('')}</div>
            </div>
        `;
    }).join('');

    // 添加点击事件
    document.querySelectorAll('.video-item').forEach(item => {
        item.addEventListener('click', () => {
            const filename = item.dataset.filename;
            selectVideo(filename);
        });
    });
}

// 筛选视频
function filterVideos() {
    const searchTerm = searchInput.value.toLowerCase();
    
    let filtered = videos.filter(video => {
        // 搜索筛选
        if (searchTerm && !video.filename.toLowerCase().includes(searchTerm)) {
            return false;
        }
        
        // 状态筛选
        if (currentFilter === 'all') {
            return true;
        } else if (currentFilter === 'none') {
            // 未标注：没有预标注也没有已标注
            return !video.has_pre_annotation && !video.has_annotation;
        } else if (currentFilter === 'pre') {
            // 预标注：有预标注但没有已标注
            return video.has_pre_annotation && !video.has_annotation;
        } else if (currentFilter === 'done') {
            // 已标注：有已标注
            return video.has_annotation;
        }
        
        return true;
    });

    return filtered;
}

// 选择视频
async function selectVideo(filename) {
    // 切换视频前，立即保存当前标注（如果有未保存的更改）
    if (currentVideo && currentVideo !== filename) {
        await flushAutoSave();
    }

    currentVideo = filename;

    // 更新列表选中状态
    document.querySelectorAll('.video-item').forEach(item => {
        item.classList.remove('active');
        if (item.dataset.filename === filename) {
            item.classList.add('active');
        }
    });

    // 加载视频到 Video.js 播放器
    if (player) {
        player.src({
            type: 'video/mp4',
            src: `/api/video/${encodeURIComponent(filename)}`
        });
        player.load();
    }
    currentVideoName.textContent = filename;

    // 加载标注
    await loadAnnotation(filename);
    
    // 加载模型标注（如果已配置）
    if (config && config.model_annotation_dir && config.model_annotation_dir !== '') {
        await loadModelAnnotation(filename);
    }
}

// 加载标注
async function loadAnnotation(filename) {
    try {
        const stem = filename.replace(/\.mp4$/, '');
        const response = await fetch(`/api/annotation/${stem}.txt`);
        currentAnnotation = await response.json();
        renderEditor();
    } catch (error) {
        console.error('Failed to load annotation:', error);
        currentAnnotation = {
            title: '',
            is_tutorial: true,
            steps: []
        };
        renderEditor();
    }
    // 重置自动保存状态，加载后的数据视为已保存
    lastSavedAnnotationJSON = JSON.stringify(currentAnnotation);
    updateAutoSaveStatus('');
}

// 渲染编辑器
function renderEditor() {
    if (!currentAnnotation) {
        return;
    }

    // 非教学视频
    if (!currentAnnotation.is_tutorial) {
        notTutorialGroup.style.display = 'block';
        notTutorialCheck.checked = true;
        tutorialTitle.value = '';
        tutorialTitle.disabled = true;
        stepsContainer.innerHTML = '<p style="color: #666; padding: 1rem;">Non-tutorial video</p>';
        selectedStepIndex = -1;
        return;
    }

    // 教学视频
    notTutorialGroup.style.display = 'block';
    notTutorialCheck.checked = false;
    tutorialTitle.disabled = false;
    tutorialTitle.value = currentAnnotation.title || '';

    // 渲染步骤
    stepsContainer.innerHTML = '';
    if (currentAnnotation.steps && currentAnnotation.steps.length > 0) {
        currentAnnotation.steps.forEach((step, index) => {
            addStepElement(step, index);
        });
        // 如果之前有选中的步骤，保持选中状态
        if (selectedStepIndex >= 0 && selectedStepIndex < currentAnnotation.steps.length) {
            selectStep(selectedStepIndex);
        } else if (currentAnnotation.steps.length > 0) {
            // 默认选中第一个步骤
            selectStep(0);
        }
    } else {
        selectedStepIndex = -1;
    }
}

// 添加步骤
function addStep() {
    if (!currentAnnotation || !currentAnnotation.is_tutorial) {
        return;
    }

    const stepNumber = (currentAnnotation.steps?.length || 0) + 1;
    const newStep = {
        number: stepNumber,
        timestamp: '00:00.000',
        description: ''
    };

    if (!currentAnnotation.steps) {
        currentAnnotation.steps = [];
    }
    currentAnnotation.steps.push(newStep);
    addStepElement(newStep, currentAnnotation.steps.length - 1);
    scheduleAutoSave();
}

// 添加步骤元素
// 选中步骤
function selectStep(index) {
    selectedStepIndex = index;
    
    // 更新所有步骤的选中状态
    document.querySelectorAll('.step-item').forEach((item, idx) => {
        if (idx === index) {
            item.classList.add('selected');
        } else {
            item.classList.remove('selected');
        }
    });
}

function addStepElement(step, index) {
    const stepDiv = document.createElement('div');
    stepDiv.className = 'step-item';
    stepDiv.draggable = true;
    stepDiv.dataset.index = index;
    stepDiv.innerHTML = `
        <div class="step-header">
            <span class="step-drag-handle" title="Drag to reorder">⋮⋮</span>
            <span class="step-number">Step ${step.number}</span>
            <input type="text" class="step-timestamp clickable" value="${step.timestamp}" 
                   placeholder="00:00.000" data-index="${index}" title="Click to seek to this time">
        </div>
        <textarea class="step-description" data-index="${index}" 
                  placeholder="Enter step description...">${step.description || ''}</textarea>
        <button class="step-remove" data-index="${index}">×</button>
    `;

    // 点击步骤容器选中该步骤
    stepDiv.addEventListener('click', (e) => {
        // 如果点击的是删除按钮或拖动手柄，不选中
        if (!e.target.classList.contains('step-remove') && 
            !e.target.classList.contains('step-drag-handle')) {
            selectStep(index);
        }
    });

    // 拖拽事件
    stepDiv.addEventListener('dragstart', handleDragStart);
    stepDiv.addEventListener('dragover', handleDragOver);
    stepDiv.addEventListener('drop', handleDrop);
    stepDiv.addEventListener('dragend', handleDragEnd);

    // 时间戳输入
    const timestampInput = stepDiv.querySelector('.step-timestamp');
    
    // 点击时间戳跳转到视频对应位置
    timestampInput.addEventListener('click', (e) => {
        e.stopPropagation(); // 防止触发步骤选中
        const timestamp = e.target.value;
        seekToTimestamp(timestamp);
    });
    
    timestampInput.addEventListener('change', (e) => {
        const idx = parseInt(e.target.dataset.index);
        if (currentAnnotation.steps[idx]) {
            currentAnnotation.steps[idx].timestamp = e.target.value;
        }
        scheduleAutoSave();
    });

    // 时间戳实时输入同步（包括粘贴）
    timestampInput.addEventListener('input', (e) => {
        const idx = parseInt(e.target.dataset.index);
        if (currentAnnotation.steps[idx]) {
            currentAnnotation.steps[idx].timestamp = e.target.value;
        }
        scheduleAutoSave();
    });

    // 阻止时间戳输入框的点击事件冒泡
    timestampInput.addEventListener('focus', (e) => {
        e.stopPropagation();
    });

    // 描述输入
    const descriptionInput = stepDiv.querySelector('.step-description');
    descriptionInput.addEventListener('input', (e) => {
        const idx = parseInt(e.target.dataset.index);
        if (currentAnnotation.steps[idx]) {
            currentAnnotation.steps[idx].description = e.target.value;
        }
        scheduleAutoSave();
    });

    // 阻止描述输入框的点击事件冒泡
    descriptionInput.addEventListener('focus', (e) => {
        e.stopPropagation();
    });

    // 删除按钮
    const removeBtn = stepDiv.querySelector('.step-remove');
    removeBtn.addEventListener('click', (e) => {
        e.stopPropagation(); // 防止触发步骤选中
        const idx = parseInt(removeBtn.dataset.index);
        currentAnnotation.steps.splice(idx, 1);
        // 重新编号
        currentAnnotation.steps.forEach((s, i) => {
            s.number = i + 1;
        });
        // 调整选中索引
        if (selectedStepIndex === idx) {
            // 如果删除的是当前选中的步骤
            if (currentAnnotation.steps.length > 0) {
                selectedStepIndex = Math.min(idx, currentAnnotation.steps.length - 1);
            } else {
                selectedStepIndex = -1;
            }
        } else if (selectedStepIndex > idx) {
            // 如果删除的步骤在选中步骤之前
            selectedStepIndex--;
        }
        renderEditor();
        scheduleAutoSave();
    });

    stepsContainer.appendChild(stepDiv);
}

// 拖拽相关变量
let draggedElement = null;
let draggedIndex = null;

// 拖拽开始
function handleDragStart(e) {
    draggedElement = this;
    draggedIndex = parseInt(this.dataset.index);
    this.style.opacity = '0.4';
    e.dataTransfer.effectAllowed = 'move';
}

// 拖拽经过
function handleDragOver(e) {
    if (e.preventDefault) {
        e.preventDefault();
    }
    e.dataTransfer.dropEffect = 'move';
    
    // 添加视觉反馈
    if (this !== draggedElement) {
        this.style.borderTop = '2px solid #2196f3';
    }
    
    return false;
}

// 放下
function handleDrop(e) {
    if (e.stopPropagation) {
        e.stopPropagation();
    }
    
    if (draggedElement !== this) {
        const targetIndex = parseInt(this.dataset.index);
        
        // 交换步骤
        const draggedStep = currentAnnotation.steps[draggedIndex];
        currentAnnotation.steps.splice(draggedIndex, 1);
        currentAnnotation.steps.splice(targetIndex, 0, draggedStep);
        
        // 重新编号
        currentAnnotation.steps.forEach((s, i) => {
            s.number = i + 1;
        });
        
        // 重新渲染
        renderEditor();
        scheduleAutoSave();
    }
    
    return false;
}

// 拖拽结束
function handleDragEnd(e) {
    this.style.opacity = '1';
    
    // 移除所有视觉反馈
    document.querySelectorAll('.step-item').forEach(item => {
        item.style.borderTop = '';
    });
}

// 保存标注（手动保存）
async function saveAnnotation() {
    // 取消待执行的自动保存
    if (autoSaveTimer) {
        clearTimeout(autoSaveTimer);
        autoSaveTimer = null;
    }

    if (!currentVideo) {
        alert('Please select a video first');
        return;
    }

    if (!currentAnnotation) {
        alert('No annotation to save');
        return;
    }

    // 更新标注数据
    if (currentAnnotation.is_tutorial) {
        currentAnnotation.title = tutorialTitle.value.trim();
    }

    // 验证
    if (currentAnnotation.is_tutorial) {
        if (!currentAnnotation.title) {
            alert('Please enter tutorial title');
            return;
        }
        if (!currentAnnotation.steps || currentAnnotation.steps.length === 0) {
            alert('At least one step is required');
            return;
        }
        // 验证时间戳格式
        for (const step of currentAnnotation.steps) {
            // 支持两种格式：mm:ss.SSS 或 mm:ss
            if (!/^\d{2}:\d{2}(\.\d{3})?$/.test(step.timestamp)) {
                alert(`Step ${step.number} has invalid timestamp format. Should be mm:ss.SSS (e.g., 12:32.766)`);
                return;
            }
            if (!step.description.trim()) {
                alert(`Step ${step.number} description cannot be empty`);
                return;
            }
        }
    }

    try {
        const stem = currentVideo.replace(/\.mp4$/, '');
        const response = await fetch(`/api/annotation/${stem}.txt`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(currentAnnotation)
        });

        if (response.ok) {
            lastSavedAnnotationJSON = JSON.stringify(currentAnnotation);
            updateAutoSaveStatus('saved');
            alert('Saved successfully');
            loadVideos(); // 刷新列表状态
        } else {
            const error = await response.text();
            alert('Failed to save: ' + error);
        }
    } catch (error) {
        console.error('Failed to save annotation:', error);
        alert('Failed to save: ' + error.message);
    }
}

// 删除标注
async function deleteAnnotation() {
    if (!currentVideo) {
        alert('Please select a video first');
        return;
    }

    // 确认删除
    if (!confirm('Are you sure you want to delete the saved annotation?\n\nPre-annotation will be reloaded if available.')) {
        return;
    }

    try {
        const stem = currentVideo.replace(/\.mp4$/, '');
        const response = await fetch(`/api/annotation/${stem}.txt`, {
            method: 'DELETE'
        });

        if (response.ok) {
            alert('Annotation deleted');
            // 重新加载标注（会自动从预标注加载）
            await loadAnnotation(currentVideo);
            // 刷新列表状态
            loadVideos();
        } else {
            const error = await response.text();
            alert('Failed to delete: ' + error);
        }
    } catch (error) {
        console.error('Failed to delete annotation:', error);
        alert('Failed to delete: ' + error.message);
    }
}

// 更新时间显示
function updateTimeDisplay() {
    if (!player) return;
    currentTime.textContent = formatTimestamp(player.currentTime());
}

// 格式化时间为 mm:ss.SSS
// 使用 Math.round 转为整数毫秒后再拆分，避免浮点精度丢失导致的1ms偏差
function formatTimestamp(timeInSeconds) {
    const totalMs = Math.round(timeInSeconds * 1000);
    const minutes = Math.floor(totalMs / 60000);
    const seconds = Math.floor((totalMs % 60000) / 1000);
    const millis = totalMs % 1000;
    return `${String(minutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}.${String(millis).padStart(3, '0')}`;
}

// 解析时间戳为秒数
function parseTimestamp(timestamp) {
    // 支持两种格式：mm:ss.SSS 和 mm:ss
    const parts = timestamp.split(':');
    if (parts.length !== 2) return 0;
    
    const minutes = parseInt(parts[0], 10);
    if (isNaN(minutes)) return 0;
    
    // 解析秒和毫秒
    const secondsParts = parts[1].split('.');
    const seconds = parseInt(secondsParts[0], 10);
    const millis = secondsParts.length > 1 ? parseInt(secondsParts[1].padEnd(3, '0'), 10) : 0;
    
    if (isNaN(seconds) || isNaN(millis)) return 0;
    
    return minutes * 60 + seconds + millis / 1000;
}

// 插入当前时间戳
function insertCurrentTimestamp() {
    if (!player || !currentAnnotation || !currentAnnotation.is_tutorial) {
        return;
    }

    const time = player.currentTime();
    const timestamp = formatTimestamp(time);

    if (!currentAnnotation.steps) {
        currentAnnotation.steps = [];
    }

    // 确定插入位置：如果有选中的步骤，插入到其下方；否则插入到末尾
    let insertIndex;
    if (selectedStepIndex >= 0 && selectedStepIndex < currentAnnotation.steps.length) {
        insertIndex = selectedStepIndex + 1;
    } else {
        insertIndex = currentAnnotation.steps.length;
    }

    // 创建新步骤
    const newStep = {
        number: insertIndex + 1, // 临时编号，后面会重新编号
        timestamp: timestamp,
        description: ''
    };

    // 在指定位置插入步骤
    currentAnnotation.steps.splice(insertIndex, 0, newStep);
    
    // 重新编号所有步骤
    currentAnnotation.steps.forEach((s, i) => {
        s.number = i + 1;
    });
    
    // 更新选中索引为新插入的步骤
    selectedStepIndex = insertIndex;
    
    // 重新渲染编辑器
    renderEditor();
    
    // 聚焦到新步骤的描述框
    setTimeout(() => {
        const stepItems = document.querySelectorAll('.step-item');
        if (stepItems.length > insertIndex) {
            const newStepItem = stepItems[insertIndex];
            const descriptionInput = newStepItem.querySelector('.step-description');
            if (descriptionInput) {
                descriptionInput.focus();
            }
        }
    }, 100);

    // 暂停视频，方便用户输入描述
    player.pause();

    // 触发自动保存
    scheduleAutoSave();
}

// 跳转到指定时间戳
function seekToTimestamp(timestamp) {
    if (!player) return;
    
    const totalSeconds = parseTimestamp(timestamp);
    if (totalSeconds === 0 && timestamp !== '00:00.000' && timestamp !== '00:00') return;
    
    player.currentTime(totalSeconds);
    
    // 跳转后保持暂停状态
    player.pause();
}

// 复制当前时间戳到剪贴板
function copyCurrentTimestamp() {
    if (!player) return;
    
    const timestamp = currentTime.textContent;
    
    // 使用Clipboard API复制
    if (navigator.clipboard && navigator.clipboard.writeText) {
        navigator.clipboard.writeText(timestamp).then(() => {
            showCopyFeedback();
        }).catch(err => {
            console.error('Failed to copy timestamp:', err);
            // Fallback method
            fallbackCopyTimestamp(timestamp);
        });
    } else {
        // Fallback for older browsers
        fallbackCopyTimestamp(timestamp);
    }
}

// Fallback copy method for older browsers
function fallbackCopyTimestamp(text) {
    const textarea = document.createElement('textarea');
    textarea.value = text;
    textarea.style.position = 'fixed';
    textarea.style.opacity = '0';
    document.body.appendChild(textarea);
    textarea.select();
    try {
        document.execCommand('copy');
        showCopyFeedback();
    } catch (err) {
        console.error('Fallback copy failed:', err);
    }
    document.body.removeChild(textarea);
}

// 显示复制成功提示
function showCopyFeedback() {
    // 创建提示元素
    const feedback = document.createElement('div');
    feedback.textContent = 'Copied!';
    feedback.style.cssText = `
        position: fixed;
        bottom: 100px;
        left: 50%;
        transform: translateX(-50%);
        background-color: #4CAF50;
        color: white;
        padding: 12px 24px;
        border-radius: 4px;
        font-size: 14px;
        font-weight: 500;
        box-shadow: 0 2px 8px rgba(0,0,0,0.2);
        z-index: 10000;
        animation: fadeInOut 2s ease-in-out;
    `;
    
    // 添加动画
    const style = document.createElement('style');
    style.textContent = `
        @keyframes fadeInOut {
            0% { opacity: 0; transform: translateX(-50%) translateY(10px); }
            10% { opacity: 1; transform: translateX(-50%) translateY(0); }
            90% { opacity: 1; transform: translateX(-50%) translateY(0); }
            100% { opacity: 0; transform: translateX(-50%) translateY(-10px); }
        }
    `;
    document.head.appendChild(style);
    
    document.body.appendChild(feedback);
    
    // 2秒后移除
    setTimeout(() => {
        document.body.removeChild(feedback);
        document.head.removeChild(style);
    }, 2000);
}

// 加载模型标注
async function loadModelAnnotation(filename) {
    try {
        const stem = filename.replace(/\.mp4$/, '');
        const response = await fetch(`/api/model-annotation/${stem}.txt`);
        const data = await response.json();
        
        if (data.available) {
            currentModelAnnotation = data.annotation;
        } else {
            currentModelAnnotation = null;
        }
        
        renderModelPanel();
    } catch (error) {
        console.error('Failed to load model annotation:', error);
        currentModelAnnotation = null;
        renderModelPanel();
    }
}

// 渲染模型标注面板
function renderModelPanel() {
    const modelContent = document.getElementById('modelContent');
    
    if (!currentModelAnnotation) {
        modelContent.innerHTML = '<div class="model-unavailable">No model annotation available for this video</div>';
        return;
    }
    
    // 非教学视频
    if (!currentModelAnnotation.is_tutorial) {
        modelContent.innerHTML = '<div style="padding: 1rem;"><p style="color: #666;">Non-tutorial video</p></div>';
        return;
    }
    
    // 教学视频 - 显示标题和步骤
    let html = '<div style="padding: 1rem;">';
    
    // 标题
    html += '<div class="form-group">';
    html += '<label>Tutorial Title:</label>';
    html += `<input type="text" value="${escapeHtml(currentModelAnnotation.title || '')}" readonly style="background-color: #e9ecef;">`;
    html += '</div>';
    
    // 步骤
    if (currentModelAnnotation.steps && currentModelAnnotation.steps.length > 0) {
        html += '<div class="form-group">';
        html += '<label>Steps:</label>';
        html += '<div class="steps-container">';
        
        currentModelAnnotation.steps.forEach((step, index) => {
            html += '<div class="step-item" title="Click timestamp to seek video, click description to copy">';
            html += '<div class="step-header">';
            html += `<span class="step-number">Step ${step.number}</span>`;
            html += `<input type="text" class="step-timestamp model-timestamp" value="${escapeHtml(step.timestamp)}" readonly style="background-color: #e9ecef;" data-timestamp="${escapeHtml(step.timestamp)}" title="Click to seek video">`;
            html += '</div>';
            html += `<textarea class="step-description model-description" readonly style="background-color: #e9ecef;" data-description="${escapeHtml(step.description || '')}" title="Click to copy description">${escapeHtml(step.description || '')}</textarea>`;
            html += '</div>';
        });
        
        html += '</div>';
        html += '</div>';
    } else {
        html += '<p style="color: #666; padding: 1rem;">No steps found</p>';
    }
    
    html += '</div>';
    modelContent.innerHTML = html;
    
    // 添加事件监听器
    setupModelPanelListeners();
}

// 设置模型面板的事件监听器
function setupModelPanelListeners() {
    // 时间戳点击跳转
    document.querySelectorAll('.model-timestamp').forEach(timestampInput => {
        timestampInput.addEventListener('click', (e) => {
            e.stopPropagation();
            const timestamp = e.target.dataset.timestamp;
            seekToTimestamp(timestamp);
        });
    });
    
    // 描述点击复制（只在没有选择文字时触发）
    document.querySelectorAll('.model-description').forEach(descriptionArea => {
        descriptionArea.addEventListener('click', (e) => {
            // 检查是否有文字被选中
            const selection = window.getSelection();
            const selectedText = selection.toString();
            
            // 如果没有选中文字，则复制整个描述
            if (!selectedText) {
                const description = e.target.dataset.description;
                copyToClipboard(description, 'Description copied!');
            }
        });
        
        // 允许双击选择文字
        descriptionArea.addEventListener('dblclick', (e) => {
            e.stopPropagation();
        });
    });
}

// 通用复制到剪贴板函数
function copyToClipboard(text, message) {
    if (navigator.clipboard && navigator.clipboard.writeText) {
        navigator.clipboard.writeText(text).then(() => {
            showCopyMessage(message);
        }).catch(err => {
            console.error('Failed to copy:', err);
            fallbackCopy(text, message);
        });
    } else {
        fallbackCopy(text, message);
    }
}

// 备用复制方法
function fallbackCopy(text, message) {
    const textarea = document.createElement('textarea');
    textarea.value = text;
    textarea.style.position = 'fixed';
    textarea.style.opacity = '0';
    document.body.appendChild(textarea);
    textarea.select();
    try {
        document.execCommand('copy');
        showCopyMessage(message);
    } catch (err) {
        console.error('Fallback copy failed:', err);
    }
    document.body.removeChild(textarea);
}

// 显示复制消息
function showCopyMessage(message) {
    const feedback = document.createElement('div');
    feedback.textContent = message;
    feedback.style.cssText = `
        position: fixed;
        bottom: 100px;
        left: 50%;
        transform: translateX(-50%);
        background-color: #4CAF50;
        color: white;
        padding: 12px 24px;
        border-radius: 4px;
        font-size: 14px;
        font-weight: 500;
        box-shadow: 0 2px 8px rgba(0,0,0,0.2);
        z-index: 10000;
        animation: fadeInOut 2s ease-in-out;
    `;
    
    document.body.appendChild(feedback);
    
    setTimeout(() => {
        if (feedback.parentNode) {
            document.body.removeChild(feedback);
        }
    }, 2000);
}

// HTML 转义函数
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}


// Open native file/folder dialog and set the selected path.
async function openBrowser(inputId, mode) {
    try {
        const response = await fetch(`/api/dialog?mode=${encodeURIComponent(mode)}`);
        if (!response.ok) {
            const errorText = await response.text();
            alert(`Failed to open dialog: ${errorText}`);
            return;
        }

        const data = await response.json();
        if (data.path) {
            document.getElementById(inputId).value = data.path;
        }
    } catch (error) {
        console.error('Failed to open dialog:', error);
        alert('Failed to open dialog. Please paste the path manually.');
    }
}

// ========== 自动保存功能 ==========

// 调度自动保存（防抖）
function scheduleAutoSave() {
    if (!currentVideo || !currentAnnotation) return;

    updateAutoSaveStatus('unsaved');

    if (autoSaveTimer) {
        clearTimeout(autoSaveTimer);
    }

    autoSaveTimer = setTimeout(() => {
        performAutoSave();
    }, AUTO_SAVE_DELAY);
}

// 从表单同步数据到 currentAnnotation
function syncAnnotationFromForm() {
    if (!currentAnnotation) return;
    if (currentAnnotation.is_tutorial) {
        currentAnnotation.title = tutorialTitle.value.trim();
    }
    // 步骤数据已通过 input/change 事件实时同步
}

// 静默验证标注（不弹窗）
function validateAnnotationSilently() {
    if (!currentAnnotation) return false;

    if (!currentAnnotation.is_tutorial) {
        return true; // 非教学视频始终有效
    }

    if (!currentAnnotation.title) return false;
    if (!currentAnnotation.steps || currentAnnotation.steps.length === 0) return false;

    for (const step of currentAnnotation.steps) {
        if (!/^\d{2}:\d{2}(\.\d{3})?$/.test(step.timestamp)) return false;
        if (!step.description.trim()) return false;
    }

    return true;
}

// 执行自动保存
async function performAutoSave() {
    if (!currentVideo || !currentAnnotation) return;

    syncAnnotationFromForm();

    if (!validateAnnotationSilently()) {
        return; // 数据不完整，暂不保存
    }

    const annotationJSON = JSON.stringify(currentAnnotation);
    if (annotationJSON === lastSavedAnnotationJSON) {
        updateAutoSaveStatus('saved');
        return; // 无变化
    }

    updateAutoSaveStatus('saving');

    try {
        const stem = currentVideo.replace(/\.mp4$/, '');
        const response = await fetch(`/api/annotation/${stem}.txt`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: annotationJSON
        });

        if (response.ok) {
            lastSavedAnnotationJSON = annotationJSON;
            updateAutoSaveStatus('saved');
            loadVideos(); // 刷新列表状态，更新"已标注"标签
        } else {
            updateAutoSaveStatus('error');
        }
    } catch (error) {
        console.error('Auto-save failed:', error);
        updateAutoSaveStatus('error');
    }
}

// 立即执行自动保存（用于切换视频前）
async function flushAutoSave() {
    if (autoSaveTimer) {
        clearTimeout(autoSaveTimer);
        autoSaveTimer = null;
    }
    await performAutoSave();
}

// 更新自动保存状态指示器
function updateAutoSaveStatus(status) {
    const statusEl = document.getElementById('autoSaveStatus');
    if (!statusEl) return;

    statusEl.className = 'auto-save-status ' + status;

    switch(status) {
        case 'unsaved':
            statusEl.textContent = '● Unsaved';
            break;
        case 'saving':
            statusEl.textContent = '⏳ Saving...';
            break;
        case 'saved':
            statusEl.textContent = '✓ Saved';
            break;
        case 'error':
            statusEl.textContent = '✕ Save failed';
            break;
        default:
            statusEl.textContent = '';
    }
}

// 设置面板可调整大小
function setupResizableHandles() {
    const handles = document.querySelectorAll('.resize-handle');
    
    handles.forEach(handle => {
        let isResizing = false;
        let startX = 0;
        let startWidth = 0;
        let panel = null;
        
        handle.addEventListener('mousedown', (e) => {
            isResizing = true;
            startX = e.clientX;
            
            const panelId = handle.dataset.panel;
            panel = document.getElementById(panelId);
            startWidth = panel.offsetWidth;
            
            // 如果是视频面板，需要移除 flex 样式，改为固定宽度
            if (panelId === 'videoPanel') {
                panel.style.flex = 'none';
            }
            
            handle.classList.add('resizing');
            document.body.style.cursor = 'col-resize';
            document.body.style.userSelect = 'none';
            
            e.preventDefault();
        });
        
        document.addEventListener('mousemove', (e) => {
            if (!isResizing || !panel) return;
            
            const deltaX = e.clientX - startX;
            const newWidth = startWidth + deltaX;
            
            // 获取面板的最小和最大宽度
            const minWidth = parseInt(getComputedStyle(panel).minWidth) || 200;
            const maxWidth = parseInt(getComputedStyle(panel).maxWidth) || 2000;
            
            if (newWidth >= minWidth && newWidth <= maxWidth) {
                panel.style.width = newWidth + 'px';
            }
        });
        
        document.addEventListener('mouseup', () => {
            if (isResizing) {
                isResizing = false;
                handle.classList.remove('resizing');
                document.body.style.cursor = '';
                document.body.style.userSelect = '';
                panel = null;
            }
        });
    });
}
