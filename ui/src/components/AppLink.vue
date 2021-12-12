<template>
  <div>
    <v-list nav dense>
      <v-list-item v-for="link in userLinks" :key="link.to" link :to="{ name: link.to }">
        <v-list-item-icon>
          <v-badge
              v-if="'count' in link && link.count"
              :content="link.count"
              color="red"
              left
              offset-x="35"
              offset-y="8"
              bottom>
            <v-icon>{{ link.icon }}</v-icon>
          </v-badge>
          <v-icon v-else>{{ link.icon }}</v-icon>
        </v-list-item-icon>
        <v-list-item-title>{{ link.name }}</v-list-item-title>
      </v-list-item>
    </v-list>
    <v-divider v-if="userLinks.length > 0"></v-divider>
  </div>
</template>

<script lang="ts">
import Vue from "vue";

export default Vue.extend({
  name: "AppLink",
  props: ["links"],
  computed: {
    userLinks: function (): Array<any> {
      return this.lodash.filter(this.links, link => {
        return this.hasRole(link) && this.hasTier(link)
      })
    }
  },
  methods: {
    hasRole: function (link: any) {
      if (!("role" in link)) {
        return true;
      }
      let has = false;
      if (this.$store.state.user.roles) {
        this.lodash.forEach(this.$store.state.user.roles, (userRole) => {
          if (link.role === userRole || this.lodash.startsWith(link.role, userRole + ":")) {
            has = true;
          }
        })
      }
      return has;
    },
    hasTier: function (link: any): boolean {
      if ("tier" in link) {
        if (this.$store.state.settings.tier) {
          return this.$store.state.settings.tier == link.tier;
        }
        return false;
      }
      return true;
    }
  },
});
</script>
