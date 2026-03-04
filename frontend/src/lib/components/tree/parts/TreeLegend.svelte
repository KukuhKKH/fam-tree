<script lang="ts">
  let { showLegend = $bindable() } = $props();
</script>

{#if showLegend}
  <div
    class="absolute bottom-8 left-8 p-8 bg-white/70 dark:bg-zinc-900/70 backdrop-blur-3xl border border-zinc-200/50 dark:border-zinc-800/50 rounded-[2.5rem] shadow-2xl space-y-6 max-w-xs animate-in fade-in slide-in-from-bottom-8 duration-500 ring-1 ring-zinc-950/5 dark:ring-white/5"
  >
    <h4 class="font-black text-[10px] uppercase tracking-[0.2em] text-zinc-400">
      Petunjuk Peta
    </h4>
    <div class="space-y-4">
      <div class="flex items-center gap-4">
        <div
          class="size-5 rounded-full bg-[#0ea5e9] shadow-[0_0_15px_-5px_#0ea5e9]"
        ></div>
        <span class="text-xs font-black text-zinc-700 dark:text-zinc-300"
          >Anggota Laki-laki</span
        >
      </div>
      <div class="flex items-center gap-4">
        <div
          class="size-5 rounded-full bg-[#f43f5e] shadow-[0_0_15px_-5px_#f43f5e]"
        ></div>
        <span class="text-xs font-black text-zinc-700 dark:text-zinc-300"
          >Anggota Perempuan</span
        >
      </div>
      <div class="flex items-center gap-4">
        <div class="h-1 w-8 rounded-full bg-[#94a3b8]"></div>
        <span class="text-xs font-black text-zinc-700 dark:text-zinc-300"
          >Keturunan (Vertikal)</span
        >
      </div>
      <div class="flex items-center gap-4">
        <div class="h-1 w-8 border-t-2 border-dashed border-[#fbbf24]"></div>
        <span class="text-xs font-black text-zinc-700 dark:text-zinc-300"
          >Pasangan (Horizontal)</span
        >
      </div>
    </div>
    <p
      class="text-[11px] text-zinc-500 leading-relaxed italic border-t border-zinc-100 dark:border-zinc-800/50 pt-5"
    >
      Scroll untuk zoom, drag untuk geser. Klik anggota untuk detail.
    </p>
  </div>
{/if}
