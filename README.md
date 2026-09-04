# libgen-cli

[![Build & Test](https://github.com/anilpdv/libgen-cli/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/anilpdv/libgen-cli)
[![Go Report Card](https://goreportcard.com/badge/github.com/anilpdv/libgen-cli)](https://goreportcard.com/report/github.com/anilpdv/libgen-cli)
[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

A fast, lightweight command-line interface and Go library for searching and downloading books, articles, and research papers from Library Genesis mirrors.

---

## 📑 Table of Contents
- [Features](#features)
- [Installation](#installation)
- [Command Reference](#command-reference)
  - [Search](#search)
  - [Download](#download)
  - [Bulk Download (`download-all`)](#bulk-download-download-all)
  - [Link Resolution](#link-resolution)
  - [Mirror Health Status](#mirror-health-status)
  - [Database Dumps](#database-dumps)
- [Go Library Usage](#go-library-usage)
- [Disclaimer & License](#disclaimer--license)

---

## ⚡ Features

- **Multi-Mirror Search**: Searches across active Library Genesis mirrors (`libgen.is`, `libgen.rs`, `libgen.st`, `libgen.li`).
- **IPFS Gateway Support**: Resolves and downloads via IPFS gateways when web mirrors are slow or unreachable.
- **HTTP Range Resumption**: Automatically resumes interrupted downloads without starting over from scratch.
- **Filtering & Sorting**: Filter and sort results by title, author, publisher, year, format (PDF, EPUB, MOBI), and filesize.
- **Go Library Bindings**: Embed the core search and download engine directly in Go applications via `github.com/anilpdv/libgen-cli/libgen`.

---

## 📦 Installation

### Pre-built Binaries
Download pre-compiled binaries for macOS, Linux, and Windows from the [Releases](https://github.com/anilpdv/libgen-cli/releases) page.

### Go Install
```bash
go install github.com/anilpdv/libgen-cli@latest
```

### Build from Source
```bash
git clone https://github.com/anilpdv/libgen-cli.git
cd libgen-cli
make build
# Binary created at ./bin/libgen-cli
```

---

## 🛠️ Command Reference

### Search
Search Library Genesis by book title, author name, or ISBN:

```bash
# Basic search by title
libgen-cli search "The Go Programming Language"

# Search and resolve download URLs via IPFS mirrors
libgen-cli search "Computer Networking" -i

# Filter results by file format (e.g., EPUB or PDF)
libgen-cli search "Design Patterns" -e "epub,pdf"

# Limit result count (1–100)
libgen-cli search "Introduction to Algorithms" -r 10

# Sort results by publication year (descending)
libgen-cli search "Database Systems" --sort-by year --sort-asc=false

# Filter by publisher and year
libgen-cli search "Operating Systems" -p "Addison-Wesley" -y 2020

# Save search results to a specific directory
libgen-cli search "Calculus" -o ~/Documents/Books
```

### Download
Download a specific book directly by its MD5 hash:

```bash
# Download a single book by MD5
libgen-cli download 2F2DBA2A621B693BB95601C16ED680F8

# Download via IPFS gateway
libgen-cli download 2F2DBA2A621B693BB95601C16ED680F8 --ipfs-mirrors

# Bulk download multiple MD5 hashes
libgen-cli download 6B4B4F0073B92248EFAB34F100CA20D4 FAA323B98939EE385BB33A1A3B88AFCA

# Pipe MD5 hashes from a text file
cat hashes.txt | xargs libgen-cli download -o ~/Downloads/Books
```

### Bulk Download (`download-all`)
Search and automatically download all matching items:

```bash
# Download top 5 results for a query
libgen-cli download-all "Discrete Mathematics" -r 5 -o ~/Downloads/Books

# Download only EPUB versions
libgen-cli download-all "Data Structures" -e epub -o ~/Downloads/Books
```

### Link Resolution
Get the direct download URL for an MD5 hash without downloading:

```bash
# Get direct HTTPS download link
libgen-cli link 2F2DBA2A621B693BB95601C16ED680F8

# Get IPFS gateway download link
libgen-cli link 2F2DBA2A621B693BB95601C16ED680F8 -i
```

### Mirror Health Status
Check connectivity and response time across LibGen search and download mirrors:

```bash
# Check all mirrors
libgen-cli status

# Check search mirrors only
libgen-cli status -m search

# Check download mirrors only
libgen-cli status -m download
```

### Database Dumps
List and download compiled LibGen database dumps:

```bash
libgen-cli dbdumps -o ~/Downloads/LibGenDumps
```

---

## 💻 Go Library Usage

Use the core `libgen` package in your own Go projects:

```go
package main

import (
	"fmt"
	"log"

	"github.com/anilpdv/libgen-cli/libgen"
)

func main() {
	opts := &libgen.SearchOptions{
		Query:        "The C Programming Language",
		SearchMirror: libgen.SearchMirrors[0],
		Results:      5,
		SearchWith:   libgen.Title,
	}

	books, err := libgen.Search(opts)
	if err != nil {
		log.Fatalf("search error: %v", err)
	}

	for _, b := range books {
		fmt.Printf("[%s] %s by %s (%s, %s)\n", b.Id, b.Title, b.Author, b.Extension, b.Filesize)
	}
}
```

---

## ⚖️ Disclaimer
This project is intended for research, educational, and interoperability purposes. Please verify compliance with local copyright laws before downloading copyrighted materials.

---

## 📄 License
Licensed under the [Apache License, Version 2.0](LICENSE).
