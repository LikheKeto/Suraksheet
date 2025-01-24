<script>
	import SideNav from '$lib/components/ui/SideNav.svelte';
	import { DarkMode } from 'flowbite-svelte';
	import '../app.css';
	import ToastRenderer from '$lib/components/popups/ToastRenderer.svelte';
	import SearchBar from '$lib/components/SearchBar.svelte';
	import Tooltips from '$lib/components/ui/Tooltips.svelte';
</script>

<svelte:head>
	<title>Suraksheet - A safe place for your documents</title>
</svelte:head>

<SideNav />
<Tooltips />
<ToastRenderer position="top-left" />
<div class="fixed right-5 top-5">
	<DarkMode
		size="lg"
		btnClass="text-gray-500 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg text-xl p-2"
	/>
</div>
<SearchBar />
<div class="container m-auto w-full p-5 text-gray-700 dark:text-gray-300">
	<slot></slot>
</div>
