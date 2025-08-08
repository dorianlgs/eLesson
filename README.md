# eLesson

A Learning Management System built with PocketBase (Go backend) and Svelte 5 frontend.

## Features

- **Course Management**: Create and manage courses with automatic user assignment
- **Progress Tracking**: Automatic progress record creation and management
- **User Assignment**: Support for individual and "assign to everyone" functionality
- **Video Lessons**: Integrated video player with Plyr
- **Internationalization**: Multi-language support with svelte-i18n
- **Real-time Updates**: Live data updates via PocketBase subscriptions

## Quick Start

### Prerequisites
- Go 1.24+
- Node.js and npm

### Installation

```bash
# Install Go dependencies
go mod tidy

# Install frontend dependencies
cd ui && npm install
```

### Development

```bash
# Run development server with hot reload
go run . serve

# Frontend development (separate terminal)
cd ui && npm run dev
```

### Testing

```bash
# Run backend tests
go test ./hooks -v

# Run all tests with coverage
go test ./... -cover
```

### Build for Production

```bash
# Build frontend and embed into Go binary
go generate ./...

# Build production binary
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w"

# Run production server
./eLesson serve
```

## Architecture

### Backend (Go + PocketBase)
- **PocketBase Framework**: SQLite database with built-in auth and real-time subscriptions
- **Refactored Hooks**: Modular business logic with proper error handling and comprehensive tests
- **Embedded Frontend**: UI built into Go binary for easy deployment

### Frontend (Svelte 5)
- **Modern Framework**: Latest Svelte with runes and reactive state management
- **TailwindCSS**: Utility-first CSS framework for styling
- **Component Architecture**: Reusable UI components with proper separation of concerns

### Key Improvements (Recent)
- **Refactored Hooks**: `hooks/event.go` refactored into maintainable `CourseService` with individual methods
- **Comprehensive Tests**: Full test coverage in `hooks/event_test.go` with unit tests and benchmarks
- **Better Error Handling**: Proper error wrapping and descriptive error messages
- **Code Organization**: Separated concerns and improved maintainability

## Database Collections

- **courses**: Course information and assignee management
- **lessons**: Individual lesson content and resources  
- **users**: User authentication and profiles
- **progress**: User progress tracking through courses
- **resources**: Course/lesson attachments
- **lesson_faqs**: FAQ content for lessons
- **lesson_resources**: Resource associations

## Deployment

### Server Deployment
```bash
# Copy binary to server
rsync -av /Users/doriangonzalez/Workspace/eLesson/eLesson root@192.168.100.175:/var/www/pocketbase

# Copy migrations
rsync -av /Users/doriangonzalez/Workspace/eLesson/pb_migrations/ root@192.168.100.175:/var/www/pb_migrations/
```

### Development Commands

```bash
# Module creation (if starting fresh)
go mod init github.com/dorianlgs/eLesson

# Update all Go modules
go get -u -t ./...
go mod tidy

# Generate and build
go generate ./...
go build
```

## Contributing

1. Follow the existing code patterns and architecture
2. Run tests before submitting changes: `go test ./... -v`
3. Ensure frontend builds successfully: `cd ui && npm run build`
4. Update documentation when adding new features

## License

Private project - All rights reserved.
