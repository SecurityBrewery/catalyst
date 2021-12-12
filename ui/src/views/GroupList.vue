<template>
  <v-main style="min-height: 100vh;">
    <List
        :items="groups"
        routername="Group"
        itemid="id"
        itemname="name"
        singular="Group"
        plural="Groups"
    ></List>
  </v-main>
</template>

<script lang="ts">
import Vue from "vue";

import {Group} from "../client";
import {API} from "../services/api";
import List from "../components/List.vue";

interface State {
  groups: Array<Group>;
}

export default Vue.extend({
  name: "GroupList",
  components: {List},
  data: (): State => ({
    groups: [],
  }),
  methods: {
    loadGroups() {
      API.listGroups().then((response) => {
        this.groups = response.data;
      });
    },
    deleteGroup(name: string) {
      // API.deleteU(name).then(() => {
      //   this.loadAutomations();
      // });
    },
  },
  mounted() {
    this.loadGroups();

    this.$store.subscribeAction((action, state) => {
      if (!action.payload || !(this.lodash.has(action.payload, "ids")) || !action.payload["ids"]) {
        return
      }
      let reload = false;
      Vue.lodash.forEach(action.payload["ids"], (id) => {
        if (this.lodash.startsWith(id, "groups/")) {
          reload = true;
        }
      });
      if (reload) {
        this.loadGroups()
      }
    })
  },
});
</script>
