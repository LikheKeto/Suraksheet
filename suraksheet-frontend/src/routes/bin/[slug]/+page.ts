import { PUBLIC_SERVER_URL } from '$env/static/public';
import { binsStore, documentsStore, loadingDocuments, token } from '$lib/store';
import type { Bin, Document } from '$lib/types.js';

/** @type {import('./$types').PageLoad} */
export async function load({ params, fetch, parent }) {
	await parent();

	let tok;
	let bins: Bin[] = [];
	token.subscribe((val) => (tok = val));
	binsStore.subscribe((vals) => (bins = vals));

	let binID = 0;
	try {
		binID = parseInt(params.slug);
	} catch {
		return { error: 'Invalid Bin' };
	}

	let bin = bins.find((b) => b.id === binID);

	if (!bin) {
		return { error: 'Invalid Bin' };
	}

	let preexistingDocs: Record<string, Document[]> = {};
	documentsStore.subscribe((d) => (preexistingDocs = d));

	if (!Object.keys(preexistingDocs).includes(bin.name)) {
		let docsRes = await fetch(PUBLIC_SERVER_URL + '/bins/' + binID, {
			headers: {
				Authorization: 'Bearer ' + tok
			}
		});

		if (!docsRes.ok) {
			return { error: 'Unable to fetch documents' };
		}

		let docs: Document[] = await docsRes.json();
		documentsStore.update((vals) => {
			vals[bin.name] = docs;
			return vals;
		});
	}

	return {
		binID: bin.id
	};
}
