<style scoped>
.selector {
  cursor: pointer;
}

.selected {
  color: yellow !important;
}

.alarm {
  position: absolute;
  left: 20px;
}

.accent {
  border: 1px solid rgb(var(--v-theme-secondary));
  background-color: rgba(var(--v-theme-secondary), 0.3);
  border-radius: 7px;
}
</style>

<template>
  <v-card class="mb-3 pa-3 d-flex flex-row justify-center ga-3" rounded="xl">
    <div class="alarm d-flex flex-row justify-center">
      <v-chip>
        <v-icon class="mr-1" color="yellow">mdi-alert-circle-outline</v-icon>
        <span class="accent px-2 mr-1">archived_</span> cannot be modified!
      </v-chip>
    </div>

    <div class="ml-16" v-if="loading">loading cookbook versions...</div>
    <div class="ml-16" v-else-if="error">can't load cookbook versions</div>
    <template v-else v-for="(version, i) in versions" v-bind:key="i">
      <v-chip
        class="selector"
        color="primary"
        :class="version.version === model ? 'selected' : ''"
        @click="() => handleSwitchVersion(version.version)"
      >
        {{ version.version }}
      </v-chip>
    </template>
  </v-card>
</template>

<script setup lang="ts">
import { useCookbookVersions } from '@/composables/useCookbookVersion'
import { onBeforeMount } from 'vue'

const { loading, error, versions, get } = useCookbookVersions()

const model = defineModel<string | null>({ default: null })

function handleSwitchVersion(version: string) {
  model.value = version
}

onBeforeMount(async () => {
  await get()

  const defaultVersion = versions.value?.filter(
    (v) => !v.isArchived && v.version.includes('w161'),
  )[0]?.version

  if (defaultVersion !== undefined) {
    model.value = defaultVersion
  }
})
</script>
