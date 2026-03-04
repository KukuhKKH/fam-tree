<script lang="ts">
  import * as Command from "$lib/components/ui/command/index.js";
  import * as Popover from "$lib/components/ui/popover/index.js";
  import { Button } from "$lib/components/ui/button/index.js";
  import { Check, ChevronsUpDown } from "lucide-svelte";
  import { cn } from "$lib/utils.js";
  import { tick } from "svelte";

  let {
    persons = [],
    value = $bindable(""),
    placeholder = "Pilih orang...",
    class: className = "",
  } = $props();

  let open = $state(false);

  const selectedPerson = $derived(
    persons.find((p: any) => String(p.id) === value),
  );

  const label = $derived(selectedPerson?.full_name || placeholder);

  function closeAndSetValue(val: string) {
    value = val;
    open = false;
  }
</script>

<Popover.Root bind:open>
  <Popover.Trigger>
    {#snippet child({ props })}
      <Button
        variant="outline"
        role="combobox"
        aria-expanded={open}
        class={cn(
          "w-full flex justify-between rounded-xl h-12 border-zinc-200 dark:border-zinc-700",
          className,
        )}
        style={`${props.style ?? ""}; width: 100%;`}
        {...props}
      >
        <span class="truncate">{label}</span>
        <ChevronsUpDown class="ml-2 h-4 w-4 shrink-0 opacity-50" />
      </Button>
    {/snippet}
  </Popover.Trigger>
  <Popover.Content
    style="width: var(--radix-popper-anchor-width, var(--radix-popover-trigger-width, var(--bits-popover-anchor-width)));"
    class="p-0 rounded-2xl shadow-2xl border-zinc-200 dark:border-zinc-700 overflow-hidden"
  >
    <Command.Root class="rounded-2xl bg-popover">
      <Command.Input
        placeholder="Cari orang..."
        class="h-12 border-none focus:ring-0"
      />
      <hr class="border-zinc-100 dark:border-zinc-800" />
      <Command.Empty class="py-6 text-center text-sm text-zinc-500">
        Orang tidak ditemukan.
      </Command.Empty>
      <Command.List class="max-h-[300px] overflow-y-auto p-1.5">
        <Command.Group>
          {#each persons as person}
            <Command.Item
              value={person.full_name}
              onSelect={() => closeAndSetValue(String(person.id))}
              class="flex items-center justify-between py-3 px-4 rounded-xl cursor-pointer data-[highlighted]:bg-zinc-100 dark:data-[highlighted]:bg-zinc-800 transition-colors"
            >
              <div class="flex flex-col">
                <span class="font-bold text-sm text-zinc-900 dark:text-zinc-100"
                  >{person.full_name}</span
                >
                {#if person.nickname}
                  <span class="text-[10px] text-zinc-500 italic"
                    >"{person.nickname}"</span
                  >
                {/if}
              </div>
              <Check
                class={cn(
                  "ml-auto h-4 w-4 text-primary transition-opacity",
                  value === String(person.id) ? "opacity-100" : "opacity-0",
                )}
              />
            </Command.Item>
          {/each}
        </Command.Group>
      </Command.List>
    </Command.Root>
  </Popover.Content>
</Popover.Root>
