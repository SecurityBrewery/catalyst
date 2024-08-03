<script setup lang="ts">
import CatalystLogo from '@/components/common/CatalystLogo.vue'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Button, buttonVariants } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'

import { ref } from 'vue'

import { pb } from '@/lib/pocketbase'
import { cn } from '@/lib/utils'

interface AlertData {
  variant: 'default' | 'destructive'
  title: string
  message: string
}

const mail = ref('')
const alert = ref<AlertData | null>(null)

const resetPassword = () => {
  pb.collection('users')
    .requestPasswordReset(mail.value)
    .then(() => {
      alert.value = {
        variant: 'default',
        title: 'Password reset',
        message: 'Password reset email sent'
      }
    })
    .catch((error) => {
      alert.value = {
        variant: 'destructive',
        title: 'Password reset failed',
        message: error.message
      }
    })
}
</script>

<template>
  <div class="flex h-full w-full flex-1 items-center justify-center">
    <Card class="m-auto w-96">
      <CardHeader class="flex flex-row justify-between">
        <CardTitle class="flex flex-row">
          <CatalystLogo class="size-12" />
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
        <div v-else class="flex flex-col gap-4">
          <Input
            v-model="mail"
            type="text"
            placeholder="Email"
            class="w-full"
            @keydown.enter="resetPassword"
          />
          <Button variant="outline" class="w-full" @click="resetPassword">Reset Password</Button>
        </div>
        <RouterLink
          :to="{ name: 'login' }"
          :class="
            cn(buttonVariants({ variant: 'link', size: 'default' }), 'w-full text-foreground')
          "
          >Back to Login
        </RouterLink>
      </CardContent>
    </Card>
  </div>
</template>
