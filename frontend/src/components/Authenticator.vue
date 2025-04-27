<script setup>
import { ref, onMounted, computed, onUnmounted, nextTick } from 'vue';
import {GetSecretsList, InsertSecret, DeleteSecret, RecognizeQRCode} from '../../wailsjs/go/main/App';
import ConfirmationDialog from './ConfirmationDialog.vue';
import ManualEntryDialog from './ManualEntryDialog.vue';
import AddOptionsMenu from './AddOptionsMenu.vue';
import QrCodeUploadDialog from './QrCodeUploadDialog.vue';
import NotificationDialog from './NotificationDialog.vue';

// --- 状态 ---
const accounts = ref([]);

// --- 辅助函数 ---
// 计算当前时间在 30 秒周期内的剩余秒数
function calculateTimeLeft() {
    const nowSeconds = Math.floor(Date.now() / 1000);
    return 30 - (nowSeconds % 30);
}

// --- 通知弹窗状态 ---
const notification = ref({
  show: false,
  type: 'info',
  title: '',
  message: '',
  buttonText: '确定'
});

// 显示通知弹窗
function showNotification(type, title, message, buttonText = '确定') {
  notification.value = {
    show: true,
    type,
    title,
    message,
    buttonText
  };
}

// 关闭通知弹窗
function closeNotification() {
  notification.value.show = false;
}

