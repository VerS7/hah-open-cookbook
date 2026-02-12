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
</style>

<template>
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

    <div ref="recipesRef">
      <v-data-table-server
        class="data-table"
        item-value="id"
        sort-asc-icon="mdi-chevron-down"
        sort-desc-icon="mdi-chevron-up"
        v-model:items-per-page="itemsPerPage"
        :headers="headers"
        :items="serverItems"
        :items-length="totalItems"
        :loading="loading"
        :search="search"
        @update:sort-by="handleSort"
        hide-default-footer
      >
        <template #[`item.resourceName`]="{ item }">
          <v-icon size="2rem">
            <img
              :src="'https://www.havenandhearth.com/mt/r/' + item.resourceName"
              style="width: 100%; height: 100%"
            />
          </v-icon>
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
                load()
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

import FepStack from './FepStack.vue'
import FepBar from './FepBar.vue'
import PaginationBar from './PaginationBar.vue'
import IngredientsList from './IngredientsList.vue'

import { useRecipes, type Category, type FoodRecipe, type Order } from '@/composables/useRecipe'
import { useAuth } from '@/composables/useAuth'
import { useDebounce } from '@/composables/useDebounce'
import { useScreenshot } from '@/composables/useScreenshot'
import HelpBtn from './HelpBtn.vue'

const itemsPerPage = ref(50)
const search = ref('')
const page = ref(1)
const sort = ref<Order>('desc')
const category = ref<Category>('total')
const serverItems = ref<FoodRecipe[]>([])
const totalItems = ref(0)

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

const {
  recipes,
  total,
  pages,
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
  await loadItems('total>0', 'desc', 'total', 1, 50)
}

async function loadItems(
  filter: string,
  sort: Order,
  category: Category,
  page: number,
  pageSize: number,
): Promise<void> {
  await getRecipes(filter, sort, category, page, pageSize)
  totalItems.value = total.value!
  serverItems.value = recipes.value!
}

function handleInput() {
  filterInput.value = search.value
  filterUpdate(search.value)
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

function switchScreenshotingState() {
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

async function handleSort(s: Array<{ key: string; order: string }>) {
  if (s.length == 0 || s[0] === undefined) {
    return
  }
  const sValue = s[0]

  sort.value = sValue.order as Order
  category.value = sValue.key as Category

  await load()
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
  await loadDefault()
})

onBeforeUnmount(() => {
  window.removeEventListener('scroll', handleScroll, true)
})
</script>
