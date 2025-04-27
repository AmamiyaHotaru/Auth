<script setup>
import { defineProps, defineEmits } from 'vue';

const props = defineProps({
  title: {
    type: String,
    default: '确认操作'
  },
  message: {
    type: String,
    required: true
  },
  confirmText: {
    type: String,
    default: '确认'
  },
  cancelText: {
    type: String,
    default: '取消'
  }
});

const emit = defineEmits(['confirm', 'cancel']);

function handleConfirm() {
  emit('confirm');
}

function handleCancel() {
  emit('cancel');
}
</script>

<template>
  <div class="dialog-overlay" @click.self="handleCancel"> <!-- Close on overlay click -->
    <div class="dialog-box">
      <h3 class="dialog-title">{{ title }}</h3>
      <p class="dialog-message">{{ message }}</p>
      <div class="dialog-actions">
        <button @click="handleCancel" class="btn btn-cancel">{{ cancelText }}</button>
        <button @click="handleConfirm" class="btn btn-confirm">{{ confirmText }}</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.6); /* Dim background */
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000; /* Ensure it's on top */
}

.dialog-box {
  background-color: white;
  padding: 25px;
  border-radius: 8px;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.3);
  width: 80%;
  max-width: 400px; /* Limit maximum width */
  text-align: center;
  color: #333; /* Text color inside dialog */
}

.dialog-title {
  margin-top: 0;
  margin-bottom: 15px;
  font-size: 1.3em;
  color: #343a40;
}

.dialog-message {
  margin-bottom: 25px;
  font-size: 1em;
  color: #6c757d;
}

.dialog-actions {
  display: flex;
  justify-content: space-around; /* Space out buttons */
}

.btn {
  padding: 10px 20px;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  font-size: 1em;
  transition: background-color 0.2s ease;
}

.btn-confirm {
  background-color: #dc3545; /* Red for confirm (delete) */
  color: white;
}

.btn-confirm:hover {
  background-color: #c82333;
}

.btn-cancel {
  background-color: #6c757d; /* Grey for cancel */
  color: white;
}

.btn-cancel:hover {
  background-color: #5a6268;
}
</style>