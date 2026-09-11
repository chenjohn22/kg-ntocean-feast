<script setup>
import { onBeforeUnmount, onMounted, ref } from 'vue'

const navigation = [
  { key: 'event', label: '活動專區', image: '/assets/nav-event.png', to: '/event' },
  { key: 'map', label: '海派地圖', image: '/assets/nav-map.png', to: '/map' },
  { key: 'protect', label: '海派護照', image: '/assets/nav-protect.png', to: '/protect' },
  { key: 'gift', label: '海派好禮', image: '/assets/nav-gift.png', to: '/gift' },
  { key: 'food', label: '海派美食', image: '/assets/nav-food.png', to: '/food' },
]

const desktopVideoId = 'siHjTgOEVK8'
const mobileVideoId = 'qPuN5CSND80'
const showIntro = ref(true)
const canCloseIntro = ref(false)
const introVideoSrc = ref('')
let closeTimer
let mobileQuery

function setVideoSource() {
  const videoId = mobileQuery?.matches ? mobileVideoId : desktopVideoId
  introVideoSrc.value = `https://www.youtube.com/embed/${videoId}?autoplay=1&mute=1&playsinline=1&rel=0&modestbranding=1`
}

function closeIntro() {
  if (!canCloseIntro.value) return
  showIntro.value = false
  introVideoSrc.value = ''
  document.body.classList.remove('intro-video-open')
}

function handleKeydown(event) {
  if (event.key === 'Escape') closeIntro()
}

onMounted(() => {
  mobileQuery = window.matchMedia('(max-width: 700px), (orientation: portrait) and (max-aspect-ratio: 3/4)')
  setVideoSource()
  mobileQuery.addEventListener('change', setVideoSource)
  document.addEventListener('keydown', handleKeydown)
  document.body.classList.add('intro-video-open')
  closeTimer = window.setTimeout(() => {
    canCloseIntro.value = true
  }, 10000)
})

onBeforeUnmount(() => {
  window.clearTimeout(closeTimer)
  mobileQuery?.removeEventListener('change', setVideoSource)
  document.removeEventListener('keydown', handleKeydown)
  document.body.classList.remove('intro-video-open')
})
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

    <div v-if="showIntro" class="intro-video" role="dialog" aria-modal="true" aria-label="活動開場影片">
      <div class="intro-video__frame">
        <iframe
          :src="introVideoSrc"
          title="2026 新北海派活動影片"
          allow="autoplay; encrypted-media; picture-in-picture"
          referrerpolicy="strict-origin-when-cross-origin"
          allowfullscreen
        ></iframe>
      </div>
      <Transition name="intro-close">
        <button
          v-if="canCloseIntro"
          class="intro-video__close"
          type="button"
          aria-label="關閉影片"
          @click="closeIntro"
        >
          <span aria-hidden="true"></span>
        </button>
      </Transition>
    </div>
  </main>
</template>
