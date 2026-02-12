import { createRouter, createWebHistory } from 'vue-router'

export enum ROUTES {
  LOGIN = 'login',
  MAIN = 'main',
  ADMIN = 'admin',
}

const routes = [
  {
    name: ROUTES.MAIN,
    path: '/',
    component: () => import('@/views/MainView.vue'),
  },
  {
    name: ROUTES.LOGIN,
    path: '/login',
    component: () => import('@/views/LoginView.vue'),
  },
  {
    name: ROUTES.ADMIN,
    path: '/admin',
    component: () => import('@/views/AdminView.vue'),
  },
]

export const router = createRouter({
  history: createWebHistory(),
  routes,
})
