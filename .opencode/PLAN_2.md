# Scan Organizer - Project Plan

## Overview

A desktop application for organizing scanned PDF documents. Users drag files/folders onto the app, process them one-by-one (rename, append, trash), and export the organized files to a chosen destination.

## Current Status

### Completed ✓
1. Initial Wails + React setup
2. Backend service restructuring with staging directory
3. Drag-and-drop homepage with mode selection (replace/append/skip)
4. Basic organizer view (preview, rename, append, trash)
5. Frontend builds successfully

### Pending
1. Regenerate Go bindings (run `wails dev`)
2. Test the complete flow
3. Add Export button UI to organizer
4. Any additional features as requested

---

## Architecture

### Directory Structure
```
/home/josh/code/products/scan-organizer/
├── staging/                    # Working directory (gitignored)
│   ├── input/                 # PDFs to process
│   ├── output/                # Processed PDFs
│   └── previews/              # PNG previews
├── app.go                     # Wails app struct
├── main.go                    # Entry point
├── internal/scanorganizer/
│   └── service.go             # Core business logic
└── frontend/
    ├── src/
    │   ├── App.tsx            # Main app with routing
    │   ├── pages/Home.tsx     # Drag-drop homepage
    │   ├── hooks/
    │   │   ├── usePreview.ts  # State for processing
    │   │   └── useDropzone.ts # Drag-drop handling
    │   └── components/        # UI components
    └── bindings/              # Generated Go bindings
```

### Data Flow
```
[Drag files to Home]
    ↓
[Ask: Replace / Append / Skip?]
    ↓
[Copy to ./staging/input]
    ↓
[Organizer view: process files one-by-one]
    ↓
[Move input → output]
    ↓
[Export: move output to chosen location]
    ↓
[Clear exported output, keep input]
```

---

## Backend API

### Service Methods

| Method | Parameters | Returns | Description |
|--------|-----------|---------|-------------|
| `GetStatus` | - | `{inputCount, outputCount, isFullyProcessed}` | Check staging state |
| `AddFiles` | `paths []string, mode string` | `error` | Add files (mode: replace/append/skip) |
| `GetPreview` | - | `{preview, previousRenamed, folders}` | Get current PDF preview |
| `Rename` | `newName, folder string` | `error` | Move current PDF to output |
| `Append` | `target string` | `error` | Merge current PDF into existing |
| `Trash` | - | `error` | Move current PDF to trash |
| `GetOutputFoldersRecursive` | - | `[]string` | List output folders |
| `GetInputFiles` | - | `[]string` | List input files |
| `GetOutputFiles` | - | `[]string` | List output files |
| `Export` | `destination string` | `error` | Move output to destination |

---

## Frontend Pages

### Home Page
- Large drag-drop zone
- Visual feedback on dragover
- Mode dialog when input already has files:
  - **Replace**: Clear input, add new files
  - **Append**: Keep existing, add new files
  - **Skip**: Don't add, keep existing

### Organizer Page
- Left panel: Rename form, Trash button, Append form
- Right panel: PDF preview
- Keyboard shortcuts (Ctrl+1, Ctrl+2, Ctrl+3)

---

## Next Steps

1. **Run `wails dev`** to regenerate bindings and test
2. **Verify** the drag-drop flow works end-to-end
3. **Add Export button** to the organizer view
4. **Test** export moves files correctly

---

## Notes

- Files are copied (not moved) from source to staging
- Input should be empty after full processing (all moved to output)
- Staging persists between sessions - can restart anytime
- Export clears output folder but keeps input for re-export if needed
