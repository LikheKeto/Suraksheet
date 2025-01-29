<script lang="ts">
	import type { Bin, Document } from '$lib/types';
	import { binsStore, documentsStore } from '$lib/store';
	import {
		Modal,
		ImagePlaceholder,
		Skeleton,
		Tabs,
		TabItem,
		Breadcrumb,
		BreadcrumbItem
	} from 'flowbite-svelte';
	import { onMount } from 'svelte';
	import { downloadFileFromUrl, getAsset } from '$lib';
	import Compressor from '$lib/components/documents/Compressor.svelte';
	import DocumentDetails from '$lib/components/documents/DocumentDetails.svelte';
	import { DownloadSolid, TrashBinSolid } from 'flowbite-svelte-icons';
	import Confirmer from '$lib/components/popups/Confirmer.svelte';
	import { deleteDocument } from '$lib/crud';
	import { goto } from '$app/navigation';
	import Button from '$lib/components/ui/Button.svelte';
	import { Popup } from '$lib/components/popups/popup';

	export let data: {
		error?: string;
		doc?: Document;
		docID?: number;
	};
	let docs: Record<string, Document[]> = {};
	documentsStore.subscribe((vals) => (docs = vals));

	let bins: Bin[] = [];
	binsStore.subscribe((b) => (bins = b));

	let bin: Bin | undefined = undefined;
	let doc: Document | null = null;

	let docSize = 0;
	let docType = '';
	let showDeletion = false;
	$: {
		if (!data.doc && !data.error && data.docID) {
			Object.keys(docs).forEach((docName) => {
				let d = docs[docName].find((d) => d.id === data.docID);
				if (d) {
					doc = d;
				}
			});
		} else if (data.doc) {
			doc = data.doc;
		} else {
			console.log(data.error);
		}
		if (doc) {
			bin = bins.find((b) => b.id === doc?.bin);
		}
	}

	onMount(async () => {
		if (doc) {
			let { url, size, type } = await getAsset(doc.id);
			doc.url = url;
			docSize = size;
			docType = type;
		}
	});

	let deleting = false;
	async function handleDeletion() {
		showDeletion = false;
		deleting = true;
		if (doc) {
			try {
				await deleteDocument(doc);
				goto('/bin/' + doc.bin);
				Popup('Success', 'Document deleted successfully');
			} catch (e) {
				console.log(e);
			} finally {
				deleting = false;
			}
		}
	}

	let showEditor = false;
</script>

<Modal title="Edit document asset" bind:open={showEditor}>
	<p>This feature is currently not implemented</p>
</Modal>

{#if doc}
	<Breadcrumb class="mb-2" solid>
		<BreadcrumbItem href="/" home>Home</BreadcrumbItem>
		<BreadcrumbItem href={'/bin/' + doc.bin}>{bin?.name}</BreadcrumbItem>
		<BreadcrumbItem>{doc.referenceName}</BreadcrumbItem>
	</Breadcrumb>
	<p class="pb-4 text-2xl font-bold">{doc?.referenceName}</p>
	<div class="flex flex-wrap gap-4">
		{#if doc.url}
			<img src={doc.url} alt={doc.name} class="w-2xl max-h-96 object-contain" />
		{:else}
			<ImagePlaceholder imgOnly class="m-auto mt-8 w-full" />
		{/if}
		<div class="flex flex-row flex-wrap gap-2 md:flex-col">
			<Button on:click={() => downloadFileFromUrl(doc, (doc && doc.url) || null)}>
				<DownloadSolid class="me-2 h-5 w-5" /> Download
			</Button>
			<Button
				loading={deleting}
				loadingText="Deleting document..."
				on:click={() => {
					showDeletion = true;
				}}
				outline
				color="red"
			>
				<TrashBinSolid class="me-2 h-5 w-5" /> Delete document
			</Button>
			<Confirmer
				message="Are you sure you want to delete this document?"
				onConfirmation={handleDeletion}
				bind:showConfirmer={showDeletion}
			/>
		</div>
	</div>
	<Tabs tabStyle="underline">
		<TabItem open title="Details">
			<DocumentDetails {doc} {docSize} {docType} {bin} />
		</TabItem>
		<TabItem title="Edit">
			<p class="text-sm text-gray-500 dark:text-gray-400">This Feature is not yet implemented</p>
		</TabItem>
		<TabItem title="Compress">
			<Compressor {doc} />
		</TabItem>
	</Tabs>
{:else}
	<Skeleton size="sm" class="h-4 overflow-clip" />
	<ImagePlaceholder class="mt-8" />
{/if}
