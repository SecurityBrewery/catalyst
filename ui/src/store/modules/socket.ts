import Vue from "vue";
import Vuex, {ActionContext} from "vuex";
import lodash from "lodash";

Vue.use(Vuex);

interface State {
    socket: any;
}

export const socketStore = {
    state: (): State => ({
        socket: {
            isConnected: false,
            message: '',
            reconnectError: false,
        }
    }),
    mutations: {
        SOCKET_ONOPEN (state: State, event: any)  {
            // console.log("SOCKET_ONOPEN");
            Vue.prototype.$socket = event.currentTarget;
            state.socket.isConnected = true;
        },
        SOCKET_ONCLOSE (state: State, event: any)  {
            // console.log("SOCKET_ONCLOSE");
            state.socket.isConnected = false;
        },
        SOCKET_ONERROR (state: State, event: any)  {
            // console.log("SOCKET_ONERROR");
            console.error(state, event);
        },
        // default handler called for all methods
        SOCKET_ONMESSAGE (state: State, message: any)  {
            // console.log("SOCKET_ONMESSAGE");
            state.socket.message = message;
        },
        // mutations for reconnect methods
        SOCKET_RECONNECT(state: State, count: any) {
            // console.log("SOCKET_RECONNECT");
            console.info(state, count);
        },
        SOCKET_RECONNECT_ERROR(state: State) {
            // console.log("SOCKET_RECONNECT_ERROR");
            state.socket.reconnectError = true;
        },
    },
    actions: {
        sendMessage: function(context: ActionContext<any, any>, msg: any) {
            Vue.prototype.$socket.send(msg);
        },
        update: function (context: ActionContext<any, any>, msg: any) {
            // console.log("update", msg);
            if (!msg || !(lodash.has(msg, "ids")) || !msg["ids"]) {
                return
            }
            Vue.lodash.forEach(msg["ids"], (id) => {
                if (lodash.startsWith(id, "settings/")) {
                    context.dispatch("getSetting")
                }
            });
        }
    }
}
