import { createRouter, createWebHistory } from 'vue-router'

import { useAuthStore } from '@/store/auth'
import LoginView from '@/views/LoginView.vue'
import PasswordResetView from '@/views/PasswordResetView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/dashboard',
      meta: { requiresAuth: true }
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: () => import('@/views/DashboardView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/tickets/:type/:id?',
      name: 'tickets',
      component: () => import('@/views/TicketView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/reactions/:id?',
      name: 'reactions',
      component: () => import('@/views/ReactionView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/users/:id?',
      name: 'users',
      component: () => import('@/views/UserView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/groups/:id?',
      name: 'groups',
      component: () => import('@/views/GroupView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/types/:id?',
      name: 'types',
      component: () => import('@/views/TypeView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/login',
      name: 'login',
      component: LoginView,
      meta: { requiresAuth: false }
    },
    {
      path: '/password-reset',
      name: 'password-reset',
      component: PasswordResetView,
      meta: { requiresAuth: false }
    }
  ]
})

router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next('/login')
  } else {
    next()
  }
})

export default router
