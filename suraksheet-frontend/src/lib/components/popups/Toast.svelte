<script lang="ts">
	import { Toast } from 'flowbite-svelte';
	import {
		CheckCircleSolid,
		ExclamationCircleSolid,
		FireOutline,
		CloseCircleSolid
	} from 'flowbite-svelte-icons';
	import { onMount } from 'svelte';
	import { fade, slide } from 'svelte/transition';

	export let type: 'Error' | 'Success' | 'Warning' = 'Success';
	export let message: string = '';
	export let duration: number = 3000; // Auto-destroy timeout
	export let onClose: () => void; // Callback for manual close

	let colors = {
		Error: 'red',
		Success: 'green',
		Warning: 'yellow'
	};
	let icons = {
		Error: CloseCircleSolid,
		Success: CheckCircleSolid,
		Warning: ExclamationCircleSolid
	};
	let Icon;

	// Assign icon based on type
	$: Icon = icons[type];

	let timeout;
	onMount(() => {
		if (duration > 0) {
			timeout = setTimeout(() => {
				handleClose();
			}, duration);
		}
	});

	function handleClose() {
		clearTimeout(timeout);
		if (onClose) onClose(); // Trigger callback
	}
</script>

<div class="w-full" transition:slide>
	<Toast color={colors[type]}>
		<svelte:fragment slot="icon">
			<Icon class="h-5 w-5" />
			<span class="sr-only">{type} icon</span>
		</svelte:fragment>
		<div class="flex items-center justify-between">
			<span>{message}</span>
			<!-- <button class="ml-4 text-lg" on:click={handleClose}>✖</button> -->
		</div>
	</Toast>
</div>
