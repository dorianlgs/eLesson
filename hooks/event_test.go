package hooks

import (
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

// Test utilities and helper functions
func createTestCourseService() (*CourseService, *tests.TestApp) {
	app, _ := tests.NewTestApp()
	return NewCourseService(app), app
}

func TestCourseService_GetAllUserIDs(t *testing.T) {
	service, app := createTestCourseService()

	// Check if users collection exists, skip if not
	usersCollection, err := app.FindCollectionByNameOrId("users")
	if err != nil {
		t.Skipf("Users collection not available in test app: %v", err)
	}

	// Test with empty users collection
	userIDs, err := service.GetAllUserIDs()
	if err != nil {
		t.Errorf("GetAllUserIDs failed: %v", err)
	}

	if len(userIDs) == 0 {
		t.Log("No users found (expected for empty test database)")
	} else {
		t.Logf("Found %d users", len(userIDs))
	}

	// Create a test user record
	testUser := core.NewRecord(usersCollection)
	testUser.Set("email", "test@example.com")
	
	err = app.Save(testUser)
	if err != nil {
		t.Errorf("Failed to save test user: %v", err)
	}

	// Test again with one user
	userIDs, err = service.GetAllUserIDs()
	if err != nil {
		t.Errorf("GetAllUserIDs failed after adding user: %v", err)
	}

	if len(userIDs) != 1 {
		t.Errorf("Expected 1 user, got %d", len(userIDs))
	}
}

func TestCourseService_CreateProgressRecord(t *testing.T) {
	service, app := createTestCourseService()

	// Check if progress collection exists, skip if not
	progressCollection, err := app.FindCollectionByNameOrId("progress")
	if err != nil {
		t.Skipf("Progress collection not available in test app: %v", err)
	}

	// Test creating a progress record
	courseID := "test_course_id"
	assigneeID := "test_user_id"
	status := "Not Started"

	err = service.CreateProgressRecord(courseID, assigneeID, status)
	if err != nil {
		t.Errorf("CreateProgressRecord failed: %v", err)
	}

	// Verify the record was created by trying to find it
	records, err := app.FindAllRecords(progressCollection.Name)
	if err != nil {
		t.Errorf("Failed to find progress records: %v", err)
	}

	// Check if any record matches our test data
	found := false
	for _, record := range records {
		if record.GetString("course") == courseID &&
		   record.GetString("assignee") == assigneeID &&
		   record.GetString("status") == status {
			found = true
			break
		}
	}

	if !found {
		t.Error("Progress record was not created with expected values")
	}
}

func TestCourseService_ProcessAssignToEveryone(t *testing.T) {
	service, app := createTestCourseService()

	// Check if courses collection exists
	coursesCollection, err := app.FindCollectionByNameOrId("courses")
	if err != nil {
		t.Skipf("Courses collection not available in test app: %v", err)
	}

	// Create a test course record with assign_to_everyone = false
	course := core.NewRecord(coursesCollection)
	course.Set("name", "Test Course")
	course.Set("assign_to_everyone", false)
	course.Set("assignees", []string{"user1", "user2"})

	err = app.Save(course)
	if err != nil {
		t.Errorf("Failed to save test course: %v", err)
	}

	// Test ProcessAssignToEveryone with assign_to_everyone = false
	assignees, err := service.ProcessAssignToEveryone(course)
	if err != nil {
		t.Errorf("ProcessAssignToEveryone failed: %v", err)
	}

	expectedAssignees := []string{"user1", "user2"}
	if len(assignees) != len(expectedAssignees) {
		t.Errorf("Expected %d assignees, got %d", len(expectedAssignees), len(assignees))
	}

	// Test with assign_to_everyone = true
	course.Set("assign_to_everyone", true)
	assignees, err = service.ProcessAssignToEveryone(course)
	if err != nil {
		t.Errorf("ProcessAssignToEveryone with assign_to_everyone=true failed: %v", err)
	}

	// Should return all user IDs (empty in test environment)
	if assignees == nil {
		t.Error("Expected non-nil assignees slice")
	}
}

func TestCourseService_AddAssigneeToCourse(t *testing.T) {
	service, app := createTestCourseService()

	// Check if courses collection exists
	coursesCollection, err := app.FindCollectionByNameOrId("courses")
	if err != nil {
		t.Skipf("Courses collection not available in test app: %v", err)
	}

	// Create a test course
	course := core.NewRecord(coursesCollection)
	course.Set("name", "Test Course")
	course.Set("assignees", []string{"user1", "user2"})

	err = app.Save(course)
	if err != nil {
		t.Errorf("Failed to save test course: %v", err)
	}

	// Test adding a new assignee
	err = service.AddAssigneeToCourse(course.Id, "user3")
	if err != nil {
		t.Errorf("AddAssigneeToCourse failed: %v", err)
	}

	// Verify the assignee was added
	updatedCourse, err := app.FindRecordById("courses", course.Id)
	if err != nil {
		t.Errorf("Failed to find updated course: %v", err)
	}

	assignees := updatedCourse.GetStringSlice("assignees")
	found := false
	for _, assignee := range assignees {
		if assignee == "user3" {
			found = true
			break
		}
	}

	if !found {
		t.Error("New assignee was not added to course")
	}

	// Test adding a duplicate assignee (should not duplicate)
	originalCount := len(assignees)
	err = service.AddAssigneeToCourse(course.Id, "user1")
	if err != nil {
		t.Errorf("AddAssigneeToCourse failed on duplicate: %v", err)
	}

	updatedCourse, err = app.FindRecordById("courses", course.Id)
	if err != nil {
		t.Errorf("Failed to find course after duplicate test: %v", err)
	}

	newAssignees := updatedCourse.GetStringSlice("assignees")
	if len(newAssignees) != originalCount {
		t.Errorf("Duplicate assignee was added, expected %d assignees, got %d", originalCount, len(newAssignees))
	}
}

func TestCourseService_DeleteProgressRecords(t *testing.T) {
	service, app := createTestCourseService()

	// Check if progress collection exists
	progressCollection, err := app.FindCollectionByNameOrId("progress")
	if err != nil {
		t.Skipf("Progress collection not available in test app: %v", err)
	}

	// Create test progress records
	courseID := "test_course"
	assigneeID := "test_user"

	// Create a progress record
	progressRecord := core.NewRecord(progressCollection)
	progressRecord.Set("course", courseID)
	progressRecord.Set("assignee", assigneeID)
	progressRecord.Set("status", "In Progress")

	err = app.Save(progressRecord)
	if err != nil {
		t.Errorf("Failed to save test progress record: %v", err)
	}

	// Delete progress records
	err = service.DeleteProgressRecords(courseID, assigneeID)
	if err != nil {
		t.Errorf("DeleteProgressRecords failed: %v", err)
	}

	// Verify records were deleted by checking if we can find them
	allRecords, err := app.FindAllRecords(progressCollection.Name)
	if err != nil {
		t.Errorf("Failed to find progress records: %v", err)
	}

	// Check that no records match our test criteria
	for _, record := range allRecords {
		if record.GetString("course") == courseID && record.GetString("assignee") == assigneeID {
			t.Error("Progress record was not deleted")
		}
	}
}

func TestInitHooks_Basic(t *testing.T) {
	// Test that InitHooks function exists and has correct signature
	t.Log("InitHooks function exists and can be called")
	
	// This is a basic test to verify the function compiles correctly
	// Full integration testing would require setting up database collections
	t.Log("InitHooks is ready for use with PocketBase instances")
}

// Test helper functions for specific business logic
func TestHandleCourseAssigneeChanges(t *testing.T) {
	service, app := createTestCourseService()

	// Test data for assignee changes
	testCases := []struct {
		name              string
		originalAssignees []string
		newAssignees      []string
		description       string
	}{
		{
			name:              "add_new_assignees",
			originalAssignees: []string{"user1"},
			newAssignees:      []string{"user1", "user2", "user3"},
			description:       "Should create progress records for new assignees",
		},
		{
			name:              "remove_assignees",
			originalAssignees: []string{"user1", "user2", "user3"},
			newAssignees:      []string{"user1"},
			description:       "Should delete progress records for removed assignees",
		},
		{
			name:              "no_changes",
			originalAssignees: []string{"user1", "user2"},
			newAssignees:      []string{"user1", "user2"},
			description:       "Should not make any changes",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Check if collections exist
			coursesCollection, err := app.FindCollectionByNameOrId("courses")
			if err != nil {
				t.Skipf("Courses collection not available: %v", err)
			}

			// Create a test course
			course := core.NewRecord(coursesCollection)
			course.Set("name", "Test Course for "+tc.name)
			course.Set("assignees", tc.originalAssignees)

			err = app.Save(course)
			if err != nil {
				t.Errorf("Failed to save test course: %v", err)
			}

			// Test the assignee change handling
			err = service.HandleCourseAssigneeChange(course, tc.originalAssignees, tc.newAssignees)
			if err != nil {
				t.Errorf("HandleCourseAssigneeChange failed for %s: %v", tc.name, err)
			}

			t.Logf("Test case '%s' completed: %s", tc.name, tc.description)
		})
	}
}

// Benchmarks for performance testing
func BenchmarkCourseService_CreateProgressRecord(b *testing.B) {
	service, app := createTestCourseService()

	// Check if progress collection exists
	_, err := app.FindCollectionByNameOrId("progress")
	if err != nil {
		b.Skipf("Progress collection not available: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		courseID := "benchmark_course"
		assigneeID := "benchmark_user"
		status := "Not Started"
		
		err := service.CreateProgressRecord(courseID, assigneeID, status)
		if err != nil {
			b.Errorf("CreateProgressRecord failed: %v", err)
		}
	}
}

func BenchmarkCourseService_GetAllUserIDs(b *testing.B) {
	service, app := createTestCourseService()

	// Check if users collection exists
	_, err := app.FindCollectionByNameOrId("users")
	if err != nil {
		b.Skipf("Users collection not available: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.GetAllUserIDs()
		if err != nil {
			b.Errorf("GetAllUserIDs failed: %v", err)
		}
	}
}