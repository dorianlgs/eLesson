<script>
  import Icon from "@iconify/svelte";
  import { isSidebarVisible } from "../lib/store";
</script>

<div class="flex w-full flex-col gap-5 p-5">
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
    <div class="w-1/2 animate-pulse rounded-full bg-white/10 p-1"></div>
  </div>
  <div
    class="flex w-full flex-grow animate-pulse items-center justify-center rounded-md bg-white/10 p-5"
  >
    <Icon
      class="flex-shrink-0 text-6xl text-white/10"
      icon="svg-spinners:bars-scale-fade"
    />
  </div>
</div> 