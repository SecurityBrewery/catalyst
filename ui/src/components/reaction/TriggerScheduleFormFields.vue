<script setup lang="ts">
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
  <FormField name="triggerdata.expression" v-slot="{ componentField }" validate-on-input>
    <FormItem>
      <FormLabel for="expression" class="text-right"> Cron Expression </FormLabel>
      <FormControl>
        <Input
          id="expression"
          class="col-span-3"
          v-bind="componentField"
          placeholder="Enter a cron expression"
        />
      </FormControl>
      <FormDescription>
        <p class="text-sm text-gray-500">
          A cron expression or macro to schedule the trigger. Example: <code>0 * * * *</code> (every
          hour) or <code>@daily</code> (every day).
          <a
            href="https://en.wikipedia.org/wiki/Cron#CRON_expression"
            target="_blank"
            rel="noopener noreferrer"
            class="text-blue-600 underline"
          >
            Learn more
          </a>
        </p>
      </FormDescription>
      <FormMessage />
    </FormItem>
  </FormField>
</template>
