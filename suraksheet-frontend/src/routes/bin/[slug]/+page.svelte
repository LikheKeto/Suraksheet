<script lang="ts">
	import CreateDocument from '$lib/components/CreateDocument.svelte';
	import Documents from '$lib/components/Documents.svelte';
	import { binsStore } from '$lib/store';
	import type { Bin } from '$lib/types';
	import { Breadcrumb, BreadcrumbItem, Skeleton } from 'flowbite-svelte';
	import { CirclePlusSolid } from 'flowbite-svelte-icons';

	export let data: {
		error?: string;
		binID?: number;
	};
	let bins: Bin[] = [];
	binsStore.subscribe((vals) => (bins = vals));
	let bin = bins.find((b) => b.id === data.binID);

	let hideBinCreator = true;
</script>

{#if bin && bin.id}
	<Breadcrumb class="mb-2" solid>
		<BreadcrumbItem href="/" home>Home</BreadcrumbItem>
		<BreadcrumbItem>{bin.name}</BreadcrumbItem>
	</Breadcrumb>
	<div class="flex items-end gap-2">
		<p class="text-2xl font-bold">{bin?.name}</p>
		<button
			on:click={() => (hideBinCreator = false)}
			class="text-primary-400 hover:text-primary-500 dark:text-primary-600 dark:hover:text-primary-500 rounded-lg"
		>
			<CirclePlusSolid size="xl" />
		</button>
	</div>
	<CreateDocument bind:binID={bin.id} bind:hidden={hideBinCreator} />
{:else}
	<Skeleton size="sm" class="h-4 overflow-clip" />
{/if}
<Documents bin={bin?.name} />
