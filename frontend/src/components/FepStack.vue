<style scoped>
.feps {
  width: fit-content;
  border-radius: 8px;
  overflow: hidden;
  flex-wrap: nowrap;
}

.fep {
  width: 2.5rem;
  height: 2.5rem;
  text-align: center;
  font-size: 0.75rem;
  color: black;
  justify-content: center;
}

.high {
  font-style: italic;
}
</style>

<template>
  <v-row class="feps" no-gutters>
    <template v-for="name in fepNames" :key="name">
      <div class="fep d-flex flex-column" :class="name">
        <span v-if="getFEPValue(name, 1) !== null">
          {{ getFEPValue(name, 1) }}
        </span>
        <span v-else>&nbsp;</span>

        <span v-if="getFEPValue(name, 2) !== null" class="high">
          {{ getFEPValue(name, 2) }}
        </span>
        <span v-else>&nbsp;</span>
      </div>
    </template>
  </v-row>
</template>

<script setup lang="ts">
import type { FoodFEP } from '@/composables/useRecipe'

const props = defineProps<{ feps: FoodFEP[] }>()

const fepNames = ['str', 'agi', 'int', 'con', 'prc', 'csm', 'dex', 'wil', 'psy'] as const

function getFEPValue(name: string, multiplier: number): string | null {
  const fep = props.feps.find((f) => f.name === name + multiplier)
  return fep ? fep.value.toFixed(2) : null
}
</script>
