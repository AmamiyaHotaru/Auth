<script setup>
import { ref, onMounted, computed, onUnmounted, nextTick } from 'vue';
import { GetSecretsList, InsertSecret, DeleteSecret, RecognizeQRCode } from '../../wailsjs/go/main/App';
import ConfirmationDialog from './ConfirmationDialog.vue';
import ManualEntryDialog from './ManualEntryDialog.vue';
import AddOptionsMenu from './AddOptionsMenu.vue';
import QrCodeUploadDialog from './QrCodeUploadDialog.vue';
import * as runtime from "../../wailsjs/runtime/runtime";

// --- 状态 ---
const accounts = ref([]);
// 存储显示验证码的账户ID
const visibleCodes = ref(new Set());
// 用于追踪复制按钮的状态
const copiedCodes = ref(new Set());
// 全局显示/隐藏所有验证码的状态
const showAllCodes = ref(false);

// 显示或隐藏验证码
function toggleCodeVisibility(accountId, event) {
  event.stopPropagation(); // 阻止事件冒泡
  if (visibleCodes.value.has(accountId)) {
    visibleCodes.value.delete(accountId);
  } else {
    visibleCodes.value.add(accountId);
  }
}

// 判断验证码是否可见
function isCodeVisible(accountId) {
  return showAllCodes.value || visibleCodes.value.has(accountId);
}

// 切换显示/隐藏所有验证码
function toggleAllCodesVisibility(event) {
  event.stopPropagation(); // 阻止事件冒泡
  showAllCodes.value = !showAllCodes.value;
}

// 复制验证码到剪贴板
async function copyCodeToClipboard(code, accountId, event) {
  event.stopPropagation(); // 阻止事件冒泡
  
  // 添加视觉反馈
  const targetElement = event.currentTarget;
  targetElement.classList.add('copy-active');
  
  // 500毫秒后移除视觉反馈
  setTimeout(() => {
    targetElement.classList.remove('copy-active');
  }, 500);
  
  console.log('复制验证码:', code);
  const success = await runtime.ClipboardSetText(code);
  console.log('复制结果:', success);
  // 添加短暂的"已复制"状态
  copiedCodes.value.add(accountId);
  setTimeout(() => {
    copiedCodes.value.delete(accountId);
  }, 1500);
      
}

