<template>
  <v-main class="ma-4">
    <div v-if="settings !== undefined">
      <v-text-field label="Time Format" v-model="settings.timeformat"></v-text-field>

      <v-subheader class="mx-0 px-0">Artifact States</v-subheader>
      <v-card v-for="state in settings.artifactStates" :key="state.id" class="d-flex mb-2">
        <v-row class="px-4 pt-2" dense>
          <v-col>
            <v-text-field label="ID" v-model="state.id"></v-text-field>
          </v-col>
          <v-col>
            <v-text-field label="Name" v-model="state.name"></v-text-field>
          </v-col>
          <v-col>
            <v-text-field label="Icon" v-model="state.icon"></v-text-field>
          </v-col>
          <v-col>
            <v-select label="Color" v-model="state.color" :items="['info', 'error', 'success', 'warning']" clearable></v-select>
          </v-col>
        </v-row>
        <v-btn icon class="mt-6 mr-4" @click="removeState(state.id)"><v-icon>mdi-close</v-icon></v-btn>
      </v-card>
      <v-card class="d-flex mb-2">
        <v-row class="px-4 pt-2" dense>
          <v-col>
            <v-text-field label="ID" v-model="newState.id"></v-text-field>
          </v-col>
          <v-col>
            <v-text-field label="Name" v-model="newState.name"></v-text-field>
          </v-col>
          <v-col>
            <v-text-field label="Icon" v-model="newState.icon"></v-text-field>
          </v-col>
          <v-col>
            <v-select label="Color" v-model="newState.color" :items="['info', 'error', 'success', 'warning']" clearable></v-select>
          </v-col>
        </v-row>
        <v-btn icon class="mt-6 mr-4" @click="addState"><v-icon>mdi-plus</v-icon></v-btn>
      </v-card>

      <v-subheader class="mx-0 px-0">Artifact Types</v-subheader>
      <v-card v-for="state in settings.artifactKinds" :key="state.id" class="d-flex mb-2">
        <v-row class="px-4 pt-2" dense>
          <v-col>
            <v-text-field label="ID" v-model="state.id"></v-text-field>
          </v-col>
          <v-col>
            <v-text-field label="Name" v-model="state.name"></v-text-field>
          </v-col>
          <v-col>
            <v-text-field label="Icon" v-model="state.icon"></v-text-field>
          </v-col>
          <v-col>
            <v-select label="Color" v-model="state.color" :items="['info', 'error', 'success', 'warning']" clearable></v-select>
          </v-col>
        </v-row>
        <v-btn icon class="mt-6 mr-4" @click="removeKind(state.id)"><v-icon>mdi-close</v-icon></v-btn>
      </v-card>
      <v-card class="d-flex mb-2">
        <v-row class="px-4 pt-2" dense>
          <v-col>
            <v-text-field label="ID" v-model="newKind.id"></v-text-field>
          </v-col>
          <v-col>
            <v-text-field label="Name" v-model="newKind.name"></v-text-field>
          </v-col>
          <v-col>
            <v-text-field label="Icon" v-model="newKind.icon"></v-text-field>
          </v-col>
          <v-col>
            <v-select label="Color" v-model="newKind.color" :items="['info', 'error', 'success', 'warning']" clearable></v-select>
          </v-col>
        </v-row>
        <v-btn icon class="mt-6 mr-4" @click="addKind"><v-icon>mdi-plus</v-icon></v-btn>
      </v-card>

      <v-btn color="success" @click="save" outlined class="mt-2">
        <v-icon>mdi-content-save</v-icon>
        Save
      </v-btn>
    </div>
  </v-main>
</template>

<script lang="ts">
import {DateTime} from "luxon";
import Vue from "vue";
import {Settings, SettingsResponse, Type} from "@/client";
import {API} from "@/services/api";
import {AxiosResponse} from "axios";

interface State {
  valid: boolean;
  settings?: Settings;
  newState: Type;
  newKind: Type;
}

export default Vue.extend({
  name: "Settings",
  data: (): State => ({
    valid: true,
    settings: undefined,
    newState: {} as Type,
    newKind: {} as Type,
  }),
  methods: {
    save: function () {
      if (this.settings === undefined) {
        return
      }

      API.saveSettings(this.settings).then((response) => {
        this.settings = response.data;
        this.$store.dispatch("getSettings");
      })
    },
    addState: function () {
      if (this.settings === undefined) {
        return
      }

      this.settings.artifactStates.push(this.newState);
      this.newState = {} as Type;
    },
    removeState: function (id: string) {
      if (this.settings === undefined) {
        return
      }

      this.settings.artifactStates = this.lodash.filter(this.settings.artifactStates, function (t) { return t.id !== id });
    },
    addKind: function () {
      if (this.settings === undefined) {
        return
      }

      this.settings.artifactKinds.push(this.newKind);
      this.newKind = {} as Type;
    },
    removeKind: function (id: string) {
      if (this.settings === undefined) {
        return
      }

      this.settings.artifactKinds = this.lodash.filter(this.settings.artifactKinds, function (t) { return t.id !== id });
    },
    timeformat: function (s: string) {
      let format = this.$store.state.settings.timeformat;
      if (!format) {
        return DateTime.fromISO(s).toLocaleString(DateTime.DATETIME_SHORT);
      }
      return DateTime.fromISO(s).toFormat(format);
    },
    dateformat: function (s: string) {
      let format = this.$store.state.settings.timeformat;
      if (!format) {
        return DateTime.fromISO(s).toLocaleString(DateTime.DATETIME_SHORT);
      }
      return DateTime.fromISO(s).toFormat(format);
    },
    datetimeformat: function (s: string) {
      let format = this.$store.state.settings.timeformat;
      if (!format) {
        return DateTime.fromISO(s).toLocaleString(DateTime.DATETIME_SHORT);
      }
      return DateTime.fromISO(s).toFormat(format);
    },
  },
  mounted() {
    API.getSettings().then((response: AxiosResponse<SettingsResponse>) => {
      this.settings = response.data;
    })
  }
});
</script>
