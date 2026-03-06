<script lang="ts">
  import { page } from "$app/state";
  import { Button } from "$lib/components/ui/button/index.js";
  import {
    ChevronLeft,
    Share2,
    TreeDeciduous,
    Eye,
    Globe,
  } from "lucide-svelte";
  import Network from "$lib/components/tree/Network.svelte";
  import { toast } from "svelte-sonner";

  const { data } = $props();
  const treeData = $derived(data.treeData);
  const family = $derived(data.family);
  const slug = $derived(data.slug);

  async function copyShareLink() {
    const url = window.location.href;
    try {
      await navigator.clipboard.writeText(url);
      toast.success("Link berhasil disalin!");
    } catch {
      toast.error("Gagal menyalin link.");
    }
  }
</script>

<svelte:head>
  <title>Pohon Silsilah {family?.name || ""} — Publik | Silsilah</title>
  <meta
    name="description"
    content={`Visualisasi pohon silsilah keluarga ${family?.name || ""}. Telusuri hubungan antar anggota keluarga secara interaktif.`}
  />

  <!-- Open Graph / Facebook -->
  <meta property="og:type" content="article" />
  <meta property="og:url" content={page.url.href} />
  <meta
    property="og:title"
    content={`Pohon Silsilah ${family?.name || ""} — Publik | Silsilah`}
  />
  <meta
    property="og:description"
    content={`Visualisasi interaktif pohon silsilah keluarga ${family?.name || ""}.`}
  />
  <meta property="og:image" content="/og-image.png" />

  <!-- Twitter -->
  <meta property="twitter:card" content="summary_large_image" />
  <meta property="twitter:url" content={page.url.href} />
  <meta
    property="twitter:title"
    content={`Pohon Silsilah ${family?.name || ""} — Publik | Silsilah`}
  />
  <meta
    property="twitter:description"
    content={`Visualisasi interaktif pohon silsilah keluarga ${family?.name || ""}.`}
  />
  <meta property="twitter:image" content="/og-image.png" />
</svelte:head>

<div class="h-[calc(100vh-8rem)] flex flex-col max-w-[100vw] overflow-hidden">
  <!-- Top Navigation / Header -->
  <div
    class="flex flex-col lg:flex-row items-start lg:items-center justify-between gap-4 px-4 sm:px-6 py-4"
  >
    <div class="flex items-center gap-4">
      <Button
        variant="ghost"
        href={`/s/${slug}`}
        class="size-12 p-0 rounded-2xl bg-white dark:bg-zinc-950 border border-zinc-200 dark:border-zinc-900 shadow-sm hover:bg-zinc-100 dark:hover:bg-zinc-900 cursor-pointer"
      >
        <ChevronLeft size={24} />
      </Button>
      <div>
        <div class="flex items-center gap-3">
          <h1
            class="text-2xl sm:text-3xl font-black tracking-tight text-zinc-900 dark:text-white leading-none"
          >
            Pohon Silsilah {family.name}
          </h1>
          <div
            class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-emerald-50 dark:bg-emerald-900/20 border border-emerald-200 dark:border-emerald-800/40 text-emerald-700 dark:text-emerald-400 text-[10px] font-bold uppercase tracking-widest"
          >
            <Eye size={10} />
            Read Only
          </div>
        </div>
        <p class="text-xs sm:text-sm font-medium text-zinc-500 mt-1">
          Tampilan publik — hanya untuk melihat.
        </p>
      </div>
    </div>

    <div class="flex flex-wrap items-center gap-2">
      <Button
        variant="outline"
        class="rounded-2xl h-12 px-5 font-bold gap-2 border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-950 cursor-pointer shadow-sm"
        onclick={copyShareLink}
      >
        <Share2 size={18} />
        <span class="hidden sm:inline">Bagikan</span>
      </Button>
    </div>
  </div>

  <!-- Main Network Visualizer -->
  <div class="flex-1 min-h-0 px-4 sm:px-6 pb-4 sm:pb-6">
    <div
      class="h-full rounded-2xl overflow-hidden border border-zinc-200 dark:border-zinc-800"
    >
      <Network {treeData} familySlug={slug} readOnly={true} />
    </div>
  </div>
</div>
