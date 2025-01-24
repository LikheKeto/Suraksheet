import { PUBLIC_SERVER_URL } from "$env/static/public";
import { assetsStore, token } from "./store";
import type { Document } from "./types";

export async function getAsset(id: number) {
    let assets: Record<number, Blob> = {}
    assetsStore.subscribe(ass => assets = ass)
    let asset = assets[id]
    if (asset) {
        let objectURL = URL.createObjectURL(asset)
        return { url: objectURL, size: asset.size, type: asset.type }
    }

    let tok: string = "";
    token.subscribe(t => tok = t)
    let res = await fetch(PUBLIC_SERVER_URL + '/document/' + id + '/asset', {
        headers: {
            Authorization: 'Bearer ' + tok
        }
    });
    if (res.ok) {
        let blob = await res.blob();
        assetsStore.update(ass => {
            ass[id] = blob
            return ass
        })
        let objectURL = URL.createObjectURL(blob);
        return { url: objectURL, size: blob.size, type: blob.type };
    } else {
        throw new Error('Failed to fetch image');
    }
}


export function getReadableFileSizeString(fileSizeInBytes: number) {
    var i = -1;
    var byteUnits = [' KB', ' MB', ' GB', ' TB', 'PB', 'EB', 'ZB', 'YB'];
    do {
        fileSizeInBytes /= 1024;
        i++;
    } while (fileSizeInBytes > 1024);

    return Math.max(fileSizeInBytes, 0.1).toFixed(1) + byteUnits[i];
}

export function downloadFileFromUrl(doc: Document | null, url: string | null, compressed = false) {
    if (doc && url) {
        let extension = doc.name.split('.')[doc.name.split('.').length - 1]
        let fileName = doc.referenceName + (compressed ? "_compressed" : "") + '.' + extension
        const link = document.createElement('a');
        link.href = url;
        link.download = fileName;
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
    }
}