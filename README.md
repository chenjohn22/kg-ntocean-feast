# 2026 新北海派活動網站

以 Vue 3、Vue Router 與 Vite 製作的響應式活動網站首頁。

## 開發

```bash
npm install
npm run dev
```

## 建置

```bash
npm run build
```

## Docker Compose

```bash
docker compose up --build -d
```

啟動後開啟 <http://localhost:8080>。停止服務可執行：

```bash
docker compose down
```

首頁會依畫面方向切換桌機版（1920×1080）與手機版（1080×1920）視覺。五個首頁入口已分別預留 `/event`、`/map`、`/protect`、`/gift`、`/food` 路由。
