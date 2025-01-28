<script lang="ts">
	import { Label, Input, Helper, Checkbox, Spinner } from 'flowbite-svelte';
	import { EnvelopeSolid, UserSolid, LockSolid } from 'flowbite-svelte-icons';
	import { PUBLIC_SERVER_URL } from '$env/static/public';
	import { Popup } from '$lib/components/popups/popup';
	import Button from '$lib/components/ui/Button.svelte';
	import { goto } from '$app/navigation';

	let email = '';
	let firstName = '';
	let lastName = '';
	let password = '';
	let confirmPassword = '';
	let agreeTerms = false;

	let emailError = '';
	let passwordError = '';
	let confirmPasswordError = '';

	let loading = false;

	function validateEmail(email: string) {
		const re = /^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,6}$/;
		return re.test(email);
	}

	async function handleSubmit() {
		emailError = !validateEmail(email) ? 'Please enter a valid email address' : '';
		passwordError = password.length < 8 ? 'Password must be at least 8 characters long' : '';
		confirmPasswordError = password !== confirmPassword ? 'Passwords do not match' : '';

		if (!emailError && !passwordError && !confirmPasswordError && agreeTerms) {
			loading = true;
			let resp = await fetch(PUBLIC_SERVER_URL + '/register', {
				method: 'post',
				body: JSON.stringify({
					email,
					firstName,
					lastName,
					password
				})
			});
			if (!resp.ok) {
				const errorData = await resp.json();
				Popup('Error', errorData.error, 10000);
			} else {
				Popup('Success', 'Account created successfully');
				goto('/login');
			}
			loading = false;
		}
	}
</script>

<div class="mb-4">
	<h1 class="text-2xl font-semibold">Welcome to Suraksheet</h1>
	<p class="text-sm text-gray-500">A safe place for your documents.</p>
</div>

<hr class="mb-4 border-slate-500" />

<div class="max-w-lg">
	<form on:submit|preventDefault={handleSubmit} class="flex flex-col space-y-6">
		<h3 class="text-xl font-semibold">Create an account</h3>

		<Label name="email" class="space-y-2">
			<span>Email</span>
			<Input type="email" name="email" placeholder="name@company.com" required bind:value={email}>
				<EnvelopeSolid slot="left" class="h-5 w-5" />
			</Input>
			{#if emailError}
				<Helper color="red">{emailError}</Helper>
			{/if}
		</Label>

		<div class="grid grid-cols-2 gap-4">
			<Label class="space-y-2">
				<span>First Name</span>
				<Input type="text" name="firstName" placeholder="John" required bind:value={firstName}>
					<UserSolid slot="left" class="h-5 w-5 text-gray-500 dark:text-gray-400" />
				</Input>
			</Label>

			<Label class="space-y-2">
				<span>Last Name</span>
				<Input type="text" name="lastName" placeholder="Doe" required bind:value={lastName}>
					<UserSolid slot="left" class="h-5 w-5 text-gray-500 dark:text-gray-400" />
				</Input>
			</Label>
		</div>

		<Label class="space-y-2">
			<span>Password</span>
			<Input type="password" name="password" placeholder="••••••••" required bind:value={password}>
				<LockSolid slot="left" class="h-5 w-5 text-gray-500 dark:text-gray-400" />
			</Input>
			{#if passwordError}
				<Helper color="red">{passwordError}</Helper>
			{/if}
		</Label>

		<Label class="space-y-2">
			<span>Confirm Password</span>
			<Input
				type="password"
				name="confirmPassword"
				placeholder="••••••••"
				required
				bind:value={confirmPassword}
			>
				<LockSolid slot="left" class="h-5 w-5 text-gray-500 dark:text-gray-400" />
			</Input>
			{#if confirmPasswordError}
				<Helper color="red">{confirmPasswordError}</Helper>
			{/if}
		</Label>

		<Checkbox required bind:checked={agreeTerms}>
			<div>
				I agree with the <a href="/terms" class="text-primary-500 hover:underline"
					>terms and conditions</a
				>
			</div>
		</Checkbox>

		<Button
			bind:loading
			loadingText="Creating account..."
			type="submit"
			class="bg-primary-500 dark:bg-primary-600 w-full">Create account</Button
		>
	</form>

	<div class="mt-5 text-sm font-medium">
		Already have an account? <a
			href="/login"
			class="text-primary-500 dark:text-primary-500 hover:underline">Login here</a
		>
	</div>
</div>
