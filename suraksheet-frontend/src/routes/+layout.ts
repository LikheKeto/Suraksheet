import { goto } from '$app/navigation';
import { PUBLIC_SERVER_URL } from '$env/static/public';
import { binsStore, documentsStore, loadingBins, loadingDocuments, token } from '$lib/store';
import type { Bin, Document } from '$lib/types.js';
import { get } from 'svelte/store';

/** @type {import('./$types').PageLoad} */
export async function load({ fetch }) {
	const tok = get(token);
	let fetched = false;
	if (tok) {
		await fetchData(tok, fetch);
		fetched = true;
	}

	token.subscribe(async (newToken) => {
		if (newToken && !fetched) {
			await fetchData(newToken, fetch);
		}
	});
}

async function fetchData(tok: string, fetch: typeof window.fetch) {
	let binsRes = await fetch(PUBLIC_SERVER_URL + '/bins', {
		headers: {
			Authorization: 'Bearer ' + tok
		}
	});

	if (binsRes.ok) {
		let bins: Bin[] = await binsRes.json();
		binsStore.set(bins);
		loadingBins.set(false);

		const noBin = bins.find((bin) => bin.name === 'No Bin');
		if (noBin) {
			let docsRes = await fetch(PUBLIC_SERVER_URL + '/bins/' + noBin.id, {
				headers: {
					Authorization: 'Bearer ' + tok
				}
			});

			if (docsRes.ok) {
				let docs: Document[] = await docsRes.json();
				documentsStore.set({ 'No Bin': docs.map((d) => ({ ...d, url: null })) });
				loadingDocuments.set(false);
			}
		}
	}
	if (binsRes.status === 403) {
		localStorage.removeItem('token');
		window.location.replace('/login');
		// goto('/login');
	}
}
