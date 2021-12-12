<template>
  <v-main style="min-height: 100vh;">
    <List
        :items="playbooks"
        routername="Playbook"
        itemid="id"
        itemname="name"
        singular="Playbook"
        plural="Playbooks"
        @delete="deletePlaybook"
        writepermission="engineer:playbook:write"
    ></List>
  </v-main>
</template>

<script lang="ts">
import Vue from "vue";

import {PlaybookTemplate} from "../client";
import {API} from "../services/api";
import List from "../components/List.vue";

interface State {
  playbooks: Array<PlaybookTemplate>;
}

export default Vue.extend({
  name: "PlaybookList",
  components: {List},
  data: (): State => ({
    playbooks: [],
  }),
  methods: {
    loadPlaybooks() {
      API.listPlaybooks().then((response) => {
        if (response.data) {
          this.playbooks = response.data;
        }
      });
    },
    deletePlaybook(name: string) {
      API.deletePlaybook(name).then(() => {
        this.loadPlaybooks();
      });
    },
  },
  mounted() {
    this.loadPlaybooks();

    this.$store.subscribeAction((action, state) => {
      if (!action.payload || !(this.lodash.has(action.payload, "ids")) || !action.payload["ids"]) {
        return
      }
      let reload = false;
      Vue.lodash.forEach(action.payload["ids"], (id) => {
        if (this.lodash.startsWith(id, "playbooks/")) {
          reload = true;
        }
      });
      if (reload) {
        this.loadPlaybooks()
      }
    })
  },
});
</script>
