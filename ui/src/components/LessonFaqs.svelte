<script>
  import { slide } from "svelte/transition";
  import { quintOut } from "svelte/easing";
  import Icon from "@iconify/svelte";
  import { lesson_faqs } from "../lib/pocketbase";

  let { lessonId } = $props();

  let lessonFaqs = $derived(
    $lesson_faqs.filter((faq) => faq.lesson.includes(lessonId)),
  );

  // Track which FAQs are open by their ID
  let openFaqs = $state([]);

  function toggleFaq(faqId) {
    if (openFaqs.includes(faqId)) {
      var index = openFaqs.indexOf(faqId);
      if (index !== -1) {
        openFaqs.splice(index, 1);
      }
    } else {
      openFaqs.push(faqId);
    }
  }
</script>

{#if lessonFaqs.length > 0}
  <div class="flex-1 space-y-4 md:w-full">
    <h2 class="flex items-center gap-2 text-base">
      <Icon class="flex-shrink-0" icon="ph:chats" />
      FAQs
    </h2>
    {#each lessonFaqs as faq}
      <button
        onclick={() => toggleFaq(faq.id)}
        class="w-full space-y-2 rounded-md bg-white/10 p-2 outline outline-[1.5px] outline-white/20 transition hover:bg-white/20"
      >
        <div class="flex items-center justify-between gap-2 text-start">
          <h3 class="line-clamp-1">
            {faq.question}
          </h3>
          <Icon
            class={openFaqs.includes(faq.id)
              ? "flex-shrink-0 rotate-45 transition"
              : "flex-shrink-0 transition"}
            icon="ph:plus"
          />
        </div>

        {#if openFaqs.includes(faq.id)}
          {faq.answer}
        {/if}
      </button>
    {/each}
  </div>
{/if}
