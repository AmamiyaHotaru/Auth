<script setup>
import { ref } from 'vue';
import * as runtime from "../../../wailsjs/runtime/runtime";

// 定义props
const props = defineProps({
  show: {
    type: Boolean,
    default: false
  }
});

// 发射事件
const emit = defineEmits(['close']);

// 关闭弹窗
function closeDialog() {
  emit('close');
}

// 打开项目仓库链接
function openGithubRepo() {
  runtime.BrowserOpenURL("https://github.com/AmamiyaHotaru/Auth");
}
</script>

<template>
  <div v-if="show" class="about-dialog-overlay" @click.self="closeDialog">
    <div class="about-dialog">
      <div class="about-dialog-header">
        <h3>关于 Euthenticator</h3>
        <button class="close-button" @click="closeDialog">&times;</button>
      </div>
      <div class="about-dialog-body">
        <div class="about-logo">
          <img src="../../assets/images/logo-universal.png" alt="Euthenticator" class="about-logo-image">
        </div>
        <div class="about-version">
          <span class="version-label">版本</span>
          <span class="version-number">1.1.0</span>
        </div>
        <p class="about-description">
          Euthenticator 是一个纯本地的 TOTP 安全令牌生成器，支持 Steam 和 Google Authenticator 导入。
          所有数据均存储在本地，无需任何网络连接，保证您的账户安全。
        </p>
        <div class="about-details">
          <div class="detail-item">
            <div class="detail-icon">
              <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M9 19c-5 1.5-5-2.5-7-3m14 6v-3.87a3.37 3.37 0 0 0-.94-2.61c3.14-.35 6.44-1.54 6.44-7A5.44 5.44 0 0 0 20 4.77 5.07 5.07 0 0 0 19.91 1S18.73.65 16 2.48a13.38 13.38 0 0 0-7 0C6.27.65 5.09 1 5.09 1A5.07 5.07 0 0 0 5 4.77a5.44 5.44 0 0 0-1.5 3.78c0 5.42 3.3 6.61 6.44 7A3.37 3.37 0 0 0 9 18.13V22"></path>
              </svg>
            </div>
            <div class="detail-content">
              <div class="detail-label">项目仓库</div>
              <a @click.prevent="openGithubRepo" class="detail-link" style="cursor: pointer;">
                github.com/AmamiyaHotaru/Auth
              </a>
            </div>
          </div>
          <div class="detail-item">
            <div class="detail-icon">
              <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
                <circle cx="12" cy="7" r="4"></circle>
              </svg>
            </div>
            <div class="detail-content">
              <div class="detail-label">开发者</div>
              <div class="detail-value">AmamiyaHotaru</div>
            </div>
          </div>
          <div class="detail-item">
            <div class="detail-icon">
              <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"></path>
                <polyline points="22,6 12,13 2,6"></polyline>
              </svg>
            </div>
            <div class="detail-content">
              <div class="detail-label">联系邮箱</div>
              <a href="mailto:eunie@eunie.online" class="detail-link">eunie@eunie.online</a>
            </div>
          </div>
        </div>
      </div>
      <div class="about-dialog-footer">
        <button @click="closeDialog" class="about-close-button">
          关闭
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.about-dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.about-dialog {
  background-color: white;
  border-radius: 8px;
  width: 400px;
  max-width: 90%;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
  overflow: hidden;
}

.about-dialog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px;
  background-color: #f8f9fa;
  border-bottom: 1px solid #dee2e6;
}

.about-dialog-header h3 {
  margin: 0;
  font-size: 1.2em;
}

.close-button {
  background: none;
  border: none;
  font-size: 1.5em;
  cursor: pointer;
  color: #6c757d;
}

.close-button:hover {
  color: #212529;
}

.about-dialog-body {
  padding: 15px;
}

.about-logo {
  text-align: center;
  margin-bottom: 15px;
}

.about-logo-image {
  width: 80px;
  height: 80px;
}

.about-version {
  display: flex;
  justify-content: space-between;
  margin-bottom: 10px;
}

.version-label {
  font-weight: bold;
}

.version-number {
  color: #6c757d;
}

.about-description {
  margin-bottom: 15px;
  line-height: 1.5;
}

.about-details {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.detail-item {
  display: flex;
  align-items: center;
  gap: 10px;
}

.detail-icon {
  color: #6c757d;
}

.detail-content {
  flex-grow: 1;
}

.detail-label {
  font-weight: bold;
}

.detail-link {
  color: #4285F4;
  text-decoration: none;
}

.detail-link:hover {
  text-decoration: underline;
}

.about-dialog-footer {
  padding: 15px;
  text-align: right;
  background-color: #f8f9fa;
  border-top: 1px solid #dee2e6;
}

.about-close-button {
  background-color: #34A853;
  color: white;
  border: none;
  padding: 8px 15px;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.about-close-button:hover {
  background-color: #2c8a44;
}
</style>