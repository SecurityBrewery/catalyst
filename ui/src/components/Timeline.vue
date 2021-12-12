<template>
  <div>
    <v-textarea
      v-model="comment"
      hide-details
      flat
      label="Add a comment..."
      solo
      auto-grow
      rows="2"
      class="py-2"
    >
      <template v-slot:append>
        <v-btn class="mx-0 mt-n1" text @click="addTicketLog">
          <v-icon>mdi-send</v-icon>
        </v-btn>
      </template>
    </v-textarea>
    <div
      v-for="(log, id) in internalLogs"
      :key="id"
      :icon="icon(log.type)"
      :small="small(log.type)"
      class="pb-2"
    >
      <v-card v-if="log.type === 'comment'" elevation="0" color="cards">
        <v-card-subtitle class="pb-0">
          <strong> {{ log.creator }}</strong>
          <span class="text--disabled ml-3" :title="log.created">
            {{ relDate(log.created) }}
          </span>
        </v-card-subtitle>
        <v-card-text class="mb-0 mt-2">
          <!--{{ log.message }}-->
          <vue-markdown v-if="show">{{ log.message }}</vue-markdown>
        </v-card-text>
      </v-card>
      <div v-else style="line-height: 24px" class="d-flex flex-row">
        <v-divider class="mt-3 mr-3"></v-divider>
        {{ log.message }}
        <span class="text--disabled ml-1" :title="log.created">
          &middot;
          {{ log.creator }} &middot;
          {{ relDate(log.created) }}
        </span>
        <v-divider class="mt-3 ml-3"></v-divider>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import { LogEntry } from "@/client";
import VueMarkdown from "vue-markdown";
import { API } from "@/services/api";

interface State {
  comment: string
  internalLogs: Array<LogEntry>
  show: boolean
}

export default Vue.extend({
  name: "Timeline",
  components: {
    "vue-markdown": VueMarkdown
  },
  props: ["id", "logs"],
  data: (): State => ({
    comment: "",
    internalLogs: [],
    show: true
  }),
  watch: {
    logs: function () {
      // this.internalLogs = this.logs;
      this.reload(this.logs);
    }
  },
  methods: {
    reload: function(newlogs: Array<LogEntry>) {
      if (newlogs === undefined) {
        return
      }
      this.show = false;
      Vue.nextTick(() => {
        this.internalLogs = newlogs;
        this.show = true;
      })
    },
    icon: function(s: string) {
      switch (s) {
        case "comment":
          return "mdi-comment";
      }
      return "";
    },
    small: function(s: string) {
      switch (s) {
        case "comment":
          return false;
      }
      return true;
    },
    relDate: function(date: string) {
      let rtf = new Intl.RelativeTimeFormat("en", { numeric: "auto" });
      let deltaDays =
        (new Date(date).getTime() - new Date().getTime()) / (1000 * 3600 * 24);
      let relDate = rtf.format(Math.round(deltaDays), "days");
      if (deltaDays > -3) {
        relDate +=
          ", " +
          new Date(date).toLocaleTimeString([], {
            hour: "2-digit",
            minute: "2-digit"
          });
      }
      return relDate;
    },
    addTicketLog() {
      // API.addLog({ id: this.id, message: this.comment }).then(
      //   response => {
      //     this.$store.dispatch("alertSuccess", { name: "Log saved", type: "success" });
      //     if (this.internalLogs === undefined) {
      //       this.reload([response.data]);
      //     } else {
      //       this.internalLogs.unshift(response.data);
      //       this.reload(this.internalLogs);
      //     }
      //   }
      // );
    }
  },
  mounted() {
    this.internalLogs = this.logs;
  }
});
</script>
