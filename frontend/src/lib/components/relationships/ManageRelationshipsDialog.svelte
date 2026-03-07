<script lang="ts">
  import * as Dialog from "$lib/components/ui/dialog/index.js";
  import { Button } from "$lib/components/ui/button/index.js";
  import { Input } from "$lib/components/ui/input/index.js";
  import { Label } from "$lib/components/ui/label/index.js";
  import {
    Link,
    Trash2,
    Heart,
    HeartOff,
    Save,
    Loader2,
    Info,
    XCircle,
    Settings2,
  } from "lucide-svelte";
  import { apiFetch } from "$lib/api";
  import { toast } from "svelte-sonner";
  import { invalidateAll } from "$app/navigation";
  import { Badge } from "$lib/components/ui/badge/index.js";

  let { person, familySlug, relatives = [], readOnly = false } = $props();

  let open = $state(false);
  let loading = $state<string | null>(null);

  async function deleteRelationship(relId: number) {
    if (!confirm("Hapus hubungan ini? Tindakan ini tidak dapat dibatalkan."))
      return;

    loading = `delete-${relId}`;
    try {
      await apiFetch(`/families/${familySlug}/relationships/${relId}`, {
        method: "DELETE",
      });
      toast.success("Hubungan berhasil dihapus");
      await invalidateAll();
    } catch (err: any) {
      toast.error(err.message || "Gagal menghapus hubungan");
    } finally {
      loading = null;
    }
  }

  async function updateRelationship(relId: number, data: any) {
    loading = `update-${relId}`;
    try {
      await apiFetch(`/families/${familySlug}/relationships/${relId}`, {
        method: "PATCH",
        body: JSON.stringify(data),
      });
      toast.success("Hubungan diperbarui");
      await invalidateAll();
    } catch (err: any) {
      toast.error(err.message || "Gagal memperbarui hubungan");
    } finally {
      loading = null;
    }
  }

  function formatDate(iso: string | null | undefined): string {
    if (!iso) return "";
    return iso.split("T")[0];
  }

  function getHumanLabel(rel: any): string {
    const type = rel.relationship_type?.replace("_", "-");
    const other = rel.person;

    if (type === "parent-child") {
      return rel.from === person.id
        ? `Orang Tua dari ${other.nickname || other.full_name}`
        : `Anak dari ${other.nickname || other.full_name}`;
    } else if (type === "spouse") {
      return `Pasangan dengan ${other.nickname || other.full_name}`;
    } else if (type === "sibling") {
      return `Saudara dari ${other.nickname || other.full_name}`;
    }
    return rel.displayLabel;
  }
</script>

