<template>
  <div v-if="dashboard">
    <h2 class="d-flex">
      <span v-if="!editmode">{{ dashboard.name }}</span>
      <v-text-field v-else v-model="dashboard.name" outlined dense class="mb-0" hide-details></v-text-field>

      <v-spacer></v-spacer>
      <v-btn v-if="editmode" small outlined @click="addWidget" class="mr-1">
        <v-icon>mdi-plus</v-icon>
        Add Widget
      </v-btn>
      <v-btn v-if="editmode" small outlined @click="save" class="mr-1">
        <v-icon>mdi-content-save</v-icon>
        Save
      </v-btn>
      <v-btn v-if="editmode && $route.params.id !== 'new'" small outlined @click="cancel">
        <v-icon>mdi-cancel</v-icon>
        Cancel
      </v-btn>
      <v-btn v-if="!editmode" small outlined @click="edit">
        <v-icon>mdi-pencil</v-icon>
        Edit
      </v-btn>
    </h2>

    <v-row>
      <v-col v-for="(widget, index) in dashboard.widgets" :key="index" :cols="widget.width">
        <v-card class="mb-2">
          <v-card-title>
            <span v-if="!editmode">{{ widget.name }}</span>
            <v-text-field v-else outlined dense hide-details v-model="widget.name" class="mr-1"></v-text-field>
            <v-btn v-if="editmode" outlined @click="removeWidget(index)">
              <v-icon>mdi-close</v-icon>
              Remove
            </v-btn>
          </v-card-title>

          <v-card-text v-if="editmode">
            <v-row>
              <v-col cols="8">
                <v-select label="Type" v-model="widget.type" :items="['line', 'bar', 'pie']"></v-select>
              </v-col>
              <v-col cols="4">
                <v-text-field label="Width" type="number" v-model="widget.width"></v-text-field>
              </v-col>
            </v-row>
            <v-text-field label="Aggregation" v-model="widget.aggregation"></v-text-field>
            <v-text-field label="Filter" v-model="widget.filter" clearable></v-text-field>

          </v-card-text>

          <v-card-text v-if="data[index] === null">
            {{ widgetErrors[index] }}
          </v-card-text>
          <div v-else>
            <line-chart
                v-if="widget.type === 'line' && data[index]"
                :chart-data="data[index]"
                :styles="{ width: '100%', position: 'relative' }"
                :chart-options="{
                    responsive: true,
                    maintainAspectRatio: false,
                    legend: false,
                    scales: { yAxes: [ { ticks: { beginAtZero: true, precision: 0 } } ] }
                  }"
            >
            </line-chart>

            <pie-chart
                v-if="widget.type === 'pie' && data[index]"
                :chart-data="data[index]"
                :styles="{ width: '100%', position: 'relative' }"
                :chart-options="{
                    responsive: true,
                    maintainAspectRatio: false,
                  }"
            >
            </pie-chart>

            <bar-chart
                v-if="widget.type === 'bar' && data[index]"
                :chart-data="data[index]"
                :styles="{
                    width: '100%',
                    'max-height': '400px',
                    position: 'relative'
                  }"
                :chart-options="{
                    responsive: true,
                    maintainAspectRatio: false,
                    legend: false,
                    scales: { xAxes: [ { ticks: { beginAtZero: true, precision: 0 } } ] },
                  }"
            ></bar-chart>
          </div>
        </v-card>
      </v-col>
    </v-row>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import {DashboardResponse, Widget} from "@/client";
import { API } from "@/services/api";
import {createHash} from "crypto";
import {colors} from "@/plugins/vuetify";
import LineChart from "../components/charts/Line";
import BarChart from "../components/charts/Bar";
import PieChart from "../components/charts/Doughnut";
import {ChartData} from "chart.js";
import {AxiosError, AxiosTransformer} from "axios";

interface State {
  dashboard?: DashboardResponse;
  undodashboard?: DashboardResponse;
  data: Record<string, any>;
  editmode: boolean;
  widgetErrors: Record<number, string>;
}

