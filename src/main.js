import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import App from './App.vue'
import HomeView from './views/HomeView.vue'
import PlaceholderView from './views/PlaceholderView.vue'
import './style.css'

const sections = {
  event: '活動專區',
  map: '海派地圖',
  protect: '海派護照',
  gift: '海派好禮',
  food: '海派美食',
}

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', name: 'home', component: HomeView },
    ...Object.entries(sections).map(([path, title]) => ({
      path: `/${path}`,
      name: path,
      component: PlaceholderView,
      props: { title },
    })),
    { path: '/:pathMatch(.*)*', redirect: '/' },
  ],
  scrollBehavior: () => ({ top: 0 }),
})

createApp(App).use(router).mount('#app')
