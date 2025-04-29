<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue';

// 导入组件
import ConfirmationDialog from './ConfirmationDialog.vue';
import ManualEntryDialog from './ManualEntryDialog.vue';
import AddOptionsMenu from './AddOptionsMenu.vue';
import QrCodeUploadDialog from './QrCodeUploadDialog.vue';
import Alert from './Alert.vue';

// 导入模块化组件
import TitleBar from './TitleBar.vue';
import SideBar from './SideBar.vue';
import AccountCard from './AccountCard.vue';
import AboutDialog from './AboutDialog.vue';

// 导入服务
import * as accountService from './accountService.js';
import * as timerService from './timerService.js';
import * as codeUtils from './codeUtils.js';
import * as selectionModeService from './selectionMode.js';

// --- 状态 ---
const accounts = ref([]);
// 存储显示验证码的账户ID
const visibleCodes = ref(new Set());
// 用于追踪复制按钮的状态
const copiedCodes = ref(new Set());
// 全局显示/隐藏所有验证码的状态
const showAllCodes = ref(false);

// --- 编辑对话框状态 ---
const showEditDialog = ref(false);
const editingAccount = ref(null);

// --- 选择模式相关状态 ---
const selectionMode = ref(false);
const selectedAccountIds = ref(new Set());
const selectedCount = computed(() => selectedAccountIds.value.size);
const ignoreNextClick = ref(false);

// --- 确认对话框状态 ---
const showDeleteConfirmation = ref(false);
const confirmationTitle = ref('');
const confirmationMessage = ref('');

// --- 添加账户菜单和对话框状态 ---
const showAddOptions = ref(false);
const showManualEntryDialog = ref(false);
const showQrCodeDialog = ref(false);
const qrCodeDialogRef = ref(null);
const addButtonPosition = ref(null);
const qrCodeImageData = ref(null);

// --- 关于按钮状态 ---
const showAboutInfo = ref(false);

// --- Alert 组件全局提示 ---
const alertRef = ref(null);

// --- 其他状态和变量 ---
let timerInterval = null; // TOTP 计时器
let longPressTimer = null; // 长按检测计时器
const longPressDuration = 500; // 长按所需毫秒数

// --- 验证码显示和复制功能 ---

// 显示或隐藏验证码
function toggleCodeVisibility(accountId, event) {
  // 添加对 event 的空值检查
  if (event) {
    event.stopPropagation(); // 阻止事件冒泡
  }
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
  
 
  const success = await codeUtils.copyToClipboard(code);
  console.log('复制结果:', success);
  
  // 添加短暂的"已复制"状态
  copiedCodes.value.add(accountId);
  setTimeout(() => {
    copiedCodes.value.delete(accountId);
  }, 1500);
}

// --- 账户列表管理 ---

// 获取账户列表
const getSecretsList = async () => {
  try {
    const response = await accountService.getSecretsList();
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
  const currentTimeLeft = timerService.calculateTimeLeft();
  
  accounts.value = accountService.processAccountList(list, currentTimeLeft);
  console.log('已更新账户列表，数量:', accounts.value.length);
}

// --- 选择模式相关功能 ---

// 进入选择模式
function enterSelectionMode(selectedId) {
  selectionModeService.enterSelectionMode(selectedId, selectionMode, selectedAccountIds);
}

// 退出选择模式
function cancelSelectionMode() {
  selectionModeService.cancelSelectionMode(selectionMode, selectedAccountIds);
}

// 切换账户的选择状态
function toggleSelection(id) {
  selectionModeService.toggleSelection(id, selectedAccountIds);
}

// 请求删除选中的账户
function requestDeleteSelected() {
  const result = selectionModeService.requestDeleteSelected({ value: selectedCount.value }, selectedAccountIds);
  if (result) {
    confirmationTitle.value = result.title;
    confirmationMessage.value = result.message;
    showDeleteConfirmation.value = true;
  }
}

// 删除选中的账户
async function deleteSelectedAccounts() {
  try {
    console.log('删除账户:', [...selectedAccountIds.value]);
    
    // 将选中的账户ID转换为数组
    const idArray = Array.from(selectedAccountIds.value);
    
    await accountService.deleteSecret(idArray);
    
    // 刷新账户列表
    await getSecretsList();
    
    // 关闭确认对话框
    showDeleteConfirmation.value = false;
    
    // 退出选择模式
    cancelSelectionMode();
    
    if (alertRef.value) alertRef.value.show('success', '删除账户成功');
    
  } catch (error) {
    console.error('删除账户失败:', error);
    if (alertRef.value) alertRef.value.show('error', '删除账户失败');
  }
}

// 取消删除操作
function handleDeletionCancel() {
  showDeleteConfirmation.value = false;
}

// 编辑选中的账户
function editSelectedAccount() {
  if (selectedCount.value !== 1) {
    if (alertRef.value) {
      alertRef.value.show('warning', '请只选择一个令牌进行编辑');
    }
    return;
  }
  
  // 获取选中的账户ID
  const selectedId = [...selectedAccountIds.value][0];
  
  // 查找对应的账户对象
  const accountToEdit = accounts.value.find(acc => acc.ID === selectedId);
  
  if (accountToEdit) {
    editingAccount.value = {
      ...accountToEdit,
      id: accountToEdit.ID // 确保ID属性一致
    };
    
    showEditDialog.value = true;
  }
}

// --- 条目交互功能 ---

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
  }, longPressDuration);
}

