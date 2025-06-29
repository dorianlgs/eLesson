<script>
  let { course } = $props();
</script>

<div class="w-full space-y-2">
  <h1 class="text-base leading-relaxed">{course.title}</h1>
  {#if course.description}
    <h3 class="leading-relaxed text-white/50">
      {course.description}
    </h3>
  {/if}
</div> 