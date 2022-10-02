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
        <h2>
          Create a new
          <span v-if="user.apikey">API Key</span>
          <span v-else>User</span>
        </h2>
        <v-form>
          <v-btn-toggle v-model="user.apikey" mandatory dense>
            <v-btn :value="false">User</v-btn>
            <v-btn :value="true">API Key</v-btn>
          </v-btn-toggle>
          <v-text-field label="ID" v-model="user.id" class="mb-2" :rules="[
                v => !!v || 'ID is required',
                v => (v && v.length < 254) || 'ID must be between 1 and 254 characters',
                v => /^[a-z\d\-]+$/.test(v) || 'Only characters a-z, 0-9 and - are allowed',
                // v => /^[A-Za-z0-9_\-\:@\(\)\+,=;\$!\*'%]+$/.test(v) || 'Only characters A-Z, a-z, 0-9, _, -, :, ., @, (, ), +, ,, =, ;, $, !, *, \', % are allowed',
              ]"></v-text-field>
          <v-text-field
              v-if="!user.apikey"
              label="Password"
              v-model="user.password"
              :append-icon="show ? 'mdi-eye' : 'mdi-eye-off'"
              :type="show ? 'text' : 'password'"
              @click:append="show = !show"
              :rules="[
                v => !!v || 'Password is required',
                v => (v && v.length >= 8) || 'Password must be at least 8 characters',
              ]"></v-text-field>
          <v-select multiple chips label="Roles" v-model="user.roles" :items="['analyst', 'engineer', 'admin']"></v-select>
          <v-btn @click="save" color="success" outlined>
            <v-icon>mdi-plus-thick</v-icon>
            Create
            <span v-if="user.apikey" class="ml-1">API Key</span>
            <span v-else class="ml-1">User</span>
          </v-btn>
        </v-form>
        <v-alert v-if="newUserResponse" color="warning" class="mt-4" dismissible>
          <b>New API secret:</b> {{ newUserResponse.secret }}<br>
          Make sure you save it - you won't be able to access it again.
        </v-alert>
      </div>
      <div v-else>
        <h2>
          {{ user.id }}
          <span v-if="user.apikey">(API Key)</span>
        </h2>

        <v-text-field v-if="!user.apikey" label="New Password (leave empty to keep)" v-model="user.password" hide-details class="mb-4"></v-text-field>
        <v-checkbox v-if="!user.apikey" label="Blocked" v-model="user.blocked" hide-details class="mb-4"></v-checkbox>

        <v-select multiple chips v-if="!user.apikey" label="Roles" v-model="user.roles" :items="['analyst', 'engineer', 'admin']"></v-select>
        <div v-else>
          <v-chip v-for="role in user.roles" :key="role" class="mr-1 mb-1">{{ role }}</v-chip>
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

import {NewUserResponse, UserForm} from "@/client";
import {API} from "@/services/api";

interface State {
  show: boolean;
  user?: UserForm;
  newUserResponse?: NewUserResponse;
}

export default Vue.extend({
  name: "User",
  components: {},
  data: (): State => ({
    show: false,
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
