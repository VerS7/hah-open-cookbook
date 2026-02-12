<style scoped>
.github-link {
  cursor: pointer;
  position: fixed;
  bottom: 1rem;
  left: 1rem;
}

.github-link {
  color: rgb(var(--v-theme-primary));
  transition: all 0.5s ease-in-out;
}

.github-link:hover {
  color: rgb(var(--v-theme-secondary)) !important;
}

.github-tooltip :deep(.v-overlay__content) {
  border-radius: 25px;
  backdrop-filter: blur(5px);
  color: rgb(var(--v-theme-secondary));
  background-color: rgba(255, 255, 255, 0.1);
}
</style>

<template>
  <v-tooltip class="github-tooltip" text="Cookbook repository">
    <template #activator="{ props }">
      <v-icon v-bind="props" class="github-link" size="3rem" color="primary" @click="handleLink">
        mdi-github
      </v-icon>
    </template>
  </v-tooltip>
</template>

<script setup lang="ts">
function handleLink() {
  window.open('https://github.com/VerS7/hah-open-cookbook')
}
</script>
