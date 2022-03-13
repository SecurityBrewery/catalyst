import Vue from "vue";
import VueRouter, { RouteConfig, RawLocation, Route } from "vue-router";
import ArtifactPopup from "../views/ArtifactPopup.vue";
import Ticket from "../views/Ticket.vue";
import TicketList from "../views/TicketList.vue";
import Graph from "../views/Graph.vue";
import Playbook from "../views/Playbook.vue";
import PlaybookList from "../views/PlaybookList.vue";
import Automation from "../views/Automation.vue";
import UserData from "../views/UserData.vue";
import Profile from "../views/Profile.vue";
import UserDataList from "../views/UserDataList.vue";
import AutomationList from "../views/AutomationList.vue";
import Rule from "../views/Rule.vue";
import RuleList from "../views/RuleList.vue";
import Template from "../views/Template.vue";
import TemplateList from "../views/TemplateList.vue";
import Dashboard from "../views/Dashboard.vue";
import DashboardList from "../views/DashboardList.vue";
import API from "../views/API.vue";
import User from '../views/User.vue';
import UserList from "@/views/UserList.vue";
import Job from '../views/Job.vue';
import JobList from "@/views/JobList.vue";
import GroupList from "@/views/GroupList.vue";
import Home from "@/views/Home.vue";
import Group from "@/views/Group.vue";
import TicketType from '../views/TicketType.vue';
import TicketTypeList from "@/views/TicketTypeList.vue";
import TaskList from "@/views/TaskList.vue";
import Settings from "@/views/Settings.vue";

Vue.use(VueRouter);

const originalPush = VueRouter.prototype.push;
VueRouter.prototype.push = function push(location: RawLocation): Promise<Route> {
  return new Promise((resolve, reject) => {
    originalPush.call(this, location, () => {
      // on complete

      resolve(this.currentRoute);
    }, (error) => {
      // on abort

      // only ignore NavigationDuplicated error
      if (error.name === 'NavigationDuplicated') {
        resolve(this.currentRoute);
      } else {
        reject(error);
      }
    });
  });
};


const routes: Array<RouteConfig> = [
  {
    path: "/",
    name: "Catalyst",
    redirect: { name: "Dashboard" },
  },

  {
    path: "/home",
    name: "Home",
    component: Home,
    meta: { title: "Home" },
  },

  {
    path: "/profile",
    name: "Profile",
    component: Profile,
    meta: { title: "Profile" },
  },

  {
    path: "/tickets",
    name: "TicketList",
    component: TicketList,
    meta: { title: "Tickets" },
    props: true,
    children: [
        {
          path: ":id",
          name: "Ticket",
          component: Ticket,
        }
    ]
  },

  {
    path: "/tickets/:type?/:id/artifact/:artifact",
    name: "ArtifactPopup",
    meta: { title: "Artifact" },
    component: ArtifactPopup,
  },

  {
    path: "/tasks",
    name: "TaskList",
    meta: { title: "Tasks" },
    component: TaskList,
  },

  {
    path: "/templates",
    name: "TemplateList",
    component: TemplateList,
    meta: { title: "Templates" },
    children: [
      {
        path: ":id",
        name: "Template",
        component: Template,
      },
    ]
  },

  {
    path: "/tickettype",
    name: "TicketTypeList",
    component: TicketTypeList,
    meta: { title: "Ticket Types" },
    children: [
      {
        path: ":id",
        name: "TicketType",
        component: TicketType,
      },
    ]
  },

  {
    path: "/playbooks",
    name: "PlaybookList",
    component: PlaybookList,
    meta: { title: "Playbooks" },
    children: [
      {
        path: ":id",
        name: "Playbook",
        component: Playbook,
      },
    ]
  },

  {
    path: "/userdata",
    name: "UserDataList",
    component: UserDataList,
    meta: { title: "User Data" },
    children: [
      {
        path: ":id",
        name: "UserData",
        component: UserData,
      },
    ]
  },

  {
    path: "/jobs",
    name: "JobList",
    component: JobList,
    meta: { title: "Jobs" },
    children: [
      {
        path: ":id",
        name: "Job",
        component: Job,
      },
    ]
  },

  {
    path: "/automations",
    name: "AutomationList",
    component: AutomationList,
    meta: { title: "Automations" },
    children: [
      {
        path: ":id",
        name: "Automation",
        component: Automation,
      },
    ]
  },

  {
    path: "/rules",
    name: "RuleList",
    component: RuleList,
    meta: { title: "Rules" },
    children: [
      {
        path: ":id",
        name: "Rule",
        component: Rule,
      },
    ]
  },

  {
    path: "/users",
    name: "UserList",
    component: UserList,
    meta: { title: "Users" },
    children: [
      {
        path: ":id",
        name: "User",
        component: User,
      },
    ]
  },

  {
    path: "/groups",
    name: "GroupList",
    component: GroupList,
    meta: { title: "Groups" },
    children: [
      {
        path: ":id",
        name: "Group",
        component: Group,
      },
    ]
  },


  {
    path: "/dashboards",
    name: "DashboardList",
    component: DashboardList,
    meta: { title: "Dashboards" },
    children: [
      {
        path: ":id",
        name: "Dashboard",
        component: Dashboard,
      },
    ]
  },

  {
    path: "/settings",
    name: "Settings",
    component: Settings,
    meta: { title: "Settings" },
  },

  {
    path: "/apidocs",
    name: "API",
    component: API,
    meta: { title: "API" },
  },

  {
    path: "/graph/:col/:id",
    name: "Graph",
    component: Graph,
    meta: { title: "Graph" },
  },
];

const router = new VueRouter({
  mode: 'history',
  routes,
});

export default router;
