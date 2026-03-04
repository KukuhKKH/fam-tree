<script lang="ts">
  import * as AlertDialog from "$lib/components/ui/alert-dialog/index.js";
  import { Button } from "$lib/components/ui/button/index.js";
  import { Loader2, AlertTriangle } from "lucide-svelte";

  let {
    open = $bindable(false),
    onConfirm,
    title = "Apakah Anda yakin?",
    description = "Tindakan ini tidak dapat dibatalkan.",
    confirmText = "Hapus",
    cancelText = "Batal",
    loading = false,
    variant = "destructive",
  } = $props();

  function handleConfirm() {
    onConfirm();
  }
</script>

<AlertDialog.Root bind:open>
  <AlertDialog.Content
    class="rounded-3xl border-zinc-200 dark:border-zinc-800 shadow-2xl"
  >
    <AlertDialog.Header>
      <div
        class="mx-auto w-12 h-12 rounded-2xl bg-rose-50 dark:bg-rose-950/30 flex items-center justify-center text-rose-500 mb-4"
      >
        <AlertTriangle size={24} />
      </div>
      <AlertDialog.Title
        class="text-2xl font-black text-center text-zinc-900 dark:text-white"
      >
        {title}
      </AlertDialog.Title>
      <AlertDialog.Description
        class="text-center text-zinc-500 dark:text-zinc-400 text-sm py-2"
      >
        {description}
      </AlertDialog.Description>
    </AlertDialog.Header>
    <AlertDialog.Footer class="flex flex-col sm:flex-row gap-3 mt-6">
      <AlertDialog.Cancel
        class="flex-1 rounded-xl h-12 font-bold bg-zinc-100 dark:bg-zinc-800 text-zinc-600 border-none hover:bg-zinc-200 dark:hover:bg-zinc-700"
      >
        {cancelText}
      </AlertDialog.Cancel>
      <AlertDialog.Action
        onclick={handleConfirm}
        class="flex-1 rounded-xl h-12 font-bold bg-rose-500 text-white hover:bg-rose-600 shadow-lg shadow-rose-500/20"
      >
        {#if loading}
          <Loader2 class="mr-2 h-4 w-4 animate-spin" />
        {/if}
        {confirmText}
      </AlertDialog.Action>
    </AlertDialog.Footer>
  </AlertDialog.Content>
</AlertDialog.Root>
