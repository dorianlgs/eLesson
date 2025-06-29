<script>
  import { pb } from "../lib/pocketbase";

  let { lesson } = $props();
</script>

{#if lesson.video}
  <video
    controls
    crossorigin
    playsinline
    id="lessonVideo"
    data-poster={pb.files.getUrl(lesson, lesson.thumbnail)}
  >
    <source src={pb.files.getUrl(lesson, lesson.video)} />
    <track
      kind="captions"
      label="English captions"
      src={pb.files.getUrl(lesson, lesson.captions)}
      srclang="en"
      default
    />
  </video>
{/if} 