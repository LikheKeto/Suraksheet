<script lang="ts">
	import { onMount } from 'svelte';
	import { Input, Modal, Kbd, Card, Tooltip, ImagePlaceholder } from 'flowbite-svelte';
	import { SearchOutline } from 'flowbite-svelte-icons';
	import { PUBLIC_SERVER_URL } from '$env/static/public';
	import { token } from '$lib/store';
	import { hydrateImages } from '$lib';
	import type { Document } from '$lib/types';

	let isOpen = false;
	let searchQuery = '';
	let searchInput: HTMLInputElement;
	let searchResults: Document[] = [];
	let isLoading = false;
	let hasSearched = false;
	let error = '';

	const handleKeydown = (event: KeyboardEvent) => {
		if ((event.ctrlKey || event.metaKey) && event.key === 'k') {
			event.preventDefault();
			isOpen = !isOpen;
		}
		if (event.key === 'Escape') {
			isOpen = false;
		}
	};

	onMount(() => {
		window.addEventListener('keydown', handleKeydown);
		return () => window.removeEventListener('keydown', handleKeydown);
	});

	const performSearch = async () => {
		if (!searchQuery.trim()) return;
		isLoading = true;
		error = '';
		hasSearched = true;
		try {
			const response = await fetch(
				`${PUBLIC_SERVER_URL}/document/search?q=${encodeURIComponent(searchQuery)}`,
				{
					headers: {
						Authorization: `Bearer ${$token}`
					}
				}
			);
			if (!response.ok) throw new Error('Failed to fetch search results');
			const data = await response.json();
			searchResults = data;
			searchResults = await hydrateImages(searchResults);
		} catch (err: any) {
			error = err['error'] || 'search failed';
		} finally {
			isLoading = false;
		}
	};

	$: isOpen && searchInput && searchInput.focus();
</script>

<div class="container relative m-auto mt-20 max-w-lg gap-2 px-4 sm:mt-8">
	<div class="pointer-events-none absolute inset-y-0 start-0 flex items-center ps-7">
		<SearchOutline class="h-4 w-4 text-slate-900 dark:text-slate-100" />
	</div>
	<Input
		id="search"
		on:click={() => (isOpen = true)}
		type="text"
		class="ps-10"
		placeholder="Search... (Ctrl + k)"
	/>
</div>

<Modal outsideclose autoclose bind:open={isOpen} size="lg" placement="center">
	<div class="space-y-4 p-4">
		<div class="md-block relative gap-2">
			<div class="pointer-events-none absolute inset-y-0 start-0 flex items-center ps-3">
				<SearchOutline class="h-4 w-4" />
			</div>
			<input
				type="text"
				bind:this={searchInput}
				class="border-1 focus:ring-primary-500 w-full rounded-md border-gray-500 bg-gray-100 ps-10 text-gray-700 placeholder:text-gray-400 focus:ring-2 dark:bg-gray-600 dark:text-gray-300"
				bind:value={searchQuery}
				placeholder="Search..."
				on:keydown={(e) => e.key === 'Enter' && performSearch()}
			/>
		</div>
		<div class="text-sm text-gray-500 dark:text-gray-400">
			Press <Kbd class="px-2 py-1.5">Enter</Kbd> to search or <Kbd class="px-2 py-1.5">Esc</Kbd> to close.
		</div>
		<hr class="border-slate-500" />

		{#if isLoading}
			<p class="text-gray-400">Loading...</p>
		{:else if error}
			<p class="text-red-500">{error}</p>
		{:else if searchResults.length === 0 && hasSearched}
			<p class="text-gray-400">No results found.</p>
		{:else}
			<div
				class="flex w-full flex-wrap gap-2 rounded-md bg-gray-100 bg-opacity-50 dark:bg-gray-800"
			>
				{#each searchResults as doc, id}
					<Card
						on:click={() => {
							isOpen = false;
						}}
						href={'/document/' + doc.id}
						class="flex h-24 w-full cursor-pointer flex-row items-center justify-between gap-2 rounded-md p-4 shadow-md transition duration-300 ease-in-out sm:h-32 sm:w-32 sm:flex-col sm:justify-center sm:gap-0 sm:p-0"
					>
						{#if doc.url}
							<img src={doc.url} alt={doc.name} class="h-20 w-24 object-contain" />
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
			</div>
		{/if}
	</div>
</Modal>
