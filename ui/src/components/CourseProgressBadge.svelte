<script>
  import { t } from "../lib/i18n";

  let { status } = $props();
</script>

<h3
  class={status === "Completed"
    ? "rounded-full bg-emerald-400/10 px-3 py-1 text-emerald-400/70"
    : status === "In Progress"
      ? "rounded-full bg-amber-400/10 px-3 py-1 text-amber-400/70"
      : "rounded-full bg-white/10 px-3 py-1 text-white/70"}
>
  {status === "Completed"
    ? $t("completed")
    : status === "In Progress"
      ? $t("inProgress")
      : $t("notStarted")}
</h3> 