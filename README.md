# GoodReadsScrape

**GoodReadsScrape** is a command-line tool written in Go for scraping reviews from the Goodreads platform and generating CSV output. The tool is designed with an efficient architecture using concurrency workers to speed up the scraping process.

## 📋 Description

GoodReadsScrape allows you to automatically collect book reviews from Goodreads. The tool supports scraping multiple URLs simultaneously with a configurable worker pool system, language filtering, and automatic CSV file output with timestamps.

## ✨ Key Features

- **Goodreads Review Scraping**: Extract reviews from Goodreads book pages
- **Concurrency Workers**: Process multiple URLs in parallel with configurable worker count
- **Flexible Input**: Support input from text files (one URL per line) or single URL as argument
- **Language Filter**: Filter reviews by language (default: Indonesian)
- **Max Reviews**: Limit the number of reviews scraped per book
- **CSV Output**: Scraped results saved in structured CSV format

## 🛠 Technologies Used

- **Go 1.25.3**: Primary programming language
- **goquery**: Library for HTML parsing and scraping
- **HTTP Client**: Native Go HTTP client for fetching data
- **CSV Writer**: Encoding/csv for writing CSV output
- **Concurrency**: Goroutines and channels for parallel processing
- **Flag Package**: Native Go flag for command-line arguments

## 📦 Installation

### Prerequisites

- Go 1.25.3 or higher
- Goodreads API key (required for authentication)

### Installation

1. Clone repository:

```bash
git clone https://github.com/rizkirmdhnnn/goodreadscrape.git
cd goodreadscrape
```

2. Build binary:

```bash
go build -o goodreadscrape cmd/bot/main.go
```

3. (Optional) Install to `$GOPATH/bin`:

```bash
go install ./cmd/bot
```

## 🚀 Usage

### Basic Format

```bash
goodreadscrape -api YOUR_API_KEY [OPTIONS] [URL]
```

### Example Commands

**1. Scraping Single URL:**

```bash
goodreadscrape -api YOUR_API_KEY https://www.goodreads.com/book/show/123456
```

**2. Scraping from File with Custom Output:**

```bash
goodreadscrape -api YOUR_API_KEY -f urls.txt -o results/my_reviews.csv
```

**3. Scraping with Full Configuration:**

```bash
goodreadscrape -api YOUR_API_KEY \
  -f urls.txt \
  -c 10 \
  -m 200 \
  -l en \
  -verbose \
  -o results/reviews_english.csv
```

## 📝 Command-Line Flags

| Flag       | Type   | Default | Description                                                               |
| ---------- | ------ | ------- | ------------------------------------------------------------------------- |
| `-api`     | string | -       | **Required.** API key for Goodreads authentication                        |
| `-c`       | int    | 5       | Number of concurrent workers for parallel processing                      |
| `-verbose` | bool   | false   | Enable verbose logging for debugging                                      |
| `-f`       | string | -       | Text file containing Goodreads URLs (one URL per line)                    |
| `-m`       | int    | 100     | Maximum number of reviews to scrape per book                              |
| `-o`       | string | auto    | Output CSV file. Default: `results/goodreads_reviews_YYYYMMDD_HHMMSS.csv` |
| `-l`       | string | "id"    | Language filter for reviews (examples: "id", "en", "es")                  |

### Important Notes

- The `-api` flag is **required** and must be provided
- If not using `-f`, you must provide a URL as a positional argument
- Output file will be automatically created in the `results/` directory if not specified
- If the output file already exists, new data will be appended to it

## 💡 Usage Examples

### Scenario 1: Scraping Single Book

Scrape reviews from a single book with default configuration:

```bash
goodreadscrape -api YOUR_API_KEY https://www.goodreads.com/book/show/123456
```

Output will be saved to `results/goodreads_reviews_20240101_120000.csv`

### Scenario 2: Batch Scraping from File

Prepare `urls.txt` file:

```
https://www.goodreads.com/book/show/123456
https://www.goodreads.com/book/show/789012
https://www.goodreads.com/book/show/345678
```

Run:

```bash
goodreadscrape -api YOUR_API_KEY -f urls.txt -c 10 -verbose
```

### Scenario 3: Scraping with English Language Filter

Scrape English reviews with a maximum of 200 reviews per book:

```bash
goodreadscrape -api YOUR_API_KEY \
  -f urls.txt \
  -l en \
  -m 200 \
  -o results/english_reviews.csv \
  -verbose
```

## 📊 CSV Output Format

The generated CSV file has the following structure:

| Column         | Description                  | Example                                      |
| -------------- | ---------------------------- | -------------------------------------------- |
| `BookURL`      | Full book URL on Goodreads   | `https://www.goodreads.com/book/show/123456` |
| `BookTitle`    | Book title                   | `The Great Gatsby`                           |
| `ReviewID`     | Unique identifier for review | `review_12345`                               |
| `ReviewerName` | Reviewer name                | `John Doe`                                   |
| `Rating`       | Rating given (1-5)           | `5`                                          |
| `ReviewText`   | Full review text             | `This is an amazing book...`                 |
| `ReviewDate`   | Review date                  | `2024-01-15`                                 |
| `Language`     | Review language              | `id` or `en`                                 |

### Example CSV Output

```csv
BookURL,BookTitle,ReviewID,ReviewerName,Rating,ReviewText,ReviewDate,Language
https://www.goodreads.com/book/show/123456,The Great Gatsby,review_1,John Doe,5,"Amazing book!",2024-01-15,en
https://www.goodreads.com/book/show/123456,The Great Gatsby,review_2,Jane Smith,4,"Good read",2024-01-16,en
```

---

**Note**: This tool is created for educational and research purposes. Please ensure you comply with Goodreads Terms of Service and use this tool responsibly.
