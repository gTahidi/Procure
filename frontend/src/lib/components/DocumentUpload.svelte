<script lang="ts">
	import { createEventDispatcher } from 'svelte';

	let files: FileList | null = null;
	let selectedDocumentType: string = '';
	const dispatch = createEventDispatcher();

	// Mock document types - these would eventually come from a config or API
	const documentTypes = [
		{ value: '', label: 'Select document type...' },
		{ value: 'quote', label: 'Quotation' },
		{ value: 'specification', label: 'Specification Sheet' },
		{ value: 'contract_draft', label: 'Contract Draft' },
		{ value: 'supporting_document', label: 'Supporting Document' },
		{ value: 'other', label: 'Other' }
	];

	function handleFileSelect(event: Event) {
		const input = event.target as HTMLInputElement;
		if (input.files) {
			files = input.files;
		}
	}

	function handleAttachFiles() {
		if (!files || files.length === 0) {
			alert('Please select at least one file.');
			return;
		}
		if (!selectedDocumentType) {
			alert('Please select a document type.');
			return;
		}

		// In a real scenario, you'd start the upload process here.
		// For now, we'll dispatch an event with the file details.
		const fileDetails = Array.from(files).map(file => ({
			file: file,
			name: file.name,
			size: file.size,
			type: file.type,
			documentType: selectedDocumentType
		}));

		dispatch('filesAttached', fileDetails);
		
		// Reset after 'upload'
		files = null;
		selectedDocumentType = '';
		// Clear the file input visually (this is a bit tricky, might need a more robust solution or component reset from parent)
		const fileInput = document.getElementById('fileInput') as HTMLInputElement;
		if(fileInput) fileInput.value = '';

		alert(`${fileDetails.length} file(s) 'attached' with type: ${documentTypes.find(dt => dt.value === fileDetails[0].documentType)?.label}. Check console for details.`);
		console.log('Files Attached:', fileDetails);
	}

	function formatFileSize(bytes: number): string {
        if (bytes === 0) return '0 Bytes';
        const k = 1024;
        const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    }
</script>

<div class="border border-gray-300 p-4 rounded-lg shadow-sm bg-white">
	<h3 class="text-lg font-medium text-gray-800 mb-3">Upload Documents</h3>
	
	<div class="mb-4">
		<label for="fileInput" class="block text-sm font-medium text-gray-700 mb-1">Select files:</label>
		<input 
			type="file" 
			id="fileInput"
			multiple 
			on:change={handleFileSelect} 
			class="block w-full text-sm text-gray-900 border border-gray-300 rounded-lg cursor-pointer bg-gray-50 focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 p-2"
		/>
	</div>

	{#if files && files.length > 0}
		<div class="mb-4">
			<p class="text-sm font-medium text-gray-700 mb-1">Selected files:</p>
			<ul class="list-disc list-inside pl-2 text-sm text-gray-600 max-h-32 overflow-y-auto border border-gray-200 p-2 rounded-md">
				{#each Array.from(files) as file}
					<li>{file.name} ({formatFileSize(file.size)})</li>
				{/each}
			</ul>
		</div>
		
		<div class="mb-4">
			<label for="documentType" class="block text-sm font-medium text-gray-700 mb-1">Document type:</label>
			<select 
				id="documentType" 
				bind:value={selectedDocumentType} 
				class="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md"
				required
			>
				{#each documentTypes as docType}
					<option value={docType.value} disabled={docType.value === ''}>{docType.label}</option>
				{/each}
			</select>
		</div>

		<button 
			on:click={handleAttachFiles} 
			class="w-full bg-indigo-600 hover:bg-indigo-700 text-white font-semibold py-2 px-4 rounded-md focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
			disabled={!files || files.length === 0 || !selectedDocumentType}
		>
			Attach Selected Files
		</button>
	{:else}
		<p class="text-sm text-gray-500 italic">No files selected.</p>
	{/if}
</div>
