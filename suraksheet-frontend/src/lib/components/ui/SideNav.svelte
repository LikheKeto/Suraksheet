<script>
	import { page } from '$app/stores';
	import { token } from '$lib/store';
	import {
		Drawer,
		Sidebar,
		SidebarBrand,
		SidebarGroup,
		SidebarItem,
		SidebarWrapper,
		SidebarCta,
		ArrowKeyLeft,
		Tooltip
	} from 'flowbite-svelte';
	import {
		ChartPieSolid,
		UserSolid,
		ArrowRightToBracketOutline,
		ProfileCardSolid,
		UserSettingsOutline,
		BarsFromLeftOutline,
		CaretLeftOutline
	} from 'flowbite-svelte-icons';
	import { sineIn } from 'svelte/easing';
	import { goto } from '$app/navigation';
	import { fade, fly, slide } from 'svelte/transition';
	import { tweened } from 'svelte/motion';

	function goBack() {
		window.history.back(); // Use browser history to go back
	}

	$: isHomePage = $page.url.pathname === '/';

	$: activeUrl = $page.url.pathname;

	let hidden2 = true;

	let site = {
		name: 'Suraksheet',
		href: '/',
		img: '/favicon.png'
	};

	let transitionParams = {
		x: -320,
		duration: 200,
		easing: sineIn
	};

	let brandSpanClass =
		'self-center text-gray-700 text-2xl font-bold whitespace-nowrap dark:text-gray-300';
	let spanClass = 'ms-3 text-gray-600 dark:text-gray-400 font-semibold';
	let iconClass =
		'w-6 h-6 text-gray-500 transition duration-75 dark:text-gray-400 group-hover:text-gray-900 dark:group-hover:text-white';

	function signOut() {
		token.set('');
		window.location.replace('/login');
	}
</script>

<div class="fixed left-5 top-5 flex flex-row text-center sm:flex-col">
	<button
		id="opensidebar"
		class="rounded-lg p-2 text-xl text-gray-500 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-700"
		on:click={() => (hidden2 = false)}
	>
		<BarsFromLeftOutline size="xl" />
	</button>

	{#if !isHomePage}
		<button
			transition:fly={{ x: -100, duration: 500 }}
			id="goback"
			class="rounded-lg p-2 text-xl text-gray-500 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-700"
			on:click={goBack}
		>
			<CaretLeftOutline size="xl" />
		</button>
	{/if}
</div>
<Drawer
	width="w-64"
	divClass="overflow-y-auto z-50 bg-white dark:bg-gray-800"
	transitionType="fly"
	{transitionParams}
	bind:hidden={hidden2}
	id="sidebar2"
>
	<Sidebar asideClass="max-w-64" {activeUrl}>
		<SidebarWrapper>
			<SidebarGroup>
				<SidebarBrand spanClass={brandSpanClass} {site} />
				<hr class="h-px border-0 bg-gray-300 dark:bg-gray-700" />
				<SidebarItem href="/" {spanClass} label="Dashboard">
					<svelte:fragment slot="icon">
						<ChartPieSolid class={iconClass} />
					</svelte:fragment>
				</SidebarItem>
				<SidebarItem {spanClass} label="Details">
					<svelte:fragment slot="icon">
						<ProfileCardSolid class={iconClass} />
					</svelte:fragment>
				</SidebarItem>
				<SidebarItem {spanClass} label="Profile">
					<svelte:fragment slot="icon">
						<UserSolid class={iconClass} />
					</svelte:fragment>
				</SidebarItem>
				<SidebarItem {spanClass} label="Settings">
					<svelte:fragment slot="icon">
						<UserSettingsOutline class={iconClass} />
					</svelte:fragment>
				</SidebarItem>
				<SidebarItem on:click={signOut} {spanClass} label="Sign Out">
					<svelte:fragment slot="icon">
						<ArrowRightToBracketOutline class={iconClass} />
					</svelte:fragment>
				</SidebarItem>
			</SidebarGroup>
		</SidebarWrapper>
	</Sidebar>
</Drawer>
