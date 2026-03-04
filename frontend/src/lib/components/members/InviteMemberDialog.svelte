<script lang="ts">
  import * as Dialog from "$lib/components/ui/dialog/index.js";
  import { Button } from "$lib/components/ui/button/index.js";
  import { Input } from "$lib/components/ui/input/index.js";
  import { Label } from "$lib/components/ui/label/index.js";
  import * as Select from "$lib/components/ui/select/index.js";
  import { UserPlus, Loader2, Mail, Shield } from "lucide-svelte";
  import { apiFetch } from "$lib/api";
  import { toast } from "svelte-sonner";
  import { invalidateAll } from "$app/navigation";

  let { familySlug } = $props();

  let open = $state(false);
  let loading = $state(false);
  let email = $state("");
  let role = $state("viewer");

  const roleOptions = [
    { value: "admin", label: "Admin" },
    { value: "editor", label: "Editor" },
    { value: "viewer", label: "Viewer" },
  ];

  async function inviteMember() {
    if (!email) return toast.error("Email wajib diisi");

    loading = true;
    try {
      await apiFetch(`/families/${familySlug}/members`, {
        method: "POST",
        body: JSON.stringify({ email, role }),
      });

      toast.success(`Undangan berhasil dikirim ke ${email}`);
      open = false;
      email = "";
      role = "viewer";
      await invalidateAll();
    } catch (err: any) {
      toast.error(err.message || "Gagal mengirim undangan");
    } finally {
      loading = false;
    }
  }
</script>

<Dialog.Root bind:open>
  <Dialog.Trigger>
    {#snippet child({ props })}
      <Button
        {...props}
        class="gap-2 shadow-lg shadow-primary/10 rounded-xl cursor-pointer"
      >
        <UserPlus size={18} />
        Undang Anggota
      </Button>
    {/snippet}
  </Dialog.Trigger>
  <Dialog.Content class="sm:max-w-[425px]">
    <Dialog.Header>
      <Dialog.Title class="text-2xl font-bold flex items-center gap-2">
        <UserPlus class="text-primary" />
        Undang Anggota
      </Dialog.Title>
      <Dialog.Description>
        Masukkan alamat email orang yang ingin Anda undang untuk bergabung dalam
        silsilah keluarga ini.
      </Dialog.Description>
    </Dialog.Header>

    <div class="grid gap-6 py-6">
      <div class="grid gap-2">
        <Label for="email" class="font-bold flex items-center gap-2">
          <Mail size={14} class="text-zinc-400" />
          Alamat Email
        </Label>
        <Input
          id="email"
          type="email"
          placeholder="nama@email.com"
          bind:value={email}
          class="rounded-xl border-zinc-200 dark:border-zinc-800"
        />
      </div>
      <div class="grid gap-2">
        <Label for="role" class="font-bold flex items-center gap-2">
          <Shield size={14} class="text-zinc-400" />
          Peran (Role)
        </Label>
        <Select.Root type="single" bind:value={role}>
          <Select.Trigger
            class="rounded-xl border-zinc-200 dark:border-zinc-800 capitalize"
          >
            <Select.Value placeholder="Pilih Role" />
          </Select.Trigger>
          <Select.Content class="backdrop-blur-xl">
            {#each roleOptions as opt}
              <Select.Item
                value={opt.value}
                label={opt.label}
                class="capitalize cursor-pointer"
              >
                {opt.label}
              </Select.Item>
            {/each}
          </Select.Content>
        </Select.Root>
        <p class="text-[11px] text-zinc-500 italic px-1">
          * Admin bisa mengelola anggota & profil. Viewer hanya bisa melihat.
        </p>
      </div>
    </div>

    <Dialog.Footer>
      <Button
        type="submit"
        onclick={inviteMember}
        disabled={loading}
        class="w-full sm:w-auto rounded-xl px-8"
      >
        {#if loading}
          <Loader2 class="mr-2 h-4 w-4 animate-spin" />
          Mengirim...
        {:else}
          Kirim Undangan
        {/if}
      </Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
