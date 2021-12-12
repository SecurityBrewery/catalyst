import Vue from "vue";
import Vuex, {ActionContext} from "vuex";
import {API} from "@/services/api";
import {UserData, TicketList, Settings, UserResponse} from "@/client";
import {AxiosResponse} from "axios";
import {Alert} from "@/types/types";
import {templateStore} from "./modules/templates";
import {socketStore} from "@/store/modules/socket";

Vue.use(Vuex);

export default new Vuex.Store({
  modules: {
    templates: templateStore,
    socket: socketStore,
  },
  state: {
    user: {} as UserResponse,
    counts: {} as Record<string, number>,
    task_count: 0 as number,

    settings: {} as Settings,
    userdata: {} as UserData,

    alert: {} as Alert,
    showAlert: false as boolean,
  },
  getters: {
    timeformat: (state) => {
      if ('timeformat' in state.settings && state.settings.timeformat) {
        return state.settings.timeformat
      }
      return 'dd-MM-yyyy'
    }
  },
  mutations: {
    setUser (state, msg) {
      state.user = msg;
    },
    setCount (state, msg) {
      Vue.set(state.counts, msg.name, msg.count);
    },
    setTaskCount (state, msg) {
      state.task_count = msg;
    },
    setUserData (state, msg: UserData) {
      state.userdata = msg
    },
    setSettings (state, msg: Settings) {
      state.settings = msg
    },
    setAlert (state, msg: Alert) {
      state.showAlert = false;
      state.showAlert = true;
      state.alert = msg;
    }
  },
  actions: {
    getUser (context: ActionContext<any, any>) {
      API.currentUser().then((response) => {
        context.commit("setUser", response.data);
        context.dispatch("fetchCount");
      })
    },
    getUserData (context: ActionContext<any, any>) {
      API.currentUserData().then((response: AxiosResponse<UserData>) => {
        context.commit("setUserData", response.data);
      })
    },
    getSettings (context: ActionContext<any, any>) {
      API.getSettings().then((response: AxiosResponse<Settings>) => {
        context.commit("setSettings", response.data);
        context.dispatch("fetchCount");
      })
    },
    fetchCount (context: ActionContext<any, any>) {
      if (!context.state.user.id || !context.state.settings.ticketTypes) {
        return
      }

      const username = context.state.user.id;
      Vue.lodash.forEach(context.state.settings.ticketTypes, (t) => {

        API.listTickets(t.id,0,10,[],[], "status == 'open' AND (owner == '"+username+"' OR !owner)")
            .then((response: AxiosResponse<TicketList>) => {
              context.commit("setCount", {"name": t.id, "count": response.data.count});
            });
        API.listTasks().then((response) => {
          if (response.data) {
            context.commit("setTaskCount", response.data.length );
          }
        })
      })
    },
    alertSuccess(context: ActionContext<any, any>, msg) {
      msg.type = "success"
      context.commit("setAlert", msg)
    },
    alertError(context: ActionContext<any, any>, msg) {
      msg.type = "error"
      context.commit("setAlert", msg)
    },
  },
});
