<script lang="ts">
  import { Button } from "$lib/components/ui/button/index.js";
  import {
    ChevronLeft,
    Download,
    Maximize,
    Focus,
    Share2,
    Settings2,
  } from "lucide-svelte";
  import Network from "$lib/components/tree/Network.svelte";
  import AddRelationshipDialog from "$lib/components/relationships/AddRelationshipDialog.svelte";

  const { data } = $props();
  const treeData = $derived(data.treeData);
  const family = $derived(data.family);
  const slug = $derived(data.slug);
</script>

<svelte:head>
  <title>Pohon Silsilah {family.name} | Silsilah</title>
</svelte:head>

<div class="h-[calc(100vh-6rem)] flex flex-col space-y-6">
  <!-- Top Navigation / Header -->
  <div
    class="flex flex-col lg:flex-row items-start lg:items-center justify-between gap-4 px-2"
  >
    <div class="flex items-center gap-4">
      <Button
        variant="ghost"
        href={`/families/${slug}`}
        class="size-12 p-0 rounded-2xl bg-white dark:bg-zinc-950 border border-zinc-200 dark:border-zinc-900 shadow-sm hover:bg-zinc-100 dark:hover:bg-zinc-900 cursor-pointer"
      >
        <ChevronLeft size={24} />
      </Button>
      <div>
        <h1
          class="text-2xl sm:text-3xl font-black tracking-tight text-zinc-900 dark:text-white leading-none mb-1"
        >
          Pohon Silsilah {family.name}
        </h1>
        <p class="text-xs sm:text-sm font-medium text-zinc-500">
          Visualisasi interaktif hubungan antar generasi keluarga.
        </p>
      </div>
    </div>

    <div class="flex flex-wrap items-center gap-2">
      <AddRelationshipDialog persons={treeData.nodes} familySlug={slug} />
      <div
        class="h-8 w-px bg-zinc-200 dark:bg-zinc-800 mx-1 hidden sm:block"
      ></div>
      <Button
        variant="outline"
        class="rounded-2xl h-12 px-5 font-bold gap-2 border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-950 cursor-pointer shadow-sm"
      >
        <Download size={18} />
        <span class="hidden sm:inline">Export</span>
      </Button>
      <Button
        variant="outline"
        class="rounded-2xl h-12 px-5 font-bold gap-2 border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-950 cursor-pointer shadow-sm"
      >
        <Share2 size={18} />
      </Button>
    </div>
  </div>

  <!-- Main Network Visualizer -->
  <div class="flex-1 min-h-0">
    <Network {treeData} familySlug={slug} />
  </div>
</div>
```
