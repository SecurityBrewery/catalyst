<template>
  <div v-if="automation !== undefined" class="flex-grow-1 flex-column d-flex fill-height pa-8">
    <v-row class="flex-grow-0 flex-shrink-0">
      <v-spacer></v-spacer>
      <v-btn href="https://catalyst-soar.com/docs/catalyst/engineer/automations" target="_blank" outlined rounded small>
        <v-icon>mdi-book-open</v-icon> Handbook
      </v-btn>
    </v-row>

    <v-alert v-if="readonly" type="info">You do not have write access to automations.</v-alert>
    <h2 v-if="readonly">Automation: {{ automation.id }}</h2>
    <h2 v-else-if="this.$route.params.id === 'new'">New Automation: {{ automation.id }}</h2>
    <h2 v-else>Edit Automation: {{ automation.id }}</h2>

    <v-row class="flex-grow-0 flex-shrink-0">
      <v-col cols="12">
        <v-text-field :readonly="readonly" label="ID" v-model="automation.id" class="flex-grow-0 flex-shrink-0"></v-text-field>

        <v-select
            label="Type"
            :items="types"
            item-text="id"
            return-object
            multiple
            v-model="automation.type"
        ></v-select>

        <v-text-field :readonly="readonly" label="Docker Image" v-model="automation.image" class="flex-grow-0 flex-shrink-0"></v-text-field>

        <v-subheader class="pl-0 py-0" style="height: 20px; font-size: 12px">Script</v-subheader>
        <div class="flex-grow-1 flex-shrink-1 overflow-scroll">
          <Editor v-model="automation.script" lang="python" :readonly="readonly" ></Editor>
        </div>

        <AdvancedJSONSchemaEditor
            @save="save"
            :schema="automation.schema ? JSON.parse(automation.schema) : {}"
            hidepreview
            class="mt-4"
            :readonly="readonly"></AdvancedJSONSchemaEditor>
      </v-col>
    </v-row>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Editor from "../components/Editor.vue";
import AdvancedJSONSchemaEditor from "../components/AdvancedJSONSchemaEditor.vue";
import {API} from "@/services/api";
import {AutomationResponse, AutomationForm, AutomationResponseTypeEnum} from "@/client";
import {DateTime} from "luxon";

interface State {
  automation?: AutomationResponse;
  data?: any;
  valid: boolean;
}

export default Vue.extend({
  name: "Automation",
  components: { Editor, AdvancedJSONSchemaEditor },
  data: (): State => ({
    automation: undefined,
    data: undefined,
    valid: true,
  }),
  watch: {
    $route: function () {
      this.loadAutomation();
    },
  },
  computed: {
    readonly: function (): boolean {
      return !this.hasRole("engineer:automation:write");
    },
    types: function (): Array<string> {
      return [ AutomationResponseTypeEnum.Global, AutomationResponseTypeEnum.Playbook, AutomationResponseTypeEnum.Artifact ]
    }
  },
  methods: {
    save(schema) {
      if (this.automation === undefined) {
        return;
      }
      let automation = this.automation as any as AutomationForm;
      automation.schema = schema;
      if (this.$route.params.id == "new") {
        API.createAutomation(automation);
      } else {
        API.updateAutomation(this.$route.params.id, automation);
      }
    },
    loadAutomation() {
      if (!this.$route.params.id) {
        return;
      }
      if (this.$route.params.id == "new") {
        this.automation = { id: "my-automation", image: "docker.io/ubuntu", script: "", type: [ AutomationResponseTypeEnum.Global, AutomationResponseTypeEnum.Playbook, AutomationResponseTypeEnum.Artifact ] };
      } else {
        API.getAutomation(this.$route.params.id).then((response) => {
          this.automation = response.data;
        });
      }
    },
    hasRole: function (s: string): boolean {
      if (this.$store.state.user.roles) {
        return this.lodash.includes(this.$store.state.user.roles, s);
      }
      return false;
    },
    timeformat: function(s: string, locale: string) {
      let format = this.$store.state.settings.timeformat;
      if (!format) {
        return DateTime.fromISO(s).toLocaleString(DateTime.DATETIME_SHORT);
      }
      return DateTime.fromISO(s).toFormat(format);
    },
    dateformat: function(s: string, locale: string) {
      let format = this.$store.state.settings.timeformat;
      if (!format) {
        return DateTime.fromISO(s).toLocaleString(DateTime.DATETIME_SHORT);
      }
      return DateTime.fromISO(s).toFormat(format);
    },
    datetimeformat: function(s: string, locale: string) {
      let format = this.$store.state.settings.timeformat;
      if (!format) {
        return DateTime.fromISO(s).toLocaleString(DateTime.DATETIME_SHORT);
      }
      return DateTime.fromISO(s).toFormat(format);
    },
  },
  mounted() {
    this.loadAutomation();
  },
});
</script>
