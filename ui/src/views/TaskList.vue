<template>
  <v-main>
    <v-container>
      <v-data-table
          :headers="headers"
          :items="tasks ? tasks : []"
          item-key="name"
          multi-sort
          class="elevation-1 cards clickable"
          :loading="loading"
          :footer-props="{ 'items-per-page-options': [10, 25, 50, 100] }"
          @click:row="open"
      >
      </v-data-table>
    </v-container>
  </v-main>
</template>

<script lang="ts">
import Vue from "vue";
import {TaskResponse, TaskWithContext} from "@/client";
import {API} from "@/services/api";

interface State {
  tasks: Array<TaskResponse>;
  loading: boolean;
}

export default Vue.extend({
  name: "TaskList",
  data: (): State => ({
    tasks: [],
    loading: true,
  }),
  watch: {
    $route: function () {
      this.loadTasks();
    },
  },
  computed: {
    headers() {
      return [
        {
          text: "Ticket",
          align: "start",
          value: "ticket_name"
        },
        {
          text: "Playbook",
          align: "start",
          value: "playbook_name"
        },
        {
          text: "Task",
          align: "start",
          value: "task.name"
        },
        {
          text: "Owner",
          align: "start",
          value: "task.owner"
        }
      ];
    }
  },
  methods: {
    open(task: TaskWithContext) {
      this.$router.push({
        name: "Ticket",
        params: { id: task.ticket_id.toString(), type: "-" }
      });
    },
    loadTasks() {
      this.loading = true;
      API.listTasks().then((reponse) => {
        this.tasks = reponse.data;
        this.loading = false;
      })
    }
  },
  mounted() {
    this.loadTasks();
  }
});
</script>
