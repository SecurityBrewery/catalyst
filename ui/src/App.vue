<template>
  <div>
    <v-app v-if="!authenticated" id="app" class="background-dark">
      <v-container class="login d-flex flex-column justify-center">
        <v-form ref="form" v-model="valid" @submit.prevent="login">
          <v-card class="pa-4">
            <div class="d-flex justify-center">
              <v-img src="/ui/static/flask.svg" height="100" width="100" class="flex-grow-0"></v-img>
            </div>
            <v-card-title class="text-center justify-center">
              Catalyst Login
            </v-card-title>
            <v-card-text v-if="hassimple" class="text-center">
              <v-text-field id="username" name="username" label="Name" v-model="username" :rules="[
                v => !!v || 'Name is required',
              ]"></v-text-field>
              <v-text-field
                  id="password"
                  name="password"
                  label="Password"
                  :append-icon="show ? 'mdi-eye' : 'mdi-eye-off'"
                  :type="show ? 'text' : 'password'"
                  @click:append="show = !show"
                  v-model="password"
                  :rules="[
                    v => !!v || 'Password is required',
                    // v => (v && v.length > 8) || 'Password must be more than 8 characters',
                  ]"></v-text-field>
            </v-card-text>
            <v-card-actions class="justify-center">
              <v-btn v-if="hasoidc" text href="/auth/oidclogin">
                Login with OIDC
              </v-btn>
              <v-spacer v-if="hassimple"></v-spacer>
              <v-btn v-if="hassimple" type="submit" color="primary" elevation="0" :disabled="!valid">
                Login
              </v-btn>
            </v-card-actions>
          </v-card>
        </v-form>
      </v-container>
    </v-app>
    <v-app v-else class="background">
      <v-navigation-drawer dark permanent :mini-variant="mini" :expand-on-hover="mini" app color="statusbar">
        <v-list>
          <v-list-item class="px-2" :to="{ name: 'Home' }">
            <v-list-item-avatar rounded="0">
              <v-img src="/ui/static/flask_white.svg" :width="40"></v-img>
            </v-list-item-avatar>
            <v-list-item-content>
              <v-list-item-title class="title">
                Catalyst
              </v-list-item-title>
            </v-list-item-content>
          </v-list-item>
        </v-list>

        <v-list nav dense>
          <v-list-item>
            <v-list-item-icon>
              <v-icon class="my-1">mdi-arrow-right-bold</v-icon>
            </v-list-item-icon>
            <v-list-item-title>
              <v-text-field
                  placeholder="Goto"
                  outlined
                  dense
                  hide-details
                  v-on:keyup.enter="enter"
                  clearable
                  color="#fff"
                  v-model="goto"></v-text-field>
            </v-list-item-title>
          </v-list-item>
        </v-list>
        <v-divider></v-divider>

        <AppLink :links="internal"></AppLink>

        <v-list nav dense v-if="$store.state.settings.ticketTypes">
          <v-list-item
              v-for="customType in $store.state.settings.ticketTypes"
              :key="customType.id"
              link
              :class="{ 'v-list-item--active': ($route.params.type === customType.id) }"
              @click="openTicketList(customType.id)">
            <v-list-item-icon>
              <v-badge
                  v-if="customType.id in counts && counts[customType.id] > 0"
                  :content="counts[customType.id]"
                  color="red"
                  left
                  offset-x="35"
                  offset-y="8"
                  bottom>
                <v-icon>{{ customType.icon }}</v-icon>
              </v-badge>
              <v-icon v-else>{{ customType.icon }}</v-icon>
            </v-list-item-icon>
            <v-list-item-title>{{ customType.name }}</v-list-item-title>
          </v-list-item>
        </v-list>

        <v-divider></v-divider>

        <AppLink :links="settings"></AppLink>

        <template v-slot:append>

          <v-list nav dense>
            <v-list-item class="version" dense style="min-height: 20px">
              <v-list-item-content>
                <v-list-item-title style="text-align: center; opacity: 0.5;">
                  {{ $store.state.settings.tier }} v{{ $store.state.settings.version }}
                </v-list-item-title>
              </v-list-item-content>
            </v-list-item>
          </v-list>
          <v-divider></v-divider>
          <v-list nav dense>

            <v-list-item :to="{ name: 'API' }">
              <v-list-item-icon>
                <v-icon>mdi-share-variant</v-icon>
              </v-list-item-icon>

              <v-list-item-content>
                <v-list-item-title>API Documentation</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
          </v-list>
        </template>
      </v-navigation-drawer>
      <v-app-bar app dense flat absolute color="transparent">
        <v-btn icon @click="mini = !mini">
          <v-icon id="toggle_menu" color="primary">mdi-menu</v-icon>
        </v-btn>

        <v-breadcrumbs :items="crumbs">
          <template v-slot:item="{ item }">
            <v-breadcrumbs-item
                :to="item.to"
                class="text-subtitle-2 crumb-item"
                :disabled="item.disabled"
                exact
            >
              {{ item.text }}
            </v-breadcrumbs-item>
          </template>
        </v-breadcrumbs>

        <v-spacer></v-spacer>

        <v-menu left bottom offset-y>
          <template v-slot:activator="{ on, attrs }">
            <v-btn
                icon
                v-bind="attrs"
                v-on="on"
            >
              <v-avatar v-if="$store.state.userdata.image" size="32">
                <v-img :src="$store.state.userdata.image"></v-img>
              </v-avatar>
              <v-icon v-else>mdi-account-circle</v-icon>
            </v-btn>
          </template>

          <v-list>
            <v-list-item :to="{ name: 'Profile' }">
              <v-list-item-title>Account</v-list-item-title>
              <v-list-item-icon><v-icon>mdi-account-circle</v-icon></v-list-item-icon>
            </v-list-item>
            <v-list-item @click="logout">
              <v-list-item-title>Logout</v-list-item-title>
              <v-list-item-icon><v-icon>mdi-logout</v-icon></v-list-item-icon>
            </v-list-item>
          </v-list>
        </v-menu>

      </v-app-bar>
      <router-view></router-view>
      <v-snackbar id="alert" v-model="snackbar" :color="$store.state.alert.type" :timeout="$store.state.alert.type === 'error' ? -1 : 5000" outlined>
        <b style="display: block">{{ $store.state.alert.name | capitalize }}</b>
        {{ $store.state.alert.detail }}
        <template v-slot:action="{ attrs }">
          <v-btn text v-bind="attrs" @click="snackbar = false">Close</v-btn>
        </template>
      </v-snackbar>
    </v-app>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import AppLink from "./components/AppLink.vue";
