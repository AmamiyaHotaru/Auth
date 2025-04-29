<template>
  <transition name="alert-fade">
    <div v-if="visible" :class="['custom-alert', type]">
      <span class="alert-icon" v-if="icon">{{ icon }}</span>
      <span class="alert-message">{{ message }}</span>
      <span class="alert-close" @click="close">×</span>
    </div>
  </transition>
</template>

<script setup>
import { ref, watch, computed } from 'vue';

const props = defineProps({
  type: { type: String, default: 'info' }, // success/warning/error/info
  message: { type: String, default: '' },
  duration: { type: Number, default: 2500 }, // ms
  show: { type: Boolean, default: false },
});

const visible = ref(props.show);
const message = ref(props.message);
const type = ref(props.type);
let timer = null;

const iconMap = {
  success: '✔',
  warning: '⚠',
  error: '✖',
  info: 'ℹ',
};
const icon = computed(() => iconMap[type.value] || 'ℹ');

function show(newType, newMsg, newDuration) {
  type.value = newType || 'info';
  message.value = newMsg || '';
  visible.value = true;
  if (timer) clearTimeout(timer);
  timer = setTimeout(() => {
    visible.value = false;
  }, newDuration || props.duration);
}

function close() {
  visible.value = false;
  if (timer) clearTimeout(timer);
}

// 支持外部通过ref调用show
defineExpose({ show, close });

// 响应props变化
watch(() => props.show, (val) => { visible.value = val; });
watch(() => props.message, (val) => { message.value = val; });
watch(() => props.type, (val) => { type.value = val; });
</script>

<style scoped>
.custom-alert {
  position: fixed;
  top: 30px;
  left: 50%;
  transform: translateX(-50%);
  min-width: 220px;
  max-width: 80vw;
  padding: 12px 32px 12px 18px;
  border-radius: 4px;
  color: #fff;
  font-size: 15px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.12);
  display: flex;
  align-items: center;
  z-index: 9999;
  opacity: 0.98;
}
.custom-alert .alert-icon {
  margin-right: 8px;
  font-size: 18px;
}
.custom-alert .alert-message {
  flex: 1;
}
.custom-alert .alert-close {
  margin-left: 12px;
  cursor: pointer;
  font-size: 18px;
  opacity: 0.7;
}
.custom-alert.success { background: #67c23a; }
.custom-alert.warning { background: #e6a23c; }
.custom-alert.error { background: #f56c6c; }
.custom-alert.info { background: #409eff; }

.alert-fade-enter-active, .alert-fade-leave-active {
  transition: opacity 0.3s;
}
.alert-fade-enter-from, .alert-fade-leave-to {
  opacity: 0;
}
</style>
