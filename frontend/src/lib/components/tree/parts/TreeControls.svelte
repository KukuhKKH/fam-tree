<script lang="ts">
  import { Search, Maximize2, Info } from "lucide-svelte";
  import { Button } from "$lib/components/ui/button/index.js";

  let {
    searchQuery = $bindable(),
    showLegend = $bindable(),
    onSearch,
    onReset,
  } = $props();
</script>

<div
  class="absolute top-4 left-4 right-4 sm:top-8 sm:left-8 sm:right-8 flex flex-col sm:flex-row gap-3 sm:gap-4 pointer-events-none z-40"
>
  <div
    class="flex items-center gap-2 p-1.5 bg-white/50 dark:bg-zinc-900/50 backdrop-blur-2xl border border-zinc-200/50 dark:border-zinc-800/50 rounded-[2rem] shadow-2xl pointer-events-auto max-w-sm w-full ring-1 ring-zinc-950/5 dark:ring-white/5"
  >
    <div class="flex-1 flex items-center gap-3 px-4">
      <Search size={18} class="text-zinc-400" />
      <input
        bind:value={searchQuery}
        onkeydown={(e) => e.key === "Enter" && onSearch()}
        placeholder="Cari nama..."
        class="bg-transparent border-none focus:ring-0 text-sm font-bold w-full text-zinc-900 dark:text-white placeholder:text-zinc-400"
      />
    </div>
    <Button
      size="sm"
      class="rounded-2xl h-10 px-6 font-black tracking-tight cursor-pointer"
      onclick={onSearch}>Cari</Button
    >
  </div>

  <div class="flex items-center gap-3 pointer-events-auto sm:ml-auto">
    <Button
      variant="secondary"
      size="icon"
      class="rounded-full h-12 w-12 shadow-2xl bg-white/50 dark:bg-zinc-900/50 backdrop-blur-2xl border border-zinc-200/50 dark:border-zinc-800/50 ring-1 ring-zinc-950/5 dark:ring-white/5 cursor-pointer hover:scale-110 active:scale-95 transition-all"
      onclick={onReset}
      title="Reset View"
    >
      <Maximize2 size={20} />
    </Button>
    <Button
      variant="secondary"
      class="rounded-[1.5rem] h-12 px-6 shadow-2xl bg-white/50 dark:bg-zinc-900/50 backdrop-blur-2xl border border-zinc-200/50 dark:border-zinc-800/50 ring-1 ring-zinc-950/5 dark:ring-white/5 gap-3 font-black cursor-pointer hover:scale-105 active:scale-95 transition-all"
      onclick={() => (showLegend = !showLegend)}
    >
      <Info size={18} />
      <span class="hidden sm:inline">Legenda</span>
    </Button>
  </div>
</div>
