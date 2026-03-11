# Scan Organizer

A privacy-focused, local-only desktop application for organizing scanned PDF documents.

## Features

- Preview PDF pages
- Rename and organize PDFs into folders
- Merge/append multiple PDFs
- Drag-and-drop import
- OCR text extraction (coming soon)
- Smart suggestions based on document content (coming soon)

## Philosophy

- **Local-only** - No cloud, no network, no telemetry
- **Privacy-first** - Your documents never leave your machine
- **Simple but powerful** - Easy for anyone, powerful enough for power users

## Tech Stack

- [Wails](https://wails.io/) - Go desktop framework
- React/Vue (frontend - coming soon)
- Tesseract OCR (planned)

## Installation

Download the latest release for your platform from the releases page.

## Development

### Requirements

- Go 1.21+
- Node.js (for frontend)
- pdftoppm (for PDF previews)
- pdfcpu (for PDF merging)

### Running

```bash
go run main.go
```

### Building

```bash
wails build
```

## Folder Structure

```
project-root/
├── input/         # Place PDFs to organize
├── output/        # Organized PDFs
├── static/
│   └── previews/  # Auto-generated previews
├── trash/         # Discarded files
└── main.go        # Application entry point
```

## License

MIT
