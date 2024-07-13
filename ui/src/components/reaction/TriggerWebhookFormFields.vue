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
  <FormField name="triggerdata.token" v-slot="{ componentField }" validate-on-input>
    <FormItem>
      <FormLabel for="token" class="text-right">Token</FormLabel>
      <FormControl>
        <Input
          id="token"
          class="col-span-3"
          v-bind="componentField"
          placeholder="Enter a token (e.g. 'xyz...')"
        />
      </FormControl>
      <FormDescription>
        Optional. Include an authorization token in the request headers.
      </FormDescription>
      <FormMessage />
    </FormItem>
  </FormField>

  <FormField name="triggerdata.path" v-slot="{ componentField }" validate-on-input>
    <FormItem>
      <FormLabel for="path" class="text-right">Path</FormLabel>
      <FormControl>
        <Input
          id="path"
          class="col-span-3"
          v-bind="componentField"
          placeholder="Enter a path (e.g. 'action1')"
        />
      </FormControl>
      <FormDescription> Specify the path to trigger the reaction. </FormDescription>
      <FormMessage />
    </FormItem>
  </FormField>
</template>