export default Vue.extend({
  name: "Dashboard",
  components: {
    LineChart,
    BarChart,
    PieChart,
  },
  data: (): State => ({
    dashboard: undefined,
    undodashboard: undefined,
    data: {},
    editmode: false,
    widgetErrors: {},
  }),
  watch: {
    $route: function () {
      this.loadDashboard();
    },
  },
  methods: {
    edit: function () {
      this.undodashboard = this.lodash.cloneDeep(this.dashboard);
      this.editmode = true;
    },
    save: function () {
      if (!this.dashboard) {
        return
      }

      let widgets = [] as Array<Widget>;
      this.lodash.forEach(this.dashboard.widgets, (widget) => {
        widget.width = this.lodash.toInteger(widget.width);
        if (!widget.filter) {
          this.lodash.unset(widget, "filter")
        }
        widgets.push(widget);
      })
      this.dashboard.widgets = widgets;

      if (this.$route.params.id === 'new') {
        API.createDashboard(this.dashboard).then((response) => {
          this.loadWidgetData(response.data.widgets);

          this.dashboard = response.data;
          this.editmode = false;

          this.$router.push({ name: "Dashboard", params: { id: response.data.id }})
        })
      } else {
        API.updateDashboard(this.dashboard.id, this.dashboard).then((response) => {
          this.loadWidgetData(response.data.widgets);

          this.dashboard = response.data;
          this.editmode = false;
        })
      }
    },
    cancel: function () {
      this.dashboard = this.lodash.cloneDeep(this.undodashboard);
      this.editmode = false;
    },
    addWidget: function () {
      if (!this.dashboard) {
        return
      }

      this.dashboard.widgets.push({name: "new widget", width: 6, aggregation: "", type: "line"})
    },
    removeWidget: function (id: number) {
      if (!this.dashboard) {
        return
      }

      console.log(id);
      let widgets = this.lodash.cloneDeep(this.dashboard.widgets);
      this.lodash.pullAt(widgets, [id]);
      Vue.set(this.dashboard, "widgets", widgets);
    },
    loadDashboard: function () {
      if (this.$route.params.id === 'new') {
        this.dashboard = {
          name: "New dashboard",
          widgets: [{name: "new widget", width: 6, aggregation: "", type: "line"}],
        } as DashboardResponse
        this.editmode = true;
      } else {
        API.getDashboard(this.$route.params.id).then((response) => {
          this.loadWidgetData(response.data.widgets);

          this.dashboard = response.data;
        });
      }
    },
    loadWidgetData: function (widgets: Array<Widget>) {
      this.lodash.forEach(widgets, (widget: Widget, index: number) => {
        let widgetErrors = {};
        let defaultTransformers = this.axios.defaults.transformResponse as AxiosTransformer[]
        let transformResponse = defaultTransformers.concat((data) => {
          data.notoast = true;
          return data
        });
        API.dashboardData(widget.aggregation, widget.filter, {transformResponse: transformResponse}).then((response) => {
          let d = { labels: [], datasets: [{data: [], backgroundColor: []}] } as ChartData;
          this.lodash.forEach(response.data, (v: any, k: string) => {
            // @ts-expect-error T2532
            d.labels.push(k)
            // @ts-expect-error T2532
            d.datasets[0].data.push(v)

            if (widget.type !== 'line') {
              // @ts-expect-error T2532
              d.datasets[0].backgroundColor.push(this.color(this.lodash.toString(v)));
            }
          })

          Vue.set(this.data, index, d);
        }).catch((err: AxiosError) => {
          widgetErrors[index] = this.lodash.toString(err.response?.data.error);
          Vue.set(this.data, index, null);
        })
        Vue.set(this, 'widgetErrors', widgetErrors);
      })
    },
    color: function (s: string): string {
      let pos = createHash('md5').update(s).digest().readUInt32BE(0) % colors.length;
      return colors[pos];
    },
  },
  mounted() {
    this.loadDashboard();
  }
});
</script>
