<script setup>
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'

const props = defineProps({
  title: { type: String, required: true },
  slides: { type: Array, required: true },
  mobileSlides: { type: Array, required: true },
  cta: { type: Object, default: null },
  scrolling: { type: Boolean, default: false },
})

const sectionLinks = [
  { label: '首頁', to: '/' },
  { label: '活動專區', to: '/event' },
  { label: '海派地圖', to: '/map' },
  { label: '海派護照', to: '/protect' },
  { label: '海派好禮', to: '/gift' },
  { label: '海派美食', to: '/food' },
]

const current = ref(0)
const menu = ref(null)
const mobileMenu = ref(null)
const now = ref(Date.now())
const activeSlide = computed(() => props.slides[current.value])
const showCta = computed(() => (
  props.cta && (!props.cta.expiresAt || now.value < Date.parse(props.cta.expiresAt))
))
let clockTimer

watch(() => props.title, () => {
  current.value = 0
})

onMounted(() => {
  clockTimer = window.setInterval(() => {
    now.value = Date.now()
  }, 1000)
})

onBeforeUnmount(() => {
  window.clearInterval(clockTimer)
})

function previous() {
  current.value = (current.value - 1 + props.slides.length) % props.slides.length
}

function next() {
  current.value = (current.value + 1) % props.slides.length
}

function handleKeydown(event) {
  if (props.scrolling || props.slides.length < 2) return
  if (event.key === 'ArrowLeft') previous()
  if (event.key === 'ArrowRight') next()
}

function closeMenu() {
  menu.value?.removeAttribute('open')
  mobileMenu.value?.removeAttribute('open')
}

function hotspotStyle(link) {
  return {
    left: `${link.left}%`,
    top: `${link.top}%`,
    width: `${link.width}%`,
    height: `${link.height}%`,
  }
}
</script>

<template>
  <main class="content-page" :aria-label="title" tabindex="-1" @keydown="handleKeydown">
    <section :class="['content-desktop', { 'content-desktop--scroll': scrolling }]">
      <h1 class="sr-only">{{ title }}</h1>

      <div v-if="scrolling" class="content-scroll">
        <figure v-for="(slide, index) in slides" :key="slide.src" class="content-scroll__page content-image-stage">
          <img class="content-artwork" :src="slide.src" :alt="slide.alt" />
          <a
            v-for="link in slide.links"
            :key="`${link.name}-${link.address}`"
            class="content-map-hotspot"
            :href="link.href"
            target="_blank"
            rel="noopener noreferrer"
            :style="hotspotStyle(link)"
            :aria-label="`${link.name} Google Maps（另開新視窗）`"
            :title="`${link.name}｜開啟 Google Maps`"
          />
          <figcaption class="sr-only">{{ title }}第 {{ index + 1 }} 頁</figcaption>
        </figure>
      </div>

      <Transition v-else name="content-slide" mode="out-in">
        <div
          :key="activeSlide.src"
          class="content-image-stage"
        >
          <img class="content-artwork" :src="activeSlide.src" :alt="activeSlide.alt" />
          <a
            v-for="link in activeSlide.links"
            :key="`${link.name}-${link.address}`"
            class="content-map-hotspot"
            :href="link.href"
            target="_blank"
            rel="noopener noreferrer"
            :style="hotspotStyle(link)"
            :aria-label="`${link.name} Google Maps（另開新視窗）`"
            :title="`${link.name}｜開啟 Google Maps`"
          />
        </div>
      </Transition>

      <details ref="menu" class="content-menu">
        <summary><span aria-hidden="true">☰</span> 網站選單</summary>
        <nav aria-label="內頁導覽">
          <RouterLink v-for="link in sectionLinks" :key="link.to" :to="link.to" @click="closeMenu">
            {{ link.label }}
          </RouterLink>
        </nav>
      </details>

      <a
        v-if="showCta && cta.external"
        :class="['content-cta', cta.variant && `content-cta--${cta.variant}`]"
        :href="cta.to"
        :aria-label="cta.label"
        target="_blank"
        rel="noopener noreferrer"
      >
        <img :src="cta.image" :alt="cta.label" />
      </a>
      <RouterLink
        v-else-if="showCta"
        :class="['content-cta', cta.variant && `content-cta--${cta.variant}`]"
        :to="cta.to"
        :aria-label="cta.label"
      >
        <img :src="cta.image" :alt="cta.label" />
      </RouterLink>

      <template v-if="!scrolling && slides.length > 1">
        <button class="content-arrow content-arrow--previous" type="button" aria-label="上一頁" @click="previous">‹</button>
        <button class="content-arrow content-arrow--next" type="button" aria-label="下一頁" @click="next">›</button>

        <nav class="content-pagination" :aria-label="`${title}頁面選擇`">
          <button
            v-for="(slide, index) in slides"
            :key="slide.src"
            type="button"
            :class="{ active: index === current }"
            :aria-label="`前往第 ${index + 1} 頁`"
            :aria-current="index === current ? 'page' : undefined"
            @click="current = index"
          >
            {{ index + 1 }}
          </button>
        </nav>
      </template>
    </section>

    <section :class="['content-mobile', { 'content-mobile--scroll': mobileSlides.length > 1 }]">
      <h1 class="sr-only">{{ title }}</h1>

      <div v-if="mobileSlides.length > 1" class="content-mobile-scroll">
        <figure v-for="(slide, index) in mobileSlides" :key="slide.src" class="content-mobile-scroll__page content-image-stage">
          <img class="content-mobile-artwork" :src="slide.src" :alt="slide.alt" />
          <a
            v-for="link in slide.links"
            :key="`${link.name}-${link.address}`"
            class="content-map-hotspot"
            :href="link.href"
            target="_blank"
            rel="noopener noreferrer"
            :style="hotspotStyle(link)"
            :aria-label="`${link.name} Google Maps（另開新視窗）`"
            :title="`${link.name}｜開啟 Google Maps`"
          />
          <figcaption class="sr-only">{{ title }}手機版第 {{ index + 1 }} 頁</figcaption>
        </figure>
      </div>

      <img
        v-else
        class="content-mobile-artwork"
        :src="mobileSlides[0].src"
        :alt="mobileSlides[0].alt"
      />

      <details ref="mobileMenu" class="content-menu content-menu--mobile">
        <summary><span aria-hidden="true">☰</span> 網站選單</summary>
        <nav aria-label="手機版內頁導覽">
          <RouterLink v-for="link in sectionLinks" :key="link.to" :to="link.to" @click="closeMenu">
            {{ link.label }}
          </RouterLink>
        </nav>
      </details>

      <a
        v-if="showCta && cta.external"
        :class="['content-cta', 'content-cta--mobile', cta.variant && `content-cta--${cta.variant}`]"
        :href="cta.to"
        :aria-label="cta.label"
        target="_blank"
        rel="noopener noreferrer"
      >
        <img :src="cta.mobileImage || cta.image" :alt="cta.label" />
      </a>
      <RouterLink
        v-else-if="showCta"
        :class="['content-cta', 'content-cta--mobile', cta.variant && `content-cta--${cta.variant}`]"
        :to="cta.to"
        :aria-label="cta.label"
      >
        <img :src="cta.mobileImage || cta.image" :alt="cta.label" />
      </RouterLink>
    </section>
  </main>
</template>
