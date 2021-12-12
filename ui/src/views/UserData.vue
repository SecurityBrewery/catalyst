<template>
  <div v-if="userdata">
    <user-data-editor :userdata="userdata" @save="saveUserData"></user-data-editor>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import { UserData } from "@/client";
import { API } from "@/services/api";
import UserDataEditor from "@/components/UserDataEditor.vue";

interface State {
  userdata?: UserData;
}

export default Vue.extend({
  name: "UserData",
  data: (): State => ({
    userdata: undefined,
  }),
  components: { UserDataEditor },
  watch: {
    $route: function () {
      this.loadUserData();
    },
  },
  methods: {
    saveUserData: function(userdata: UserData) {
      API.updateUserData(this.$route.params.id, userdata).then(() => {
        this.$store.dispatch("alertSuccess", { name: "User data saved" });
      });
    },
    loadUserData: function () {
      API.getUserData(this.$route.params.id).then((response) => {
        this.userdata = response.data;
      });
    }
  },
  mounted() {
    this.loadUserData();
  }
});
</script>
