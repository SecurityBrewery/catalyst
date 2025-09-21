<script setup lang="ts">
import CatalystLogo from '@/components/common/CatalystLogo.vue'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Button, buttonVariants } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'

import { useQuery } from '@tanstack/vue-query'
import { ref, watch } from 'vue'
import { useRouter } from 'vue-router'

import { useAPI } from '@/api'
import { cn } from '@/lib/utils'
import { useAuthStore } from '@/store/auth'

const api = useAPI()

const authStore = useAuthStore()
const router = useRouter()

const mail = ref('')
const password = ref('')
const errorTitle = ref('')
const errorMessage = ref('')

const login = () => {
  fetch('/auth/local/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ email: mail.value, password: password.value })
  })
    .then((response) => {
      if (response.ok) {
        response.json().then((data) => {
          authStore.setToken(data.token)
          router.push({ name: 'dashboard' })
        })
      } else {
        errorTitle.value = 'Login failed'
        errorMessage.value = 'Invalid username or password'
      }
    })
    .catch((error) => {
      errorTitle.value = 'Login failed'
      errorMessage.value = error.message
    })
}

const { data: config } = useQuery({
  queryKey: ['config'],
  queryFn: () => api.getConfig()
})

watch(
  () => config.value,
  () => {
    if (!config.value) return
    if (config.value.flags.includes('demo') || config.value.flags.includes('dev')) {
      mail.value = 'user@catalyst-soar.com'
      password.value = '1234567890'
    }
  },
  { immediate: true }
)
</script>

<template>
  <div class="flex h-full w-full flex-1 items-center justify-center">
    <Card class="m-auto w-96">
      <CardHeader class="flex flex-row justify-between">
        <CardTitle class="flex flex-row">
          <CatalystLogo :size="12" />
          <div>
            <h1 class="text-lg font-bold">Catalyst</h1>
            <div class="text-muted-foreground">Login</div>
          </div>
        </CardTitle>
      </CardHeader>
      <CardContent class="flex flex-col gap-4">
        <Alert v-if="errorTitle || errorMessage" variant="destructive" class="border-4 p-4">
          <AlertTitle v-if="errorTitle">{{ errorTitle }}</AlertTitle>
          <AlertDescription v-if="errorMessage">{{ errorMessage }}</AlertDescription>
        </Alert>
        <Input
          v-model="mail"
          type="text"
          placeholder="Email"
          class="w-full"
          @keydown.enter="login"
        />
        <Input
          v-model="password"
          type="password"
          placeholder="Password"
          class="w-full"
          @keydown.enter="login"
        />
        <Button variant="outline" class="w-full" @click="login">Login</Button>
        <RouterLink
          :to="{ name: 'password-reset' }"
          :class="
            cn(buttonVariants({ variant: 'link', size: 'default' }), 'w-full text-foreground')
          "
          >Reset Password
        </RouterLink>
      </CardContent>
    </Card>
  </div>
</template>
