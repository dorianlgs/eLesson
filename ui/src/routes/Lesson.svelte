<script>
  import { onMount, tick } from "svelte";
  import {
    pb,
    lessons,
    courses,
    progress,
    lesson_faqs,
    lesson_resources,
    currentUser,
    fetchRecords,
    updateProgressStatus,
  } from "../lib/pocketbase";
  import { navigate, useLocation } from "svelte5-router";
  import Sidebar from "../components/Sidebar.svelte";
  import {
    isLoading,
    getStoredLessons,
    storeLessons,
    showAlert,
  } from "../lib/store";
  import Plyr from "plyr";
  import slugify from "slugify";
  import NotFound from "./NotFound.svelte";
  import Title from "../components/Title.svelte";
  import { t } from "../lib/i18n";
  import LessonHeader from "../components/LessonHeader.svelte";
  import LessonVideo from "../components/LessonVideo.svelte";
  import LessonContent from "../components/LessonContent.svelte";
  import LessonFooter from "../components/LessonFooter.svelte";
  import LessonLoadingState from "../components/LessonLoadingState.svelte";

  let { lessonTitle } = $props();

  let loading = $state({});
  let lessonVideo;
  let currentCourseStatus = $state("");
  let currentLessonTitle = $state("");

  const lessonLocation = useLocation();

  onMount(async () => {
    if ($currentUser) {
      $isLoading = true;
      await fetchRecords();
      $isLoading = false;
    } else {
      navigate("/login");
    }
  });

  $effect(() => {
    lessonVideo = new Plyr("#lessonVideo", {
      invertTime: false,
      toggleInvert: false,
      captions: {
        active: true,
        update: true,
      },
    });
  });

  // find the current course status
  $effect(() => {
    const currentLesson = $lessons.find(
      (lesson) =>
        slugify(lesson.title, { lower: true, strict: true }) ===
        slugify(lessonTitle, { lower: true, strict: true }),
    );
    if (currentLesson) {
      currentLessonTitle = currentLesson.title;
      const currentCourse = $courses.find(
        (course) => course.id === currentLesson.course,
      );
      if (currentCourse) {
        const currentStatus = $progress.find(
          (progressRecord) => progressRecord.course === currentCourse.id,
        );
        if (currentStatus) {
          currentCourseStatus = currentStatus.status;
        }
      }
    }
  });

  // function to get lessons of the current course
  function getCourseLessons(courseId) {
    return $lessons.filter((lesson) => lesson.course === courseId);
  }

  // function to find the index of the current lesson within its course
  function findCurrentLessonIndex(courseLessons) {
    return courseLessons.findIndex(
      (lesson) =>
        slugify(lesson.title, { lower: true, strict: true }) ===
        slugify(lessonTitle, { lower: true, strict: true }),
    );
  }

  // function to navigate to the next lesson within the same course
  function goToNextLesson() {
    const currentLesson = $lessons.find(
      (lesson) =>
        slugify(lesson.title, { lower: true, strict: true }) ===
        slugify(lessonTitle, { lower: true, strict: true }),
    );
    if (currentLesson) {
      const courseLessons = getCourseLessons(currentLesson.course);
      const currentLessonIndex = findCurrentLessonIndex(courseLessons);
      if (
        currentLessonIndex >= 0 &&
        currentLessonIndex < courseLessons.length - 1
      ) {
        const nextLesson = courseLessons[currentLessonIndex + 1];
        navigate(
          `/${slugify(nextLesson.title, { lower: true, strict: true })}`,
        );

        const lessonsByCourse = getStoredLessons();
        lessonsByCourse[nextLesson.course] = nextLesson;
        storeLessons(lessonsByCourse);
      }
    }
  }

  // function to navigate to the previous lesson within the same course
  function goToPreviousLesson() {
    const currentLesson = $lessons.find(
      (lesson) =>
        slugify(lesson.title, { lower: true, strict: true }) ===
        slugify(lessonTitle, { lower: true, strict: true }),
    );
    if (currentLesson) {
      const courseLessons = getCourseLessons(currentLesson.course);
      const currentLessonIndex = findCurrentLessonIndex(courseLessons);
      if (currentLessonIndex > 0) {
        const previousLesson = courseLessons[currentLessonIndex - 1];
        navigate(
          `/${slugify(previousLesson.title, { lower: true, strict: true })}`,
        );

        const lessonsByCourse = getStoredLessons();
        lessonsByCourse[previousLesson.course] = previousLesson;
        storeLessons(lessonsByCourse);
      }
    }
  }

  // function to complete a course and update the progress status to "Completed"
  async function completeCourse(lessonId) {
    const currentLesson = $lessons.find(
      (lesson) =>
        slugify(lesson.title, { lower: true, strict: true }) ===
        slugify(lessonTitle, { lower: true, strict: true }),
    );
    const currentCourse = $courses.find(
      (course) => course.id === currentLesson.course,
    );
    const progressRecord = $progress.find(
      (progressRecord) => progressRecord.course === currentCourse.id,
    );

    if (progressRecord.status === "In Progress") {
      loading[lessonId] = true;
      const updatedProgressRecord = await updateProgressStatus(
        progressRecord.id,
        "Completed",
      );

      if (!updatedProgressRecord) {
        loading[lessonId] = false;
      }

      if (updatedProgressRecord) {
        await tick();

        $progress = $progress.map((progressRecord) => {
          if (progressRecord.course === currentCourse.id) {
            return { ...progressRecord, status: "Completed" };
          }
          return progressRecord;
        });

        loading[lessonId] = false;

        navigate("/");

        showAlert(
          `${currentCourse.title.length > 30 ? currentCourse.title.slice(0, 30) + "..." : currentCourse.title} completed successfully`,
          "success",
        );
      }
    }
  }

  // Helper functions to determine lesson position
  function getCurrentLesson() {
    return $lessons.find(
      (lesson) =>
        slugify(lesson.title, { lower: true, strict: true }) ===
        slugify(lessonTitle, { lower: true, strict: true }),
    );
  }

  function getLessonPosition() {
    const currentLesson = getCurrentLesson();
    if (!currentLesson) return { isFirst: true, isLast: true };
    
    const courseLessons = getCourseLessons(currentLesson.course);
    const currentIndex = findCurrentLessonIndex(courseLessons);
    
    return {
      isFirst: currentIndex <= 0,
      isLast: currentIndex >= courseLessons.length - 1
    };
  }
