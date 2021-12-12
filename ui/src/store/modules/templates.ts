import Vue from "vue";
import Vuex, {ActionContext} from "vuex";
import {TicketTemplate} from "@/client";
import {API} from "@/services/api";
import {AxiosResponse} from "axios";

Vue.use(Vuex);

interface State {
    templates: Array<TicketTemplate>;
}

export const templateStore = {
    state: (): State => ({
        templates: [],
    }),
    mutations: {
        setTemplates(state: State, msg: Array<TicketTemplate>) {
            state.templates = msg;
        },
    },
    actions: {
        listTemplates(context: ActionContext<any, any>) {
            API.listTemplates().then((response: AxiosResponse<Array<TicketTemplate>>) => {
                context.commit("setTemplates", response.data)
            });
        },
        getTemplate(context: ActionContext<any, any>, id: string) {
            return new Promise((resolve) => {
                API.getTemplate(id).then((response: AxiosResponse<TicketTemplate>) => {
                    resolve(response.data);
                }).catch(error => {
                    context.dispatch("alertError", {name: "Template could not be loaded", details: error});
                });
            });
        },
        addTemplate(context: ActionContext<any, any>, template: TicketTemplate) {
            return new Promise((resolve) => {
                API.createTemplate(template).then(() => {
                    context.dispatch("listTemplates").then(() => {
                        context.dispatch("alertSuccess", {name: "Template created"}).then(() => {
                            resolve({})
                        });
                    }).catch(error => {
                        context.dispatch("alertError", {name: "Template created, but reload failed", details: error});
                    });
                }).catch(error => {
                    context.dispatch("alertError", {name: "Template could not be created", details: error});
                });
            });
        },
        updateTemplate(context: ActionContext<any, any>, msg: any) {
            API.updateTemplate(msg.id, msg.template).then(() => {
                context.dispatch("alertSuccess", {name: "Template updated"});
            }).catch(error => {
                context.dispatch("alertError", {name: "Template could not be updated", details: error});
            });
        },
        deleteTemplate(context: ActionContext<any, any>, id: string) {
            return new Promise((resolve) => {
                API.deleteTemplate(id).then(() => {
                    context.dispatch("listTemplates").then(() => {
                        context.dispatch("alertSuccess", {name: "Template deleted"}).then(() => {
                            resolve({});
                        });
                    }).catch(error => {
                        context.dispatch("alertError", {name: "Template deleted, but reload failed", details: error});
                    });
                }).catch(error => {
                    context.dispatch("alertError", {name: "Template could not be deleted", details: error});
                });
            });
        },
    }
}
