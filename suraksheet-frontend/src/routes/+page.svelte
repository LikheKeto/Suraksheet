<script lang="ts">
	import Bins from '$lib/components/bins/Bins.svelte';
	import CreateBin from '$lib/components/bins/CreateEditBin.svelte';
	import CreateDocument from '$lib/components/documents/CreateDocument.svelte';
	import Documents from '$lib/components/documents/Documents.svelte';
	import { binsStore } from '$lib/store';
	import { Breadcrumb, BreadcrumbItem } from 'flowbite-svelte';
	import { CirclePlusSolid } from 'flowbite-svelte-icons';

	let hideBinCreator = true;
	let hideDocumentCreator = true;
</script>

<Breadcrumb class="mb-2" solid>
	<BreadcrumbItem home>Home</BreadcrumbItem>
	<BreadcrumbItem></BreadcrumbItem>
</Breadcrumb>
<div class="flex items-center gap-2">
	<p class="text-2xl font-semibold">Bins</p>
	<button
		on:click={() => (hideBinCreator = false)}
		class="text-primary-400 hover:text-primary-500 dark:text-primary-600 dark:hover:text-primary-500 rounded-lg"
	>
		<CirclePlusSolid size="xl" />
	</button>
</div>
<Bins />
<CreateBin bind:hidden={hideBinCreator} />

<div class="flex items-center gap-2">
	<p class="text-2xl font-semibold">Documents</p>
	<button
		on:click={() => (hideDocumentCreator = false)}
		class="text-primary-400 hover:text-primary-500 dark:text-primary-600 dark:hover:text-primary-500 rounded-lg"
	>
		<CirclePlusSolid size="xl" />
	</button>
</div>
<Documents />
<CreateDocument
	bind:hidden={hideDocumentCreator}
	binID={$binsStore.find((bin) => bin.name === 'No Bin')?.id || undefined}
/>
