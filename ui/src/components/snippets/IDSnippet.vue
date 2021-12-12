<template>
  <div>
    <TicketSnippet v-if="ticket !== undefined" :ticket="ticket" :to="{ name: 'Ticket', params: { type: ticket.type, id: ticket.id } }"></TicketSnippet>
    <ArtifactSnippet v-if="artifact !== undefined" :artifact="artifact" :to="{ name: 'Artifact', params: { artifact: artifact.name } }"></ArtifactSnippet>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import {Artifact, TicketResponse} from "../../client";
import {API} from "@/services/api";
import TicketSnippet from "./TicketSnippet.vue";
import ArtifactSnippet from "./ArtifactSnippet.vue";

interface State {
  ticket?: TicketResponse;
  artifact?: Artifact;
}

export default Vue.extend({
  name: "IDSnippet",
  props: ["id"],
  data: (): State => ({
    ticket: undefined,
    artifact: undefined,
  }),
  components: {
    ArtifactSnippet,
    TicketSnippet
  },
  methods: {
    loadSnippet() {
      if (this.id.startsWith("tickets/")) {
        this.artifact = undefined;
        let ticketID = this.id.replace("tickets/", "")
        API.getTicket(ticketID).then(response => {
          this.ticket = response.data;
        });
      }
      if (this.id.startsWith("artifacts/")) {
        this.ticket = undefined;
        // TODO
        // let artifactID = this.id.replace("artifacts/", "")
        // API.getArtifact(artifactID).then(response => {
        //   this.artifact = response.data;
        // });
      }
    }
  },
  watch: {
    "id": "loadSnippet",
    $route: "loadSnippet"
  },
  mounted() {
    this.loadSnippet();
  }
});
</script>

<style scoped></style>
