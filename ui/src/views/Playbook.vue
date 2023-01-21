<template>
  <div v-if="playbook !== undefined" class="fill-height d-flex flex-column pa-8">
    <v-row class="flex-grow-0 flex-shrink-0">
      <v-spacer></v-spacer>
      <v-btn href="https://catalyst-soar.com/docs/catalyst/engineer/playbook" target="_blank" outlined rounded small>
        <v-icon>mdi-book-open</v-icon> Handbook
      </v-btn>
    </v-row>

    <v-alert v-if="readonly" type="info">You do not have write access to playbooks.</v-alert>
    <h2 v-if="this.$route.params.id === 'new'">New Playbook</h2>
    <h2 v-else>Edit Playbook: {{ playbook.name }}</h2>

    <v-alert v-if="formaterrors.length" color="warning">
      <div v-for="(formaterror, index) in formaterrors" :key="index">
        {{ formaterror.instancePath }}{{ formaterror.instancePath ? ": " : "" }}{{ formaterror.message }}
      </div>
    </v-alert>
    <v-alert v-else-if="error" color="warning">
      {{ error }}
    </v-alert>

    <div class="flex-grow-1 flex-shrink-1 overflow-scroll">
      <v-tabs
          v-model="tab"
          background-color="transparent"
      >
        <v-tab>Graph (experimental)</v-tab>
        <v-tab>YAML</v-tab>
      </v-tabs>
      <v-tabs-items v-model="tab" style="background: transparent">
        <v-tab-item>
          <v-text-field
              v-model="playbook.name"
              label="Name"
              outlined
              dense
              :readonly="readonly"
              class="mt-4"
          />
          <PlaybookEditor
              v-if="playbookJSON"
              v-model="playbookJSON" />
        </v-tab-item>
        <v-tab-item>
          <v-card class="py-2">
            <Editor v-model="playbookYAML" lang="yaml" :readonly="readonly"></Editor>
          </v-card>
        </v-tab-item>
      </v-tabs-items>
    </div>

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
import Vue from "vue";

import {Playbook, PlaybookTemplate, Task, TaskResponse} from "../client";
import { API } from "@/services/api";
import Editor from "../components/Editor.vue";
import {alg, Graph} from "graphlib";
import yaml from 'yaml';
import Ajv from "ajv";
import PlaybookEditor from "@/components/playbookeditor/PlaybookEditor.vue";

const playbookSchema = {
  type: "object",
  required: ["name", "tasks"],
  properties: {
    name: { type: "string" },
    tasks: {
      type: "object",
      additionalProperties: { $ref: "#/definitions/Task" }
    }
  },
  // additionalProperties: false,
  $id: "#/definitions/Playbook"
};

const taskSchema = {
  type: "object",
  required: ["name", "type"],
  properties: {
    automation: { type: "string" },
    join: { type: "boolean" },
    payload: { type: "object", additionalProperties: { type: "string" } },
    name: { type: "string" },
    next: {
      type: "object",
      additionalProperties: { type: ["string", "null"] }
    },
    schema: { type: "object" },
    type: { type: "string", enum: ["task", "input", "automation"] }
  },
  // additionalProperties: false,
  $id: "#/definitions/Task"
};

interface State {
  playbook?: PlaybookTemplate;
  g: Record<string, any>;
  selected: any;
  error: string;
  tab: number;
  playbookYAML: string;
  playbookJSON: any;
}

interface TaskWithID {
  id: string;
  task: Task;
}

const inityaml = "name: VirusTotal hash check\n" +
    "tasks:\n" +
    "  input:\n" +
    "    name: Please enter a word\n" +
    "    type: input\n" +
    "    schema:\n" +
    "      title: Word\n" +
    "      type: object\n" +
    "      properties:\n" +
    "        word:\n" +
    "          type: string\n" +
    "          title: Enter a Word\n" +
    "          default: \"\"\n" +
    "    next:\n" +
    "      hash: \"word != ''\"\n" +
    "\n" +
    "  hash:\n" +
    "    name: Hash the word\n" +
    "    type: automation\n" +
    "    automation: hash.sha1\n" +
    "    payload:\n" +
    "      default: \"playbook.tasks['input'].data['word']\"\n" +
    "    next:\n" +
    "      end:\n" +
    "\n" +
    "  end:\n" +
    "    name: Finish the incident\n" +
    "    type: task\n"

