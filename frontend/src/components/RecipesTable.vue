<style scoped>
.gap-1 {
  gap: 4px;
}

.bottom-bar {
  width: fit-content;
  border-radius: 20px;
  left: auto;
  bottom: 10px;
  position: fixed;
  flex-wrap: nowrap;
}

.bar__content {
  flex-grow: 1;
  background-color: rgba(0, 0, 0, 0);
}

.table-container {
  margin-bottom: 100px;
}

.data-table :deep(table) {
  padding-bottom: 100px;
}

.search {
  opacity: 0.8;
  width: 900px;
}

.recipes {
  color: rgb(var(--v-theme-primary));
}

.screenshot {
  margin-left: 150px;
  color: rgb(var(--v-theme-primary));
}

.screenshot-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.1);
  z-index: 1000;
  display: flex;
  align-items: start;
  justify-content: center;
  backdrop-filter: blur(15px);
}

.item-name-cell {
  position: relative;
  max-width: 320px;
  min-width: 0;
  min-height: calc(2 * 1.25em);
}

.item-name-text {
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  overflow: hidden;
  text-overflow: ellipsis;
  line-height: 1.25;
  max-height: calc(2 * 1.25em);
  white-space: normal;
  word-break: normal;
  overflow-wrap: normal;
  hyphens: none;
  cursor: text;
  user-select: text;
}

.item-name-cell.is-expanded {
  z-index: 20;
}

.item-name-cell.is-expanded .item-name-text {
  position: absolute;
  top: 0;
  left: 0;
  z-index: 20;
  display: block;
  max-height: none;
  overflow: visible;
  text-overflow: clip;
  line-clamp: unset;
  -webkit-line-clamp: unset;
  min-width: 100%;
  width: max-content;
  max-width: min(560px, 60vw);
  padding: 10px 14px;
  border-radius: 15px;
  border: 1px solid rgba(255, 255, 255, 0.14);
  background-color: rgba(0, 0, 0, 0.5);
  box-shadow: 0 12px 30px rgba(0, 0, 0, 0.24);
  backdrop-filter: blur(10px);
  animation: item-name-fade-in 0.2s ease;
  cursor: text;
  user-select: text;
}

