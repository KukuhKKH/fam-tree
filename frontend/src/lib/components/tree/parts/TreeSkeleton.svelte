<script lang="ts">
  import { Loader2 } from "lucide-svelte";
</script>

<div
  class="relative w-full h-[calc(100vh-12rem)] min-h-[450px] sm:min-h-[600px] rounded-[2rem] sm:rounded-[3.5rem] bg-zinc-100/50 dark:bg-zinc-900/20 border border-dashed border-zinc-200 dark:border-zinc-800 shadow-inner overflow-hidden flex flex-col items-center justify-center space-y-6 animate-pulse"
>
  <!-- Decorative background circles -->
  <div
    class="absolute top-1/4 left-1/4 size-64 bg-primary/5 rounded-full blur-[100px]"
  ></div>
  <div
    class="absolute bottom-1/4 right-1/4 size-80 bg-rose-500/5 rounded-full blur-[120px]"
  ></div>

  <div class="relative">
    <div
      class="size-24 rounded-[2.5rem] bg-white dark:bg-zinc-900 shadow-xl flex items-center justify-center ring-1 ring-zinc-200/50 dark:ring-zinc-800/50"
    >
      <Loader2 class="size-10 text-primary animate-spin" />
    </div>
  </div>

  <div class="space-y-2 text-center">
    <div class="h-6 w-48 bg-zinc-200 dark:bg-zinc-800 rounded-lg mx-auto"></div>
    <div class="h-3 w-64 bg-zinc-100 dark:bg-zinc-900 rounded-lg mx-auto"></div>
  </div>

  <!-- Dummy controls skeleton -->
  <div class="absolute bottom-12 flex gap-4">
    <div
      class="size-12 rounded-2xl bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 shadow-sm"
    ></div>
    <div
      class="size-12 rounded-2xl bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 shadow-sm"
    ></div>
    <div
      class="size-12 rounded-2xl bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 shadow-sm"
    ></div>
  </div>
</div>
