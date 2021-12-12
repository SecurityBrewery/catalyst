<template>
  <div class="fill-height">
    <v-card
        v-if="selected !== undefined"
        class="mt-3 ml-3 px-0"
        style="position: absolute; width: 33%; left: 60px; top: 60px; z-index: 5">
      <v-card-title>
        {{ selected.name }}
        <v-spacer></v-spacer>
        <v-btn
            text
            :to="{ name: 'Graph', params: { col: col(selected.id), 'id': id(selected.id) } }"
        >
          <v-icon>mdi-bullseye</v-icon>
        </v-btn>
        <v-btn
            vv-if="selected.id.startsWith('artifacts/')"
            text
            @click="to(selected.id)"
        >
          <v-icon>mdi-open-in-new</v-icon>
        </v-btn>
      </v-card-title>
      <v-card-text class="ma-0 pa-0">
        <!--TicketSnippet v-if="selected.id.startsWith('tickets/')"></TicketSnippet-->
        <!--ArtifactSnippet v-if="selected.id.startsWith('artifacts/')"></ArtifactSnippet-->
        <IDSnippet :id="selected.id"></IDSnippet>
      </v-card-text>
    </v-card>
    <v-main class="fill-height overflow-hidden">
      <div class="d-flex flex-column fill-height overflow-hidden">
        <v-row class="flex-grow-0 pb-4">
          <v-col cols="4" class="pl-8">
            <v-slider
                class="mt-4 mb-4"
                v-model="depth"
                dense
                :label="'Depth ' + depth"
                hide-details
                max="10"
                min="0"
            ></v-slider>
          </v-col>
          <v-col cols="4" class="pl-4">
          <v-slider
              class="mt-4 mb-4"
              v-model="force"
              dense
              label="Node Distance"
              hide-details
              max="10000"
              min="0"
          ></v-slider>
          </v-col>
          <v-col cols="2">
            <v-switch label="Hide Labels"
              v-model="hidelabel" class="mt-6 mb-4" hide-details>
            </v-switch>
          </v-col>
        </v-row>
        <v-row class="flex-grow-1 mt-0">
          <d3-network
              v-if="g !== undefined && g.links !== undefined"
              :net-nodes="nodes"
              :net-links="g.links"
              :options="options"
              @node-click="nodeclick"
              class="pt-0 mt-n4"
              style="width: 100%; height: 100%"
          />
        </v-row>
      </div>
    </v-main>
  </div>
</template>

<script lang="ts">
import Vue from "vue";

import D3Network from "vue-d3-network";
import {API} from "../services/api";
import {Graph, Node} from "../client";
// import TicketSnippet from "../components/snippets/TicketSnippet.vue";
// import ArtifactSnippet from "../components/snippets/ArtifactSnippet.vue";
import IDSnippet from "../components/snippets/IDSnippet.vue";

interface State {
  g?: Graph;
  force: number;
  depth: number;
  hidelabel: boolean;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  selected?: any;
}

export default Vue.extend({
  name: "Graph",
  components: {
    IDSnippet,
    D3Network
  },
  data: (): State => ({
    g: undefined,
    force: 3000,
    depth: 2,
    hidelabel: false,
    selected: undefined,
  }),
  computed: {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    options: function (): any {
      return {
        canvas: false,
        nodeLabels: !this.hidelabel,
        nodeSize: 20,
        linkWidth: 2,
        fontSize: 16,
        force: this.force
      }
    },
    nodes: function (): Array<any> {
      if (this.g === undefined || this.g.nodes === undefined) {
        return []
      }
      return this.lodash.map(this.g.nodes, (node: Node) => {
        if (node.id === this.$route.params.col + "/" + this.$route.params.id) {
          return {id: node.id, name: node.name, _size: 40, _cssClass: "center", _labelClass: "center"}
        }
        if (node.id.startsWith("tickets/")) {
          return {id: node.id, name: node.name, _size: 30, _cssClass: "ticket", _labelClass: "event"}
        }
        return {id: node.id, name: node.name}
      })
    }
  },
  watch: {
    depth: function () {
      this.fetchGraph();
    },
    $route: function () {
      this.fetchGraph();
    },
  },
  methods: {
    col: function (id: string): string {
      let parts = id.split("/");
      return parts[0];
    },
    id: function (id: string): string {
      let parts = id.split("/");
      return parts[1];
    },
    to: function (fid: string) {
      let col = this.col(fid);
      let id = this.id(fid);

      if (col === 'tickets') {
        this.$router.push({
          name: "Ticket",
          params: { id: id, type: "-" }
        });
        return
      }

      this.$router.push({
        name: "Artifact",
        params: { artifact: id }
      });
    },
    fetchGraph: function(): void {
      API.graph(this.$route.params.col, this.$route.params.id, this.depth).then((response) => {
        this.g = response.data;
      });
    },
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    nodeclick(e: any, node: any) {
      this.selected = node;
    }
  },
  mounted() {
    this.fetchGraph();
  },
});
</script>

<style>
.node,
.node.selected {
  stroke: #388E3C !important;
}
.node.event {
  stroke: #D32F2F !important;
}
.node.center {
  stroke: #FFEB3B !important;
  fill: #FFEB3B !important;
}

.theme--dark .node-label,
.theme--dark .node-label.event {
  fill: #ffffff !important;
}

.node-label,
.node-label.event {
  fill: #000000 !important;
}

.link {
  stroke: #424242 !important;
}
.link.selected,
.link:hover,.node:hover{
  stroke: #FFEB3B !important;
}
</style>

