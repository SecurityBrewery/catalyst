<script setup lang="ts">
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import {
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '@/components/ui/select'
import { Switch } from '@/components/ui/switch'
import { toast } from '@/components/ui/toast'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip'

import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { defineRule, useForm } from 'vee-validate'
import { ref, watch } from 'vue'

import { useAPI } from '@/api'
import type { Settings } from '@/client/models'
import { handleError } from '@/lib/utils'

const api = useAPI()
const queryClient = useQueryClient()

const submitDisabledReason = ref<string>('')
const isDemo = ref(false)

const { data: config } = useQuery({
  queryKey: ['config'],
  queryFn: () => api.getConfig()
})

watch(
  () => config.value,
  () => {
    if (!config.value) return
    if (config.value.flags.includes('demo')) {
      isDemo.value = true
    }
  },
  { immediate: true }
)

const { data: settings } = useQuery({
  queryKey: ['settings'],
  queryFn: (): Promise<Settings> => api.getSettings()
})

defineRule('required', (value: string) => {
  if (!value || !value.length) {
    return 'This field is required'
  }

  return true
})

const { handleSubmit, validate, values, setValues } = useForm({
  initialValues: {
    meta: {
      appName: '',
      appUrl: '',
      senderName: '',
      senderAddress: ''
    },
    smtp: {
      enabled: false,
      host: '',
      port: 0,
      username: '',
      password: '',
      authMethod: '',
      tls: false,
      localName: ''
    }
  },
  validationSchema: {
    'meta.appName': 'required',
    'meta.appUrl': 'required',
    'meta.senderName': 'required',
    'meta.senderAddress': 'required',
    'smtp.host': (val: string, values: any) =>
      values.smtp && values.smtp.enabled ? (!val ? 'This field is required' : true) : true,
    'smtp.port': (val: number, values: any) =>
      values.smtp && values.smtp.enabled ? (!val ? 'This field is required' : true) : true,
    'smtp.username': (val: string, values: any) =>
      values.smtp && values.smtp.enabled ? (!val ? 'This field is required' : true) : true,
    // 'smtp.password': (val: string, values: any) => true,
    'smtp.authMethod': (val: string, values: any) =>
      values.smtp && values.smtp.enabled ? (!val ? 'This field is required' : true) : true
    // 'smtp.localName': (val: string, values: any) => true
  }
})

watch(
  () => settings.value,
  () => {
    if (!settings.value) return
    setValues({
      meta: {
        appName: settings.value.meta.appName,
        appUrl: settings.value.meta.appUrl,
        senderName: settings.value.meta.senderName,
        senderAddress: settings.value.meta.senderAddress
      },
      smtp: {
        enabled: settings.value.smtp.enabled,
        host: settings.value.smtp.host,
        port: settings.value.smtp.port,
        username: settings.value.smtp.username,
        password: settings.value.smtp.password,
        authMethod: settings.value.smtp.authMethod,
        tls: settings.value.smtp.tls,
        localName: settings.value.smtp.localName
      }
    })

    updateSubmitDisabledReason()
  },
  { immediate: true }
)

const updateSettingsMutation = useMutation({
  mutationFn: (vals: Settings): Promise<Settings> => api.updateSettings({ settings: vals }),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['settings'] })
    toast({
      title: 'Settings updated',
      description: 'The settings have been updated successfully'
    })
  },
  onError: handleError('Failed to update settings')
})

function equalSettings(values: any, value?: Settings) {
  if (!value) return false

  return (
    values.meta.appName === value.meta.appName &&
    values.meta.appUrl === value.meta.appUrl &&
    values.meta.senderName === value.meta.senderName &&
    values.meta.senderAddress === value.meta.senderAddress &&
    values.smtp.enabled === value.smtp.enabled &&
    values.smtp.host === value.smtp.host &&
    values.smtp.port === value.smtp.port &&
    values.smtp.username === value.smtp.username &&
    values.smtp.password === value.smtp.password &&
    values.smtp.authMethod === value.smtp.authMethod &&
    values.smtp.tls === value.smtp.tls &&
    values.smtp.localName === value.smtp.localName
  )
}

const updateSubmitDisabledReason = () => {
  if (isDemo.value) {
    submitDisabledReason.value = 'Settings cannot be edited in demo mode'

    return
  }

  if (equalSettings(values, settings.value)) {
    submitDisabledReason.value = 'Make changes to save'

    return
  }

  validate({ mode: 'silent' }).then((res) => {
    if (res.valid) {
      submitDisabledReason.value = ''
    } else {
      submitDisabledReason.value = 'Please fix the errors'
    }
  })
}

