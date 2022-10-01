<template>
  <v-main style="min-height: 100vh;">
    <List
        :items="$store.state.templates.templates"
        routername="Template"
        itemid="id"
        itemname="name"
        singular="Template"
        plural="Templates"
        writepermission="template:write"
        @delete="deleteTemplate"
    ></List>
  </v-main>
</template>

<script lang="ts">
import Vue from "vue";

import List from "../components/List.vue";

export default Vue.extend({
  name: "TemplateList",
  components: {List},
  methods: {
    deleteTemplate(title: string) {
      this.$store.dispatch("deleteTemplate", title).then(() => {
        this.$router.push({name: "TemplateList"});
      });
    }
  },
  mounted() {
    this.$store.dispatch("listTemplates");

    this.$store.subscribeAction((action, state) => {
      if (!action.payload || !(this.lodash.has(action.payload, "ids")) || !action.payload["ids"]) {
        return
      }
      let reload = false;
      Vue.lodash.forEach(action.payload["ids"], (id) => {
        if (this.lodash.startsWith(id, "templates/")) {
          reload = true;
        }
      });
      if (reload) {
        this.$store.dispatch("listTemplates");
      }
    })
  }
});
</script>
