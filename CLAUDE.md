# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Go Backend
- **Development**: `go run . serve` - Runs the PocketBase server with hot reload
- **Install dependencies**: `go mod tidy`
- **Production build**: `GOOS=linux GOARCH=amd64 go build -ldflags "-s -w"`
- **Update dependencies**: `go get -u -t ./... && go mod tidy`
- **Run tests**: `go test ./hooks -v` - Test the refactored hooks business logic
- **Run all tests**: `go test ./... -cover` - Run all tests with coverage report
- **Benchmarks**: `go test ./hooks -bench=.` - Run performance benchmarks

### Frontend (UI)
The frontend is located in the `ui/` directory and uses Svelte 5 with Vite.

- **Development**: `cd ui && npm run dev`
- **Build**: `cd ui && npm run build` (also triggered by `go generate ./...`)
- **Install dependencies**: `cd ui && npm install`
- **Preview**: `cd ui && npm run preview`

### Build Process
The frontend is embedded into the Go binary using `go:embed`. Running `go generate ./...` will:
1. Install npm dependencies (`npm install --force`)
2. Build the frontend (`npm run build`)
3. Embed the `dist` folder into the Go binary

## Architecture

### Backend Structure
- **PocketBase Framework**: Built on PocketBase, a Go backend with SQLite database
- **Main Entry**: `main.go` - Initializes PocketBase app with plugins and static file serving
- **Database Hooks**: `hooks/event.go` - Contains business logic hooks for record operations
- **Embedded UI**: `ui/embed.go` - Embeds the frontend build into the Go binary
- **Migrations**: `pb_migrations/` - Database schema migrations in JavaScript

### Frontend Structure (Svelte 5)
- **Router**: Uses `svelte5-router` with routes in `src/routes/`
- **Components**: Reusable UI components in `src/components/`
- **State Management**: PocketBase client in `src/lib/pocketbase.js`
- **Internationalization**: `svelte-i18n` with translations in `src/lib/i18n.js`
- **Styling**: TailwindCSS with custom configuration

### Core Features
- **Learning Management System**: Course creation, lesson management, progress tracking
- **User Management**: Authentication via PocketBase, user assignment to courses
- **Media Support**: Video lessons with Plyr video player
- **Progress Tracking**: Automatic progress record creation/deletion based on course assignments

### Key Hooks Logic (Refactored)
The `hooks/event.go` has been refactored for better maintainability:
- **CourseService**: Encapsulates all course-related business logic with proper error handling
- **Separated Functions**: Individual methods for each operation:
  - `GetAllUserIDs()` - Retrieves all user IDs from the database
  - `CreateProgressRecord()` - Creates new progress tracking records
  - `DeleteProgressRecords()` - Removes progress records for course/assignee pairs
  - `HandleCourseAssigneeChange()` - Manages assignee additions and removals
  - `ProcessAssignToEveryone()` - Handles "assign to everyone" logic with all users
  - `RemoveAssigneeFromCourse()` - Safely removes assignees using transactions
  - `AddAssigneeToCourse()` - Adds assignees while preventing duplicates
  - `AssignUserToAllEveryCourses()` - Assigns new users to all public courses
- **Comprehensive Tests**: Located in `hooks/event_test.go` with unit tests, integration tests, and performance benchmarks
- **Error Handling**: Consistent error wrapping with descriptive messages using `fmt.Errorf`
- **Business Logic**: Auto-creates progress records when courses are assigned to users, handles "assign to everyone" functionality, maintains referential integrity between courses, users, and progress records

### Database Collections
- **courses**: Course information with assignee management
- **lessons**: Individual lesson content and resources
- **users**: User authentication and profile data
- **progress**: Tracking user progress through courses
- **resources**: Course/lesson attachments and files
- **lesson_faqs**: FAQ content for lessons
- **lesson_resources**: Resource associations for lessons

## Development Notes

### PocketBase Integration
- Database operations use PocketBase's ORM and query builder
- All frontend API calls go through PocketBase client (`src/lib/pocketbase.js`)
- Real-time updates available via PocketBase subscriptions
- **Refactored Hooks**: Business logic moved to testable service layer for better maintainability

### Frontend Patterns
- Components follow Svelte 5 runes and syntax
- State management via PocketBase reactive stores
- Responsive design with TailwindCSS utilities
- Video playback handled by Plyr with custom styling

### Testing Strategy
- **Unit Tests**: Each service method has individual tests with success/error scenarios
- **Integration Tests**: End-to-end testing of hook functionality 
- **Performance Tests**: Benchmarks for critical operations to identify bottlenecks
- **Test Utilities**: Helper functions and mock data for consistent testing
- **Coverage**: Run `go test ./... -cover` to check test coverage

### Code Quality Standards
- **Error Handling**: All service methods use proper error wrapping with context
- **Single Responsibility**: Each method handles one specific operation
- **Testability**: Business logic separated from PocketBase hooks for easier testing
- **Documentation**: Function names and tests serve as living documentation

### Deployment
Production deployment involves:
1. Building the Go binary with embedded UI: `go generate ./... && go build`
2. Running tests to ensure quality: `go test ./... -v`
3. Copying binary and migrations to server
4. PocketBase handles database initialization and migrations automatically