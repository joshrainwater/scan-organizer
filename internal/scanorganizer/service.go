package scanorganizer

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

type Service struct {
	stagingDir string
	inputDir   string
	outputDir  string
	previewDir string

	mu              sync.Mutex
	inputFiles      []string
	outputFiles     []string
	currentIndex    int
	previousRenamed []string
	folders         []string
}

type PreviewData struct {
	Preview         string   `json:"preview"`
	PreviousRenamed []string `json:"previousRenamed"`
	Folders         []string `json:"folders"`
}

type StatusData struct {
	InputCount       int  `json:"inputCount"`
	OutputCount      int  `json:"outputCount"`
	IsFullyProcessed bool `json:"isFullyProcessed"`
}

func NewService(stagingDir string) (*Service, error) {
	s := &Service{
		stagingDir: stagingDir,
		inputDir:   filepath.Join(stagingDir, "input"),
		outputDir:  filepath.Join(stagingDir, "output"),
		previewDir: filepath.Join(stagingDir, "previews"),
	}

	if err := os.MkdirAll(s.inputDir, 0755); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(s.outputDir, 0755); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(s.previewDir, 0755); err != nil {
		return nil, err
	}

	if err := s.loadInputFiles(); err != nil {
		return nil, err
	}
	if err := s.loadOutputFiles(); err != nil {
		return nil, err
	}
	if err := s.loadFolders(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Service) loadInputFiles() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.inputFiles = nil

	files, err := os.ReadDir(s.inputDir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if strings.HasSuffix(strings.ToLower(file.Name()), ".pdf") {
			s.inputFiles = append(s.inputFiles, file.Name())
		}
	}
	sort.Strings(s.inputFiles)
	if s.currentIndex >= len(s.inputFiles) {
		s.currentIndex = 0
	}
	return nil
}

func (s *Service) loadOutputFiles() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.outputFiles = nil

	files, err := os.ReadDir(s.outputDir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if strings.HasSuffix(strings.ToLower(file.Name()), ".pdf") {
			s.outputFiles = append(s.outputFiles, file.Name())
		}
	}
	sort.Strings(s.outputFiles)
	return nil
}

func (s *Service) loadFolders() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.folders = []string{}
	err := filepath.Walk(s.outputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != s.outputDir {
			relPath, err := filepath.Rel(s.outputDir, path)
			if err != nil {
				return err
			}
			s.folders = append(s.folders, relPath)
		}
		return nil
	})
	if err != nil {
		return err
	}
	sort.Strings(s.folders)
	return nil
}

func (s *Service) AddFiles(paths []string, mode string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch mode {
	case "replace":
		entries, err := os.ReadDir(s.inputDir)
		if err != nil {
			return err
		}
		for _, entry := range entries {
			if err := os.RemoveAll(filepath.Join(s.inputDir, entry.Name())); err != nil {
				return err
			}
		}
	case "append":
	case "skip":
	default:
		return fmt.Errorf("invalid mode: %s (use: replace, append, skip)", mode)
	}

	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			log.Printf("failed to stat %s: %v", path, err)
			continue
		}

		if info.IsDir() {
			if err := s.copyDirectoryContents(path, s.inputDir); err != nil {
				log.Printf("failed to copy directory %s: %v", path, err)
			}
		} else {
			if err := s.copyFile(path, s.inputDir); err != nil {
				log.Printf("failed to copy file %s: %v", path, err)
			}
		}
	}

	return nil
}

func (s *Service) copyFile(src string, destDir string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !strings.HasSuffix(strings.ToLower(srcInfo.Name()), ".pdf") {
		return nil
	}

	dest := filepath.Join(destDir, srcInfo.Name())
	return copyFileContents(src, dest)
}

func copyFileContents(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	sourceInfo, err := sourceFile.Stat()
	if err != nil {
		return err
	}

	destFile, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, sourceInfo.Mode())
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

func (s *Service) copyDirectoryContents(srcDir, destDir string) error {
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(srcDir, entry.Name())
		if entry.IsDir() {
			if err := s.copyDirectoryContents(srcPath, destDir); err != nil {
				log.Printf("failed to copy subdir %s: %v", srcPath, err)
			}
		} else {
			if err := s.copyFile(srcPath, destDir); err != nil {
				log.Printf("failed to copy %s: %v", srcPath, err)
			}
		}
	}
	return nil
}

func (s *Service) Export(destination string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := os.MkdirAll(destination, 0755); err != nil {
		return err
	}

	entries, err := os.ReadDir(s.outputDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		src := filepath.Join(s.outputDir, entry.Name())
		dest := filepath.Join(destination, entry.Name())

		counter := 1
		for {
			if _, err := os.Stat(dest); os.IsNotExist(err) {
				break
			}
			ext := filepath.Ext(entry.Name())
			baseName := strings.TrimSuffix(entry.Name(), ext)
			dest = fmt.Sprintf("%s (%d)%s", baseName, counter, ext)
			counter++
		}

		if err := os.Rename(src, dest); err != nil {
			return err
		}
	}

	s.outputFiles = nil
	return nil
}

