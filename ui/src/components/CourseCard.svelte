<script>
  import { run } from 'svelte/legacy';
  import { lessons, progress } from "../lib/pocketbase";
  import CourseProgressBadge from "./CourseProgressBadge.svelte";
  import CourseLessonCount from "./CourseLessonCount.svelte";
  import CourseActions from "./CourseActions.svelte";
  import CourseInfo from "./CourseInfo.svelte";
  import CourseLessonItem from "./CourseLessonItem.svelte";

  let { 
    course, 
    isOpen, 
    loading, 
    onToggleCourse, 
    onResetProgress, 
    onStartCourse 
  } = $props();

  let progressRecord = $derived($progress.find((p) => p.course === course.id));
  let courseLessons = $derived($lessons.filter((lesson) => lesson.course === course.id));
</script>

<div
  id={course.id}
  class={isOpen
    ? "w-full rounded-md outline outline-[1.5px] outline-white/20 transition-all hover:outline-white/20"
    : "w-full rounded-md outline outline-[1.5px] outline-white/10 transition-all hover:outline-white/20"}
>
  <div
    aria-hidden="true"
    onclick={() => onToggleCourse(course.id)}
    class={isOpen
      ? "w-full cursor-pointer space-y-5 rounded-b-none rounded-t-md bg-white/5 p-5"
      : "w-full cursor-pointer space-y-5 rounded-md bg-white/5 p-5"}
  >
    {#if progressRecord}
      <div
        class="flex w-full items-center justify-between gap-5 sm:flex-col"
      >
        <div
          class="flex items-center gap-3 sm:w-full xs:flex-col xs:items-start"
        >
          <CourseProgressBadge status={progressRecord.status} />
          <CourseLessonCount courseId={course.id} />
        </div>
        <CourseActions
          courseId={course.id}
          status={progressRecord.status}
          isLoading={loading[course.id]}
          onResetProgress={onResetProgress}
          onStartCourse={onStartCourse}
        />
      </div>
    {/if}
    <CourseInfo {course} />
  </div>
  
  {#if isOpen}
    {#each courseLessons as lesson (lesson.id)}
      <CourseLessonItem {lesson} />
    {/each}
  {/if}
</div> 