<Dialog.Root bind:open>
  <Dialog.Trigger>
    {#snippet child({ props })}
      <Button
        {...props}
        variant="outline"
        class="w-full rounded-[1.8rem] h-14 font-black gap-2 transition-all border-zinc-200 dark:border-zinc-800 hover:bg-zinc-50 dark:hover:bg-zinc-900 cursor-pointer text-base active:scale-95"
      >
        <Settings2 size={18} />
        Kelola Hubungan
      </Button>
    {/snippet}
  </Dialog.Trigger>

  <Dialog.Content
    class="sm:max-w-xl p-0 overflow-hidden rounded-[2.5rem] border-zinc-200/70 dark:border-zinc-800/70 shadow-2xl"
  >
    <div class="flex flex-col h-full max-h-[85vh]">
      <!-- Header -->
      <div class="px-8 pt-8 pb-4 bg-white dark:bg-zinc-950">
        <Dialog.Header>
          <div class="flex items-start gap-4">
            <div
              class="p-3 rounded-2xl bg-primary/10 text-primary ring-1 ring-primary/20 shadow-sm"
            >
              <Settings2 size={24} />
            </div>
            <div>
              <Dialog.Title class="text-2xl font-black tracking-tight"
                >Kelola Relasi</Dialog.Title
              >
              <Dialog.Description class="text-sm font-medium text-zinc-500">
                Atur atau putus hubungan silsilah untuk <strong
                  >{person.full_name}</strong
                >.
              </Dialog.Description>
            </div>
          </div>
        </Dialog.Header>
      </div>

      <!-- Content -->
      <div
        class="px-8 py-6 overflow-y-auto bg-zinc-50/50 dark:bg-zinc-950/30 space-y-6 custom-scrollbar"
      >
        {#if relatives.length === 0}
          <div
            class="flex flex-col items-center justify-center py-12 text-center space-y-4"
          >
            <div
              class="size-16 rounded-3xl bg-zinc-100 dark:bg-zinc-900 grid place-items-center text-zinc-400"
            >
              <HeartOff size={32} />
            </div>
            <div class="space-y-1">
              <p class="font-black text-zinc-800 dark:text-zinc-200">
                Tidak ada hubungan
              </p>
              <p class="text-xs text-zinc-500 max-w-[200px]">
                Orang ini belum memiliki koneksi silsilah yang tercatat.
              </p>
            </div>
          </div>
        {:else}
          <div class="space-y-4">
            {#each relatives as rel}
              <div
                class="p-5 rounded-[2rem] bg-white dark:bg-zinc-950 border border-zinc-200/60 dark:border-zinc-800/60 shadow-sm space-y-4"
              >
                <div class="flex items-start justify-between gap-4">
                  <div class="flex items-start gap-3 min-w-0">
                    <div
                      class="size-12 rounded-2xl overflow-hidden ring-1 ring-zinc-200 shrink-0"
                    >
                      <img
                        src={rel.person.photo_url ||
                          `https://ui-avatars.com/api/?name=${encodeURIComponent(rel.person.full_name)}&background=${rel.person.gender === "male" ? "0ea5e9" : "f43f5e"}&color=fff&bold=true&size=64`}
                        alt={rel.person.full_name}
                        class="w-full h-full object-cover"
                      />
                    </div>
                    <div class="min-w-0">
                      <p
                        class="text-sm font-black text-zinc-800 dark:text-zinc-100 truncate"
                      >
                        {rel.person.full_name}
                      </p>
                      <p
                        class="text-[10px] font-black uppercase tracking-widest text-primary"
                      >
                        {rel.displayLabel}
                      </p>
                    </div>
                  </div>

                  <Button
                    variant="ghost"
                    size="icon"
                    class="size-9 rounded-xl text-rose-500 hover:bg-rose-50 hover:text-rose-600 cursor-pointer"
                    disabled={loading === `delete-${rel.id}`}
                    onclick={() => deleteRelationship(rel.id)}
                  >
                    {#if loading === `delete-${rel.id}`}
                      <Loader2 size={16} class="animate-spin" />
                    {:else}
                      <Trash2 size={16} />
                    {/if}
                  </Button>
                </div>

                <!-- Spouse specific: Divorce settings -->
                {#if rel.relationship_type === "spouse"}
                  <div
                    class="pt-2 mt-4 border-t border-zinc-50 dark:border-zinc-900 space-y-4"
                  >
                    <div class="flex items-center gap-2">
                      <div
                        class="size-7 rounded-lg bg-zinc-100 dark:bg-zinc-900 grid place-items-center text-zinc-400"
                      >
                        <HeartOff size={14} />
                      </div>
                      <span
                        class="text-[10px] font-black uppercase tracking-widest text-zinc-500"
                        >Status Pernikahan</span
                      >
                    </div>

                    <div class="grid gap-3">
                      <div class="flex flex-col gap-2">
                        <Label class="text-[11px] font-bold text-zinc-500 px-1"
                          >Tanggal Perceraian (Jika ada)</Label
                        >
                        <div class="flex gap-2">
                          <Input
                            type="date"
                            class="rounded-xl h-10 text-sm font-medium"
                            value={formatDate(rel.divorce_date)}
                            onchange={(e) => {
                              updateRelationship(rel.id, {
                                divorce_date: e.currentTarget.value,
                              });
                            }}
                          />
                          {#if rel.divorce_date}
                            <Button
                              variant="outline"
                              size="icon"
                              class="size-10 rounded-xl"
                              onclick={() =>
                                updateRelationship(rel.id, {
                                  divorce_date: "",
                                })}
                              title="Hapus status cerai"
                            >
                              <XCircle size={16} class="text-zinc-400" />
                            </Button>
                          {/if}
                        </div>
                        {#if rel.divorce_date}
                          <p
                            class="text-[10px] font-bold text-rose-500 flex items-center gap-1 px-1"
                          >
                            <Info size={10} /> Status: Bercerai
                          </p>
                        {/if}
                      </div>
                    </div>
                  </div>
                {/if}
              </div>
            {/each}
          </div>
        {/if}
      </div>

      <!-- Footer -->
      <div
        class="p-6 bg-white dark:bg-zinc-950 border-t border-zinc-100 dark:border-zinc-800"
      >
        <Button
          class="w-full rounded-[1.8rem] h-12 font-black cursor-pointer"
          onclick={() => (open = false)}
        >
          Selesai
        </Button>
      </div>
    </div>
  </Dialog.Content>
</Dialog.Root>

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
