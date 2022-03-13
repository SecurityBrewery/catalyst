<template>
  <v-list-item link dense>
    <v-list-item-content @click="goto">
      <v-list-item-title class="d-flex">
        <v-icon small class="mr-1">mdi-gauge</v-icon>
        <span class="text-truncate">{{ artifact.name }}</span>
      </v-list-item-title>
      <v-list-item-subtitle class="d-flex">
        <v-icon small class="mr-1" :color="statusColor">{{ statusIcon }}</v-icon>
        <span :class="statusColor + '--text'">{{ artifact.status | capitalize }}</span>

        <v-icon small class="mx-1" :color="kindColor">{{ kindIcon }}</v-icon>
        <span :class="kindColor + '--text'">{{ artifact.kind | capitalize }}</span>

        <v-spacer></v-spacer>
        <v-icon small class="mr-1">mdi-information</v-icon>
        <span class="mr-1">{{ artifact.enrichments ? lodash.size(artifact.enrichments) : 0 }}</span>
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
import {Type, TypeColorEnum} from "../../client";

export default Vue.extend({
  name: "ArtifactSnippet",
  props: ["artifact", "to", "action"],
  computed: {
    statusIcon: function () {
      let icon = "mdi-help";
      this.lodash.forEach(this.$store.state.settings.artifactStates, (state: Type) => {
        if (this.artifact.status === state.id) {
          icon = state.icon;
        }
      })
      return icon;
    },
    statusColor: function () {
      let color = TypeColorEnum.Info as TypeColorEnum;
      this.lodash.forEach(this.$store.state.settings.artifactStates, (state: Type) => {
        if (this.artifact.status === state.id && state.color) {
          color = state.color;
        }
      })
      return color;
    },
    kindIcon: function () {
      let icon = "mdi-help";
      this.lodash.forEach(this.$store.state.settings.artifactKinds, (state: Type) => {
        if (this.artifact.kind === state.id) {
          icon = state.icon;
        }
      })
      return icon;
    },
    kindColor: function () {
      let color = TypeColorEnum.Info as TypeColorEnum;
      this.lodash.forEach(this.$store.state.settings.artifactKinds, (state: Type) => {
        if (this.artifact.kind === state.id && state.color) {
          color = state.color;
        }
      })
      return color;
    }
  },
  methods: {
    actionClick: function () {
      this.$emit("actionClick", this.artifact)
    },
    goto: function () {
      this.$emit("click", this.artifact)
      if (this.to) {
        this.$router.push(this.to);
      }
    }
  }
});
</script>
