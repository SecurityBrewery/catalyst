<template>
  <v-main style="min-height: 100vh;">
    <List
        :items="dashboards"
        routername="Dashboard"
        itemid="id"
        itemname="name"
        singular="Dashboard"
        plural="Dashboards"
        writepermission="dashboard:write"
        @delete="deleteDashboard"
    ></List>
  </v-main>
</template>

<script lang="ts">
import Vue from "vue";

import {Dashboard} from "@/client";
import {API} from "@/services/api";
import List from "../components/List.vue";

interface State {
  dashboards: Array<Dashboard>;
}

export default Vue.extend({
  name: "DashboardList",
  components: {List},
  data: (): State => ({
    dashboards: [],
  }),
  methods: {
    loadDashboards() {
      API.listDashboards().then((response) => {
        this.dashboards = response.data;
      });
    },
    deleteDashboard(id: string) {
      API.deleteDashboard(id).then(() => {
        this.loadDashboards();
      })
    }
  },
  mounted() {
    this.loadDashboards();

    this.$store.subscribeAction((action, state) => {
      if (!action.payload || !(this.lodash.has(action.payload, "ids")) || !action.payload["ids"]) {
        return
      }
      let reload = false;
      Vue.lodash.forEach(action.payload["ids"], (id) => {
        if (this.lodash.startsWith(id, "dashboard/")) {
          reload = true;
        }
      });
      if (reload) {
        this.loadDashboards()
      }
    })
  },
});
</script>
