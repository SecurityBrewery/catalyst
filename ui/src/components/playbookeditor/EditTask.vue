<template>
  <v-card data-app outlined>
    <v-card-title class="d-flex">
      <span>Edit Task</span>
      <v-spacer/>
      <v-dialog v-model="deleteDialog" max-width="400">
        <template #activator="{ on }">
          <v-btn outlined v-on="on" class="mr-2" small>
            <v-icon>mdi-delete</v-icon>
          </v-btn>
        </template>
        <v-card>
          <v-card-title class="headline">Delete Task</v-card-title>
          <v-card-text>Are you sure you want to delete this task?</v-card-text>
          <v-card-actions>
            <v-spacer/>
            <v-btn color="blue darken-1" text @click="deleteDialog = false">Cancel</v-btn>
            <v-btn color="blue darken-1" text @click="deleteTask">Delete</v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
      <v-btn @click="close" outlined small>
        <v-icon>mdi-close</v-icon>
      </v-btn>
    </v-card-title>
    <v-card-text>
      <v-text-field
          v-model="task.name"
          label="Name"
          variant="underlined"/>
      <v-textarea
          v-model="task.description"
          label="Description"
          auto-grow
          rows="1"
          variant="underlined"/>
      <v-select
          v-model="task.type"
          :items="['input','automation','task']"
          label="Type"
          variant="underlined"/>

      <AdvancedJSONSchemaEditor
          v-if="task.type === 'input'"
          :schema="task.schema"
          @save="task.schema = JSON.parse($event)" />

      <v-select
          v-if="task.type === 'automation'"
          v-model="task.automation"
          :items="automations"
          item-text="id"
          item-value="id"
          label="Automation"
          variant="underlined"/>

      <v-list v-if="task.type === 'automation'">
        <v-subheader class="pa-0" style="padding-inline-start: 0 !important;">Payload Mapping</v-subheader>
        <v-toolbar v-for="(expr, key) in task.payload" :key="key" class="next-row" flat dense>
          {{ key }}:
          <v-text-field
              v-model="task.payload[key]"
              label="Expression"
              variant="solo"
              clearable
              hide-details
              density="compact"
              bg-color="surface"
          />
          <v-btn @click="deletePayloadMapping(key)" color="error" class="pa-0 ma-0" icon>
            <v-icon>mdi-delete</v-icon>
          </v-btn>
        </v-toolbar>
        <v-toolbar class="next-row" flat dense>
          <v-text-field
              v-model="newPayloadMapping"
              label="Payload Field"
              variant="solo"
              bg-color="surface"
              hide-details
              density="compact"
          />:
          <v-text-field
              v-model="newExpression"
              label="CAQL Expression"
              variant="solo"
              hide-details
              density="compact"
              bg-color="surface"
          />
          <v-btn
              @click="addPayloadMapping"
              :disabled="!newPayloadMapping || !newExpression"
              class="pa-0 ma-0"
              icon>
            <v-icon>mdi-plus</v-icon>
          </v-btn>
        </v-toolbar>
      </v-list>

      <v-list v-if="task.next || possibleNexts.length > 1">
        <v-subheader class="pa-0" style="padding-inline-start: 0 !important;">Next Task(s)</v-subheader>
        <v-toolbar v-for="(expr, key) in task.next" :key="key" class="next-row" flat dense>
          If
          <v-text-field
              v-model="task.next[key]"
              label="Condition (leave empty to always run)"
              variant="solo"
              clearable
              hide-details
              density="compact"
              bg-color="surface"
          />
          run
          <span class="font-weight-black">
            {{ playbook.tasks[key].name }}
          </span>
          <v-btn @click="deleteNext(key)" color="error" class="pa-0 ma-0" icon>
            <v-icon>mdi-delete</v-icon>
          </v-btn>
        </v-toolbar>
        <v-toolbar v-if="possibleNexts.length > 0" class="next-row" flat dense>
          If
          <v-text-field
              v-model="newCondition"
              label="Condition (leave empty to always run)"
              variant="solo"
              clearable
              hide-details
              density="compact"
              bg-color="surface"
          />
          run
          <v-select
              v-model="newNext"
              item-text="name"
              item-value="key"
              :items="possibleNexts"
              variant="solo"
              bg-color="surface"
              hide-details
              density="compact"
          />
          <v-btn
              @click="addNext"
              :disabled="!newNext"
              class="pa-0 ma-0"
              icon>
            <v-icon>mdi-plus</v-icon>
          </v-btn>
        </v-toolbar>
      </v-list>
      <v-switch
          v-if="parents.length > 1"
          label="Join (Require all previous tasks to be completed)"
          v-model="task.join"
          color="primary"/>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import {defineProps, ref, watch, defineEmits, del, set, onMounted, computed} from "vue";
import AdvancedJSONSchemaEditor from "@/components/AdvancedJSONSchemaEditor.vue";
import {API} from "@/services/api";
import {AutomationResponse} from "@/client";

interface Task {
  name: string;
  description: string;
  type: string;
  next: Record<string, string>;
  payload: Record<string, string>;
  join: boolean;
}

const props = defineProps<{
  value: Task;
  possibleNexts: Array<Record<string, string>>;
  parents: Array<string>;
  playbook: object;
}>();

const emit = defineEmits(["input", "delete", "close"]);

const deleteDialog = ref(false);
const deleteTask = () => emit("delete");

const close = () => {
  emit("close");
};

const task = ref(props.value);
watch(() => props.value, (value) => {
  task.value = value;
});
watch(task, (value) => {
  emit("input", value);
});

// const task = computed({
//   get: () => {
//     console.log("get", props.value);
//     return props.value;
//   },
//   set: (value) => {
//     console.log("set", value);
//     emit("input", value);
//   }
// });

const deleteNext = (key: string) => {
  del(task.value.next, key);
};

const deletePayloadMapping = (key: string) => {
  del(task.value.payload, key);
};

const newNext = ref('');
const newCondition = ref('');

const newPayloadMapping = ref('');
const newExpression = ref('');

watch(() => props.possibleNexts, () => {
  if (props.possibleNexts.length > 0) {
    newNext.value = props.possibleNexts[0].key;
  }
}, {deep: true, immediate: true});


const addNext = () => {
  if (task.value.next === undefined) {
    // task.value.next = {};
    set(task.value, 'next', {});
  }
  // task.value.next[newNext.value] = newCondition.value;
  set(task.value.next, newNext.value, newCondition.value);
  newNext.value = "";
  newCondition.value = "";
};

const addPayloadMapping = () => {
  if (task.value.payload === undefined) {
    // task.value.payload = {};
    set(task.value, 'payload', {});
  }
  // task.value.payload[newPayloadMapping.value] = newExpression.value;
  set(task.value.payload, newPayloadMapping.value, newExpression.value);
  newPayloadMapping.value = "";
  newExpression.value = "";
};

const automations = ref<Array<AutomationResponse>>([]);
onMounted(() => {
  API.listAutomations().then((response) => {
    automations.value = response.data;
  });
});
</script>

<style>
.next-row {
  padding-top: 10px;
  background: none !important;
}

.next-row .v-toolbar__content {
  gap: 5px;
  padding: 0;
}
</style>
