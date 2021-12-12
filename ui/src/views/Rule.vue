<template>
  <div v-if="rule !== undefined" class="flex-grow-1 flex-column d-flex fill-height pa-8">
    <h2 v-if="this.$route.params.id === 'new'">New Rule: {{ rule.name }}</h2>
    <h2 v-else>Edit Rule: {{ rule.name }}</h2>

    <v-text-field label="Name" v-model="rule.name" class="flex-grow-0 flex-shrink-0"></v-text-field>
    <v-text-field label="Condition" v-model="rule.condition" class="flex-grow-0 flex-shrink-0"></v-text-field>

    <v-subheader class="pl-0 py-0" style="height: 20px; font-size: 12px">
      Update
    </v-subheader>
    <div class="flex-grow-1 flex-shrink-1 overflow-scroll">
      <Editor v-model="update" lang="json"></Editor>
    </div>

    <v-row class="px-3 my-6 flex-grow-0 flex-shrink-0">
      <v-btn v-if="this.$route.params.id === 'new'" color="success" @click="save" outlined>
        <v-icon>mdi-plus-thick</v-icon>
        Create
      </v-btn>
      <v-btn v-else color="success" outlined @click="save">
        <v-icon>mdi-content-save</v-icon>
        Save
      </v-btn>
    </v-row>
  </div>
</template>

<script lang="ts">
import Vue from "vue";

import { Rule } from "@/client";
import { API } from "@/services/api";
import Editor from "../components/Editor.vue";

interface State {
  rule?: Rule;
  update: string;
}

export default Vue.extend({
  name: "Rule",
  components: { Editor },
  data: (): State => ({
    rule: undefined,
    update: "",
  }),
  watch: {
    $route: function () {
      this.loadRule();
    },
  },
  methods: {
    save() {
      if (this.rule === undefined) {
        return;
      }
      let rule = this.rule;
      // rule.id = kebabCase(rule.name);
      rule.update = JSON.parse(this.update);
      if (this.$route.params.id == 'new') {
        API.createRule(rule).then(() => {
          this.$store.dispatch("alertSuccess", { name: "Rule created" });
        });
      } else {
        API.updateRule(this.$route.params.id, rule).then(() => {
          this.$store.dispatch("alertSuccess", { name: "Rule saved" });
        });
      }
    },
    loadRule() {
      if (!this.$route.params.id) {
        return
      }
      if (this.$route.params.id == 'new') {
        this.rule = { name: "MyPlaybook", condition: "type == alert", update: { details: { status: "ignored" } } }
      } else {
        API.getRule(this.$route.params.id).then((response) => {
          this.rule = response.data;
          this.update = JSON.stringify(response.data.update, null, 4)
        });
      }
    },
  },
  mounted() {
    this.loadRule();
  },
});
</script>
