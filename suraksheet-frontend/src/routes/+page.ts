import { browser } from '$app/environment';
import { token } from '$lib/store.js';
import { redirect } from '@sveltejs/kit';

export async function load({parent}) {
    if(browser){
        if(!localStorage.getItem('token')){
            throw redirect(302, '/login')
        }
    }
}
