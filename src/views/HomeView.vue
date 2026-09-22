<script setup>
import { nextTick, onBeforeUnmount, onMounted, ref } from 'vue'

const navigation = [
  { key: 'event', label: '活動專區', image: '/assets/nav-event.png', to: '/event' },
  { key: 'map', label: '海派地圖', image: '/assets/nav-map.png', to: '/map' },
  { key: 'protect', label: '海派護照', image: '/assets/nav-protect.png', to: '/protect' },
  { key: 'gift', label: '海派好禮', image: '/assets/nav-gift.png', to: '/gift' },
  { key: 'food', label: '海派美食', image: '/assets/nav-food.png', to: '/food' },
]

const featuredLinks = [
  {
    key: 'crab-news',
    label: '螃蟹快報',
    image: '/assets/nav-crab-news.png',
    href: 'https://www.facebook.com/share/1ERmgwB6Qi/?mibextid=wwXIfr',
  },
  {
    key: 'wanli-crab',
    label: '萬里蟹官網',
    image: '/assets/nav-wanli-crab.png',
    href: 'https://wanlicrab.tw',
  },
]

const desktopVideoId = 'yDoNhantIcM'
const mobileVideoId = 'qPuN5CSND80'
const showIntro = ref(true)
const canCloseIntro = ref(false)
const introVideoSrc = ref('')
const introVideoFrame = ref(null)
let closeTimer
let mobileQuery
let introPlayer

function loadYouTubeIframeAPI() {
  if (window.YT?.Player) return Promise.resolve(window.YT)

  return new Promise((resolve, reject) => {
    const previousReady = window.onYouTubeIframeAPIReady
    window.onYouTubeIframeAPIReady = () => {
      previousReady?.()
      resolve(window.YT)
    }

    if (document.querySelector('script[src="https://www.youtube.com/iframe_api"]')) return

    const script = document.createElement('script')
    script.src = 'https://www.youtube.com/iframe_api'
    script.async = true
    script.addEventListener('error', reject, { once: true })
    document.head.appendChild(script)
  })
}

function setVideoSource() {
  const videoId = mobileQuery?.matches ? mobileVideoId : desktopVideoId
  const origin = encodeURIComponent(window.location.origin)
  introVideoSrc.value = `https://www.youtube.com/embed/${videoId}?autoplay=1&mute=1&playsinline=1&controls=1&fs=0&rel=0&enablejsapi=1&origin=${origin}`
}

function dismissIntro() {
  showIntro.value = false
  introVideoSrc.value = ''
  document.body.classList.remove('intro-video-open')
}

function closeIntro() {
  if (!canCloseIntro.value) return
  dismissIntro()
}

function handlePlayerStateChange(event) {
  if (event.data === window.YT?.PlayerState.ENDED) dismissIntro()
}

async function setupIntroPlayer() {
  try {
    const YT = await loadYouTubeIframeAPI()
    await nextTick()
    if (!showIntro.value || !introVideoFrame.value) return

    introPlayer = new YT.Player(introVideoFrame.value, {
      events: { onStateChange: handlePlayerStateChange },
    })
  } catch (error) {
    console.error('YouTube IFrame API 載入失敗', error)
  }
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
  setupIntroPlayer()
  closeTimer = window.setTimeout(() => {
    canCloseIntro.value = true
  }, 10000)
})

onBeforeUnmount(() => {
  window.clearTimeout(closeTimer)
  mobileQuery?.removeEventListener('change', setVideoSource)
  document.removeEventListener('keydown', handleKeydown)
  introPlayer?.destroy()
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
      <a
        v-for="item in featuredLinks"
        :key="item.key"
        :class="['featured-link', `featured-link--${item.key}`]"
        :href="item.href"
        :aria-label="`${item.label}（另開新視窗）`"
        target="_blank"
        rel="noopener noreferrer"
      >
        <img :src="item.image" :alt="item.label" draggable="false" />
      </a>
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
      <a
        v-for="item in featuredLinks"
        :key="item.key"
        :class="['mobile-featured-link', `mobile-featured-link--${item.key}`]"
        :href="item.href"
        :aria-label="`${item.label}（另開新視窗）`"
        target="_blank"
        rel="noopener noreferrer"
      >
        <img :src="item.image" :alt="item.label" draggable="false" />
      </a>
    </nav>

    <div v-if="showIntro" class="intro-video" role="dialog" aria-modal="true" aria-label="活動開場影片">
      <div class="intro-video__frame">
        <iframe
          ref="introVideoFrame"
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
