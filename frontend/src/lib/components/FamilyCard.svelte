<script lang="ts">
  import * as Card from "$lib/components/ui/card/index.js";
  import { Badge } from "$lib/components/ui/badge/index.js";
  import { Users, Globe, Lock, Link } from "lucide-svelte";

  let { family } = $props<{
    family: {
      name: string;
      slug: string;
      description: string;
      visibility: string;
      my_role: string;
    };
  }>();

  const roleColors: Record<string, string> = {
    owner: "bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400",
    admin:
      "bg-orange-100 text-orange-700 dark:bg-orange-900/30 dark:text-orange-400",
    editor: "bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400",
    viewer: "bg-zinc-100 text-zinc-700 dark:bg-zinc-800 dark:text-zinc-400",
  };

  const visibilityIcons = {
    public: Globe,
    private: Lock,
    link_only: Link,
  };

  const VisibilityIcon = $derived(
    visibilityIcons[family.visibility as keyof typeof visibilityIcons] || Lock,
  );
</script>

<a href="/families/{family.slug}" class="block group h-full">
  <Card.Root
    class="h-full transition-all duration-500 hover:shadow-xl hover:shadow-primary/5 hover:border-primary/20 group-hover:-translate-y-1.5 overflow-hidden border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900"
  >
    <Card.Header>
      <div class="flex items-start justify-between">
        <div
          class="p-2.5 bg-primary/10 rounded-xl text-primary border border-primary/10"
        >
          <Users size={22} />
        </div>
        <Badge
          variant="secondary"
          class="capitalize font-bold tracking-tight rounded-lg px-2.5 py-0.5 {roleColors[
            family.my_role
          ] || ''}"
        >
          {family.my_role}
        </Badge>
      </div>
      <Card.Title
        class="mt-5 text-xl font-bold group-hover:text-primary transition-colors line-clamp-1 leading-tight"
        >{family.name}</Card.Title
      >
      <Card.Description
        class="mt-2 line-clamp-2 h-10 text-zinc-500 dark:text-zinc-400 text-sm leading-relaxed"
      >
        {family.description ||
          "Silsilah keluarga yang belum memiliki deskripsi."}
      </Card.Description>
    </Card.Header>
    <Card.Footer
      class="border-t border-zinc-100 dark:border-zinc-800 p-5 bg-zinc-50/50 dark:bg-zinc-950/20 flex items-center justify-between text-[11px] font-bold uppercase tracking-wider text-zinc-400 dark:text-zinc-500"
    >
      <div class="flex items-center gap-2">
        <VisibilityIcon size={14} class="text-primary/60" />
        <span
          >{family.visibility === "link_only"
            ? "Hanya Link"
            : family.visibility}</span
        >
      </div>
      <div
        class="font-mono lowercase opacity-40 group-hover:opacity-100 transition-opacity"
      >
        /{family.slug}
      </div>
    </Card.Footer>
  </Card.Root>
</a>
