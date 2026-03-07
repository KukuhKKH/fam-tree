<script lang="ts">
  import { cn } from "$lib/utils.js";
  import { Popover as PopoverPrimitive } from "bits-ui";

  let {
    ref = $bindable(null),
    class: className,
    ...restProps
  }: PopoverPrimitive.TriggerProps = $props();
</script>

<PopoverPrimitive.Trigger
  bind:ref
  data-slot="popover-trigger"
  class={cn("", className)}
  {...restProps}
>
  {#snippet child({ props })}
    {#if restProps.child}
      {@render restProps.child({ props })}
    {:else}
      <button {...props}>
        {@render restProps.children?.()}
      </button>
    {/if}
  {/snippet}
</PopoverPrimitive.Trigger>
