<style scoped>
.token__copy {
  color: rgb(var(--v-theme-primary));
}

.alert {
  background-color: rgba(var(--v-theme-success), 0.8) !important;
  border-top-left-radius: 25px !important;
  border-bottom-left-radius: 25px !important;
}
</style>

<template>
  <v-btn
    class="token__copy"
    width="200"
    rounded="xl"
    variant="flat"
    density="comfortable"
    prepend-icon="mdi-content-copy"
    @click="tokenToClipboard"
  >
    copy token
  </v-btn>

  <alert-bar v-model="alertBarActivator" :timeout="1500">
    <v-card class="alert morphism text-button px-5" rounded="0" elevation="5">Token copied!</v-card>
  </alert-bar>
</template>

<script setup lang="ts">
import { ref } from 'vue'

import AlertBar from './AlertBar.vue'

import { useAuth } from '@/composables/useAuth'

const { token } = useAuth()

const alertBarActivator = ref(false)

async function tokenToClipboard() {
  await navigator.clipboard.writeText(token.value!)
  alertBarActivator.value = true
}
</script>
