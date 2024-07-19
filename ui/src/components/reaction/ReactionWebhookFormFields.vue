<script setup lang="ts">
import GrowListTextarea from '@/components/form/ListInput.vue'
import {
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
</script>

<template>
  <FormField name="actiondata.headers" v-slot="{ value, handleChange }">
    <FormItem>
      <FormLabel for="headers" class="text-right">Headers</FormLabel>
      <FormControl>
        <GrowListTextarea
          id="headers"
          :modelValue="value"
          @update:modelValue="handleChange"
          placeholder="Content-Type: application/json"
        />
      </FormControl>
      <FormDescription> Specify the headers to include in the request. </FormDescription>
      <FormMessage />
    </FormItem>
  </FormField>
  <FormField name="actiondata.url" v-slot="{ componentField }" validate-on-input>
    <FormItem>
      <FormLabel for="url" class="text-right">URL</FormLabel>
      <FormControl>
        <Input id="url" v-bind="componentField" placeholder="https://example.com/webhook" />
      </FormControl>
      <FormDescription> Specify the URL to send the request to. </FormDescription>
      <FormMessage />
    </FormItem>
  </FormField>
</template>
