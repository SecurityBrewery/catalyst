<template>
  <v-container>
    <h1>{{ name }}</h1>

    <v-divider class="my-2"></v-divider>

    <h2 class="text--disabled" style="font-size: 12pt" v-if="artifact">
      Status:
      <v-menu offset-y class="mr-2">
        <template v-slot:activator="{ on, attrs }">
          <span v-bind="attrs" v-on="on">
            <v-icon small class="mr-1" :color="statusColor(artifact.status)">{{ statusIcon(artifact.status) }}</v-icon>
            <span :class="statusColor(artifact.status) + '--text'">{{ artifact.status | capitalize }}</span>
          </span>
        </template>
        <v-list>
          <v-list-item dense link v-for="state in otherStates" :key="state.id" @click="setStatus(state.id)">
            <v-list-item-title>
              Set status to <v-icon small>{{ statusIcon(state.id) }}</v-icon> {{ state.name }}
            </v-list-item-title>
          </v-list-item>
        </v-list>
      </v-menu>
      &middot;
      Type:
      <v-menu
          :close-on-content-click="false"
          offset-y
          class="mr-2">
        <template v-slot:activator="{ on, attrs }">
          <span v-bind="attrs" v-on="on">
            {{ artifact.type ? artifact.type : "empty" }}
          </span>
        </template>
        <v-combobox
            class="pt-6 pb-0 px-3"
            style="background-color: white"
            v-model="artifacttype"
            :items="['ip', 'domain', 'md5', 'sha1', 'sha256', 'filename', 'url']"
            label="Artifact Type"
            @change="setArtifactType"
        ></v-combobox>
      </v-menu>
    </h2>

    <v-divider class="mt-0 mb-4"></v-divider>

    <div v-if="automations">
      <v-card
          v-for="automation in artifactautomations" :key="automation.id" class="mb-4" elevation="0" outlined>
        <v-card-title>
          {{ automation.id }}
          <v-spacer></v-spacer>
          <v-btn
              v-if="artifact && artifact.enrichments && automation.id in artifact.enrichments"
              @click="enrich(automation.id)"
              elevation="0"
              >
            <v-icon>mdi-sync</v-icon>
          </v-btn>
          <v-btn v-else @click="enrich(automation.id)" elevation="0">
            <v-icon>mdi-play</v-icon>
          </v-btn>
        </v-card-title>
        <v-card-text v-if="artifact && artifact.enrichments && automation.id in artifact.enrichments && artifact.enrichments[automation.id].data">
          <div class="template" v-if="artifact.enrichments[automation.id].html" v-html="artifact.enrichments[automation.id].html"></div>
          <JSONHTML v-else :json="artifact.enrichments[automation.id].data"></JSONHTML>
        </v-card-text>
      </v-card>
    </div>

    <div v-if="artifact">
      <div v-for="enrichment in artifact.enrichments" :key="enrichment.name" >
        <v-card v-if="!hasautomation(enrichment.name)" outlined>
          <v-card-title>
            {{ enrichment.name }}
          </v-card-title>
          <v-card-text>
            <vue-markdown>
              {{ enrichment.data }}
            </vue-markdown>
          </v-card-text>
        </v-card>
      </div>
    </div>
  </v-container>
</template>

<script lang="ts">
import Vue from "vue";
import {
  Artifact,
  AutomationResponse, AutomationResponseTypeEnum, AutomationTypeEnum,
  Type,
  TypeColorEnum
} from "@/client";
import VueMarkdown from "vue-markdown";
import { API } from "@/services/api";
import JSONHTML from "../components/JSONHTML.vue";


interface State {
  artifact?: Artifact;
  automations?: Array<AutomationResponse>;
  artifacttype: string;
}

