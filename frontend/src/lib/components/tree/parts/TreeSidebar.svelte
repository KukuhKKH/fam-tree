<script lang="ts">
  import {
    Plus,
    Calendar,
    MapPin,
    HeartOff,
    ChevronRight,
  } from "lucide-svelte";
  import { Button } from "$lib/components/ui/button/index.js";
  import { Badge } from "$lib/components/ui/badge/index.js";
  import { goto } from "$app/navigation";
  import EditPersonDialog from "$lib/components/persons/EditPersonDialog.svelte";

  let {
    selectedPersonId = $bindable(),
    selectedPerson,
    relatives,
    familySlug,
    network,
    getGenderColor,
    getGenderLabel,
    readOnly = false,
  } = $props();

  function selectAnotherPerson(id: number) {
    selectedPersonId = id;
    network?.focus(id, {
      scale: 1,
      animation: { duration: 1000, easingFunction: "easeInOutQuart" },
    });
    network?.selectNodes([id]);
  }
</script>

{#if selectedPerson}
  <div
    class="absolute inset-4 sm:inset-auto sm:top-8 sm:right-8 sm:bottom-8 sm:w-[26rem] flex flex-col pointer-events-none z-50"
  >
    <div
      class="bg-white/90 dark:bg-zinc-900/90 backdrop-blur-3xl border border-zinc-200/50 dark:border-zinc-800/50 rounded-[2rem] sm:rounded-[3rem] shadow-[0_50px_100px_-20px_rgba(0,0,0,0.15)] overflow-hidden flex flex-col pointer-events-auto ring-1 ring-zinc-950/5 dark:ring-white/5 animate-in fade-in slide-in-from-bottom-8 sm:slide-in-from-right-12 duration-700 h-full"
    >
      <!-- Header Image & Banner -->
      <div class="relative h-24 bg-zinc-100 dark:bg-zinc-800 rounded-t-[3rem]">
        <div
          class="absolute inset-0 bg-gradient-to-br {selectedPerson.gender ===
          'male'
            ? 'from-blue-600/20 via-primary/10'
            : 'from-rose-500/20 via-primary/10'} to-transparent rounded-t-[3rem] overflow-hidden"
        >
          <!-- Decorative Shapes inside clipped area -->
          <div
            class="absolute -top-10 -left-10 size-40 rounded-full bg-white/20 blur-3xl"
          ></div>
          <div
            class="absolute top-0 right-0 size-32 bg-primary/10 blur-2xl"
          ></div>
        </div>

        <button
          class="absolute top-6 right-6 z-30 size-10 rounded-full bg-white/50 dark:bg-zinc-800/50 backdrop-blur-md hover:bg-white dark:hover:bg-zinc-700 text-zinc-600 dark:text-zinc-300 transition-all active:scale-90 flex items-center justify-center border border-white/20 cursor-pointer shadow-sm"
          onclick={() => (selectedPersonId = null)}
        >
          <Plus size={20} class="rotate-45" />
        </button>

        <!-- Avatar -->
        <div class="absolute -bottom-10 left-8 flex items-end gap-5 z-20">
          <div
            class="size-28 p-1.5 rounded-[2.2rem] bg-white dark:bg-zinc-900 shadow-[0_20px_50px_rgba(0,0,0,0.1)] ring-1 ring-zinc-100 dark:ring-zinc-800"
          >
            <div class="w-full h-full rounded-[1.8rem] overflow-hidden">
              <img
                src={selectedPerson.photo_url ||
                  `https://ui-avatars.com/api/?name=${encodeURIComponent(selectedPerson.full_name)}&background=${selectedPerson.gender === "male" ? "0ea5e9" : "f43f5e"}&color=fff&bold=true&size=128`}
                alt={selectedPerson.full_name}
                class="w-full h-full object-cover hover:scale-110 transition-transform duration-700"
              />
            </div>
          </div>
        </div>
      </div>

      <!-- Content Area -->
      <div
        class="px-8 pt-14 pb-8 flex-1 overflow-y-auto space-y-8 custom-scrollbar"
      >
        <!-- Name & Primary Info -->
        <div class="space-y-3">
          <div class="space-y-1">
            <h3
              class="text-2xl font-black text-zinc-900 dark:text-white leading-tight tracking-tight"
            >
              {selectedPerson.full_name}
            </h3>
            {#if selectedPerson.nickname}
              <p class="text-base font-bold text-primary italic">
                "{selectedPerson.nickname}"
              </p>
            {/if}
          </div>

          <div class="flex flex-wrap gap-2">
            <Badge
              variant="outline"
              class="rounded-full px-3 py-0.5 text-[10px] font-black uppercase tracking-widest border-{getGenderColor(
                selectedPerson.gender,
              )}/20 bg-{getGenderColor(
                selectedPerson.gender,
              )}/10 text-{getGenderColor(selectedPerson.gender)}"
            >
              {getGenderLabel(selectedPerson.gender)}
            </Badge>
            <Badge
              variant="outline"
              class="rounded-full px-3 py-0.5 text-[10px] font-black uppercase tracking-widest {selectedPerson.is_alive
                ? 'bg-emerald-500/10 text-emerald-600 border-emerald-500/20'
                : 'bg-zinc-500/10 text-zinc-600 border-zinc-500/20'}"
            >
              {selectedPerson.is_alive ? "Hidup" : "Wafat"}
            </Badge>
          </div>
        </div>

        <!-- Quick Stats Grid -->
        <div class="grid grid-cols-2 gap-4">
          <div
            class="p-5 rounded-[2rem] bg-zinc-50 dark:bg-zinc-800/40 border border-zinc-100 dark:border-zinc-800/50 space-y-2"
          >
            <div class="flex items-center gap-2 text-zinc-400">
              <Calendar size={14} />
              <span class="text-[9px] font-black uppercase tracking-widest"
                >Kelahiran</span
              >
            </div>
            <p class="text-sm font-black text-zinc-800 dark:text-zinc-100">
              {selectedPerson.birth_date
                ? new Date(selectedPerson.birth_date).getFullYear()
                : "—"}
            </p>
          </div>
          <div
            class="p-5 rounded-[2rem] bg-zinc-50 dark:bg-zinc-800/40 border border-zinc-100 dark:border-zinc-800/50 space-y-2"
          >
            <div class="flex items-center gap-2 text-zinc-400">
              <MapPin size={14} />
              <span class="text-[9px] font-black uppercase tracking-widest"
                >Tempat</span
              >
            </div>
            <p
              class="text-sm font-black text-zinc-800 dark:text-zinc-100 truncate"
            >
              {selectedPerson.birth_place || "—"}
            </p>
          </div>
        </div>

        <!-- Detailed Info Section -->
        <div class="space-y-6">
          <div class="flex items-center gap-4">
            <div class="h-px flex-1 bg-zinc-100 dark:bg-zinc-800"></div>
            <span
              class="text-[10px] font-black uppercase tracking-[0.2em] text-zinc-400"
              >Arsip Personal</span
            >
            <div class="h-px flex-1 bg-zinc-100 dark:bg-zinc-800"></div>
          </div>

          <div class="space-y-5">
            <!-- Born Info -->
            <div class="flex gap-4 group">
              <div
                class="size-10 rounded-2xl bg-emerald-500/10 text-emerald-600 flex items-center justify-center shrink-0 border border-emerald-500/10"
              >
                <Calendar size={18} />
              </div>
              <div class="space-y-1 pt-0.5">
                <p
                  class="text-[14px] font-black text-zinc-800 dark:text-zinc-200"
                >
                  {#if selectedPerson.birth_date}
                    {new Date(selectedPerson.birth_date).toLocaleDateString(
                      "id-ID",
                      { dateStyle: "long" },
                    )}
                  {:else}
                    Tanggal lahir tidak dicatat
                  {/if}
                </p>
                <p
                  class="text-[12px] font-bold text-zinc-400 flex items-center gap-1"
                >
                  <MapPin size={12} />
                  {selectedPerson.birth_place || "Lokasi tidak diketahui"}
                </p>
              </div>
            </div>

            <!-- Pass Away Info -->
            {#if !selectedPerson.is_alive}
              <div class="flex gap-4 group">
                <div
                  class="size-10 rounded-2xl bg-zinc-500/10 text-zinc-500 flex items-center justify-center shrink-0 border border-zinc-500/10"
                >
                  <HeartOff size={18} />
                </div>
                <div class="space-y-1 pt-0.5">
                  <p
                    class="text-[14px] font-black text-zinc-800 dark:text-zinc-200"
                  >
                    {#if selectedPerson.death_date}
                      {new Date(selectedPerson.death_date).toLocaleDateString(
                        "id-ID",
                        { dateStyle: "long" },
                      )}
                    {:else}
                      Wafat (Tanggal tidak dicatat)
                    {/if}
                  </p>
                  <p
                    class="text-[12px] font-bold text-zinc-400 flex items-center gap-1"
                  >
                    <MapPin size={12} />
                    {selectedPerson.death_place || "Lokasi tidak diketahui"}
                  </p>
                </div>
              </div>
            {/if}
          </div>
        </div>

        <!-- Relatives Section -->
        {#if relatives.length > 0}
          <div class="space-y-4">
            <div class="flex items-center gap-4">
              <div class="h-px flex-1 bg-zinc-100 dark:bg-zinc-800"></div>
              <span
                class="text-[10px] font-black uppercase tracking-[0.2em] text-zinc-400"
                >Hubungan Terdekat</span
              >
              <div class="h-px flex-1 bg-zinc-100 dark:bg-zinc-800"></div>
            </div>

            <div class="grid grid-cols-1 gap-2">
              {#each relatives as rel}
                <button
                  class="flex items-center gap-3 p-3 rounded-2xl bg-zinc-50 dark:bg-zinc-800/40 border border-zinc-100 dark:border-zinc-800/50 hover:bg-white dark:hover:bg-zinc-800 transition-all group cursor-pointer"
                  onclick={() => selectAnotherPerson(rel.person.id)}
                >
                  <div
                    class="size-10 rounded-xl overflow-hidden ring-1 ring-zinc-200 dark:ring-zinc-700"
                  >
                    <img
                      src={rel.person.photo_url ||
                        `https://ui-avatars.com/api/?name=${encodeURIComponent(rel.person.full_name)}&background=${rel.person.gender === "male" ? "0ea5e9" : "f43f5e"}&color=fff&bold=true&size=64`}
                      alt={rel.person.full_name}
                      class="w-full h-full object-cover"
                    />
                  </div>
                  <div class="flex-1 min-w-0 text-left">
                    <p
                      class="text-[13px] font-black text-zinc-800 dark:text-zinc-200 truncate group-hover:text-primary transition-colors"
                    >
                      {rel.person.full_name}
                    </p>
                    <p
                      class="text-[10px] font-bold text-zinc-400 uppercase tracking-tight"
                    >
                      {rel.displayLabel}
                    </p>
                  </div>
                  <ChevronRight
                    size={14}
                    class="text-zinc-300 group-hover:text-primary group-hover:translate-x-1 transition-all"
                  />
                </button>
              {/each}
            </div>
          </div>
        {/if}

        <!-- Action Buttons -->
        {#if !readOnly}
          <div class="pt-6 flex flex-col gap-3">
            <Button
              class="w-full rounded-[1.8rem] h-14 font-black gap-2 transition-all hover:shadow-2xl hover:shadow-primary/20 cursor-pointer text-base active:scale-95"
              onclick={() =>
                goto(`/families/${familySlug}/persons/${selectedPerson.id}`)}
            >
              Lihat Profil Lengkap
              <ChevronRight size={18} />
            </Button>

            <EditPersonDialog person={selectedPerson} {familySlug} />
          </div>
        {/if}
      </div>
    </div>
  </div>
{/if}

<style>
  .custom-scrollbar::-webkit-scrollbar {
    width: 4px;
  }
  .custom-scrollbar::-webkit-scrollbar-track {
    background: transparent;
  }
  .custom-scrollbar::-webkit-scrollbar-thumb {
    background: rgba(0, 0, 0, 0.05);
    border-radius: 10px;
  }
</style>
