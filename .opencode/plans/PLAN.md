# Scan Organizer - Product Plan

## Vision
A privacy-focused, local-only desktop application for organizing scanned PDF documents. Targets individuals and small business users (accounting, billing, medical offices).

## Core Principles
- **Local-only** - No cloud, no network calls, no telemetry
- **Privacy-first** - Documents never leave the user's machine
- **Cross-platform** - Windows, macOS, Linux support
- **Simple but powerful**

---

## Phase 1: Foundation & Polish (v1.0)

### 1.1 Modern UI/UX
- [ ] Replace inline HTML with React or Vue
- [ ] Modern styling, consistent components
- [ ] Dark/light theme support
- [ ] Loading states, progress indicators
- [ ] Toast notifications

### 1.2 Reliability
- [ ] Comprehensive error handling
- [ ] Undo functionality
- [ ] Logging system
- [ ] Crash recovery

### 1.3 Core Workflow
- [ ] Drag-and-drop import
- [ ] Keyboard shortcuts
- [ ] Batch operations
- [ ] File browser for output

### 1.4 Cross-Platform Build
- [ ] CI/CD for Windows/macOS/Linux
- [ ] Native installers (NSIS, DMG)

---

## Phase 2: Smart Features (v1.1)

### 2.1 OCR Integration
- [ ] Local OCR (Tesseract or pure Go)
- [ ] Text extraction from PDFs
- [ ] Auto-detect document type

### 2.2 Smart Suggestions
- [ ] Pattern matching for dates, vendors, amounts
- [ ] Auto-suggest folder names
- [ ] Configurable naming rules

---

## Phase 3: Advanced (v2.0)
- [ ] Local LLM for intelligent categorization (optional)
- [ ] Data export (CSV/JSON)
- [ ] Plugin system

---

## Out of Scope (v1.0)
- Cloud sync / backup
- Multi-user collaboration
- Mobile app
- Web version
