<style scoped>
.fep-control {
  width: 2.5rem;
  height: 50px;
  cursor: pointer;
  transition: all 0.2s ease;
  user-select: none;
  position: relative;
  overflow: hidden;
}

.icon-wrapper {
  width: 2.5rem;
  height: 15px;
  position: relative;
  overflow: hidden;
}

.sort-icon {
  position: absolute;
}

.sort-icon.visible {
  opacity: 1;
}

.fep-control.has-sort:hover .sort-icon.visible {
  transform: translateY(-100%);
  opacity: 0;
}

.fep-control.has-sort:hover .sort-icon.hover-icon {
  transform: translateY(0);
  opacity: 1;
}

.fep-control:not(.has-sort):hover .sort-icon.hover-icon {
  transform: translateY(0);
  opacity: 0.7;
}

.sort-icon.hover-icon {
  transform: translateY(100%);
  opacity: 0;
}
</style>

<template>
  <div class="d-flex flex-row flex-nowrap fep-header">
    <template v-for="fep in FEPS" :key="fep">
      <div class="fep-control" :class="{ active: model?.fep === fep }" @click="toggleSort(fep)">
        <div class="icon-wrapper d-flex justify-center">
          <v-icon
            size="small"
            class="sort-icon"
            :class="{
              visible: model?.fep === fep && model.order,
              'hover-icon': !(model?.fep === fep && model.order),
            }"
          >
            {{ model?.order === 'asc' ? 'mdi-chevron-up' : 'mdi-chevron-down' }}
          </v-icon>
        </div>
        <span class="control-label">{{ fep.toUpperCase() }}</span>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import type { Order } from '@/composables/useRecipe'

const model = defineModel<{ fep: string; order: Order | null }>()

const FEPS = ['str', 'agi', 'int', 'con', 'prc', 'csm', 'dex', 'wil', 'psy'] as const

function toggleSort(fep: string) {
  if (model.value?.fep === fep) {
    model.value = {
      fep,
      order: model.value.order === 'desc' ? 'asc' : 'desc',
    }
  } else {
    model.value = {
      fep,
      order: 'desc',
    }
  }
}
</script>
