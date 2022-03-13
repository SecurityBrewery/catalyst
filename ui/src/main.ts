import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";

import JsonSchemaEditor from "json-schema-editor-vue";
import vuetify from "./plugins/vuetify";
import VuePipeline from "vue-pipeline";
import VueLodash from "vue-lodash";
import lodash from "lodash";
import axios from "axios";
import { DateTime } from 'luxon';
import VueNativeSock from 'vue-native-websocket';
import antInputDirective from 'ant-design-vue/es/_util/antInputDirective'
import antDirective from 'ant-design-vue/es/_util/antDirective'

import VueAxios from "vue-axios";
import VueLuxon from "vue-luxon";

import "./registerServiceWorker";

import "json-schema-editor-vue/lib/json-schema-editor-vue.css";
import "@mdi/font/css/materialdesignicons.css";
import "vue-d3-network/dist/vue-d3-network.css";
import '@koumoul/vjsf/dist/main.css'

import { Problem } from "@/types/types";
Vue.use(VueLodash, { lodash: lodash });
Vue.use(antDirective);
Vue.use(antInputDirective);
Vue.use(JsonSchemaEditor);
Vue.use(VuePipeline);
Vue.use(VueAxios, axios);
Vue.use(VueLuxon);

// import VJsf from '@koumoul/vjsf'
import VJsf from '@koumoul/vjsf/lib/VJsf.js';
import '@koumoul/vjsf/lib/deps/third-party.js';

Vue.component('VJsf', VJsf)

Vue.config.productionTip = false;

// eslint-disable-next-line @typescript-eslint/no-explicit-any
Vue.filter("capitalize", function(value: any) {
  if (!value) return "";
  return lodash.startCase(value.toString());
});

Vue.filter("formatdate", function(s: string, format: string) {
    if (!format) {
        return DateTime.fromISO(s).toLocaleString(DateTime.DATETIME_SHORT);
    }
    return DateTime.fromISO(s).toFormat(format);
});

let protocol = "ws"
if (location.protocol === "https:") {
    protocol = "wss"
}
Vue.use(VueNativeSock, protocol + '://' + location.hostname + ':'+ location.port +'/wss', { store: store, format: 'json' })

const v = new Vue({
  router,
  vuetify,
  store,
  render: h => h(App)
}).$mount("#app");

axios.interceptors.response.use(
  // response => response,
  response => {
      lodash.unset(response.data, 'notoast');

      return Promise.resolve(response);
  },
  error => {
      if (!lodash.has(error.response.data, 'notoast')) {
          if (error.response.data && 'title' in error.response.data && 'detail' in error.response.data) {
              const problem = error.response.data as Problem;
              v.$store.dispatch("alertError", { name: problem.title, detail: problem.detail });
              return Promise.reject(error);
          }
          if (error.response.data && 'error' in error.response.data) {
              v.$store.dispatch("alertError", { name: "Error", detail: error.response.data.error });
              return Promise.reject(error);
          }
          v.$store.dispatch("alertError", { name: "Error", detail: JSON.stringify(error.response.data) });
      }

      return Promise.reject(error);
  }
);
