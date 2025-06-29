<script>
  import Icon from "@iconify/svelte";
  import { navigate } from "svelte5-router";
  import { isSidebarVisible } from "../lib/store";
  import { t } from "../lib/i18n";

  let { 
    lessonTitle, 
    currentCourseStatus, 
    onGoToPreviousLesson, 
    onGoToNextLesson, 
    onCompleteCourse,
    isLoading,
    isFirstLesson,
    isLastLesson,
    isCompleted
  } = $props();
</script>

<div class="space-y-5">
  <button
    onclick={() => navigate("/")}
    class="flex items-center gap-2 text-white/50 transition hover:text-white"
  >
    <Icon class="flex-shrink-0" icon="ph:arrow-left" />
    {$t("myCourses")}
  </button>
  
  <div
    class="flex items-end justify-between gap-5 sm:w-full sm:flex-col sm:items-start sm:gap-3"
  >
    <div class="flex items-center gap-3">
      <button
        onclick={() => ($isSidebarVisible = !$isSidebarVisible)}
        class="group flex items-center justify-center rounded-full bg-transparent p-2 text-xl transition hover:bg-white/10"
      >
        <Icon
          class="flex-shrink-0 text-white/50 transition group-hover:text-white"
          icon="ph:list"
        />
      </button>
      <h1 class="text-balance text-xl lg:text-lg">
        {lessonTitle}
      </h1>
    </div>

    {#if currentCourseStatus === "In Progress" || currentCourseStatus === "Completed"}
      <div class="flex items-center gap-3 sm:w-full sm:flex-col">
        {#if !isFirstLesson}
          <button
            onclick={onGoToPreviousLesson}
            class="line-clamp-1 flex items-center justify-center gap-2 truncate rounded-md bg-white/10 px-4 py-2 outline outline-[1.5px] outline-white/20 transition hover:bg-white/20 sm:order-last sm:w-full"
          >
            <Icon class="flex-shrink-0" icon="ph:arrow-left" />
            {$t("previousLesson")}
          </button>
        {/if}

        {#if isLastLesson}
          <button
            onclick={onCompleteCourse}
            class={isLoading || isCompleted
              ? "pointer-events-none line-clamp-1 flex items-center justify-center gap-2 truncate rounded-md bg-emerald-400/60 px-4 py-2 opacity-50 transition hover:bg-emerald-400/50 sm:order-first sm:w-full"
              : "line-clamp-1 flex items-center justify-center gap-2 truncate rounded-md bg-emerald-400/60 px-4 py-2 transition hover:bg-emerald-400/50 sm:order-first sm:w-full"}
          >
            {isCompleted ? $t("courseCompleted") : $t("completeCourse")}
            {#if isLoading}
              <Icon
                class="flex-shrink-0 animate-spin text-base"
                icon="fluent:spinner-ios-16-regular"
              />
            {:else}
              <Icon class="flex-shrink-0" icon="ph:check" />
            {/if}
          </button>
        {:else}
          <button
            onclick={onGoToNextLesson}
            class="line-clamp-1 flex items-center justify-center gap-2 truncate rounded-md bg-main px-4 py-2 transition hover:bg-main/80 sm:order-first sm:w-full"
          >
            {$t("nextLesson")}
            <Icon class="flex-shrink-0" icon="ph:arrow-right" />
          </button>
        {/if}
      </div>
    {/if}
  </div>
</div> 