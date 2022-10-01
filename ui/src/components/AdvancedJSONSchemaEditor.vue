<template>
  <div>
    <div class="d-flex" >
      <v-spacer></v-spacer>
      <v-switch
          id="advanced"
          v-model="advanced"
          label="Advanced"
          class="float-right mt-0"
      ></v-switch>
    </div>
    <div class="d-flex" >
      <v-spacer></v-spacer>
      <span v-if="advanced" class="float-right">
        See
        <a target="_blank" href="https://koumoul-dev.github.io/vuetify-jsonschema-form/latest/">
          vuetify-jsonschema documentation
        </a>
        for styling.
      </span>
    </div>

    <v-row class="flex-grow-0 flex-shrink-0">
      <v-col :cols="hidepreview ? 12 : 7">
        <v-subheader class="pl-0 py-0" style="height: 20px; font-size: 12px">
          Schema
        </v-subheader>
        <div v-if="!advanced">
          <json-schema-editor :disabled="readonly" :value="{ root: internalSchema }" lang="en_US" style="border: 1px solid #393a3f" class="mb-3 rounded" />
        </div>
        <div v-else class="flex-grow-1 flex-shrink-1 overflow-scroll">
          <Editor v-model="schemaString" lang="json" :readonly="readonly"></Editor>
        </div>
      </v-col>
      <v-col v-if="!hidepreview" cols="5">
        <v-subheader class="pl-0 py-0" style="height: 20px; font-size: 12px">
          Form output preview
        </v-subheader>
        <v-form v-model="valid">
          <v-jsf
              v-model="details"
              :schema="advanced ? parsedSchemaString : internalSchema"
              :options="{ readonly: true, formats: { time: timeformat, date: dateformat, 'date-time': datetimeformat } }"
          />
        </v-form>
      </v-col>
    </v-row>

    <v-row v-if="!readonly" class="px-3 my-6 flex-grow-0 flex-shrink-0">
      <v-btn v-if="this.$route.params.id === 'new'" color="success" @click="save" outlined>
        <v-icon>mdi-plus-thick</v-icon>
        Create
      </v-btn>
      <v-btn v-else color="success" @click="save" outlined>
        <v-icon>mdi-content-save</v-icon>
        Save
      </v-btn>
    </v-row>
  </div>
</template>

<script lang="ts">
import Editor from "./Editor.vue";
import Vue from "vue";
import {DateTime} from "luxon";

interface State {
  advanced: boolean;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  internalSchema: any;
  schemaString: string;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  details: any;
  valid: boolean;
}

export default Vue.extend({
  name: "AdvancedJSONSchemaEditor",
  components: { Editor },
  props: {
    schema: {
      type: Object,
          required: true
    },
    hidepreview: {
      type: Boolean,
      default: true
    },
    readonly: {
      type: Boolean,
      default: false
    },
  },
  data: (): State => ({
    details: {},
    advanced: false,
    internalSchema: {},
    schemaString: "{}",
    valid: true,
  }),
  watch: {
    schema: function () {
      this.internalSchema = this.schema;
      this.schemaString = JSON.stringify(this.internalSchema, null, 2);
    },
    advanced: function (advanced) {
      if (advanced) {
        this.schemaString = JSON.stringify(this.internalSchema, null, 2);
      } else {
        this.internalSchema = JSON.parse(this.schemaString);
      }
    }
  },
  computed: {
    parsedSchemaString: function() {
      try {
        return JSON.parse(this.schemaString);
      }
      catch (e) {
        return {}
      }
    },
  },
  methods: {
    save: function () {
      let schema = this.schemaString;
      if (!this.advanced) {
        schema = JSON.stringify(this.internalSchema);
      }
      this.$emit("save", schema);
    },
    timeformat: function(s: string, locale: string) {
      let format = this.$store.state.settings.timeformat;
      if (!format) {
        return DateTime.fromISO(s).toLocaleString(DateTime.DATETIME_SHORT);
      }
      return DateTime.fromISO(s).toFormat(format);
    },
    dateformat: function(s: string, locale: string) {
      let format = this.$store.state.settings.timeformat;
      if (!format) {
        return DateTime.fromISO(s).toLocaleString(DateTime.DATETIME_SHORT);
      }
      return DateTime.fromISO(s).toFormat(format);
    },
    datetimeformat: function(s: string, locale: string) {
      let format = this.$store.state.settings.timeformat;
      if (!format) {
        return DateTime.fromISO(s).toLocaleString(DateTime.DATETIME_SHORT);
      }
      return DateTime.fromISO(s).toFormat(format);
    },
    hasRole: function (s: string): boolean {
      if (this.$store.state.settings.roles) {
        return this.lodash.includes(this.$store.state.settings.roles, s);
      }
      return false;
    }
  },
  mounted() {
    this.internalSchema = this.schema;
    this.schemaString = JSON.stringify(this.internalSchema, null, 2);
  },
});
</script>

<style>
.theme--dark .ant-btn,
.theme--dark .ant-select-selection,
.theme--dark .ant-input,
.theme--dark .ant-input-number,
.theme--dark .ant-modal-header,
.theme--dark .ant-modal-title,
.theme--dark .ant-form-item,
.theme--dark .ant-modal-close-x,
.theme--dark .ant-select-dropdown,
.theme--dark .ant-checkbox-inner {
  color: white !important;
  background: none !important;
}

.theme--dark .ant-select-selection,
.theme--dark .ant-input,
.theme--dark .ant-input-number,
.theme--dark .ant-modal-header,
.theme--dark .ant-modal-footer,
.theme--dark .ant-checkbox-inner {
  border-color: #424242 !important;
}

.theme--dark .ant-modal-content,
.theme--dark .ant-select-dropdown {
  color: white !important;
  background: #303030 !important;
}
</style>
