<script setup lang="ts">
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'

import { useQuery } from '@tanstack/vue-query'
import { ref, watch } from 'vue'

import { pb } from '@/lib/pocketbase'

const mail = ref('')
const password = ref('')
const errorTitle = ref('')
const errorMessage = ref('')

const login = () => {
  pb.collection('users')
    .authWithPassword(mail.value, password.value)
    .then(() => {
      window.location.href = '/ui/'
    })
    .catch((error) => {
      errorTitle.value = 'Login failed'
      errorMessage.value = error.message
    })
}

const resetPassword = () => {
  pb.collection('users')
    .requestPasswordReset(mail.value)
    .then(() => {
      errorTitle.value = 'Password reset'
      errorMessage.value = 'Password reset email sent'
    })
    .catch((error) => {
      errorTitle.value = 'Password reset failed'
      errorMessage.value = error.message
    })
}

const { data: config } = useQuery({
  queryKey: ['config'],
  queryFn: (): Promise<Record<string, Array<String>>> => pb.send('/api/config', {})
})

watch(
  () => config.value,
  () => {
    if (!config.value) return
    if (config.value['flags'].includes('demo') || config.value['flags'].includes('dev')) {
      mail.value = 'user@catalyst-soar.com'
      password.value = '1234567890'
    }
  }
)
</script>

<template>
  <div class="flex h-full w-full flex-1 items-center justify-center">
    <Card class="m-auto w-96">
      <CardHeader class="flex flex-row justify-between">
        <CardTitle>Catalyst</CardTitle>
      </CardHeader>
      <CardContent class="flex flex-col gap-4">
        <Alert v-if="errorTitle || errorMessage" variant="destructive" class="border-4 p-4">
          <AlertTitle v-if="errorTitle">{{ errorTitle }}</AlertTitle>
          <AlertDescription v-if="errorMessage">{{ errorMessage }}</AlertDescription>
        </Alert>
        <Input
          v-model="mail"
          type="text"
          placeholder="Username"
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
        <Button variant="outline" class="w-full" @click="resetPassword">Reset Password</Button>
      </CardContent>
    </Card>
  </div>
</template>
