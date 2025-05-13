<svelte:head>
	<title>Login - Procurement System</title>
	<meta name="description" content="Login to the Procurement System" />
</svelte:head>

<div class="max-w-md mx-auto mt-10 bg-white p-8 rounded-lg shadow-md">
	<h1 class="text-2xl font-semibold text-center mb-6">Login</h1>
	<form on:submit|preventDefault={() => alert('Login form submitted (mock)')}>
		<div class="mb-4">
			<label for="email" class="block text-sm font-medium text-gray-700 mb-1">Email Address</label>
			<input type="email" id="email" name="email" placeholder="you@example.com" required 
			       class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm">
		</div>
		<div class="mb-6">
			<label for="password" class="block text-sm font-medium text-gray-700 mb-1">Password</label>
			<input type="password" id="password" name="password" placeholder="••••••••" required
			       class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm">
		</div>
		<button type="submit" 
		        class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
			Sign In
		</button>
	</form>
	<p class="text-center text-sm text-gray-600 mt-4">
		Authentication will be handled by Auth0.
	</p>
</div>
