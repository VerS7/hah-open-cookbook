<style scoped>
.alert {
  position: fixed;
  right: 0;
  top: 1rem;
  transition: all 0.5s ease-in-out;
  opacity: 0.1;
  transform: translateX(100%);
}

.alert.show {
  opacity: 1;
  transform: translateX(0);
}
</style>

<template>
  <teleport to="body">
    <div v-show="isVisible" class="alert" :class="{ show: isActive }">
      <slot></slot>
    </div>
  </teleport>
</template>

<script setup lang="ts">
import { ref, watch, nextTick, onUnmounted } from 'vue'

const activator = defineModel<boolean>({ default: false })

const props = withDefaults(
  defineProps<{
    timeout: number
  }>(),
  {
    timeout: 1000,
  },
)

const isVisible = ref(false)
const isActive = ref(false)
let hideTimeout: ReturnType<typeof setTimeout> | null = null

function showAlert() {
  if (hideTimeout) {
    clearTimeout(hideTimeout)
    hideTimeout = null
  }

  isVisible.value = true

  nextTick(() => {
    isActive.value = true

    hideTimeout = setTimeout(() => {
      hideAlert()
    }, props.timeout)
  })
}

function hideAlert() {
  isActive.value = false

  setTimeout(() => {
    isVisible.value = false
    hideTimeout = null
  }, 500)
}

watch(activator, (newValue) => {
  if (newValue) {
    showAlert()
    activator.value = false
  }
})

onUnmounted(() => {
  if (hideTimeout) {
    clearTimeout(hideTimeout)
  }
})
</script>
