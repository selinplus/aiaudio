import Vue from 'vue'
import Vuex from 'vuex'
import db from './localstorage'
Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    token: '',
    username: '',
  },
  mutations: {
    setToken(state, val) {
      state.token = val
    },
    setUsername(state, val) {
      state.username = val
    },
    getUsername(state) {
      return state.username
    },
    getToken(state){
        return state.username
    }
  },
  actions: {
    saveToken(context,token){
        db.save('token', token).then(()=>
        context.commit('setToken',token)
        )
    },
    saveUsername(context, username){
        db.save('username', username).then(
            () => context.commit('setUsername', username)
        )
    },    
  }
})
