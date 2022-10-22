<template>
  <span>
    <span v-if="id === undefined">
      <v-icon small class="mr-1">mdi-account</v-icon>
      unassigned
    </span>
    <span v-else-if="user === undefined">
      <v-icon small class="mr-1">mdi-account</v-icon>
      {{ id }}
    </span>
    <span v-else>
      <v-avatar v-if="user.image" :size="lodash.isInteger(size) ? size : 24">
        <v-img :src="user.image"></v-img>
      </v-avatar>
      <v-icon v-else small class="">mdi-account</v-icon>
      {{ user.name ? user.name : id }}
    </span>
  </span>
</template>

<script lang="ts">
import Vue from "vue";
import {UserData} from "@/client";
import {API} from "@/services/api";
import {AxiosResponseTransformer} from "axios";

interface State {
  user?: UserData,
}

export default Vue.extend({
  name: "User",
  props: ["id", "size"],
  data: (): State => ({
    user: undefined,
  }),
  watch: {
    id: function(): void {
      this.loadUserData();
    }
  },
  methods: {
    loadUserData: function () {
      if (this.id === undefined) {
        this.user = undefined;
        return
      }

      let defaultTransformers = this.axios.defaults.transformResponse as AxiosResponseTransformer[]
      let transformResponse = defaultTransformers.concat((data) => {
        data.notoast = true;
        return data
      });
      API.getUserData(this.id, {transformResponse: transformResponse}).then(response => {
        this.user = response.data;
      })
    }
  },
  mounted() {
    this.loadUserData();
  }
});
</script>
