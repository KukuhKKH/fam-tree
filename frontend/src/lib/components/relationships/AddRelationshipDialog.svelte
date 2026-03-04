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
        class="gap-2 rounded-xl border border-dashed border-zinc-300/70 dark:border-zinc-700/70 hover:border-zinc-400 dark:hover:border-zinc-600 hover:bg-zinc-50 dark:hover:bg-zinc-900/40 transition-colors cursor-pointer"
      >
        <Link size={16} />
        <span>Hubungkan Orang</span>
      </Button>
    {/snippet}
  </Dialog.Trigger>

  <!-- Content -->
  <Dialog.Content
    class="sm:max-w-[560px] p-0 overflow-hidden border-zinc-200/70 dark:border-zinc-800/70"
  >
    <form onsubmit={handleSubmit} class="flex flex-col">
      <!-- Header -->
      <div class="px-6 pt-6 pb-4 bg-white dark:bg-zinc-950">
        <Dialog.Header class="space-y-1">
          <div class="flex items-start gap-3">
            <div
              class="p-2.5 rounded-2xl bg-primary/10 text-primary ring-1 ring-primary/15"
            >
              <Link size={22} />
            </div>
            <div class="min-w-0">
              <Dialog.Title
                class="text-xl sm:text-2xl font-semibold tracking-tight"
              >
                Pembangun Hubungan
              </Dialog.Title>
              <Dialog.Description
                class="text-sm text-zinc-500 dark:text-zinc-400"
              >
                Tentukan koneksi silsilah antar anggota keluarga.
              </Dialog.Description>
            </div>
          </div>
        </Dialog.Header>
      </div>

      <!-- Body -->
      <div class="px-6 py-5 bg-zinc-50/60 dark:bg-zinc-950/30">
        <div class="grid gap-5">
          <!-- People Selection -->
          <div
            class="relative rounded-2xl bg-white dark:bg-zinc-950 border border-zinc-200/70 dark:border-zinc-800/70 p-5"
          >
            <div
              class="grid grid-cols-1 md:grid-cols-[1fr,72px,1fr] items-end gap-4"
            >
              <div class="space-y-2">
                <Label
                  class="text-[11px] font-semibold uppercase tracking-wider text-zinc-500 dark:text-zinc-400"
                >
                  Orang Pertama
                </Label>
                <PersonSelect
                  {persons}
                  bind:value={personAId}
                  placeholder="Pilih Orang"
                  class="w-full rounded-xl border-zinc-200 dark:border-zinc-800 focus-within:ring-2 focus-within:ring-primary/30"
                />
              </div>

              <!-- Connector -->
              <div class="flex md:flex-col items-center justify-center">
                <div
                  class="hidden md:block w-px h-6 bg-zinc-200 dark:bg-zinc-800"
                ></div>
                <div
                  class="mt-2 mb-2 md:mt-0 md:mb-0 h-11 w-11 rounded-2xl bg-primary text-primary-foreground grid place-items-center shadow-sm ring-1 ring-primary/20"
                >
                  <ArrowRight size={18} class="md:rotate-0 rotate-90" />
                </div>
                <div
                  class="hidden md:block w-px h-6 bg-zinc-200 dark:bg-zinc-800"
                ></div>
              </div>

              <div class="space-y-2">
                <Label
                  class="text-[11px] font-semibold uppercase tracking-wider text-zinc-500 dark:text-zinc-400"
                >
                  Orang Kedua
                </Label>
                <PersonSelect
                  {persons}
                  bind:value={personBId}
                  placeholder="Pilih Orang"
                  class="w-full rounded-xl border-zinc-200 dark:border-zinc-800 focus-within:ring-2 focus-within:ring-primary/30"
                />
              </div>
            </div>
          </div>

          <!-- Relationship Type -->
          <div
            class="rounded-2xl bg-white dark:bg-zinc-950 border border-zinc-200/70 dark:border-zinc-800/70 p-5"
          >
            <Label
              class="font-semibold flex items-center gap-2 text-zinc-800 dark:text-zinc-200"
            >
              <Users size={18} class="text-primary" />
              Pilih Jenis Hubungan
            </Label>

            <div class="mt-4 grid grid-cols-1 gap-2.5">
              {#each relationshipOptions as opt}
                <button
                  type="button"
                  onclick={() => (relationshipType = opt.value)}
                  class="group w-full text-left rounded-2xl border transition-colors focus:outline-none focus-visible:ring-2 focus-visible:ring-primary/30
                    {relationshipType === opt.value
                    ? 'border-primary/40 bg-primary/5'
                    : 'border-zinc-200/70 dark:border-zinc-800/70 bg-zinc-50/40 dark:bg-zinc-900/20 hover:bg-white dark:hover:bg-zinc-900/35 hover:border-zinc-300 dark:hover:border-zinc-700'}"
                >
                  <div class="flex items-center justify-between p-4">
                    <div class="flex items-center gap-3">
                      <div
                        class="h-10 w-10 rounded-xl grid place-items-center ring-1 transition-colors
                          {relationshipType === opt.value
                          ? 'bg-primary text-primary-foreground ring-primary/20'
                          : 'bg-white dark:bg-zinc-900 text-zinc-500 dark:text-zinc-400 ring-zinc-200/70 dark:ring-zinc-800/70 group-hover:text-primary'}"
                      >
                        {#if opt.value === "parent_child"}
                          <Users size={18} />
                        {:else if opt.value === "spouse"}
                          <Heart size={18} />
                        {:else}
                          <Link size={18} />
                        {/if}
                      </div>

                      <div class="leading-tight">
                        <div
                          class="font-semibold text-sm
                            {relationshipType === opt.value
                            ? 'text-primary'
                            : 'text-zinc-800 dark:text-zinc-200'}"
                        >
                          {opt.label}
                        </div>
                        <div
                          class="mt-1 text-xs text-zinc-500 dark:text-zinc-400"
                        >
                          {#if opt.value === "parent_child"}
                            Orang pertama adalah Ayah/Ibu
                          {:else if opt.value === "spouse"}
                            Keduanya adalah suami & istri
                          {:else}
                            Keduanya adalah saudara sekandung
                          {/if}
                        </div>
                      </div>
                    </div>

                    <!-- Radio indicator (lebih “radio-like”) -->
                    <div
                      class="h-5 w-5 rounded-full border grid place-items-center transition-colors
                        {relationshipType === opt.value
                        ? 'border-primary bg-primary'
                        : 'border-zinc-300 dark:border-zinc-700 bg-transparent'}"
                    >
                      {#if relationshipType === opt.value}
                        <Check size={12} class="text-white" strokeWidth={3.5} />
                      {/if}
                    </div>
                  </div>
                </button>
              {/each}
            </div>

            <!-- Info -->
            <div
              class="mt-4 rounded-2xl border border-zinc-200/70 dark:border-zinc-800/70 bg-zinc-50 dark:bg-zinc-900/20 p-4"
            >
              <p class="text-xs text-zinc-600 dark:text-zinc-400 flex gap-2">
                <span class="mt-[1px] text-primary">ℹ️</span>
                <span>
                  {#if relationshipType === "parent_child"}
                    <strong class="text-zinc-900 dark:text-zinc-100">
                      {selectedPersonA?.full_name || "..."}
                    </strong>
                    akan tercatat sebagai Orang Tua dari
                    <strong class="text-zinc-900 dark:text-zinc-100">
                      {selectedPersonB?.full_name || "..."}</strong
                    >.
                  {:else if relationshipType === "spouse"}
                    <strong class="text-zinc-900 dark:text-zinc-100">
                      {selectedPersonA?.full_name || "..."}
                    </strong>
                    dan
                    <strong class="text-zinc-900 dark:text-zinc-100">
                      {selectedPersonB?.full_name || "..."}
                    </strong>
                    akan tercatat sebagai pasangan.
                  {:else}
                    <strong class="text-zinc-900 dark:text-zinc-100">
                      {selectedPersonA?.full_name || "..."}
                    </strong>
                    dan
                    <strong class="text-zinc-900 dark:text-zinc-100">
                      {selectedPersonB?.full_name || "..."}
                    </strong>
                    akan tercatat sebagai saudara kandung.
                  {/if}
                </span>
              </p>
            </div>
          </div>
        </div>
      </div>

      <!-- Footer -->
      <div
        class="px-6 py-5 bg-white dark:bg-zinc-950 border-t border-zinc-200/70 dark:border-zinc-800/70"
      >
        <Dialog.Footer class="w-full">
          <Button
            type="submit"
            class="w-full rounded-2xl h-12 font-semibold text-base shadow-sm cursor-pointer"
            disabled={loading}
          >
            {#if loading}
              <Loader2 class="mr-2 h-5 w-5 animate-spin" />
              Menghubungkan...
            {:else}
              Simpan Hubungan
            {/if}
          </Button>
        </Dialog.Footer>
      </div>
    </form>
  </Dialog.Content>
</Dialog.Root>
