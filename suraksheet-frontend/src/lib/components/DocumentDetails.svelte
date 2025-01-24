<script lang="ts">
	import { getReadableFileSizeString } from '$lib';
	import { binsStore } from '$lib/store';
	import type { Bin, Document } from '$lib/types';
	import {
		Table,
		TableBody,
		TableBodyCell,
		TableBodyRow,
	} from 'flowbite-svelte';
	import {
		ArchiveOutline,
		BarsFromLeftOutline,
		CalendarMonthOutline,
		FileCloneOutline,
		FloppyDiskAltOutline,
		FileImageOutline,

		FilePdfSolid,

		BookOpenOutline


	} from 'flowbite-svelte-icons';

	export let doc: Document;
	export let docSize: number;
    export let docType:string;
	export let bin: Bin|undefined;

	let readingAll = false
</script>

<style>
    .text-wrap {
        white-space: pre-wrap; /* Allows text to wrap */
        word-wrap: break-word; /* Prevents overflow */
    }

    .collapsed {
        max-height: 1rem; /* Adjust the height as needed */
        overflow: hidden;
    }

    .expanded {
        max-height: none;
    }
</style>

<Table>
	<TableBody tableBodyClass="divide-y">
		<TableBodyRow>
			<TableBodyCell class="flex items-center gap-1"
				><FileCloneOutline class="text-primary-600 dark:text-primary-300" size="sm" /> Name</TableBodyCell
			>
			<TableBodyCell>{doc.referenceName}</TableBodyCell>
		</TableBodyRow>
		<TableBodyRow>
			<TableBodyCell class="flex items-center gap-1"
				><FileImageOutline class="text-primary-600 dark:text-primary-300" size="sm" />Original name</TableBodyCell
			>
			<TableBodyCell>{doc.name}</TableBodyCell>
		</TableBodyRow>
		<TableBodyRow>
			<TableBodyCell class="flex items-center gap-1"
				><ArchiveOutline
					class="text-primary-600 dark:text-primary-300"
					size="sm"
				/>Bin</TableBodyCell
			>
			<TableBodyCell>
				<a class="underline" href={'/bin/' + bin?.id}>
					{bin?.name}
				</a>
			</TableBodyCell>
		</TableBodyRow>
		<TableBodyRow>
			<TableBodyCell class="flex items-center gap-1"
				><FloppyDiskAltOutline class="text-primary-600 dark:text-primary-300" size="sm" />File size</TableBodyCell
			>
			<TableBodyCell>{getReadableFileSizeString(docSize)}</TableBodyCell>
		</TableBodyRow>
        <TableBodyRow>
			<TableBodyCell class="flex items-center gap-1"
				><FilePdfSolid class="text-primary-600 dark:text-primary-300" size="sm" />File type</TableBodyCell
			>
			<TableBodyCell>{docType}</TableBodyCell>
		</TableBodyRow>
		<TableBodyRow>
			<TableBodyCell class="flex items-center gap-1"
				><CalendarMonthOutline class="text-primary-600 dark:text-primary-300" size="sm" />Created at</TableBodyCell
			>
			<TableBodyCell>{new Date(doc.createdAt).toLocaleString()}</TableBodyCell>
		</TableBodyRow>
		<TableBodyRow>
			<TableBodyCell class="flex items-center gap-1"
				><BookOpenOutline class="text-primary-600 dark:text-primary-300" size="sm" />Extract</TableBodyCell
			>
			<TableBodyCell class="text-wrap">
				<!-- Conditionally toggle between collapsed and expanded content -->
				<div class={readingAll ? 'expanded' : 'collapsed'}>
					{doc.extract || "Extract not found"}
				</div>
				<button class="text-blue-500 underline" on:click={() => {
					readingAll = !readingAll;
				}}>
					{readingAll ? "show less" : "show more"}
				</button>
			</TableBodyCell>
		</TableBodyRow>
	</TableBody>
</Table>
