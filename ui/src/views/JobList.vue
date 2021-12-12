<template>
  <v-main style="min-height: 100vh;">
    <List
        :items="jobs"
        routername="Job"
        itemid="id"
        itemname="id"
        singular="Job"
        plural="Jobs"
        writepermission="admin:job:write"
        :deletable="false"
    ></List>
  </v-main>
</template>

<script lang="ts">
import Vue from "vue";

import {JobResponse} from "@/client";
import {API} from "@/services/api";
import List from "../components/List.vue";

interface State {
  jobs: Array<JobResponse>;
}

export default Vue.extend({
  name: "JobList",
  components: {List},
  data: (): State => ({
    jobs: [],
  }),
  methods: {
    loadJobs() {
      API.listJobs().then((response) => {
        if (response.data) {
          this.jobs = response.data;
        }
      });
    },
  },
  mounted() {
    this.loadJobs();

    this.$store.subscribeAction((action, state) => {
      if (!action.payload || !(this.lodash.has(action.payload, "ids")) || !action.payload["ids"]) {
        return
      }
      let reload = false;
      Vue.lodash.forEach(action.payload["ids"], (id) => {
        if (this.lodash.startsWith(id, "jobs/")) {
          reload = true;
        }
      });
      if (reload) {
        this.loadJobs()
      }
    })
  },
});
</script>