</script>

<Title title={currentLessonTitle} />

{#if $currentUser}
  <main class="flex h-dvh justify-between lg:overflow-x-hidden">
    <Sidebar isCoursesVisible={false} />
    {#if $isLoading}
      <LessonLoadingState />
    {:else if $lessons.length === 0 || $lessons.every((lesson) => slugify( lesson.title, { lower: true, strict: true }, ) !== $lessonLocation.pathname.slice(1))}
      <NotFound page="lesson" />
    {:else}
      {#each $lessons as lesson (lesson.id)}
        {#if slugify( lesson.title, { lower: true, strict: true }, ) === slugify( lessonTitle, { lower: true, strict: true }, )}
          {@const lessonPosition = getLessonPosition()}
          {@const hasFooterContent = $lesson_faqs.filter((faq) => faq.lesson.includes(lesson.id)).length > 0 || $lesson_resources.filter((resource) => resource.lesson.includes(lesson.id)).length > 0 || lesson.downloads.length > 0}
          
          <section
            class={hasFooterContent
              ? "flex flex-1 flex-col justify-between gap-5 overflow-y-scroll bg-dark p-5"
              : "flex flex-1 flex-col justify-between overflow-y-scroll bg-dark p-5"}
          >
            <div class="space-y-5">
              <LessonHeader
                lessonTitle={lesson.title}
                currentCourseStatus={currentCourseStatus}
                onGoToPreviousLesson={goToPreviousLesson}
                onGoToNextLesson={goToNextLesson}
                onCompleteCourse={() => completeCourse(lesson.id)}
                isLoading={loading[lesson.id]}
                isFirstLesson={lessonPosition.isFirst}
                isLastLesson={lessonPosition.isLast}
                isCompleted={currentCourseStatus === "Completed"}
              />

              <LessonVideo {lesson} />
              <LessonContent {lesson} />
            </div>

            <LessonFooter {lesson} />
          </section>
        {/if}
      {/each}
    {/if}
  </main>
{/if}
