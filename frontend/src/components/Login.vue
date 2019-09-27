<template>
  <v-app>
    <v-content>
      <v-container fluid fill-height>
        <v-layout align-center justify-center>
          <v-flex xs12 sm8 md4>
            <v-card class="elevation-12">
              <v-toolbar dark>
                <v-toolbar-title>{{title}}</v-toolbar-title>
                <v-spacer></v-spacer>
                <v-tooltip bottom>
                  <template v-slot:activator="{ on }">
                    <v-btn :href="source" icon large target="_blank" v-on="on">
                      <v-icon large>code</v-icon>
                    </v-btn>
                  </template>
                  <span>Source</span>
                </v-tooltip>
              </v-toolbar>
              <v-card-text>
                <v-form>
                  <v-text-field
                    prepend-icon="person"
                    v-model="username"
                    name="login"
                    label="用户名"
                    type="text"
                  ></v-text-field>
                  <v-text-field
                    id="password"
                    prepend-icon="lock"
                    v-model="password"
                    name="password"
                    label="密码"
                    type="password"
                  ></v-text-field>
                </v-form>
              </v-card-text>
              <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn color="amber" @click="doLogin">登录</v-btn>
              </v-card-actions>
            </v-card>
            <v-card class="elevation-12">
              <v-toolbar dark>
                <v-spacer></v-spacer>
                <v-toolbar-title>{{version}}</v-toolbar-title>                
              </v-toolbar>              
            </v-card>
          </v-flex>
          <v-snackbar
            v-model="snackbar"
            :color="color"
            :multi-line="mode === 'multi-line'"
            :timeout="timeout"
            :vertical="mode === 'vertical'"
          >
            {{ text }}
            <v-btn dark flat @click="snackbar = false">Close</v-btn>
          </v-snackbar>
        </v-layout>
      </v-container>
    </v-content>
  </v-app>
</template>

<script>
import { mapMutations } from "vuex";
export default {
  data: () => ({
    drawer: null,
    username: "",
    password: "",
    snackbar: false,
    color: "",
    mode: "",
    timeout: 2000,
    text: ""
  }),

  props: {
    title: String,
    version: String
  },

  methods: {
    ...mapMutations({
      addToken: "setToken",
      addUsername: "setUsername"
    }),
    doLogin() {
      if (this.username && this.password) {
        let self = this;
        window.backend.Login(this.username, this.password)
          .then(res => {
            let data = JSON.parse(res);
            if (data.code == 200) {
              // self.addToken(data.data.token);
              // self.addUsername(self.username);
              self.text = data.msg;
              self.snackbar = true;
              self.$router.push("home/index");
            } else {
              self.text = data.msg;
              self.snackbar = true;
            }
          });
      } else {
        this.text = "用户名密码未填写";
        this.snackbar = true;
      }
    }
  }
};
</script>