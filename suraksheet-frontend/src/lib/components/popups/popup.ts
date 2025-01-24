import { writable } from 'svelte/store';

export const toasts = writable<any[]>([]);

export function Popup(
	type: 'Error' | 'Success' | 'Warning',
	message: string,
	duration: number = 3000
) {
	const id = Math.random().toString(36).substring(2);
	const toast = { id, type, message, duration };
	toasts.update((current) => [...current, toast]);

	// Remove the toast after duration
	if (duration > 0) {
		setTimeout(() => {
			removeToast(id);
		}, duration);
	}
}

export function removeToast(id: string) {
	toasts.update((current) => current.filter((toast) => toast.id !== id));
}
