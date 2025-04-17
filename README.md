# PDF Manager Web Tool

This is a lightweight local web-based tool written in Go that helps you manage a folder of PDF files. It allows you to:

- Preview the first page of each PDF
- Rename PDFs
- Append a PDF to the end of a previously renamed file (merging)
- Iterate through PDFs one-by-one in a browser interface

---

## 🔧 Requirements

Make sure you have the following installed:

### 1. **Go** (1.18 or later)

Install from: [https://golang.org/dl/](https://golang.org/dl/)

### 2. **pdftoppm**

Used to generate a PNG preview of the first page.

Install via:
```bash
sudo apt install poppler-utils
```

### 3. **pdfcpu**

Used to merge PDF files.

Install via:

```bash
go install github.com/pdfcpu/pdfcpu/cmd/pdfcpu@latest
```

Make sure `pdfcpu` is available in your `$PATH`:

```bash
pdfcpu version
```

Note that the default directory for go binaries is $HOME/go/bin

---

## 📂 Folder Structure

```
project-root/
├── input/         # Place your PDFs here
├── output/        # Renamed and merged PDFs go here
├── static/
│   └── previews/  # Auto-generated PNG previews
├── main.go        # The Go web app
```

---

## 🚀 Running the Tool

From the project root directory, run:

```bash
go run main.go
```

Then open your browser and go to:

```
http://localhost:8080
```

---

## 🖱 How to Use

1. Place all your PDFs into the `input/` folder.
2. Start the server and open the web UI.
3. For each file:
   - View a preview of the first page.
   - Use the input field to rename the file. It will be moved to the `output/` folder.
   - Or, use the dropdown to append it to a previously renamed file.
4. Use the **Next** / **Prev** buttons to navigate.

---

## 📌 Notes

- Only PDFs are recognized in the `input/` folder.
- Preview is auto-generated for the first page only.
- Renamed files are stored in `output/`.
- Merging appends the entire current PDF to the selected file.

---

## 📄 License

MIT License

---

## 🙏 Acknowledgements

- [pdfcpu](https://github.com/pdfcpu/pdfcpu)
- [Poppler Utils](https://poppler.freedesktop.org/)


