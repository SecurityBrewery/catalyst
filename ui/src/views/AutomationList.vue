<template>
  <v-main style="min-height: 100vh">
    <List
        :items="automations"
        routername="Automation"
        itemid="id"
        itemname="id"
        singular="Automation"
        plural="Automations"
        writepermission="automation:write"
        @delete="deleteAutomation"
    ></List>
  </v-main>
</template>

<script lang="ts">
import Vue from "vue";

import {AutomationResponse} from "@/client";
import {API} from "@/services/api";
import List from "../components/List.vue";

interface State {
  automations: Array<AutomationResponse>;
}

export default Vue.extend({
  name: "AutomationList",
  components: {List},
  data: (): State => ({
    automations: [],
  }),
  methods: {
    loadAutomations() {
      API.listAutomations().then((response) => {
        this.automations = response.data;
      });
    },
    deleteAutomation(name: string) {
      API.deleteAutomation(name).then(() => {
        this.loadAutomations();
      });
    },
  },
  mounted() {
    this.loadAutomations();

    this.$store.subscribeAction((action, state) => {
      if (!action.payload || !(this.lodash.has(action.payload, "ids")) || !action.payload["ids"]) {
        return
      }
      let reload = false;
      Vue.lodash.forEach(action.payload["ids"], (id) => {
        if (this.lodash.startsWith(id, "automations/")) {
          reload = true;
        }
      });
      if (reload) {
        this.loadAutomations()
      }
    })
  },
});
</script>
