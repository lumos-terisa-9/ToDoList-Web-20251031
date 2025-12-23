<template>
  <div class="activity-detail-member">
    <!-- 公共部分 -->
    <ActivityDetailCommon
        ref="commonComponent"
        :activity="activity"
        :org-id="orgId"
        :user-role="userRole"
        @close="$emit('close')"
        @review-submitted="handleReviewSubmitted"
        @review-failed="handleReviewFailed"
    />

    <!-- 成功/错误提示 -->
    <div v-if="showSuccessMessage" class="success-message">
      <div class="success-icon">✅</div>
      <div class="success-text">{{ successMessage }}</div>
    </div>

    <div v-if="showErrorMessage" class="error-message">
      <div class="error-icon">❌</div>
      <div class="error-text">{{ errorMessage }}</div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import ActivityDetailCommon from './ActivityDetailCommon.vue'

const props = defineProps({
  activity: {
    type: Object,
    required: true
  },
  orgId: {
    type: [String, Number],
    required: true
  },
  orgName: {
    type: String,
    required: true
  }
})

const emit = defineEmits(['close'])

// 引用
const commonComponent = ref(null)
const userRole = 'Member'

// 提示消息
const showSuccessMessage = ref(false)
const successMessage = ref('')
const showErrorMessage = ref(false)
const errorMessage = ref('')

// 处理评价提交成功
function handleReviewSubmitted() {
  showSuccess('评价提交成功')
}

// 处理评价提交失败
function handleReviewFailed(errorMessage) {
  // 如果错误信息包含"参与"，显示特定提示
  if (errorMessage.includes('参与')) {
    showError('请先参与活动再进行评价')
  } else {
    showError(errorMessage)
  }
}

// 显示成功消息
function showSuccess(message) {
  successMessage.value = message
  showSuccessMessage.value = true

  // 如果已经有错误消息，先关闭它
  if (showErrorMessage.value) {
    showErrorMessage.value = false
  }

  setTimeout(() => {
    showSuccessMessage.value = false
  }, 3000)
}

// 显示失败消息
function showError(message) {
  errorMessage.value = message
  showErrorMessage.value = true

  // 如果已经有成功消息，先关闭它
  if (showSuccessMessage.value) {
    showSuccessMessage.value = false
  }

  setTimeout(() => {
    showErrorMessage.value = false
  }, 3000)
}
</script>

<style scoped>
.activity-detail-member {
  position: relative;
}

/* 成功消息 */
.success-message {
  position: fixed;
  top: 84px;
  right: 24px;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px 24px;
  background: rgba(100, 200, 100, 0.15);
  border: 1px solid rgba(100, 200, 100, 0.3);
  border-radius: 12px;
  backdrop-filter: blur(10px);
  z-index: 5000;
  animation: slideIn 0.3s ease, fadeOut 0.3s ease 2.7s forwards;
}

.success-icon {
  font-size: 20px;
}

.success-text {
  font-weight: 600;
  color: #64c864;
}

.error-message {
  position: fixed;
  top: 84px;
  right: 24px;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px 24px;
  background: rgba(255, 100, 100, 0.15);
  border: 1px solid rgba(255, 100, 100, 0.3);
  border-radius: 12px;
  backdrop-filter: blur(10px);
  z-index: 5000;
  animation: slideIn 0.3s ease, fadeOut 0.3s ease 2.7s forwards;
}

.error-icon {
  font-size: 20px;
  color: #ff6464;
}

.error-text {
  font-weight: 600;
  color: #ff6464;
}

@keyframes slideIn {
  from {
    transform: translateX(100%);
    opacity: 0;
  }
  to {
    transform: translateX(0);
    opacity: 1;
  }
}

@keyframes fadeOut {
  from {
    opacity: 1;
  }
  to {
    opacity: 0;
  }
}
</style>
