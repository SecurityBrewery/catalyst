<template>
  <v-container>
    <v-card style="background-color: #fff;" class="mt-8">
      <div class="swagger" id="swagger" ></div>
    </v-card>
  </v-container>
</template>

<script lang="ts">
import SwaggerUI from 'swagger-ui';
import 'swagger-ui/dist/swagger-ui.css';
import Vue from "vue";
import spec from '../../../generated/community.json';

export default Vue.extend({
  name: "API",
  components: {},
  mounted() {
    SwaggerUI({
      spec: spec,
      dom_id: '#swagger',
      oauth2RedirectUrl: location.href
    })
  }
});
</script>

<style>
#swagger a,
#swagger h1,
#swagger h2,
#swagger h3,
#swagger h4,
#swagger h5,
#swagger h6 {
  color: black !important
}

#swagger .info {
  padding: 0 10px 10px;
  border-radius: 4px;
}

#swagger .servers-title,
#swagger .servers {
  display: none;
}
</style>
