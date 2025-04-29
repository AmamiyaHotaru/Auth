<script setup>
import { ref, onMounted } from 'vue';
import * as runtime from "../../../wailsjs/runtime/runtime";

// 窗口状态
const isMaximized = ref(false);

// 窗口最小化
async function minimizeWindow() {
  await runtime.WindowMinimise();
}

// 窗口最大化/还原
async function toggleMaximizeWindow() {
  isMaximized.value = !isMaximized.value;
  await runtime.WindowToggleMaximise();
}

// 关闭窗口（退出应用）
async function closeWindow() {
  await runtime.Quit();
}

// 检查窗口是否最大化
async function checkWindowState() {
  isMaximized.value = await runtime.WindowIsMaximised();
}

// 组件挂载时检查窗口状态
onMounted(async () => {
  // 检查窗口状态
  await checkWindowState();
  
  // 监听窗口大小变化
  window.addEventListener('resize', async () => {
    await checkWindowState();
  });
});

// 在组件卸载时清理事件监听
defineExpose({
  checkWindowState
});
</script>

<template>
  <!-- 窗口标题栏 - 可拖动区域 -->
  <div class="window-titlebar" style="--wails-draggable:drag">
    <div class="titlebar-drag-region">
      <div class="app-icon"></div>
    </div>
    <div class="window-controls">
      <button @click="minimizeWindow" class="window-control-button minimize-button" title="最小化">
        <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 12 12">
          <rect x="1" y="6" width="10" height="1" fill="currentColor" />
        </svg>
      </button>
      <button @click="toggleMaximizeWindow" class="window-control-button maximize-button" title="最大化/还原">
        <svg v-if="isMaximized" xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 12 12">
          <path d="M3 1v3H1v7h7V9h3V1H3zm1 2h5v5H8V4H4V3z" fill="currentColor" />
        </svg>
        <svg v-else xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 12 12">
          <rect x="1" y="1" width="10" height="10" stroke="currentColor" fill="none" stroke-width="1" />
        </svg>
      </button>
      <button @click="closeWindow" class="window-control-button close-button" title="关闭">
        <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 12 12">
          <path d="M 1,1 L 11,11 M 1,11 L 11,1" stroke="currentColor" stroke-width="1.5" fill="none" />
        </svg>
      </button>
    </div>
  </div>
</template>

<style scoped>
.window-titlebar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: #343a40;
  color: white;
  padding: 5px 10px;
  flex-shrink: 0;
}

.titlebar-drag-region {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-grow: 1;
  cursor: grab;
}

.titlebar-icon {
  width: 20px;
  height: 20px;
}

.window-controls {
  display: flex;
  gap: 5px;
}

.window-control-button {
  background: none;
  border: none;
  color: white;
  width: 30px;
  height: 30px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.window-control-button:hover {
  background-color: rgba(255, 255, 255, 0.2);
}

.window-control-button.close-button:hover {
  background-color: #EA4335;
}
</style>