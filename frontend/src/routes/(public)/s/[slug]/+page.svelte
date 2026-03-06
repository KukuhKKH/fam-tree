<script lang="ts">
  import { page } from "$app/state";
  import { Button } from "$lib/components/ui/button/index.js";
  import {
    TreeDeciduous,
    Users,
    Globe,
    ArrowRight,
    Eye,
    Calendar,
    MapPin,
    Heart,
    Skull,
    User,
    Share2,
    ExternalLink,
  } from "lucide-svelte";
  import { toast } from "svelte-sonner";

  const { data } = $props();
  const family = $derived(data.family);
  const persons = $derived(data.persons);
  const treeData = $derived(data.treeData);
  const slug = $derived(data.slug);

  function formatDate(dateStr: string | null | undefined): string {
    if (!dateStr) return "";
    try {
      return new Date(dateStr).toLocaleDateString("id-ID", {
        day: "numeric",
        month: "long",
        year: "numeric",
      });
    } catch {
      return "";
    }
  }

  function getAge(
    birthDate: string | null | undefined,
    isAlive: boolean,
    deathDate?: string | null,
  ): string {
    if (!birthDate) return "";
    const born = new Date(birthDate);
    const end = isAlive ? new Date() : deathDate ? new Date(deathDate) : null;
    if (!end) return "";
    const age = Math.floor(
      (end.getTime() - born.getTime()) / (365.25 * 24 * 60 * 60 * 1000),
    );
    return isAlive ? `${age} tahun` : `Meninggal usia ${age} tahun`;
  }

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
  <title>Keluarga {family?.name || ""} — Silsilah Publik</title>
  <meta
    name="description"
    content={family?.description ||
      `Lihat silsilah keluarga ${family?.name || ""} secara publik. Temukan leluhur, kerabat, dan sejarah keluarga Anda.`}
  />

  <!-- Open Graph / Facebook -->
  <meta property="og:type" content="article" />
  <meta property="og:url" content={page.url.href} />
  <meta
    property="og:title"
    content={`Keluarga ${family?.name || ""} — Silsilah Publik`}
  />
  <meta
    property="og:description"
    content={family?.description ||
      `Lihat silsilah keluarga ${family?.name || ""} secara publik di platform Silsilah.`}
  />
  <meta property="og:image" content="/og-image.png" />

  <!-- Twitter -->
  <meta property="twitter:card" content="summary_large_image" />
  <meta property="twitter:url" content={page.url.href} />
  <meta
    property="twitter:title"
    content={`Keluarga ${family?.name || ""} — Silsilah Publik`}
  />
  <meta
    property="twitter:description"
    content={family?.description ||
      `Lihat silsilah keluarga ${family?.name || ""} secara publik di platform Silsilah.`}
  />
  <meta property="twitter:image" content="/og-image.png" />
</svelte:head>

