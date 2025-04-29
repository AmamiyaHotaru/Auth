<script setup>
import { ref, onMounted, onBeforeUnmount, watch } from 'vue';

const props = defineProps({
  show: Boolean,
  position: Object
});

const emit = defineEmits(['select', 'close']);

// 调试显示状态和位置
watch(() => props.show, (newValue) => {
  console.log('AddOptionsMenu show变化:', newValue);
  if (newValue) {
    console.log('AddOptionsMenu position:', props.position);
  }
});

// 处理选项选择
function handleSelect(option) {
  console.log('选择选项:', option);
  emit('select', option);
}

// 点击外部关闭菜单
function handleClickOutside(event) {
  // 添加阻止冒泡，避免组件在外层点击时立即被关闭
  event.stopPropagation();
  
  // 检查点击事件是否发生在菜单之外
  const optionsMenu = document.querySelector('.options-menu');
  if (props.show && optionsMenu && !optionsMenu.contains(event.target)) {
    console.log('点击在菜单外部，关闭菜单');
    emit('close');
  }
}

// 监听点击外部事件
onMounted(() => {
  console.log('AddOptionsMenu 组件已挂载');
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
      top: position?.top ? `${position.top}px` : '70px',
      left: position?.left ? `${position.left}px` : '20px'
    }"
    @click.stop
  >
    <div class="menu-item" @click="handleSelect('scan')">
      <div class="icon">📷</div>
      <div class="text">
        解析二维码
        <div class="subtext">支持谷歌身份验证器</div>
      </div>
    </div>
    <div class="menu-item" @click="handleSelect('manual')">
      <div class="icon">✏️</div>
      <div class="text">手动输入</div>
    </div>
  </div>
</template>

<style scoped>
.options-menu {
  position: absolute; /* 修改为 absolute 以避免可能的定位问题 */
  background-color: white;
  border-radius: 8px;
  width: 180px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  overflow: hidden;
  z-index: 1000; /* 增加 z-index 确保在其他元素之上 */
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

.subtext {
  font-size: 0.8em;
  color: #999;
  margin-top: 3px;
}

.menu-item:not(:last-child) {
  border-bottom: 1px solid #eee;
}
</style>