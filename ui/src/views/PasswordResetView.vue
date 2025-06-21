<script setup lang="ts">
import CatalystLogo from '@/components/common/CatalystLogo.vue'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Button, buttonVariants } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'

import { ref } from 'vue'
import { RouterLink, useRoute } from 'vue-router'

import { cn } from '@/lib/utils'
import router from '@/router'

const route = useRoute()

interface AlertData {
  variant: 'default' | 'destructive'
  title: string
  message: string
}

const mode = ref(route.query.mail && route.query.token ? 'reset-password' : 'enter-email')
const alert = ref<AlertData | null>(null)
const mail = ref<string>((route.query.mail as string) || '')
const token = ref<string>((route.query.token as string) || '')
const password = ref('')
const confirmPassword = ref('')

const sendMail = () => {
  fetch('/auth/local/reset-password-mail', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ email: mail.value })
  }).then((response) => {
    if (!response.ok) {
      alert.value = {
        variant: 'destructive',
        title: 'Password reset failed',
        message: 'Failed to send password reset email'
      }

      return
    }

    alert.value = {
      variant: 'default',
      title: 'Password reset',
      message: 'Password reset email sent'
    }

    mode.value = 'reset-password'
  })
}

const resetPassword = () => {
  if (password.value !== confirmPassword.value) {
    alert.value = {
      variant: 'destructive',
      title: 'Password reset failed',
      message: 'Passwords do not match'
    }

    return
  }

  fetch('/auth/local/reset-password', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      email: mail.value,
      token: token.value,
      password: password.value,
      password_confirm: confirmPassword.value
    })
  }).then((response) => {
    if (!response.ok) {
      alert.value = {
        variant: 'destructive',
        title: 'Password reset failed',
        message: 'Failed to reset password'
      }
      return
    }
    alert.value = {
      variant: 'default',
      title: 'Password reset',
      message: 'Password has been reset successfully'
    }
    mode.value = 'enter-email'
    mail.value = ''
    token.value = ''
    password.value = ''
    confirmPassword.value = ''

    router.push({ name: 'login' })
  })
}
</script>

<template>
  <div class="flex h-full w-full flex-1 items-center justify-center">
    <Card class="m-auto w-96">
      <CardHeader class="flex flex-row justify-between">
        <CardTitle class="flex flex-row">
          <CatalystLogo :size="12" />
          <div>
            <h1 class="text-lg font-bold">Catalyst</h1>
            <div class="text-muted-foreground">Password Reset</div>
          </div>
        </CardTitle>
      </CardHeader>
      <CardContent class="flex flex-col gap-4">
        <Alert v-if="alert" :variant="alert.variant" class="border-4 p-4">
          <AlertTitle>{{ alert.title }}</AlertTitle>
          <AlertDescription>{{ alert.message }}</AlertDescription>
        </Alert>
        <div v-if="mode === 'enter-email'" class="flex flex-col items-stretch gap-4">
          <Input v-model="mail" type="text" placeholder="Email" @keydown.enter="sendMail" />
          <Button variant="outline" class="w-full" :disabled="!mail" @click="sendMail">
            Send Password Reset Email
          </Button>
        </div>
        <div v-else-if="mode === 'reset-password'" class="flex flex-col items-stretch gap-4">
          <Input v-model="mail" type="text" placeholder="Email" />
          <Input v-model="token" type="text" placeholder="Enter received token" />
          <Input v-model="password" type="password" placeholder="New Password" />
          <Input v-model="confirmPassword" type="password" placeholder="Confirm New Password" />
          <Button
            variant="outline"
            class="w-full"
            :disabled="!mail || !token || !password || !confirmPassword"
            @click="resetPassword"
            >Reset Password</Button
          >
        </div>
        <RouterLink
          :to="{ name: 'login' }"
          :class="
            cn(buttonVariants({ variant: 'link', size: 'default' }), 'text-foreground w-full')
          "
          >Back to Login
        </RouterLink>
      </CardContent>
    </Card>
  </div>
</template>
