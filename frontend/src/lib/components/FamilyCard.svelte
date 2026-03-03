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
		}
	}>();

	const roleColors: Record<string, string> = {
		owner: 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400',
		admin: 'bg-orange-100 text-orange-700 dark:bg-orange-900/30 dark:text-orange-400',
		editor: 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400',
		viewer: 'bg-zinc-100 text-zinc-700 dark:bg-zinc-800 dark:text-zinc-400',
	};

	const visibilityIcons = {
		public: Globe,
		private: Lock,
		link_only: Link
	};

	const VisibilityIcon = $derived(visibilityIcons[family.visibility as keyof typeof visibilityIcons] || Lock);
</script>

<a href="/families/{family.slug}" class="block group h-full">
	<Card.Root class="h-full transition-all duration-300 hover:shadow-md hover:border-primary/30 group-hover:-translate-y-1">
		<Card.Header>
			<div class="flex items-start justify-between">
				<div class="p-2 bg-primary/5 rounded-lg text-primary">
					<Users size={20} />
				</div>
				<Badge variant="outline" class="capitalize {roleColors[family.my_role] || ''}">
					{family.my_role}
				</Badge>
			</div>
			<Card.Title class="mt-4 group-hover:text-primary transition-colors line-clamp-1">{family.name}</Card.Title>
			<Card.Description class="line-clamp-2 h-10">
				{family.description || 'Tidak ada deskripsi.'}
			</Card.Description>
		</Card.Header>
		<Card.Footer class="border-t border-zinc-100 dark:border-zinc-800 pt-4 flex items-center justify-between text-xs text-zinc-500">
			<div class="flex items-center gap-1.5 capitalize">
				<VisibilityIcon size={14} />
				{family.visibility.replace('_', ' ')}
			</div>
			<div class="font-mono opacity-0 group-hover:opacity-100 transition-opacity">
				/{family.slug}
			</div>
		</Card.Footer>
	</Card.Root>
</a>
