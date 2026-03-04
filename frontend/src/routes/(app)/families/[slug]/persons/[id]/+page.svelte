<script lang="ts">
  import { Button } from "$lib/components/ui/button/index.js";
  import { Badge } from "$lib/components/ui/badge/index.js";
  import {
    ArrowLeft,
    Heart,
    HeartOff,
    Calendar,
    MapPin,
    Edit,
    Share2,
    Plus,
    Shield,
    History,
    CalendarDays,
    Info,
    ChevronRight,
  } from "lucide-svelte";
  import { cn } from "$lib/utils.js";
  import AddRelationshipDialog from "$lib/components/relationships/AddRelationshipDialog.svelte";
  import EditPersonDialog from "$lib/components/persons/EditPersonDialog.svelte";

  const { data } = $props();
  const person = $derived(data.person);
  const family = $derived(data.family);
  const persons = $derived(data.persons);
  const relationships = $derived(data.relationships || []);
  const slug = $derived(data.slug);

  const canEdit = $derived(family.my_role !== "viewer");

  function calculateAge(
    birthDate: string,
    isAlive: boolean,
    deathDate?: string,
  ) {
    if (!birthDate) return null;
    const start = new Date(birthDate);
    const end = isAlive
      ? new Date()
      : deathDate
        ? new Date(deathDate)
        : new Date();
    let age = end.getFullYear() - start.getFullYear();
    const m = end.getMonth() - start.getMonth();
    if (m < 0 || (m === 0 && end.getDate() < start.getDate())) age--;
    return age;
  }

  const age = $derived(
    calculateAge(person.birth_date, person.is_alive, person.death_date),
  );

  // Compute relationships for this person
  const personRelationships = $derived(
    relationships
      .filter(
        (r: any) => r.person_a_id === person.id || r.person_b_id === person.id,
      )
      .map((r: any) => {
        const isA = r.person_a_id === person.id;
        const relatedId = isA ? r.person_b_id : r.person_a_id;
        const relatedPerson = persons.find((p: any) => p.id === relatedId);
        if (!relatedPerson) return null;

        let label = "";
        if (r.relationship_type === "parent_child") {
          if (isA) {
            // person is parent → related is child
            label = "Anak";
          } else {
            // person is child → related is parent
            label = relatedPerson.gender === "male" ? "Ayah" : "Ibu";
          }
        } else if (r.relationship_type === "spouse") {
          label = relatedPerson.gender === "male" ? "Suami" : "Istri";
        } else if (r.relationship_type === "sibling") {
          label = "Saudara";
        }

        return { ...r, relatedPerson, label };
      })
      .filter(Boolean),
  );
</script>

<svelte:head>
  <title>{person.full_name} | Silsilah</title>
</svelte:head>

