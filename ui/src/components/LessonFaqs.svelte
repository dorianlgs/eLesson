<script>
  import { slide } from "svelte/transition";
  import { quintOut } from "svelte/easing";
  import Icon from "@iconify/svelte";
  import { lesson_faqs } from "../lib/pocketbase";

  let { lessonId } = $props();

  let lessonFaqs = $derived($lesson_faqs.filter((faq) => faq.lesson.includes(lessonId)));
</script>

{#if lessonFaqs.length > 0}
  <div class="flex-1 space-y-4 md:w-full">
    <h2 class="flex items-center gap-2 text-base">
      <Icon class="flex-shrink-0" icon="ph:chats" />
      FAQs
    </h2>
    {#each lessonFaqs as faq}
      <button
        onclick={() => (faq.isOpen = !faq.isOpen)}
        class="w-full space-y-2 rounded-md bg-white/10 p-2 outline outline-[1.5px] outline-white/20 transition hover:bg-white/20"
      >
        <div
          class="flex items-center justify-between gap-2 text-start"
        >
          <h3 class="line-clamp-1">
            {faq.question}
          </h3>
          <Icon
            class={faq.isOpen
              ? "flex-shrink-0 rotate-45 transition"
              : "flex-shrink-0 transition"}
            icon="ph:plus"
          />
        </div>
        {#if faq.isOpen}
          <p
            transition:slide={{
              duration: 250,
              easing: quintOut,
            }}
            class="text-start text-white/60"
          >
            {faq.answer}
          </p>
        {/if}
      </button>
    {/each}
  </div>
{/if} 