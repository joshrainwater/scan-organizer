package main

import (
	"context"

	"github.com/joshrainwater/scan-organizer/internal/scanorganizer"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type App struct {
	ctx     context.Context
	service *scanorganizer.Service
}

func NewApp() *App {
	service, err := scanorganizer.NewService("./staging")
	if err != nil {
		panic(err)
	}
	return &App{
		service: service,
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) shutdown(_ context.Context) {
}

func (a *App) GetPreview() (*scanorganizer.PreviewData, error) {
	return a.service.GetPreview()
}

func (a *App) Rename(newName, folder string) error {
	return a.service.Rename(newName, folder)
}

func (a *App) Append(target string) error {
	return a.service.Append(target)
}

func (a *App) Trash() error {
	return a.service.Trash()
}

func (a *App) GetOutputFoldersRecursive() []string {
	return a.service.GetOutputFoldersRecursive()
}

func (a *App) GetInputFiles() []string {
	return a.service.GetInputFiles()
}

func (a *App) GetOutputFiles() []string {
	return a.service.GetOutputFiles()
}

func (a *App) AddFiles(paths []string, mode string) error {
	return a.service.AddFiles(paths, mode)
}

func (a *App) Export(destination string) error {
	return a.service.Export(destination)
}

func (a *App) GetStatus() (*scanorganizer.StatusData, error) {
	return a.service.GetStatus()
}

func (a *App) SelectExportDirectory() (string, error) {
	path, err := application.Get().Dialog.OpenFile().
		SetTitle("Select Export Destination").
		CanChooseDirectories(true).
		CanChooseFiles(false).
		PromptForSingleSelection()

	if err != nil {
		return "", err
	}
	if path == "" {
		return "", nil
	}
	return path, nil
}
