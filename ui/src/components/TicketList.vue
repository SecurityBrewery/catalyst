<template>
  <div>
    <v-row v-if="selectedtype" class="mx-0 my-2" dense>
      <v-col :cols="this.$route.params.id ? 4 : 2">
        <v-select
            v-model="selectedtype"
            :items="tickettypes"
            item-text="name"
            item-value="id"
            solo
            rounded
            label="Type"
            dense
            flat
            hide-details
            height="48px"></v-select>
      </v-col>
      <v-col :cols="this.$route.params.id ? 8 : 10">
        <v-btn elevation="0" rounded class="float-right mb-2" @click="opennew({ name: 'Ticket', params: { type: selectedtype, id: 'new' } })">
          <v-icon class="mr-1">mdi-plus</v-icon>
          New {{ selectedtype | capitalize }}
        </v-btn>
      </v-col>
    </v-row>

    <v-toolbar
        id="caqlbar"
        rounded
        filled
        dense
        flat
        elevation="0"
        style="border-radius: 40px !important;"
    >
      <v-btn-toggle dense v-model="defaultcaql">
        <v-btn
            text
            color="primary"
            @click="caql = !caql"
            rounded
            style="border-radius: 40px !important;">
          CAQL
        </v-btn>
      </v-btn-toggle>

      <v-text-field
          placeholder="Search term or query (e.g. name == 'malware' AND 'wannacry')"
          v-model="term"
          dense
          solo
          flat
          hide-details
          clearable
          append-icon="mdi-magnify"
          @click:clear="clear"
          @click:append="loadTickets"
          @keydown.enter="loadTickets"
          :rules="[validate]"
          @focus="focus = true"
          @blur="blur"
      ></v-text-field>
    </v-toolbar>

    <span v-if="focus && caql">CAQL Query Suggestions</span>
    <v-list class="mb-2" v-if="focus && caql">
      <v-list-item v-for="example in examples" :key="example.q" dense link @click="term = example.q; caql = true; defaultcaql = 0; loadTickets()">
        <v-list-item-content>
          <v-list-item-title>
            <v-row>
            <span class="col-6">{{ example.q }}</span> <span class="text--disabled col-6">{{ example.desc }}</span>
            </v-row>
          </v-list-item-title>
        </v-list-item-content>
      </v-list-item>
      <v-list-item>
        <v-list-item-content>
          Fields: {{ lodash.join(fields, ", ") }}
        </v-list-item-content>
      </v-list-item>
    </v-list>

    <v-data-table
      :headers="headers"
      :items="tickets"
      item-key="name"
      multi-sort
      class="elevation-0 cards clickable mt-2"
      :options.sync="options"
      :server-items-length="totalTickets"
      :loading="loading"
      :footer-props="{ 'items-per-page-options': [10, 25, 50, 100] }"
    >
      <template v-slot:item="{ item }">
        <tr @click="open(item)">
          <td colspan="5" class="pa-0">
            <v-list-item :to="{ name: 'Ticket', params: { type: item.type, id: item.id } }" class="pa-0" style="background: none">
              <ticketSnippet :ticket="item"></ticketSnippet>
            </v-list-item>
          </td>
        </tr>
      </template>
    </v-data-table>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import { Ticket, TicketType } from "@/client";
import { API } from "@/services/api";
import TicketSnippet from "../components/snippets/TicketSnippet.vue";
import {validateCAQL} from "@/suggestions/suggestions";
import {DateTime} from "luxon";

interface State {
  term: string;
  loading: boolean;
  tickets: Array<Ticket>;
  totalTickets: number;
  options: {
    page?: number;
    itemsPerPage?: number;
    sortBy?: string[];
    sortDesc?: boolean[];
    groupBy?: string[];
    groupDesc?: boolean[];
    multiSort?: boolean;
    mustSort?: boolean;
  };
  focus: boolean;
  caql: boolean;
  defaultcaql?: number;

  tickettypes: Array<TicketType>;
  selectedtype: string;
}

interface QuerySuggestion {
  q: string;
  desc: string;
}

