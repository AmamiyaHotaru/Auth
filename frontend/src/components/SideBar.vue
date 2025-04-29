<script setup>
import { ref } from 'vue';
import * as runtime from "../../wailsjs/runtime/runtime";

// 定义props
const props = defineProps({
  accountsCount: {
    type: Number,
    default: 0
  }
});

// 发射事件
const emit = defineEmits(['show-about']);

// 打开项目仓库链接
function openGithubRepo() {
  runtime.BrowserOpenURL("https://github.com/AmamiyaHotaru/Auth");
}

// 显示关于信息
function showAbout() {
  emit('show-about');
}
</script>

<template>
  <aside class="sidebar">
    <div class="sidebar-header">
      <div class="app-logo">
        <img src="../assets/images/logo-universal.png" alt="Euthenticator" class="logo-image">
      </div>
      <h1 class="app-title">Euthenticator</h1>
    </div>
    <div class="sidebar-content">
      <!-- 这里可以未来添加导航选项或过滤器等 -->
      <div class="sidebar-section">
        <h2 class="section-title">账户列表</h2>
        <!-- 账户统计信息 -->
        <div class="account-stats">
          <span class="stat-item">总计: {{ accountsCount }}</span>
        </div>
      </div>
    </div>
    <div class="sidebar-footer">
      <!-- 侧边栏底部添加关于按钮 -->
      <button @click="showAbout" class="about-button">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <circle cx="12" cy="12" r="10"></circle>
          <line x1="12" y1="16" x2="12" y2="12"></line>
          <line x1="12" y1="8" x2="12.01" y2="8"></line>
        </svg>
        <span>关于软件</span>
      </button>
    </div>
  </aside>
</template>

<style scoped>
.sidebar {
  width: 220px;
  background-color: #343a40;
  color: white;
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
}

.sidebar-header {
  padding: 20px;
  text-align: center;
  border-bottom: 1px solid #495057;
}

.app-logo {
  margin-bottom: 10px;
}

.logo-image {
  width: 50px;
  height: 50px;
}

.app-title {
  font-size: 1.5em;
  margin: 0;
}

.sidebar-content {
  flex-grow: 1;
  padding: 15px;
}

.sidebar-section {
  margin-bottom: 20px;
}

.section-title {
  font-size: 1.2em;
  margin-bottom: 10px;
}

.account-stats {
  font-size: 0.9em;
}

.sidebar-footer {
  padding: 15px;
  border-top: 1px solid #495057;
}

.about-button {
  display: flex;
  align-items: center;
  gap: 5px;
  background: none;
  border: none;
  color: white;
  font-size: 0.9em;
  cursor: pointer;
  transition: color 0.2s ease;
}

.about-button:hover {
  color: #adb5bd;
}
</style>