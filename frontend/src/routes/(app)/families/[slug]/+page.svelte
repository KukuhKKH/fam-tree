<script lang="ts">
  import { page } from "$app/state";
  import { Button } from "$lib/components/ui/button/index.js";
  import MemberTable from "$lib/components/members/MemberTable.svelte";
  import InviteMemberDialog from "$lib/components/members/InviteMemberDialog.svelte";
  import FamilySettingsDialog from "$lib/components/families/FamilySettingsDialog.svelte";
  import PersonGrid from "$lib/components/persons/PersonGrid.svelte";
  import {
    TreeDeciduous,
    Users,
    ArrowLeft,
    Globe,
    Lock,
    Info,
    Calendar,
    BookOpen,
    Share2,
    Link,
  } from "lucide-svelte";
  import { toast } from "svelte-sonner";

  const { data } = $props();
  const family = $derived(data.family);
  const members = $derived(data.members);
  const persons = $derived(data.persons);
  const myRole = $derived(family.my_role);
  const canManage = $derived(myRole === "owner" || myRole === "admin");

  import AddRelationshipDialog from "$lib/components/relationships/AddRelationshipDialog.svelte";

  let activeTab = $state("persons");
</script>

<svelte:head>
  <title>{family.name} | Silsilah</title>
</svelte:head>

