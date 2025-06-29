package hooks

import (
	"slices"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func InitHooks(app *pocketbase.PocketBase) error {

	// create progress records for every assignee added when a course record is created
	app.OnRecordCreateRequest("courses").BindFunc(func(e *core.RecordRequestEvent) error {
		e.Next()
		record := e.Record

		progressCollection, err := app.FindCollectionByNameOrId("progress")
		if err != nil {
			return err
		}

		usersCollection, err := app.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}

		assign_to_everyone := record.GetBool("assign_to_everyone")
		assignees := record.GetStringSlice("assignees")

		if assign_to_everyone {
			allUsers, err := app.FindAllRecords(usersCollection.Name)
			if err != nil {
				return err
			}

			assignees = []string{}
			for _, user := range allUsers {
				assignees = append(assignees, user.Id)
			}

			record.Set("assignees", assignees)
			if err := app.Save(record); err != nil {
				return err
			}
		}

		for _, assignee := range assignees {
			progressRecord := core.NewRecord(progressCollection)
			progressRecord.Set("course", record.Id)
			progressRecord.Set("assignee", assignee)
			progressRecord.Set("status", "Not Started")
			if err := app.Save(progressRecord); err != nil {
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
		updatedAssignees := updatedRecord.GetStringSlice("assignees")

		progressCollection, err := app.FindCollectionByNameOrId("progress")
		if err != nil {
			return err
		}

		usersCollection, err := app.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}

		assign_to_everyone := updatedRecord.GetBool("assign_to_everyone")

		if assign_to_everyone {
			allUsers, err := app.FindAllRecords(usersCollection.Name)
			if err != nil {
				return err
			}

			updatedAssignees := []string{}
			for _, user := range allUsers {
				updatedAssignees = append(updatedAssignees, user.Id)
			}
			updatedRecord.Set("assignees", updatedAssignees)
			if err := app.Save(updatedRecord); err != nil {
				return err
			}
		}

		newAssignees := []string{}
		for _, assignee := range updatedAssignees {
			if !slices.Contains(originalAssignees, assignee) {
				newAssignees = append(newAssignees, assignee)
			}
		}

		removedAssignees := []string{}
		for _, assignee := range originalAssignees {
			if !slices.Contains(updatedAssignees, assignee) {
				removedAssignees = append(removedAssignees, assignee)
			}
		}

		for _, assignee := range newAssignees {
			progressRecord := core.NewRecord(progressCollection)
			progressRecord.Set("course", updatedRecord.Id)
			progressRecord.Set("assignee", assignee)
			progressRecord.Set("status", "Not Started")
			if err := app.Save(progressRecord); err != nil {
				return err
			}
		}

		for _, assignee := range removedAssignees {
			progressRecords, err := app.FindAllRecords(progressCollection.Name,
				dbx.HashExp{
					"assignee": assignee,
					"course":   updatedRecord.Id})

			if err != nil {
				return err
			}

			for _, progressRecord := range progressRecords {
				if err := app.Delete(progressRecord); err != nil {
					return err
				}
			}
		}

		return nil
	})

	// remove assignees from course records when their corresponding progress records are deleted
	app.OnRecordDeleteRequest("progress").BindFunc(func(e *core.RecordRequestEvent) error {
		e.Next()

		deletedProgressRecord := e.Record
		courseId := deletedProgressRecord.GetString("course")

		if courseId != "" {
			err := app.RunInTransaction(func(txApp core.App) error {
				courseRecord, err := txApp.FindRecordById("courses", courseId)
				if err != nil {
					return err
				}

				if courseRecord != nil {
					assignees := courseRecord.GetStringSlice("assignees")
					assigneeToRemove := deletedProgressRecord.GetString("assignee")

					updatedAssignees := []string{}
					for _, assignee := range assignees {
						if assignee != assigneeToRemove {
							updatedAssignees = append(updatedAssignees, assignee)
						}
					}

					courseRecord.Set("assignees", updatedAssignees)
					if err := txApp.Save(courseRecord); err != nil {
						return err
					}

					return nil
				}

				return nil
			})

			if err != nil {
				return err
			}
		}

		return nil
	})

	// reset the course and assignee fields to their original values when they get updated
	app.OnRecordUpdateRequest("progress").BindFunc(func(e *core.RecordRequestEvent) error {
		e.Next()

		updatedRecord := e.Record
		originalRecord := updatedRecord.Original()

		if updatedRecord.GetString("course") != originalRecord.GetString("course") ||
			updatedRecord.GetString("assignee") != originalRecord.GetString("assignee") {

			updatedRecord.Set("course", originalRecord.GetString("course"))
			updatedRecord.Set("assignee", originalRecord.GetString("assignee"))

			if err := app.Save(updatedRecord); err != nil {
				return err
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

		if courseId != "" {
			courseRecord, err := app.FindRecordById("courses", courseId)
			if err != nil {
				return err
			}

			if courseRecord != nil {
				assignees := courseRecord.GetStringSlice("assignees")

				if !slices.Contains(assignees, assignee) {
					assignees = append(assignees, assignee)
					courseRecord.Set("assignees", assignees)
					if err := app.Save(courseRecord); err != nil {
						return err
					}
				}
			}
		}

		return nil
	})

	// add new users to courses that are assigned to everyone and create progress records for them
	app.OnRecordCreateRequest("users").BindFunc(func(e *core.RecordRequestEvent) error {
		e.Next()

		newUser := e.Record
		progressCollection, err := app.FindCollectionByNameOrId("progress")
		if err != nil {
			return err
		}
		coursesCollection, err := app.FindCollectionByNameOrId("courses")
		if err != nil {
			return err
		}

		assignedToEveryoneCourses, err := app.FindAllRecords(
			coursesCollection.Name,
			dbx.HashExp{"assign_to_everyone": true},
		)
		if err != nil {
			return err
		}

		for _, course := range assignedToEveryoneCourses {
			assignees := course.GetStringSlice("assignees")
			if !slices.Contains(assignees, newUser.Id) {

				assignees = append(assignees, newUser.Id)
				course.Set("assignees", assignees)
				if err := app.Save(course); err != nil {
					return err
				}

				progressRecord := core.NewRecord(progressCollection)
				progressRecord.Set("course", course.Id)
				progressRecord.Set("assignee", newUser.Id)
				progressRecord.Set("status", "Not Started")

				if err := app.Save(progressRecord); err != nil {
					return err
				}
			}
		}

		return nil
	})

	return nil
}
