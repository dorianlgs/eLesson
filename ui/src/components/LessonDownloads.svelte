<script>
  import Icon from "@iconify/svelte";
  import { pb } from "../lib/pocketbase";
  import { cleanFileName } from "../lib/strConverter";
  import { t } from "../lib/i18n";

  let { lesson } = $props();
</script>

{#if lesson.downloads.length > 0}
  <div class="flex-1 space-y-4 md:w-full">
    <h2 class="flex items-center gap-2 text-base">
      <Icon class="flex-shrink-0" icon="ph:file" />
      {$t("downloads")}
    </h2>
    {#each lesson.downloads as download}
      <a
        href={pb.files.getUrl(lesson, download)}
        download
        class="block w-full rounded-md bg-white/10 p-2 outline outline-[1.5px] outline-white/20 transition hover:bg-white/20"
      >
        <div class="flex items-center justify-between gap-2">
          <h3 class="line-clamp-1">{cleanFileName(download)}</h3>
          <Icon class="flex-shrink-0" icon="ph:download-simple" />
        </div>
      </a>
    {/each}
  </div>
{/if} 