export default Vue.extend({
  name: "ArtifactPopup",
  components: {
    "vue-markdown": VueMarkdown,
    JSONHTML,
  },
  props: ["id", "name"],
  data: (): State => ({
    artifact: undefined,
    automations: undefined,
    artifacttype: "unknown",
  }),
  watch: {
    id: function(): void {
      this.load();
    },
    name: function(): void {
      this.load();
    }
  },
  computed: {
    artifactautomations: function (): Array<AutomationResponse> {
      if (!this.automations) {
        return [];
      }
      return this.lodash.filter(this.automations, (automation: AutomationResponse) => {
        if (!automation || !automation.type) {
          return true;
        }
        return this.lodash.includes(automation.type, AutomationResponseTypeEnum.Artifact)
      })
    },
    ticketID(): number {
      return parseInt(this.id, 10);
    },
    otherStates: function (): Array<Type> {
      return this.lodash.filter(this.$store.state.settings.artifactStates, (state: Type) => {
        if (!this.artifact || !this.artifact.status) {
          return true;
        }
        return state.id !== this.artifact.status;
      })
    },
  },
  methods: {
    setArtifactType() {
      if (!this.artifact || !this.artifact.name || this.ticketID === undefined) {
        return;
      }

      let artifact = this.artifact
      artifact.type = this.artifacttype
      API.setArtifact(this.ticketID, this.artifact.name, artifact).then((response) => {
        this.$store.dispatch("alertSuccess", { name: "Artifact type changed", type: "success" })
        if (response.data.artifacts) {
          this.lodash.forEach(response.data.artifacts, (artifact) => {
            if (artifact.name == this.name) {
              this.artifact = artifact;
            }
          })
        }
      });
    },
    hasautomation(name: string): boolean {
      let found = false;
      this.lodash.forEach(this.automations, (automation) => {
        if (automation.id === name) {
          found = true
        }
      })
      return found;
    },
    loadArtifact(id: number, artifact: string) {
      API.getArtifact(id, artifact).then((response) => {
        this.artifact = response.data;
      });
    },
    loadAutomations() {
      API.listAutomations().then((response) => {
        this.automations = response.data;
      });
    },
    enrich(automation: string) {
      if (this.artifact === undefined || this.ticketID === undefined) {
        return
      }
      API.runArtifact(this.ticketID, this.name, automation);
    },
    statusIcon: function (status: string): string {
      let icon = "mdi-help";
      this.lodash.forEach(this.$store.state.settings.artifactStates, (state: Type) => {
        if (status === state.id) {
          icon = state.icon;
        }
      })
      return icon;
    },
    statusColor: function (status: string) {
      let color = TypeColorEnum.Info as TypeColorEnum;
      this.lodash.forEach(this.$store.state.settings.artifactStates, (state: Type) => {
        if (status === state.id && state.color) {
          color = state.color
        }
      })
      return color;
    },
    setStatus(status: string) {
      if (!this.artifact || !this.artifact.name || this.ticketID === undefined) {
        return;
      }

      let artifact = this.artifact
      artifact.status = status
      API.setArtifact(this.ticketID, this.artifact.name, artifact).then((response) => {
        this.$store.dispatch("alertSuccess", { name: "Artifact status changed", type: "success" })
        if (response.data.artifacts) {
          this.lodash.forEach(response.data.artifacts, (artifact) => {
            if (artifact.name == this.name) {
              this.artifact = artifact;
            }
          })
        }
      });
    },
    load() {
      this.loadArtifact(this.ticketID, this.name);
    }
  },
  mounted(): void {
    this.load();
    this.loadAutomations();

    this.$store.subscribeAction((action, state) => {
      if (!action.payload || !(this.lodash.has(action.payload, "ids")) || !action.payload["ids"]) {
        return
      }
      Vue.lodash.forEach(action.payload["ids"], (id) => {
        if (id === "tickets/" + this.ticketID) {
          this.load();
        }
      });
    })
  },
});
</script>

<style>
.template {
  overflow: auto;
}
.template table, .template th, .template td {
  border: 1px solid #777 !important;
  border-left: 0 !important;
  border-right: 0 !important;
  padding: 4px !important;
}
</style>
