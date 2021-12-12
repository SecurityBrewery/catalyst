<template>
  <div>
    <div v-if="$route.params.id == 'new'" class="fill-height d-flex flex-column pa-8">
      <v-alert v-if="readonly" type="info">
        You do not have write access to jobs.
      </v-alert>
      <h2>Start new job</h2>

      <v-select
          label="Automations"
          :items="globalautomations"
          item-text="id"
          return-object
          v-model="automation"
      ></v-select>

      <v-form v-if="automation" v-model="valid">
        <v-jsf
            v-model="data"
            :schema="JSON.parse(automation.schema)"
            :options="{ readonly: true, formats: { time: timeformat, date: dateformat, 'date-time': datetimeformat } }"
        />
        <v-btn @click="run" color="success" outlined :disabled="!valid">
          Run
        </v-btn>
      </v-form>
    </div>
    <div v-else-if="job !== undefined">
      <h2>Job: {{ job.id }}</h2>

      <div v-if="job.automation" class="flex-grow-0 flex-shrink-0">
        <v-subheader class="pl-0 py-0" style="height: 20px; font-size: 12px">Automation</v-subheader>
        {{ job.automation }}
      </div>

      <div v-if="job.container" class="flex-grow-0 flex-shrink-0">
        <v-subheader class="pl-0 py-0" style="height: 20px; font-size: 12px">Container</v-subheader>
        {{ job.container }}
      </div>

      <div v-if="job.status" class="flex-grow-0 flex-shrink-0">
        <v-subheader class="pl-0 py-0" style="height: 20px; font-size: 12px">Status</v-subheader>
        {{ job.status }}
      </div>

      <div v-if="job.payload" class="flex-grow-0 flex-shrink-0">
        <v-subheader class="pl-0 py-0" style="height: 20px; font-size: 12px">Input Payload</v-subheader>
        <Editor :value="JSON.stringify(job.payload, null, 2)" lang="json" :readonly="true"></Editor>
      </div>

      <div v-if="job.log" class="flex-grow-0 flex-shrink-0">
        <v-subheader class="pl-0 py-0" style="height: 20px; font-size: 12px">Log</v-subheader>
        <Editor :value="job.log" lang="log" :readonly="true"></Editor>
      </div>

      <div v-if="job.output" class="flex-grow-0 flex-shrink-0">
        <v-subheader class="pl-0 py-0" style="height: 20px; font-size: 12px">Output</v-subheader>
        <Editor :value="JSON.stringify(job.output, null, 2)" lang="json" :readonly="true"></Editor>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import Vue from "vue";

import {AutomationResponse, AutomationResponseTypeEnum, JobResponse} from "@/client";
import {API} from "@/services/api";
import {DateTime} from "luxon";
import Editor from "@/components/Editor.vue";

interface State {
  job?: JobResponse,
  data?: any;
  valid: boolean;

  automation?: AutomationResponse;
  automations: Array<AutomationResponse>;
}

export default Vue.extend({
  name: "Job",
  components: { Editor },
  data: (): State => ({
    job: undefined,
    data: undefined,
    valid: true,

    automation: undefined,
    automations: [],
  }),
  watch: {
    $route: 'loadJob',
  },
  computed: {
    readonly: function (): boolean {
      return !this.hasRole("admin:job:write");
    },
    globalautomations: function (): Array<AutomationResponse> {
      if (!this.automations) {
        return [];
      }
      return this.lodash.filter(this.automations, (automation: AutomationResponse) => {
        if (!automation || !automation.type) {
          return true;
        }
        return this.lodash.includes(automation.type, AutomationResponseTypeEnum.Global)
      })
    },
  },
  methods: {
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
    run: function () {
      if (!this.automation) {
        return;
      }
      API.runJob({ automation: this.automation.id, payload: this.data }).then(() => {
        this.$store.dispatch("alertSuccess", { name: "Job started." });
      })
    },
    loadJob() {
      if (!this.$route.params.id) {
        return
      }
      if (this.$route.params.id == 'new') {
        // this.template = { name: "MyTemplate", schema: '{ "type": "object", "name": "Incident" }' }
        // this.schema = { type: "object", name: "Incident" }
      } else {
        API.getJob(this.$route.params.id).then((response) => {
          this.job = response.data;
        });
      }
    },
    loadAutomations() {
      API.listAutomations(this.$route.params.id).then((response) => {
        this.automations = response.data;
      });
    },
    hasRole: function (s: string): boolean {
      if (this.$store.state.user.roles) {
        return this.lodash.includes(this.$store.state.user.roles, s);
      }
      return false;
    }
  },
  mounted() {
    this.loadJob();
    this.loadAutomations();

    this.$store.subscribeAction((action, state) => {
      if (!action.payload || !(this.lodash.has(action.payload, "ids")) || !action.payload["ids"]) {
        return
      }
      let reload = false;
      Vue.lodash.forEach(action.payload["ids"], (id) => {
        if (this.lodash.startsWith(id, "jobs/")) {
          reload = true;
        }
      });
      if (reload) {
        this.loadJob()
      }
    })
  },
});
</script>