@keyframes item-name-fade-in {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

.screenshot-mode .item-name-cell {
  max-width: none;
  min-height: auto;
}

.screenshot-mode .item-name-text {
  display: block;
  -webkit-line-clamp: unset;
  line-clamp: unset;
  max-height: none;
  overflow: visible;
  text-overflow: clip;
  white-space: nowrap;
}
</style>

<template>
  <cookbook-versions-bar
    v-model="cookbookVersion"
    @update:model-value="handleCookbookVersionChange"
  >
  </cookbook-versions-bar>

  <v-card rounded="xl" class="pa-5 table-container">
    <div class="d-flex justify-space-between flex-nowrap">
      <v-chip>
        <span class="text-primary mr-1">{{ total ? total : 0 }}</span>
        recipe(s) on
        <span class="text-primary mx-1">{{ pages ? pages : 0 }}</span>
        page(s)
      </v-chip>

      <v-btn
        class="screenshot"
        variant="outlined"
        density="comfortable"
        rounded="xl"
        prepend-icon="mdi-monitor-screenshot"
        :disabled="loading || total === 0"
        :loading="screenshotCapturing"
        @click="handleScreenshot"
      >
        screenshot
      </v-btn>

      <v-chip>
        from
        <span class="text-primary mx-1">{{
          timestamps ? formatDate(timestamps.first) : 'Unknown date'
        }}</span>
        to
        <span class="text-primary ml-1">{{
          timestamps ? formatDate(timestamps.last) : 'Unknown date'
        }}</span>
      </v-chip>
    </div>

    <div v-if="screenshotCapturing" class="screenshot-overlay">
      <v-progress-circular
        class="mt-16"
        size="128"
        :style="stickyLoader"
        indeterminate
      ></v-progress-circular>
    </div>

    <div ref="recipesRef" :class="{ 'screenshot-mode': screenshotCapturing }">
      <v-data-table-server
        class="data-table"
        item-value="id"
        must-sort
        sort-asc-icon="mdi-chevron-down"
        sort-desc-icon="mdi-chevron-up"
        v-model:items-per-page="itemsPerPage"
        :headers="headers"
        :items="serverItems"
        :items-length="totalItems"
        :loading="loading"
        :search="search"
        :sort-by="sortBy"
        @update:sort-by="handleSort"
        hide-default-footer
      >
        <template #[`header.feps`]>
          <fep-control v-model="fep" @update:model-value="handleFepControl"></fep-control>
        </template>

        <template #[`item.resourceName`]="{ item }">
          <v-icon size="2rem">
            <img
              :src="'https://www.havenandhearth.com/mt/r/' + item.resourceName"
              style="width: 100%; height: 100%"
            />
          </v-icon>
        </template>

        <template #[`item.itemName`]="{ item }">
          <div
            class="item-name-cell"
            :class="{ 'is-expanded': expandedNameId === item.id }"
            @mouseenter="handleNameEnter(item.id, $event)"
            @mouseleave="handleNameLeave"
          >
            <span class="item-name-text">{{ item.itemName }}</span>
          </div>
        </template>

        <template #[`item.feps`]="{ item }">
          <fep-stack :feps="item.feps"></fep-stack>
        </template>

        <template #[`item.ingredients`]="{ item }">
          <ingredients-list
            class="align-center"
            :ingredients="item.ingredients"
            :ingredient-width="ingredientWidth"
          ></ingredients-list>
        </template>

        <template #[`item.hunger`]="{ item }">
          <span>{{ item.hunger }}%</span>
        </template>
        <template #[`item.feph`]="{ item }">
          <span>{{ (item.feps.reduce((s, e) => s + e.value, 0) / item.hunger).toFixed(2) }}</span>
        </template>

        <template #[`item.total`]="{ item }">
          <span>{{ item.feps.reduce((s, e) => s + e.value, 0).toFixed(2) }}</span>
        </template>

        <template #[`item.energy`]="{ item }">
          <span>{{ item.energy }}%</span>
        </template>

        <template #[`item.bar`]="{ item }">
          <div class="d-flex align-center">
            <fep-bar :feps="item.feps"></fep-bar>
          </div>
        </template>

        <template #[`item.id`]="{ item }">
          <template v-if="isAdmin && !screenshotCapturing">
            <v-btn icon="mdi-close" size="xs" variant="flat" @click="removeRecipe(item.id)"></v-btn>
          </template>
        </template>

        <template #[`no-data`]>No recipes found</template>
        <template #[`loading`]>Loading recipes...</template>
      </v-data-table-server>
    </div>
  </v-card>

  <!-- Bottom bar (search, pagination) -->
  <div v-if="!screenshotCapturing" class="d-flex flex-row justify-center">
    <v-container class="bottom-bar d-flex morphism px-10">
      <v-card class="bar__content" elevation="0">
        <v-row class="px-3 py-1 flex-nowrap">
          <v-text-field
            class="search mt-5"
            bg-color="surface"
            rounded="lg"
            density="compact"
            variant="outlined"
            v-model="search"
            :loading="loading"
            :error="error !== null"
            @input="handleInput"
          >
            <template #append-inner>
              <help-btn />
            </template>
          </v-text-field>

          <pagination-bar
            v-model="page"
            :count="total"
            :page-size="itemsPerPage"
            @update:pageSize="
              (e: number) => {
                itemsPerPage = e
              }
            "
            @next="load"
            @prev="load"
          ></pagination-bar>
        </v-row>
      </v-card>
    </v-container>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch, type Ref } from 'vue'
import type { DataTableHeader } from 'vuetify'
import { useRoute, useRouter, type LocationQueryRaw } from 'vue-router'

import FepStack from './FepStack.vue'
import FepBar from './FepBar.vue'
import PaginationBar from './PaginationBar.vue'
import IngredientsList from './IngredientsList.vue'
import HelpBtn from './HelpBtn.vue'
import FepControl from './FepControl.vue'

