<template>
  <v-main style="min-height: 100vh;">
    <List
        :items="users"
        routername="User"
        itemid="id"
        itemname="id"
        singular="User / API Key"
        plural="Users / API Keys"
        @delete="deleteUser"
        writepermission="user:write"
    ></List>
  </v-main>
</template>

<script lang="ts">
import Vue from "vue";

import {UserResponse} from "@/client";
import {API} from "@/services/api";
import List from "../components/List.vue";

interface State {
  users: Array<UserResponse>;
}

export default Vue.extend({
  name: "UserList",
  components: {List},
  data: (): State => ({
    users: [],
  }),
  methods: {
    loadUsers() {
      API.listUsers().then((response) => {
        if (response.data) {
          this.users = response.data;
        }
      });
    },
    deleteUser(name: string) {
      API.deleteUser(name).then(() => {
        this.loadUsers();
      });
    },
  },
  mounted() {
    this.loadUsers();

    this.$store.subscribeAction((action, state) => {
      if (!action.payload || !(this.lodash.has(action.payload, "ids")) || !action.payload["ids"]) {
        return
      }
      let reload = false;
      Vue.lodash.forEach(action.payload["ids"], (id) => {
        if (this.lodash.startsWith(id, "users/")) {
          reload = true;
        }
      });
      if (reload) {
        this.loadUsers()
      }
    })
  },
});
</script>
