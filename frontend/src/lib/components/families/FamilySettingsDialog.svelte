<script lang="ts">
  import * as Dialog from "$lib/components/ui/dialog/index.js";
  import { Button } from "$lib/components/ui/button/index.js";
  import { Input } from "$lib/components/ui/input/index.js";
  import { Label } from "$lib/components/ui/label/index.js";
  import { Textarea } from "$lib/components/ui/textarea/index.js";
  import * as Select from "$lib/components/ui/select/index.js";
  import { Settings, Loader2, Trash2, Save, Globe, Lock } from "lucide-svelte";
  import { apiFetch } from "$lib/api";
  import { toast } from "svelte-sonner";
  import { invalidateAll, goto } from "$app/navigation";

  let { family } = $props();

  let open = $state(false);
  let loading = $state(false);
  let deleting = $state(false);

  let form = $state({
    name: "",
    description: "",
    visibility: "private",
  });

  // Initialize/Sync form state when prop changes or modal opens
  $effect(() => {
    if (open && family) {
      form.name = family.name || "";
      form.description = family.description || "";
      form.visibility = family.visibility || "private";
    }
  });

  const visibilityOptions = [
    { value: "public", label: "Publik (Terlihat semua orang)", icon: Globe },
    { value: "private", label: "Privat (Hanya anggota)", icon: Lock },
  ];

  async function updateFamily() {
    if (!form.name) return toast.error("Nama keluarga wajib diisi");

    loading = true;
    try {
      await apiFetch(`/families/${family.slug}`, {
        method: "PATCH",
        body: JSON.stringify(form),
      });

      toast.success("Info keluarga berhasil diperbarui");
      open = false;
      await invalidateAll();
    } catch (err: any) {
      toast.error(err.message || "Gagal memperbarui info keluarga");
    } finally {
      loading = false;
    }
  }

  async function deleteFamily() {
    if (
      !confirm(
        `TINDAKAN BERBAHAYA!\n\nApakah Anda yakin ingin menghapus keluarga "${family.name}"? Semua data silsilah di dalamnya akan hilang selamanya.`,
      )
    )
      return;

    deleting = true;
    try {
      await apiFetch(`/families/${family.slug}`, {
        method: "DELETE",
      });

      toast.success("Keluarga berhasil dihapus");
      open = false;
      await goto("/dashboard");
    } catch (err: any) {
      toast.error(err.message || "Gagal menghapus keluarga");
    } finally {
      deleting = false;
    }
  }
</script>

<Dialog.Root bind:open>
  <Dialog.Trigger>
    {#snippet child({ props })}
      <Button
        {...props}
        variant="outline"
        class="gap-2 rounded-xl border-zinc-200 dark:border-zinc-800 cursor-pointer"
      >
        <Settings size={18} />
        Pengaturan
      </Button>
    {/snippet}
  </Dialog.Trigger>
  <Dialog.Content class="sm:max-w-[500px]">
    <Dialog.Header>
      <Dialog.Title class="text-2xl font-bold flex items-center gap-2">
        <Settings class="text-primary" />
        Pengaturan Keluarga
      </Dialog.Title>
      <Dialog.Description>
        Sesuaikan informasi utama keluarga dan pengaturan privasi di sini.
      </Dialog.Description>
    </Dialog.Header>

    <div class="grid gap-6 py-6">
      <div class="grid gap-2">
        <Label for="name" class="font-bold">Nama Keluarga</Label>
        <Input
          id="name"
          bind:value={form.name}
          placeholder="Contoh: Bani Adam"
          class="rounded-xl border-zinc-200 dark:border-zinc-800"
        />
      </div>

      <div class="grid gap-2">
        <Label for="description" class="font-bold">Deskripsi</Label>
        <Textarea
          id="description"
          bind:value={form.description}
          placeholder="Ceritakan sedikit tentang asal-usul keluarga ini..."
          class="rounded-xl border-zinc-200 dark:border-zinc-800 min-h-[100px]"
        />
      </div>

      <div class="grid gap-2">
        <Label for="visibility" class="font-bold">Visibilitas</Label>
        <Select.Root
          type="single"
          value={form.visibility}
          onValueChange={(v) => (form.visibility = v)}
        >
          <Select.Trigger
            class="rounded-xl border-zinc-200 dark:border-zinc-800 capitalize"
          >
            <Select.Value placeholder="Pilih Visibilitas" />
          </Select.Trigger>
          <Select.Content class="backdrop-blur-xl">
            {#each visibilityOptions as opt}
              <Select.Item
                value={opt.value}
                label={opt.label}
                class="cursor-pointer gap-2"
              >
                <opt.icon size={14} class="inline mr-2" />
                {opt.label}
              </Select.Item>
            {/each}
          </Select.Content>
        </Select.Root>
      </div>

      <div class="pt-4 border-t border-zinc-100 dark:border-zinc-800">
        <h4
          class="text-sm font-bold text-destructive mb-2 uppercase tracking-wider"
        >
          Zona Bahaya
        </h4>
        <p class="text-xs text-zinc-500 mb-4">
          Menghapus keluarga bersifat permanen dan tidak dapat dibatalkan.
        </p>
        <Button
          variant="destructive"
          class="w-full sm:w-auto rounded-xl gap-2 cursor-pointer"
          onclick={deleteFamily}
          disabled={deleting}
        >
          {#if deleting}
            <Loader2 size={16} class="animate-spin" />
            Menghapus...
          {:else}
            <Trash2 size={16} />
            Hapus Keluarga Ini
          {/if}
        </Button>
      </div>
    </div>

    <Dialog.Footer>
      <Button
        type="submit"
        onclick={updateFamily}
        disabled={loading}
        class="w-full sm:w-auto rounded-xl px-8 bg-primary hover:bg-primary/90"
      >
        {#if loading}
          <Loader2 class="mr-2 h-4 w-4 animate-spin" />
          Menyimpan...
        {:else}
          <Save size={18} class="mr-2" />
          Simpan Perubahan
        {/if}
      </Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
