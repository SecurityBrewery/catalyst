<template>
  <v-list-item link dense>
    <v-list-item-content @click="goto">
      <v-list-item-title class="d-flex">
        <v-icon small class="mr-1">{{ typeIcon }}</v-icon>
        <span class="text-truncate">{{ ticket.type | capitalize }} #{{ ticket.id }}: {{ ticket.name }}</span>
        <v-spacer></v-spacer>
        <v-icon small class="mr-1">mdi-calendar-plus</v-icon>
        {{ ticket.created | formatdate($store.state.settings.timeformat) }}
      </v-list-item-title>
      <v-list-item-subtitle class="d-flex">
        <v-icon small class="mr-1" :color="statusColor">{{ statusIcon }}</v-icon>
        <span :class="statusColor + '--text'">{{ ticket.status | capitalize }}</span>

        <v-icon small class="mx-1">mdi-account</v-icon>
        {{ ticket.owner ? ticket.owner : 'unassigned' }}
        <v-spacer></v-spacer>
        <v-icon small class="mr-1">mdi-source-branch</v-icon>
        <span class="mr-1">{{ ticket.playbooks ? lodash.size(ticket.playbooks) : 0 }}</span>
        <v-icon small class="mr-1">mdi-checkbox-multiple-marked-outline</v-icon>
        <span class="mr-1">{{ opentaskcount }}</span>
        <v-icon small class="mr-1">mdi-comment</v-icon>
        <span class="mr-1">{{ ticket.comments ? ticket.comments.length : 0 }}</span>
        <v-icon small class="mr-1">mdi-file</v-icon>
        <span class="mr-1">{{ ticket.files ? ticket.files.length : 0 }}</span>
        <v-icon small class="mr-1">mdi-link</v-icon>
        <span class="mr-1">{{ ticket.references ? ticket.references.length : 0 }}</span>
      </v-list-item-subtitle>
    </v-list-item-content>
    <v-list-item-action v-if="action !== ''">
      <v-btn icon small>
        <v-icon small @click="actionClick">{{ action }}</v-icon>
      </v-btn>
    </v-list-item-action>
  </v-list-item>
</template>

<script lang="ts">
import Vue from "vue";
import {Playbook, Task, Type, TypeColorEnum} from "@/client";

export default Vue.extend({
  name: "TicketSnippet",
  props: ["ticket", "to", "action"],
  computed: {
    opentaskcount: function() {
      let count = 0;
      this.lodash.forEach(this.ticket.playbooks, (playbook: Playbook) => {
        this.lodash.forEach(playbook.tasks, (task: Task) => {
          if (task.done) {
            count++;
          }
        })
      })
      return count;
    },
    typeIcon: function () {
      let icon = "mdi-help";
      this.lodash.forEach(this.$store.state.settings.ticketTypes, (ticketType: Type) => {
        if (this.ticket.type === ticketType.id) {
          icon = ticketType.icon
        }
      })
      return icon;
    },
    statusIcon: function() {
      if (this.ticket.status === 'closed') {
        return "mdi-checkbox-marked-circle-outline";
      }
      return "mdi-arrow-right-drop-circle-outline";
    },
    statusColor: function() {
      if (this.ticket.status === 'closed') {
        return TypeColorEnum.Success;
      }
      return TypeColorEnum.Info;
    },
  },
  methods: {
    actionClick: function () {
      this.$emit("actionClick", this.ticket)
    },
    goto: function () {
      if (this.to === undefined) {
        return
      }
      this.$router.push(this.to);
    }
  }
});
</script>
