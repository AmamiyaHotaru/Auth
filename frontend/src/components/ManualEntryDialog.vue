<script setup>
import { ref } from 'vue';

const props = defineProps({
  show: Boolean
});

const emit = defineEmits(['confirm', 'cancel']);

// 表单数据
const accountName = ref('');
const serverName = ref('');
const secret = ref('');
const accountType = ref(0); // 默认为 TOTP 类型 (0)

// 验证状态
const formError = ref('');

// 处理表单提交
function handleSubmit() {
  // 简单验证
  if (!accountName.value.trim()) {
    formError.value = '请输入账户名称';
    return;
  }
  
  if (!serverName.value.trim()) {
    formError.value = '请输入服务商名称';
    return;
  }
  
  if (!secret.value.trim()) {
    formError.value = '请输入密钥';
    return;
  }
  
  // 提交数据
  emit('confirm', {
    accountName: accountName.value.trim(),
    serverName: serverName.value.trim(),
    secret: secret.value.trim(),
    accountType: parseInt(accountType.value) // 确保是数字类型
  });
  
  // 重置表单
  resetForm();
}

// 取消操作
function handleCancel() {
  resetForm();
  emit('cancel');
}

// 重置表单
function resetForm() {
  accountName.value = '';
  serverName.value = '';
  secret.value = '';
  accountType.value = 0;
  formError.value = '';
}

</script>

<template>
  <div v-if="show" class="dialog-overlay">
    <div class="dialog-container">
      <div class="dialog-header">
        <h3>手动添加账户</h3>
      </div>
      
      <div class="dialog-body">
        <div v-if="formError" class="error-message">{{ formError }}</div>
        
        <div class="form-group">
          <label for="accountName">账户名称</label>
          <input 
            type="text" 
            id="accountName" 
            v-model="accountName"
            placeholder="例如：username@example.com"
            autofocus
          />
        </div>
        
        <div class="form-group">
          <label for="serverName">服务商名称</label>
          <input 
            type="text" 
            id="serverName" 
            v-model="serverName"
            placeholder="例如：Google、GitHub、Steam"
          />
        </div>
        
        <div class="form-group">
          <label for="secret">密钥</label>
          <input 
            type="text" 
            id="secret" 
            v-model="secret"
            placeholder="请输入密钥字符串"
          />
        </div>
        
        <div class="form-group">
          <label for="accountType">账户类型</label>
          <select id="accountType" v-model="accountType" class="account-type-select">
            <option value="0">TOTP（常规验证码）</option>
            <option value="1">Steam</option>
          </select>
        </div>
      </div>
      
      <div class="dialog-footer">
        <button @click="handleCancel" class="btn-cancel">取消</button>
        <button @click="handleSubmit" class="btn-confirm">添加</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.dialog-container {
  background-color: white;
  border-radius: 8px;
  width: 90%;
  max-width: 400px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
  overflow: hidden;
}

.dialog-header {
  padding: 15px;
  background-color: #4285F4;
  color: white;
  text-align: center;
}

.dialog-header h3 {
  margin: 0;
  font-size: 1.2em;
}

.dialog-body {
  padding: 20px;
}

.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: 500;
  color: #333;
}

.form-group input, .form-group select {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 16px;
  box-sizing: border-box;
}

.form-group input:focus, .form-group select:focus {
  border-color: #4285F4;
  outline: none;
  box-shadow: 0 0 0 2px rgba(66, 133, 244, 0.2);
}

.account-type-select {
  background-color: white;
  cursor: pointer;
}

.error-message {
  color: #d9534f;
  margin-bottom: 15px;
  padding: 8px 12px;
  background-color: rgba(217, 83, 79, 0.1);
  border-left: 3px solid #d9534f;
  border-radius: 2px;
}

.dialog-footer {
  padding: 10px 20px;
  display: flex;
  justify-content: flex-end;
  background-color: #f5f5f5;
  border-top: 1px solid #ddd;
}

.btn-cancel, .btn-confirm {
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
  margin-left: 10px;
}

.btn-cancel {
  background-color: #f5f5f5;
  color: #333;
  border: 1px solid #ddd;
}

.btn-confirm {
  background-color: #4285F4;
  color: white;
}

.btn-cancel:hover {
  background-color: #e5e5e5;
}

.btn-confirm:hover {
  background-color: #3b78e7;
}
</style>