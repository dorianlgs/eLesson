<script>
  import Icon from "@iconify/svelte";
  import { t } from "../lib/i18n";
</script>

<div class="flex w-full flex-1 items-center justify-center">
  <div
    class="flex w-[500px] flex-col items-center gap-5 text-center sm:w-full"
  >
    <Icon
      class="flex-shrink-0 text-6xl text-main"
      icon="ph:book-open-text-fill"
    />
    <div class="space-y-3">
      <h1 class="text-balance text-30px leading-snug">
        {$t("noCourseAssigned")}
      </h1>
      <p class="text-base text-white/50">
        {$t("courseAssignedCheck")}
      </p>
    </div>
  </div>
</div> 