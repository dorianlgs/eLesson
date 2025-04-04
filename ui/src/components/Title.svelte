<script>
  import { run } from 'svelte/legacy';

  import customize from "../../customize.json";
  const { name } = customize;

  /**
   * @typedef {Object} Props
   * @property {any} [title]
   * @property {string} [suffix]
   */

  /** @type {Props} */
  let { title = name, suffix = "" } = $props();

  let documentTitle = $state();

  run(() => {
    if (suffix === "") {
      documentTitle = title;
    } else {
      documentTitle = `${title} - ${suffix}`;
    }
  });
</script>

<svelte:head>
  <title>{documentTitle}</title>
</svelte:head>
