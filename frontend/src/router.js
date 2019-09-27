import Vue from 'vue'
import Router from 'vue-router'
import Index from './views/Index.vue'
import Home from './views/Home.vue'
import Welcome from './components/Welcome'
import InComing from './components/InComing'
import Statistic from './components/Statistic'
import Settings from './components/Settings'
import HistoryComing from './components/HistoryComing'
Vue.use(Router)

// export default new Router({
//   routes: 
// })

let constRouter = [
  {
    path: '/',
    name: '登录页',
    component: Index
  },
  {
    path: '/home',
    name: '首页',
    component: Home,
    children: [
      {
        path: 'index',
        name: 'index',
        component: Welcome,
      },
      {
        path: 'incoming',
        name: '来电',
        component: InComing,
      },
      {
        path: 'history',
        name: '历史记录',
        component: HistoryComing,
      },
      {
        path: 'statistic',
        name: '统计信息',
        component: Statistic,
      },
      {
        path: 'settings',
        name: '设置',
        component: Settings,
      },
    ],
  }
]

let router = new Router({
  routes: constRouter
})

// const whiteList = '/'

//let asyncRouter

// 导航守卫，路由判断是否有TOKEN
// router.beforeEach((to, from, next) => {
//   if (whiteList === to.path) {
//     next()
//   }
//   let token = store.state.token
//   if (token.length) {
//     next(true)
//   } else {
//     next('/')
//   }
// });

// function go (to, next) {
//   asyncRouter = filterAsyncRouter(asyncRouter)
//   router.addRoutes(asyncRouter)
//   next({...to, replace: true})
// }

// function save (name, data) {
//   localStorage.setItem(name, JSON.stringify(data))
// }

// function get (name) {
//   return JSON.parse(localStorage.getItem(name))
// }

// function filterAsyncRouter (routes) {
//   return routes.filter((route) => {
//     let component = route.component
//     if (component) {
//       switch (route.component) {
//         case 'MenuView':
//           route.component = MenuView
//           break
//         case 'PageView':
//           route.component = PageView
//           break
//         case 'EmptyPageView':
//           route.component = EmptyPageView
//           break
//         case 'HomePageView':
//           route.component = HomePageView
//           break
//         default:
//           route.component = view(component)
//       }
//       if (route.children && route.children.length) {
//         route.children = filterAsyncRouter(route.children)
//       }
//       return true
//     }
//   })
// }

// function view (path) {
//   return function (resolve) {
//     import(`@/views/${path}.vue`).then(mod => {
//       resolve(mod)
//     })
//   }
// }

export default router
