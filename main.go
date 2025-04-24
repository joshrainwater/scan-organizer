package main

import (
	"bytes"
	"fmt"
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
	trashDir  = "./trash"

	pdfFiles        []string
	currentIndex    int
	previousRenamed []string
	folders         []string
	mu              sync.Mutex
)

func main() {
	os.MkdirAll(imgDir, 0755)
	os.MkdirAll(outputDir, 0755)
	os.MkdirAll(trashDir, 0755)

	loadPDFFiles()
	loadFolders()

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/rename", handleRename)
	http.HandleFunc("/append", handleAppend)
	http.HandleFunc("/trash", handleTrash)
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

func loadFolders() {
	folders = []string{} // Clear existing folders
	filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != outputDir {
			// Get relative path from outputDir
			relPath, err := filepath.Rel(outputDir, path)
			if err != nil {
				return err
			}
			folders = append(folders, relPath)
		}
		return nil
	})
	sort.Strings(folders)
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
		"Preview":         preview,
		"PreviousRenamed": previousRenamed,
		"Folders":         folders,
	})
}

func handleRename(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	newName := r.FormValue("newname")
	folder := r.FormValue("folder")

	// Ensure newName has .pdf extension
	if !strings.HasSuffix(strings.ToLower(newName), ".pdf") {
		newName = newName + ".pdf"
	}

	mu.Lock()
	defer mu.Unlock()
	oldPath := filepath.Join(pdfDir, pdfFiles[currentIndex])

	// Create the full path including the folder
	newPath := filepath.Join(outputDir, folder, newName)
	os.MkdirAll(filepath.Dir(newPath), 0755) // Create folder if it doesn't exist

	os.Rename(oldPath, newPath)
	previousRenamed = append(previousRenamed, newName)

	// Reload folders after rename
	loadFolders()

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

	// Print the full command and paths for debugging
	log.Printf("Attempting to merge PDFs:\nTarget: %s\nCurrent: %s\n", targetPath, current)

	cmd := exec.Command("pdfcpu", "merge", "tmp.pdf", targetPath, current)

	// Capture both stdout and stderr
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Printf("Command failed with error: %v\n", err)
		log.Printf("Stdout: %s\n", stdout.String())
		log.Printf("Stderr: %s\n", stderr.String())
		http.Error(w, fmt.Sprintf("Failed to merge: %v\nStdout: %s\nStderr: %s", err, stdout.String(), stderr.String()), 500)
		return
	}

	// Print success message with command output
	log.Printf("Merge successful. Output:\n%s\n", stdout.String())

	os.Rename("tmp.pdf", targetPath)
	// Remove from list
	pdfFiles = append(pdfFiles[:currentIndex], pdfFiles[currentIndex+1:]...)
	if currentIndex >= len(pdfFiles) && currentIndex > 0 {
		currentIndex--
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleTrash(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	oldPath := filepath.Join(pdfDir, pdfFiles[currentIndex])

	newPath := filepath.Join(trashDir, pdfFiles[currentIndex])
	os.Rename(oldPath, newPath)
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
	<title>Scan Organizer</title>
</head>
<body style="background-color: #eee">
	<div style="display: flex;">
	<div style="flex:none; width: 33%; padding: 2rem">
		<form action="/rename" method="post">
			<input id="folderInput" 
				style="display: block; width: 100%; margin-bottom: 1rem; padding: 0.5rem;" 
				name="folder" 
				type="text" 
				list="folders"
				placeholder="Select or type folder name..."
				autocomplete="off"
				autofocus/>
			<datalist id="folders">
				{{range .Folders}}
				<option value="{{.}}">{{.}}</option>
				{{end}}
			</datalist>
			<input style="display: block; width: 100%; margin-bottom: 1rem; padding: 0.5rem;" 
				type="text" 
				name="newname" 
				placeholder="Rename to..." 
				required>
			<button type="submit">Rename</button>
		</form>
		<hr>
		<form action="/trash" method="post">
			<button id="trashButton" type="submit">Trash</button>
		</form>
		<hr>
		<form action="/append" method="post">
			<select id="targetSelect" name="target">
				{{range .PreviousRenamed}}
				<option value="{{.}}">{{.}}</option>
				{{end}}
			</select>
			<button type="submit">Append to selected</button>
		</form>
	</div>
	<div class="display: block; background-color: white">
		<img src="{{.Preview}}" style="max-height: 50rem; border: 1px solid #ddd">
	</div>
	</div>

	<script>
		document.addEventListener('keydown', function(e) {
			if (e.ctrlKey) {
				switch(e.key) {
					case 'u':
						e.preventDefault();
						document.getElementById('folderInput').focus();
						break;
					case 'i':
						e.preventDefault();
						document.getElementById('trashButton').focus();
						break;
					case 'o':
						e.preventDefault();
						document.getElementById('targetSelect').focus();
						break;
				}
			}
		});
	</script>
</body>
</html>
`
