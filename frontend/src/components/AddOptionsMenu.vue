<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue';

const props = defineProps({
  show: Boolean,
  position: Object
});

const emit = defineEmits(['select', 'close']);

// 处理选项选择
function handleSelect(option) {
  emit('select', option);
}

// 点击外部关闭菜单
function handleClickOutside(event) {
  if (props.show) {
    emit('close');
  }
}

// 监听点击外部事件
onMounted(() => {
  document.addEventListener('click', handleClickOutside);
});

onBeforeUnmount(() => {
  document.removeEventListener('click', handleClickOutside);
});
</script>

<template>
  <div 
    v-if="show" 
    class="options-menu"
    :style="{ 
      bottom: `${position?.bottom || 70}px`,
      right: `${position?.right || 20}px`
    }"
    @click.stop
  >
    <div class="menu-item" @click="handleSelect('scan')">
      <div class="icon">📷</div>
      <div class="text">解析图片</div>
    </div>
    <div class="menu-item" @click="handleSelect('manual')">
      <div class="icon">✏️</div>
      <div class="text">手动输入</div>
    </div>
  </div>
</template>

<style scoped>
.options-menu {
  position: fixed;
  background-color: white;
  border-radius: 8px;
  width: 180px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  overflow: hidden;
  z-index: 990;
}

.menu-item {
  display: flex;
  align-items: center;
  padding: 15px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.menu-item:hover {
  background-color: #f5f5f5;
}

.menu-item .icon {
  margin-right: 12px;
  font-size: 1.2em;
}

.menu-item .text {
  font-size: 0.95em;
}

.menu-item:not(:last-child) {
  border-bottom: 1px solid #eee;
}
</style>