<script lang="ts">
	import CreateDocument from '$lib/components/documents/CreateDocument.svelte';
	import EditBin from '$lib/components/bins/CreateEditBin.svelte';
	import Documents from '$lib/components/documents/Documents.svelte';
	import { binsStore } from '$lib/store';
	import type { Bin } from '$lib/types';
	import {
		Button,
		Dropdown,
		DropdownItem,
		Breadcrumb,
		BreadcrumbItem,
		Skeleton
	} from 'flowbite-svelte';
	import {
		CirclePlusSolid,
		DotsHorizontalOutline,
		DotsVerticalOutline
	} from 'flowbite-svelte-icons';

	export let data: {
		error?: string;
		binID?: number;
	};
	let bin: Bin | undefined = undefined;
	binsStore.subscribe((vals) => (bin = vals.find((b) => b.id === data.binID)));

	let hideDocumentCreator = true;
	let hideBinEditor = true;
</script>

{#if bin && bin.id}
	<Breadcrumb class="mb-2" solid>
		<BreadcrumbItem href="/" home>Home</BreadcrumbItem>
		<BreadcrumbItem>{bin.name}</BreadcrumbItem>
	</Breadcrumb>
	<div class="flex items-end gap-2">
		<p class="text-2xl font-bold">{bin?.name}</p>
		<button
			on:click={() => (hideDocumentCreator = false)}
			class="text-primary-400 hover:text-primary-500 dark:text-primary-600 dark:hover:text-primary-500 rounded-lg"
		>
			<CirclePlusSolid size="xl" />
		</button>
		<button class="rounded-md p-1 hover:bg-gray-100 hover:dark:bg-gray-800"
			><DotsHorizontalOutline class="h-6 w-6" /></button
		>
		<Dropdown trigger="hover">
			<DropdownItem
				on:click={() => {
					console.log('here');

					hideBinEditor = false;
				}}>Edit Bin</DropdownItem
			>
			<DropdownItem>Delete Bin</DropdownItem>
		</Dropdown>
	</div>
	<CreateDocument bind:binID={bin.id} bind:hidden={hideDocumentCreator} />
	<EditBin bind:hidden={hideBinEditor} editingBin={bin} />
{:else}
	<Skeleton size="sm" class="h-4 overflow-clip" />
{/if}
<Documents bin={bin?.name} />