export default Vue.extend({
  name: "Playbook",
  components: { Editor, PlaybookEditor },
  data: (): State => ({
    playbook: undefined,
    g: {},
    selected: undefined,
    error: "",
    tab: 1,
    playbookJSON: undefined,
    playbookYAML: inityaml
  }),
  watch: {
    '$route': function () {
      this.loadPlaybook();
    },
    tab: function (value) {
      if (value === 0) {
        this.playbookJSON = yaml.parse(this.playbookYAML);
      } else {
        this.playbookYAML = yaml.stringify(this.playbookJSON);
      }
    }
  },
  computed: {
    formaterrors: function (): Array<any>  {
      if (!this.playbook) {
        return [];
      }
      try {
        let playbook = yaml.parse(this.playbook.yaml);
        const ajv = new Ajv({validateFormats: false});
        ajv.addSchema(taskSchema, "#/definitions/Task")
        const validate = ajv.compile(playbookSchema);
        const valid = validate(playbook)
        if (!valid && validate.errors) {
          return validate.errors;
        }
        return [];
      }
      catch (e) {
        return [e];
      }
    },
    readonly: function (): boolean {
      return !this.hasRole("playbook:write");
    },
  },
  methods: {
    gstatus: function(task: TaskResponse) {
      if (task.active) {
        return "open"
      }
      return "inactive"
    },
    tasks: function(g: any, playbook: Playbook): Array<TaskWithID> {
      let taskKeys = alg.topsort(g);
      let tasks = [] as Array<TaskWithID>;
      for (const tasksKey in taskKeys) {
        let taskWithID = {} as TaskWithID;
        if (playbook.tasks[taskKeys[tasksKey]] === undefined) {
          continue; // TODO
        }
        taskWithID.task = playbook.tasks[taskKeys[tasksKey]];
        taskWithID.id = taskKeys[tasksKey];
        tasks.push(taskWithID);
      }
      return tasks;
    },
    save() {
      if (this.playbook === undefined) {
        return;
      }
      let playbook = this.playbook;
      if (this.tab === 0) {
        let jsonData = this.playbookJSON;
        jsonData["name"] = playbook.name;
        playbook.yaml = yaml.stringify(jsonData);
      } else {
        playbook.yaml = this.playbookYAML;
      }

      if (this.$route.params.id == 'new') {
        // playbook.id = kebabCase(playbook.name);
        API.createPlaybook(playbook).then(() => {
          this.$store.dispatch("alertSuccess", { name: "Playbook created" });
        });
      } else {
        API.updatePlaybook(this.$route.params.id, playbook).then(() => {
          this.$store.dispatch("alertSuccess", { name: "Playbook saved" });
        });
      }
    },
    loadPlaybook() {
      if (!this.$route.params.id) {
        return
      }
      if (this.$route.params.id == 'new') {
        this.playbook = { name: "MyPlaybook", yaml: inityaml }
        this.playbookJSON = yaml.parse(this.playbook.yaml);
        this.playbookYAML = this.playbook.yaml;
      } else {
        API.getPlaybook(this.$route.params.id).then((response) => {
          this.playbook = response.data;
          this.playbookJSON = yaml.parse(this.playbook.yaml);
          this.playbookYAML = this.playbook.yaml;
        });
      }
    },
    hasRole: function (s: string): boolean {
      if (this.$store.state.settings.roles) {
        return this.lodash.includes(this.$store.state.settings.roles, s);
      }
      return false;
    },
  },
  mounted() {
    this.loadPlaybook();
  },
});
</script>

<style>
.my-code {
  background: #2d2d2d;
  color: #ccc;

  width: inherit;
  height: inherit;

  font-family: Fira code, Fira Mono, Consolas, Menlo, Courier, monospace;
  font-size: 14px;
  line-height: 1.5;
  padding: 15px;
}

.pipeline-node-label {
  fill: #333 !important;
}
</style>
