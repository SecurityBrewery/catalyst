<template>
  <v-main style="min-height: 100vh;">
    <List
        :items="userdatas"
        routername="UserData"
        itemid="id"
        itemname="name"
        singular="User Data"
        plural="User Data"
        :show-new="false"
        :deletable="false"
        writepermission="userdata:write"
    ></List>
  </v-main>
</template>

<script lang="ts">
import Vue from "vue";

import {UserData} from "@/client";
import {API} from "@/services/api";
import List from "../components/List.vue";

interface State {
  userdatas: Array<UserData>;
}

export default Vue.extend({
  name: "UserDataList",
  components: {List},
  data: (): State => ({
    userdatas: [],
  }),
  methods: {
    loadSettings() {
      API.listUserData().then((response) => {
        this.userdatas = response.data;
      });
    },
  },
  mounted() {
    this.loadSettings();

    this.$store.subscribeAction((action, state) => {
      if (!action.payload || !(this.lodash.has(action.payload, "ids")) || !action.payload["ids"]) {
        return
      }
      let reload = false;
      Vue.lodash.forEach(action.payload["ids"], (id) => {
        if (this.lodash.startsWith(id, "settings/")) {
          reload = true;
        }
      });
      if (reload) {
        this.loadSettings()
      }
    })
  },
});
</script>
