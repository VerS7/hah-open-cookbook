<style scoped>
.user-icon {
  color: rgb(var(--v-theme-surface-variant));
  outline: 1px solid;
}

.token {
  outline: 1px solid rgb(var(--v-theme-surface-variant));
}
</style>

<template>
  <v-card class="mb-3 mt-1 px-5 py-2" rounded="xl">
    <v-row class="d-flex justify-space-between align-center flex-nowrap" no-gutters>
      <!-- User -->
      <div class="d-flex flex-row justify-center align-center">
        <v-icon class="user-icon mr-3 rounded-circle pa-5" size="2rem">mdi-account</v-icon>
        <div class="d-flex flex-column">
          <div class="text-h5">{{ username ? username : 'NO VALUE' }}</div>
          <span v-if="isAdmin" class="text-caption text-error">ADMIN</span>
          <span v-else class="text-caption text-primary">USER</span>
        </div>
      </div>

      <div class="d-flex flex-row ga-1">
        <copy-token-btn></copy-token-btn>
        <export-btn></export-btn>
      </div>

      <!-- Logout -->
      <v-btn
        color="error"
        variant="outlined"
        rounded="pill"
        @click="ensureLogout"
        :loading="loading"
      >
        logout
      </v-btn>
    </v-row>
  </v-card>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

import ExportBtn from './exportBtn.vue'
import CopyTokenBtn from './CopyTokenBtn.vue'

import { ROUTES } from '@/router'
import { useAuth } from '@/composables/useAuth'

const router = useRouter()

const { username, isAdmin, loading, logout } = useAuth()

function ensureLogout() {
  loading.value = true

  setTimeout(() => {
    logout()
    loading.value = false
    router.push({ name: ROUTES.LOGIN })
  }, 1000)
}
</script>
