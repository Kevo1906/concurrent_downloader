
# 🚀 Concurrent File Downloader CLI (Go)

A blazing-fast, concurrent file downloader written in Go, designed for software engineers and power users who need efficient and reliable downloading from multiple URLs.

## ✨ Features

- 🧠 Concurrent downloads with configurable workers
- 🧰 Retry mechanism with exponential backoff
- 📂 Supports input via URL list or file
- 📊 Real-time progress bar for each file
- ⏱️ Context-aware cancellation (graceful interrupt)
- 💬 Clean, informative logging

## 🛠️ Installation

```bash
git clone https://github.com/yourusername/concurrent-downloader.git
cd concurrent-downloader
go build -o downloader
```

## 📦 Usage

### Download from a list of URLs in a file

```bash
./downloader --file=urls.txt --output=./downloads --workers=5
```

### Download from command-line argument

```bash
./downloader --urls="http://example.com/file1.jpg,http://example.com/file2.jpg" --output=./downloads
```

## 📁 Example `urls.txt`

```
https://example.com/file1.jpg
https://example.com/file2.jpg
```

## ⚙️ Flags

| Flag       | Description                                 |
|------------|---------------------------------------------|
| `--file`   | Path to file containing URLs (one per line) |
| `--urls`   | Comma-separated URLs                        |
| `--output` | Directory to save downloaded files          |
| `--workers`| Number of concurrent downloads (default: 5) |

## 🧹 .gitignore

```gitignore
/downloaded_files/
*.exe
*.log
*.tmp
```

## 🧑‍💻 Author

Crafted with care by Kevin

---