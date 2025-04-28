<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue';

const props = defineProps({
  show: Boolean
});

const emit = defineEmits(['update:show', 'qrcode-detected', 'close']);

// 状态变量
const selectedFile = ref(null);
const errorMessage = ref('');
const isProcessing = ref(false);
const isDragging = ref(false);
const fileInput = ref(null);

// 生命周期钩子
onMounted(() => {
  // 添加粘贴事件监听
  document.addEventListener('paste', handlePaste);
});

onBeforeUnmount(() => {
  // 移除粘贴事件监听
  document.removeEventListener('paste', handlePaste);
});

// 方法
function triggerFileInput() {
  if (isProcessing.value) return; // 处理中时禁止选择新文件
  fileInput.value.click();
}

function handleFileSelect(event) {
  const file = event.target.files[0];
  if (file) {
    validateAndProcessFile(file);
  }
}

function handleFileDrop(event) {
  if (isProcessing.value) return; // 处理中时禁止拖放
  isDragging.value = false;
  const file = event.dataTransfer.files[0];
  if (file) {
    validateAndProcessFile(file);
  }
}

function handlePaste(event) {
  if (!props.show || isProcessing.value) return; // 处理中时禁止粘贴

  const items = (event.clipboardData || event.originalEvent.clipboardData).items;
  
  for (const item of items) {
    if (item.kind === 'file') {
      const blob = item.getAsFile();
      if (blob.type.indexOf('image/') === 0) {
        // 创建文件对象
        const pasteFile = new File([blob], `paste-${Date.now()}.png`, { 
          type: blob.type 
        });
        
        validateAndProcessFile(pasteFile);
        event.preventDefault();
        break;
      }
    }
  }
}

function validateAndProcessFile(file) {
  // 清除之前的错误信息
  errorMessage.value = '';
  
  // 验证文件是否为图片
  if (!file.type.startsWith('image/')) {
    errorMessage.value = '只支持上传图片文件';
    return;
  }
  
  // 立即设置处理状态为加载中
  isProcessing.value = true;
  
  // 处理文件
  processSelectedFile(file);
}

function processSelectedFile(file) {
  selectedFile.value = file;
  
  // 直接读取文件内容为二进制数组，不再创建预览
  readFileAsArrayBuffer(file);
}

function readFileAsArrayBuffer(file) {
  const reader = new FileReader();
  
  reader.onload = (e) => {
    const arrayBuffer = e.target.result;
    // 将ArrayBuffer转为Uint8Array
    const uint8Array = new Uint8Array(arrayBuffer);
    
    // 直接将二进制数据作为数组发送
    emit('qrcode-detected', Array.from(uint8Array));
  };
  
  reader.onerror = () => {
    isProcessing.value = false;
    errorMessage.value = '读取文件失败，请重试';
  };
  
  // 以二进制格式读取文件
  reader.readAsArrayBuffer(file);
}

// 外部调用，用于设置处理失败的错误信息
function setProcessingError(error) {
  isProcessing.value = false;
  errorMessage.value = error ? `二维码识别失败: ${error}` : '处理失败，请重试';
}

// 外部调用，用于重置处理状态
function resetProcessingState() {
  isProcessing.value = false;
}

// 暴露方法给父组件
defineExpose({
  setProcessingError,
  resetProcessingState
});

function handleClose() {
  // 如果正在处理，先询问用户是否确定要关闭
  if (isProcessing.value) {
    if (!confirm('正在处理图片，确定要取消吗？')) {
      return;
    }
  }
  
  // 重置表单
  selectedFile.value = null;
  errorMessage.value = '';
  isProcessing.value = false;
  
  // 清空文件输入元素
  if (fileInput.value) {
    fileInput.value.value = '';
  }
  
  // 关闭对话框
  emit('close');
}

function closeIfOverlay(event) {
  // 只有点击背景时才关闭
  if (event.target.classList.contains('upload-dialog-overlay')) {
    handleClose();
  }
}
</script>

