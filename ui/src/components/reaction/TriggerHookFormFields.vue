<script setup lang="ts">
import TriggerHookFormFieldCollections from '@/components/reaction/TriggerHookFormFieldCollections.vue'
import TriggerHookFormFieldEvents from '@/components/reaction/TriggerHookFormFieldEvents.vue'
import {
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage
} from '@/components/ui/form'
</script>

<template>
  <FormField name="triggerdata.collections" v-slot="{ componentField }" validate-on-input>
    <FormItem>
      <FormLabel for="collections" class="text-right">Collections</FormLabel>
      <FormControl>
        <TriggerHookFormFieldCollections id="collections" v-bind="componentField" />
      </FormControl>
      <FormDescription> Specify the collections to trigger the reaction. </FormDescription>
      <FormMessage />
    </FormItem>
  </FormField>

  <FormField name="triggerdata.events" v-slot="{ componentField }" validate-on-input>
    <FormItem>
      <FormLabel for="events" class="text-right">Events</FormLabel>
      <FormControl>
        <TriggerHookFormFieldEvents id="events" v-bind="componentField" />
      </FormControl>
      <FormDescription> Specify the events to trigger the reaction. </FormDescription>
      <FormMessage />
    </FormItem>
  </FormField>
</template>
