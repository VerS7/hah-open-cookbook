<style scoped>
.frame {
  width: 100vw;
  height: 100vh;
}

.background {
  width: 90vw;
  height: 90vh;
  background-size: cover;
  background-position: center;
  background-image: url('hah.png');
}

.login {
  background-color: rgba(0, 0, 0, 0.3);
  border: 1px solid rgba(var(--v-theme-surface), 0.5);
  backdrop-filter: blur(5px);
}
</style>

<template>
  <div class="frame d-flex justify-center align-center">
    <div class="background d-flex justify-center align-center rounded-xl">
      <v-card class="login d-flex flex-column pa-5" width="400px" rounded="xl" elevation="16">
        <span class="text-h2 my-5 text-center">HaH Cookbook</span>
        <v-text-field
          placeholder="username"
          variant="outlined"
          rounded="lg"
          density="compact"
          v-model="username"
        ></v-text-field>
        <v-text-field
          placeholder="password"
          type="password"
          variant="outlined"
          rounded="lg"
          density="compact"
          v-model="password"
        ></v-text-field>
        <v-btn
          color="primary"
          variant="outlined"
          rounded="xl"
          @click="handleLogin"
          :loading="loading"
          :error="error !== null"
        >
          login
        </v-btn>
      </v-card>
    </div>
  </div>

  <github-link></github-link>
</template>

<script setup lang="ts">
import { onBeforeMount, ref } from 'vue'

import { useAuth } from '@/composables/useAuth'
import { useRouter } from 'vue-router'
import { ROUTES } from '@/router'
import GithubLink from '@/components/GithubLink.vue'

const { loading, error, isAuthenticated, login, loadFromStorage } = useAuth()

const username = ref('')
const password = ref('')

const router = useRouter()

async function handleLogin() {
  loading.value = true

  setTimeout(async () => {
    await login(username.value, password.value)
    if (error.value !== null) return
    router.push({ name: 'main' })
  }, 1000)
}

onBeforeMount(() => {
  loadFromStorage()

  if (isAuthenticated.value) router.push({ name: ROUTES.MAIN })
})
</script>
