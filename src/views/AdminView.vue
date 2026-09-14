<script setup>
import { computed, onMounted, reactive, ref } from 'vue'

const authChecked = ref(false)
const loggedIn = ref(false)
const loginForm = reactive({ username: 'admin', password: '' })
const loginError = ref('')
const loading = ref(false)
const rows = ref([])
const summary = reactive({ total: 0, registered: 0, available: 0 })
const pagination = reactive({ page: 1, pageSize: 20, total: 0, totalPages: 0 })
const filters = reactive({ q: '', name: '', phone: '', email: '', code: '', dateFrom: '', dateTo: '', status: '' })
const dataError = ref('')
const activePanel = ref('registrations')
const selectedSuggestion = ref('')
const satisfaction = reactive({ total: 0, average: 0, ratings: [] })

const ratingLabels = ['非常不滿意', '不滿意', '普通', '滿意', '非常滿意']

const startRow = computed(() => pagination.total ? (pagination.page - 1) * pagination.pageSize + 1 : 0)
const endRow = computed(() => Math.min(pagination.page * pagination.pageSize, pagination.total))

function queryString(includePaging = true) {
  const params = new URLSearchParams()
  const mapping = { q: 'q', name: 'name', phone: 'phone', email: 'email', code: 'code', dateFrom: 'date_from', dateTo: 'date_to', status: 'status' }
  Object.entries(mapping).forEach(([key, param]) => {
    if (filters[key]) params.set(param, filters[key])
  })
  if (includePaging) {
    params.set('page', pagination.page)
    params.set('page_size', pagination.pageSize)
  }
  return params.toString()
}

async function api(url, options) {
  const response = await fetch(url, options)
  let data = null
  if (response.status !== 204) data = await response.json()
  if (response.status === 401) {
    loggedIn.value = false
    throw new Error(data?.error || '登入已逾時')
  }
  if (!response.ok) throw new Error(data?.error || '操作失敗')
  return data
}

async function checkSession() {
  try {
    await api('/api/admin/session')
    loggedIn.value = true
    await loadDashboard()
  } catch {
    loggedIn.value = false
  } finally {
    authChecked.value = true
  }
}

async function login() {
  loginError.value = ''
  loading.value = true
  try {
    await api('/api/admin/login', {
      method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(loginForm),
    })
    loggedIn.value = true
    loginForm.password = ''
    await loadDashboard()
  } catch (error) {
    loginError.value = error.message
  } finally {
    loading.value = false
  }
}

async function logout() {
  try { await api('/api/admin/logout', { method: 'POST' }) } catch { /* session is already gone */ }
  loggedIn.value = false
}

async function loadDashboard() {
  loading.value = true
  dataError.value = ''
  try {
    const [summaryData, listData, satisfactionData] = await Promise.all([
      api('/api/admin/summary'),
      api(`/api/admin/codes?${queryString()}`),
      api('/api/admin/satisfaction'),
    ])
    Object.assign(summary, summaryData)
    Object.assign(satisfaction, satisfactionData)
    rows.value = listData.items
    pagination.total = listData.total
    pagination.totalPages = listData.total_pages
    if (pagination.page > Math.max(1, pagination.totalPages)) {
      pagination.page = Math.max(1, pagination.totalPages)
      return loadDashboard()
    }
  } catch (error) {
    dataError.value = error.message
  } finally {
    loading.value = false
  }
}

function search() {
  pagination.page = 1
  loadDashboard()
}

function resetFilters() {
  Object.assign(filters, { q: '', name: '', phone: '', email: '', code: '', dateFrom: '', dateTo: '', status: '' })
  pagination.page = 1
  loadDashboard()
}