export default Vue.extend({
  name: "TicketList",
  components: {
    TicketSnippet,
  },
  props: ["type", "query"],
  data: (): State => ({
    term: "status == 'open'",
    loading: true,
    tickets: [],
    totalTickets: 0,
    options: {
      itemsPerPage: 10,
    },
    focus: false,
    caql: true,
    defaultcaql: 0,
    tickettypes: [],

    selectedtype: "",
  }),
  computed: {
    fields(): Array<string> {
      return [
        "type", "id", "name",
        "status" , "owner",
        "created", "modified",
        "details", "details.description", "details.…",
        "schema",
        "comments", "comments.#.created", "comments.#.creator", "comments.#.message",
        "playbooks", "playbooks.#.name", "playbooks.#.tasks",
        "references", "references.#.href", "references.#.name",
        "artifacts", "artifacts.#.name", "artifacts.#.status", "artifacts.#.type",
        "files", "files.#.name"
      ];
    },
    user(): string {
      return this.$store.state.user.id
    },
    headers() {
      return [
        {
          text: "Name",
          align: "start",
          value: "name",
        },
        {
          text: "Status",
          align: "start",
          value: "status",
        },
        {
          text: "Owner",
          align: "start",
          value: "owner",
        },
        {
          text: "Creation",
          align: "start",
          value: "created",
        },
        {
          text: "Last Modification",
          align: "start",
          value: "modified",
        },
      ];
    },
    examples (): Array<QuerySuggestion> {
      let twoWeeksAgo = DateTime.utc().minus({weeks: 2}).toFormat("yyyy-MM-dd");

      let ex: Array<QuerySuggestion> = [];

      if (this.user) {
        ex.push({q: "status == 'open' AND (owner == '" + this.user + "' OR !owner)", desc: "Select all open tickets by you and unassigned"})
        ex.push({q: "status == 'closed' AND owner == '" + this.user + "'", desc: "Select completed tickets by you"})
      } else {
        ex.push({q: "status == 'open'", desc: "Select all open tickets"})
        ex.push({q: "status == 'closed'", desc: "Select completed tickets"})
      }

      ex.push({q: "created > \""+twoWeeksAgo+"\"", desc: "Select tickets created in the last two weeks"})

      if (this.term && this.term.match(/^[A-Za-z ]+$/)) {
        ex.unshift({q: "'" + this.term + "'", desc: "Full text search for '" + this.term + "'"})
      }
      return ex;
    },
  },
  watch: {
    options: {
      handler() {
        this.loadTickets();
      },
      deep: true,
    },
    selectedtype: function () {
      this.loadTickets();
    },
    $route: function () {
      this.selectedtype = this.type;
    },
  },
  methods: {
    blur: function () {
      setTimeout(()=>{this.focus = false}, 200)
    },
    clear: function () {
      this.caql = false;
      this.defaultcaql = undefined;
      this.term = "";
      this.loadTickets();
    },
    open: function (ticket: Ticket) {
      this.$emit("click", ticket);
    },
    opennew: function (to) {
      this.$router.push(to).then(() => {
        this.$emit("new");
      })
    },
    select: function (e: string) {
      this.loadTerm(e);
    },
    loadTickets() {
      let term = this.term;
      if (!term) {
        term = "";
      }
      this.loadTerm(term);
    },
    loadTicketTypes() {
      API.listTicketTypes().then((reponse) => {
        this.tickettypes = reponse.data;
      })
    },
    loadTerm(term: string) {
      this.loading = true;
      let offset = 0;
      let count = 25;
      let sortBy: Array<string> = [];
      let sortDesc: Array<boolean> = [];
      if (this.options.itemsPerPage !== undefined) {
        count = this.options.itemsPerPage;
        if (this.options.page !== undefined) {
          offset = (this.options.page - 1) * this.options.itemsPerPage;
        }
      }
      if (this.options.sortBy !== undefined) {
        sortBy = this.options.sortBy;
      }
      if (this.options.sortDesc !== undefined) {
        sortDesc = this.options.sortDesc;
      }

      let ticketType = this.selectedtype;
      if (!ticketType) {
        ticketType = "";
      }

      if (!this.caql && term.length > 0) {
        term = "'" + this.lodash.join(this.lodash.split(term, " "), "'&&'") + "'"
      }

      API.listTickets(ticketType, offset, count, sortBy, sortDesc, term)
          .then((response) => {
            if (response.data.tickets) {
              this.tickets = response.data.tickets;
            } else {
              this.tickets = [];
            }
            this.totalTickets = response.data.count;
            this.loading = false;
          })
          .catch(() => {
            this.loading = false;
          });

    },
    validate: function () {
      if (!this.term) {
        return true
      }
      let err = validateCAQL(this.term);
      if (err !== null) {
        return err.message;
      }
      return true;
    }
  },
  mounted() {
    if (this.user) {
      this.term = "status == 'open' AND (owner == '" + this.user + "' OR !owner)";
    } else {
      this.term = "status == 'open'";
    }

    if (this.query) {
      this.term = this.query;
    }

    this.selectedtype = this.type;
    this.loadTicketTypes();
    this.loadTickets();
  },
});
</script>

<style>
.clickable td {
  cursor: pointer !important;
}
.vue-simple-suggest.designed .input-wrapper input {
  background: none !important;
  border: none !important;
  color: #fff;
}
.vue-simple-suggest.designed .suggestions {
  background-color: #333 !important;
  top: 60px !important;
  border: 1px solid #ddd !important;
}
</style>
