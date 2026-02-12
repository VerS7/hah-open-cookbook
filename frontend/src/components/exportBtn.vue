<style scoped>
.export__wrapper {
  display: inline-block;
  position: relative;
}

.export {
  color: rgb(var(--v-theme-primary));
  min-width: v-bind(menuWidth + 'px') !important;
  max-width: v-bind(menuWidth + 'px') !important;
  width: v-bind(menuWidth + 'px') !important;
  justify-content: center !important;
}

/* Стили для элементов меню */
.export__item {
  min-height: 40px !important;
  padding: 0 12px !important;
}

/* Обеспечиваем одинаковую ширину */
:deep(.v-list) {
  min-width: v-bind(menuWidth + 'px') !important;
  width: v-bind(menuWidth + 'px') !important;
}

:deep(.v-list-item) {
  min-height: 40px !important;
  padding: 0 12px !important;
}
</style>

<template>
  <div class="export__wrapper">
    <!-- Кнопка экспорта -->
    <v-btn
      v-if="!showMenu"
      class="export"
      variant="flat"
      density="comfortable"
      rounded="xl"
      prepend-icon="mdi-download"
      @click="toggleMenu"
      :width="menuWidth"
      :loading="loading"
      :disabled="loading"
    >
      Export As
    </v-btn>

    <!-- Меню выбора формата -->
    <v-menu
      v-else
      location="bottom"
      v-model="showMenu"
      :close-on-content-click="false"
      @update:model-value="onMenuStateChange"
    >
      <template #activator="{ props }">
        <v-btn
          v-bind="props"
          variant="flat"
          density="comfortable"
          rounded="xl"
          prepend-icon="mdi-download"
          :width="menuWidth"
          class="export"
        >
          Export As
        </v-btn>
      </template>

      <v-list :width="menuWidth">
        <v-list-item
          v-for="format in exportFormats"
          :key="format.value"
          :value="format.value"
          @click="handleExport(format.value)"
          class="export__item"
        >
          <template #prepend>
            <v-icon :icon="format.icon" size="small"></v-icon>
          </template>
          <v-list-item-title>{{ format.title }}</v-list-item-title>
        </v-list-item>
      </v-list>
    </v-menu>
  </div>
</template>

<script setup lang="ts">
import { useAuth } from '@/composables/useAuth'
import { useExport, type Export } from '@/composables/UseExport'
import { ref, onUnmounted } from 'vue'

interface ExportFormat {
  value: string
  title: string
  icon: string
}

const exportFormats: ExportFormat[] = [
  { value: 'json', title: 'JSON', icon: 'mdi-code-json' },
  { value: 'default', title: 'Custom DB', icon: 'mdi-database-cog' },
]

const { token } = useAuth()
const { loading, exportAs } = useExport(token.value!)

const showMenu = ref(false)
const menuWidth = 180

function onMenuStateChange(value: boolean) {
  if (value) {
    disableScroll()
  } else {
    enableScroll()
  }
}

function disableScroll() {
  const scrollY = window.scrollY

  document.body.style.position = 'fixed'
  document.body.style.top = `-${scrollY}px`
  document.body.style.width = '100%'
  document.body.style.overflow = 'hidden'
}

function enableScroll() {
  const scrollY = document.body.style.top
  document.body.style.position = ''
  document.body.style.top = ''
  document.body.style.width = ''
  document.body.style.overflow = ''

  if (scrollY) {
    window.scrollTo(0, parseInt(scrollY || '0') * -1)
  }
}

function toggleMenu() {
  showMenu.value = !showMenu.value
  if (showMenu.value) {
    disableScroll()
  }
}

async function handleExport(format: string) {
  showMenu.value = false

  enableScroll()

  loading.value = true
  setTimeout(async () => {
    await exportAs(format as Export)
  }, 3000)
}

onUnmounted(() => {
  enableScroll()
})
</script>
