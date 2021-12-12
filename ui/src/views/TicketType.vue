<template>
  <div v-if="tickettype !== undefined" class="flex-grow-1 flex-column d-flex fill-height pa-8">
    <h2 v-if="this.$route.params.id === 'new'">New Ticket Type: {{ tickettype.name }}</h2>
    <h2 v-else>Edit Ticket Type: {{ tickettype.name }}</h2>

    <v-text-field label="Name" v-model="tickettype.name" class="flex-grow-0 flex-shrink-0"></v-text-field>

    <v-subheader class="pl-0 py-0" style="height: 20px; font-size: 12px">
      Icon (see <a href="https://materialdesignicons.com" class="mx-1"> materialdesignicons.com </a> and prefix with "mdi-")
    </v-subheader>
    <v-text-field v-model="tickettype.icon" placeholder="e.g. mdi-alert" class="flex-grow-0 flex-shrink-0 mt-n3"></v-text-field>


    <v-select
        v-model="selectedTemplate"
        :items="templates"
        item-text="name"
        label="Default Template"
        return-object
        class="flex-grow-0 flex-shrink-0"
    ></v-select>

    <v-select
        v-model="selectedPlaybooks"
        :items="playbooks"
        item-text="name"
        label="Default Playbooks"
        return-object
        multiple
        class="flex-grow-0 flex-shrink-0"
    ></v-select>

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

import { PlaybookTemplateResponse, TicketTemplateResponse, TicketType } from "@/client";
import { API } from "@/services/api";

interface State {
  tickettype?: TicketType;
  selectedPlaybooks: Array<PlaybookTemplateResponse>;
  playbooks?: Array<PlaybookTemplateResponse>;
  selectedTemplate?: TicketTemplateResponse;
  templates?: Array<TicketTemplateResponse>;
}

export default Vue.extend({
  name: "TicketType",
  data: (): State => ({
    tickettype: undefined,
    selectedPlaybooks: [],
    playbooks: undefined,
    selectedTemplate: undefined,
    templates: undefined,
  }),
  watch: {
    $route: function () {
      this.loadRule();
    },
  },
  methods: {
    save() {
      if (this.tickettype === undefined) {
        return;
      }
      let tickettype = this.tickettype;
      if (this.selectedTemplate) {
        tickettype.default_template = this.selectedTemplate.id;
      }
      tickettype.default_playbooks = [];
      this.lodash.forEach(this.selectedPlaybooks, (playbook) => {
        tickettype.default_playbooks.push(playbook.id)
      })

      if (this.$route.params.id == 'new') {
        API.createTicketType(tickettype).then(() => {
          this.$store.dispatch("alertSuccess", { name: "Ticket type created" });
        });
      } else {
        API.updateTicketType(this.$route.params.id, tickettype).then(() => {
          this.$store.dispatch("alertSuccess", { name: "Ticket type saved" });
        });
      }
    },
    loadRule() {
      if (!this.$route.params.id) {
        return
      }

      API.listTemplates().then((response) => {
        this.templates = response.data;
        if (response.data.length > 0) {
          this.selectedTemplate = response.data[0];
        }
      })
      API.listPlaybooks().then((response) => {
        this.playbooks = response.data;
      })

      if (this.$route.params.id == 'new') {
        this.tickettype = { name: "TI Ticket", default_groups: [], default_playbooks: [], default_template: "", icon: "" }
      } else {
        API.getTicketType(this.$route.params.id).then((response) => {
          this.tickettype = response.data;
        });
      }
    },
  },
  mounted() {
    this.loadRule();
  },
});
</script>