// --- 辅助函数 ---
// 计算当前时间在 30 秒周期内的剩余秒数
function calculateTimeLeft() {
    const nowSeconds = Math.floor(Date.now() / 1000);
    return 30 - (nowSeconds % 30);
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

// 获取账户列表
const getSecretsList = async () => {
  try {
    const response = await GetSecretsList();
    processAccountList(response);
  } catch (error) {
    console.error('获取账户列表时出错:', error);
    accounts.value = [];
  }
};

// 处理账户列表数据，添加 timeLeft 字段
function processAccountList(list) {
  if (!list || list.length === 0) {
    console.log('获取到的账户列表为空');
    accounts.value = [];
    return;
  }
  
  // 计算当前时间周期剩余秒数
  const currentTimeLeft = calculateTimeLeft();
  
  accounts.value = list.map(account => ({
    ...account, 
    timeLeft: currentTimeLeft
  }));
  
  console.log('已更新账户列表，数量:', accounts.value.length);
}

// --- 其他状态和变量 ---
let timerInterval = null; // TOTP 计时器
let longPressTimer = null; // 长按检测计时器
const longPressDuration = 500; // 长按所需毫秒数

// --- 选择状态管理 ---
const selectionMode = ref(false); 
const selectedAccountIds = ref(new Set()); 
const selectedCount = computed(() => selectedAccountIds.value.size); 
const ignoreNextClick = ref(false); 

// --- 确认对话框状态 ---
const showDeleteConfirmation = ref(false); // 是否显示删除确认对话框
const confirmationTitle = ref(''); // 对话框标题
const confirmationMessage = ref(''); // 对话框消息

// --- 添加账户菜单和对话框状态 ---
const showAddOptions = ref(false); // 是否显示添加选项菜单
const showManualEntryDialog = ref(false); // 是否显示手动添加对话框
const showQrCodeDialog = ref(false); // 是否显示二维码扫描对话框
const qrCodeDialogRef = ref(null); // 二维码上传对话框引用
const addButtonPosition = ref(null); // 添加按钮位置
const qrCodeImageData = ref(null); // 二维码图片数据

// --- 关于按钮状态 ---
const showAboutInfo = ref(false); // 是否显示关于信息

// --- 事件处理 ---
// 处理条目交互开始（鼠标按下或触摸开始）
function handleItemInteractionStart(event, account) {
  // 对触摸事件，不立即阻止默认行为，以允许滚动
  if (event.type === 'touchstart') {
      // 允许滚动
  }
  // 清除任何现有的长按计时器
  clearTimeout(longPressTimer);
  ignoreNextClick.value = false; // 重置忽略标志

  // 启动长按计时器
  longPressTimer = setTimeout(() => {
    enterSelectionMode(account.ID); // 使用 ID 而非 id
    ignoreNextClick.value = true; // 设置标志，因为长按后通常会紧跟一个 click 事件
    // 可选：提供触觉反馈
    // navigator.vibrate(50);
  }, longPressDuration);
}

// 处理条目交互结束（鼠标松开或触摸结束）
function handleItemInteractionEnd(event, account) {
  // 清除长按计时器
  clearTimeout(longPressTimer);
  longPressTimer = null;
}

// 处理条目点击事件
function handleItemClick(account) {
    // 如果设置了忽略标志，则消耗它并直接返回
    if (ignoreNextClick.value) {
        ignoreNextClick.value = false; // 消耗标志
        return;
    }

    // 如果处于选择模式，则切换选中状态
    if (selectionMode.value) {
        toggleSelection(account.ID); // 使用 ID 而非 id
    }
    // 否则（非选择模式下的简单点击），不执行任何操作
}

// 进入选择模式
function enterSelectionMode(accountID) {
  if (!selectionMode.value) { // 确保只在首次进入时执行
      selectionMode.value = true;
      selectedAccountIds.value.clear();
      selectedAccountIds.value.add(accountID);
  }
}

// 取消选择模式
function cancelSelectionMode() {
  selectionMode.value = false;
  selectedAccountIds.value.clear();
  ignoreNextClick.value = false; // 退出时也重置标志
}

// 切换账户的选中状态
function toggleSelection(accountID) {
  if (selectedAccountIds.value.has(accountID)) {
    selectedAccountIds.value.delete(accountID);
    // 如果取消选中后没有选中的项目了，则退出选择模式
    if (selectedAccountIds.value.size === 0) {
      cancelSelectionMode();
    }
  } else {
    selectedAccountIds.value.add(accountID);
  }
}

// --- 删除逻辑 ---
// 请求删除选中的账户
function requestDeleteSelected() {
  if (selectedCount.value === 0) return;
  confirmationTitle.value = '删除确认';
  confirmationMessage.value = `确定要删除选中的 ${selectedCount.value} 个账户吗？此操作无法撤销。`;
  showDeleteConfirmation.value = true;
}

// 执行删除选中的账户操作
async function deleteSelectedAccounts() {
  showDeleteConfirmation.value = false; // 先关闭对话框
  const idsToDelete = Array.from(selectedAccountIds.value);
  console.log(`尝试删除账户，ID: ${idsToDelete.join(', ')}`);

  try {
    // 调用Go后端的DeleteSecret方法删除账户
    await DeleteSecret(idsToDelete);
    console.log(`成功删除账户，ID: ${idsToDelete.join(', ')}`);
    
    // 刷新账户列表，确保显示最新数据
    await getSecretsList();
    
  } catch (error) {
    console.error(`删除账户 ${idsToDelete.join(', ')} 时出错:`, error);
  } finally {
     // 删除尝试后退出选择模式
     cancelSelectionMode();
  }
}

// 处理取消删除操作
function handleDeletionCancel() {
    showDeleteConfirmation.value = false;
}

// --- 添加账户相关逻辑 ---
// 显示添加选项菜单
function showAddOptionsMenu(event) {
  if (selectionMode.value) return; // 选择模式下不允许添加

  console.log('显示添加选项菜单...');
  event.stopPropagation(); // 阻止事件冒泡
  
  // 计算并设置菜单位置
  const button = event.currentTarget; // 使用currentTarget代替target以确保获取的是按钮元素
  const rect = button.getBoundingClientRect();
  
  // 调整位置计算，确保菜单显示在按钮正下方
  addButtonPosition.value = {
    top: rect.bottom + 5, // 菜单顶部位置为按钮底部加5px间距
    left: rect.left-80 // 菜单左侧对齐按钮左侧
  };
  
  // 显示选项菜单
  showAddOptions.value = true;
  
  console.log('菜单位置:', addButtonPosition.value);
}

// 处理添加选项选择
function handleAddOptionSelect(option) {
  showAddOptions.value = false;
  
  if (option === 'manual') {
    // 显示手动添加对话框
    showManualEntryDialog.value = true;
  } else if (option === 'scan') {
    // 显示二维码扫描对话框
    showQrCodeDialog.value = true;
  }
}

// 关闭选项菜单
function closeAddOptions() {
  showAddOptions.value = false;
}

// 处理手动添加账户
async function handleManualAdd(formData) {
  try {
    console.log('提交的账户数据:', formData);
    
    // 调用Go后端添加秘钥，传递账户类型参数
    await InsertSecret(
      formData.accountName, 
      formData.serverName, 
      formData.secret, 
      formData.accountType
    );
    
    // 刷新账户列表
    await getSecretsList();
    
    // 关闭对话框
    showManualEntryDialog.value = false;
    
  } catch (error) {
    console.error('添加账户失败:', error);
  }
}

// 关闭手动添加对话框
function closeManualEntryDialog() {
  showManualEntryDialog.value = false;
}

// 关闭二维码扫描对话框
function closeQrCodeDialog() {
  showQrCodeDialog.value = false;
  qrCodeImageData.value = null;
}

// 处理二维码检测事件
async function handleQrCodeDetected(data) {
  console.log('检测到二维码图片数据，长度:', data.length);
  
  try {
    // 设置处理状态（通过引用调用组件方法）
    if (qrCodeDialogRef.value) {
      // 这里不关闭状态，因为需要等待后端处理完成
    }
    
    // 直接传递图像数据数组，不需要额外处理
    await RecognizeQRCode(data);
    
    
    // 检测成功后刷新账户列表
    await getSecretsList();
    
    // 重置处理状态
    if (qrCodeDialogRef.value) {
      qrCodeDialogRef.value.resetProcessingState();
    }
    
    // 关闭对话框
    showQrCodeDialog.value = false;
    
  } catch (error) {
    console.error('二维码识别失败:', error);
    
    // 设置错误状态
    if (qrCodeDialogRef.value) {
      qrCodeDialogRef.value.setProcessingError(error);
    }
    
    
    // 注意：不关闭对话框，让用户可以重试
  }
}

// --- TOTP 计时器逻辑 ---
// 每 30 秒更新一次验证码
async function updateCodes() {
  console.log("30 秒周期结束，获取新验证码...");
  await getSecretsList(); 
}

// 处理计时器每秒的 tick
function handleTimerTick() {
    const newTimeLeft = calculateTimeLeft();
    
    // 检测是否需要更新验证码（如果 timeLeft 从 1 变为 30，说明过了一个周期）
    let needsUpdate = false;
    if (accounts.value.length > 0) {
      needsUpdate = newTimeLeft === 30 && accounts.value[0].timeLeft === 1;
    }
    
    // 更新所有账户的 timeLeft
    accounts.value = accounts.value.map(acc => ({
      ...acc,
      timeLeft: newTimeLeft
    }));

    // 如果需要更新验证码，无论是否在选择模式下都刷新
    if (needsUpdate) {
      updateCodes();
    }
}

// 启动计时器 interval
function startTimerInterval() {
    if (timerInterval) clearInterval(timerInterval);
    // 先立即执行一次，确保 timeLeft 是正确的
    handleTimerTick();
    timerInterval = setInterval(handleTimerTick, 1000); // 每秒执行一次
}

// --- 生命周期钩子 ---
onMounted(async () => {
  console.log('组件已挂载，获取初始数据...');
  await getSecretsList(); // 获取初始数据
  startTimerInterval(); // 启动计时器
});

onUnmounted(() => {
  if (timerInterval) clearInterval(timerInterval); // 组件卸载时清除计时器
  clearTimeout(longPressTimer); // 清除可能存在的长按计时器
});

// 计算圆形进度条的样式
const getProgressStyle = (timeLeft) => {
  // 计算百分比，确保不低于 0
  const percentage = Math.max(0, timeLeft / 30) * 100;
  return {
    background: `radial-gradient(white 60%, transparent 61%), conic-gradient(#4285F4 ${percentage}%, #e0e0e0 ${percentage}%)`
  };
};

</script>

<template>
  <div id="authenticator-container" :class="{ 'selection-mode': selectionMode }">
    <!-- 侧边栏 -->
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
            <span class="stat-item">总计: {{ accounts.length }}</span>
          </div>
        </div>
      </div>
      <div class="sidebar-footer">
        <!-- 侧边栏底部添加关于按钮 -->
        <button @click="showAboutInfo = true" class="about-button">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="12" cy="12" r="10"></circle>
            <line x1="12" y1="16" x2="12" y2="12"></line>
            <line x1="12" y1="8" x2="12.01" y2="8"></line>
          </svg>
          <span>关于软件</span>
        </button>
      </div>
    </aside>

    <!-- 主内容区域 -->
    <main class="main-content">
      <!-- 条件渲染的头部 -->
      <header v-if="!selectionMode" class="content-header">
        <div class="header-left">
          <h2>安全令牌</h2>
        </div>
        <div class="header-right">
          <button @click="toggleAllCodesVisibility($event)" class="toolbar-button" :class="{'active': showAllCodes}" title="显示/隐藏所有验证码">
            <svg v-if="showAllCodes" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-eye-off">
              <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"></path>
              <line x1="1" y1="1" x2="23" y2="23"></line>
            </svg>
            <svg v-else xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-eye">
              <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path>
              <circle cx="12" cy="12" r="3"></circle>
            </svg>
            <span class="button-label">{{showAllCodes ? '隐藏所有' : '显示所有'}}</span>
          </button>
          <button @click="showAddOptionsMenu($event)" class="toolbar-button primary-button">
            <span class="button-icon">+</span>
            <span class="button-label">添加账户</span>
          </button>
        </div>
      </header>
      
      <header v-else class="selection-header">
        <div class="header-left">
          <button @click="cancelSelectionMode" class="cancel-selection-btn">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M19 12H5M12 19l-7-7 7-7"/>
            </svg>
            返回
          </button>
          <span class="selection-count">{{ selectedCount }} 项已选中</span>
        </div>
        <div class="header-right">
          <button @click="requestDeleteSelected" :disabled="selectedCount === 0" class="toolbar-button danger-button">
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <polyline points="3 6 5 6 21 6"></polyline>
              <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1-2-2h4a2 2 0 0 1-2 2v1"></path>
            </svg>
            <span class="button-label">删除选中</span>
          </button>
        </div>
      </header>

      <!-- 账户列表部分 -->
      <div class="accounts-container">
        <!-- 没有账户时的提示 -->
        <div v-if="accounts.length === 0" class="no-accounts">
          <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
            <rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
            <path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
          </svg>
          <p>还没有添加任何验证器账户</p>
          <button @click="showAddOptionsMenu($event)" class="desktop-button primary-button">
            <span class="button-icon">+</span>
            <span>添加账户</span>
          </button>
        </div>

        <!-- 账户列表项 -->
        <div v-else class="accounts-list">
          <div v-for="account in accounts" :key="account.ID"
              class="account-card"
              :class="{ 'selected': selectedAccountIds.has(account.ID), 'selection-active': selectionMode }"
              @mousedown="handleItemInteractionStart($event, account)"
              @mouseup="handleItemInteractionEnd($event, account)"
              @touchstart.passive="handleItemInteractionStart($event, account)"
              @touchend="handleItemInteractionEnd($event, account)"
              @click="handleItemClick(account)"
              @contextmenu.prevent> <!-- 阻止条目上的右键菜单 -->

            <!-- 选择指示器 -->
            <div v-if="selectionMode" class="selection-indicator">
                <input type="checkbox" :checked="selectedAccountIds.has(account.ID)" @click.stop="toggleSelection(account.ID)" class="selection-checkbox">
            </div>

            <!-- 账户详情 -->
            <div class="card-content">
              <div class="account-header">
                <div class="service-info">
                  <div class="service-icon">{{ account.ServerName.charAt(0).toUpperCase() }}</div>
                  <div class="service-details">
                    <span class="server-name">{{ account.ServerName }}</span>
                    <span class="account-name">{{ account.AccountName }}</span>
                  </div>
                </div>
                <div class="time-counter">
                  <div class="progress-ring-wrapper">
                    <svg class="progress-ring" width="36" height="36">
                      <circle class="progress-ring-bg" r="15" cx="18" cy="18" />
                      <circle class="progress-ring-circle" 
                              r="15" 
                              cx="18" 
                              cy="18" 
                              :stroke-dashoffset="calculateProgressOffset(account.timeLeft)" />
                    </svg>
                    <span class="time-left">{{ account.timeLeft }}</span>
                  </div>
                </div>
              </div>
              
              <div class="code-section">
                <span v-if="isCodeVisible(account.ID)" class="code" @click="copyCodeToClipboard(account.Code, account.ID, $event)">{{ formatCode(account.Code) }}</span>
                <span v-else class="hidden-code" @click="copyCodeToClipboard(account.Code, account.ID, $event)">••• •••</span>
                <div class="code-actions">
                  <button @click="toggleCodeVisibility(account.ID, $event)" class="action-button" :class="{'active': isCodeVisible(account.ID)}" title="显示/隐藏">
                    <svg v-if="isCodeVisible(account.ID)" xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                      <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"></path>
                      <line x1="1" y1="1" x2="23" y2="23"></line>
                    </svg>
                    <svg v-else xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                      <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path>
                      <circle cx="12" cy="12" r="3"></circle>
                    </svg>
                  </button>
                  <button @click="copyCodeToClipboard(account.Code, account.ID, $event)" class="action-button" title="复制">
                    <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                      <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
                      <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1-2 2v1"></path>
                    </svg>
                  </button>
                </div>
              </div>
              
              <!-- 复制成功提示 -->
              <div v-if="copiedCodes.has(account.ID)" class="copied-indicator">
                已复制!
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>

    <!-- 添加选项菜单 -->
    <AddOptionsMenu
      :show="showAddOptions"
      :position="addButtonPosition"
      @select="handleAddOptionSelect"
      @close="closeAddOptions"
    />

    <!-- 手动添加账户对话框 -->
    <ManualEntryDialog
      :show="showManualEntryDialog"
      @confirm="handleManualAdd"
      @cancel="closeManualEntryDialog"
    />

    <!-- 二维码上传对话框 -->
    <QrCodeUploadDialog
      ref="qrCodeDialogRef"
      :show="showQrCodeDialog"
      @qrcode-detected="handleQrCodeDetected"
      @close="closeQrCodeDialog"
    />

    <!-- 确认对话框 -->
    <ConfirmationDialog
      v-if="showDeleteConfirmation"
      :title="confirmationTitle"
      :message="confirmationMessage"
      confirmText="删除"
      @confirm="deleteSelectedAccounts"
      @cancel="handleDeletionCancel"
    />

    <!-- 关于信息弹窗 -->
    <div v-if="showAboutInfo" class="about-dialog-overlay" @click.self="showAboutInfo = false">
      <div class="about-dialog">
        <div class="about-dialog-header">
          <h3>关于 Euthenticator</h3>
          <button class="close-button" @click="showAboutInfo = false">&times;</button>
        </div>
        <div class="about-dialog-body">
          <div class="about-logo">
            <img src="../assets/images/logo-universal.png" alt="Euthenticator" class="about-logo-image">
          </div>
          <div class="about-version">
            <span class="version-label">版本</span>
            <span class="version-number">1.0.0</span>
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
                <a href="https://github.com/AmamiyaHotaru/Auth" target="_blank" class="detail-link">
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
          <button @click="showAboutInfo = false" class="about-close-button">
            关闭
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* 基础样式 */
#authenticator-container {
  display: flex;
  height: 100vh;
  background-color: #f8f9fa;
  color: #212529;
  overflow: hidden; /* 防止列表滚动时 body 也滚动 */
}

