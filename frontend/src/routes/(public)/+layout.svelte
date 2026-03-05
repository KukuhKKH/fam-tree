<script lang="ts">
  import { Button } from "$lib/components/ui/button/index.js";
  import { TreeDeciduous } from "lucide-svelte";
  import ThemeToggle from "$lib/components/ThemeToggle.svelte";

  let { children, data } = $props();
</script>

<div class="min-h-screen bg-zinc-50 dark:bg-zinc-950 flex flex-col">
  <header
    class="h-16 border-b border-zinc-200 dark:border-zinc-800 bg-white/80 dark:bg-zinc-900/80 backdrop-blur-md sticky top-0 z-50 px-6 flex items-center justify-between"
  >
    <a
      href="/"
      class="flex items-center gap-2 hover:opacity-80 transition-opacity"
    >
      <TreeDeciduous class="text-primary" size={24} />
      <span class="font-bold text-xl tracking-tight">Silsilah</span>
    </a>
    <div class="flex items-center gap-3">
      <ThemeToggle />
      {#if data.user}
        <Button
          href="/dashboard"
          variant="outline"
          class="rounded-xl font-bold cursor-pointer"
        >
          Ke Dashboard
        </Button>
      {:else}
        <Button
          href="/login"
          variant="outline"
          class="rounded-xl font-bold cursor-pointer"
        >
          Masuk
        </Button>
      {/if}
    </div>
  </header>

  <main class="flex-1 overflow-x-hidden">
    {@render children()}
  </main>

  <footer
    class="p-6 border-t border-zinc-200 dark:border-zinc-800 text-center text-sm text-zinc-500"
  >
    &copy; {new Date().getFullYear()} Silsilah Keluarga — Platform Silsilah Mandiri
  </footer>
</div>
