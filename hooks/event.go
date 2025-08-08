package hooks

import (
	"fmt"
	"slices"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

type CourseService struct {
	app core.App
}

func NewCourseService(app core.App) *CourseService {
	return &CourseService{app: app}
}

func (cs *CourseService) GetAllUserIDs() ([]string, error) {
	usersCollection, err := cs.app.FindCollectionByNameOrId("users")
	if err != nil {
		return nil, fmt.Errorf("failed to find users collection: %w", err)
	}

	allUsers, err := cs.app.FindAllRecords(usersCollection.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to find all users: %w", err)
	}

	userIDs := make([]string, 0, len(allUsers))
	for _, user := range allUsers {
		userIDs = append(userIDs, user.Id)
	}
	return userIDs, nil
}

func (cs *CourseService) CreateProgressRecord(courseID, assigneeID, status string) error {
	progressCollection, err := cs.app.FindCollectionByNameOrId("progress")
	if err != nil {
		return fmt.Errorf("failed to find progress collection: %w", err)
	}

	progressRecord := core.NewRecord(progressCollection)
	progressRecord.Set("course", courseID)
	progressRecord.Set("assignee", assigneeID)
	progressRecord.Set("status", status)

	if err := cs.app.Save(progressRecord); err != nil {
		return fmt.Errorf("failed to save progress record: %w", err)
	}
	return nil
}

func (cs *CourseService) DeleteProgressRecords(courseID, assigneeID string) error {
	progressCollection, err := cs.app.FindCollectionByNameOrId("progress")
	if err != nil {
		return fmt.Errorf("failed to find progress collection: %w", err)
	}

	progressRecords, err := cs.app.FindAllRecords(progressCollection.Name,
		dbx.HashExp{
			"assignee": assigneeID,
			"course":   courseID,
		})
	if err != nil {
		return fmt.Errorf("failed to find progress records: %w", err)
	}

	for _, progressRecord := range progressRecords {
		if err := cs.app.Delete(progressRecord); err != nil {
			return fmt.Errorf("failed to delete progress record: %w", err)
		}
	}
	return nil
}

func (cs *CourseService) HandleCourseAssigneeChange(courseRecord *core.Record, originalAssignees, newAssignees []string) error {
	toAdd := make([]string, 0)
	toRemove := make([]string, 0)

	for _, assignee := range newAssignees {
		if !slices.Contains(originalAssignees, assignee) {
			toAdd = append(toAdd, assignee)
		}
	}

	for _, assignee := range originalAssignees {
		if !slices.Contains(newAssignees, assignee) {
			toRemove = append(toRemove, assignee)
		}
	}

	for _, assignee := range toAdd {
		if err := cs.CreateProgressRecord(courseRecord.Id, assignee, "Not Started"); err != nil {
			return err
		}
	}

	for _, assignee := range toRemove {
		if err := cs.DeleteProgressRecords(courseRecord.Id, assignee); err != nil {
			return err
		}
	}

	return nil
}

func (cs *CourseService) ProcessAssignToEveryone(record *core.Record) ([]string, error) {
	if !record.GetBool("assign_to_everyone") {
		return record.GetStringSlice("assignees"), nil
	}

	allUserIDs, err := cs.GetAllUserIDs()
	if err != nil {
		return nil, err
	}

	record.Set("assignees", allUserIDs)
	if err := cs.app.Save(record); err != nil {
		return nil, fmt.Errorf("failed to save course with all users: %w", err)
	}

	return allUserIDs, nil
}

func (cs *CourseService) RemoveAssigneeFromCourse(courseID, assigneeToRemove string) error {
	return cs.app.RunInTransaction(func(txApp core.App) error {
		courseRecord, err := txApp.FindRecordById("courses", courseID)
		if err != nil {
			return fmt.Errorf("failed to find course: %w", err)
		}

		if courseRecord == nil {
			return nil
		}

		assignees := courseRecord.GetStringSlice("assignees")
		updatedAssignees := make([]string, 0, len(assignees))

		for _, assignee := range assignees {
			if assignee != assigneeToRemove {
				updatedAssignees = append(updatedAssignees, assignee)
			}
		}

		courseRecord.Set("assignees", updatedAssignees)
		if err := txApp.Save(courseRecord); err != nil {
			return fmt.Errorf("failed to save course after removing assignee: %w", err)
		}

		return nil
	})
}

func (cs *CourseService) AddAssigneeToCourse(courseID, assigneeID string) error {
	courseRecord, err := cs.app.FindRecordById("courses", courseID)
	if err != nil {
		return fmt.Errorf("failed to find course: %w", err)
	}

	if courseRecord == nil {
		return nil
	}

	assignees := courseRecord.GetStringSlice("assignees")
	if !slices.Contains(assignees, assigneeID) {
		assignees = append(assignees, assigneeID)
		courseRecord.Set("assignees", assignees)
		if err := cs.app.Save(courseRecord); err != nil {
			return fmt.Errorf("failed to save course with new assignee: %w", err)
		}
	}
	return nil
}

func (cs *CourseService) AssignUserToAllEveryCourses(userID string) error {
	_, err := cs.app.FindCollectionByNameOrId("progress")
	if err != nil {
		return fmt.Errorf("failed to find progress collection: %w", err)
	}

	coursesCollection, err := cs.app.FindCollectionByNameOrId("courses")
	if err != nil {
		return fmt.Errorf("failed to find courses collection: %w", err)
	}

	assignedToEveryoneCourses, err := cs.app.FindAllRecords(
		coursesCollection.Name,
		dbx.HashExp{"assign_to_everyone": true},
	)
	if err != nil {
		return fmt.Errorf("failed to find courses assigned to everyone: %w", err)
	}

	for _, course := range assignedToEveryoneCourses {
		assignees := course.GetStringSlice("assignees")
		if !slices.Contains(assignees, userID) {
			assignees = append(assignees, userID)
			course.Set("assignees", assignees)
			if err := cs.app.Save(course); err != nil {
				return fmt.Errorf("failed to save course with new user: %w", err)
			}

			if err := cs.CreateProgressRecord(course.Id, userID, "Not Started"); err != nil {
				return err
			}
		}
	}

	return nil
}

func InitHooks(app *pocketbase.PocketBase) error {
	courseService := NewCourseService(app)

	// create progress records for every assignee added when a course record is created
	app.OnRecordCreateRequest("courses").BindFunc(func(e *core.RecordRequestEvent) error {
		e.Next()
		record := e.Record

		assignees, err := courseService.ProcessAssignToEveryone(record)
		if err != nil {
			return err
		}

		for _, assignee := range assignees {
			if err := courseService.CreateProgressRecord(record.Id, assignee, "Not Started"); err != nil {
				return err
			}
		}

		return nil
	})

	// create/delete progress records for every assignee added/removed when a course record is updated
	app.OnRecordUpdateRequest("courses").BindFunc(func(e *core.RecordRequestEvent) error {
		e.Next()

		updatedRecord := e.Record
		originalRecord := updatedRecord.Original()
		originalAssignees := originalRecord.GetStringSlice("assignees")

		// Process assign_to_everyone and get final assignees
		updatedAssignees, err := courseService.ProcessAssignToEveryone(updatedRecord)
		if err != nil {
			return err
		}

		// Handle assignee changes
		return courseService.HandleCourseAssigneeChange(updatedRecord, originalAssignees, updatedAssignees)
	})

	// remove assignees from course records when their corresponding progress records are deleted
	app.OnRecordDeleteRequest("progress").BindFunc(func(e *core.RecordRequestEvent) error {
		e.Next()

		deletedProgressRecord := e.Record
		courseId := deletedProgressRecord.GetString("course")
		assigneeToRemove := deletedProgressRecord.GetString("assignee")

		if courseId != "" && assigneeToRemove != "" {
			return courseService.RemoveAssigneeFromCourse(courseId, assigneeToRemove)
		}

		return nil
	})

	// reset the course and assignee fields to their original values when they get updated
	app.OnRecordUpdateRequest("progress").BindFunc(func(e *core.RecordRequestEvent) error {
		e.Next()

		updatedRecord := e.Record
		originalRecord := updatedRecord.Original()

		originalCourse := originalRecord.GetString("course")
		originalAssignee := originalRecord.GetString("assignee")

		if updatedRecord.GetString("course") != originalCourse || updatedRecord.GetString("assignee") != originalAssignee {
			updatedRecord.Set("course", originalCourse)
			updatedRecord.Set("assignee", originalAssignee)

			if err := app.Save(updatedRecord); err != nil {
				return fmt.Errorf("failed to revert progress record changes: %w", err)
			}
		}

		return nil
	})

	// add assignee to the corresponding course record when a progress record is created
	app.OnRecordCreateRequest("progress").BindFunc(func(e *core.RecordRequestEvent) error {
		e.Next()

		progressRecord := e.Record
		assignee := progressRecord.GetString("assignee")
		courseId := progressRecord.GetString("course")

		if courseId != "" && assignee != "" {
			return courseService.AddAssigneeToCourse(courseId, assignee)
		}

		return nil
	})

	// add new users to courses that are assigned to everyone and create progress records for them
	app.OnRecordCreateRequest("users").BindFunc(func(e *core.RecordRequestEvent) error {
		e.Next()

		newUser := e.Record
		return courseService.AssignUserToAllEveryCourses(newUser.Id)
	})

	return nil
}
