module.exports = {
  publicPath: "/ui/static/",
  transpileDependencies: ["vuetify", "@koumoul/vjsf"],
  pwa: {
    name: "Catalyst",
    themeColor: "#FFC107",
    msTileColor: "#000000",
    appleMobileWebAppCapable: "yes",
    appleMobileWebAppStatusBarStyle: "black",

    // configure the workbox plugin
    // workboxPluginMode: "InjectManifest",
    // workboxOptions: {
    //   // swSrc is required in InjectManifest mode.
    //   swSrc: "dev/sw.js",
    //   // ...other Workbox options...
    // },
    iconPaths: {
      favicon32: "img/icons/favicon-32x32.png",
      favicon16: "img/icons/favicon-16x16.png",
      appleTouchIcon: "img/icons/apple-touch-icon-152x152.png",
      maskIcon: "img/icons/safari-pinned-tab.svg",
      msTileImage: "img/icons/msapplication-icon-144x144.png",
    },
    manifestCrossorigin: 'use-credentials'
  },
  devServer: {
    disableHostCheck: true
  }
};
