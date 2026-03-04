<script lang="ts">
  import PersonCard from "./PersonCard.svelte";
  import AddPersonDialog from "./AddPersonDialog.svelte";
  import { User, Search, Filter } from "lucide-svelte";
  import { Input } from "$lib/components/ui/input/index.js";

  let { persons = [], familySlug, canManage } = $props();

  let searchQuery = $state("");
  let genderFilter = $state("all");

  const filteredPersons = $derived(
    persons.filter((p: any) => {
      const matchesSearch =
        p.full_name.toLowerCase().includes(searchQuery.toLowerCase()) ||
        p.nickname?.toLowerCase().includes(searchQuery.toLowerCase());
      const matchesGender = genderFilter === "all" || p.gender === genderFilter;
      return matchesSearch && matchesGender;
    }),
  );
</script>

<div class="space-y-6">
  <!-- Filter & Search Bar -->
  <div
    class="flex flex-col md:flex-row gap-4 items-center justify-between pb-2"
  >
    <div class="relative w-full md:max-w-md group">
      <Search
        class="absolute left-3 top-1/2 -translate-y-1/2 text-zinc-400 group-focus-within:text-primary transition-colors"
        size={18}
      />
      <Input
        placeholder="Cari nama atau panggilan..."
        bind:value={searchQuery}
        class="pl-10 rounded-2xl border-zinc-200 dark:border-zinc-800 focus:ring-primary/20 transition-all"
      />
    </div>

    <div class="flex items-center gap-2 w-full md:w-auto">
      <div
        class="flex bg-zinc-100 dark:bg-zinc-800 p-1 rounded-xl border border-zinc-200 dark:border-zinc-700"
      >
        <button
          onclick={() => (genderFilter = "all")}
          class="px-4 py-1.5 rounded-lg text-xs font-bold transition-all {genderFilter ===
          'all'
            ? 'bg-white dark:bg-zinc-900 shadow-sm text-primary'
            : 'text-zinc-500 hover:text-zinc-700 dark:hover:text-zinc-300'}"
        >
          Semua
        </button>
        <button
          onclick={() => (genderFilter = "male")}
          class="px-4 py-1.5 rounded-lg text-xs font-bold transition-all {genderFilter ===
          'male'
            ? 'bg-white dark:bg-zinc-900 shadow-sm text-blue-500'
            : 'text-zinc-500 hover:text-zinc-700 dark:hover:text-zinc-300'}"
        >
          Laki-laki
        </button>
        <button
          onclick={() => (genderFilter = "female")}
          class="px-4 py-1.5 rounded-lg text-xs font-bold transition-all {genderFilter ===
          'female'
            ? 'bg-white dark:bg-zinc-900 shadow-sm text-rose-500'
            : 'text-zinc-500 hover:text-zinc-700 dark:hover:text-zinc-300'}"
        >
          Perempuan
        </button>
      </div>

      {#if canManage}
        <AddPersonDialog {familySlug} />
      {/if}
    </div>
  </div>

  <!-- Results Grid -->
  {#if filteredPersons.length > 0}
    <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
      {#each filteredPersons as person (person.id)}
        <PersonCard {person} {familySlug} {canManage} />
      {/each}
    </div>
  {:else}
    <div
      class="py-20 flex flex-col items-center justify-center text-center space-y-4 rounded-3xl border-2 border-dashed border-zinc-100 dark:border-zinc-800"
    >
      <div
        class="p-4 rounded-full bg-zinc-50 dark:bg-zinc-900 text-zinc-300 dark:text-zinc-700"
      >
        <User size={48} />
      </div>
      <div class="space-y-1">
        <h3 class="font-bold text-lg text-zinc-900 dark:text-zinc-100">
          Person tidak ditemukan
        </h3>
        <p class="text-sm text-zinc-500 max-w-xs mx-auto">
          {#if searchQuery || genderFilter !== "all"}
            Coba ubah kriteria pencarian atau filter Anda.
          {:else}
            Belum ada data silsilah di keluarga ini. Mulai dengan menambahkan
            person pertama.
          {/if}
        </p>
      </div>
      {#if canManage && !searchQuery && genderFilter === "all"}
        <AddPersonDialog {familySlug} />
      {/if}
    </div>
  {/if}
</div>
