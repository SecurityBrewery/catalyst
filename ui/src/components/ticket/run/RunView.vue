<script setup lang="ts">
import RunAddDialog from '@/components/ticket/run/RunAddDialog.vue'
import RunContextMenu from '@/components/ticket/run/RunContextMenu.vue'
import StatusIcon from '@/components/ticket/StatusIcon.vue'
import StepView from '@/components/ticket/run/StepView.vue'
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger
} from '@/components/ui/accordion'
import { Card } from '@/components/ui/card'

import { ChevronDown } from 'lucide-vue-next'

import type { Run, Ticket } from '@/lib/types'

const props = defineProps<{
  ticket: Ticket
  runs: Array<Run>
}>()

const status = (run: Run) => {
  if (run.steps.some((step) => step.status === 'failed')) {
    return 'failed'
  }

  if (run.steps.every((step) => step.status === 'completed')) {
    return 'completed'
  }

  if (run.steps.every((step) => step.status === 'open')) {
    return 'open'
  }

  return 'pending'
}
</script>

<template>
  <div class="mt-2 flex flex-col gap-2">
    <Card
      v-if="!runs || runs.length === 0"
      class="flex h-10 items-center p-4 text-muted-foreground"
    >
      No runs added yet.
    </Card>
    <Accordion
      v-else
      type="single"
      collapsible
      class="flex flex-col gap-2"
      :default-value="String(runs[0].id)"
    >
      <AccordionItem
        v-for="run in runs"
        :key="String(run.id)"
        :value="String(run.id)"
        class="border-none"
      >
        <Card>
          <div class="flex flex-row items-center p-2">
            <div class="flex-1">
              <AccordionTrigger class="px-4 py-2">
                <ChevronDown class="h-4 w-4 shrink-0 transition-transform duration-200" />
                {{ run.name }}
                <div class="flex flex-row items-center space-x-2 text-lg font-semibold">
                  <div class="ml-auto"></div>
                </div>
              </AccordionTrigger>
            </div>
            <div class="mr-2 flex flex-row items-center">
              <RunContextMenu :run="run" />
              <StatusIcon :status="status(run)" class="-mr-1 size-7" />
            </div>
          </div>
          <AccordionContent>
            <StepView
              v-for="(step, index) in run.steps"
              :key="step.name"
              :run="run"
              :index="index"
              :step="step"
            />
          </AccordionContent>
        </Card>
      </AccordionItem>
    </Accordion>
    <RunAddDialog :ticket="props.ticket" class="w-full" />
  </div>
</template>
