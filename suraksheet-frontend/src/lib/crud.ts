import { PUBLIC_SERVER_URL } from '$env/static/public';
import { binsStore, documentsStore, token } from './store';
import type { Bin, Document } from './types';

function deletionError(str: string) {
	return 'Unable to delete: ' + str;
}
export async function deleteDocument(doc: Document) {
	if (!doc || !doc.id) {
		throw deletionError('Invalid document');
	}
	let bins: Bin[] = [];
	binsStore.subscribe((b) => (bins = b));
	let bin = bins.find((b) => b.id === doc.bin);
	if (!bin) {
		throw deletionError('Invalid bin');
	}

	let tok: string = '';
	token.subscribe((t) => (tok = t));

	let res = await fetch(PUBLIC_SERVER_URL + '/document', {
		method: 'DELETE',
		body: JSON.stringify({ id: doc.id }),
		headers: {
			Authorization: 'Bearer ' + tok
		}
	});
	if (!res.ok) {
		throw deletionError((await res.json())['error']);
	}
	if (Object.keys(documentsStore).includes(bin.name)) {
		documentsStore.update((docs) => {
			docs[bin.name] = docs[bin.name].filter((d) => d.id !== doc.id);
			return docs;
		});
	}
	return true;
}
