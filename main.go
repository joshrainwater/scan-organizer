package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

var (
	pdfDir    = "./input"
	imgDir    = "./static/previews"
	outputDir = "./output"

	pdfFiles     []string
	currentIndex int
	previousRenamed []string
	mu           sync.Mutex
)

func main() {
	os.MkdirAll(imgDir, 0755)
	os.MkdirAll(outputDir, 0755)
	
	loadPDFFiles()
	
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/next", handleNext)
	http.HandleFunc("/prev", handlePrev)
	http.HandleFunc("/rename", handleRename)
	http.HandleFunc("/append", handleAppend)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func loadPDFFiles() {
	files, _ := os.ReadDir(pdfDir)
	for _, file := range files {
		if strings.HasSuffix(strings.ToLower(file.Name()), ".pdf") {
			pdfFiles = append(pdfFiles, file.Name())
		}
	}
	sort.Strings(pdfFiles)
}

func renderPreview(pdf string) string {
	img := strings.TrimSuffix(pdf, ".pdf") + ".png"
	imgPath := filepath.Join(imgDir, img)
	if _, err := os.Stat(imgPath); os.IsNotExist(err) {
		cmd := exec.Command("pdftoppm", filepath.Join(pdfDir, pdf), filepath.Join(imgDir, strings.TrimSuffix(pdf, ".pdf")), "-png", "-singlefile")
		cmd.Run()
	}
	return "/static/previews/" + img
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	if len(pdfFiles) == 0 {
		http.Error(w, "No PDFs found", 404)
		return
	}
	current := pdfFiles[currentIndex]
	preview := renderPreview(current)
	tmpl := template.Must(template.New("index").Parse(indexHTML))
	tmpl.Execute(w, map[string]interface{}{
		"Current": current,
		"Preview": preview,
		"PreviousRenamed": previousRenamed,
	})
}

func handleNext(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	if currentIndex < len(pdfFiles)-1 {
		currentIndex++
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handlePrev(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	if currentIndex > 0 {
		currentIndex--
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleRename(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	newName := r.FormValue("newname")
	mu.Lock()
	defer mu.Unlock()
	oldPath := filepath.Join(pdfDir, pdfFiles[currentIndex])
	newPath := filepath.Join(outputDir, newName)
	os.Rename(oldPath, newPath)
	previousRenamed = append(previousRenamed, newName)
	// Remove from list
	pdfFiles = append(pdfFiles[:currentIndex], pdfFiles[currentIndex+1:]...)
	if currentIndex >= len(pdfFiles) && currentIndex > 0 {
		currentIndex--
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleAppend(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	target := r.FormValue("target")
	mu.Lock()
	defer mu.Unlock()
	current := filepath.Join(pdfDir, pdfFiles[currentIndex])
	targetPath := filepath.Join(outputDir, target)
	cmd := exec.Command("pdfcpu", "merge", "tmp.pdf", targetPath, current)
	err := cmd.Run()
	if err != nil {
		http.Error(w, "Failed to merge", 500)
		return
	}
	os.Rename("tmp.pdf", targetPath)
	// Remove from list
	pdfFiles = append(pdfFiles[:currentIndex], pdfFiles[currentIndex+1:]...)
	if currentIndex >= len(pdfFiles) && currentIndex > 0 {
		currentIndex--
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

const indexHTML = `
<!DOCTYPE html>
<html>
<head>
	<title>PDF Tool</title>
</head>
<body>
	<h1>Viewing: {{.Current}}</h1>
	<img src="{{.Preview}}" style="max-width:100%; height:auto;">
	<form action="/rename" method="post">
		<input type="text" name="newname" placeholder="Rename to..." required>
		<button type="submit">Rename</button>
	</form>
	<br>
	<form action="/append" method="post">
		<select name="target">
			{{range .PreviousRenamed}}
			<option value="{{.}}">{{.}}</option>
			{{end}}
		</select>
		<button type="submit">Append to selected</button>
	</form>
	<br>
	<a href="/prev">⬅ Prev</a>
	<a href="/next">Next ➡</a>
</body>
</html>
`
