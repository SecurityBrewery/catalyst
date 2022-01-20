<template>
  <v-main style="min-height: 100vh">
    <splitpanes class="default-theme" @resize="paneSize = $event[0].size">
      <pane class="pa-3" :size="paneSize">
        <TicketListComponent :type="$route.params.type" :query="query" @click="open" @new="opennew"></TicketListComponent>
      </pane>
      <pane v-if="this.$route.params.id" class="pa-3" :size="100 - paneSize">
        <v-row>
          <v-spacer></v-spacer>
          <v-btn @click="close" outlined rounded class="mt-3 mr-3">
            <v-icon>mdi-close</v-icon>
            Close
          </v-btn>
        </v-row>
        <router-view></router-view>
      </pane>
    </splitpanes>
  </v-main>
</template>

<script lang="ts">
import Vue from "vue";
import {TicketResponse} from "@/client";
import { Splitpanes, Pane } from 'splitpanes';
import 'splitpanes/dist/splitpanes.css';

import TicketListComponent from "../components/TicketList.vue";

export default Vue.extend({
  name: "TicketList",
  components: {
    TicketListComponent,
    Splitpanes, Pane,
  },
  props: ['query', 'type'],
  data: () => ({
    paneSize: 30,
  }),
  methods: {
    numeric: function (n: any): boolean {
      return !isNaN(parseFloat(n)) && !isNaN(n - 0);
    },
    open: function (ticket: TicketResponse) {
      this.$router.push({
        name: "Ticket",
        params: {type: this.$route.params.type, id: ticket.id.toString()}
      }).then(() => {
        this.paneSize = 30;
      });
    },
    opennew: function () {
      this.paneSize = 30;
    },
    close: function () {
      this.$router.push({
        name: "TicketList",
        params: {type: this.$route.params.type},
      }).then(() => {
        this.paneSize = 100;
      });
    },
    hidelist: function () {
      this.paneSize = 0;
    }
  }
});
</script>

<style lang="scss">
.splitpanes.default-theme .splitpanes__pane {
  background: none;
}
.splitpanes.default-theme .splitpanes__splitter {
  background: none;
  border: none;
}
</style>
