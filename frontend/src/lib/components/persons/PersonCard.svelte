<script lang="ts">
  import { Button } from "$lib/components/ui/button/index.js";
  import { Badge } from "$lib/components/ui/badge/index.js";
  import {
    User,
    Heart,
    HeartOff,
    Calendar,
    MapPin,
    ExternalLink,
    Trash2,
  } from "lucide-svelte";
  import { cn } from "$lib/utils.js";
  import { apiFetch } from "$lib/api";
  import { toast } from "svelte-sonner";
  import { invalidateAll } from "$app/navigation";

  import ConfirmDialog from "$lib/components/ui/ConfirmDialog.svelte";

  let { person, familySlug, canManage } = $props();

  let deleting = $state(false);
  let confirmOpen = $state(false);

  async function deletePerson() {
    deleting = true;
    try {
      await apiFetch(`/families/${familySlug}/persons/${person.id}`, {
        method: "DELETE",
      });
      toast.success(`${person.full_name} berhasil dihapus dari silsilah`);
      await invalidateAll();
    } catch (err: any) {
      toast.error(err.message || "Gagal menghapus data person");
    } finally {
      deleting = false;
      confirmOpen = false;
    }
  }

  const birthYear = $derived(
    person.birth_date ? new Date(person.birth_date).getFullYear() : "?",
  );
  const deathYear = $derived(
    !person.is_alive && person.death_date
      ? new Date(person.death_date).getFullYear()
      : person.is_alive
        ? ""
        : "?",
  );
</script>

<div
  class="group relative p-6 rounded-3xl bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 shadow-sm hover:shadow-xl hover:-translate-y-1 transition-all duration-300"
>
  <div class="flex items-start justify-between mb-4">
    <div class="flex items-center gap-4">
      <div
        class={cn(
          "w-12 h-12 rounded-2xl flex items-center justify-center text-xl font-bold uppercase",
          person.gender === "male"
            ? "bg-blue-100/50 text-blue-600 dark:bg-blue-900/20"
            : "bg-rose-100/50 text-rose-600 dark:bg-rose-900/20",
        )}
      >
        {person.full_name.substring(0, 2)}
      </div>
      <div class="space-y-1">
        <h4
          class="font-bold text-lg group-hover:text-primary transition-colors line-clamp-1"
        >
          {person.full_name}
        </h4>
        {#if person.nickname}
          <p class="text-xs text-zinc-500 italic font-medium">
            "{person.nickname}"
          </p>
        {/if}
      </div>
    </div>

    <div class="flex items-center gap-1">
      {#if person.is_alive}
        <div
          class="px-2 py-0.5 rounded-full bg-emerald-100/50 dark:bg-emerald-900/20 text-emerald-600 dark:text-emerald-400 text-[10px] font-bold uppercase tracking-wider flex items-center gap-1"
        >
          <Heart size={10} fill="currentColor" /> Hidup
        </div>
      {:else}
        <div
          class="px-2 py-0.5 rounded-full bg-zinc-100/50 dark:bg-zinc-800 text-zinc-500 text-[10px] font-bold uppercase tracking-wider flex items-center gap-1"
        >
          <HeartOff size={10} /> Wafat
        </div>
      {/if}
    </div>
  </div>

  <div class="space-y-3 pt-2">
    <div class="flex items-center gap-3 text-sm text-zinc-500">
      <Calendar size={14} class="text-primary/60 shrink-0" />
      <span class="font-medium">{birthYear} — {deathYear || "Skg"}</span>
    </div>

    {#if person.birth_place}
      <div class="flex items-center gap-3 text-sm text-zinc-500">
        <MapPin size={14} class="text-primary/60 shrink-0" />
        <span class="truncate">{person.birth_place}</span>
      </div>
    {/if}
  </div>

  <div class="flex items-center gap-2 mt-6">
    <Button
      variant="secondary"
      class="flex-1 rounded-xl h-10 text-xs font-bold gap-2 cursor-pointer"
      href={`/families/${familySlug}/persons/${person.id}`}
    >
      <ExternalLink size={14} />
      Detail Profil
    </Button>

    {#if canManage}
      <Button
        variant="ghost"
        size="icon"
        class="rounded-xl h-10 w-10 text-zinc-400 hover:text-rose-500 hover:bg-rose-50 dark:hover:bg-rose-950/20 cursor-pointer transition-colors"
        disabled={deleting}
        onclick={() => (confirmOpen = true)}
      >
        <Trash2 size={16} />
      </Button>
    {/if}
  </div>
</div>

<ConfirmDialog
  bind:open={confirmOpen}
  loading={deleting}
  onConfirm={deletePerson}
  title={`Hapus ${person.full_name}?`}
  description="Data person dan semua hubungan keluarganya akan dihapus permanen dari silsilah ini."
/>
