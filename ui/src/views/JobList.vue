<template>
  <v-main style="min-height: 100vh;">
    <v-row class="fill-height ma-0">
      <v-col cols="3" class="listnav" style="">
        <v-list nav color="background">
          <v-list-item
              v-if="canWrite"
              :to="{ name: 'Job', params: { id: 'new' } }"
              class="mt-4 mx-4 text-center newbutton"
          >
            <v-list-item-content>
              <v-list-item-title>
                <v-icon small>mdi-plus</v-icon> New Job
              </v-list-item-title>
            </v-list-item-content>
          </v-list-item>
          <v-subheader class="pl-4">Jobs</v-subheader>
          <v-list-item
              v-for="item in (jobs ? jobs : [])"
              :key="item.id"
              link
              :to="{ name: 'Job', params: { id: item.id } }"
              class="mx-2"
          >
            <v-list-item-icon>
              <v-icon v-if="item.status === 'created'">mdi-star</v-icon>
              <v-icon v-else-if="item.status === 'running'">mdi-run-fast</v-icon>
              <v-icon v-else-if="item.status === 'paused'">mdi-pause</v-icon>
              <v-icon v-else-if="item.status === 'restarting'">mdi-restart</v-icon>
              <v-icon v-else-if="item.status === 'removing'">mdi-close</v-icon>
              <v-icon v-else-if="item.status === 'exited'">mdi-exit-to-app</v-icon>
              <v-icon v-else-if="item.status === 'dead'">mdi-skull</v-icon>
              <v-icon v-else-if="item.status === 'completed'">mdi-check</v-icon>
            </v-list-item-icon>
            <v-list-item-content>
              <v-list-item-subtitle>
                {{ item.id }}
              </v-list-item-subtitle>
              <v-list-item-title>
                {{ item.automation }} ({{ item.status }})
              </v-list-item-title>
            </v-list-item-content>
          </v-list-item>
        </v-list>
      </v-col>
      <v-col cols="9">
        <router-view></router-view>
      </v-col>
    </v-row>
  </v-main>
</template>

<script lang="ts">
import Vue from "vue";

import {JobResponse} from "@/client";
import {API} from "@/services/api";

interface State {
  jobs: Array<JobResponse>;
}

export default Vue.extend({
  name: "JobList",
  components: {},
  data: (): State => ({
    jobs: [],
  }),
  computed: {
    canWrite: function (): boolean {
      return this.hasRole("admin:job:write");
    },
  },
  methods: {
    loadJobs() {
      API.listJobs().then((response) => {
        if (response.data) {
          this.jobs = response.data;
        }
      });
    },
    hasRole: function (s: string): boolean {
      if (this.$store.state.user.roles) {
        return this.lodash.includes(this.$store.state.user.roles, s);
      }
      return false;
    }
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
