<script setup>
import { reactive, ref } from 'vue'

const form = reactive({ name: '', phone: '', email: '', code: '' })
const loading = ref(false)
const error = ref('')
const success = ref(null)

async function submit() {
  error.value = ''
  loading.value = true
  try {
    const response = await fetch('/api/registrations', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(form),
    })
    const data = await response.json()
    if (!response.ok) throw new Error(data.error || '登記失敗，請稍後再試')
    success.value = data
  } catch (requestError) {
    error.value = requestError.message
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <main class="campaign-page registration-page">
    <section class="campaign-card form-card">
      <template v-if="success">
        <div class="success-mark" aria-hidden="true">✓</div>
        <p class="campaign-kicker">登記完成</p>
        <h1>{{ success.message }}</h1>
        <p class="confirmation-code">您的登錄編號：<strong>{{ success.code }}</strong></p>
        <RouterLink class="primary-action" to="/">返回首頁</RouterLink>
      </template>

      <template v-else>
        <p class="campaign-kicker">海派好禮</p>
        <h1>抽獎資料登記</h1>
        <p class="form-hint"><span aria-hidden="true">＊</span>以下欄位皆為必填，且每組登錄編號僅能使用一次。</p>

        <form class="registration-form" @submit.prevent="submit">
          <label>
            <span>姓名</span>
            <input v-model.trim="form.name" name="name" autocomplete="name" maxlength="100" required placeholder="請輸入姓名" />
          </label>
          <label>
            <span>聯絡電話</span>
            <input v-model.trim="form.phone" name="phone" type="tel" autocomplete="tel" maxlength="30" required placeholder="例：0912345678" />
          </label>
          <label>
            <span>電子信箱</span>
            <input v-model.trim="form.email" name="email" type="email" autocomplete="email" maxlength="254" required placeholder="example@email.com" />
          </label>
          <label>
            <span>登錄編號</span>
            <input v-model.trim="form.code" class="code-input" name="code" maxlength="32" required placeholder="請輸入登錄編號" autocapitalize="characters" />
          </label>
          <p v-if="error" class="form-error" role="alert">{{ error }}</p>
          <button class="primary-action primary-action--button" type="submit" :disabled="loading">
            {{ loading ? '登記中…' : '確認送出' }}
          </button>
        </form>
        <RouterLink class="text-link" to="/gift">返回海派好禮</RouterLink>
      </template>
    </section>
  </main>
</template>
