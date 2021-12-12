<template>
  <v-main style="min-height: 100vh;">
    <List
        :items="rules"
        routername="Rule"
        itemid="id"
        itemname="name"
        singular="Rule"
        plural="Rules"
        @delete="deleteRule"
    ></List>
  </v-main>
</template>

<script lang="ts">
import Vue from "vue";

import {Rule} from "@/client";
import {API} from "@/services/api";
import List from "../components/List.vue";

interface State {
  rules: Array<Rule>;
}

export default Vue.extend({
  name: "RuleList",
  components: {List},
  data: (): State => ({
    rules: [],
  }),
  methods: {
    loadRules() {
      API.listRules().then((response) => {
        this.rules = response.data;
      });
    },
    deleteRule(name: string) {
      API.deleteRule(name).then(() => {
        this.loadRules();
      });
    },
  },
  mounted() {
    this.loadRules();

    this.$store.subscribeAction((action, state) => {
      if (!action.payload || !(this.lodash.has(action.payload, "ids")) || !action.payload["ids"]) {
        return
      }
      let reload = false;
      Vue.lodash.forEach(action.payload["ids"], (id) => {
        if (this.lodash.startsWith(id, "rules/")) {
          reload = true;
        }
      });
      if (reload) {
        this.loadRules()
      }
    })
  },
});
</script>
