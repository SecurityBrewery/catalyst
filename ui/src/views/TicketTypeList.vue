<template>
  <v-main style="min-height: 100vh;">
    <List
        :items="tickettypes"
        routername="TicketType"
        itemid="id"
        itemname="name"
        singular="Ticket Type"
        plural="Ticket Types"
        @delete="deleteTicketType"
        writepermission="tickettype:write"
    ></List>
  </v-main>
</template>

<script lang="ts">
import Vue from "vue";

import { TicketType} from "@/client";
import {API} from "@/services/api";
import List from "../components/List.vue";

interface State {
  tickettypes: Array<TicketType>;
}

export default Vue.extend({
  name: "TicketTypeList",
  components: {List},
  data: (): State => ({
    tickettypes: [],
  }),
  methods: {
    loadTicketType() {
      API.listTicketTypes().then((response) => {
        this.tickettypes = response.data;
      });
    },
    deleteTicketType(id: string) {
      API.deleteTicketType(id).then(() => {
        this.loadTicketType();
      });
    },
  },
  mounted() {
    this.loadTicketType();

    this.$store.subscribeAction((action, state) => {
      if (!action.payload || !(this.lodash.has(action.payload, "ids")) || !action.payload["ids"]) {
        return
      }
      let reload = false;
      Vue.lodash.forEach(action.payload["ids"], (id) => {
        if (this.lodash.startsWith(id, "tickettypes/")) {
          reload = true;
        }
      });
      if (reload) {
        this.loadTicketType()
      }
    })
  },
});
</script>
