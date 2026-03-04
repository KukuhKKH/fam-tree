<script lang="ts">
  import { Calendar } from "$lib/components/ui/calendar/index.js";
  import * as Popover from "$lib/components/ui/popover/index.js";
  import { Button } from "$lib/components/ui/button/index.js";
  import { cn } from "$lib/utils.js";
  import { CalendarIcon } from "lucide-svelte";
  import {
    DateFormatter,
    type DateValue,
    getLocalTimeZone,
    parseDate,
    today,
  } from "@internationalized/date";

  let {
    value = $bindable(""),
    placeholder = "Pilih tanggal",
    class: className = "",
  } = $props();

  const df = new DateFormatter("id-ID", {
    dateStyle: "long",
  });

  function safeParse(val: string) {
    if (!val) return undefined;
    try {
      // Extract YYYY-MM-DD from ISO string or full date string
      const datePart = val.split("T")[0];
      return parseDate(datePart);
    } catch (e) {
      return undefined;
    }
  }

  let dateValue = $state<DateValue | undefined>(undefined);
  let initialized = false;

  // Sync external value → internal dateValue
  $effect(() => {
    const parsed = safeParse(value);
    if (parsed && (!dateValue || dateValue.toString() !== parsed.toString())) {
      dateValue = parsed;
    } else if (!value && dateValue) {
      dateValue = undefined;
    }
    // Mark as initialized after first sync
    initialized = true;
  });

  // Sync internal dateValue → external value (only after initialization)
  $effect(() => {
    if (!initialized) return;
    const newVal = dateValue ? dateValue.toString() : "";
    if (newVal !== value) {
      value = newVal;
    }
  });

  const yearRange = Array.from(
    { length: new Date().getFullYear() + 10 - 1890 + 1 },
    (_, i) => 1890 + i,
  );
</script>

<Popover.Root>
  <Popover.Trigger>
    {#snippet child({ props })}
      <Button
        variant="outline"
        class={cn(
          "w-full justify-start text-left font-normal rounded-xl border-zinc-200 dark:border-zinc-800 h-10 shadow-sm",
          !dateValue && "text-zinc-500",
          className,
        )}
        {...props}
      >
        <CalendarIcon class="mr-2 h-4 w-4 opacity-50" />
        {dateValue
          ? df.format(dateValue.toDate(getLocalTimeZone()))
          : placeholder}
      </Button>
    {/snippet}
  </Popover.Trigger>
  <Popover.Content
    class="w-auto p-0 rounded-2xl shadow-2xl border-zinc-200 dark:border-zinc-800 backdrop-blur-xl bg-white/90 dark:bg-zinc-950/90"
  >
    <Calendar
      type="single"
      bind:value={dateValue}
      initialFocus
      captionLayout="dropdown"
      years={yearRange}
      locale="id-ID"
      class="rounded-2xl"
    />
  </Popover.Content>
</Popover.Root>
