<script setup>
import { ref, computed } from 'vue';
import * as runtime from "../../../wailsjs/runtime/runtime";

// 定义props
const props = defineProps({
  account: {
    type: Object,
    required: true
  },
  selectionMode: {
    type: Boolean,
    default: false
  },
  isSelected: {
    type: Boolean,
    default: false
  },
  isCodeVisible: {
    type: Boolean,
    default: false
  }
});

// 发射事件
const emit = defineEmits([
  'toggle-selection',
  'toggle-code-visibility',
  'copy-code',
  'interaction-start',
  'interaction-end',
  'click'
]);

// 复制状态
const isCopied = ref(false);

// 触发选择切换
function toggleSelection(event) {
  event.stopPropagation();
  emit('toggle-selection', props.account.ID);
}

// 触发代码可见性切换
function toggleCodeVisibility(event) {
  // 添加对 event 的空值检查
  if (event) {
    event.stopPropagation();
  }
  emit('toggle-code-visibility', props.account.ID, event);
}

// 复制验证码到剪贴板
async function copyCode(event) {
  event.stopPropagation();
  emit('copy-code', props.account.Code, props.account.ID, event);
  
  // 添加视觉反馈
  const targetElement = event.currentTarget;
  targetElement.classList.add('copy-active');
  
  // 500毫秒后移除视觉反馈
  setTimeout(() => {
    targetElement.classList.remove('copy-active');
  }, 500);
}

// 处理交互开始
function handleInteractionStart(event) {
  emit('interaction-start', event, props.account);
}

// 处理交互结束
function handleInteractionEnd(event) {
  emit('interaction-end', event, props.account);
}

// 处理点击
function handleClick() {
  emit('click', props.account);
}

// 计算进度条偏移量
function calculateProgressOffset(timeLeft) {
  const radius = 15;
  const circumference = 2 * Math.PI * radius;
  const percentage = timeLeft / 30;
  return circumference * (1 - percentage);
}

// 格式化验证码
function formatCode(code) {
  // 将验证码格式化为 3-3 结构，例如：123 456
  if (code && code.length >= 6) {
    return code.substring(0, 3) + ' ' + code.substring(3, 6);
  }
  return code;
}
</script>

<template>
  <div 
    class="account-card"
    :class="{ 
      'selected': isSelected, 
      'selection-active': selectionMode 
    }"
    @mousedown="handleInteractionStart"
    @mouseup="handleInteractionEnd"
    @touchstart.passive="handleInteractionStart"
    @touchend="handleInteractionEnd"
    @click="handleClick"
    @contextmenu.prevent
  >
    <!-- 选择指示器 -->
    <div v-if="selectionMode" class="selection-indicator">
      <input type="checkbox" :checked="isSelected" @click.stop="toggleSelection" class="selection-checkbox">
    </div>

    <!-- 账户详情 -->
    <div class="card-content">
      <div class="account-header">
        <div class="service-info">
          <div class="service-icon">
            {{ account.ServerName && account.ServerName.length > 0 
              ? account.ServerName.charAt(0).toUpperCase() 
              : account.AccountName.charAt(0).toUpperCase() }}
          </div>
          <div class="service-details">
            <span class="server-name">{{ account.ServerName }}</span>
            <span class="account-name">{{ account.AccountName }}</span>
          </div>
        </div>
        <div class="time-counter">
          <div class="progress-ring-wrapper">
            <svg class="progress-ring" width="36" height="36">
              <circle class="progress-ring-bg" r="15" cx="18" cy="18" />
              <circle 
                class="progress-ring-circle" 
                r="15" 
                cx="18" 
                cy="18" 
                :stroke-dashoffset="calculateProgressOffset(account.timeLeft)" 
              />
            </svg>
            <span class="time-left">{{ account.timeLeft }}</span>
          </div>
        </div>
      </div>
      
      <div class="code-section">
        <span 
          v-if="isCodeVisible" 
          class="code" 
          @click="copyCode"
        >{{ formatCode(account.Code) }}</span>
        <span 
          v-else 
          class="hidden-code" 
          @click="copyCode"
        >••• •••</span>
        <div class="code-actions">
          <button @click="toggleCodeVisibility" class="action-button" :class="{'active': isCodeVisible}" title="显示/隐藏">
            <svg v-if="isCodeVisible" xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"></path>
              <line x1="1" y1="1" x2="23" y2="23"></line>
            </svg>
            <svg v-else xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path>
              <circle cx="12" cy="12" r="3"></circle>
            </svg>
          </button>
          <button @click="copyCode" class="action-button" title="复制">
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
              <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1-2 2v1"></path>
            </svg>
          </button>
        </div>
      </div>
      
      <!-- 复制成功提示 -->
      <div v-if="account.isCopied" class="copied-indicator">
        已复制!
      </div>
    </div>
  </div>
