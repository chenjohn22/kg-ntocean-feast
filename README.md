# 2026 新北海派活動網站

Vue 3 活動網站，加上 Go / MySQL 的海派好禮登記系統與管理後台。

## 功能

- `/gift`：海派好禮入口
- `/gift/register`：顧客填寫基本資料、抽獎券來源、合作餐廳、滿意度及活動建議
- `/kg-manager-admin`：序號管理後台（預設帳密 `admin` / `admin`）
- 60,000 筆序號匯入、資料庫唯一限制防止重複登記
- 後台支援姓名、電話、信箱、登錄編號、日期與狀態篩選
- 分頁、滿意度圖表、可拖曳檢視的建議欄位與依目前搜尋條件匯出 UTF-8 CSV
- 每日自動備份序號與登記資料，僅保留最近 72 小時

## 一鍵啟動（Docker Compose）

```bash
docker compose up --build -d
```

啟動過程會建立 MySQL database、執行 migration，並將
`random_codes_60000.json` 的 60,000 筆序號及 `backend/data/activity_restaurants.json`
的活動餐廳，以可重複執行的方式匯入資料庫。

完成後開啟：

- 活動網站：<http://localhost:8080>
- 登記頁：<http://localhost:8080/gift/register>
- 管理後台：<http://localhost:8080/kg-manager-admin>

停止服務：

```bash
docker compose down
```

資料儲存在 `mysql-data` volume，一般的 `docker compose down` 不會刪除資料。
正式環境務必透過環境變數更換 `ADMIN_PASSWORD`，HTTPS 環境並設定
`ADMIN_COOKIE_SECURE=true`。

## 資料庫備份與還原

`backup` 服務預設為關閉；部署到 server 後，使用下列指令啟用：

```bash
docker compose --profile backup up -d
```

啟用後會在台北時間每天凌晨 03:00 備份 `lottery_codes`、
`activity_restaurants` 與 `registrations`，檔案存放在專案的 `backups/` 目錄。每次成功備份後會刪除
超過 72 小時的 `.sql.gz`；若備份失敗，不會清除任何舊備份。

若服務在當日 03:00 後才啟動，且當天尚無備份，會立即補做一次。可在
`.env` 調整執行時間與保存時間：

```dotenv
BACKUP_HOUR=3
BACKUP_RETENTION_HOURS=72
```

手動立即備份：

```bash
docker compose --profile backup exec backup /bin/sh /opt/backup/backup-db.sh
```

查看備份：

```bash
ls -lh backups/
```

還原指定備份前，請先確認這會覆蓋同名資料表中的現有資料：

```bash
gunzip -c backups/ntocean-feast-YYYYMMDD-HHMMSS.sql.gz | \
  docker compose exec -T db mysql -u root -p ntocean_feast
```

## 本機開發

預設資料庫設定為：

```dotenv
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=rootpassword
DB_NAME=ntocean_feast
```

請先建立空的 database（Compose 使用者不需手動建立）：

```bash
mysql -h 127.0.0.1 -P 3306 -u root -p -e \
  'CREATE DATABASE IF NOT EXISTS ntocean_feast CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci'
```

執行 schema migration 與序號匯入：

```bash
make migrate
```

分別啟動 Go API 與 Vite 開發伺服器：

```bash
npm install
make dev-api
# 另一個終端機
make dev-web
```

Vite 會將 `/api` 代理到 `127.0.0.1:8081`。

## 建置

```bash
make build
```

## 設定

可參考 `.env.example`。主要環境變數如下：

- `DB_HOST`、`DB_PORT`、`DB_USER`、`DB_PASSWORD`、`DB_NAME`
- `APP_ADDR`（API 預設 `:8081`）
- `ADMIN_USER`、`ADMIN_PASSWORD`
- `ADMIN_COOKIE_SECURE`（使用 HTTPS 時設為 `true`）
- `BACKUP_HOUR`（台北時間每日執行小時，預設 `3`）
- `BACKUP_RETENTION_HOURS`（保存時數，預設 `72`）

Migration 檔案位於 `backend/migrations`。回滾最新一版可執行
`make migrate-down`；回滾會移除最新 migration 建立的欄位或資料表，執行前應先備份。

## Demo Server 部署

Server 使用 `compose.server.yaml` 加上 Caddy 提供 HTTPS。Server 的 `.env` 請將
`WEB_BIND_ADDRESS` 設為 `127.0.0.1`，避免略過 HTTPS 直接連入 Web 容器：

```bash
docker compose -f compose.yaml -f compose.server.yaml --profile backup up --build -d
```

正式網址 `www.ntocean-feast.com.tw` 目前暫停對外服務；測試網址為
<https://demo.ntocean-feast.com.tw>。正式上線時將 `www` 加回 `Caddyfile`，
並確認兩個 DNS 的 A 記錄都指向 Server。

首頁會依畫面方向切換桌機版（1920×1080）與手機版（1080×1920）視覺。五個首頁入口已分別預留 `/event`、`/map`、`/protect`、`/gift`、`/food` 路由。
