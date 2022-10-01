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
    <div v-else class="px-4 overflow-scroll">
        <vue-pipeline
            v-if="pipelineData"
            ref="pipeline"
            :x="50"
            :y="55"
            :data="pipelineData"
            :showArrow="true"
            :ystep="70"
            :xstep="100"
            lineStyle="default"
        />
    </div>

    <v-subheader class="pl-0 py-0" style="height: 20px; font-size: 12px">
      Playbook
    </v-subheader>
    <div class="flex-grow-1 flex-shrink-1 overflow-scroll">
      <Editor v-model="playbook.yaml" @input="updatePipeline" lang="yaml" :readonly="readonly"></Editor>
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
    msg: { type: "object", additionalProperties: { type: "string" } },
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
  pipelineData: any;
  error: string;
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
    "    msg:\n" +
    "      payload: \"playbook.tasks['input'].data['word']\"\n" +
    "    next:\n" +
    "      end:\n" +
    "\n" +
    "  end:\n" +
    "    name: Finish the incident\n" +
    "    type: task\n"

export default Vue.extend({
  name: "Playbook",
  components: { Editor },
  data: (): State => ({
    playbook: undefined,
    g: {},
    selected: undefined,
    pipelineData: undefined,
    error: "",
  }),
  watch: {
    '$route': function () {
      this.loadPlaybook();
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
    updatePipeline: function () {
      if (this.playbook) {
        this.pipeline(this.playbook.yaml);
      }
    },
    pipeline: function(playbookYAML: string) {
      try {
        let playbook = yaml.parse(playbookYAML);

        this.error = "";

        let g = new Graph();

        for (const stepKey in playbook.tasks) {
          g.setNode(stepKey);
        }

        this.lodash.forEach(playbook.tasks, (task: Task, stepKey: string) => {
          if ("next" in task) {
            this.lodash.forEach(task.next, (condition, nextKey) => {
              g.setEdge(stepKey, nextKey);
            });
          }
        });

        let tasks = this.tasks(g, playbook);
        let elements = [] as Array<any>;
        this.lodash.forEach(tasks, task => {
          elements.push({
            id: task.id,
            name: task.task.name,
            next: [],
            status: "unknown"
          });
        });

        this.lodash.forEach(tasks, (task: TaskWithID) => {
          if ("next" in task.task) {
            this.lodash.forEach(task.task.next, (condition, nextKey) => {
              let nextID = this.lodash.findIndex(elements, ["id", nextKey]);
              let stepID = this.lodash.findIndex(elements, ["id", task.id]);
              if (nextID !== -1) {
                // TODO: invalid schema
                elements[stepID].next.push({index: nextID});
              }
            });
          }
        });

        this.pipelineData = undefined;
        this.$nextTick(() => {
          this.pipelineData = this.lodash.values(elements);
        })
      }
      catch (e: unknown) {
        console.log(e);
        this.error = this.lodash.toString(e);
      }
    },
    save() {
      if (this.playbook === undefined) {
        return;
      }
      if (this.$route.params.id == 'new') {
        let playbook = this.playbook;
        // playbook.id = kebabCase(playbook.name);
        API.createPlaybook(playbook).then(() => {
          this.$store.dispatch("alertSuccess", { name: "Playbook created" });
        });
      } else {
        API.updatePlaybook(this.$route.params.id, this.playbook).then(() => {
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
      } else {
        API.getPlaybook(this.$route.params.id).then((response) => {
          this.playbook = response.data;
        });
      }
    },
    hasRole: function (s: string): boolean {
      if (this.$store.state.settings.roles) {
        return this.lodash.includes(this.$store.state.settings.roles, s);
      }
      return false;
    }
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
