<script lang="ts">
	import { updated } from '$app/stores';
	import { PUBLIC_SERVER_URL } from '$env/static/public';
	import { binsStore, token } from '$lib/store';
	import type { Bin } from '$lib/types';
	import { Drawer, Button, CloseButton, Label, Input, Textarea, Helper, P } from 'flowbite-svelte';
	import { InfoCircleSolid, UserAddOutline, CalendarEditSolid } from 'flowbite-svelte-icons';
	import { sineIn } from 'svelte/easing';
	import { Popup } from '../popups/popup';

	export let hidden = true;
	export let editingBin: Bin | undefined = undefined; // Prop for the bin being edited

	let transitionParams = {
		x: -320,
		duration: 200,
		easing: sineIn
	};

	let binName = editingBin?.name || ''; // Pre-fill with editing bin's name
	let nameError = '';

	async function handleCreateEditBin() {
		if (binName.length < 3) {
			nameError = 'Bin name must be at least 3 characters.';
			return;
		}
		nameError = '';

		const url = `${PUBLIC_SERVER_URL}/bins`;

		const method = editingBin ? 'PATCH' : 'POST';
		const body = editingBin ? { id: editingBin.id, name: binName } : { name: binName };

		let resp = await fetch(url, {
			method,
			headers: {
				'Content-Type': 'application/json',
				Authorization: 'Bearer ' + $token
			},
			body: JSON.stringify(body)
		});

		if (!resp.ok) {
			nameError = (await resp.json())['error'];
		} else {
			if (editingBin) {
				binsStore.update((bins) => {
					const updatedBin: Bin = { ...editingBin, name: binName };
					return bins.map((b) => (b.id === editingBin.id ? updatedBin : b));
				});
				Popup('Success', 'Bin updated successfully');
			} else {
				let newBin = await resp.json();
				binsStore.update((bins) => [...bins, newBin]);
				binName = '';
				Popup('Success', 'Bin created successfully');
			}
			hidden = true;
		}
	}
</script>

<Drawer transitionType="fly" {transitionParams} bind:hidden id="sidebar4">
	<div class="flex items-center">
		<h5
			id="drawer-label"
			class="mb-6 inline-flex items-center text-base font-semibold uppercase text-gray-500 dark:text-gray-400"
		>
			<InfoCircleSolid class="me-2.5 h-5 w-5" />
			{editingBin ? 'Edit Bin' : 'New Bin'}
		</h5>
		<CloseButton on:click={() => (hidden = true)} class="mb-4" />
	</div>
	<form on:submit|preventDefault={handleCreateEditBin} class="mb-6">
		<div class="mb-6">
			<Label for="Name" class="mb-2 block">Name</Label>
			<Input bind:value={binName} id="name" name="name" required placeholder="College docs" />
			{#if nameError}
				<Helper color="red">{nameError}</Helper>
			{/if}
		</div>
		<Button color="primary" type="submit" class="w-full">
			<CalendarEditSolid class="me-2.5 h-3.5 w-3.5" />
			{editingBin ? 'Update Bin' : 'Create Bin'}
		</Button>
	</form>
</Drawer>
