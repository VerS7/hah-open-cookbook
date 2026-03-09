// Recipes

import { computed, ref } from 'vue'
import { BASE_API_URL, useAPI } from './useAPI'
import { buildAuthTokenHeader } from './useAuth'

export type FEP =
  | 'str'
  | 'str1'
  | 'str2'
  | 'agi'
  | 'agi1'
  | 'agi2'
  | 'int'
  | 'int1'
  | 'int2'
  | 'con'
  | 'con1'
  | 'con2'
  | 'prc'
  | 'prc1'
  | 'prc2'
  | 'csm'
  | 'csm1'
  | 'csm2'
  | 'dex'
  | 'dex1'
  | 'dex2'
  | 'wil'
  | 'wil1'
  | 'wil2'
  | 'psy'
  | 'psy1'
  | 'psy2'

export type Order = 'asc' | 'desc'
export type Category = 'name' | 'hunger' | 'total' | 'energy' | 'feph' | FEP

export interface FoodFEP {
  name: FEP
  value: number
}

export interface FoodIngredient {
  name: string
  percentage: string
}

export interface FoodRecipe {
  id: number
  itemName: string
  resourceName: string
  hunger: number
  energy: number
  feps: FoodFEP[]
  ingredients: FoodIngredient[]
}

export interface FoodRecipeTimestamps {
  first: string
  last: string
}

export interface FoodRecipesPage {
  recipes: FoodRecipe[]
  count: number
  total: number
  page: number
  pages: number
  timestamps: FoodRecipeTimestamps
}

export interface FoodRecipesTimestamps {
  first: string
  last: string
}

export function useRecipes(token: string) {
  const loading = ref(false)
  const error = ref<string | null>(null)
  const recipesPage = ref<FoodRecipesPage | null>(null)

  const recipes = computed<FoodRecipe[] | null>(() =>
    recipesPage.value ? recipesPage.value.recipes : null,
  )
  const total = computed<number | null>(() => (recipesPage.value ? recipesPage.value.total : null))
  const count = computed<number | null>(() => (recipesPage.value ? recipesPage.value.count : null))
  const page = computed<number | null>(() => (recipesPage.value ? recipesPage.value.page : null))
  const pages = computed<number | null>(() => (recipesPage.value ? recipesPage.value.pages : null))
  const timestamps = computed<FoodRecipesTimestamps | null>(() =>
    recipesPage.value ? recipesPage.value.timestamps : null,
  )

  async function get(
    filter: string,
    sort: Order,
    by: Category,
    page: number,
    length: number,
    cookbookVersion: string,
  ) {
    const query = new URLSearchParams({
      filter: filter,
      sort: sort,
      by: by,
      p: page.toString(),
      l: length.toString(),
    })

    loading.value = true
    error.value = null

    const {
      data,
      error: APIerror,
      execute,
    } = useAPI<FoodRecipesPage>(BASE_API_URL + `/${cookbookVersion}/recipes?${query}`, {
      method: 'GET',
      headers: buildAuthTokenHeader(token),
    })

    await execute()
    if (APIerror.value) error.value = APIerror.value
    if (data.value) recipesPage.value = data.value
    loading.value = false
  }

  async function remove(id: number) {
    loading.value = true
    error.value = null

    const { error: APIerror, execute } = useAPI<FoodRecipesPage>(BASE_API_URL + `/recipe/` + id, {
      method: 'DELETE',
      headers: buildAuthTokenHeader(token),
    })

    await execute()

    if (APIerror.value) error.value = APIerror.value
    loading.value = false
  }

  return { recipes, total, count, page, pages, timestamps, loading, error, get, remove }
}
