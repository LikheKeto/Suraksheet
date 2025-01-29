<script>
	import SideNav from '$lib/components/ui/SideNav.svelte';
	import { DarkMode } from 'flowbite-svelte';
	import '../app.css';
	import ToastRenderer from '$lib/components/popups/ToastRenderer.svelte';
	import SearchBar from '$lib/components/SearchBar.svelte';
	import Tooltips from '$lib/components/ui/Tooltips.svelte';
	import { token } from '$lib/store';
	import { onMount } from 'svelte';

	let loggedIn = false;
	onMount(() => {
		if (localStorage.getItem('token')) {
			loggedIn = true;
		}
	});
	token.subscribe((tok) => {
		loggedIn = tok.length > 5;
	});
</script>

<svelte:head>
	<title>Suraksheet - A safe place for your documents</title>
</svelte:head>

{#if loggedIn}
	<SideNav />
	<SearchBar />
{:else}
	<div
		class="container m-auto my-8 flex items-end justify-center text-3xl font-bold text-gray-800 dark:text-gray-200"
	>
		<img src="/favicon.png" class="w-12" alt="suraksheet" />
		Suraksheet
	</div>
	<!-- <hr class="container m-auto mt-4 border-2 border-gray-500" /> -->
{/if}
<Tooltips />
<ToastRenderer position="top-left" />
<div class="fixed right-5 top-5">
	<DarkMode
		size="lg"
		btnClass="text-gray-500 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg text-xl p-2"
	/>
</div>
<div class="container m-auto w-full p-5 text-gray-700 dark:text-gray-300">
	<slot></slot>
</div>