<template>
  <div v-if="show" class="upload-dialog-overlay" @click="closeIfOverlay">
    <div class="upload-dialog">
      <div class="upload-dialog-header">
        <h3>扫描二维码</h3>
        <button class="close-button" @click="handleClose">&times;</button>
      </div>
      
      <div class="upload-dialog-body">
        <div 
          class="upload-drop-zone"
          @dragover.prevent
          @drop.prevent="handleFileDrop"
          @click="triggerFileInput"
          :class="{ 'active': isDragging, 'disabled': isProcessing }"
          @dragenter.prevent="!isProcessing && (isDragging = true)"
          @dragleave.prevent="!isProcessing && (isDragging = false)"
        >
          <input 
            type="file" 
            ref="fileInput" 
            class="file-input" 
            accept="image/*" 
            @change="handleFileSelect"
            :disabled="isProcessing"
          >
          <div v-if="!isProcessing" class="upload-icon">📷</div>
          <div v-else class="spinner"></div>
          <p v-if="!isProcessing">拖放二维码图片到此处，或点击选择文件</p>
          <p v-else>正在处理二维码，请稍候...</p>
          <p v-if="!isProcessing" class="upload-tip">
            支持粘贴上传，可以直接截图后 Ctrl+V 粘贴<br>
            <span class="support-note">支持谷歌身份验证器二维码导入</span>
          </p>
        </div>
        
        <div v-if="errorMessage" class="error-message">
          <div class="error-icon">!</div>
          {{ errorMessage }}
          <button class="retry-button" @click="errorMessage = ''; triggerFileInput()">重试</button>
        </div>
      </div>
      
      <div class="upload-dialog-footer">
        <button class="cancel-button" @click="handleClose">取消</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.upload-dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.upload-dialog {
  background-color: #fff;
  border-radius: 8px;
  width: 90%;
  max-width: 500px;
  max-height: 80vh;
  overflow-y: auto;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.15);
}

.upload-dialog-header {
  padding: 16px;
  border-bottom: 1px solid #e8f5e9;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.upload-dialog-header h3 {
  margin: 0;
  color: #4285F4;
}

.close-button {
  background: none;
  border: none;
  font-size: 20px;
  cursor: pointer;
  color: #999;
}

.close-button:hover {
  color: #333;
}

.upload-dialog-body {
  padding: 16px;
}

.upload-drop-zone {
  border: 2px dashed #4285F4;
  border-radius: 8px;
  padding: 30px;
  text-align: center;
  cursor: pointer;
  transition: border-color 0.3s, background-color 0.3s;
  position: relative;
}

.upload-drop-zone.active {
  border-color: #3367d6;
  background-color: #f0f5ff;
}

.upload-drop-zone.disabled {
  border-color: #ccc;
  cursor: not-allowed;
  opacity: 0.8;
}

.upload-drop-zone:hover:not(.disabled) {
  background-color: #f8f8f8;
}

.file-input {
  display: none;
}

.upload-icon {
  font-size: 40px;
  margin-bottom: 10px;
  color: #4285F4;
}

.upload-tip {
  font-size: 12px;
  color: #999;
  margin-top: 10px;
}

.support-note {
  color: #4285F4;
  font-weight: 500;
  display: inline-block;
  margin-top: 5px;
}

.error-message {
  background-color: #ffebee;
  color: #ff5252;
  padding: 15px;
  border-radius: 4px;
  margin-top: 15px;
  font-size: 14px;
  display: flex;
  align-items: center;
  animation: shake 0.5s ease-in-out;
  position: relative;
}

.error-icon {
  display: inline-flex;
  justify-content: center;
  align-items: center;
  min-width: 24px;
  height: 24px;
  border-radius: 50%;
  background-color: #ff5252;
  color: white;
  font-weight: bold;
  margin-right: 8px;
}

.retry-button {
  position: absolute;
  right: 10px;
  padding: 4px 10px;
  background-color: #ff5252;
  color: white;
  border: none;
  border-radius: 3px;
  cursor: pointer;
}

.retry-button:hover {
  background-color: #ff3333;
}

@keyframes shake {
  0%, 100% { transform: translateX(0); }
  10%, 30%, 50%, 70%, 90% { transform: translateX(-5px); }
  20%, 40%, 60%, 80% { transform: translateX(5px); }
}

.upload-dialog-footer {
  padding: 16px;
  border-top: 1px solid #e8f5e9;
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.cancel-button {
  padding: 8px 16px;
  background-color: #f5f5f5;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.cancel-button:hover {
  background-color: #e5e5e5;
}

/* 加载动画 */
.spinner {
  display: inline-block;
  width: 40px;
  height: 40px;
  margin-bottom: 10px;
  border: 4px solid rgba(66, 133, 244, 0.3);
  border-radius: 50%;
  border-top-color: #4285F4;
  animation: spin 1s ease-in-out infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>