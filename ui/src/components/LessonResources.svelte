<script>
  import Icon from "@iconify/svelte";
  import { lesson_resources } from "../lib/pocketbase";

  let { lessonId } = $props();

  let lessonResources = $derived($lesson_resources.filter((resource) => resource.lesson.includes(lessonId)));
</script>

{#if lessonResources.length > 0}
  <div class="flex-1 space-y-4 md:w-full">
    <h2 class="flex items-center gap-2 text-base">
      <Icon class="flex-shrink-0" icon="ph:link" />
      Resources
    </h2>
    {#each lessonResources as resource}
      <a
        href={resource.link}
        target="_blank"
        class="block w-full rounded-md bg-white/10 p-2 outline outline-[1.5px] outline-white/20 transition hover:bg-white/20"
      >
        <div class="flex items-center justify-between gap-2">
          <h3 class="line-clamp-1">{resource.name}</h3>
          <Icon
            class="flex-shrink-0"
            icon="ph:arrow-up-right"
          />
        </div>
      </a>
    {/each}
  </div>
{/if} 