import { useRecipes, type Category, type FoodRecipe, type Order } from '@/composables/useRecipe'
import { useAuth } from '@/composables/useAuth'
import { useDebounce } from '@/composables/useDebounce'
import { useScreenshot } from '@/composables/useScreenshot'
import CookbookVersionsBar from './CookbookVersionsBar.vue'

const route = useRoute()
const router = useRouter()

const cookbookVersion = ref<string | null>(null)
const itemsPerPage = ref(50)
const search = ref('')
const page = ref(1)
const sort = ref<Order>('desc')
const sortBy = ref<Array<{ key: Category; order: Order }>>([])
const category = ref<Category>('total')
const serverItems = ref<FoodRecipe[]>([])
const totalItems = ref(0)
const fep = ref<{ fep: string; order: Order | null }>()
const expandedNameId = ref<number | null>(null)

const headers: DataTableHeader[] = [
  { title: '', key: 'resourceName', sortable: false },
  { title: 'Name', key: 'itemName', sortable: false },
  { title: 'FEPs', key: 'feps', sortable: false, align: 'center', width: 0 },
  { title: 'Total', key: 'total', sortable: true },
  { title: 'Ingredients', key: 'ingredients', sortable: false, width: '35%', align: 'center' },
  { title: 'Hunger', key: 'hunger', sortable: true },
  { title: 'FEP/H', key: 'feph', sortable: true },
  { title: 'Energy', key: 'energy', sortable: true },
  { title: 'FEP bar', key: 'bar', sortable: false, align: 'center' },
  { title: '', key: 'id', sortable: false },
]

const { isAdmin, token } = useAuth()

const initialIngredientWidth = 120
const ingredientWidth = ref(initialIngredientWidth)

const scroll = ref(0)
const stickyLoader = computed(() => ({
  top: `${scroll.value}px`,
}))

const recipesRef = ref<HTMLElement | null>(null)
const screenshotCapturing = ref(false)
const { capture, download } = useScreenshot(recipesRef as Ref<HTMLElement>)
const managedQueryKeys = ['v', 'q', 'l', 'p'] as const
const supportedPageSizes = new Set([10, 20, 50, 100])

const {
  recipes,
  total,
  pages,
  page: loadedPage,
  timestamps,
  loading,
  error,
  remove: removeRecipe,
  get: getRecipes,
} = useRecipes(token.value!)

const {
  value: filterInput,
  debouncedValue: filterDebounce,
  updateDebouncedValue: filterUpdate,
} = useDebounce('', 1000)

const hasRestoredQueryState = managedQueryKeys.some((key) => typeof route.query[key] === 'string')
const pendingInitialQueryLoad = ref(hasRestoredQueryState)

const parsePositiveInt = (value: unknown, fallback: number): number => {
  if (typeof value !== 'string') return fallback

  const parsed = Number.parseInt(value, 10)
  return Number.isInteger(parsed) && parsed > 0 ? parsed : fallback
}

const restoredQuery = typeof route.query.q === 'string' ? route.query.q : ''
const restoredVersion = typeof route.query.v === 'string' ? route.query.v : null
const restoredPageLength = parsePositiveInt(route.query.l, itemsPerPage.value)
const restoredPage = parsePositiveInt(route.query.p, page.value)

if (restoredVersion) {
  cookbookVersion.value = restoredVersion
}

search.value = restoredQuery
page.value = restoredPage
itemsPerPage.value = supportedPageSizes.has(restoredPageLength) ? restoredPageLength : 50

function buildManagedQuery() {
  if (search.value.length === 0 || cookbookVersion.value == null) {
    return {}
  }

  return {
    v: cookbookVersion.value,
    q: search.value,
    l: itemsPerPage.value.toString(),
    p: page.value.toString(),
  }
}

async function syncQueryState() {
  const nextManagedQuery = buildManagedQuery()
  const nextQuery: LocationQueryRaw = { ...route.query }

  for (const key of managedQueryKeys) {
    delete nextQuery[key]
  }

  Object.assign(nextQuery, nextManagedQuery)

  const keys = Object.keys(nextQuery)
  const currentKeys = Object.keys(route.query)
  const isSameQuery =
    keys.length === currentKeys.length &&
    keys.every((key) => {
      const nextValue = nextQuery[key]
      const currentValue = route.query[key]

      if (Array.isArray(nextValue) || Array.isArray(currentValue)) {
        return JSON.stringify(nextValue) === JSON.stringify(currentValue)
      }

      return nextValue === currentValue
    })

  if (isSameQuery) {
    return
  }

  await router.replace({ query: nextQuery })
}