</template>

<style scoped>
.account-card {
  display: flex;
  flex-direction: column;
  background-color: white;
  border-radius: 8px;
  padding: 15px;
  width: calc(33.333% - 10px);
  box-shadow: 0 2px 4px rgba(0,0,0,0.05);
  border: 1px solid #dee2e6;
  transition: all 0.2s ease;
  position: relative;
}

.account-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0,0,0,0.1);
}

.account-card.selected {
  background-color: #e8f0fe;
  border: 1px solid #4285F4;
}

.selection-mode .account-card:hover {
  background-color: #f1f3f5;
}

.selection-indicator {
  position: absolute;
  top: 10px;
  left: 10px;
}

.selection-checkbox {
  width: 20px;
  height: 20px;
  cursor: pointer;
}

.card-content {
  display: flex;
  flex-direction: column;
}

.account-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.service-info {
  display: flex;
  align-items: center;
}

.service-icon {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background-color: #4285F4;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5em;
  font-weight: bold;
  margin-right: 10px;
}

.service-details {
  display: flex;
  flex-direction: column;
}

.server-name {
  font-weight: bold;
  font-size: 1.2em;
}

.account-name {
  font-size: 0.9em;
  color: #6c757d;
}

.time-counter {
  display: flex;
  align-items: center;
}

.progress-ring-wrapper {
  position: relative;
  width: 36px;
  height: 36px;
}

.progress-ring {
  position: absolute;
  top: 0;
  left: 0;
}

.progress-ring-bg {
  fill: none;
  stroke: #e0e0e0;
  stroke-width: 2;
}

.progress-ring-circle {
  fill: none;
  stroke: #4285F4;
  stroke-width: 2;
  stroke-dasharray: 94.2; /* 2 * Math.PI * 15 */
  transition: stroke-dashoffset 0.2s linear;
}

.time-left {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-size: 0.85em;
  font-weight: bold;
  color: #4285F4;
}

.code-section {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.code, .hidden-code {
  font-size: 1.5em;
  font-weight: bold;
  color: #0d47a1;
  letter-spacing: 2px;
  font-family: 'Courier New', Courier, monospace;
  cursor: pointer;
  padding: 5px 8px;
  border-radius: 4px;
  background-color: #f8f9fa;
  transition: background-color 0.2s;
}

.code:hover, .hidden-code:hover {
  background-color: #e9ecef;
}

.code-actions {
  display: flex;
  gap: 10px;
}

.action-button {
  background: none;
  border: none;
  cursor: pointer;
  color: #6c757d;
  padding: 5px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
}

.action-button:hover {
  background-color: rgba(0, 0, 0, 0.05);
  color: #4285F4;
}

.action-button.active {
  color: #4285F4;
}

.copied-indicator {
  position: absolute;
  bottom: 5px;
  right: 10px;
  background-color: #34A853;
  color: white;
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 0.75em;
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(5px); }
  to { opacity: 1; transform: translateY(0); }
}

/* 调整卡片样式以适应不同分辨率 */
@media (max-width: 1200px) {
  .account-card {
    width: calc(50% - 10px);
  }
}

@media (max-width: 800px) {
  .account-card {
    width: 100%;
  }
}
</style>