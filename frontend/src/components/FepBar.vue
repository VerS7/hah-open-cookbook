<style scoped>
.bar {
  width: 8rem;
  height: 1rem;
  border-radius: 10px;
  overflow: hidden;
}

.empty {
  background-color: rgba(105, 105, 105, 0.3);
}
</style>

<template>
  <v-row class="bar">
    <template v-if="total !== 0">
      <template v-for="fep in feps" v-bind:key="fep.name">
        <div :class="fep.name.slice(0, -1)" :style="`width: ${(fep.value / total) * 100}%`"></div>
      </template>
    </template>

    <template v-else>
      <div class="empty" style="width: 100%"></div>
    </template>
  </v-row>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { FoodFEP } from '@/composables/useRecipe'

const props = defineProps<{ feps: FoodFEP[] }>()
const total = computed(() => props.feps.reduce((s, v) => s + v.value, 0))
</script>
