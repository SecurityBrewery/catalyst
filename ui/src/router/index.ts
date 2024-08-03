import { createRouter, createWebHistory } from 'vue-router'

import DashboardView from '@/views/DashboardView.vue'
import LoginView from '@/views/LoginView.vue'
import PasswordResetView from '@/views/PasswordResetView.vue'
import ReactionView from '@/views/ReactionView.vue'
import TicketView from '@/views/TicketView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/dashboard'
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: DashboardView
    },
    {
      path: '/tickets/:type/:id?',
      name: 'tickets',
      component: TicketView
    },
    {
      path: '/reactions/:id?',
      name: 'reactions',
      component: ReactionView
    },
    {
      path: '/login',
      name: 'login',
      component: LoginView
    },
    {
      path: '/password-reset',
      name: 'password-reset',
      component: PasswordResetView
    }
  ]
})

export default router
