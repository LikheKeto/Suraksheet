<script lang="ts">
	import { PUBLIC_SERVER_URL } from '$env/static/public';
	import { binsStore, token } from '$lib/store';
	import type { Bin } from '$lib/types';
	import { Drawer, CloseButton, Label, Input, Select, Fileupload, Helper } from 'flowbite-svelte';
	import { InfoCircleSolid, UserAddOutline, CalendarEditSolid } from 'flowbite-svelte-icons';
	import { sineIn } from 'svelte/easing';
	import Button from './ui/Button.svelte';

	export let hidden = true;
	export let binID: number | undefined = undefined;
	let transitionParams = {
		x: -320,
		duration: 200,
		easing: sineIn
	};
	let bins: {
		value: string;
		name: string;
	}[] = [];
	binsStore.subscribe(
		(vals) =>
			(bins = vals.map((v) => {
				return {
					value: v.id.toString(),
					name: v.name
				};
			}))
	);

	let loading = false;
	let error = '';
	let name = '';
	let selectedBin = binID?.toString() || '';
	let files: FileList | undefined;

	async function handleCreateDocument() {
		loading = true;
		let file = files?.item(0);
		if (file) {
			if (file.size > 10 * 1024 * 1024) {
				error = 'File too big!';
				loading = false;
				return;
			}
			const formData = new FormData();
			formData.append('referenceName', name);
			formData.append('binID', selectedBin.toString());
			formData.append('file', file);

			let res = await fetch(PUBLIC_SERVER_URL + '/document', {
				method: 'POST',
				body: formData,
				headers: {
					Authorization: 'Bearer ' + $token
				}
			});
			if (!res.ok) {
				error = (await res.json())['error'];
				return;
			}
			error = '';
			loading = false;
			window.location.reload();
		}
	}
</script>

<Drawer transitionType="fly" {transitionParams} bind:hidden>
	<div class="flex items-center">
		<h5
			id="drawer-label"
			class="mb-6 inline-flex items-center text-base font-semibold uppercase text-gray-500 dark:text-gray-400"
		>
			<InfoCircleSolid class="me-2.5 h-5 w-5" />New document
		</h5>
		<CloseButton on:click={() => (hidden = true)} class="mb-4 dark:text-white" />
	</div>
	<form on:submit|preventDefault={handleCreateDocument} class="mb-6">
		<div class="mb-6">
			<Label for="name" class="mb-2 block">Name</Label>
			<Input bind:value={name} id="name" name="name" required placeholder="+2 Marksheet" />
		</div>
		<div class="mb-6">
			<Label class="pb-2"
				>Upload document <span class="text-xs text-gray-400">(png/jpeg/pdf)</span></Label
			>
			<Fileupload
				required
				bind:files
				accept="image/png, image/jpeg, image/gif, application/pdf, application/vnd.openxmlformats-officedocument.wordprocessingml.document"
			/>
		</div>
		<div class="mb-6">
			<Label>
				Bin
				<Select required class="mt-2" items={bins} bind:value={selectedBin} />
			</Label>
		</div>
		<Button bind:loading loadingText="Creating document..." type="submit" class="w-full">
			<CalendarEditSolid class="me-2.5 h-3.5 w-3.5 text-white" /> Create document
		</Button>
		{#if error}
			<Helper class="mt-4 text-center" color="red">{error}</Helper>
		{/if}
	</form>
</Drawer>
