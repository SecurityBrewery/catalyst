<template>
  <v-row style="position: relative; min-height: 600px">
    <v-col :cols="(selectedStep && selectedStep in playbook.tasks) || showNewDialog ? 8 :12">
      <v-card style="overflow: hidden" outlined>
        <v-card-text>
          <PanZoom ref="panZoomPanel" :config="panZoomConfig">
            <template #default>
              <PlaybookGraph
                  v-if="playbook"
                  :playbook="playbook"
                  :horizontal="horizontal"
                  :selected="selectedStep"
                  @update:selected="showNewDialog = false; selectedStep = $event"
              />
            </template>
            <template #actionbar>
              <v-btn @click="showNewDialog = true; selectedStep = ''" large rounded>
                <v-icon color="#000">mdi-plus</v-icon>
                New Step
              </v-btn>
            </template>
            <template #toolbar>
              <v-btn @click="toggleOrientation" label="Toggle Orientation" icon>
                <v-icon color="#000">mdi-format-rotate-90</v-icon>
              </v-btn>
            </template>
          </PanZoom>
        </v-card-text>
      </v-card>
    </v-col>
    <v-col
        v-if="(selectedStep && selectedStep in playbook.tasks) || showNewDialog"
        cols="4"
    >
        <EditTask
            v-if="selectedStep && selectedStep in playbook.tasks"
            :value="playbook.tasks[selectedStep]"
            @input="updateTask"
            :possibleNexts="possibleNexts"
            :parents="parents"
            :playbook="playbook"
            @delete="deleteTask"
            @close="unselectStep"/>
        <NewTask
            v-else-if="showNewDialog"
            :playbook="playbook"
            @createTask="createTask"
            @close="showNewDialog = false"/>
    </v-col>
  </v-row>
</template>

<script lang="ts" setup>
import {computed, nextTick, ref, defineProps, set, defineEmits, del} from "vue";
import PlaybookGraph from "@/components/playbookeditor/PlaybookGraph.vue";
import PanZoom from "@/components/playbookeditor/PanZoom.vue";
import EditTask from "@/components/playbookeditor/EditTask.vue";
import NewTask from "@/components/playbookeditor/NewTask.vue";

const props = defineProps({
  'value': {
    type: Object,
    required: true
  }
});

const emit = defineEmits(['input']);

const playbook = computed({
  get: () => props.value,
  set: (value) => {
    emit('input', value);
  }
});

// selected step
const selectedStep = ref("")
const unselectStep = () => {
  selectedStep.value = "";
}

const updateTask = (task: any) => {
  set(playbook.value.tasks, selectedStep.value, task);
  emit('input', playbook.value);
}

const deleteTask = () => {
  const parents = Array<any>();
  for (const task in playbook.value.tasks) {
    if (playbook.value.tasks[task].next && playbook.value.tasks[task].next[selectedStep.value]) {
      parents.push(task);
    }
  }

  const children = Array<any>();
  if (playbook.value.tasks[selectedStep.value].next) {
    for (const next in playbook.value.tasks[selectedStep.value].next) {
      children.push(next);
    }
  }

  for (const parent of parents) {
    del(playbook.value.tasks[parent].next, selectedStep.value);
    for (const child of children) {
      set(playbook.value.tasks[parent].next, child, playbook.value.tasks[selectedStep.value].next[child]);
    }
  }

  del(playbook.value.tasks, selectedStep.value);

  // for (const task in playbook.value.tasks) {
  //   if (playbook.value.tasks[task].next && playbook.value.tasks[task].next[selectedStep.value]) {
  //     del(playbook.value.tasks[task].next, selectedStep.value);
  //   }
  // }
  // del(playbook.value.tasks, selectedStep.value);
  emit('input', playbook.value);
}

const panZoomConfig = ref({})

const panZoomPanel = ref(null);

const horizontal = ref(false);
const toggleOrientation = () => {
  horizontal.value = !horizontal.value;

  nextTick(() => {
    panZoomPanel.value?.reset(false);
  });
}

const showNewDialog = ref(false);
const createTask = (task: any) => {
  const t = {
    name: task.name,
    description: task.description,
    type: 'task',
    next: {}
  };
  set(playbook.value.tasks, task.key, t);
  selectedStep.value = task.key;
};

// edit task
const possibleNexts = computed(() => {
  if (!selectedStep.value) {
    return [];
  }

  let nexts = Object.keys(playbook.value.tasks);
  nexts = nexts.filter((n) => n !== selectedStep.value);

  // remove any nexts that are already in the list
  if (playbook.value.tasks[selectedStep.value] && 'next' in playbook.value.tasks[selectedStep.value]) {
    for (const next in playbook.value.tasks[selectedStep.value].next) {
      nexts = nexts.filter((n) => n !== next);
    }
  }

  // remove parents recursively
  const parents = findAncestor(selectedStep.value);
  for (const parent of parents) {
    nexts = nexts.filter((n) => n !== parent);
  }

  const result: Array<Record<string, string>> = [];
  for (const next of nexts) {
    if (next && playbook.value.tasks[next].name) {
      result.push({"key": next, "name": playbook.value.tasks[next].name});
    }
  }

  return result;
});

const findAncestor = (step: string): Array<string> => {
  const parents: Array<string> = [];
  for (const task in playbook.value.tasks) {
    for (const next in playbook.value.tasks[task].next) {
      if (next === step) {
        if (!parents.includes(task)) {
          parents.push(task);
          parents.push(...findAncestor(task));
        }
      }
    }
  }
  return parents;
};

const parents = computed(() => {
  const parents: Array<string> = [];
  for (const task in playbook.value.tasks) {
    for (const next in playbook.value.tasks[task].next) {
      if (next === selectedStep.value) {
        if (!parents.includes(task)) {
          parents.push(task);
        }
      }
    }
  }
  return parents;
});
</script>

<style>
#graphwrapper {
  padding: 20px;
}

#gab {
  position: absolute;
  left: 10px;
  top: 20px;
  border-radius: 30px;
  z-index: 200;
}

#gtb {
  position: absolute;
  left: 0;
  bottom: 20px;
  border-radius: 30px;
  z-index: 200;
}

#gtb .v-toolbar__content > .v-btn:first-child {
  margin-inline-start: 0;
}

#gtb .v-toolbar__content > .v-btn:last-child {
  margin-inline-end: 0;
}
</style>
