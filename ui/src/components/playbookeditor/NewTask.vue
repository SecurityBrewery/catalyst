<template>
  <v-card outlined>
    <v-card-title>
      Create a new step
      <v-spacer />
      <v-btn @click="close" outlined small>
        <v-icon>mdi-close</v-icon>
      </v-btn>
    </v-card-title>
    <v-card-text>
      <v-form v-model="valid">
        <v-text-field
            v-model="newTask.name"
            label="Name"
            :rules="[val => (val || '').length > 0 || 'This field is required']"
            variant="underlined"/>
        <v-textarea
            v-model="newTask.description"
            label="Description"
            auto-grow
            rows="1"
            variant="underlined"/>
        <v-text-field
            v-model="newTask.key"
            label="Key (generated automatically)"
            readonly
            disabled
            :rules="[val => (val || '').length > 0 || 'This field is required']"
            variant="underlined"/>
      </v-form>
    </v-card-text>
    <v-card-actions>
      <v-spacer/>
      <v-btn color="primary" @click="createTask" :disabled="!valid">Create</v-btn>
    </v-card-actions>
  </v-card>
</template>

<script setup lang="ts">
import {defineProps, ref, watch, defineEmits} from "vue";

const props = defineProps<{
  playbook: any;
}>();

const valid = ref(false);
const newTask = ref({
  name: '',
  description: '',
  key: '',
  next: {},
});

watch(newTask, (val) => {
  if (val.name) {
    const newKeyBase = val.name.toLowerCase().replace(/ /g, '_');
    let newKey = newKeyBase;
    if (!(newKey in props.playbook.tasks)) {
      newTask.value.key = newKey;
    } else {
      let i = 1;
      while (newKey in props.playbook.tasks) {
        newKey = newKeyBase + '_' + i;
        i++;
      }
      newTask.value.key = newKey;
    }
  }
}, {deep: true});

const emit = defineEmits(["createTask", "close"]);

const createTask = () => {
  emit('createTask', newTask.value);
};

const close = () => {
  emit("close");
};
</script>
