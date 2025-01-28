import { PUBLIC_SERVER_URL } from '$env/static/public';
import { documentsStore, loadingDocuments, token } from '$lib/store';
import type { Document } from '$lib/types.js';

export async function load({ params, fetch, parent }) {
	await parent();
	loadingDocuments.set(true);
	let tok;
	let docs: Record<string, Document[]> = {};
	token.subscribe((val) => (tok = val));
	documentsStore.subscribe((vals) => (docs = vals));

	let docID = 0;
	try {
		docID = parseInt(params.slug);
	} catch {
		return { error: 'Invalid Document' };
	}
	let doQuery = true;

	Object.keys(docs).forEach((docName) => {
		let found = docs[docName].find((d) => d.id === docID);
		if (found) {
			loadingDocuments.set(false);
			doQuery = false;
		}
	});

	if (!doQuery) {
		return { doc: null, docID: docID };
	}
	let docRes = await fetch(PUBLIC_SERVER_URL + '/document/' + docID, {
		headers: {
			Authorization: 'Bearer ' + tok
		}
	});
	if (!docRes.ok) {
		return { error: 'Unable to fetch document' };
	}
	let res = await docRes.json();
	loadingDocuments.set(false);
	return { doc: { url: null, ...res } };
}
