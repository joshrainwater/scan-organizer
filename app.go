package main

import (
	"context"

	"github.com/joshrainwater/scan-organizer/internal/scanorganizer"
)

type App struct {
	ctx     context.Context
	service *scanorganizer.Service
}

func NewApp() *App {
	service, err := scanorganizer.NewService("./input", "./static/previews", "./output", "./trash")
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
	// no-op for now
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

