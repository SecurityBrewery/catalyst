import "@mdi/font/css/materialdesignicons.css";
import Vue from "vue";
import Vuetify from "vuetify/lib";

Vue.use(Vuetify);

export const colors = [
  '#c62828', '#6a1b9a',
  '#283593', '#0277bd',
  '#00695c', '#558b2f',
  '#f9a825', '#ef6c00',
  '#ad1457', '#4527a0',
  '#1565c0', '#00838f',
  '#2e7d32', '#9e9d24',
  '#ff8f00', '#d84315',
]

export default new Vuetify({
  theme: {
    dark: false,
    themes: {
      light: {
        statusbar: "#212121", // "#1c313a",
        appbar: "#eeeeee", //'#28282e',
        background: "#f5f5f5", //'#393a3f',
        cards: "#e0e0e0", //'#393a3f',

        primary: "#607d8b",
        yellow: "#FFC107", // "#FFEB3B",
        error: "#C62828", // "#d32f2f",
        info: "#1565C0", // "#2196F3",
        success: "#2E7D32", // "#689f38",
        warning: "#D84315", // "#fbc02d",

        red: "#C62828", // "#d32f2f",
      },
      dark: {
        // statusbar: "#d35400",
        // statusbar: "#f03f24", //'#18181c', 240 63 36
        // appbar: "#212121", //'#28282e',
        // background: "#303030", //'#393a3f',
        // cards: "#424242", //'#393a3f',

        statusbar: "#121212",
        appbar: "#212121", //'#28282e',
        background: "#303030", //'#393a3f',
        cards: "#424242", //'#393a3f',

        // alertmanager: "#ef6c00",
        // catalyst: "#ad1457",
        // blocklist: "#2e7d32",
        // uploadportal: "#283593",
        // search: "#c62828",

        // darksteel: "#1c313a",
        // steel: "#455a64",
        // lightsteel: "#718792",
        primary: "#FFC107", // "#00bcd4", // "#FFEB3B",
        yellow: "#FFC107", // "#FFEB3B",
        // accent: "#82B1FF",
        error: "#ef9a9a", // "#d32f2f",
        info: "#90caf9", // "#2196F3",
        success: "#a5d6a7", // "#689f38",
        warning: "#ffab91", // "#fbc02d",

        red: "#C62828", // "#d32f2f",
      },
    },
  },
  icons: {
    values: {
    },
  },
});