import router from "vue-router";
import {API} from "@/services/api";

export default Vue.extend({
  name: "App",
  components: {AppLink},
  data: () => ({
    show: false,
    hassimple: false,
    hasoidc: false,
    username: "",
    password: "",
    valid: true,
    authenticated: false,
    settings: [
      { icon: "mdi-format-list-bulleted-type", name: "Ticket Types", to: "TicketTypeList", role: "engineer:tickettype:write" },
      { icon: "mdi-file-hidden", name: "Templates", to: "TemplateList", role: "analyst:template:read" },
      { icon: "mdi-file-cog-outline", name: "Playbooks", to: "PlaybookList", role: "analyst:playbook:read" },
      { icon: "mdi-flash", name: "Automations", to: "AutomationList", role: "analyst:automation:read" },
      { icon: "mdi-filter", name: "Ingestion Rules", to: "RuleList", role: "analyst:rule:read", tier: "enterprise" },
      { icon: "mdi-account", name: "Users & API Keys", to: "UserList", role: "admin:user:write" },
      { icon: "mdi-account-group", name: "Groups", to: "GroupList", role: "admin:group:write", tier: "enterprise" },
      { icon: "mdi-cogs", name: "User Data", to: "UserDataList", role: "admin:userdata:write" },
      { icon: "mdi-format-list-checks", name: "Jobs", to: "JobList", role: "admin:job:write" },
      { icon: "mdi-cog", name: "Settings", to: "Settings", role: "admin:settings:write" },
    ],
    mini: true,
    goto: "",

    snackbar: false,
  }),
  watch: {
    showAlert: function () {
      this.snackbar = true
    }
  },
  computed: {
    counts: function (): number {
      return this.$store.state.counts
    },
    internal: function (): Array<any> {
      return  [
        { icon: "mdi-view-dashboard", name: "Dashboards", to: "DashboardList", role: "analyst:dashboard:read" },
        { icon: "mdi-check-bold", name: "Open Tasks", to: "TaskList", count: this.$store.state.task_count },
      ]
    },
    showAlert: function (): boolean {
      return this.$store.state.showAlert
    },
    crumbs: function() {
      this.$route.name

      let pathArray = this.$route.path.split("/")
      pathArray.shift()

      return this.lodash.reduce(pathArray, (breadcrumbs, path, idx) => {
        let to = {};
        let text = path;

        let toPath = breadcrumbs[idx - 1] ? "/" + breadcrumbs[idx - 1].xpath + "/" + path : "/" + path;
        let resolved = this.$router.resolve(toPath);
        if (resolved) {
          to = { name: resolved.resolved.name, params: resolved.resolved.params };
          text = resolved.resolved.meta && resolved.resolved.meta.title ? resolved.resolved.meta.title : text;
        }

        breadcrumbs.push({ xpath: path, to: to, text: text });
        return breadcrumbs;
      }, [] as Array<any>);
    }
  },

  methods: {
    enter: function () {
      if (!this.goto) {
        return
      }
      this.$router.push({
        name: "Ticket",
        params: { id: this.goto.toString() }
      });
    },
    openTicketList: function (type: string) {
      (this.$router as any).history.current = router.START_LOCATION;
      this.$router.push({ name: "TicketList", params: { type: type } });
    },
    hasRole: function (s: string) {
      if (this.$store.state.user.roles) {
        return this.lodash.includes(this.$store.state.user.roles, s);
      }
      return false;
    },
    login: function () {
      this.axios.post(
          "/auth/login",
          {username: this.username, password: this.password},
      ).then((response) => {
        console.log(response.data);
        if (!this.lodash.isObject(response.data)) {
          return
        }

        this.$store.dispatch("getUser");
        this.$store.dispatch("getUserData");
        this.$store.dispatch("getSettings");

        this.authenticated = true;
      }).catch(() => {
        this.valid = false;
      })
    },
    logout: function () {
      this.axios.post("/auth/logout").then(() => {
        this.authenticated = false;
      })
    }
  },
  mounted() {
    this.axios.get("/auth/config").then((response) => {
      this.hassimple = response.data.simple;
      this.hasoidc = response.data.oidc;

      API.currentUser().then((response) => {
        if (!this.lodash.isObject(response.data)) {
          if (!this.hassimple && this.hasoidc) {
            window.location.href = "/auth/oidclogin";
          }

          return
        }

        this.authenticated = true;

        this.$store.dispatch("getUser");
        this.$store.dispatch("getUserData");
        this.$store.dispatch("getSettings");
      })
    }).catch(() => {
      this.hassimple = false;
      this.hasoidc = false;
    })
  },
});
</script>

