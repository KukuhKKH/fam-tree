<script lang="ts">
  import * as Dialog from "$lib/components/ui/dialog/index.js";
  import { Button } from "$lib/components/ui/button/index.js";
  import { Label } from "$lib/components/ui/label/index.js";
  import {
    Link,
    Loader2,
    Heart,
    Users,
    ArrowRight,
    Check,
    Info,
    Save,
  } from "lucide-svelte";
  import PersonSelect from "$lib/components/persons/PersonSelect.svelte";
  import { apiFetch } from "$lib/api";
  import { toast } from "svelte-sonner";
  import { invalidateAll } from "$app/navigation";

  let {
    familySlug,
    persons = [],
    currentPersonId = null,
    style = "",
  } = $props();

  let open = $state(false);
  let loading = $state(false);

  let personAId = $state("");
  let personBId = $state("");
  let relationshipType = $state("parent_child");

  $effect(() => {
    if (currentPersonId && !personAId) {
      personAId = String(currentPersonId);
    }
  });

  const relationshipOptions = [
    { value: "parent_child", label: "Orang Tua -> Anak" },
    { value: "spouse", label: "Pasangan (Suami/Istri)" },
    { value: "sibling", label: "Saudara Kandung" },
  ];

  async function handleSubmit(e: SubmitEvent) {
    e.preventDefault();
    if (!personAId || !personBId)
      return toast.error("Pilih kedua orang terlebih dahulu");
    if (personAId === personBId)
      return toast.error("Tidak bisa menghubungkan orang yang sama");

    loading = true;
    try {
      await apiFetch(`/families/${familySlug}/relationships`, {
        method: "POST",
        body: JSON.stringify({
          person_a_id: Number(personAId),
          person_b_id: Number(personBId),
          relationship_type: relationshipType,
        }),
      });

      toast.success("Hubungan berhasil ditambahkan");
      open = false;
      personBId = "";
      await invalidateAll();
    } catch (err: any) {
      toast.error(err.message || "Gagal menambahkan hubungan");
    } finally {
      loading = false;
    }
  }

  const selectedPersonA = $derived(
    persons.find((p: any) => String(p.id) === personAId),
  );
  const selectedPersonB = $derived(
    persons.find((p: any) => String(p.id) === personBId),
  );
</script>

