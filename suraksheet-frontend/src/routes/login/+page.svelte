<script lang="ts">
	import { Card, Label, Input, Helper, Toast } from 'flowbite-svelte';
	import { EnvelopeSolid, LockSolid } from 'flowbite-svelte-icons';
	import { PUBLIC_SERVER_URL } from '$env/static/public';
	import { token } from '$lib/store';
	import { Popup } from '$lib/components/popups/popup';
	import { load } from '../+page';
	import Button from '$lib/components/ui/Button.svelte';

	let email = '';
	let password = '';

	let emailError = '';
	let passwordError = '';
	let loading = false;

	function validateEmail(email: string) {
		const re = /^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,6}$/;
		return re.test(email);
	}

	async function handleSubmit() {
		emailError = !validateEmail(email) ? 'Please enter a valid email address' : '';
		passwordError = password.length < 8 ? 'Password must be at least 8 characters long' : '';

		if (!emailError && !passwordError) {
			loading = true;
			const resp = await fetch(PUBLIC_SERVER_URL + '/login', {
				method: 'post',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					email,
					password
				})
			});

			if (!resp.ok) {
				const errorData = await resp.json();
				Popup('Error', 'Unable to login: ' + errorData.error, 10000);
			} else {
				let res = await resp.json();
				token.set(res['token']);
				window.location.replace('/');
				Popup('Success', 'Logged in successfully');
			}
			loading = false;
		}
	}
</script>

<div class="mb-4">
	<h1 class="text-2xl font-semibold">Welcome to Suraksheet</h1>
	<p class="text-sm text-gray-500">Securely login to access your documents.</p>
</div>

<hr class="mb-4 border-slate-500" />

<div class="max-w-lg">
	<form on:submit|preventDefault={handleSubmit} class="flex flex-col space-y-6">
		<h3 class="text-xl font-semibold">Login to your account</h3>
		<Label class="space-y-2">
			<span>Email</span>
			<Input type="email" name="email" placeholder="name@company.com" required bind:value={email}>
				<EnvelopeSolid slot="left" class="h-5 w-5" />
			</Input>
			{#if emailError}
				<Helper color="red">{emailError}</Helper>
			{/if}
		</Label>

		<Label class="space-y-2">
			<span>Password</span>
			<Input type="password" name="password" placeholder="••••••••" required bind:value={password}>
				<LockSolid slot="left" class="h-5 w-5 text-gray-500 dark:text-gray-400" />
			</Input>
			{#if passwordError}
				<Helper color="red">{passwordError}</Helper>
			{/if}
		</Label>
		<Button bind:loading loadingText="Logging in..." type="submit">Login</Button>
	</form>

	<div class="mt-5 text-sm font-medium">
		Don't have an account? <a
			href="/register"
			class="text-primary-500 dark:text-primary-500 hover:underline">Register here</a
		>
	</div>
</div>