.sidebar {
  width: 250px;
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

.main-content {
  flex-grow: 1;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
}

.content-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px;
  background-color: #f8f9fa;
  border-bottom: 1px solid #dee2e6;
  position: sticky;
  top: 0;
  z-index: 10;
}

.selection-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px;
  background-color: #f8f9fa;
  border-bottom: 1px solid #dee2e6;
  position: sticky;
  top: 0;
  z-index: 10;
}

.header-left h2 {
  margin: 0;
  font-size: 1.5em;
}

.header-right {
  display: flex;
  gap: 10px;
}

.toolbar-button {
  display: flex;
  align-items: center;
  padding: 10px 15px;
  border: none;
  border-radius: 4px;
  background-color: #e9ecef;
  color: #212529;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.toolbar-button:hover {
  background-color: #dee2e6;
}

.toolbar-button.primary-button {
  background-color: #34A853;
  color: white;
}

.toolbar-button.primary-button:hover {
  background-color: #2c8a44;
}

.toolbar-button.danger-button {
  background-color: #EA4335;
  color: white;
}

.toolbar-button.danger-button:hover {
  background-color: #c9302c;
}

.toolbar-button.active {
  background-color: #e8f0fe;
  color: #4285F4;
}

.button-icon {
  font-size: 1.5em;
  margin-right: 5px;
}

.button-label {
  font-size: 0.9em;
}

.selection-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px;
  background-color: #f8f9fa;
  border-bottom: 1px solid #dee2e6;
}

.selection-count {
  font-size: 1.1em;
  margin-left: 15px;
}

.cancel-selection-btn {
  display: flex;
  align-items: center;
  gap: 5px;
  background: none;
  border: none;
  color: #212529;
  font-size: 1em;
  padding: 5px 10px;
  cursor: pointer;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.cancel-selection-btn:hover {
  background-color: #e9ecef;
}

.accounts-container {
  flex-grow: 1;
  padding: 15px;
}

.no-accounts {
  text-align: center;
  margin-top: 50px;
  color: #6c757d;
}

.desktop-button {
  display: inline-flex;
  align-items: center;
  padding: 10px 15px;
  border: none;
  border-radius: 4px;
  background-color: #34A853;
  color: white;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.desktop-button:hover {
  background-color: #2c8a44;
}

.accounts-list {
  display: flex;
  flex-wrap: wrap;
  gap: 15px;
}

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
  
  .sidebar {
    width: 200px;
  }
}

/* 关于弹窗样式 */
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