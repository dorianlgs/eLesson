<script>
  import Icon from "@iconify/svelte";
  import { lessons } from "../lib/pocketbase";
  import { t } from "../lib/i18n";

  let { courseId } = $props();

  let lessonCount = $derived($lessons.filter((lesson) => lesson.course === courseId).length);
</script>

<h3 class="flex items-center gap-2 text-white/50">
  <Icon class="flex-shrink-0 text-lg" icon="ph:book-open" />
  {lessonCount}
  {lessonCount === 1
    ? $t("lessonInThisCourse")
    : $t("lessonsInThisCourse")}
</h3> 