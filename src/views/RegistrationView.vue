<script setup>
import { computed, onMounted, reactive, ref, watch } from 'vue'

const form = reactive({
  name: '', phone: '', email: '', code: '', ticketSource: '', restaurantId: '', satisfaction: null, suggestion: '',
})
const loading = ref(false)
const error = ref('')
const success = ref(null)
const restaurants = ref([])
const restaurantsLoading = ref(false)
const restaurantsError = ref('')

const restaurantGroups = computed(() => {
  const groups = new Map()
  restaurants.value.forEach((restaurant) => {
    if (!groups.has(restaurant.category)) groups.set(restaurant.category, [])
    groups.get(restaurant.category).push(restaurant)
  })
  return [...groups.entries()].map(([category, items]) => ({ category, items }))
})

watch(() => form.ticketSource, (source) => {
  if (source !== 'partner_restaurant') form.restaurantId = ''
})

async function loadRestaurants() {
  restaurantsLoading.value = true
  restaurantsError.value = ''
  try {
    const response = await fetch('/api/restaurants')
    const data = await response.json()
    if (!response.ok) throw new Error(data.error || '無法載入活動餐廳')
    restaurants.value = data.items
  } catch (requestError) {
    restaurantsError.value = requestError.message
  } finally {
    restaurantsLoading.value = false
  }
}

async function submit() {
  error.value = ''
  loading.value = true
  try {
    const response = await fetch('/api/registrations', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        name: form.name,
        phone: form.phone,
        email: form.email,
        code: form.code,
        ticket_source: form.ticketSource,
        restaurant_id: form.ticketSource === 'partner_restaurant' ? Number(form.restaurantId) : null,
        satisfaction: form.satisfaction,
        suggestion: form.suggestion,
      }),
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

onMounted(loadRestaurants)
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

          <fieldset class="survey-fieldset">
            <legend>獲得抽獎券的來源</legend>
            <div class="radio-options radio-options--source">
              <label class="radio-card"><input v-model="form.ticketSource" type="radio" name="ticket-source" value="fuji_banquet" required /><span>富基海派宴</span></label>
              <label class="radio-card"><input v-model="form.ticketSource" type="radio" name="ticket-source" value="guihou_fair" /><span>龜吼園遊會</span></label>
              <label class="radio-card"><input v-model="form.ticketSource" type="radio" name="ticket-source" value="partner_restaurant" /><span>活動合作餐廳</span></label>
            </div>
          </fieldset>

          <label v-if="form.ticketSource === 'partner_restaurant'" class="conditional-field">
            <span>活動合作餐廳</span>
            <select v-model="form.restaurantId" name="restaurant" required :disabled="restaurantsLoading || !!restaurantsError">
              <option value="" disabled>{{ restaurantsLoading ? '餐廳載入中…' : '請選擇餐廳' }}</option>
              <optgroup v-for="group in restaurantGroups" :key="group.category" :label="group.category">
                <option v-for="restaurant in group.items" :key="restaurant.id" :value="restaurant.id">
                  {{ restaurant.name }}{{ restaurant.location ? `（${restaurant.location}）` : '' }}
                </option>
              </optgroup>
            </select>
            <small v-if="restaurantsError" class="field-error">{{ restaurantsError }}，請重新整理頁面再試。</small>
          </label>

          <fieldset class="survey-fieldset">
            <legend>對於活動整體滿意度</legend>
            <div class="radio-options radio-options--rating">
              <label v-for="(label, rating) in ['非常不滿意', '不滿意', '普通', '滿意', '非常滿意']" :key="rating" class="radio-card radio-card--rating">
                <input v-model="form.satisfaction" type="radio" name="satisfaction" :value="rating + 1" :required="rating === 0" />
                <span><strong>{{ rating + 1 }}</strong>{{ label }}</span>
              </label>
            </div>
          </fieldset>

          <label>
            <span>對於新北海派活動留下您的寶貴建議</span>
            <textarea v-model.trim="form.suggestion" name="suggestion" rows="5" maxlength="1000" required placeholder="請輸入您的建議（最多 1000 字）" />
            <small class="character-count">{{ form.suggestion.length }} / 1000</small>
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
