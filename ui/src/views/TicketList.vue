<template>
  <v-main>
    <TicketListComponent :type="this.$route.params.type" :query="query" @click="open"></TicketListComponent>
  </v-main>
</template>

<script lang="ts">
import Vue from "vue";
import {TicketResponse} from "@/client";

import TicketListComponent from "../components/TicketList.vue";

export default Vue.extend({
  name: "TicketList",
  components: {TicketListComponent},
  props: ['query'],
  methods: {
    open: function (ticket: TicketResponse) {
      if (ticket.id === undefined) {
        return;
      }

      this.$router.push({
        name: "Ticket",
        params: {type: this.$route.params.type, id: ticket.id.toString()}
      });
    }
  }
});
</script>
