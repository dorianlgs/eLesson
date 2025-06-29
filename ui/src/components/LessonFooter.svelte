<script>
  import LessonFaqs from "./LessonFaqs.svelte";
  import LessonResources from "./LessonResources.svelte";
  import LessonDownloads from "./LessonDownloads.svelte";

  let { lesson } = $props();
</script>

<div
  class="flex w-full items-start justify-between gap-5 md:flex-col"
>
  <LessonFaqs lessonId={lesson.id} />
  <LessonResources lessonId={lesson.id} />
  <LessonDownloads {lesson} />
</div> 