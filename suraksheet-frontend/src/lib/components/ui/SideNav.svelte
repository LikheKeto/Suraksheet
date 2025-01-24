<script>
    import {page} from '$app/stores'
	import { token } from '$lib/store';
    import { Drawer, CloseButton, Button, Sidebar, SidebarBrand, SidebarGroup, SidebarItem, SidebarWrapper, SidebarCta } from 'flowbite-svelte';
    import {ChartPieSolid, GridSolid, MailBoxSolid, UserSolid, ArrowRightToBracketOutline, EditOutline, CloseOutline, ProfileCardSolid, UserSettingsOutline, BarsFromLeftOutline } from 'flowbite-svelte-icons';
    import {sineIn} from 'svelte/easing'
    
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

  let brandSpanClass = 'self-center text-gray-700 text-2xl font-bold whitespace-nowrap dark:text-gray-300'
  let spanClass = 'ms-3 text-gray-600 dark:text-gray-400 font-semibold'
  let iconClass ='w-6 h-6 text-gray-500 transition duration-75 dark:text-gray-400 group-hover:text-gray-900 dark:group-hover:text-white'

  function signOut(){
    token.set("")
    window.location.replace('/login')
  }
</script>

<div class="text-center">
    <button class="text-gray-500 fixed left-5 top-5 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg text-xl p-2" on:click={() => (hidden2 = false)}>
        <BarsFromLeftOutline size="xl" />
    </button>
</div>
<Drawer width='w-64' divClass="overflow-y-auto z-50 bg-white dark:bg-gray-800" transitionType="fly" {transitionParams} bind:hidden={hidden2} id="sidebar2">
<Sidebar asideClass="max-w-64" {activeUrl}>
    <SidebarWrapper>
        <SidebarGroup>
            <SidebarBrand spanClass={brandSpanClass} {site} />
            <hr class="h-px bg-gray-300 border-0 dark:bg-gray-700">
            <SidebarItem href='/' {spanClass} label="Dashboard" >
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