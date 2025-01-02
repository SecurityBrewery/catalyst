<script setup lang="ts">
import { Checkbox } from '@/components/ui/checkbox'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '@/components/ui/select'

import isEqual from 'lodash.isequal'
import { onMounted, ref, watch } from 'vue'

import type { JSONSchema } from '@/lib/types'

const model = defineModel<Record<string, any>>()

const props = defineProps<{
  schema: JSONSchema
}>()

const formdata = ref<Record<string, any>>({})

onMounted(() => {
  if (!model.value) return

  for (const key in props.schema.properties) {
    formdata.value[key] = model.value[key]
  }
})

watch(
  () => formdata.value,
  () => {
    const normFormdata = JSON.parse(JSON.stringify(formdata.value))
    const normModel = JSON.parse(JSON.stringify(model.value))

    if (isEqual(normFormdata, normModel)) return

    model.value = { ...formdata.value }
  },
  { deep: true }
)
</script>

<template>
  <div v-for="(property, key) in schema.properties" :key="key">
    <FormField v-if="property.enum" :name="key" v-slot="{ componentField }" v-model="formdata[key]">
      <FormItem>
        <FormLabel :for="key" class="text-right">
          {{ property.title }}
        </FormLabel>
        <Select :id="key" class="col-span-3" v-bind="componentField">
          <SelectTrigger class="font-medium">
            <SelectValue :placeholder="'Select a ' + property.title" />
          </SelectTrigger>
          <SelectContent>
            <SelectGroup>
              <SelectItem v-for="option in property.enum" :key="option" :value="option">
                {{ option }}
              </SelectItem>
            </SelectGroup>
          </SelectContent>
        </Select>
        <FormMessage />
      </FormItem>
    </FormField>
    <FormField
      v-if="property.type === 'string'"
      :name="key"
      v-slot="{ componentField }"
      v-model="formdata[key]"
    >
      <FormItem>
        <FormLabel :for="key" class="text-right">
          {{ property.title }}
        </FormLabel>
        <Input :id="key" class="col-span-3" v-bind="componentField" />
        <FormMessage />
      </FormItem>
    </FormField>
    <FormField
      v-else-if="property.type === 'boolean'"
      :name="key"
      v-slot="{ value, handleChange }"
      type="checkbox"
      v-model="formdata[key]"
    >
      <FormItem class="flex flex-row items-start gap-x-3 space-y-0 py-4">
        <FormControl>
          <Checkbox :checked="value" @update:checked="handleChange" />
        </FormControl>
        <div class="space-y-1 leading-none">
          <FormLabel>
            {{ property.title }}
          </FormLabel>
          <FormMessage />
        </div>
      </FormItem>
    </FormField>
    <FormField
      v-else-if="property.type === 'integer'"
      :name="key"
      v-slot="{ componentField }"
      v-model="formdata[key]"
    >
      <FormItem>
        <FormLabel :for="key" class="text-right">
          {{ property.title }}
        </FormLabel>
        <Input :id="key" class="col-span-3" type="number" v-bind="componentField" />
        <FormMessage />
      </FormItem>
    </FormField>
  </div>
</template>
