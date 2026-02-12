<style scoped>
.page {
  min-width: 50px;
  min-height: 50px;
}

.arrow {
  background-color: rgba(0, 0, 0, 0);
}
</style>

<template>
  <v-row class="d-flex flex-row align-center ml-5 flex-nowrap" no-gutters>
    <v-btn
      class="arrow morphism"
      icon="mdi-chevron-left"
      elevation="5"
      border="md"
      rounded="circle"
      @click="handlePrev"
    />

    <div class="page mx-2 d-flex flex-row justify-center align-center rounded-circle elevation-5">
      {{ page }}
    </div>

    <v-btn
      class="arrow morphism"
      icon="mdi-chevron-right"
      elevation="5"
      border="md"
      rounded="circle"
      @click="handleNext"
    />

    <v-select
      class="ml-4"
      rounded="lg"
      variant="outlined"
      density="compact"
      suffix="&nbsp;/&nbsp;page"
      :items="pageSizes"
      :model-value="pageSize"
      @update:model-value="updatePageSize"
      hide-details
    />
  </v-row>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    count: number | null
    pageSize: number
    pageSizes?: number[]
  }>(),
  {
    pageSizes: () => [10, 20, 50, 100],
  },
)

const page = defineModel<number>({ default: 1 })

const emit = defineEmits<{
  next: []
  prev: []
  'update:pageSize': [value: number]
}>()

const maxPages = computed(() => {
  if (props.count === null) return 0
  return Math.ceil(props.count / props.pageSize)
})

const canGoPrev = computed(() => page.value > 1)
const canGoNext = computed(() => page.value < maxPages.value)

function handlePrev() {
  if (!canGoPrev.value) return

  if (page.value > maxPages.value) {
    page.value = maxPages.value
  } else {
    page.value -= 1
  }

  emit('prev')
}

function handleNext() {
  if (!canGoNext.value) return

  if (page.value + 1 >= maxPages.value) {
    page.value = maxPages.value
  } else {
    page.value += 1
  }

  emit('next')
}

function updatePageSize(value: number) {
  emit('update:pageSize', value)
}
</script>
