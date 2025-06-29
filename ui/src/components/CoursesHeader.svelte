<script>
  import Icon from "@iconify/svelte";
  import { isSidebarVisible, isLoading } from "../lib/store";
  import { courses } from "../lib/pocketbase";
  import { t } from "../lib/i18n";
</script>

<div
  class="flex w-full items-center justify-between gap-5 sm:flex-col sm:items-start"
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
    <h1 class="flex items-center gap-2 text-base">
      <Icon class="flex-shrink-0" icon="ph:graduation-cap" />
      {$t("myCourses")}
    </h1>
  </div>
  {#if $isLoading}
    <div
      class="w-full max-w-56 animate-pulse rounded-full bg-white/10 p-1 sm:hidden"
    ></div>
  {:else if $courses.length > 0}
    <h2 class="text-white/50 sm:hidden">
      {$courses.length}
      {$courses.length === 1 ? $t("courseAssigned") : $t("coursesAssigned")}
    </h2>
  {/if}
</div> 