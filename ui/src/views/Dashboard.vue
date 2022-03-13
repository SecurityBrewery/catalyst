<template>
  <div v-if="dashboard">
    <h2>{{ dashboard.name }}</h2>

    <v-row>
      <v-col v-for="widget in dashboard.widgets" :key="widget.name" :cols="widget.width">
        <v-card class="mb-2">
          <v-card-title>{{ widget.name }}</v-card-title>
          <line-chart
              v-if="widget.type === 'line' && data[widget.name]"
              :chart-data="data[widget.name]"
              :styles="{ width: '100%', position: 'relative' }"
              :chart-options="{
                    responsive: true,
                    maintainAspectRatio: false,
                    legend: undefined,
                    scales: { yAxes: [ { ticks: { beginAtZero: true, precision: 0 } } ] }
                  }"
          >
          </line-chart>

          <pie-chart
              v-if="widget.type === 'pie' && data[widget.name]"
              :chart-data="data[widget.name]"
              :styles="{ width: '100%', position: 'relative' }"
          >
          </pie-chart>

          <bar-chart
              v-if="widget.type === 'bar' && data[widget.name]"
              :chart-data="data[widget.name]"
              :styles="{
                    width: '100%',
                    'max-height': '400px',
                    position: 'relative'
                  }"
              :chart-options="{
                    responsive: true,
                    maintainAspectRatio: false,
                    legend: undefined,
                    scales: { xAxes: [ { ticks: { beginAtZero: true, precision: 0 } } ] },
                  }"
          ></bar-chart>
        </v-card>
      </v-col>
    </v-row>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import { Dashboard, Widget } from "@/client";
import { API } from "@/services/api";
import {createHash} from "crypto";
import {colors} from "@/plugins/vuetify";
import LineChart from "../components/charts/Line";
import BarChart from "../components/charts/Bar";
import PieChart from "../components/charts/Doughnut";
import {ChartData} from "chart.js";

interface State {
  dashboard?: Dashboard;
  data: Record<string, any>;
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
    data: {},
  }),
  watch: {
    $route: function () {
      this.loadDashboard();
    },
  },
  methods: {
    saveUserData: function(dashboard: Dashboard) {
      API.updateUserData(this.$route.params.id, dashboard).then(() => {
        this.$store.dispatch("alertSuccess", { name: "Dashboard saved" });
      });
    },
    loadDashboard: function () {
      API.getDashboard(this.$route.params.id).then((response) => {
        this.lodash.forEach(response.data.widgets, (widget: Widget) => {
          API.dashboardData(widget.aggregation, widget.filter).then((response) => {
            let d = { labels: [], datasets: [{data: [], backgroundColor: []}] } as ChartData;
            this.lodash.forEach(response.data, (v: any, k: string) => {
              // @ts-expect-error T2532
              d.labels.push(k)
              // @ts-expect-error T2532
              d.datasets[0].data.push(v)
              // @ts-expect-error T2532
              d.datasets[0].backgroundColor.push(v);
            })

            Vue.set(this.data, widget.name, d);
            // this.data[widget.name] = d;
          })
        })

        this.dashboard = response.data;
      });
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