func (s *Service) GetStatus() (*StatusData, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return &StatusData{
		InputCount:       len(s.inputFiles),
		OutputCount:      len(s.outputFiles),
		IsFullyProcessed: len(s.inputFiles) == 0 && len(s.outputFiles) > 0,
	}, nil
}

func (s *Service) renderPreview(pdf string) (string, error) {
	img := strings.TrimSuffix(pdf, ".pdf") + ".png"
	imgPath := filepath.Join(s.previewDir, img)
	if _, err := os.Stat(imgPath); os.IsNotExist(err) {
		cmd := exec.Command("pdftoppm", filepath.Join(s.inputDir, pdf), filepath.Join(s.previewDir, strings.TrimSuffix(pdf, ".pdf")), "-png", "-singlefile")
		_ = cmd.Run()
	}

	data, err := os.ReadFile(imgPath)
	if err != nil {
		return "", err
	}
	encoded := base64.StdEncoding.EncodeToString(data)
	return "data:image/png;base64," + encoded, nil
}

func (s *Service) GetPreview() (*PreviewData, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.inputFiles) == 0 {
		return nil, fmt.Errorf("no PDFs found in %s", s.inputDir)
	}

	current := s.inputFiles[s.currentIndex]
	preview, err := s.renderPreview(current)
	if err != nil {
		return nil, err
	}

	return &PreviewData{
		Preview:         preview,
		PreviousRenamed: append([]string(nil), s.previousRenamed...),
		Folders:         append([]string(nil), s.folders...),
	}, nil
}

func (s *Service) Rename(newName, folder string) error {
	if newName == "" {
		return fmt.Errorf("new name is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.inputFiles) == 0 {
		return fmt.Errorf("no PDFs to rename")
	}

	if !strings.HasSuffix(strings.ToLower(newName), ".pdf") {
		newName = newName + ".pdf"
	}

	oldPath := filepath.Join(s.inputDir, s.inputFiles[s.currentIndex])
	newPath := filepath.Join(s.outputDir, folder, newName)
	if err := os.MkdirAll(filepath.Dir(newPath), 0755); err != nil {
		return err
	}

	if err := os.Rename(oldPath, newPath); err != nil {
		return err
	}

	s.previousRenamed = append(s.previousRenamed, newName)

	if err := s.loadFolders(); err != nil {
		log.Printf("failed to reload folders: %v", err)
	}

	s.inputFiles = append(s.inputFiles[:s.currentIndex], s.inputFiles[s.currentIndex+1:]...)
	if s.currentIndex >= len(s.inputFiles) && s.currentIndex > 0 {
		s.currentIndex--
	}

	return nil
}

func (s *Service) Append(target string) error {
	if target == "" {
		return fmt.Errorf("target is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.inputFiles) == 0 {
		return fmt.Errorf("no PDFs to append")
	}

	current := filepath.Join(s.inputDir, s.inputFiles[s.currentIndex])
	targetPath := filepath.Join(s.outputDir, target)

	log.Printf("Attempting to merge PDFs:\nTarget: %s\nCurrent: %s\n", targetPath, current)

	cmd := exec.Command("pdfcpu", "merge", "tmp.pdf", targetPath, current)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		log.Printf("Command failed with error: %v\n", err)
		log.Printf("Stdout: %s\n", stdout.String())
		log.Printf("Stderr: %s\n", stderr.String())
		return fmt.Errorf("failed to merge: %v\nStdout: %s\nStderr: %s", err, stdout.String(), stderr.String())
	}

	log.Printf("Merge successful. Output:\n%s\n", stdout.String())

	if err := os.Rename("tmp.pdf", targetPath); err != nil {
		return err
	}

	s.inputFiles = append(s.inputFiles[:s.currentIndex], s.inputFiles[s.currentIndex+1:]...)
	if s.currentIndex >= len(s.inputFiles) && s.currentIndex > 0 {
		s.currentIndex--
	}

	return nil
}

func (s *Service) Trash() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.inputFiles) == 0 {
		return fmt.Errorf("no PDFs to trash")
	}

	oldPath := filepath.Join(s.inputDir, s.inputFiles[s.currentIndex])
	newPath := filepath.Join(s.stagingDir, "trash", s.inputFiles[s.currentIndex])

	if err := os.MkdirAll(filepath.Dir(newPath), 0755); err != nil {
		return err
	}

	if err := os.Rename(oldPath, newPath); err != nil {
		return err
	}

	s.inputFiles = append(s.inputFiles[:s.currentIndex], s.inputFiles[s.currentIndex+1:]...)
	if s.currentIndex >= len(s.inputFiles) && s.currentIndex > 0 {
		s.currentIndex--
	}

	return nil
}

func (s *Service) GetOutputFoldersRecursive() []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	return append([]string(nil), s.folders...)
}

func (s *Service) GetInputFiles() []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	return append([]string(nil), s.inputFiles...)
}

func (s *Service) GetOutputFiles() []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	return append([]string(nil), s.outputFiles...)
}
