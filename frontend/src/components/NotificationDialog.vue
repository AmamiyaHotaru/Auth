<script setup>
import { defineProps, defineEmits } from 'vue';

const props = defineProps({
  show: {
    type: Boolean,
    default: false
  },
  type: {
    type: String,
    default: 'success', // 'success' | 'error' | 'warning' | 'info'
    validator: (value) => ['success', 'error', 'warning', 'info'].includes(value)
  },
  title: {
    type: String,
    default: '提示'
  },
  message: {
    type: String,
    required: true
  },
  buttonText: {
    type: String,
    default: '确定'
  }
});

const emit = defineEmits(['close']);

function handleClose() {
  emit('close');
}

// 根据类型获取不同的图标和颜色
const getTypeConfig = () => {
  switch (props.type) {
    case 'success':
      return {
        icon: '✓',
        iconBg: '#d4edda',
        iconColor: '#28a745',
        titleColor: '#28a745',
        animation: 'bounceIn'
      };
    case 'error':
      return {
        icon: '!',
        iconBg: '#f8d7da',
        iconColor: '#dc3545',
        titleColor: '#dc3545',
        animation: 'shake'
      };
    case 'warning':
      return {
        icon: '⚠',
        iconBg: '#fff3cd',
        iconColor: '#ffc107',
        titleColor: '#856404',
        animation: 'pulse'
      };
    case 'info':
    default:
      return {
        icon: 'i',
        iconBg: '#d1ecf1',
        iconColor: '#0dcaf0',
        titleColor: '#0c5460',
        animation: 'fadeIn'
      };
  }
};
</script>

<template>
  <div v-if="show" class="notification-dialog-overlay" @click.self="handleClose">
    <div class="notification-dialog-box">
      <div class="notification-icon-container">
        <div class="notification-icon" :style="{
          backgroundColor: getTypeConfig().iconBg,
          color: getTypeConfig().iconColor
        }">
          <span>{{ getTypeConfig().icon }}</span>
        </div>
      </div>
      <h3 class="notification-title" :style="{ color: getTypeConfig().titleColor }">
        {{ title }}
      </h3>
      <p class="notification-message">{{ message }}</p>
      <div class="notification-actions">
        <button @click="handleClose" class="btn btn-primary">{{ buttonText }}</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.notification-dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.6);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
  animation: fadeIn 0.2s ease;
}

.notification-dialog-box {
  background-color: white;
  padding: 25px;
  border-radius: 10px;
  box-shadow: 0 5px 20px rgba(0, 0, 0, 0.25);
  width: 85%;
  max-width: 400px;
  text-align: center;
  color: #333;
  animation: slideIn 0.3s ease;
}

.notification-icon-container {
  margin-bottom: 15px;
}

.notification-icon {
  display: inline-flex;
  justify-content: center;
  align-items: center;
  width: 60px;
  height: 60px;
  border-radius: 50%;
  font-size: 35px;
  font-weight: bold;
  animation: iconAnim 1s ease-in-out;
}

.notification-title {
  margin-top: 0;
  margin-bottom: 15px;
  font-size: 1.5em;
}

.notification-message {
  margin-bottom: 25px;
  font-size: 1.1em;
  color: #555;
  line-height: 1.5;
  word-break: break-word;
}

.notification-actions {
  display: flex;
  justify-content: center;
}

.btn {
  padding: 12px 24px;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  font-size: 1.1em;
  transition: all 0.2s ease;
  font-weight: 500;
}

.btn-primary {
  background-color: #4285F4;
  color: white;
  box-shadow: 0 2px 5px rgba(66, 133, 244, 0.3);
}

.btn-primary:hover {
  background-color: #3367d6;
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(66, 133, 244, 0.4);
}

.btn-primary:active {
  transform: translateY(0);
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes slideIn {
  from { transform: translateY(-20px); opacity: 0; }
  to { transform: translateY(0); opacity: 1; }
}

@keyframes iconAnim {
  0% { transform: scale(0); }
  50% { transform: scale(1.1); }
  100% { transform: scale(1); }
}

@keyframes shake {
  0%, 100% { transform: translateX(0); }
  10%, 30%, 50%, 70%, 90% { transform: translateX(-5px); }
  20%, 40%, 60%, 80% { transform: translateX(5px); }
}

@keyframes pulse {
  0% { transform: scale(1); }
  50% { transform: scale(1.05); }
  100% { transform: scale(1); }
}

@keyframes bounceIn {
  0% { transform: scale(0.3); opacity: 0; }
  50% { transform: scale(1.05); opacity: 1; }
  70% { transform: scale(0.9); }
  100% { transform: scale(1); }
}
</style>