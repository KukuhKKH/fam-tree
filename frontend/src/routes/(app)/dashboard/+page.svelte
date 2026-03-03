<script lang="ts">
	import { page } from '$app/state';
	import * as Card from "$lib/components/ui/card/index.js";
	import { Users, TreeDeciduous, Heart } from "lucide-svelte";
	import FamilyCard from '$lib/components/FamilyCard.svelte';
	import CreateFamilyDialog from '$lib/components/CreateFamilyDialog.svelte';

	let families = $derived(page.data.families || []);
	
	const stats = $derived([
		{ label: 'Keluarga', value: families.length.toString(), icon: Users, color: 'text-blue-500' },
		{ label: 'Anggota Silsilah', value: '-', icon: TreeDeciduous, color: 'text-green-500' },
		{ label: 'Relasi', value: '-', icon: Heart, color: 'text-red-500' },
	]);
</script>

<div class="space-y-8">
	<div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
		<div>
			<h1 class="text-3xl font-bold tracking-tight">Selamat datang, {page.data.user?.name || 'User'}!</h1>
			<p class="text-zinc-500 dark:text-zinc-400">Ikhtisar silsilah keluarga Anda hari ini.</p>
		</div>
		<CreateFamilyDialog />
	</div>

	<div class="grid gap-4 md:grid-cols-3">
		{#each stats as stat}
			<Card.Root>
				<Card.Header class="flex flex-row items-center justify-between space-y-0 pb-2">
					<Card.Title class="text-sm font-medium">{stat.label}</Card.Title>
					<stat.icon class="h-4 w-4 {stat.color}" />
				</Card.Header>
				<Card.Content>
					<div class="text-2xl font-bold">{stat.value}</div>
				</Card.Content>
			</Card.Root>
		{/each}
	</div>

	<div class="space-y-4">
		<h2 class="text-xl font-semibold tracking-tight">Daftar Keluarga Saya</h2>
		
		{#if families.length === 0}
			<Card.Root class="bg-primary/5 border-primary/10 border-dashed border-2">
				<Card.Content class="p-12 flex flex-col items-center text-center space-y-4">
					<div class="p-4 bg-primary/10 rounded-full text-primary">
						<Users size={48} />
					</div>
					<div class="space-y-2">
						<h3 class="text-lg font-semibold">Belum ada keluarga</h3>
						<p class="text-sm text-zinc-600 dark:text-zinc-400 max-w-xs">
							Anda belum bergabung atau membuat keluarga. Mulai dengan membuat keluarga baru!
						</p>
					</div>
					<CreateFamilyDialog />
				</Card.Content>
			</Card.Root>
		{:else}
			<div class="grid gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
				{#each families as family}
					<FamilyCard {family} />
				{/each}
			</div>
		{/if}
	</div>
</div>
