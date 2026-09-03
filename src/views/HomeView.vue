<script setup>
const navigation = [
  { key: 'event', label: '活動專區', image: '/assets/nav-event.png', to: '/event' },
  { key: 'map', label: '海派地圖', image: '/assets/nav-map.png', to: '/map' },
  { key: 'protect', label: '海派護照', image: '/assets/nav-protect.png', to: '/protect' },
  { key: 'gift', label: '海派好禮', image: '/assets/nav-gift.png', to: '/gift' },
  { key: 'food', label: '海派美食', image: '/assets/nav-food.png', to: '/food' },
]
</script>

<template>
  <main class="home" aria-labelledby="page-title">
    <h1 id="page-title" class="sr-only">2026 新北海派－漁你共遊，鮮味港覺</h1>

    <div class="desktop-scene" aria-hidden="true"></div>
    <div class="mobile-scene" aria-hidden="true"></div>

    <nav class="desktop-navigation" aria-label="網站主要選單">
      <RouterLink
        v-for="item in navigation"
        :key="item.key"
        :class="['nav-art', `nav-art--${item.key}`]"
        :to="item.to"
        :aria-label="item.label"
      >
        <img :src="item.image" :alt="item.label" draggable="false" />
      </RouterLink>
    </nav>

    <nav class="mobile-navigation" aria-label="網站主要選單">
      <RouterLink
        v-for="item in navigation"
        :key="item.key"
        :class="['mobile-hotspot', `mobile-hotspot--${item.key}`]"
        :to="item.to"
        :aria-label="item.label"
      >
        <span class="sr-only">{{ item.label }}</span>
      </RouterLink>
    </nav>
  </main>
</template>
