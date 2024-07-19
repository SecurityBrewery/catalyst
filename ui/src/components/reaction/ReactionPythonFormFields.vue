<script setup lang="ts">
import GrowTextarea from '@/components/form/GrowTextarea.vue'
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
  <FormField name="actiondata.requirements" v-slot="{ componentField }">
    <FormItem>
      <FormLabel for="requirements" class="text-right">requirements.txt</FormLabel>
      <FormControl>
        <GrowTextarea id="requirements" class="col-span-3" v-bind="componentField" />
      </FormControl>
      <FormDescription> Specify the Python packages required to run the script. </FormDescription>
      <FormMessage />
    </FormItem>
  </FormField>
  <FormField name="actiondata.script" v-slot="{ componentField }" validate-on-input>
    <FormItem>
      <FormLabel for="script" class="text-right">Script</FormLabel>
      <FormControl>
        <GrowTextarea id="script" class="col-span-3" v-bind="componentField" />
      </FormControl>
      <FormDescription>
        Write a Python script to run when the reaction is triggered.
      </FormDescription>
      <FormMessage />
    </FormItem>
  </FormField>
</template>