// 处理条目交互结束（鼠标松开或触摸结束）
function handleItemInteractionEnd() {
  // 清除长按计时器
  clearTimeout(longPressTimer);
}

// 处理条目点击
function handleItemClick(account) {
  // 如果需要忽略这次点击（因为已经触发了长按），则不执行操作
  if (ignoreNextClick.value) {
    ignoreNextClick.value = false; // 重置标志
    return;
  }

  // 如果已处于选择模式，则切换该条目的选择状态
  if (selectionMode.value) {
    toggleSelection(account.ID);
  }
  // 否则可以执行其他操作
  else {
    // 直接切换代码可见性，不需要stopPropagation
    if (visibleCodes.value.has(account.ID)) {
      visibleCodes.value.delete(account.ID);
    } else {
      visibleCodes.value.add(account.ID);
    }
  }
}

// --- 添加账户相关逻辑 ---

// 显示添加选项菜单
function showAddOptionsMenu(event) {

  
  // 计算按钮位置
  const rect = event.currentTarget.getBoundingClientRect();

  
  // 阻止事件冒泡，防止菜单立即关闭
  event.stopPropagation();
  
  // 获取窗口宽度
  const windowWidth = window.innerWidth;
  
  // 计算菜单合适的水平位置
  // 菜单宽度约为180px，添加一些边距
  const menuWidth = 180;
  const margin = 10;
  
  // 如果按钮靠近右边缘，则将菜单向左侧显示
  let left = rect.left;
  if (left + menuWidth + margin > windowWidth) {
    left = rect.right - menuWidth;
  }
  
  // 设置菜单位置 - 在按钮下方
  addButtonPosition.value = {
    left: left,
    top: rect.bottom + 5 // 添加5px的间距
  };


  // 显示菜单
  showAddOptions.value = true;

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
    await accountService.addManualSecret(
      formData.accountName, 
      formData.serverName, 
      formData.secret, 
      formData.accountType
    );
    
    // 刷新账户列表
    await getSecretsList();
    
    // 关闭对话框
    showManualEntryDialog.value = false;
    
    if (alertRef.value) alertRef.value.show('success', '添加账户成功');
    
  } catch (error) {
    console.error('添加账户失败:', error);
    if (alertRef.value) alertRef.value.show('error', '添加账户失败');
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
    await accountService.addQRCodeSecret(data);
    
    // 检测成功后刷新账户列表
    await getSecretsList();
    
    // 重置处理状态
    if (qrCodeDialogRef.value) {
      qrCodeDialogRef.value.resetProcessingState();
    }
    
    // 关闭对话框
    showQrCodeDialog.value = false;

    if (alertRef.value) alertRef.value.show('success', '二维码解析并添加账户成功');
    
  } catch (error) {
    console.error('二维码识别失败:', error);
    
    // 设置错误状态
    if (qrCodeDialogRef.value) {
      qrCodeDialogRef.value.setProcessingError(error);
    }

    if (alertRef.value) alertRef.value.show('error', '二维码解析或添加账户失败');
    
    // 注意：不关闭对话框，让用户可以重试
  }
}

// 编辑账户
async function handleEditAccount(formData) {
  try {
    console.log('编辑账户数据:', formData);
    
    if (!formData.ID) {
      console.error('缺少账户ID，无法编辑');
      if (alertRef.value) alertRef.value.show('error', '缺少账户ID，无法编辑');
      return;
    }
    
    // 调用Go后端更新账户
    await accountService.updateSecret(
      formData.ID,
      formData.accountName, 
      formData.serverName, 
      formData.accountType
    );
    
    // 刷新账户列表
    await getSecretsList();
    
    // 关闭编辑对话框
    closeEditDialog();
    
    // 退出选择模式
    cancelSelectionMode();
    
    if (alertRef.value) alertRef.value.show('success', '编辑账户成功');
    
  } catch (error) {
    console.error('编辑账户失败:', error);
    if (alertRef.value) alertRef.value.show('error', '编辑账户失败');
  }
}

