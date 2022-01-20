<template>
  <div class="mt-8">
    <h2>New {{ $route.params.type | capitalize }}</h2>
    <v-form class="create clearfix">
      <v-text-field label="Title" v-model="name"></v-text-field>

      <v-select
          label="Playbooks"
          :items="playbooks"
          item-text="name"
          return-object
          multiple
          v-model="selectedPlaybooks"
      ></v-select>

      <v-select
        label="Template"
        :items="templates"
        item-text="name"
        return-object
        v-model="selectedTemplate"
      ></v-select>

      <v-subheader class="pl-0 mt-4" style="height: 20px">Details</v-subheader>
      <div v-if="selectedTemplate !== undefined" class="details">
        <v-lazy>
          <v-jsf
            v-model="details"
            :options="{ fieldProps: { 'hide-details': true } }"
            :schema="selectedSchema"
          />
        </v-lazy>
      </div>

      <v-btn color="green" @click="createTicket" outlined>
        Create
      </v-btn>
    </v-form>
  </div>
</template>

<script lang="ts">
import Vue from "vue";

import {
  PlaybookTemplateForm,
  PlaybookTemplateResponse,
  TicketForm,
  TicketTemplateResponse,
  TicketType
} from "@/client";
import { API } from "@/services/api";

interface State {
  selectedTemplate?: TicketTemplateResponse;
  templates: Array<TicketTemplateResponse>;

  selectedPlaybooks: Array<PlaybookTemplateResponse>;
  playbooks: Array<PlaybookTemplateResponse>;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  details: any;
  name: string;
  ticketType?: TicketType;
}

export default Vue.extend({
  name: "TicketNew",
  // components: { VJsf },
  data: (): State => ({
    selectedTemplate: undefined,
    templates: [],
    selectedPlaybooks: [],
    playbooks: [],
    name: "",
    details: {},
    ticketType: undefined,
  }),
  computed: {
    selectedSchema() {
      if (
        this.selectedTemplate === undefined ||
        this.selectedTemplate.schema === undefined
      ) {
        return {};
      }
      return JSON.parse(this.selectedTemplate.schema);
    }
  },
  methods: {
    createTicket: function() {
      if (
        this.selectedTemplate === undefined ||
        this.selectedTemplate?.schema === undefined
      ) {
        this.$store.dispatch("alertError", { name: "No template selected" });
        return;
      }

      let ticket: TicketForm = {
        type: this.$route.params.type,
        name: this.name,
        details: this.details,
        schema: this.selectedTemplate?.schema,
        status: "open",
        owner: this.$store.state.user.id,
        playbooks: this.selectedPlaybooks as Array<PlaybookTemplateForm>,
      };

      API.createTicket(ticket).then(response => {
        if (response.data.id === undefined) {
          return;
        }
        this.$router.push({
          name: "Ticket",
          params: {
            id: response.data.id.toString(),
            type: response.data.type,
          }
        });
      });
    }
  },
  mounted() {
    API.getTicketType(this.$route.params.type).then((response) => {
      this.ticketType = response.data;

      API.listTemplates().then(response => {
        this.templates = response.data;
        this.lodash.forEach(this.templates, (template) => {
          if (this.ticketType && template.id == this.ticketType.default_template) {
            this.selectedTemplate = template;
          }
        });
      });

      API.listPlaybooks().then(response => {
        this.playbooks = response.data;
        this.lodash.forEach(this.playbooks, (playbook) => {
          if (this.ticketType && this.lodash.includes(this.ticketType.default_playbooks, playbook.id)) {
            this.selectedPlaybooks.push(playbook);
          }
        });
      });
    });
  },
});
</script>

<style>
.details > .col > .row > .row,
.details > .col > .row > div > .row {
  margin-top: 16px !important;
  padding: 16px !important;
  background-color: rgba(0, 0, 0, 0.2);
}

.v-application .create .vjsf-property {
  margin-bottom: 16px;
}

.v-application .create .error--text {
  color: white !important;
}
</style>
