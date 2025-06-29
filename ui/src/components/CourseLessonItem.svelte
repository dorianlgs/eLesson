<script>
  import Icon from "@iconify/svelte";
  import { navigate } from "svelte5-router";
  import slugify from "slugify";
  import { t } from "../lib/i18n";

  let { lesson } = $props();
</script>

<div
  class="flex w-full items-center justify-between gap-5 border-t-[1.5px] border-t-white/10 p-5"
>
  <div class="flex items-center gap-3">
    {#if lesson.video}
      <Icon
        class="flex-shrink-0 text-3xl text-main"
        icon="ph:video"
      />
    {:else if lesson.content}
      <Icon
        class="flex-shrink-0 text-3xl text-main"
        icon="ph:text-align-left"
      />
    {/if}
    <h3
      class="line-clamp-1 truncate text-wrap break-all text-base"
    >
      {lesson.title}
    </h3>
  </div>
  <button
    onclick={() =>
      navigate(
        `/${slugify(lesson.title, { lower: true, strict: true })}`,
      )}
    class="flex items-center gap-2 p-2 text-white/50 transition hover:text-white"
  >
    <Icon class="flex-shrink-0 text-lg" icon="ph:eye" />
    {$t("view")}
  </button>
</div> 