<script lang="ts">
	import DocumentUpload from '$lib/components/DocumentUpload.svelte';

	function handleFilesAttached(event: CustomEvent) {
		console.log('Event received from DocumentUpload component:', event.detail);
		// Here you would typically handle the files, e.g., add them to a list in the parent component's state
		// or prepare them for actual upload to a server.
		alert(`Received ${event.detail.length} file(s) from the component. Check console.`);
	}
</script>

<svelte:head>
	<title>Test Document Upload - Procurement System</title>
</svelte:head>

<div class="container mx-auto py-8">
	<h1 class="text-2xl font-semibold mb-6">Test Document Upload Component</h1>

	<div class="max-w-xl mx-auto bg-gray-50 p-6 rounded-lg shadow-inner">
		<p class="mb-4 text-sm text-gray-600">
			This page is for testing the <code>DocumentUpload.svelte</code> component in isolation.
		</p>
		<DocumentUpload on:filesAttached={handleFilesAttached} />
	</div>

	<div class="mt-8 p-4 bg-gray-100 rounded-md">
		<h2 class="text-lg font-medium mb-2">Instructions:</h2>
		<ul class="list-disc list-inside text-sm text-gray-700">
			<li>Select one or more files using the 'Select files' input.</li>
			<li>Choose a 'Document type' from the dropdown.</li>
			<li>Click the 'Attach Selected Files' button.</li>
			<li>An alert should confirm the action, and details will be logged to the browser console.</li>
			<li>The component should then reset.</li>
		</ul>
	</div>
</div>