watch(
  () => values,
  () => updateSubmitDisabledReason(),
  { deep: true, immediate: true }
)

watch(
  () => isDemo.value,
  () => updateSubmitDisabledReason()
)

const onSubmit = handleSubmit((vals) => {
  if (!settings.value) return
  updateSettingsMutation.mutate({
    ...settings.value,
    meta: {
      ...settings.value.meta,
      appName: vals.meta.appName,
      appUrl: vals.meta.appUrl
    },
    smtp: {
      ...settings.value.smtp,
      enabled: vals.smtp.enabled,
      host: vals.smtp.host,
      port: vals.smtp.port,
      username: vals.smtp.username,
      password: vals.smtp.password,
      authMethod: vals.smtp.authMethod,
      tls: vals.smtp.tls,
      localName: vals.smtp.localName
    }
  })
})
</script>

<template>
  <form @submit="onSubmit" class="flex w-full flex-col items-start gap-4">
    <FormField name="meta.appName" v-slot="{ componentField }" validate-on-input>
      <FormItem class="w-full">
        <FormLabel for="meta.appName" class="text-right">App Name</FormLabel>
        <Input id="meta.appName" class="col-span-3" v-bind="componentField" />
        <FormMessage />
      </FormItem>
    </FormField>

    <FormField name="meta.appUrl" v-slot="{ componentField }" validate-on-input>
      <FormItem class="w-full">
        <FormLabel for="meta.appUrl" class="text-right">App URL</FormLabel>
        <Input id="meta.appUrl" class="col-span-3" v-bind="componentField" />
        <FormMessage />
      </FormItem>
    </FormField>

    <FormField name="meta.senderName" v-slot="{ componentField }" validate-on-input>
      <FormItem class="w-full">
        <FormLabel for="meta.senderName" class="text-right">Sender Name</FormLabel>
        <Input id="meta.senderName" class="col-span-3" v-bind="componentField" />
        <FormMessage />
      </FormItem>
    </FormField>

    <FormField name="meta.senderAddress" v-slot="{ componentField }" validate-on-input>
      <FormItem class="w-full">
        <FormLabel for="meta.senderAddress" class="text-right">Sender Address</FormLabel>
        <Input id="meta.senderAddress" type="email" class="col-span-3" v-bind="componentField" />
        <FormMessage />
      </FormItem>
    </FormField>

    <!-- Card class="w-full">
      <CardHeader>
        <CardTitle>Logs</CardTitle>
      </CardHeader>
      <CardContent class="flex flex-col gap-4">
        <FormField name="maxDays" v-slot="{ componentField }" validate-on-input>
          <FormItem class="w-full">
            <FormLabel for="maxDays" class="text-right">Max Days</FormLabel>
            <Input id="maxDays" type="number" class="col-span-3" v-bind="componentField" />
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField name="logLevel" v-slot="{ componentField }" validate-on-input>
          <FormItem>
            <FormLabel for="logLevel" class="text-right">Log Level</FormLabel>
            <FormControl>
              <Select id="logLevel" class="col-span-3" v-bind="componentField">
                <SelectTrigger class="font-medium">
                  <SelectValue placeholder="Select log level" />
                </SelectTrigger>
                <SelectContent>
                  <SelectGroup>
                    <SelectItem value="debug">Debug</SelectItem>
                    <SelectItem value="info">Info</SelectItem>
                    <SelectItem value="warn">Warn</SelectItem>
                    <SelectItem value="error">Error</SelectItem>
                    <SelectItem value="fatal">Fatal</SelectItem>
                    <SelectItem value="off">Off</SelectItem>
                  </SelectGroup>
                </SelectContent>
              </Select>
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField name="logIP" v-slot="{ value, handleChange }">
          <FormItem class="w-full items-center gap-2">
            <FormLabel>Log IP</FormLabel>
            <div class="flex flex-row items-center gap-2">
              <FormControl>
                <Switch :checked="value" @update:checked="handleChange" />
              </FormControl>
              <FormDescription>
                Check to log the IP address of the user in the logs.
              </FormDescription>
            </div>
          </FormItem>
          <FormMessage />
        </FormField>
      </CardContent>
    </Card -->

    <Card class="w-full">
      <CardHeader>
        <CardTitle>SMTP</CardTitle>
      </CardHeader>
      <CardContent class="flex flex-col gap-4">
        <FormField name="smtp.enabled" v-slot="{ value, handleChange }">
          <FormItem class="w-full items-center gap-2">
            <FormLabel>Enabled</FormLabel>
            <div class="flex flex-row items-center gap-2">
              <FormControl>
                <Switch :checked="value" @update:checked="handleChange" />
              </FormControl>
              <FormDescription>
                Check to enable SMTP settings. If enabled, you must provide the host, port,
                username, and password.
              </FormDescription>
            </div>
          </FormItem>
        </FormField>

        <FormField name="smtp.host" v-slot="{ componentField }" validate-on-input>
          <FormItem class="w-full">
            <FormLabel for="smtp.host" class="text-right">Host</FormLabel>
            <Input id="smtp.host" class="col-span-3" v-bind="componentField" />
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField name="smtp.port" v-slot="{ componentField }" validate-on-input>
          <FormItem class="w-full">
            <FormLabel for="smtp.port" class="text-right">Port</FormLabel>
            <Input id="smtp.port" type="number" class="col-span-3" v-bind="componentField" />
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField name="smtp.username" v-slot="{ componentField }" validate-on-input>
          <FormItem class="w-full">
            <FormLabel for="smtp.username" class="text-right">Username</FormLabel>
            <Input id="smtp.username" class="col-span-3" v-bind="componentField" />
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField name="smtp.password" v-slot="{ componentField }" validate-on-input>
          <FormItem class="w-full">
            <FormLabel for="smtp.password" class="text-right">Password</FormLabel>
            <Input id="smtp.password" type="password" class="col-span-3" v-bind="componentField" />
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField name="smtp.authMethod" v-slot="{ componentField }" validate-on-input>
          <FormItem>
            <FormLabel for="smtp.authMethod" class="text-right"> Authentication Method </FormLabel>
            <Select id="smtp.authMethod" class="col-span-3" v-bind="componentField">
              <SelectTrigger class="font-medium">
                <SelectValue placeholder="Select authentication method" />
              </SelectTrigger>
              <SelectContent>
                <SelectGroup>
                  <SelectItem value="plain">Plain (default)</SelectItem>
                  <SelectItem value="login">Login</SelectItem>
                </SelectGroup>
              </SelectContent>
            </Select>
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField name="smtp.tls" v-slot="{ value, handleChange }">
          <FormItem class="w-full items-center gap-2">
            <FormLabel>TLS</FormLabel>
            <div class="flex flex-row items-center gap-2">
              <FormControl>
                <Switch :checked="value" @update:checked="handleChange" />
              </FormControl>
              <FormDescription> Check to enable TLS for SMTP connections.</FormDescription>
            </div>
          </FormItem>
        </FormField>

        <FormField name="smtp.localName" v-slot="{ componentField }" validate-on-input>
          <FormItem class="w-full">
            <FormLabel for="smtp.localName" class="text-right">HELO domain</FormLabel>
            <Input id="smtp.localName" class="col-span-3" v-bind="componentField" />
            <FormMessage />
          </FormItem>
        </FormField>
      </CardContent>
    </Card>

    <!-- Card class="w-full">
      <CardHeader>
        <CardTitle>Backup</CardTitle>
      </CardHeader>
      <CardContent class="flex flex-col gap-4">
        <FormField name="cron" v-slot="{ componentField }" validate-on-input>
          <FormItem class="w-full">
            <FormLabel for="cron" class="text-right">Cron Expression</FormLabel>
            <Input id="cron" class="col-span-3" v-bind="componentField" />
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField name="cronMaxKeep" v-slot="{ componentField }" validate-on-input>
          <FormItem class="w-full">
            <FormLabel for="cronMaxKeep" class="text-right">Max Keep</FormLabel>
            <Input id="cronMaxKeep" type="number" class="col-span-3" v-bind="componentField" />
            <FormMessage />
          </FormItem>
        </FormField>

      </CardContent>
    </Card -->

    <Alert v-if="isDemo" variant="destructive">
      <AlertTitle>Cannot save</AlertTitle>
      <AlertDescription>{{ submitDisabledReason }}</AlertDescription>
    </Alert>
    <div class="flex gap-4">
      <TooltipProvider :delay-duration="0">
        <Tooltip>
          <TooltipTrigger class="cursor-default">
            <Button
              type="submit"
              :variant="submitDisabledReason !== '' ? 'secondary' : 'default'"
              :disabled="submitDisabledReason !== ''"
              :title="submitDisabledReason"
            >
              Save
            </Button>
          </TooltipTrigger>
          <TooltipContent>
            <span v-if="submitDisabledReason !== ''">
              {{ submitDisabledReason }}
            </span>
            <span v-else> Save settings. </span>
          </TooltipContent>
        </Tooltip>
      </TooltipProvider>
    </div>
  </form>
</template>
