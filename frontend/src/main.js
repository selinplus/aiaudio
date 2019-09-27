import 'babel-polyfill';
import Vue from "vue";

// Setup Vuetify
import Vuetify from 'vuetify';
Vue.use(Vuetify);
import 'vuetify/dist/vuetify.min.css';
import 'material-design-icons-iconfont';
import router from './router';
import App from "./App.vue";

Vue.config.productionTip = false;
Vue.config.devtools = true;

import Bridge from "./wailsbridge";

Bridge.Start(() => {
  new Vue({
    router,
    render: h => h(App)
  }).$mount("#app");
});