// 关闭编辑对话框
function closeEditDialog() {
  showEditDialog.value = false;
  editingAccount.value = null;
}

// --- TOTP 计时器逻辑 ---

// 每 30 秒更新一次验证码
async function updateCodes() {
  console.log("30 秒周期结束，获取新验证码...");
  await getSecretsList(); 
}

// 处理计时器每秒的 tick
function handleTimerTick() {
  const newTimeLeft = timerService.calculateTimeLeft();
    
  // 检测是否需要更新验证码（如果 timeLeft 从 1 变为 30，说明过了一个周期）
  const needsUpdate = timerService.needsCodeRefresh(accounts.value, newTimeLeft);
    
  // 更新所有账户的 timeLeft
  accounts.value = timerService.updateAccountsTimeLeft(accounts.value, newTimeLeft);

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

// 禁止Tab切换焦点
function handleTabKey(e) {
  if (e.key === 'Tab') {
    e.preventDefault();
    e.stopPropagation();
    return false;
  }
}

// --- 生命周期钩子 ---
onMounted(async () => {
  console.log('组件已挂载，获取初始数据...');
  await getSecretsList(); // 获取初始数据
  startTimerInterval(); // 启动计时器
  window.addEventListener('keydown', handleTabKey, true);
});

onUnmounted(() => {
  if (timerInterval) clearInterval(timerInterval); // 组件卸载时清除计时器
  clearTimeout(longPressTimer); // 清除可能存在的长按计时器
  window.removeEventListener('keydown', handleTabKey, true);
});
</script>

<template>
  <div id="authenticator-container" :class="{ 'selection-mode': selectionMode }">
    <!-- 窗口标题栏 -->
    <TitleBar />

    <div class="app-container">
      <!-- 侧边栏 -->
      <SideBar 
        :accountsCount="accounts.length"
        @show-about="showAboutInfo = true"
      />

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
            <button @click="editSelectedAccount" :disabled="selectedCount !== 1" :class="['toolbar-button', { 'toolbar-button--disabled': selectedCount !== 1 }]">
              <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4Z"></path>
              </svg>
              <span class="button-label">编辑选中</span>
            </button>
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
            <AccountCard
              v-for="account in accounts" 
              :key="account.ID"
              :account="{
                ...account,
                isCopied: copiedCodes.has(account.ID)
              }"
              :selectionMode="selectionMode"
              :isSelected="selectedAccountIds.has(account.ID)"
              :isCodeVisible="isCodeVisible(account.ID)"
              @toggle-selection="toggleSelection"
              @toggle-code-visibility="toggleCodeVisibility"
              @copy-code="copyCodeToClipboard"
              @interaction-start="handleItemInteractionStart"
              @interaction-end="handleItemInteractionEnd"
              @click="handleItemClick"
            />
          </div>
        </div>
      </main>
    </div>

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

    <!-- 编辑对话框 -->
    <ManualEntryDialog
      v-if="showEditDialog"
      :show="showEditDialog"
      :isEditing="true"
      :accountData="editingAccount"
      @confirm="handleEditAccount"
      @cancel="closeEditDialog"
    />
    
    <!-- 关于信息弹窗 -->
    <AboutDialog
      :show="showAboutInfo"
      @close="showAboutInfo = false"
    />

    <!-- Alert 组件 -->
    <Alert ref="alertRef" />
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

.app-container {
  display: flex;
  flex-grow: 1;
  overflow: hidden;
}

.main-content {
  flex-grow: 1;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
  scrollbar-width: none; /* Firefox */
  -ms-overflow-style: none; /* IE and Edge */
}

.main-content::-webkit-scrollbar {
  display: none; /* Chrome, Safari, Opera */
}

.accounts-container {
  flex-grow: 1;
  padding: 15px;
  scrollbar-width: none;
  -ms-overflow-style: none;
}

.accounts-container::-webkit-scrollbar {
  display: none;
}

/* 防止整个文档的滚动条出现 */
:root, body {
  scrollbar-width: none;
  -ms-overflow-style: none;
}

:root::-webkit-scrollbar, body::-webkit-scrollbar {
  display: none;
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

.toolbar-button--disabled {
  background-color: #e0e0e0 !important;
  color: #bdbdbd !important;
  cursor: not-allowed !important;
  pointer-events: auto !important;
  /* 保证即使禁用也能触发点击事件 */
  opacity: 1 !important;
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
</style>