<style scoped>
.main {
  max-width: 90vw;
}
</style>

<template>
  <v-container class="main">
    <user-bar></user-bar>
    <recipes-table></recipes-table>
  </v-container>

  <github-link></github-link>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'

import RecipesTable from '@/components/RecipesTable.vue'
import UserBar from '@/components/UserBar.vue'

import { useAuth } from '@/composables/useAuth'
import GithubLink from '@/components/GithubLink.vue'

const router = useRouter()

const { isAuthenticated, loadFromStorage } = useAuth()

onMounted(() => {
  loadFromStorage()

  if (!isAuthenticated.value) {
    router.push({ name: 'login' })
  }
})
</script>
