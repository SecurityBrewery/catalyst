<template>
  <v-main>
    <v-row>
      <v-col v-if="statistics" cols="12" lg="7">
        <v-row>
          <v-col cols="4">
            <v-subheader>Unassigned tickets</v-subheader>
            <span style="font-size: 60pt; text-align: center; display: block">
              <router-link :to="{
                name: 'TicketList',
                params: { query: 'status == \'open\' AND !owner' }
              }">
                {{ statistics.unassigned }}
              </router-link>
            </span>
            <v-subheader>Your tickets</v-subheader>
            <span style="font-size: 60pt; text-align: center; display: block">
              <router-link :to="{
                name: 'TicketList',
                params: { query: 'status == \'open\' AND owner == \'' + $store.state.user.id + '\'' }
              }">
                {{ $store.state.user.id in statistics.open_tickets_per_user ? statistics.open_tickets_per_user[$store.state.user.id] : 0 }}
              </router-link>
            </span>
          </v-col>
          <v-col cols="8">
            <v-subheader>Open tickets per owner</v-subheader>
            <bar-chart
              v-if="open_tickets_per_user"
              :chart-data="open_tickets_per_user"
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
                onClick: clickUser,
                hover: {
                  onHover: function(e) {
                     var point = this.getElementAtEvent(e);
                     if (point.length) e.target.style.cursor = 'pointer';
                     else e.target.style.cursor = 'default';
                  }
               }
              }"
            ></bar-chart>
          </v-col>
        </v-row>
        <v-row>
          <v-col cols="7">
            <v-subheader>Tickets created per week</v-subheader>
            <line-chart
              v-if="tickets_per_week"
              :chart-data="tickets_per_week"
              :styles="{ width: '100%', position: 'relative' }"
              :chart-options="{
                responsive: true,
                maintainAspectRatio: false,
                legend: undefined,
                scales: { yAxes: [ { ticks: { beginAtZero: true, precision: 0 } } ] }
              }"
            >
            </line-chart>
          </v-col>
          <v-col cols="5">
            <v-subheader>Ticket Types</v-subheader>
            <pie-chart
              v-if="tickets_per_type"
              :chart-data="tickets_per_type"
              :styles="{ width: '100%', position: 'relative' }"
              :chart-options="{
                onClick: clickPie,
                hover: {
                  onHover: function(e) {
                     var point = this.getElementAtEvent(e);
                     if (point.length) e.target.style.cursor = 'pointer';
                     else e.target.style.cursor = 'default';
                  }
               }
              }"
            >
            </pie-chart>
          </v-col>
        </v-row>
      </v-col>
      <v-col cols="12" lg="5">
        <TicketList :type="this.$route.params.type" @click="open"></TicketList>
      </v-col>
    </v-row>
  </v-main>
</template>

<script lang="ts">
import Vue from "vue";
import LineChart from "../components/charts/Line";
import BarChart from "../components/charts/Bar";
import PieChart from "../components/charts/Doughnut";
import { API } from "@/services/api";
import {Statistics, TicketResponse} from "@/client";
import {DateTime} from "luxon";
import { colors } from "@/plugins/vuetify";
import TicketList from "@/components/TicketList.vue";
import { createHash } from "crypto";

export default Vue.extend({
  name: "Dashboard",
  components: {
    LineChart,
    BarChart,
    PieChart,
    TicketList
  },
  data() {
    return {
      statistics: (undefined as unknown) as Statistics
    };
  },
  computed: {
    tickets_per_type: function () {
      let data = { labels: [] as Array<string>, datasets: [{ backgroundColor: [] as Array<string>, data: [] as Array<number> }] }
      this.lodash.forEach(this.statistics.tickets_per_type, (count, type) => {
        data.labels.push(type);
        data.datasets[0].data.push(count);

        data.datasets[0].backgroundColor.push(this.color(type));
      })
      return data
    },
    open_tickets_per_user: function () {
      let data = { labels: [] as Array<string>, datasets: [{ backgroundColor: [] as Array<string>, data: [] as Array<number> }] }
      this.lodash.forEach(this.statistics.open_tickets_per_user, (count, user) => {
        if (!user) {
          data.labels.push("unassigned");
        } else {
          data.labels.push(user);
        }
        data.datasets[0].data.push(count);
        data.datasets[0].backgroundColor.push(this.color(user));
      })
      return data
    },
    tickets_per_week: function () {
      let data = {labels: [] as Array<string>, datasets: [{backgroundColor: [] as Array<string>, data: [] as Array<number> }]}
      this.lodash.forEach(this.weeks(), (week) => {
        data.labels.push(week);
        if (week in this.statistics.tickets_per_week) {
          data.datasets[0].data.push(this.statistics.tickets_per_week[week]);
        } else {
          data.datasets[0].data.push(0);
        }
        data.datasets[0].backgroundColor.push("#607d8b");
      })
      return data
    }
  },
  methods: {
    open: function (ticket: TicketResponse) {
      if (ticket.id === undefined) {
        return;
      }

      this.$router.push({
        name: "Ticket",
        params: {type: '-', id: ticket.id.toString()}
      });
    },
    clickUser: function (evt, elem) {
      let owner = this.open_tickets_per_user.labels[elem[0]._index];
      let query = 'status == \'open\' AND owner == \'' + owner + '\'';

      if (owner == 'unassigned') {
        query = 'status == \'open\' AND !owner';
      }

      this.$router.push({
        name: "TicketList",
        params: {query: query}
      });
    },
    clickPie: function (evt, elem) {
      this.$router.push({
        name: "TicketList",
        params: {type: this.tickets_per_type.labels[elem[0]._index]}
      });
    },
    color: function (s: string): string {
      let pos = createHash('md5').update(s).digest().readUInt32BE(0) % colors.length;
      return colors[pos];
    },
    fillData() {
      API.getStatistics().then(response => {
        this.statistics = response.data;
      });
    },
    weeks: function () {
      let w = [] as Array<string>;
      for (let i = 0; i < 53; i++) {
        w.push(DateTime.utc().minus({ weeks: i }).toFormat("kkkk-WW"))
      }
      this.lodash.reverse(w);
      return w
    }
  },
  mounted() {
    this.fillData();
  }
});
</script>

<style>
canvas {
  position: relative !important;
  margin: 0 auto;
}
</style>
