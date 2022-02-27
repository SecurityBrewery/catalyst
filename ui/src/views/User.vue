<template>
  <div>
    <div v-if="user === undefined" class="text-sm-center py-16">
      <v-progress-circular
          indeterminate
          color="primary"
          :size="70"
          :width="7"
          class="align-center"
      >
      </v-progress-circular>
    </div>
    <div v-else class="fill-height d-flex flex-column pa-8">
      <div v-if="$route.params.id === 'new'">
        <h2>Create new API key</h2>
        <v-form>
          <v-text-field label="ID" v-model="user.id" hide-details></v-text-field>
          <v-select multiple chips label="Roles" v-model="user.roles" :items="$store.state.settings.roles"></v-select>
          <v-btn @click="save" color="success" outlined>
            <v-icon>mdi-plus-thick</v-icon>
            Create API-Key
          </v-btn>
        </v-form>
        <v-alert v-if="newUserResponse" color="warning" class="mt-4" dismissible>
          <b>New API-Secret:</b> {{ newUserResponse.secret }}<br>
          Make sure you save it - you won't be able to access it again.
        </v-alert>
      </div>
      <div v-else>
        <h2>
          {{ user.id }}
          <span v-if="user.apikey">(API Key)</span>
        </h2>

        <v-select multiple chips v-if="!user.apikey" label="Roles" v-model="user.roles" :items="$store.state.settings.roles"></v-select>
        <div v-else>
          <v-chip v-for="role in user.roles" :key="role">{{ role }}</v-chip>
        </div>

        <v-btn v-if="!user.apikey" @click="save" color="success" outlined>
          <v-icon>mdi-content-save</v-icon>
          Save
        </v-btn>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import Vue from "vue";

import { NewUserResponse, UserResponse } from "@/client";
import {API} from "@/services/api";

interface State {
  user?: UserResponse;
  newUserResponse?: NewUserResponse;
}

export default Vue.extend({
  name: "User",
  components: {},
  data: (): State => ({
    user: undefined,
    newUserResponse: undefined
  }),
  watch: {
    '$route': function () {
      this.loadUser();
    }
  },
  methods: {
    save() {
      if (!this.user) {
        return
      }

      if (this.$route.params.id == 'new') {
        API.createUser(this.user).then(response => {
          this.newUserResponse = response.data;
        })
      } else {
        API.updateUser(this.$route.params.id, this.user).then(response => {
          this.user = response.data;
        })
      }
    },
    loadUser() {
      if (!this.$route.params.id) {
        return
      }
      if (this.$route.params.id == 'new') {
        this.user = {id: "", roles: [], blocked: false, apikey: true}
      } else {
        API.getUser(this.$route.params.id).then((response) => {
          this.user = response.data;
        });
      }
    },
  },
  mounted() {
    this.loadUser();
  },
});
</script>
