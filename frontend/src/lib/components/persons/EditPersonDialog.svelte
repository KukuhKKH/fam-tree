<script lang="ts">
  import * as Dialog from "$lib/components/ui/dialog/index.js";
  import { Button } from "$lib/components/ui/button/index.js";
  import { Input } from "$lib/components/ui/input/index.js";
  import { Label } from "$lib/components/ui/label/index.js";
  import { Textarea } from "$lib/components/ui/textarea/index.js";
  import * as Select from "$lib/components/ui/select/index.js";
  import {
    Edit,
    Loader2,
    Calendar,
    MapPin,
    Heart,
    HeartOff,
    Save,
  } from "lucide-svelte";
  import DatePicker from "$lib/components/ui/DatePicker.svelte";
  import { apiFetch } from "$lib/api";
  import { toast } from "svelte-sonner";
  import { invalidateAll } from "$app/navigation";

  let { familySlug, person } = $props();

  let open = $state(false);
  let loading = $state(false);

  // Initialize form with default empty values, will be populated by $effect
  let form = $state({
    full_name: "",
    nickname: "",
    gender: "male",
    birth_date: "",
    birth_place: "",
    is_alive: true,
    death_date: "",
    death_place: "",
    bio: "",
  });

  const genderOptions = [
    { value: "male", label: "Laki-laki" },
    { value: "female", label: "Perempuan" },
  ];

  const currentGenderLabel = $derived(
    genderOptions.find((o) => o.value === form.gender)?.label ||
      "Pilih Jenis Kelamin",
  );

  async function handleSubmit(e: SubmitEvent) {
    e.preventDefault();
    if (form.full_name.length < 2) {
      return toast.error("Nama lengkap minimal 2 karakter");
    }

    loading = true;
    try {
      await apiFetch(`/families/${familySlug}/persons/${person.id}`, {
        method: "PATCH",
        body: JSON.stringify(form),
      });

      toast.success(`Data ${form.full_name} berhasil diperbarui`);
      open = false;
      await invalidateAll();
    } catch (err: any) {
      toast.error(err.message || "Gagal memperbarui data person");
    } finally {
      loading = false;
    }
  }

  // Sync form data when dialog opens
  $effect(() => {
    if (open) {
      form.full_name = person.full_name || "";
      form.nickname = person.nickname || "";
      form.gender = person.gender || "male";
      form.birth_date = person.birth_date || "";
      form.birth_place = person.birth_place || "";
      form.is_alive = person.is_alive ?? true;
      form.death_date = person.death_date || "";
      form.death_place = person.death_place || "";
      form.bio = person.bio || "";
    }
  });
</script>

