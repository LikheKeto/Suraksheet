<script lang="ts">
	import { binsStore, loadingBins } from '$lib/store';
	import type { Bin } from '$lib/types';
	import { Card, Skeleton, Tooltip } from 'flowbite-svelte';
	let bins: Bin[] = [];
	binsStore.subscribe((values) => (bins = values.filter((bin) => bin.name != 'No Bin')));
</script>

<div
	class="my-2 flex w-full flex-wrap gap-2 rounded-md bg-gray-200 bg-opacity-50 p-4 dark:bg-gray-800"
>
	{#if $loadingBins}
		<Skeleton divClass="animate-pulse w-full text-primary-500" size="lg" />
	{:else if bins.length > 0}
		{#each bins as bin, id}
			<Card
				href={'/bin/' + bin.id}
				class="flex h-16 w-full flex-row items-center justify-between gap-2 p-4 transition duration-300 ease-in-out sm:h-24 sm:w-32 sm:flex-col sm:justify-center sm:gap-0 sm:p-0"
			>
				<svg
					class="text-primary-500 dark:text-primary-400 h-12 w-12"
					aria-hidden="true"
					xmlns="http://www.w3.org/2000/svg"
					width="24"
					height="24"
					fill="currentColor"
					viewBox="0 0 24 24"
				>
					<path
						fill-rule="evenodd"
						d="M4 4a2 2 0 0 0-2 2v12a2 2 0 0 0 .087.586l2.977-7.937A1 1 0 0 1 6 10h12V9a2 2 0 0 0-2-2h-4.532l-1.9-2.28A2 2 0 0 0 8.032 4H4Zm2.693 8H6.5l-3 8H18l3-8H6.693Z"
						clip-rule="evenodd"
					/>
				</svg>
				<p
					id={`bin-${id}`}
					class="mt-2 w-32 truncate text-ellipsis text-right text-sm text-gray-700 sm:text-center md:w-20 dark:text-gray-300"
				>
					{bin.name}
				</p>
				<Tooltip class="absolute w-max" type="auto" triggeredBy={`#bin-${id}`} placement="bottom">
					{bin.name}
				</Tooltip>
			</Card>
		{/each}
	{:else}
		<div class="flex h-24 w-full cursor-pointer flex-col items-center justify-center rounded-md">
			<p class="text-gray-400">You have no bins</p>
		</div>
	{/if}
</div>
