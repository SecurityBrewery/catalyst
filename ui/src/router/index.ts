import { createRouter, createWebHistory } from 'vue-router'

import { useAuthStore } from '@/store/auth'
import DashboardView from '@/views/DashboardView.vue'
import LoginView from '@/views/LoginView.vue'
import PasswordResetView from '@/views/PasswordResetView.vue'
import ReactionView from '@/views/ReactionView.vue'
import RoleView from '@/views/RoleView.vue'
import TicketView from '@/views/TicketView.vue'
import TypeView from '@/views/TypeView.vue'
import UserView from '@/views/UserView.vue'

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
      component: DashboardView,
      meta: { requiresAuth: true }
    },
    {
      path: '/tickets/:type/:id?',
      name: 'tickets',
      component: TicketView,
      meta: { requiresAuth: true }
    },
    {
      path: '/reactions/:id?',
      name: 'reactions',
      component: ReactionView,
      meta: { requiresAuth: true }
    },
    {
      path: '/users/:id?',
      name: 'users',
      component: UserView,
      meta: { requiresAuth: true }
    },
    {
      path: '/groups/:id?',
      name: 'groups',
      component: RoleView,
      meta: { requiresAuth: true }
    },
    {
      path: '/types/:id?',
      name: 'types',
      component: TypeView,
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