<style>
.background {
  background-color: #f5f5f5 !important;
}

.background-dark {
  background-color: #212121 !important;
}

.login {
  height: 100%;
  max-width: 400px !important;
}

.v-app-bar.v-toolbar--dense .v-toolbar__content {
  border-bottom: 1px solid #e0e0e0 !important;
}

.theme--dark.background {
  background-color: #303030 !important;
}

.theme--dark h1,
.theme--dark h2,
.theme--dark h3,
.theme--dark h4 {
  color: white !important;
}

.theme--dark.v-btn:hover::before {
  opacity: 0 !important;
}

.theme--dark .v-application .primary--text {
  text-shadow: 0 0 3px #FFC107; /* #00bcd4;/* 0 0 3px #ffff00; */
}

.theme--dark .glow,
.theme--dark .v-list .v-list-item--active,
.theme--dark.v-btn:hover,
.theme--dark a:hover {
  color: #FFC107 !important; /* #00bcd4 !important;/* #ffff00 !important; */
  text-shadow: 0 0 3px #FFC107 !important; /* #00bcd4 !important;/* 0 0 3px #ffff00 !important; */
}

/* box-shadow: 0 0 8px rgba(255, 255, 0, 0.2) !important; */

.v-navigation-drawer--mini-variant .version {
  opacity: 0;
  transition: opacity 0.2s, visibility 0.2s;
}
</style>