function changePage(page) {
  if (page < 1 || page > pagination.totalPages || page === pagination.page) return
  pagination.page = page
  loadDashboard()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

function formatDate(value) {
  if (!value) return '—'
  return new Intl.DateTimeFormat('zh-TW', { dateStyle: 'medium', timeStyle: 'short' }).format(new Date(value))
}

function exportCsv() {
  window.location.href = `/api/admin/export.csv?${queryString(false)}`
}

function sourceLabel(source) {
  return {
    fuji_banquet: '富基海派宴',
    guihou_fair: '龜吼園遊會',
    partner_restaurant: '活動合作餐廳',
  }[source] || '—'
}

function showSuggestion(suggestion) {
  selectedSuggestion.value = suggestion
}

onMounted(checkSession)
</script>

<template>
  <main class="admin-page">
    <section v-if="!authChecked" class="admin-loading">後台載入中…</section>

    <section v-else-if="!loggedIn" class="login-card">
      <p class="admin-brand">KG MANAGER</p>
      <h1>管理後台登入</h1>
      <form @submit.prevent="login">
        <label><span>帳號</span><input v-model="loginForm.username" autocomplete="username" required /></label>
        <label><span>密碼</span><input v-model="loginForm.password" type="password" autocomplete="current-password" required autofocus /></label>
        <p v-if="loginError" class="form-error" role="alert">{{ loginError }}</p>
        <button class="admin-button admin-button--primary" type="submit" :disabled="loading">{{ loading ? '登入中…' : '登入' }}</button>
      </form>
      <RouterLink class="admin-home-link" to="/">回活動首頁</RouterLink>
    </section>

    <template v-else>
      <header class="admin-header">
        <div><p class="admin-brand">KG MANAGER</p><h1>海派好禮登記管理</h1></div>
        <button class="admin-button" type="button" @click="logout">登出</button>
      </header>

      <div class="admin-content">
        <nav class="admin-tabs" aria-label="後台頁面">
          <button type="button" :class="{ active: activePanel === 'registrations' }" @click="activePanel = 'registrations'">登記資料</button>
          <button type="button" :class="{ active: activePanel === 'satisfaction' }" @click="activePanel = 'satisfaction'">滿意度統計</button>
        </nav>

        <template v-if="activePanel === 'registrations'">
          <section class="stat-grid" aria-label="登記統計">
            <article><span>全部序號</span><strong>{{ summary.total.toLocaleString() }}</strong></article>
            <article><span>已登記</span><strong>{{ summary.registered.toLocaleString() }}</strong></article>
            <article><span>未登記</span><strong>{{ summary.available.toLocaleString() }}</strong></article>
          </section>

          <section class="filter-panel">
            <form @submit.prevent="search">
              <label class="filter-wide"><span>快速搜尋</span><input v-model.trim="filters.q" placeholder="姓名、電話、信箱或登錄編號" /></label>
              <label><span>姓名</span><input v-model.trim="filters.name" /></label>
              <label><span>聯絡電話</span><input v-model.trim="filters.phone" /></label>
              <label><span>信箱</span><input v-model.trim="filters.email" /></label>
              <label><span>登錄編號</span><input v-model.trim="filters.code" /></label>
              <label><span>開始日期</span><input v-model="filters.dateFrom" type="date" /></label>
              <label><span>結束日期</span><input v-model="filters.dateTo" type="date" /></label>
              <label><span>登記狀態</span><select v-model="filters.status"><option value="">全部</option><option value="registered">已登記</option><option value="available">未登記</option></select></label>
              <div class="filter-actions">
                <button class="admin-button admin-button--primary" type="submit">搜尋</button>
                <button class="admin-button" type="button" @click="resetFilters">清除</button>
              </div>
            </form>
          </section>

          <section class="table-panel">
            <div class="table-toolbar">
              <p>顯示第 {{ startRow.toLocaleString() }}–{{ endRow.toLocaleString() }} 筆，共 {{ pagination.total.toLocaleString() }} 筆</p>
              <button class="admin-button admin-button--export" type="button" @click="exportCsv">匯出 CSV</button>
            </div>
            <p v-if="dataError" class="form-error">{{ dataError }}</p>
            <div class="table-scroll" :class="{ 'is-loading': loading }">
              <table class="registration-table">
                <thead><tr><th>登錄編號</th><th>狀態</th><th>姓名</th><th>聯絡電話</th><th>信箱</th><th>抽獎券來源</th><th>活動合作餐廳</th><th>滿意度</th><th>活動建議</th><th>登記日期</th></tr></thead>
                <tbody>
                  <tr v-for="row in rows" :key="row.id">
                    <td class="code-cell">{{ row.code }}</td>
                    <td><span :class="['status-pill', row.registered ? 'is-registered' : 'is-available']">{{ row.registered ? '已登記' : '未登記' }}</span></td>
                    <td>{{ row.name || '—' }}</td>
                    <td>{{ row.phone || '—' }}</td>
                    <td>{{ row.email || '—' }}</td>
                    <td>{{ row.registered ? sourceLabel(row.ticket_source) : '—' }}</td>
                    <td>{{ row.restaurant_name || '—' }}</td>
                    <td>{{ row.satisfaction ? `${row.satisfaction} 分` : '—' }}</td>
                    <td><button v-if="row.suggestion" class="suggestion-button" type="button" @click="showSuggestion(row.suggestion)">查看建議</button><span v-else>—</span></td>
                    <td>{{ formatDate(row.registered_at) }}</td>
                  </tr>
                  <tr v-if="!loading && !rows.length"><td colspan="10" class="empty-cell">查無符合條件的資料</td></tr>
                </tbody>
              </table>
            </div>
            <nav v-if="pagination.totalPages > 1" class="pagination" aria-label="分頁">
              <button :disabled="pagination.page === 1" @click="changePage(pagination.page - 1)">上一頁</button>
              <span>第 {{ pagination.page.toLocaleString() }} / {{ pagination.totalPages.toLocaleString() }} 頁</span>
              <button :disabled="pagination.page === pagination.totalPages" @click="changePage(pagination.page + 1)">下一頁</button>
              <select v-model.number="pagination.pageSize" aria-label="每頁筆數" @change="pagination.page = 1; loadDashboard()"><option :value="20">20 筆/頁</option><option :value="50">50 筆/頁</option><option :value="100">100 筆/頁</option></select>
            </nav>
          </section>
        </template>

        <section v-else class="satisfaction-page" aria-labelledby="satisfaction-title">
          <div class="satisfaction-heading">
            <div><p class="admin-brand">SURVEY OVERVIEW</p><h2 id="satisfaction-title">活動整體滿意度</h2></div>
            <p>統計已完成登記且填寫問卷的資料</p>
          </div>
          <div class="satisfaction-kpis">
            <article><span>平均滿意度</span><strong>{{ satisfaction.average.toFixed(2) }}</strong><small>/ 5 分</small></article>
            <article><span>問卷回覆數</span><strong>{{ satisfaction.total.toLocaleString() }}</strong><small>份</small></article>
          </div>
          <div class="rating-chart" role="img" aria-label="活動整體滿意度 1 至 5 分統計長條圖">
            <article v-for="item in satisfaction.ratings" :key="item.rating" class="rating-row">
              <div class="rating-label"><strong>{{ item.rating }} 分</strong><span>{{ ratingLabels[item.rating - 1] }}</span></div>
              <div class="rating-track"><div class="rating-bar" :style="{ width: `${item.percentage}%` }" /></div>
              <div class="rating-value"><strong>{{ item.count.toLocaleString() }}</strong><span>{{ item.percentage.toFixed(1) }}%</span></div>
            </article>
            <p v-if="!satisfaction.total" class="empty-chart">目前尚無滿意度資料</p>
          </div>
        </section>
      </div>

      <div v-if="selectedSuggestion" class="modal-backdrop" role="presentation" @click.self="selectedSuggestion = ''">
        <section class="suggestion-modal" role="dialog" aria-modal="true" aria-labelledby="suggestion-title">
          <div class="modal-header"><h2 id="suggestion-title">活動建議內容</h2><button type="button" aria-label="關閉" @click="selectedSuggestion = ''">×</button></div>
          <p>{{ selectedSuggestion }}</p>
          <button class="admin-button admin-button--primary" type="button" @click="selectedSuggestion = ''">關閉</button>
        </section>
      </div>
    </template>
  </main>
</template>
