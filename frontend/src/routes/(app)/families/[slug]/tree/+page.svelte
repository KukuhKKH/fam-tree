<script lang="ts">
  import { Button } from "$lib/components/ui/button/index.js";
  import {
    ChevronLeft,
    Download,
    Maximize,
    Focus,
    Share2,
    Settings2,
    Loader2,
    XCircle,
  } from "lucide-svelte";
  import Network from "$lib/components/tree/Network.svelte";
  import AddRelationshipDialog from "$lib/components/relationships/AddRelationshipDialog.svelte";
  import TreeSkeleton from "$lib/components/tree/parts/TreeSkeleton.svelte";

  import { toast } from "svelte-sonner";

  const { data } = $props();
  // Family is immediate, treeData is streamed
  const family = $derived(data.family);
  const slug = $derived(data.slug);
</script>

<svelte:head>
  <title>Pohon Silsilah {family.name} | Silsilah</title>
</svelte:head>

<div class="h-[calc(100vh-6rem)] flex flex-col space-y-6">
  <!-- Top Navigation / Header -->
  <div
    class="flex flex-col lg:flex-row items-start lg:items-center justify-between gap-4 px-2 sm:px-4"
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
      {#await data.streamed.treeData}
        <Button
          variant="outline"
          disabled
          class="rounded-2xl h-12 px-5 font-bold gap-2"
        >
          <Loader2 class="animate-spin size-4" />
          Memuat Data...
        </Button>
      {:then resolvedTreeData}
        <AddRelationshipDialog
          persons={resolvedTreeData.nodes}
          familySlug={slug}
        />
      {/await}

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
        onclick={() => {
          const url = `${window.location.origin}/s/${slug}`;
          navigator.clipboard.writeText(url);
          toast.success("Link publik berhasil disalin!");
        }}
      >
        <Share2 size={18} />
      </Button>
    </div>
  </div>

  <!-- Main Network Visualizer -->
  <div class="flex-1 min-h-0 px-2 sm:px-4 pb-4 sm:pb-6">
    {#await data.streamed.treeData}
      <TreeSkeleton />
    {:then resolvedTreeData}
      <Network treeData={resolvedTreeData} familySlug={slug} />
    {:catch error}
      <div
        class="w-full h-full flex flex-col items-center justify-center space-y-4 bg-rose-50 dark:bg-rose-950/10 rounded-[2rem] border border-rose-100 dark:border-rose-900/30"
      >
        <div
          class="p-4 rounded-full bg-rose-100 dark:bg-rose-900/30 text-rose-600"
        >
          <XCircle size={32} />
        </div>
        <p class="font-bold text-rose-800 dark:text-rose-200">
          Gagal memuat silsilah: {error.message}
        </p>
        <Button variant="outline" onclick={() => window.location.reload()}
          >Coba Lagi</Button
        >
      </div>
    {/await}
  </div>
</div>
