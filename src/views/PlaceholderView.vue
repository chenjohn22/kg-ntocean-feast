<script setup>
defineProps({
  title: {
    type: String,
    required: true,
  },
})
</script>

<template>
  <main class="placeholder-page">
    <section class="placeholder-card">
      <p class="placeholder-kicker">2026 新北海派</p>
      <h1>{{ title }}</h1>
      <p>此頁面已預留路由，內容可於下一階段接續製作。</p>
      <RouterLink class="back-home" to="/">返回首頁</RouterLink>
    </section>
  </main>
</template>