<div class="max-w-7xl mx-auto px-6 py-12 space-y-10">
  <!-- Hero Section -->
  <div
    class="relative rounded-3xl overflow-hidden bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 shadow-sm"
  >
    <!-- Gradient decoration -->
    <div
      class="absolute inset-0 bg-gradient-to-br from-primary/5 via-transparent to-emerald-500/5"
    ></div>

    <div class="relative p-8 md:p-12 space-y-6">
      <div
        class="flex flex-col md:flex-row md:items-center justify-between gap-6"
      >
        <div class="space-y-4 max-w-2xl">
          <div
            class="inline-flex items-center gap-2 px-3 py-1.5 rounded-full bg-emerald-50 dark:bg-emerald-900/20 border border-emerald-200 dark:border-emerald-800/40 text-emerald-700 dark:text-emerald-400 text-[11px] font-bold uppercase tracking-widest"
          >
            <Globe size={12} />
            Silsilah Publik
          </div>

          <h1
            class="text-4xl md:text-5xl font-black tracking-tight text-zinc-900 dark:text-white leading-tight"
          >
            {family.name}
          </h1>

          <p
            class="text-zinc-600 dark:text-zinc-300 text-lg leading-relaxed font-medium"
          >
            {family.description ||
              "Silsilah keluarga yang dibagikan secara publik."}
          </p>

          <!-- Stats -->
          <div class="flex flex-wrap items-center gap-8 pt-2">
            <div class="flex items-center gap-3">
              <div
                class="p-2 rounded-xl bg-primary/10 border border-primary/20"
              >
                <Users size={18} class="text-primary" />
              </div>
              <div class="flex flex-col">
                <span
                  class="text-xl font-black text-zinc-900 dark:text-white leading-tight"
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
              <div
                class="p-2 rounded-xl bg-primary/10 border border-primary/20"
              >
                <Heart size={18} class="text-primary" />
              </div>
              <div class="flex flex-col">
                <span
                  class="text-xl font-black text-zinc-900 dark:text-white leading-tight"
                >
                  {treeData.edges?.length || 0}
                </span>
                <span
                  class="text-[10px] uppercase tracking-wider text-zinc-500 font-bold"
                >
                  Relasi
                </span>
              </div>
            </div>

            <div class="flex items-center gap-3">
              <div
                class="p-2 rounded-xl bg-emerald-500/10 border border-emerald-500/20"
              >
                <Eye size={18} class="text-emerald-500" />
              </div>
              <div class="flex flex-col">
                <span
                  class="text-xl font-black text-zinc-900 dark:text-white leading-tight capitalize"
                >
                  {family.visibility === "link_only"
                    ? "Link Only"
                    : family.visibility}
                </span>
                <span
                  class="text-[10px] uppercase tracking-wider text-zinc-500 font-bold"
                >
                  Visibilitas
                </span>
              </div>
            </div>
          </div>
        </div>

        <div class="flex flex-col gap-3 self-end md:self-center">
          <Button
            href={`/s/${slug}/tree`}
            class="rounded-2xl h-14 px-8 font-bold gap-3 bg-primary hover:bg-primary/90 text-primary-foreground shadow-lg shadow-primary/20 cursor-pointer text-base"
          >
            <TreeDeciduous size={22} />
            Lihat Pohon Silsilah
            <ArrowRight size={18} />
          </Button>
          <Button
            variant="outline"
            class="rounded-2xl h-12 px-6 font-bold gap-2 border-zinc-200 dark:border-zinc-800 cursor-pointer"
            onclick={copyShareLink}
          >
            <Share2 size={18} />
            Salin Link Berbagi
          </Button>
        </div>
      </div>
    </div>
  </div>

  <!-- Persons Grid -->
  <div class="space-y-6">
    <div class="flex items-center gap-3 px-1">
      <div class="h-8 w-1.5 bg-primary rounded-full"></div>
      <h2 class="text-2xl font-black tracking-tight">
        Daftar Anggota Silsilah
      </h2>
    </div>

    {#if persons.length === 0}
      <div
        class="text-center py-20 rounded-3xl bg-zinc-50 dark:bg-zinc-950 border border-dashed border-zinc-300 dark:border-zinc-700"
      >
        <Users
          size={48}
          class="mx-auto text-zinc-300 dark:text-zinc-600 mb-4"
        />
        <p class="text-zinc-500 font-medium">
          Belum ada data anggota silsilah yang tersedia.
        </p>
      </div>
    {:else}
      <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
        {#each persons as person}
          <div
            class="group p-5 rounded-2xl bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 hover:border-primary/30 hover:shadow-lg hover:shadow-primary/5 transition-all duration-300"
          >
            <div class="flex items-start gap-4">
              <!-- Avatar -->
              <div
                class="flex-shrink-0 w-12 h-12 rounded-xl flex items-center justify-center text-lg font-black transition-all group-hover:scale-110
                {person.is_alive ? '' : 'grayscale opacity-70'}
                {person.gender === 'male'
                  ? 'bg-blue-100 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400'
                  : person.gender === 'female'
                    ? 'bg-pink-100 dark:bg-pink-900/30 text-pink-600 dark:text-pink-400'
                    : 'bg-zinc-100 dark:bg-zinc-800 text-zinc-600 dark:text-zinc-400'}"
              >
                {person.full_name[0]?.toUpperCase() || "?"}
              </div>

              <div class="flex-1 min-w-0">
                <h3
                  class="font-bold text-sm leading-tight truncate group-hover:text-primary transition-colors"
                >
                  {person.full_name}
                </h3>
                {#if person.nickname}
                  <p class="text-xs text-zinc-500 italic truncate">
                    "{person.nickname}"
                  </p>
                {/if}
                <div class="mt-2 space-y-1">
                  {#if person.birth_date}
                    <div
                      class="flex items-center gap-1.5 text-[11px] text-zinc-500"
                    >
                      <Calendar size={10} />
                      <span>Lahir: {formatDate(person.birth_date)}</span>
                    </div>
                  {/if}
                  {#if !person.is_alive}
                    <div
                      class="flex items-center gap-1.5 text-[11px] text-zinc-400"
                    >
                      <Skull size={10} />
                      <span
                        >Wafat: {person.death_date
                          ? formatDate(person.death_date)
                          : "—"}</span
                      >
                    </div>
                  {/if}
                  {#if person.birth_place || person.death_place}
                    <div
                      class="flex items-center gap-1.5 text-[11px] text-zinc-500"
                    >
                      <MapPin size={10} />
                      <span class="truncate">
                        {person.birth_place || "—"}
                        {#if person.death_place}
                          · {person.death_place}
                        {/if}
                      </span>
                    </div>
                  {/if}
                </div>
              </div>
            </div>

            <!-- Status badge -->
            <div class="mt-3 flex items-center justify-between">
              <span
                class="inline-flex items-center gap-1 px-2 py-0.5 rounded-md text-[10px] font-bold uppercase tracking-wider
                {person.is_alive
                  ? 'bg-emerald-50 dark:bg-emerald-900/20 text-emerald-600 dark:text-emerald-400'
                  : 'bg-zinc-100 dark:bg-zinc-800 text-zinc-500 dark:text-zinc-400'}"
              >
                {#if person.is_alive}
                  <Heart size={8} />
                  Hidup
                {:else}
                  <Skull size={8} />
                  Wafat
                {/if}
                {#if person.birth_date}
                  <span class="opacity-70">
                    · {getAge(
                      person.birth_date,
                      person.is_alive,
                      person.death_date,
                    )}
                  </span>
                {/if}
              </span>
              <span
                class="text-[10px] font-bold uppercase tracking-wider
                {person.gender === 'male'
                  ? 'text-blue-500'
                  : person.gender === 'female'
                    ? 'text-pink-500'
                    : 'text-zinc-400'}"
              >
                {person.gender === "male"
                  ? "Laki-laki"
                  : person.gender === "female"
                    ? "Perempuan"
                    : "Lainnya"}
              </span>
            </div>
          </div>
        {/each}
      </div>
    {/if}
  </div>

  <!-- CTA Section -->
  <div
    class="text-center p-8 rounded-3xl bg-zinc-50 dark:bg-zinc-950 border border-zinc-200 dark:border-zinc-800"
  >
    <p class="text-zinc-500 text-sm font-medium mb-4">
      Ingin membuat pohon silsilah keluarga Anda sendiri?
    </p>
    <Button
      href={page.data.user ? "/dashboard" : "/login"}
      class="rounded-2xl h-12 px-8 font-bold gap-2 bg-primary hover:bg-primary/90 text-primary-foreground cursor-pointer"
    >
      <ExternalLink size={18} />
      {page.data.user ? "Ke Dashboard" : "Masuk Sekarang"}
    </Button>
  </div>
</div>
