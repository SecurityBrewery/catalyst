import { createRouter, createWebHistory } from 'vue-router'

import DashboardView from '@/views/DashboardView.vue'
import LoginView from '@/views/LoginView.vue'
import PlaybookView from '@/views/PlaybookView.vue'
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
      path: '/login',
      name: 'login',
      component: LoginView
    },
    {
      path: '/tickets/:type/:id?',
      name: 'tickets',
      component: TicketView
    },
    {
      path: '/playbooks/:id?',
      name: 'playbooks',
      component: PlaybookView
    }
  ]
})

export default router
