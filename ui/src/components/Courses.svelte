<script>
  import { run } from 'svelte/legacy';
  import {
    lessons,
    courses,
    progress,
    updateProgressStatus,
  } from "../lib/pocketbase";
  import slugify from "slugify";
  import {
    isSidebarVisible,
    isLoading,
    getStoredLessons,
    storeLessons,
  } from "../lib/store";
  import { navigate } from "svelte5-router";
  import { tick } from "svelte";
  import { t } from "../lib/i18n";
  import CoursesHeader from "./CoursesHeader.svelte";
  import CoursesLoadingState from "./CoursesLoadingState.svelte";
  import CoursesEmptyState from "./CoursesEmptyState.svelte";
  import CourseCard from "./CourseCard.svelte";

  let isOpen = $state({});
  let loading = $state({});
  let openCourseId = $state("");
  let enableReactivity = $state(true);

  // only courses that match a progress record with "In Progress" status are set to open
  run(() => {
    if (enableReactivity) {
      $progress.forEach((progressRecord) => {
        if (progressRecord.status === "In Progress") {
          isOpen[progressRecord.course] = true;
        }
      });
    }
  });

  // scroll into view of open courses
  run(() => {
    (async () => {
      await tick();
      if (isOpen[openCourseId]) {
        const courseElement = document.getElementById(openCourseId);
        if (courseElement) {
          courseElement.scrollIntoView({ behavior: "smooth" });
        }
      }
    })();
  });

  // function to toggle opening & closing a course
  const toggleCourse = (courseId) => {
    isOpen[courseId] = !isOpen[courseId];
    openCourseId = courseId;
  };

  // function to navigate to the first lesson of a course and update the status to "In Progress"
  async function goToFirstLessonOfCourse(courseId) {
    const progressRecord = $progress.find(
      (progressRecord) => progressRecord.course === courseId,
    );

    if (progressRecord.status === "Not Started") {
      const updatedProgressRecord = await updateProgressStatus(
        progressRecord.id,
        "In Progress",
      );
      if (updatedProgressRecord) {
        await tick();
        $progress = $progress.map((progressRecord) => {
          if (progressRecord.course === courseId) {
            return { ...progressRecord, status: "In Progress" };
          }
          return progressRecord;
        });
      }
    }

    const lessonsByCourse = getStoredLessons();
    const currentLesson = lessonsByCourse[courseId];

    if (currentLesson) {
      navigate(
        `/${slugify(currentLesson.title, { lower: true, strict: true })}`,
      );
    } else {
      const firstLesson = $lessons.find((lesson) => lesson.course === courseId);
      if (firstLesson) {
        navigate(
          `/${slugify(firstLesson.title, { lower: true, strict: true })}`,
        );

        lessonsByCourse[courseId] = firstLesson;
        storeLessons(lessonsByCourse);
      }
    }
  }

  // function to reset the status of a course back to "Not Started"
  async function resetProgress(courseId) {
    const progressRecord = $progress.find(
      (progressRecord) => progressRecord.course === courseId,
    );

    if (
      progressRecord.status === "Completed" ||
      progressRecord.status === "In Progress"
    ) {
      loading[courseId] = true;
      const updatedProgressRecord = await updateProgressStatus(
        progressRecord.id,
        "Not Started",
      );

      if (!updatedProgressRecord) {
        loading[courseId] = false;
      }

      if (updatedProgressRecord) {
        await tick();

        $progress = $progress.map((progressRecord) => {
          if (progressRecord.course === courseId) {
            return { ...progressRecord, status: "Not Started" };
          }
          return progressRecord;
        });

        openCourseId = "";
        enableReactivity = false;
        loading[courseId] = false;
        isOpen[courseId] = false;

        const lessonsByCourse = getStoredLessons();
        delete lessonsByCourse[courseId];
        storeLessons(lessonsByCourse);
      }
    }
  }
</script>

<section class="flex flex-1 flex-col gap-5 overflow-y-scroll bg-dark p-5">
  <CoursesHeader />

  {#if $isLoading}
    <CoursesLoadingState />
  {:else if $courses.length === 0}
    <CoursesEmptyState />
  {:else}
    {#each $courses as course (course.id)}
      <CourseCard
        {course}
        isOpen={isOpen[course.id]}
        {loading}
        onToggleCourse={toggleCourse}
        onResetProgress={resetProgress}
        onStartCourse={goToFirstLessonOfCourse}
      />
    {/each}
  {/if}
</section>
