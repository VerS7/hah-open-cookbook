<style scoped>
.selected {
  color: yellow !important;
}
</style>

<template>
  <v-card class="mb-3 pa-3 d-flex flex-row justify-center ga-3" rounded="xl">
    <div class="ml-16" v-if="loading">loading cookbook versions...</div>
    <div class="ml-16" v-else-if="error">can't load cookbook versions</div>
    <template v-else v-for="(version, i) in versions" v-bind:key="i">
      <v-chip
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

  const defaultVersion = versions.value?.filter((v) => !v.isArchived)[0]?.version

  if (defaultVersion !== undefined) {
    model.value = defaultVersion
  }
})
</script>