async function load() {
  await loadItems(
    search.value.length > 0 ? search.value : 'total>0',
    sort.value,
    category.value,
    page.value,
    itemsPerPage.value,
  )
}

async function loadDefault() {
  page.value = 1
  sort.value = 'desc'
  category.value = 'total'
  sortBy.value = [{ key: 'total', order: 'desc' }]
  if (fep.value) fep.value.order = null

  await loadItems('total>0', sort.value, category.value, page.value, itemsPerPage.value)
}

async function loadItems(
  filter: string,
  sort: Order,
  category: Category,
  pageNumber: number,
  pageSize: number,
): Promise<void> {
  if (cookbookVersion.value == null) return

  await getRecipes(filter, sort, category, pageNumber, pageSize, cookbookVersion.value)
  pendingInitialQueryLoad.value = false

  if (error.value !== null) {
    return
  }

  if (loadedPage.value !== null) {
    page.value = loadedPage.value
  }

  totalItems.value = total.value!
  serverItems.value = recipes.value!
  await syncQueryState()
}

function handleInput() {
  filterInput.value = search.value
  filterUpdate(search.value)
}

async function handleFepControl() {
  if (!fep.value || !fep.value.order) return

  sortBy.value = []

  sort.value = fep.value.order
  category.value = fep.value.fep as Category

  await load()
}

async function handleScreenshot() {
  switchScreenshotingState()

  await new Promise((resolve) => {
    requestAnimationFrame(() => {
      requestAnimationFrame(resolve)
    })
  })

  await capture()
  download(`cookbook_p${page.value}_i${itemsPerPage.value}`)

  switchScreenshotingState()
}

function handleScroll() {
  scroll.value = window.scrollY
}

function handleNameEnter(id: number, event: MouseEvent) {
  if (screenshotCapturing.value) return

  const cell = event.currentTarget as HTMLElement | null
  if (!cell) return

  const text = cell.querySelector('.item-name-text') as HTMLElement | null
  if (!text) return

  const isTruncated =
    text.scrollHeight > text.clientHeight + 1 || text.scrollWidth > text.clientWidth + 1
  expandedNameId.value = isTruncated ? id : null
}

function handleNameLeave() {
  expandedNameId.value = null
}

function switchScreenshotingState() {
  expandedNameId.value = null

  if (screenshotCapturing.value) {
    ingredientWidth.value = initialIngredientWidth
    recipesRef.value!.style = 'width: initial'
    screenshotCapturing.value = false
  } else {
    ingredientWidth.value = 1
    recipesRef.value!.style = 'width: fit-content'
    screenshotCapturing.value = true
  }
}

async function handleSort(value: Array<{ key: Category; order: Order }>) {
  sortBy.value = value

  if (sortBy.value.length == 0 || sortBy.value[0] === undefined) {
    return
  }

  if (fep.value) fep.value.order = null

  sort.value = sortBy.value[0].order as Order
  category.value = sortBy.value[0].key as Category

  await load()
}

async function handleCookbookVersionChange() {
  if (pendingInitialQueryLoad.value) {
    await load()
    return
  }

  await loadDefault()
}

function formatDate(rawDate: string): string {
  const date = new Date(rawDate)

  const formatter = new Intl.DateTimeFormat('en-GB', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    second: undefined,
    hour12: false,
  })

  return formatter.format(date)
}

watch(filterDebounce, async () => {
  if (search.value == '') {
    await loadDefault()
  } else {
    await load()
  }
})

watch(itemsPerPage, async () => {
  await load()
})

onMounted(async () => {
  window.addEventListener('scroll', handleScroll, true)

  if (pendingInitialQueryLoad.value && cookbookVersion.value != null) {
    await load()
  }
})

onBeforeUnmount(async () => {
  window.removeEventListener('scroll', handleScroll, true)
})
</script>
