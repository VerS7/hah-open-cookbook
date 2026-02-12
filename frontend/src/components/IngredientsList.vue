<style scoped>
.ingredients-container {
  position: relative;
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  max-width: 100%;
  min-height: 32px;
}

.more-chip {
  cursor: pointer;
  flex-shrink: 0;
}

.hidden-ingredients-popup {
  position: absolute;
  top: 100%;
  left: -8px;
  z-index: 1000;
  margin-top: -3px;
  padding: 8px;
  border-radius: 15px;
  overflow-y: auto;
  background-color: rgba(0, 0, 0, 0.5);
}

.hidden-ingredients-container {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.hidden-chip-item {
  width: fit-content;
}

.ingredient {
  font-weight: bold;
  margin-left: 3px;
  color: rgb(96, 96, 255);
}

.gradient {
  transition: all 0.4s ease !important;
  position: relative;
  overflow: hidden;
  z-index: 1;
}

.gradient::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(45deg, #60daff, #ff0059);
  z-index: -1;
  opacity: 0;
  transition: opacity 0.4s ease;
}

.gradient:hover::before {
  opacity: 1;
}

.popup-fade-enter-active {
  animation: fadeIn 0.2s ease;
}

.popup-fade-leave-active {
  animation: fadeOut 0.2s ease;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(-10px) scale(0.95);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

@keyframes fadeOut {
  from {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
  to {
    opacity: 0;
    transform: translateY(-10px) scale(0.95);
  }
}

.hidden-ingredients-popup {
  transition: all 0.2s ease;
}

.hidden-ingredients-popup.leaving {
  opacity: 0;
  transform: translateY(-10px) scale(0.95);
}
</style>

<template>
  <div class="ingredients-container flex-nowrap" ref="containerRef">
    <template v-for="ingredient in localIngredients" :key="ingredient.name">
      <v-chip
        v-if="!ingredient.isHidden"
        size="small"
        class="chip-item"
        @mouseenter="handleIngredientHover(ingredient)"
        @mouseleave="handleIngredientLeave"
      >
        <span>{{ ingredient.name }}:</span>
        <span class="ingredient">{{ ingredient.percentage }}%</span>
      </v-chip>
    </template>

    <v-chip
      v-if="hiddenCount > 0"
      class="more-chip gradient"
      size="small"
      @mouseenter="handleMoreIngredientsEnter"
      @mouseleave="handleMoreIngredientsLeave"
    >
      +{{ hiddenCount }}
    </v-chip>

    <transition name="popup-fade">
      <v-card
        v-if="hiddenChipsVisible"
        class="hidden-ingredients-popup morphism"
        @mouseenter="handlePopupEnter"
        @mouseleave="handlePopupLeave"
      >
        <div class="hidden-ingredients-container pa-2">
          <v-chip
            v-for="ingredient in hiddenIngredients"
            :key="'hidden-' + ingredient.name"
            size="small"
            class="hidden-chip-item"
          >
            <span>{{ ingredient.name }}:</span>
            <span class="ingredient">{{ ingredient.percentage }}%</span>
          </v-chip>
        </div>
      </v-card>
    </transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'

import { type FoodIngredient } from '@/composables/useRecipe'

interface LocalIngredient extends FoodIngredient {
  isHidden: boolean
}

const props = withDefaults(
  defineProps<{
    ingredients: FoodIngredient[]
    maxWidth?: number
    ingredientWidth?: number
    hideDelay?: number
  }>(),
  {
    ingredients: () => [],
    maxWidth: 450,
    ingredientWidth: 120,
    hideDelay: 300,
  },
)

const containerRef = ref<HTMLElement>()
const containerWidth = ref(0)
const hiddenChipsVisible = ref(false)
const isHiding = ref(false)
const hideTimeout = ref<number | null>(null)
const resizeObserver = ref<ResizeObserver>()

const localIngredients = ref<LocalIngredient[]>([])

const hiddenIngredients = computed(() =>
  localIngredients.value.filter((ingredient) => ingredient.isHidden),
)

const hiddenCount = computed(() => hiddenIngredients.value.length)

function transformIngredients(ingredients: FoodIngredient[]): LocalIngredient[] {
  return ingredients.map((ingredient) => ({
    ...ingredient,
    isHidden: false,
  }))
}

function updateLocalIngredients() {
  const transformed = transformIngredients(props.ingredients)

  if (localIngredients.value.length > 0) {
    const maxVisible = calculateMaxVisibleIngredients()

    localIngredients.value = transformed.map((ingredient, index) => ({
      ...ingredient,
      isHidden: index >= maxVisible,
    }))
  } else {
    const maxVisible = calculateMaxVisibleIngredients()
    localIngredients.value = transformed.map((ingredient, index) => ({
      ...ingredient,
      isHidden: index >= maxVisible,
    }))
  }
}

function clearHideTimeout() {
  if (hideTimeout.value) {
    clearTimeout(hideTimeout.value)
    hideTimeout.value = null
  }
}

function showPopup() {
  clearHideTimeout()
  isHiding.value = false
  hiddenChipsVisible.value = true
}

function hidePopup() {
  clearHideTimeout()
  hideTimeout.value = setTimeout(() => {
    isHiding.value = true
    setTimeout(() => {
      hiddenChipsVisible.value = false
      isHiding.value = false
    }, 200)
  }, props.hideDelay)
}

function handleMoreIngredientsEnter() {
  clearHideTimeout()
  showPopup()
}

function handleMoreIngredientsLeave() {
  hidePopup()
}

function handlePopupEnter() {
  clearHideTimeout()
}

function handlePopupLeave() {
  hidePopup()
}

function handleIngredientHover(ingredient: LocalIngredient) {
  if (ingredient.isHidden) {
    clearHideTimeout()
    showPopup()
  }
}

function handleIngredientLeave() {
  hidePopup()
}

function calculateMaxVisibleIngredients(): number {
  if (!containerWidth.value) return localIngredients.value.length

  const availableWidth = Math.min(containerWidth.value, props.maxWidth)
  const chipTotalWidth = props.ingredientWidth + 6

  const moreButtonWidth = 50
  const maxWithoutMore = Math.floor((availableWidth - moreButtonWidth) / chipTotalWidth)
  const maxWithMore = Math.floor(availableWidth / chipTotalWidth)

  if (localIngredients.value.length <= maxWithMore) {
    return localIngredients.value.length
  }

  return Math.max(0, maxWithoutMore)
}

function updateVisibilityFlags() {
  const maxVisible = calculateMaxVisibleIngredients()

  localIngredients.value = localIngredients.value.map((ingredient, index) => ({
    ...ingredient,
    isHidden: index >= maxVisible,
  }))
}

function observeResize() {
  if (!containerRef.value) return

  resizeObserver.value = new ResizeObserver((entries) => {
    for (const entry of entries) {
      containerWidth.value = entry.contentRect.width
      nextTick(() => {
        updateVisibilityFlags()
      })
    }
  })

  resizeObserver.value.observe(containerRef.value)
}

onMounted(() => {
  nextTick(() => {
    if (containerRef.value) {
      containerWidth.value = containerRef.value.clientWidth
      observeResize()
      updateLocalIngredients()
    }
  })
})

onUnmounted(() => {
  clearHideTimeout()
  if (resizeObserver.value) {
    resizeObserver.value.disconnect()
  }
})

watch(
  () => props.ingredients,
  () => {
    nextTick(() => {
      updateLocalIngredients()
      if (containerRef.value) {
        containerWidth.value = containerRef.value.clientWidth
      }
    })
  },
  { deep: true },
)

watch(
  () => props.ingredientWidth,
  () => {
    nextTick(() => {
      updateVisibilityFlags()
    })
  },
  { immediate: true },
)

watch(containerWidth, () => {
  nextTick(() => {
    updateVisibilityFlags()
  })
})
</script>
