package scanorganizer

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

type Service struct {
	pdfDir    string
	imgDir    string
	outputDir string
	trashDir  string

	mu              sync.Mutex
	pdfFiles        []string
	currentIndex    int
	previousRenamed []string
	folders         []string
}

type PreviewData struct {
	Preview         string   `json:"preview"`
	PreviousRenamed []string `json:"previousRenamed"`
	Folders         []string `json:"folders"`
}

func NewService(pdfDir, imgDir, outputDir, trashDir string) (*Service, error) {
	s := &Service{
		pdfDir:    pdfDir,
		imgDir:    imgDir,
		outputDir: outputDir,
		trashDir:  trashDir,
	}

	if err := os.MkdirAll(s.imgDir, 0755); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(s.outputDir, 0755); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(s.trashDir, 0755); err != nil {
		return nil, err
	}

	if err := s.loadPDFFiles(); err != nil {
		return nil, err
	}
	if err := s.loadFolders(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Service) loadPDFFiles() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.pdfFiles = nil

	files, err := os.ReadDir(s.pdfDir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if strings.HasSuffix(strings.ToLower(file.Name()), ".pdf") {
			s.pdfFiles = append(s.pdfFiles, file.Name())
		}
	}
	sort.Strings(s.pdfFiles)
	if s.currentIndex >= len(s.pdfFiles) {
		s.currentIndex = 0
	}
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

func (s *Service) renderPreview(pdf string) (string, error) {
	img := strings.TrimSuffix(pdf, ".pdf") + ".png"
	imgPath := filepath.Join(s.imgDir, img)
	if _, err := os.Stat(imgPath); os.IsNotExist(err) {
		cmd := exec.Command("pdftoppm", filepath.Join(s.pdfDir, pdf), filepath.Join(s.imgDir, strings.TrimSuffix(pdf, ".pdf")), "-png", "-singlefile")
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

	if len(s.pdfFiles) == 0 {
		return nil, fmt.Errorf("no PDFs found in %s", s.pdfDir)
	}

	current := s.pdfFiles[s.currentIndex]
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

	if len(s.pdfFiles) == 0 {
		return fmt.Errorf("no PDFs to rename")
	}

	if !strings.HasSuffix(strings.ToLower(newName), ".pdf") {
		newName = newName + ".pdf"
	}

	oldPath := filepath.Join(s.pdfDir, s.pdfFiles[s.currentIndex])
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

	s.pdfFiles = append(s.pdfFiles[:s.currentIndex], s.pdfFiles[s.currentIndex+1:]...)
	if s.currentIndex >= len(s.pdfFiles) && s.currentIndex > 0 {
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

	if len(s.pdfFiles) == 0 {
		return fmt.Errorf("no PDFs to append")
	}

	current := filepath.Join(s.pdfDir, s.pdfFiles[s.currentIndex])
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

	s.pdfFiles = append(s.pdfFiles[:s.currentIndex], s.pdfFiles[s.currentIndex+1:]...)
	if s.currentIndex >= len(s.pdfFiles) && s.currentIndex > 0 {
		s.currentIndex--
	}

	return nil
}

func (s *Service) Trash() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.pdfFiles) == 0 {
		return fmt.Errorf("no PDFs to trash")
	}

	oldPath := filepath.Join(s.pdfDir, s.pdfFiles[s.currentIndex])
	newPath := filepath.Join(s.trashDir, s.pdfFiles[s.currentIndex])

	if err := os.Rename(oldPath, newPath); err != nil {
		return err
	}

	s.pdfFiles = append(s.pdfFiles[:s.currentIndex], s.pdfFiles[s.currentIndex+1:]...)
	if s.currentIndex >= len(s.pdfFiles) && s.currentIndex > 0 {
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

	return append([]string(nil), s.pdfFiles...)
}

