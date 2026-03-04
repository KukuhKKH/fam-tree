<script lang="ts">
  import { page } from "$app/state";
  import * as Card from "$lib/components/ui/card/index.js";
  import { Users, TreeDeciduous, Heart } from "lucide-svelte";
  import FamilyCard from "$lib/components/FamilyCard.svelte";
  import CreateFamilyDialog from "$lib/components/CreateFamilyDialog.svelte";

  let families = $derived(page.data.families || []);

  const stats = $derived([
    {
      label: "Keluarga",
      value: families.length.toString(),
      icon: Users,
      color: "text-blue-500",
    },
    {
      label: "Anggota Silsilah",
      value: "-",
      icon: TreeDeciduous,
      color: "text-green-500",
    },
    {
      label: "Relasi",
      value: "-",
      icon: Heart,
      color: "text-red-500",
    },
  ]);
</script>

<div class="space-y-10 pb-12">
  <div class="flex flex-col md:flex-row md:items-center justify-between gap-6">
    <div class="space-y-1">
      <h1
        class="text-3xl md:text-4xl font-black tracking-tight text-zinc-900 dark:text-white"
      >
        Selamat datang, <span class="text-primary"
          >{page.data.user?.name || "User"}</span
        >!
      </h1>
      <p class="text-zinc-500 dark:text-zinc-400 font-medium">
        Ikhtisar silsilah keluarga Anda hari ini.
      </p>
    </div>
    <CreateFamilyDialog />
  </div>

  <div class="grid gap-4 sm:gap-6 grid-cols-1 md:grid-cols-3">
    {#each stats as stat}
      <Card.Root
        class="relative overflow-hidden group border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900/50 backdrop-blur-sm transition-all duration-300 hover:shadow-lg hover:shadow-primary/5"
      >
        <div
          class="absolute top-0 right-0 p-8 opacity-[0.03] group-hover:opacity-[0.08] transition-opacity"
        >
          <stat.icon size={80} />
        </div>
        <Card.Header
          class="flex flex-row items-center justify-between space-y-0 pb-2"
        >
          <Card.Title
            class="text-[10px] uppercase tracking-widest font-black text-zinc-400"
            >{stat.label}</Card.Title
          >
          <div
            class="p-2 rounded-lg bg-zinc-100 dark:bg-zinc-800 transition-colors group-hover:bg-primary/10"
          >
            <stat.icon class="h-4 w-4 {stat.color}" />
          </div>
        </Card.Header>
        <Card.Content class="pt-2">
          <div class="text-3xl font-black tracking-tight">{stat.value}</div>
        </Card.Content>
      </Card.Root>
    {/each}
  </div>

  <div class="space-y-6">
    <div
      class="flex items-center justify-between border-b border-zinc-100 dark:border-zinc-800 pb-4"
    >
      <div class="flex items-center gap-2">
        <div class="h-8 w-1.5 bg-primary rounded-full"></div>
        <h2 class="text-xl font-black tracking-tight">Keluarga Saya</h2>
      </div>
    </div>

    {#if families.length === 0}
      <Card.Root
        class="bg-zinc-50/50 dark:bg-zinc-950/30 border-zinc-200 dark:border-zinc-800 border-dashed border-2 rounded-3xl overflow-hidden"
      >
        <Card.Content
          class="p-16 flex flex-col items-center text-center space-y-6"
        >
          <div class="relative">
            <div
              class="absolute inset-0 bg-primary/20 blur-2xl rounded-full"
            ></div>
            <div
              class="relative p-6 bg-white dark:bg-zinc-900 rounded-2xl shadow-xl text-primary border border-zinc-100 dark:border-zinc-800"
            >
              <Users size={48} />
            </div>
          </div>
          <div class="space-y-2 max-w-sm">
            <h3 class="text-xl font-bold">Mulai Perjalanan Anda</h3>
            <p class="text-sm text-zinc-500 leading-relaxed font-medium">
              Anda belum bergabung atau membuat keluarga. Buat keluarga pertama
              Anda untuk mulai mencatat sejarah keturunan.
            </p>
          </div>
          <CreateFamilyDialog />
        </Card.Content>
      </Card.Root>
    {:else}
      <div class="grid gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 pt-2">
        {#each families as family}
          <FamilyCard {family} />
        {/each}
      </div>
    {/if}
  </div>
</div>
