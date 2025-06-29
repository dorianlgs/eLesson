<script>
  import { createBubbler, stopPropagation, handlers } from 'svelte/legacy';
  import Icon from "@iconify/svelte";
  import { t } from "../lib/i18n";

  const bubble = createBubbler();

  let { 
    courseId, 
    status, 
    isLoading, 
    onResetProgress, 
    onStartCourse 
  } = $props();
</script>

<div class="flex items-center gap-3 sm:w-full">
  {#if status === "Completed" || status === "In Progress"}
    <button
      onclick={handlers(stopPropagation(bubble('click')), () => onResetProgress(courseId))}
      class={isLoading
        ? "pointer-events-none line-clamp-1 flex items-center justify-center gap-2 truncate rounded-md px-4 py-2 text-red-400 opacity-50 outline outline-[1.5px] outline-red-400/20 transition hover:bg-red-400/20 sm:w-full sm:flex-1 sm:px-0"
        : "line-clamp-1 flex items-center justify-center gap-2 truncate rounded-md px-4 py-2 text-red-400 outline outline-[1.5px] outline-red-400/20 transition hover:bg-red-400/20 sm:w-full sm:flex-1 sm:px-0"}
    >
      {$t("resetProgress")}
      {#if isLoading}
        <Icon
          class="flex-shrink-0 animate-spin text-base"
          icon="fluent:spinner-ios-16-regular"
        />
      {/if}
    </button>
  {/if}
  <button
    onclick={handlers(stopPropagation(bubble('click')), () => onStartCourse(courseId))}
    class="line-clamp-1 truncate rounded-md bg-white/10 px-4 py-2 outline outline-[1.5px] outline-white/20 transition hover:bg-white/20 sm:w-full sm:flex-1 sm:px-0"
  >
    {status === "Completed"
      ? $t("openCourse")
      : status === "In Progress"
        ? $t("continueCourse")
        : $t("startCourse")}
  </button>
</div> 