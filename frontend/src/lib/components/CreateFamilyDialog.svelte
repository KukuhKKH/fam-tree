<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { apiFetch } from "$lib/api";
  import { Button } from "$lib/components/ui/button/index.js";
  import * as Dialog from "$lib/components/ui/dialog/index.js";
  import { Input } from "$lib/components/ui/input/index.js";
  import { Textarea } from "$lib/components/ui/textarea/index.js";
  import { Label } from "$lib/components/ui/label/index.js";
  import * as Select from "$lib/components/ui/select/index.js";
  import { Plus, Loader2 } from "lucide-svelte";
  import { toast } from "svelte-sonner";

  let open = $state(false);
  let loading = $state(false);
  let error = $state("");

  let name = $state("");
  let description = $state("");
  let visibility = $state("private");

  const visibilityOptions = [
    {
      value: "public",
      label: "Publik",
      hint: "Bisa dicari oleh siapa saja",
    },
    {
      value: "link_only",
      label: "Hanya via Link",
      hint: "Hanya melalui link undangan",
    },
    {
      value: "private",
      label: "Privat",
      hint: "Hanya anggota yang diundang",
    },
  ];

  const currentLabel = $derived(
    visibilityOptions.find((o) => o.value === visibility)?.label ||
      "Pilih Visibilitas",
  );

  async function handleSubmit(e: SubmitEvent) {
    e.preventDefault();
    loading = true;
    error = "";

    try {
      await apiFetch("/families", {
        method: "POST",
        body: JSON.stringify({ name, description, visibility }),
      });

      toast.success(`Keluarga "${name}" berhasil dibuat`);

      // Reset form & close
      name = "";
      description = "";
      visibility = "private";
      open = false;

      // Refresh page data
      await invalidateAll();
    } catch (err: any) {
      error = err.message || "Gagal membuat keluarga.";
      toast.error(error);
    } finally {
      loading = false;
    }
  }
</script>

<Dialog.Root bind:open>
  <Dialog.Trigger>
    {#snippet child({ props })}
      <Button {...props} class="gap-2 shadow-sm cursor-pointer">
        <Plus size={18} />
        <span>Buat Keluarga Baru</span>
      </Button>
    {/snippet}
  </Dialog.Trigger>
  <Dialog.Content class="sm:max-w-[425px]">
    <form onsubmit={handleSubmit}>
      <Dialog.Header>
        <Dialog.Title>Buat Keluarga Baru</Dialog.Title>
        <Dialog.Description>
          Mulai silsilah keluarga baru Anda di sini. Nama dan visibilitas bisa
          diubah nanti.
        </Dialog.Description>
      </Dialog.Header>

      <div class="grid gap-4 py-4">
        {#if error}
          <div
            class="p-3 text-sm bg-red-50 text-red-600 rounded-md border border-red-100 italic"
          >
            {error}
          </div>
        {/if}

        <div class="grid gap-2">
          <Label for="name">Nama Keluarga</Label>
          <Input
            id="name"
            bind:value={name}
            placeholder="Contoh: Keluarga Besar Mbah Ngalidjo"
            required
          />
        </div>

        <div class="grid gap-2">
          <Label for="description">Deskripsi (Opsional)</Label>
          <Textarea
            id="description"
            bind:value={description}
            placeholder="Asal usul atau catatan singkat..."
          />
        </div>

        <div class="grid gap-2">
          <Label for="visibility" class="font-bold">Visibilitas</Label>
          <Select.Root type="single" bind:value={visibility}>
            <Select.Trigger
              class="w-full rounded-xl border-zinc-200 dark:border-zinc-800"
            >
              {currentLabel}
            </Select.Trigger>
            <Select.Content class="backdrop-blur-xl">
              {#each visibilityOptions as opt}
                <Select.Item value={opt.value} label={opt.label}>
                  <div class="flex flex-col text-left">
                    <span class="font-bold text-sm tracking-tight"
                      >{opt.label}</span
                    >
                    <span class="text-[10px] text-zinc-500 leading-tight"
                      >{opt.hint}</span
                    >
                  </div>
                </Select.Item>
              {/each}
            </Select.Content>
          </Select.Root>
        </div>
      </div>

      <Dialog.Footer>
        <Button type="submit" class="w-full cursor-pointer" disabled={loading}>
          {#if loading}
            <Loader2 class="mr-2 h-4 w-4 animate-spin" />
            Memproses...
          {:else}
            Simpan Keluarga
          {/if}
        </Button>
      </Dialog.Footer>
    </form>
  </Dialog.Content>
</Dialog.Root>
