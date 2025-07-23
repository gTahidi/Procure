<script>
  import { changePassword } from '$lib/authService';
  import { isLoading } from '$lib/store';
  
  // Form data
  let currentPassword = '';
  let newPassword = '';
  let confirmPassword = '';
  
  // Form state
  let errors = {
    currentPassword: '',
    newPassword: '',
    confirmPassword: '',
    form: ''
  };
  let successMessage = '';
  
  // Reset form errors when input changes
  $: if (currentPassword) errors.currentPassword = '';
  $: if (newPassword) errors.newPassword = '';
  $: if (confirmPassword) errors.confirmPassword = '';
  
  // Clear form error when any field changes
  $: if (currentPassword || newPassword || confirmPassword) {
    errors.form = '';
    successMessage = '';
  }
  
  // Handle form submission
  async function handleSubmit() {
    // Reset messages
    errors = {
      currentPassword: '',
      newPassword: '',
      confirmPassword: '',
      form: ''
    };
    successMessage = '';
    
    // Validate form
    let isValid = true;
    
    if (!currentPassword) {
      errors.currentPassword = 'Current password is required';
      isValid = false;
    }
    
    if (!newPassword) {
      errors.newPassword = 'New password is required';
      isValid = false;
    } else if (newPassword.length < 8) {
      errors.newPassword = 'Password must be at least 8 characters';
      isValid = false;
    }
    
    if (newPassword !== confirmPassword) {
      errors.confirmPassword = 'Passwords do not match';
      isValid = false;
    }
    
    if (!isValid) return;
    
    try {
      await changePassword(currentPassword, newPassword);
      
      // Clear form on success
      currentPassword = '';
      newPassword = '';
      confirmPassword = '';
      
      // Show success message
      successMessage = 'Password changed successfully';
    } catch (error) {
      errors.form = error.message || 'Failed to change password';
    }
  }
</script>

<div class="change-password-form">
  <h2>Change Password</h2>
  
  {#if errors.form}
    <div class="error-message">{errors.form}</div>
  {/if}
  
  {#if successMessage}
    <div class="success-message">{successMessage}</div>
  {/if}
  
  <form on:submit|preventDefault={handleSubmit}>
    <div class="form-group">
      <label for="currentPassword">Current Password</label>
      <input 
        type="password" 
        id="currentPassword" 
        bind:value={currentPassword} 
        disabled={$isLoading}
        required
      />
      {#if errors.currentPassword}
        <span class="field-error">{errors.currentPassword}</span>
      {/if}
    </div>
    
    <div class="form-group">
      <label for="newPassword">New Password</label>
      <input 
        type="password" 
        id="newPassword" 
        bind:value={newPassword} 
        disabled={$isLoading}
        required
      />
      {#if errors.newPassword}
        <span class="field-error">{errors.newPassword}</span>
      {/if}
    </div>
    
    <div class="form-group">
      <label for="confirmPassword">Confirm New Password</label>
      <input 
        type="password" 
        id="confirmPassword" 
        bind:value={confirmPassword} 
        disabled={$isLoading}
        required
      />
      {#if errors.confirmPassword}
        <span class="field-error">{errors.confirmPassword}</span>
      {/if}
    </div>
    
    <div class="form-actions">
      <button type="submit" disabled={$isLoading}>
        {$isLoading ? 'Changing Password...' : 'Change Password'}
      </button>
    </div>
  </form>
</div>

<style>
  .change-password-form {
    max-width: 400px;
    margin: 0 auto;
    padding: 20px;
  }
  
  h2 {
    text-align: center;
    margin-bottom: 20px;
  }
  
  .form-group {
    margin-bottom: 15px;
  }
  
  label {
    display: block;
    margin-bottom: 5px;
    font-weight: 500;
  }
  
  input {
    width: 100%;
    padding: 8px;
    border: 1px solid #ccc;
    border-radius: 4px;
  }
  
  .field-error {
    color: #e74c3c;
    font-size: 0.8rem;
    margin-top: 5px;
    display: block;
  }
  
  .error-message {
    background-color: #fdecea;
    color: #e74c3c;
    padding: 10px;
    border-radius: 4px;
    margin-bottom: 15px;
    text-align: center;
  }
  
  .success-message {
    background-color: #d4edda;
    color: #155724;
    padding: 10px;
    border-radius: 4px;
    margin-bottom: 15px;
    text-align: center;
  }
  
  button {
    width: 100%;
    padding: 10px;
    background-color: #3498db;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-weight: 500;
  }
  
  button:hover {
    background-color: #2980b9;
  }
  
  button:disabled {
    background-color: #95a5a6;
    cursor: not-allowed;
  }
  
  .form-actions {
    margin-top: 20px;
  }
</style>