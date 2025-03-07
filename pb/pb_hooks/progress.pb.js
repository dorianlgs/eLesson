// create progress records for every assignee added when a course record is created
onRecordCreateRequest((e) => {
  e.next();
  const record = e.record;
  const progressCollection = $app.findCollectionByNameOrId("progress");
  const usersCollection = $app.findCollectionByNameOrId("users");
  const assign_to_everyone = record.getBool("assign_to_everyone");
  let assignees = record.getStringSlice("assignees");

  if (assign_to_everyone) {
    const allUsers = $app.findAllRecords(usersCollection.name);
    assignees = allUsers.map((user) => user.id);
    record.set("assignees", assignees);
    $app.save(record);
  }

  assignees.forEach((assignee) => {
    const progressRecord = new Record(progressCollection, {
      course: record.id,
      assignee: assignee,
      status: "Not Started",
    });
    $app.save(progressRecord);
  });
}, "courses");

// create/delete progress records for every assignee added/removed when a course record is updated
onRecordUpdateRequest((e) => {
  e.next();
  const updatedRecord = e.record;
  const originalRecord = updatedRecord.original();
  const originalAssignees = originalRecord.getStringSlice("assignees");
  let updatedAssignees = updatedRecord.getStringSlice("assignees");
  const progressCollection = $app.findCollectionByNameOrId("progress");
  const usersCollection = $app.findCollectionByNameOrId("users");
  const assign_to_everyone = updatedRecord.getBool("assign_to_everyone");

  if (assign_to_everyone) {
    const allUsers = $app.findAllRecords(usersCollection.name);
    updatedAssignees = allUsers.map((user) => user.id);
    updatedRecord.set("assignees", updatedAssignees);
    $app.save(updatedRecord);
  }

  const newAssignees = updatedAssignees.filter(
    (assignee) => !originalAssignees.includes(assignee),
  );
  const removedAssignees = originalAssignees.filter(
    (assignee) => !updatedAssignees.includes(assignee),
  );

  newAssignees.forEach((assignee) => {
    const progressRecord = new Record(progressCollection, {
      course: updatedRecord.id,
      assignee: assignee,
      status: "Not Started",
    });

    $app.save(progressRecord);
  });

  removedAssignees.forEach((assignee) => {
    const progressRecords = $app.findAllRecords(
      progressCollection.name,
      $dbx.hashExp({
        assignee: `${assignee}`,
        course: `${updatedRecord.id}`,
      }),
    );

    progressRecords.forEach((progressRecord) => {
      $app.delete(progressRecord);
    });
  });
}, "courses");

// remove assignees from course records when their corresponding progress records are deleted
onRecordDeleteRequest((e) => {
  e.next();
  const deletedProgressRecord = e.record;
  const courseId = deletedProgressRecord.getString("course");

  if (courseId) {
    $app.runInTransaction((txApp) => {
      const courseRecord = txApp.findRecordById("courses", courseId);

      if (courseRecord) {
        const assignees = courseRecord.getStringSlice("assignees");
        const assigneeToRemove = deletedProgressRecord.get("assignee");
        const updatedAssignees = assignees.filter(
          (assignee) => assignee !== assigneeToRemove,
        );

        courseRecord.set("assignees", updatedAssignees);
        txApp.save(courseRecord);
      }
    });
  }
}, "progress");

// reset the course and assignee fields to their original values when they get updated
onRecordUpdateRequest((e) => {
  e.next();
  const updatedRecord = e.record;
  const originalRecord = updatedRecord.original();

  if (
    updatedRecord.getString("course") !== originalRecord.getString("course") ||
    updatedRecord.getString("assignee") !== originalRecord.getString("assignee")
  ) {
    updatedRecord.set("course", originalRecord.getString("course"));
    updatedRecord.set("assignee", originalRecord.getString("assignee"));

    $app.save(updatedRecord);
  }
}, "progress");

// add assignee to the corresponding course record when a progress record is created
onRecordCreateRequest((e) => {
  e.next();
  const progressRecord = e.record;
  const assignee = progressRecord.getString("assignee");
  const courseId = progressRecord.getString("course");

  if (courseId) {
    const courseRecord = $app.findRecordById("courses", courseId);

    if (courseRecord) {
      const assignees = courseRecord.getStringSlice("assignees");
      if (!assignees.includes(assignee)) {
        assignees.push(assignee);

        courseRecord.set("assignees", assignees);
        $app.save(courseRecord);
      }
    }
  }
}, "progress");

// add new users to courses that are assigned to everyone and create progress records for them
onRecordCreateRequest((e) => {
  e.next();
  const newUser = e.record;
  const progressCollection = $app.findCollectionByNameOrId("progress");
  const coursesCollection = $app.findCollectionByNameOrId("courses");

  const assignedToEveryoneCourses = $app.findAllRecords(
    coursesCollection.name,
    $dbx.hashExp({ assign_to_everyone: true }),
  );

  assignedToEveryoneCourses.forEach((course) => {
    const assignees = course.getStringSlice("assignees");
    if (!assignees.includes(newUser.id)) {
      assignees.push(newUser.id);
      course.set("assignees", assignees);
      $app.save(course);

      const progressRecord = new Record(progressCollection, {
        course: course.id,
        assignee: newUser.id,
        status: "Not Started",
      });

      $app.save(progressRecord);
    }
  });
}, "users");