<div class="space-y-8 pb-12">
  <!-- Header Section - Adaptive Design -->
  <div
    class="rounded-3xl bg-white dark:bg-zinc-900 text-zinc-900 dark:text-zinc-100 overflow-hidden relative shadow-sm border border-zinc-200 dark:border-zinc-800 transition-colors"
  >
    <!-- Background Decoration - Adaptive colors -->
    <div
      class="absolute inset-0 opacity-5 dark:opacity-20 bg-[radial-gradient(circle_at_top_right,_var(--color-primary),_transparent)]"
    ></div>

    <div
      class="p-8 md:p-12 relative z-10 flex flex-col md:flex-row gap-8 items-start md:items-center justify-between"
    >
      <div class="space-y-6 max-w-2xl">
        <div>
          <Button
            variant="ghost"
            href="/dashboard"
            class="text-zinc-500 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-white hover:bg-zinc-100 dark:hover:bg-white/10 -ml-2 rounded-xl transition-all h-8 gap-2 cursor-pointer"
          >
            <ArrowLeft size={16} />
            Kembali ke Dashboard
          </Button>
        </div>

        <div class="space-y-3">
          <div class="flex flex-wrap items-center gap-4">
            <h1
              class="text-4xl md:text-5xl font-black tracking-tight text-zinc-900 dark:text-white leading-tight"
            >
              {family.name}
            </h1>
            <div
              class="px-3 py-1 bg-zinc-100 dark:bg-white/10 backdrop-blur-md rounded-full border border-zinc-200 dark:border-white/20 text-[10px] font-bold uppercase tracking-widest flex items-center gap-1.5 text-zinc-600 dark:text-zinc-200"
            >
              {#if family.visibility === "public"}
                <Globe size={12} class="text-emerald-500" />
                <span>Publik</span>
              {:else}
                <Lock size={12} class="text-zinc-500" />
                <span>Privat</span>
              {/if}
            </div>
          </div>
          <p
            class="text-zinc-600 dark:text-zinc-300 text-lg leading-relaxed font-medium"
          >
            {family.description || "Tidak ada deskripsi untuk keluarga ini."}
          </p>
        </div>

        <div class="flex flex-wrap items-center gap-y-6 gap-x-10 pt-4">
          <div class="flex items-center gap-3">
            <div class="p-2 rounded-lg bg-primary/10 border border-primary/20">
              <Users size={18} class="text-primary font-bold" />
            </div>
            <div class="flex flex-col">
              <span
                class="text-lg font-bold text-zinc-900 dark:text-white leading-tight"
              >
                {persons.length}
              </span>
              <span
                class="text-[10px] uppercase tracking-wider text-zinc-500 font-bold"
              >
                Person
              </span>
            </div>
          </div>
          <div class="flex items-center gap-3">
            <div class="p-2 rounded-lg bg-primary/10 border border-primary/20">
              <BookOpen size={18} class="text-primary font-bold" />
            </div>
            <div class="flex flex-col">
              <span
                class="text-lg font-bold text-zinc-900 dark:text-white leading-tight"
              >
                {members.length}
              </span>
              <span
                class="text-[10px] uppercase tracking-wider text-zinc-500 font-bold"
              >
                Akun
              </span>
            </div>
          </div>
          <div class="flex items-center gap-3">
            <div class="p-2 rounded-lg bg-primary/10 border border-primary/20">
              <Info size={18} class="text-primary font-bold" />
            </div>
            <div class="flex flex-col">
              <span
                class="text-lg font-bold text-zinc-900 dark:text-white leading-tight uppercase italic"
              >
                {myRole}
              </span>
              <span
                class="text-[10px] uppercase tracking-wider text-zinc-500 font-bold"
              >
                Peran Anda
              </span>
            </div>
          </div>
        </div>
      </div>

      <div class="flex items-center gap-3 self-end md:self-auto">
        {#if family.visibility !== "private"}
          <Button
            variant="outline"
            class="rounded-xl h-10 px-4 font-bold gap-2 border-zinc-200 dark:border-zinc-800 cursor-pointer"
            onclick={() => {
              const url = `${window.location.origin}/s/${family.slug}`;
              navigator.clipboard.writeText(url);
              toast.success("Link publik berhasil disalin!");
            }}
          >
            <Share2 size={16} />
            Bagikan
          </Button>
        {/if}
        {#if canManage}
          <FamilySettingsDialog {family} />
          <InviteMemberDialog familySlug={family.slug} />
        {/if}
      </div>
    </div>
  </div>

  <!-- Tab Navigation -->
  <div
    class="flex items-center gap-1 p-1 bg-zinc-100 dark:bg-zinc-950 rounded-2xl w-fit border border-zinc-200 dark:border-zinc-900 shadow-inner"
  >
    <button
      onclick={() => (activeTab = "persons")}
      class="px-6 py-2.5 rounded-xl text-sm font-bold transition-all flex items-center gap-2 {activeTab ===
      'persons'
        ? 'bg-white dark:bg-zinc-800 shadow-sm text-primary'
        : 'text-zinc-500 hover:text-zinc-700 dark:hover:text-zinc-300 hover:bg-zinc-200 dark:hover:bg-zinc-900'}"
    >
      <Users size={16} />
      Daftar Silsilah
    </button>
    <button
      onclick={() => (activeTab = "members")}
      class="px-6 py-2.5 rounded-xl text-sm font-bold transition-all flex items-center gap-2 {activeTab ===
      'members'
        ? 'bg-white dark:bg-zinc-800 shadow-sm text-primary'
        : 'text-zinc-500 hover:text-zinc-700 dark:hover:text-zinc-300 hover:bg-zinc-200 dark:hover:bg-zinc-900'}"
    >
      <Users size={16} />
      Akun Keluarga
    </button>
  </div>

  <!-- Main Content Grid -->
  <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
    <div class="lg:col-span-2 space-y-6">
      {#if activeTab === "persons"}
        <div class="space-y-6">
          <div class="space-y-1 px-1">
            <h3 class="text-xl font-bold flex items-center gap-2">
              <Users size={20} class="text-primary" />
              Manajemen Silsilah
            </h3>
            <p class="text-sm text-zinc-500">
              Kelola data individu dan foto profil dalam pohon keluarga ini.
            </p>
          </div>
          <PersonGrid {persons} familySlug={family.slug} {canManage} />
        </div>
      {:else}
        <div class="space-y-6">
          <div class="flex items-center justify-between px-1">
            <div class="space-y-1">
              <h3 class="text-xl font-bold flex items-center gap-2">
                <Users size={20} class="text-primary" />
                Akun Keluarga
              </h3>
              <p class="text-sm text-zinc-500">
                Akun yang memiliki akses ke dashboard keluarga ini.
              </p>
            </div>
          </div>
          <MemberTable {members} familySlug={family.slug} {myRole} />
        </div>
      {/if}
    </div>

    <!-- Sidebar / Info Section -->
    <div class="space-y-6">
      <div
        class="p-6 rounded-3xl bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 shadow-sm space-y-6"
      >
        <h3 class="font-bold flex items-center gap-2">
          <Info size={18} class="text-primary" />
          Aksi Cepat
        </h3>

        <div class="pt-2 border-zinc-100 dark:border-zinc-800 space-y-3">
          <Button
            variant="outline"
            class="w-full rounded-xl gap-2 font-bold cursor-pointer h-12"
            href={`/families/${family.slug}/tree`}
          >
            <TreeDeciduous size={20} />
            Lihat Pohon Silsilah
          </Button>

          {#if canManage}
            <AddRelationshipDialog
              {persons}
              familySlug={family.slug}
              style="width: 100%;"
            />
          {/if}

          <p class="text-[10px] text-center text-zinc-400 mt-3">
            Gunakan visualisasi interaktif untuk melihat hubungan antar
            generasi.
          </p>
        </div>
      </div>

      <!-- Role Permissions Hint -->
      <div
        class="p-6 rounded-3xl bg-zinc-50 dark:bg-zinc-950 border border-dotted border-zinc-300 dark:border-zinc-700 space-y-4"
      >
        <h4 class="text-xs font-bold uppercase text-zinc-400 tracking-wider">
          Info Role
        </h4>
        <div class="space-y-3 text-xs">
          <div class="flex gap-2">
            <span class="font-bold text-red-500 min-w-[50px]">Owner:</span>
            <span class="text-zinc-500">Kontrol penuh & hapus keluarga.</span>
          </div>
          <div class="flex gap-2">
            <span class="font-bold text-orange-500 min-w-[50px]">Admin:</span>
            <span class="text-zinc-500">Kelola anggota & edit info.</span>
          </div>
          <div class="flex gap-2">
            <span class="font-bold text-blue-500 min-w-[50px]">Editor:</span>
            <span class="text-zinc-500">Edit data silsilah & person.</span>
          </div>
          <div class="flex gap-2">
            <span class="font-bold text-zinc-400 min-w-[50px]">Viewer:</span>
            <span class="text-zinc-500">Hanya melihat data silsilah.</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>
