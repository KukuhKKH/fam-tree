<script lang="ts">
  import { onMount } from "svelte";
  import { Button } from "$lib/components/ui/button/index.js";
  import { Sun, Moon } from "lucide-svelte";

  let theme = $state("light");

  onMount(() => {
    theme = document.documentElement.classList.contains("dark")
      ? "dark"
      : "light";
  });

  function toggleTheme() {
    const isDark = document.documentElement.classList.toggle("dark");
    theme = isDark ? "dark" : "light";
    localStorage.setItem("theme", theme);
  }
</script>

<Button
  variant="ghost"
  size="icon"
  onclick={toggleTheme}
  class="rounded-full cursor-pointer"
>
  {#if theme === "light"}
    <Sun class="h-5 w-5" />
  {:else}
    <Moon class="h-5 w-5 text-yellow-400" />
  {/if}
  <span class="sr-only">Toggle theme</span>
</Button>