<div class="space-y-8 mx-auto space-y-10 pb-20">
  <!-- Premium Breadcrumbs / Back Navigation -->
  <nav class="flex items-center gap-2 text-sm font-medium text-zinc-500 mb-2">
    <a
      href="/families"
      class="hover:text-primary transition-colors flex items-center gap-1"
    >
      Keluarga
    </a>
    <ChevronRight size={14} />
    <a href={`/families/${slug}`} class="hover:text-primary transition-colors">
      {family.name}
    </a>
    <ChevronRight size={14} />
    <span class="text-zinc-900 dark:text-zinc-100 font-bold truncate">
      {person.full_name}
    </span>
  </nav>

  <!-- Main Hero Profile Section -->
  <div class="relative">
    <!-- Background Decoration -->
    <div
      class="absolute inset-0 bg-gradient-to-br from-primary/10 via-transparent to-primary/5 rounded-[3rem] -z-10 blur-3xl opacity-50"
    ></div>

    <div
      class="rounded-[3rem] bg-white dark:bg-zinc-900 border border-zinc-200/60 dark:border-zinc-800/60 overflow-hidden shadow-2xl shadow-zinc-200/20 dark:shadow-none"
    >
      <!-- Cover Banner -->
      <div class="relative h-48 sm:h-64 overflow-hidden">
        <div
          class="absolute inset-0 bg-gradient-to-tr {person.gender === 'male'
            ? 'from-blue-600/20 via-primary/10 to-transparent'
            : 'from-rose-500/20 via-primary/10 to-transparent'}"
        ></div>
        <div
          class="absolute inset-0 backdrop-blur-[100px] opacity-40 dark:opacity-20"
        ></div>

        <!-- Floatting ornaments -->
        <div
          class="absolute top-10 right-10 w-32 h-32 rounded-full bg-primary/20 blur-3xl animate-pulse"
        ></div>
        <div
          class="absolute bottom-5 left-20 w-24 h-24 rounded-full bg-primary/10 blur-2xl delay-75"
        ></div>
      </div>

      <div class="px-6 sm:px-12 pb-12 -mt-20 relative">
        <div class="flex flex-col md:flex-row gap-8 items-end">
          <!-- Avatar wrapper -->
          <div class="relative group">
            <div
              class="absolute -inset-1.5 bg-gradient-to-tr from-primary to-primary-foreground rounded-[2.5rem] blur opacity-20 group-hover:opacity-40 transition duration-1000 group-hover:duration-200"
            ></div>
            <div
              class="relative size-32 sm:size-44 p-2 rounded-[2.2rem] bg-white dark:bg-zinc-900 shadow-2xl ring-1 ring-zinc-100 dark:ring-zinc-800"
            >
              <div
                class={cn(
                  "w-full h-full rounded-[1.8rem] flex items-center justify-center text-5xl sm:text-7xl font-black uppercase transition-transform duration-500 group-hover:scale-105",
                  person.gender === "male"
                    ? "bg-blue-50 text-blue-600 dark:bg-blue-950/40"
                    : "bg-rose-50 text-rose-500 dark:bg-rose-950/40",
                )}
              >
                {person.full_name?.substring(0, 1) || "?"}
              </div>
            </div>

            {#if person.is_alive}
              <div
                class="absolute bottom-4 right-4 size-6 rounded-full bg-emerald-500 ring-4 ring-white dark:ring-zinc-900 animate-pulse"
                title="Hidup"
              ></div>
            {:else}
              <div
                class="absolute bottom-4 right-4 size-6 rounded-full bg-zinc-400 ring-4 ring-white dark:ring-zinc-900 flex items-center justify-center"
                title="Wafat"
              >
                <HeartOff size={10} class="text-white" />
              </div>
            {/if}
          </div>

          <!-- Title & Basic Info -->
          <div class="flex-1 space-y-5 pb-2">
            <div class="space-y-2">
              <div class="flex flex-wrap items-center gap-3">
                <h1
                  class="text-3xl sm:text-5xl font-black tracking-tight text-zinc-900 dark:text-white"
                >
                  {person.full_name}
                </h1>
                <Badge
                  variant="outline"
                  class="rounded-full px-4 py-1.5 font-black border-primary/20 bg-primary/5 text-primary tracking-wide text-xs uppercase"
                >
                  {person.gender === "male" ? "Laki-Laki" : "Perempuan"}
                </Badge>
              </div>
              {#if person.nickname}
                <p
                  class="text-xl text-zinc-500 font-medium tracking-tight italic"
                >
                  "{person.nickname}"
                </p>
              {/if}
            </div>

            <!-- Quick Meta Info -->
            <div class="flex flex-wrap items-center gap-y-3 gap-x-8">
              {#if age !== null}
                <div class="flex items-center gap-2.5">
                  <div
                    class="p-2 rounded-xl bg-orange-50 dark:bg-orange-950/20 text-orange-600"
                  >
                    <CalendarDays size={18} />
                  </div>
                  <div class="flex flex-col">
                    <span
                      class="text-[10px] font-black uppercase tracking-widest text-zinc-400 leading-none mb-1"
                    >
                      Usia
                    </span>
                    <span class="font-bold text-zinc-900 dark:text-zinc-100"
                      >{age} Tahun {!person.is_alive ? "(Wafat)" : ""}</span
                    >
                  </div>
                </div>
              {/if}

              <div class="flex items-center gap-2.5">
                <div class="p-2 rounded-xl bg-primary/10 text-primary">
                  <Shield size={18} />
                </div>
                <div class="flex flex-col">
                  <span
                    class="text-[10px] font-black uppercase tracking-widest text-zinc-400 leading-none mb-1"
                  >
                    Akses
                  </span>
                  <span
                    class="font-bold text-zinc-900 dark:text-zinc-100 capitalize"
                  >
                    {family.my_role}
                  </span>
                </div>
              </div>

              <div class="flex items-center gap-2.5">
                <div
                  class="p-2 rounded-xl bg-emerald-50 dark:bg-emerald-950/20 text-emerald-600"
                >
                  <Heart size={18} />
                </div>
                <div class="flex flex-col">
                  <span
                    class="text-[10px] font-black uppercase tracking-widest text-zinc-400 leading-none mb-1"
                  >
                    Status
                  </span>
                  <span class="font-bold text-zinc-900 dark:text-zinc-100"
                    >{person.is_alive ? "Sehat / Hidup" : "Almarhum"}</span
                  >
                </div>
              </div>
            </div>
          </div>

          <!-- Action Buttons -->
          <div class="flex items-center gap-3 self-center md:self-end">
            {#if canEdit}
              <EditPersonDialog familySlug={slug} {person} />
            {/if}
          </div>
        </div>
      </div>
    </div>
  </div>

  <!-- Content Grid -->
  <div class="grid grid-cols-1 lg:grid-cols-12 gap-10">
    <!-- Left Column: Details -->
    <div class="lg:col-span-8 space-y-10">
      <!-- Life Timeline / Detail Information -->
      <section
        class="rounded-[2.5rem] bg-white dark:bg-zinc-900 border border-zinc-200/60 dark:border-zinc-800/60 p-8 sm:p-10 shadow-sm space-y-10"
      >
        <div class="flex items-center justify-between">
          <h3
            class="text-2xl font-black text-zinc-900 dark:text-white flex items-center gap-3"
          >
            <span
              class="size-10 rounded-2xl bg-primary/10 text-primary flex items-center justify-center"
            >
              <Info size={22} />
            </span>
            Detail Personal
          </h3>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-x-12 gap-y-10">
          <!-- Birth Side -->
          <div class="space-y-6">
            <div>
              <span
                class="text-[10px] font-black uppercase tracking-[0.2em] text-zinc-400 block mb-4"
                >Kelahiran</span
              >
              <div class="flex gap-4">
                <div
                  class="size-14 rounded-2xl bg-zinc-50 dark:bg-zinc-950 border border-zinc-100 dark:border-zinc-800 flex items-center justify-center text-primary"
                >
                  <Calendar size={28} />
                </div>
                <div class="space-y-1 pt-1">
                  <p class="text-xl font-bold text-zinc-900 dark:text-zinc-100">
                    {person.birth_date
                      ? new Date(person.birth_date).toLocaleDateString(
                          "id-ID",
                          { day: "numeric", month: "long", year: "numeric" },
                        )
                      : "—"}
                  </p>
                  <p
                    class="text-zinc-500 font-medium flex items-center gap-1.5"
                  >
                    <MapPin size={14} />
                    {person.birth_place || "Tempat lahir tidak dicatat"}
                  </p>
                </div>
              </div>
            </div>
          </div>

          <!-- Death Side -->
          {#if !person.is_alive}
            <div class="space-y-6">
              <div>
                <span
                  class="text-[10px] font-black uppercase tracking-[0.2em] text-rose-400 block mb-4"
                  >Wafat</span
                >
                <div class="flex gap-4">
                  <div
                    class="size-14 rounded-2xl bg-rose-50/50 dark:bg-rose-950/20 border border-rose-100/50 dark:border-rose-900/20 flex items-center justify-center text-rose-500"
                  >
                    <HeartOff size={28} />
                  </div>
                  <div class="space-y-1 pt-1">
                    <p
                      class="text-xl font-bold text-zinc-900 dark:text-zinc-100 text-rose-600 dark:text-rose-400"
                    >
                      {person.death_date
                        ? new Date(person.death_date).toLocaleDateString(
                            "id-ID",
                            { day: "numeric", month: "long", year: "numeric" },
                          )
                        : "—"}
                    </p>
                    <p
                      class="text-zinc-500 font-medium flex items-center gap-1.5"
                    >
                      <MapPin size={14} />
                      {person.death_place || "Tempat wafat tidak dicatat"}
                    </p>
                  </div>
                </div>
              </div>
            </div>
          {/if}
        </div>

        <!-- Biography -->
        <div class="pt-6 border-t border-zinc-100 dark:border-zinc-800">
          <span
            class="text-[10px] font-black uppercase tracking-[0.2em] text-zinc-400 block mb-4"
            >Biografi & Catatan</span
          >
          <div class="relative group">
            <div
              class="absolute -inset-1 bg-gradient-to-r from-primary/5 to-transparent rounded-3xl blur opacity-0 group-hover:opacity-100 transition duration-500"
            ></div>
            <div
              class="relative p-8 rounded-3xl bg-zinc-50/50 dark:bg-zinc-950/30 border border-zinc-100 dark:border-zinc-800 text-zinc-700 dark:text-zinc-300 leading-relaxed text-lg italic tracking-tight font-medium"
            >
              {person.bio || "Belum ada narasi biografi untuk person ini."}
            </div>
          </div>
        </div>
      </section>

      <!-- Relationships -->
      <section
        class="rounded-[2.5rem] bg-white dark:bg-zinc-900 border border-zinc-200/60 dark:border-zinc-800/60 p-8 sm:p-10 shadow-sm space-y-8"
      >
        <div class="flex items-center justify-between">
          <h3
            class="text-2xl font-black text-zinc-900 dark:text-white flex items-center gap-3"
          >
            <span
              class="size-10 rounded-2xl bg-primary/10 text-primary flex items-center justify-center"
            >
              <Heart size={22} />
            </span>
            Hubungan Keluarga
          </h3>
          {#if canEdit}
            <AddRelationshipDialog
              currentPersonId={person.id}
              familySlug={slug}
              {persons}
            />
          {/if}
        </div>

        {#if personRelationships.length === 0}
          <div
            class="bg-zinc-50 dark:bg-zinc-950 border border-dotted border-zinc-200 dark:border-zinc-800 p-12 rounded-[2rem] flex flex-col items-center text-center space-y-4"
          >
            <div
              class="p-4 rounded-full bg-white dark:bg-zinc-900 shadow-sm ring-1 ring-zinc-100 dark:ring-zinc-800"
            >
              <Plus size={32} class="text-zinc-400" />
            </div>
            <div>
              <p
                class="text-lg font-bold text-zinc-800 dark:text-zinc-200 mb-1"
              >
                Belum Ada Hubungan
              </p>
              <p class="text-sm text-zinc-500 max-w-sm mx-auto">
                Klik tombol di atas untuk menambahkan hubungan keluarga seperti
                orang tua, pasangan, atau saudara.
              </p>
            </div>
          </div>
        {:else}
          <div class="grid grid-cols-1 gap-4">
            {#each personRelationships as rel}
              <a
                href={`/families/${slug}/persons/${rel.relatedPerson.id}`}
                class="flex items-center gap-5 p-5 rounded-[1.8rem] bg-zinc-50/70 dark:bg-zinc-950/40 border border-zinc-100 dark:border-zinc-800/50 hover:bg-white dark:hover:bg-zinc-900 hover:shadow-lg hover:shadow-zinc-200/20 dark:hover:shadow-none transition-all group"
              >
                <div
                  class="size-14 rounded-2xl overflow-hidden ring-1 ring-zinc-200 dark:ring-zinc-700 shrink-0 {rel
                    .relatedPerson.gender === 'male'
                    ? 'bg-blue-50 dark:bg-blue-950/30'
                    : 'bg-rose-50 dark:bg-rose-950/30'} flex items-center justify-center text-2xl font-black uppercase"
                >
                  <span
                    class={rel.relatedPerson.gender === "male"
                      ? "text-blue-600"
                      : "text-rose-500"}
                  >
                    {rel.relatedPerson.full_name?.substring(0, 1) || "?"}
                  </span>
                </div>
                <div class="flex-1 min-w-0">
                  <p
                    class="text-base font-black text-zinc-800 dark:text-zinc-200 truncate group-hover:text-primary transition-colors"
                  >
                    {rel.relatedPerson.full_name}
                  </p>
                  <p
                    class="text-xs font-bold text-zinc-400 uppercase tracking-wide"
                  >
                    {rel.label}
                  </p>
                </div>
                <ChevronRight
                  size={18}
                  class="text-zinc-300 group-hover:text-primary group-hover:translate-x-1 transition-all shrink-0"
                />
              </a>
            {/each}
          </div>
        {/if}
      </section>
    </div>

    <!-- Right Column: Meta & Sidebar -->
    <div class="lg:col-span-4 space-y-8">
      <!-- Data History Card -->
      <div
        class="rounded-[2rem] bg-white dark:bg-zinc-900 border border-zinc-200/60 dark:border-zinc-800/60 p-8 shadow-sm space-y-6"
      >
        <h4
          class="font-black text-sm uppercase tracking-widest text-zinc-400 flex items-center gap-2"
        >
          <History size={16} class="text-primary" />
          MetaData & Riwayat
        </h4>

        <div class="space-y-6">
          <div class="flex items-center gap-4">
            <div
              class="size-10 rounded-xl bg-zinc-50 dark:bg-zinc-950 flex items-center justify-center text-zinc-400 border border-zinc-100 dark:border-zinc-800"
            >
              <Plus size={18} />
            </div>
            <div class="flex-1">
              <p
                class="text-[10px] font-black uppercase tracking-widest text-zinc-400 leading-none mb-1"
              >
                Dibuat
              </p>
              <p class="font-bold text-zinc-900 dark:text-zinc-100">
                {new Date(person.created_at).toLocaleDateString("id-ID", {
                  day: "numeric",
                  month: "long",
                  year: "numeric",
                })}
              </p>
            </div>
          </div>

          <div class="flex items-center gap-4">
            <div
              class="size-10 rounded-xl bg-zinc-50 dark:bg-zinc-950 flex items-center justify-center text-zinc-400 border border-zinc-100 dark:border-zinc-800"
            >
              <Edit size={18} />
            </div>
            <div class="flex-1">
              <p
                class="text-[10px] font-black uppercase tracking-widest text-zinc-400 leading-none mb-1"
              >
                Terakhir Update
              </p>
              <p class="font-bold text-zinc-900 dark:text-zinc-100">
                {new Date(person.updated_at).toLocaleDateString("id-ID", {
                  day: "numeric",
                  month: "long",
                  year: "numeric",
                })}
              </p>
            </div>
          </div>
        </div>

        <div class="pt-6 border-t border-zinc-100 dark:border-zinc-800">
          <div class="p-4 rounded-2xl bg-primary/5 border border-primary/10">
            <p class="text-[11px] text-zinc-500 leading-relaxed italic">
              "Setiap detail membantu menyusun puzzle sejarah besar keluarga
              kita agar tidak terlupakan oleh generasi mendatang."
            </p>
          </div>
        </div>
      </div>

      <!-- Quick Tips -->
      <div
        class="relative overflow-hidden rounded-3xl p-7 sm:p-8
         border border-zinc-200/70 dark:border-zinc-800/70
         bg-gradient-to-br from-primary/10 via-white to-white
         dark:from-primary/25 dark:via-zinc-950 dark:to-zinc-950
         shadow-sm"
      >
        <!-- blobs (aman di dua mode) -->
        <div
          class="pointer-events-none absolute -top-24 -right-24 h-56 w-56 rounded-full bg-primary/15 dark:bg-primary/18 blur-3xl"
        ></div>
        <div
          class="pointer-events-none absolute -bottom-24 -left-24 h-56 w-56 rounded-full bg-zinc-900/5 dark:bg-white/6 blur-3xl"
        ></div>

        <!-- subtle pattern -->
        <div
          class="pointer-events-none absolute inset-0 opacity-[0.25] dark:opacity-[0.18]
           [background:radial-gradient(circle_at_1px_1px,rgba(24,24,27,.22)_1px,transparent_1px)]
           dark:[background:radial-gradient(circle_at_1px_1px,rgba(255,255,255,.18)_1px,transparent_1px)]
           [background-size:18px_18px]"
        ></div>

        <div class="relative space-y-3">
          <div class="flex items-center gap-3">
            <div
              class="h-11 w-11 rounded-2xl grid place-items-center
               bg-primary/10 dark:bg-primary/18
               ring-1 ring-primary/20 dark:ring-primary/25"
            >
              <span class="text-xl">💡</span>
            </div>

            <h4
              class="text-lg font-semibold tracking-tight text-zinc-900 dark:text-zinc-100"
            >
              Tips Silsilah
            </h4>
          </div>

          <p
            class="text-sm leading-relaxed text-zinc-600 dark:text-zinc-300 max-w-[52ch]"
          >
            Menambahkan foto dan biografi yang mendetail akan membuat profil ini
            lebih berkesan bagi keturunan di masa depan.
          </p>

          <div class="pt-1">
            <span
              class="inline-flex items-center rounded-full px-3 py-1 text-xs font-semibold
               bg-zinc-900/5 text-zinc-700 ring-1 ring-zinc-900/10
               dark:bg-white/8 dark:text-zinc-200 dark:ring-white/10"
            >
              Saran cepat
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>