<Dialog.Root bind:open>
  <Dialog.Trigger>
    {#snippet child({ props })}
      <Button
        {...props}
        class="rounded-2xl gap-2 h-11 px-6 font-bold shadow-lg shadow-primary/20 bg-primary hover:bg-primary/90 transition-all hover:scale-[1.02] cursor-pointer"
      >
        <Edit size={16} />
        Edit Profil
      </Button>
    {/snippet}
  </Dialog.Trigger>
  <Dialog.Content
    class="sm:max-w-[640px] p-0 overflow-hidden rounded-[2.5rem] border-zinc-200/70 dark:border-zinc-800/70 shadow-2xl"
  >
    <form onsubmit={handleSubmit} class="flex flex-col max-h-[90vh]">
      <!-- Header Area -->
      <div class="px-8 pt-8 pb-4 bg-white dark:bg-zinc-950">
        <Dialog.Header class="space-y-1">
          <div class="flex items-start gap-4">
            <div
              class="p-3 rounded-2xl bg-primary/10 text-primary ring-1 ring-primary/20 shadow-sm"
            >
              <Edit size={24} />
            </div>
            <div class="min-w-0">
              <Dialog.Title class="text-2xl font-black tracking-tight"
                >Edit Data Person</Dialog.Title
              >
              <Dialog.Description class="text-sm font-medium text-zinc-500">
                Lengkapi atau perbarui informasi silsilah demi arsip keluarga
                yang akurat.
              </Dialog.Description>
            </div>
          </div>
        </Dialog.Header>
      </div>

      <!-- Scrollable Body -->
      <div class="px-8 py-6 overflow-y-auto bg-zinc-50/50 dark:bg-zinc-950/30">
        <div class="grid gap-8">
          <!-- Person Identity Card -->
          <div
            class="grid gap-6 p-6 rounded-[2rem] bg-white dark:bg-zinc-950 border border-zinc-200/60 dark:border-zinc-800/60 shadow-sm"
          >
            <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
              <div class="space-y-2">
                <Label
                  for="full_name"
                  class="text-xs font-black uppercase tracking-widest text-zinc-400 px-1"
                  >Nama Lengkap</Label
                >
                <Input
                  id="full_name"
                  bind:value={form.full_name}
                  placeholder="Nama Lengkap"
                  class="rounded-xl border-zinc-200 dark:border-zinc-800 h-11 focus-visible:ring-primary/20"
                  required
                />
              </div>
              <div class="space-y-2">
                <Label
                  for="nickname"
                  class="text-xs font-black uppercase tracking-widest text-zinc-400 px-1"
                  >Panggilan</Label
                >
                <Input
                  id="nickname"
                  bind:value={form.nickname}
                  placeholder="Nama Panggilan"
                  class="rounded-xl border-zinc-200 dark:border-zinc-800 h-11 focus-visible:ring-primary/20"
                />
              </div>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
              <div class="space-y-2">
                <Label
                  class="text-xs font-black uppercase tracking-widest text-zinc-400 px-1"
                  >Jenis Kelamin</Label
                >
                <Select.Root
                  type="single"
                  value={form.gender}
                  onValueChange={(v) => (form.gender = v)}
                >
                  <Select.Trigger
                    class="rounded-xl h-11 border-zinc-200 dark:border-zinc-800 shadow-none"
                  >
                    {currentGenderLabel}
                  </Select.Trigger>
                  <Select.Content
                    class="rounded-2xl border-zinc-200 dark:border-zinc-800 shadow-2xl p-1"
                  >
                    {#each genderOptions as opt}
                      <Select.Item
                        value={opt.value}
                        label={opt.label}
                        class="rounded-xl px-4 py-2 my-0.5 cursor-pointer"
                      >
                        {opt.label}
                      </Select.Item>
                    {/each}
                  </Select.Content>
                </Select.Root>
              </div>

              <div class="space-y-2">
                <Label
                  class="text-xs font-black uppercase tracking-widest text-zinc-400 px-1"
                  >Status Kehidupan</Label
                >
                <div
                  class="flex items-center gap-2 h-11 bg-zinc-50 dark:bg-zinc-900 rounded-xl px-3 border border-zinc-100 dark:border-zinc-800"
                >
                  <button
                    type="button"
                    onclick={() => (form.is_alive = true)}
                    class="flex-1 flex items-center justify-center gap-2 py-1.5 rounded-lg transition-all text-xs font-bold {form.is_alive
                      ? 'bg-white dark:bg-zinc-800 shadow-sm text-primary ring-1 ring-zinc-200/50 dark:ring-zinc-700/50'
                      : 'text-zinc-400 hover:text-zinc-600'}"
                  >
                    <Heart
                      size={14}
                      class={form.is_alive ? "text-rose-500" : ""}
                    /> Hidup
                  </button>
                  <button
                    type="button"
                    onclick={() => (form.is_alive = false)}
                    class="flex-1 flex items-center justify-center gap-2 py-1.5 rounded-lg transition-all text-xs font-bold {!form.is_alive
                      ? 'bg-white dark:bg-zinc-800 shadow-sm text-primary ring-1 ring-zinc-200/50 dark:ring-zinc-700/50'
                      : 'text-zinc-400 hover:text-zinc-600'}"
                  >
                    <HeartOff
                      size={14}
                      class={!form.is_alive ? "text-zinc-500" : ""}
                    /> Wafat
                  </button>
                </div>
              </div>
            </div>
          </div>

          <!-- Birth Information -->
          <div
            class="p-6 rounded-[2rem] bg-white dark:bg-zinc-950 border border-zinc-200/60 dark:border-zinc-800/60 shadow-sm space-y-5"
          >
            <h4
              class="text-[11px] font-black uppercase tracking-[0.2em] text-primary flex items-center gap-2"
            >
              <Calendar size={14} /> Informasi Kelahiran
            </h4>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
              <div class="space-y-2">
                <Label
                  for="birth_date"
                  class="text-[10px] font-bold text-zinc-400 px-1"
                  >Tanggal Lahir</Label
                >
                <DatePicker
                  bind:value={form.birth_date}
                  placeholder="Pilih Tanggal"
                />
              </div>
              <div class="space-y-2">
                <Label
                  for="birth_place"
                  class="text-[10px] font-bold text-zinc-400 px-1"
                  >Tempat Lahir</Label
                >
                <Input
                  id="birth_place"
                  bind:value={form.birth_place}
                  placeholder="Kota/Tempat"
                  class="rounded-xl border-zinc-200 dark:border-zinc-800 h-11"
                />
              </div>
            </div>
          </div>

          <!-- Death Information (Conditional) -->
          {#if !form.is_alive}
            <div
              class="p-6 rounded-[2rem] bg-rose-50/20 dark:bg-rose-950/5 border border-rose-100/50 dark:border-rose-900/20 shadow-sm space-y-5"
            >
              <h4
                class="text-[11px] font-black uppercase tracking-[0.2em] text-rose-500 flex items-center gap-2"
              >
                <HeartOff size={14} /> Informasi Kematian
              </h4>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
                <div class="space-y-2">
                  <Label
                    for="death_date"
                    class="text-[10px] font-bold text-rose-400 px-1"
                    >Tanggal Wafat</Label
                  >
                  <DatePicker
                    bind:value={form.death_date}
                    placeholder="Pilih Tanggal"
                  />
                </div>
                <div class="space-y-2">
                  <Label
                    for="death_place"
                    class="text-[10px] font-bold text-rose-400 px-1"
                    >Tempat Wafat</Label
                  >
                  <Input
                    id="death_place"
                    bind:value={form.death_place}
                    placeholder="Kota/Tempat"
                    class="rounded-xl border-rose-100/40 dark:border-rose-900/30 h-11 bg-white/50 dark:bg-zinc-950/50"
                  />
                </div>
              </div>
            </div>
          {/if}

          <!-- Biography Section -->
          <div class="space-y-3">
            <Label
              for="bio"
              class="text-xs font-black uppercase tracking-widest text-zinc-400 px-1"
              >Biografi & Catatan</Label
            >
            <Textarea
              id="bio"
              bind:value={form.bio}
              placeholder="Ceritakan sejarah singkat atau kenangan tentang beliau..."
              class="min-h-[140px] rounded-[1.5rem] border-zinc-200 dark:border-zinc-800 p-5 focus-visible:ring-primary/20 bg-white dark:bg-zinc-950 shadow-sm text-sm italic"
            />
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
              Simpan Perubahan
            {/if}
          </Button>
        </Dialog.Footer>
      </div>
    </form>
  </Dialog.Content>
</Dialog.Root>
