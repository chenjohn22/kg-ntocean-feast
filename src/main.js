import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import App from './App.vue'
import HomeView from './views/HomeView.vue'
import ContentPageView from './views/ContentPageView.vue'
import RegistrationView from './views/RegistrationView.vue'
import AdminView from './views/AdminView.vue'
import './style.css'

const page = (title, slides, cta = null, scrolling = false, mobileSlides = []) => ({
  title,
  slides: slides.map((src, index) => ({
    src,
    alt: `${title}${slides.length > 1 ? `第 ${index + 1} 頁` : ''}`,
  })),
  mobileSlides: mobileSlides.map((src, index) => ({
    src,
    alt: `${title}手機版${mobileSlides.length > 1 ? `第 ${index + 1} 頁` : ''}`,
  })),
  cta,
  scrolling,
})

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', name: 'home', component: HomeView },
    {
      path: '/event',
      name: 'event',
      component: ContentPageView,
      props: page('活動專區', ['/assets/pages/event.jpg'], null, false, ['/assets/pages/mobile/event.jpg']),
    },
    {
      path: '/map',
      name: 'map',
      component: ContentPageView,
      props: page('海派地圖', ['/assets/pages/map-harbor.jpg', '/assets/pages/map.jpg'], null, true, [
        '/assets/pages/mobile/map-harbor.jpg',
        '/assets/pages/mobile/map.jpg',
      ]),
    },
    {
      path: '/protect',
      name: 'protect',
      component: ContentPageView,
      props: page('海派護照', ['/assets/pages/protect-restaurants.jpg'], null, false, [
        '/assets/pages/mobile/protect-1.jpg',
        '/assets/pages/mobile/protect-2.jpg',
      ]),
    },
    {
      path: '/gift',
      name: 'gift',
      component: ContentPageView,
      props: page('海派好禮', ['/assets/pages/gift.jpg'], {
        to: '/gift/register',
        label: '我要登錄',
        image: '/assets/pages/gift-register.png',
        mobileImage: '/assets/pages/mobile/gift-register.png',
      }, false, ['/assets/pages/mobile/gift.jpg']),
    },
    {
      path: '/food',
      name: 'food',
      component: ContentPageView,
      props: page(
        '海派美食',
        Array.from({ length: 6 }, (_, index) => `/assets/pages/food-${index + 1}.jpg`),
        null,
        true,
        Array.from({ length: 6 }, (_, index) => `/assets/pages/mobile/food-${index + 1}.jpg`),
      ),
    },
    { path: '/gift/register', name: 'gift-register', component: RegistrationView },
    { path: '/kg-manager-admin', name: 'admin', component: AdminView },
    { path: '/:pathMatch(.*)*', redirect: '/' },
  ],
  scrollBehavior: () => ({ top: 0 }),
})

createApp(App).use(router).mount('#app')