<Dialog.Root bind:open>
  <Dialog.Trigger>
    {#snippet child({ props })}
      <Button
        {...props}
        variant="outline"
        {style}
        class="gap-2 rounded-2xl h-12 font-bold border-zinc-200 dark:border-zinc-800 hover:bg-zinc-50 dark:hover:bg-zinc-900 transition-all cursor-pointer shadow-sm"
      >
        <Link size={18} />
        <span>Hubungkan Orang</span>
      </Button>
    {/snippet}
  </Dialog.Trigger>

  <Dialog.Content
    class="sm:max-w-[640px] p-0 overflow-hidden rounded-[2.5rem] border-zinc-200/70 dark:border-zinc-800/70 shadow-2xl"
  >
    <form onsubmit={handleSubmit} class="flex flex-col max-h-[90vh]">
      <!-- Header -->
      <div class="px-8 pt-8 pb-4 bg-white dark:bg-zinc-950">
        <Dialog.Header class="space-y-1">
          <div class="flex items-start gap-4">
            <div
              class="p-3 rounded-2xl bg-primary/10 text-primary ring-1 ring-primary/20 shadow-sm"
            >
              <Link size={24} />
            </div>
            <div class="min-w-0">
              <Dialog.Title class="text-2xl font-black tracking-tight">
                Pembangun Hubungan
              </Dialog.Title>
              <Dialog.Description class="text-sm font-medium text-zinc-500">
                Tentukan koneksi silsilah antar anggota keluarga.
              </Dialog.Description>
            </div>
          </div>
        </Dialog.Header>
      </div>

      <!-- Scrollable Body -->
      <div class="px-8 py-6 overflow-y-auto bg-zinc-50/50 dark:bg-zinc-950/30">
        <div class="grid gap-8">
          <!-- People Selection Card -->
          <div
            class="grid gap-6 p-6 rounded-[2rem] bg-white dark:bg-zinc-950 border border-zinc-200/60 dark:border-zinc-800/60 shadow-sm"
          >
            <div
              class="grid grid-cols-1 md:grid-cols-[1fr,48px,1fr] items-end gap-4"
            >
              <div class="space-y-2">
                <Label
                  class="text-[10px] font-black uppercase tracking-widest text-zinc-400 px-1"
                  >Orang Pertama</Label
                >
                <PersonSelect
                  {persons}
                  bind:value={personAId}
                  placeholder="Pilih Orang"
                />
              </div>

              <div class="flex items-center justify-center pb-1">
                <div
                  class="h-10 w-10 rounded-xl bg-zinc-100 dark:bg-zinc-800 text-zinc-400 grid place-items-center"
                >
                  <ArrowRight size={18} class="md:rotate-0 rotate-90" />
                </div>
              </div>

              <div class="space-y-2">
                <Label
                  class="text-[10px] font-black uppercase tracking-widest text-zinc-400 px-1"
                  >Orang Kedua</Label
                >
                <PersonSelect
                  {persons}
                  bind:value={personBId}
                  placeholder="Pilih Orang"
                />
              </div>
            </div>
          </div>

          <!-- Relationship Type Selection -->
          <div class="space-y-4">
            <Label
              class="text-xs font-black uppercase tracking-widest text-zinc-400 px-1"
              >Jenis Hubungan</Label
            >
            <div class="grid grid-cols-1 gap-3">
              {#each relationshipOptions as opt}
                <button
                  type="button"
                  onclick={() => (relationshipType = opt.value)}
                  class="group relative flex items-center gap-4 p-4 rounded-2xl border transition-all text-left
                    {relationshipType === opt.value
                    ? 'bg-primary/5 border-primary shadow-sm'
                    : 'bg-white dark:bg-zinc-950 border-zinc-200 dark:border-zinc-800 hover:border-zinc-300 dark:hover:border-zinc-700'}"
                >
                  <div
                    class="size-12 rounded-xl flex items-center justify-center transition-colors
                    {relationshipType === opt.value
                      ? 'bg-primary text-white'
                      : 'bg-zinc-100 dark:bg-zinc-900 text-zinc-400 group-hover:text-primary'}"
                  >
                    {#if opt.value === "parent_child"}
                      <Users size={20} />
                    {:else if opt.value === "spouse"}
                      <Heart size={20} />
                    {:else}
                      <Link size={20} />
                    {/if}
                  </div>

                  <div class="flex-1">
                    <p
                      class="font-bold text-sm {relationshipType === opt.value
                        ? 'text-primary'
                        : 'text-zinc-900 dark:text-white'}"
                    >
                      {opt.label}
                    </p>
                    <p class="text-xs text-zinc-500 leading-tight">
                      {#if opt.value === "parent_child"}
                        Hubungan vertikal antar generasi.
                      {:else if opt.value === "spouse"}
                        Hubungan horizontal pasutri.
                      {:else}
                        Hubungan mendatar satu keturunan.
                      {/if}
                    </p>
                  </div>

                  {#if relationshipType === opt.value}
                    <div
                      class="size-6 rounded-full bg-primary text-white flex items-center justify-center"
                    >
                      <Check size={14} strokeWidth={3} />
                    </div>
                  {/if}
                </button>
              {/each}
            </div>
          </div>

          <!-- Preview & Confirmation -->
          <div
            class="p-5 rounded-2xl bg-zinc-900 text-white shadow-xl space-y-3"
          >
            <div
              class="flex items-center gap-2 text-[10px] font-black uppercase tracking-[0.2em] text-zinc-500"
            >
              <Info size={12} /> Konfirmasi Hubungan
            </div>
            <p class="text-sm leading-relaxed">
              {#if relationshipType === "parent_child"}
                <span class="font-black text-primary"
                  >{selectedPersonA?.full_name || "..."}</span
                >
                adalah <strong>Orang Tua</strong> dari
                <span class="font-black text-primary"
                  >{selectedPersonB?.full_name || "..."}</span
                >.
              {:else if relationshipType === "spouse"}
                <span class="font-black text-primary"
                  >{selectedPersonA?.full_name || "..."}</span
                >
                dan
                <span class="font-black text-primary"
                  >{selectedPersonB?.full_name || "..."}</span
                >
                adalah <strong>Pasangan</strong>.
              {:else}
                <span class="font-black text-primary"
                  >{selectedPersonA?.full_name || "..."}</span
                >
                dan
                <span class="font-black text-primary"
                  >{selectedPersonB?.full_name || "..."}</span
                >
                adalah <strong>Saudara Kandung</strong>.
              {/if}
            </p>
          </div>
        </div>
      </div>

      <!-- Footer Buttons -->
      <div
        class="p-6 bg-white dark:bg-zinc-950 border-t border-zinc-100 dark:border-zinc-800 sm:px-8 sm:py-6"
      >
        <Dialog.Footer class="flex flex-row gap-3">
          <Button
            type="button"
            variant="ghost"
            onclick={() => (open = false)}
            class="hidden sm:flex rounded-2xl h-12 px-6 font-bold text-zinc-500 cursor-pointer"
          >
            Batal
          </Button>
          <Button
            type="submit"
            class="flex-1 rounded-2xl h-12 font-bold text-base shadow-xl shadow-primary/20 cursor-pointer transition-all active:scale-[0.98]"
            disabled={loading}
          >
            {#if loading}
              <Loader2 class="mr-2 h-5 w-5 animate-spin" />
              Menyimpan...
            {:else}
              <Save size={20} class="mr-2" />
              Simpan Hubungan
            {/if}
          </Button>
        </Dialog.Footer>
      </div>
    </form>
  </Dialog.Content>
</Dialog.Root>
