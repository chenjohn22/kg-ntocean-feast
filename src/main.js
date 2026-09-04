import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import App from './App.vue'
import HomeView from './views/HomeView.vue'
import PlaceholderView from './views/PlaceholderView.vue'
import GiftView from './views/GiftView.vue'
import RegistrationView from './views/RegistrationView.vue'
import AdminView from './views/AdminView.vue'
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
    ...Object.entries(sections).filter(([path]) => path !== 'gift').map(([path, title]) => ({
      path: `/${path}`,
      name: path,
      component: PlaceholderView,
      props: { title },
    })),
    { path: '/gift', name: 'gift', component: GiftView },
    { path: '/gift/register', name: 'gift-register', component: RegistrationView },
    { path: '/kg-manager-admin', name: 'admin', component: AdminView },
    { path: '/:pathMatch(.*)*', redirect: '/' },
  ],
  scrollBehavior: () => ({ top: 0 }),
})

createApp(App).use(router).mount('#app')
