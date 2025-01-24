<script>
	import { PUBLIC_SERVER_URL } from '$env/static/public';
	import { binsStore, token } from '$lib/store';
	import { Drawer, Button, CloseButton, Label, Input, Textarea, Helper } from 'flowbite-svelte';
	import { InfoCircleSolid, UserAddOutline, CalendarEditSolid } from 'flowbite-svelte-icons';
	import { sineIn } from 'svelte/easing';

	export let hidden = true;
	let transitionParams = {
		x: -320,
		duration: 200,
		easing: sineIn
	};

	let binName = '';
	let nameError = '';

	async function handleCreateBin() {
		if (binName.length < 3) {
			nameError = 'Bin name must be at least 3 characters.';
			return;
		}
		nameError = '';
		let resp = await fetch(PUBLIC_SERVER_URL + '/bins', {
			method: 'post',
			headers: {
				'Content-Type': 'application/json',
				Authorization: 'Bearer ' + $token
			},
			body: JSON.stringify({
				name: binName
			})
		});
		if (!resp.ok) {
			nameError = (await resp.json())['error'];
		} else {
			window.location.reload();
		}
	}
</script>

<Drawer transitionType="fly" {transitionParams} bind:hidden id="sidebar4">
	<div class="flex items-center">
		<h5
			id="drawer-label"
			class="mb-6 inline-flex items-center text-base font-semibold uppercase text-gray-500 dark:text-gray-400"
		>
			<InfoCircleSolid class="me-2.5 h-5 w-5" />New Bin
		</h5>
		<CloseButton on:click={() => (hidden = true)} class="mb-4" />
	</div>
	<form on:submit|preventDefault={handleCreateBin} class="mb-6">
		<div class="mb-6">
			<Label for="Name" class="mb-2 block">Name</Label>
			<Input bind:value={binName} id="name" name="name" required placeholder="College docs" />
			{#if nameError}
				<Helper color="red">{nameError}</Helper>
			{/if}
		</div>
		<Button color="primary" type="submit" class="w-full">
			<CalendarEditSolid class="me-2.5 h-3.5 w-3.5" /> Create bin
		</Button>
	</form>
</Drawer>
