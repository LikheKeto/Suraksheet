import { PUBLIC_SERVER_URL } from "$env/static/public";
import { binsStore, documentsStore, loadingBins, loadingDocuments, token } from "$lib/store";
import type { Bin, Document } from "$lib/types.js";

/** @type {import('./$types').PageLoad} */
export async function load({ fetch }) {
    let tok;
    token.subscribe(val => tok = val)

    if (tok) {
        let binsRes = await fetch(PUBLIC_SERVER_URL + '/bins', {
            headers: {
                "Authorization": "Bearer " + tok
            }
        })
        if (binsRes.ok) {
            let bins: Bin[] = await binsRes.json()
            binsStore.set(bins)
            loadingBins.set(false)

            let docsRes = await fetch(PUBLIC_SERVER_URL + '/bins/' + bins.find(bin => bin.name === 'No Bin')?.id, {
                headers: {
                    "Authorization": "Bearer " + tok
                }
            })
            if (docsRes.ok) {
                let docs: Document[] = await docsRes.json()
                documentsStore.set({ "No Bin": docs.map(d => ({ ...d, url: null })) })
                loadingDocuments.set(false)
            }
        }
    }
}