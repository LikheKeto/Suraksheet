import { browser } from "$app/environment";
import { writable } from "svelte/store";
import type { Bin, Document } from "./types";

export const token = writable((browser && localStorage.getItem('token') || "")||"")

token.subscribe(val=>{
    if(browser)return(localStorage.setItem('token', val))
})

export const binsStore = writable<Bin[]>([])
export const loadingBins = writable(true)

export const documentsStore = writable<Record<string, Document[]>>({})
export const loadingDocuments = writable(true)

export const assetsStore = writable<Record<number, Blob>>({})