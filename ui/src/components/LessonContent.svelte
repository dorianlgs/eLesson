<script>
  import Icon from "@iconify/svelte";
  import { t } from "../lib/i18n";

  let { lesson } = $props();
</script>

{#if lesson.content}
  <div
    class="flex w-full items-center justify-center rounded-md bg-white/5 p-[rfs(50px)] sm:p-5"
  >
    <article
      class="flex w-full max-w-[1280px] flex-1 flex-col gap-5 prose-headings:leading-tight prose-h1:text-50px prose-h2:text-40px prose-h3:text-30px prose-h4:text-25px prose-h5:text-20px prose-h6:text-18px prose-p:text-base prose-p:leading-relaxed prose-p:text-white/50 prose-a:text-base prose-a:text-main prose-a:underline prose-a:underline-offset-4 prose-blockquote:w-fit prose-blockquote:rounded-md prose-blockquote:border-l-2 prose-blockquote:border-l-white/50 prose-blockquote:bg-white/5 prose-blockquote:p-5 prose-strong:text-white prose-code:text-base prose-code:text-white/50 prose-ol:list-inside prose-ol:list-decimal prose-ol:space-y-2 prose-ol:text-base
        prose-ul:list-inside prose-ul:list-disc prose-ul:space-y-2 prose-ul:text-base prose-img:w-full prose-img:rounded-md"
    >
      {@html lesson.content}
    </article>
  </div>
{/if}

{#if lesson.summary}
  <div class="space-y-2">
    <h3 class="flex items-center gap-2 text-base">
      <Icon icon="ph:text-align-left flex-shrink-0" />
      {$t("summary")}
    </h3>
    <p class="text-white/50">
      {lesson.summary}
    </p>
  </div>
{/if} 