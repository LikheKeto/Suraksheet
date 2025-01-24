<script lang="ts">
	import { downloadFileFromUrl, getAsset, getReadableFileSizeString } from '$lib';
	import type { Document } from '$lib/types';
	import imageCompression, { type Options } from 'browser-image-compression';
	import { Button, Label, Progressbar, Range, Toggle } from 'flowbite-svelte';
	import { onMount } from 'svelte';

	export let doc: Document;
	let originalImageUrl: string | null = null;
	let compressedImageUrl: string | null = null;
	let originalSize = 0;
	let compressedSize = 0;
	let progress = 0;
	let compressing = true;

	onMount(async () => {
		let { url, size } = await getAsset(doc.id);
		originalImageUrl = url;
		originalSize = size;
		compressedImageUrl = originalImageUrl;
		compressedSize = size;
	});

	let options: Options = {
		onProgress: (p) => {
			progress = p;
		},
		maxIteration: 10,
		maxSizeMB: 5,
		preserveExif: false,
		alwaysKeepResolution: false
	};

	let id: NodeJS.Timeout | null = null;
	let prevController: AbortController | null = null;
	async function handleCompression() {
		let abortController = new AbortController();
		options.signal = abortController.signal;
		if (id !== null) {
			prevController?.abort(new Error('user request'));
			clearTimeout(id);
		}
		id = setTimeout(async () => {
			if (originalImageUrl) {
				compressing = true;
				try {
					const response = await fetch(originalImageUrl);
					const blob = await response.blob();
					const file = new File([blob], doc.name, { type: blob.type });

					prevController = abortController;
					const compressedFile = await imageCompression(file, options);
					compressedSize = compressedFile.size;
					compressedImageUrl = await imageCompression.getDataUrlFromFile(compressedFile);
				} catch (error) {
					console.error('Error compressing image:', error);
				} finally {
					compressing = false;
				}
			}
		}, 500);
	}
</script>

<div
	class="flex flex-col items-center justify-center space-y-4 rounded-lg md:flex-row md:space-x-4 md:space-y-0"
>
	<!-- Original Image -->
	<div class="relative h-36 w-full max-w-md overflow-clip">
		<img
			src={originalImageUrl}
			alt="original"
			class="h-[500%] -translate-y-32 rounded-lg object-cover"
		/>
		<div
			class="absolute inset-0 flex items-center justify-center rounded-lg bg-black bg-opacity-30 p-2 text-sm text-white"
		>
			Original: {getReadableFileSizeString(originalSize)}
		</div>
	</div>

	<!-- Compressed Image -->
	<div class="relative h-36 w-full max-w-md overflow-clip">
		<img
			src={compressedImageUrl}
			alt="compressed"
			class="h-[500%] -translate-y-32 rounded-lg object-cover"
		/>
		<div
			class="pointer-events-none absolute inset-0 flex items-center justify-center rounded-lg bg-black bg-opacity-30 p-2 text-sm text-white"
		>
			Compressed: {getReadableFileSizeString(compressedSize)}
		</div>
	</div>
</div>
{#if progress != 0 && progress != 100}
	<Progressbar class="m-auto mt-4 max-w-96" {progress} size="h-4" labelInside />
{/if}
<div class="m-auto mt-4 max-w-96">
	<Label for="maxSize" class="my-2">Maximum size (MB)</Label>
	<div class="flex w-full items-center justify-between gap-4">
		<Range
			on:change={handleCompression}
			name="maxSize"
			class="w-[85%]"
			bind:value={options.maxSizeMB}
			max="10"
			min="0.1"
			step="0.1"
		/> <span class="w-[15%] whitespace-nowrap"> {options.maxSizeMB} MB </span>
	</div>
	<Label for="iters" class="my-2">Maximum iterations</Label>
	<div class="flex w-full items-center justify-between gap-4">
		<Range
			on:change={handleCompression}
			name="iters"
			class="w-[85%]"
			bind:value={options.maxIteration}
			max="20"
			min="1"
		/> <span class="w-[15%] whitespace-nowrap"> {options.maxIteration}</span>
	</div>
	<div class="my-2 flex w-full items-center justify-between gap-4">
		<Toggle
			on:change={handleCompression}
			name="preserve"
			bind:checked={options.alwaysKeepResolution}
		>
			{options.alwaysKeepResolution ? 'Persist' : 'Change'} width and height of image
		</Toggle>
	</div>
	<div class="my-2 flex w-full flex-col gap-4">
		<Toggle on:change={handleCompression} name="exif" bind:checked={options.preserveExif}>
			{options.preserveExif ? 'Preserve' : "Don't preserve"} EXIF information
		</Toggle>
		<Button
			bind:disabled={compressing}
			on:click={() => downloadFileFromUrl(doc, compressedImageUrl,true)}
			outline>Download</Button
		>
	</div>
</div>
