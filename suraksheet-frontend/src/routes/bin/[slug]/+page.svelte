<script lang="ts">
	import CreateDocument from '$lib/components/documents/CreateDocument.svelte';
	import EditBin from '$lib/components/bins/CreateEditBin.svelte';
	import Documents from '$lib/components/documents/Documents.svelte';
	import { binsStore, documentsStore, token } from '$lib/store';
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
	import Confirmer from '$lib/components/popups/Confirmer.svelte';
	import { PUBLIC_SERVER_URL } from '$env/static/public';
	import { Popup } from '$lib/components/popups/popup';
	import { goto } from '$app/navigation';

	export let data: {
		error?: string;
		binID?: number;
	};
	let bin: Bin | undefined = undefined;
	binsStore.subscribe((vals) => (bin = vals.find((b) => b.id === data.binID)));

	let hideDocumentCreator = true;
	let hideBinEditor = true;
	let showConfirmer = false;

	async function deleteBin() {
		const url = `${PUBLIC_SERVER_URL}/bins`;
		let resp = await fetch(url, {
			method: 'DELETE',
			headers: {
				'Content-Type': 'application/json',
				Authorization: 'Bearer ' + $token
			},
			body: JSON.stringify({
				id: data.binID
			})
		});

		if (!resp.ok) {
			let err = await resp.json();
			if (err.error) Popup('Error', err.error);
		} else {
			showConfirmer = false;
			goto('/');
			binsStore.update((bins) => bins.filter((b) => b.id !== data.binID));
			documentsStore.update((docs) => {
				if (bin) {
					docs[bin.name] = [];
				}
				return docs;
			});
			Popup('Success', 'Bin deleted successfully');
		}
	}
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
					hideBinEditor = false;
				}}>Edit Bin</DropdownItem
			>
			<DropdownItem
				on:click={() => {
					showConfirmer = true;
				}}>Delete Bin</DropdownItem
			>
		</Dropdown>
	</div>
	<CreateDocument bind:binID={bin.id} bind:hidden={hideDocumentCreator} />
	<EditBin bind:hidden={hideBinEditor} editingBin={bin} />
	<Confirmer
		message="Are you sure you want to delete this bin?"
		onConfirmation={deleteBin}
		bind:showConfirmer
	/>
{:else}
	<Skeleton size="sm" class="h-4 overflow-clip" />
{/if}
<Documents bin={bin?.name} />