// 获取账户列表
const getSecretsList = async () => {
  try {
    const response = await GetSecretsList();
    
    processAccountList(response);

  } catch (error) {
    console.error('获取账户列表时出错:', error);
    accounts.value = [];
    showNotification('error', '错误', '获取账户列表时出错');
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
    showNotification('error', '错误', `删除账户时出错: ${error}`);
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
  
  // 调整位置计算，确保菜单在按钮上方
  addButtonPosition.value = {
    bottom: window.innerHeight - rect.top,
    right: window.innerWidth - rect.right
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
    showNotification('error', '错误', `添加账户失败: ${error}`);
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
    
    // 显示成功提示
    showNotification('success', '成功', '二维码识别成功，账户已添加');
    
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
    
    // 显示错误通知
    showNotification('error', '错误', `二维码识别失败: ${error}`);
    
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
  <div id="authenticator-container">
    <!-- 条件渲染的头部 -->
    <header v-if="!selectionMode" class="app-header">
      <h1>身份验证器</h1>
    </header>
    <header v-else class="selection-header">
      <button @click="cancelSelectionMode" class="cancel-selection-btn">&lt;</button>
      <span class="selection-count">{{ selectedCount }} 项已选中</span>
      <button @click="requestDeleteSelected" :disabled="selectedCount === 0" class="delete-selection-btn">删除</button>
    </header>

    <main class="accounts-list">
      <!-- 没有账户时的提示 -->
      <div v-if="accounts.length === 0 && !selectionMode" class="no-accounts">
        <p>还没有添加账户</p>
      </div>
      <!-- 账户列表项 -->
      <div v-else v-for="account in accounts" :key="account.ID"
           class="account-item"
           :class="{ 'selected': selectedAccountIds.has(account.ID), 'selection-active': selectionMode }"
           @mousedown="handleItemInteractionStart($event, account)"
           @mouseup="handleItemInteractionEnd($event, account)"
           @touchstart.passive="handleItemInteractionStart($event, account)"
           @touchend="handleItemInteractionEnd($event, account)"
           @click="handleItemClick(account)"
           @contextmenu.prevent> <!-- 阻止条目上的右键菜单 -->

        <!-- 选择指示器 (可选的复选框样式) -->
         <div v-if="selectionMode" class="selection-indicator">
             <input type="checkbox" :checked="selectedAccountIds.has(account.ID)" @click.stop="toggleSelection(account.ID)" class="selection-checkbox">
         </div>

        <!-- 左侧：账户详情 -->
        <div class="account-details">
           <div class="account-info">
             <span class="issuer">{{ account.ServerName }}</span>
             <span class="name">{{ account.AccountName }}</span>
           </div>
           <div class="code-line">
             <span class="code">{{ account.Code }}</span>
           </div>
        </div>
        <!-- 右侧：进度条 (无论是否在选择模式下都显示) -->
        <div class="progress-container">
           <div class="progress-circle" :style="getProgressStyle(account.timeLeft)"></div>
        </div>
      </div>
    </main>

    <!-- 条件渲染的底部 -->
    <footer v-if="!selectionMode" class="app-footer">
       <button @click="showAddOptionsMenu" class="add-button">+</button>
    </footer>

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

    <!-- 通知弹窗 -->
    <NotificationDialog
      v-if="notification.show"
      :type="notification.type"
      :title="notification.title"
      :message="notification.message"
      :buttonText="notification.buttonText"
      @close="closeNotification"
    />

  </div>
</template>

<style scoped>
/* 基础样式 */
#authenticator-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background-color: #f8f9fa;
  color: #212529;
  overflow: hidden; /* 防止列表滚动时 body 也滚动 */
}

/* 标准头部 */
.app-header {
  background-color: #4285F4;
  color: white;
  padding: 15px;
  text-align: center;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  flex-shrink: 0;
}
.app-header h1 { margin: 0; font-size: 1.5em; }

/* 选择模式头部 */
.selection-header {
  background-color: #6c757d; /* 选择模式下使用灰色头部 */
  color: white;
  padding: 10px 15px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  flex-shrink: 0;
  height: 58px; /* 大致匹配普通头部的高度 */
}
.selection-count { font-size: 1.1em; }
.cancel-selection-btn, .delete-selection-btn {
  background: none;
  border: none;
  color: white;
  font-size: 1em;
  padding: 5px 10px;
  cursor: pointer;
}
.cancel-selection-btn { font-size: 1.5em; padding: 0 10px;} /* 让返回箭头更大 */
.delete-selection-btn:disabled { color: #adb5bd; cursor: not-allowed; }

/* 账户列表 */
.accounts-list {
  flex-grow: 1;
  overflow-y: auto;
  padding: 15px;
}
.no-accounts {
  text-align: center;
  margin-top: 50px;
  color: #6c757d;
}

/* 账户条目 */
.account-item {
  background-color: white;
  border-radius: 8px;
  margin-bottom: 15px;
  padding: 15px;
  display: flex; /* 使用 flexbox 对齐详情和进度条 */
  justify-content: space-between; /* 详情和进度条之间留空 */
  align-items: center; /* 垂直对齐项目 */
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  transition: box-shadow 0.2s ease-in-out, background-color 0.2s ease;
  position: relative; /* 如果选择删除按钮样式，则需要绝对定位 */
  user-select: none; /* 防止交互时选中文本 */
}
.account-item:hover {
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
}

/* 选择指示器 */
.selection-indicator {
    margin-right: 15px; /* 复选框与详情之间的间距 */
    flex-shrink: 0;
}
.selection-checkbox {
    width: 20px;
    height: 20px;
    cursor: pointer;
}

/* 选中状态 */
.account-item.selected {
  background-color: #e0e0e0; /* 选中项的浅灰色背景 */
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.12);
}
/* 选择模式激活时的悬停效果略有不同 */
.account-item.selection-active:hover {
    background-color: #d5d5d5;
}

/* 账户详情 */
.account-details {
    display: flex;
    flex-direction: column; /* 堆叠信息行和代码行 */
    flex-grow: 1; /* 允许详情占用可用空间 */
    margin-right: 10px; /* 减少间距 */
    min-width: 0; /* 防止 flex 项目溢出 */
    pointer-events: none; /* 防止选择模式下的指针事件 */
}
.account-item.selection-active .account-details {
    pointer-events: none;
}
.account-item:not(.selection-active) .account-details {
    pointer-events: auto; /* 非选择模式下重新启用指针事件 */
}

.account-info {
  margin-bottom: 8px; /* 第一行下方的间距 */
  white-space: nowrap; /* 防止换行 */
  overflow: hidden; /* 隐藏溢出 */
  text-overflow: ellipsis; /* 如果文本过长，添加省略号 */
}
.issuer {
  font-weight: bold;
  font-size: 1.2em; /* 更大、更粗的发行者 */
  color: #343a40;
  margin-right: 8px; /* 发行者和名称之间的间距 */
}
.name {
  font-size: 0.95em; /* 稍大的名称 */
  color: #6c757d;
}
.code-line {
 /* 容器 */
}
.code {
  font-size: 2.2em; /* 更大的代码 */
  font-weight: bold;
  color: #0d47a1; /* 深蓝色 */
  letter-spacing: 3px; /* 更大的代码间距 */
  font-family: 'Courier New', Courier, monospace; /* 等宽字体 */
}

/* 进度条 */
.progress-container {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px; /* 固定宽度 */
  height: 40px; /* 固定高度 */
  flex-shrink: 0; /* 防止容器缩小 */
  position: relative; /* 用于内部的绝对定位 */
  margin-right: 10px; /* 进度条和删除按钮之间的间距 */
}
.progress-circle {
  width: 36px; /* 圆的直径 */
  height: 36px;
  border-radius: 50%; /* 使其圆形 */
  /* 背景通过 :style 动态设置 */
  transition: background 0.1s linear; /* 平滑过渡计时器更新 */
  display: flex; /* 居中内部内容（如果有） */
  align-items: center;
  justify-content: center;
}

/* 底部 */
.app-footer {
  padding: 15px;
  text-align: center;
  border-top: 1px solid #dee2e6;
  flex-shrink: 0;
}
.add-button {
  background-color: #34A853;
  color: white;
  border: none;
  border-radius: 50%;
  width: 50px;
  height: 50px;
  font-size: 2em;
  line-height: 48px;
  cursor: pointer;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
  transition: background-color 0.2s ease;
}
.add-button:hover {
  background-color: #2c8a44;
}

/* 移除旧的删除按钮样式（如果存在） */
.delete-button { display: none; }

</style>