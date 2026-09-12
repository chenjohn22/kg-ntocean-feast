<script setup>
import { computed, ref, watch } from 'vue'

const props = defineProps({
  title: { type: String, required: true },
  slides: { type: Array, required: true },
  cta: { type: Object, default: null },
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
const activeSlide = computed(() => props.slides[current.value])

watch(() => props.title, () => {
  current.value = 0
})

function previous() {
  current.value = (current.value - 1 + props.slides.length) % props.slides.length
}

function next() {
  current.value = (current.value + 1) % props.slides.length
}

function handleKeydown(event) {
  if (props.slides.length < 2) return
  if (event.key === 'ArrowLeft') previous()
  if (event.key === 'ArrowRight') next()
}

function closeMenu() {
  menu.value?.removeAttribute('open')
}
</script>

<template>
  <main class="content-page" :aria-label="title" tabindex="-1" @keydown="handleKeydown">
    <section class="content-desktop">
      <h1 class="sr-only">{{ title }}</h1>

      <Transition name="content-slide" mode="out-in">
        <img
          :key="activeSlide.src"
          class="content-artwork"
          :src="activeSlide.src"
          :alt="activeSlide.alt"
        />
      </Transition>

      <details ref="menu" class="content-menu">
        <summary><span aria-hidden="true">☰</span> 網站選單</summary>
        <nav aria-label="內頁導覽">
          <RouterLink v-for="link in sectionLinks" :key="link.to" :to="link.to" @click="closeMenu">
            {{ link.label }}
          </RouterLink>
        </nav>
      </details>

      <RouterLink v-if="cta" class="content-cta" :to="cta.to" :aria-label="cta.label">
        <img :src="cta.image" :alt="cta.label" />
      </RouterLink>

      <template v-if="slides.length > 1">
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

    <section class="content-mobile-pending">
      <p class="content-mobile-kicker">2026 新北海派</p>
      <h1>{{ title }}</h1>
      <p>手機版視覺素材準備中，請先使用電腦版瀏覽完整內容。</p>
      <nav class="content-mobile-links" aria-label="內頁導覽">
        <RouterLink v-for="link in sectionLinks" :key="link.to" :to="link.to">
          {{ link.label }}
        </RouterLink>
      </nav>
    </section>
  </main>
</template>
