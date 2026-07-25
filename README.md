# s3Viewer

`s3Viewer` 是一個在 Terminal (命令列) 運作的 AWS S3 / MinIO 檔案管理工具。不用開啟網頁或大型 GUI 軟體，直接在小黑窗裡就能輕鬆瀏覽、上傳和下載檔案。

---

## 🌟 亮點功能

- **支援 AWS S3 & MinIO**：
  在登入頁面切換 `AWS S3` 或 `MinIO`，表單會自動切換為對應需要的輸入框。
- **密碼明碼切換**：
  Username/Password 或 AccessKey/SecretKey 旁邊都有勾選框，勾選就能直接看明碼，不怕輸錯。
- **樹狀目錄瀏覽**：
  資料夾可以點擊展開，方便快速瀏覽裡面的子目錄與檔案。
- **檔案上傳與下載**：
  內建檔案選擇器，選好本機檔案後會自動幫你帶入 `File Key`，不用手動複製檔名。
- **鍵盤操作順手**：
  可以使用 `Tab` 鍵切換按鈕，在輸入框裡也可以用左右方向鍵 (`←` / `→`) 移動游標修改文字。

---

## 🚀 快速開始

### 1. 安裝與執行

確保電腦已安裝 Go (`1.24` 以上)：

```bash
# 取得專案
git clone https://github.com/Nolions/s3Viewer.git
cd s3Viewer
```

#### 🔹 Local developer run
```bash
make run
# 或
go run ./cmd/app/main.go
```

#### 🔹 Build & run
```bash
make compile
# 或
go build -o s3Viewer ./cmd/app
./s3Viewer
```

---

### 2. 登入設定說明

#### 🔹 連線到 AWS S3
1. **Type** 選 `AWS S3`
2. 選擇 **Region** (如 `us-east-1`, `ap-northeast-1`)
3. 填入 **AccessKey** 與 **SecretKey**
4. 填入目標 **Bucket** 名稱

#### 🔹 連線到 MinIO
1. **Type** 選 `MinIO`
2. **Host**：輸入 MinIO 位址 (預設為 `http://localhost:9000`)
3. 填入 **Username** 與 **Password**
4. 填入目標 **Bucket** 名稱

---

## 🛠️ 本地開發指南 (Makefile 指令說明)

專案提供了 `makefile` 方便本地開發、啟動測試服務與跨平台編譯：

| 指令 | 說明 |
| :--- | :--- |
| `make run` | 本地開發直接執行專案 (`go run ./cmd/app/main.go`) |
| `make runMinIO` | 啟動本地測試用 MinIO 容器 (網頁 Console 於 `http://localhost:9001`，預設帳密 `admin` / `admin12345`) |
| `make compile` | 編譯預設平台的執行檔 (輸出至 `build/` 目錄) |
| `make windows-amd64` | 編譯 Windows 64 位元執行檔 |
| `make linux-amd64` | 編譯 Linux 64 位元執行檔 |
| `make linux-arm64` | 編譯 Linux ARM64 執行檔 |
| `make darwin-amd64` | 編譯 macOS Intel 執行檔 |
| `make darwin-arm64` | 編譯 macOS Apple Silicon (M1/M2/M3) 執行檔 |
| `make clean` | 清理編譯產物目錄 (`build/`) |