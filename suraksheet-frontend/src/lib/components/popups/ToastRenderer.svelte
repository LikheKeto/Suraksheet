<script lang="ts">
	import { toasts, removeToast } from './popup';
	import ToastComponent from './Toast.svelte';

	export let position: 'top-right' | 'top-left' | 'bottom-right' | 'bottom-left' | 'top-center' =
		'top-right';

	// Function to return Tailwind classes based on position
	function positionClasses(position: string) {
		return {
			'top-right': 'top-4 right-4',
			'top-left': 'top-4 left-4',
			'bottom-right': 'bottom-4 right-4',
			'bottom-left': 'bottom-4 left-4',
			'top-center': 'top-4 left-1/2 transform -translate-x-1/2' // Center alignment
		}[position];
	}
</script>

<div class={`fixed ${positionClasses(position)} space-y-2`}>
	{#each $toasts.slice(0, 3) as toast (toast.id)}
		<ToastComponent
			type={toast.type}
			message={toast.message}
			duration={toast.duration}
			onClose={() => removeToast(toast.id)}
		/>
	{/each}
</div>
