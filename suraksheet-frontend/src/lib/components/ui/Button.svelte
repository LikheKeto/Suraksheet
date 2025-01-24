<script lang="ts">
	import { Button } from 'flowbite-svelte';
	import { Spinner } from 'flowbite-svelte';
	import { createEventDispatcher } from 'svelte';

	export let loading: boolean = false; // Whether the button is in loading state
	export let loadingText: string = 'Loading...'; // Text to show while loading
	export let disabled: boolean = false; // To disable the button manually

	const dispatch = createEventDispatcher(); // Dispatcher to forward events

	// Combine loading and disabled state for the button
	$: internalDisabled = loading || disabled;

	function handleClick(event: Event) {
		// Prevent handling click events when loading
		if (!loading) {
			dispatch('click', event); // Forward the click event
		}
	}
</script>

<Button
	{...$$restProps}
	class={`flex items-center justify-center ${$$restProps.class || ''}`}
	disabled={internalDisabled}
	on:click={handleClick}
>
	{#if loading}
		<span class="flex items-center">
			<Spinner class="mr-2 h-5 w-5" />
			<span>{loadingText}</span>
		</span>
	{:else}
		<slot />
	{/if}
</Button>
