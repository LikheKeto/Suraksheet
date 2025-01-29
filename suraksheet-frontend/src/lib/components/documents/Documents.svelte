<script lang="ts">
	import { documentsStore, loadingDocuments, token } from '$lib/store';
	import type { Document } from '$lib/types';
	import { Card, ImagePlaceholder, Skeleton, Tooltip } from 'flowbite-svelte';
	import { hydrateImages } from '$lib/index';

	export let bin = 'No Bin';

	let docs: Document[] = [];
	documentsStore.subscribe(async (values) => (docs = await hydrateImages(values[bin])));
</script>

<div
	class="my-2 flex w-full flex-wrap gap-2 rounded-md bg-gray-100 bg-opacity-50 p-4 dark:bg-gray-800"
>
	{#if docs.length > 0}
		{#each docs as doc, id}
			<Card
				href={'/document/' + doc.id}
				class="flex h-24 w-full cursor-pointer flex-row items-center justify-between gap-2 rounded-md p-4 shadow-md transition duration-300 ease-in-out sm:h-32 sm:w-32 sm:flex-col sm:justify-center sm:gap-0 sm:p-0"
			>
				{#if doc.url}
					<img src={doc.url} alt={doc.name} class="h-20 w-24 object-contain" />
					<Tooltip type="auto" triggeredBy={`#doc-${id}`} placement="bottom">
						{doc.referenceName}
					</Tooltip>
				{:else}
					<ImagePlaceholder imgOnly class="w-24" imgHeight="16" />
				{/if}
				<p
					id={`doc-${id}`}
					class="mt-2 w-36 truncate text-ellipsis text-right text-sm text-gray-700 sm:w-32 sm:text-center dark:text-gray-300"
				>
					{doc.referenceName}
				</p>
			</Card>
		{/each}
	{:else if loadingDocuments}
		<Skeleton divClass="animate-pulse w-full text-primary-500" size="lg" />
	{:else}
		<div class="flex h-24 w-full cursor-pointer flex-col items-center justify-center rounded-md">
			<p class="text-gray-400">You have no documents</p>
		</div>
	{/if}
</div>
