<script lang="ts">
	import { onMount } from 'svelte';
	import { Input, Modal, Kbd } from 'flowbite-svelte';
	import { SearchOutline } from 'flowbite-svelte-icons';

	let isOpen = false;
	let searchQuery = '';
	let searchInput: HTMLInputElement;

	const handleKeydown = (event: KeyboardEvent) => {
		if ((event.ctrlKey || event.metaKey) && event.key === 'k') {
			event.preventDefault(); // Prevent browser's default search behavior
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

	const closeModal = () => {
		isOpen = false;
	};

	const performSearch = () => {
		console.log('Searching for:', searchQuery);
		searchQuery = '';
		closeModal();
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
				class="border-1 focus:ring-primary-500 w-full rounded-md border-gray-500 bg-gray-600 ps-10 text-gray-300 placeholder:text-gray-400 focus:ring-2"
				bind:value={searchQuery}
				placeholder="Search..."
				on:keydown={(e) => {
					e.key === 'Enter' && performSearch();
				}}
			/>
		</div>
		<div class="text-sm text-gray-500 dark:text-gray-400">
			Press <Kbd class="px-2 py-1.5">Enter</Kbd> to search or <Kbd class="px-2 py-1.5">Esc</Kbd> to close.
		</div>
		<hr class="border-slate-500" />
		<div>
			<h2>Recent</h2>
		</div>
	</div>
</Modal>
