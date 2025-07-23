<script>
  import { resetPassword } from '$lib/authService';
  import { isLoading } from '$lib/store';
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  
  // Form data
  let token = '';
  let newPassword = '';
  let confirmPassword = '';
  
  // Form state
  let errors = {
    token: '',
    newPassword: '',
    confirmPassword: '',
    form: ''
  };
  let successMessage = '';
  
  // Get token from URL query parameter
  onMount(() => {
    const urlToken = $page.url.searchParams.get('token');
    if (urlToken) {
      token = urlToken;
    } else {
      errors.token = 'Reset token is missing';
    }
  });
  
  // Reset form errors when input changes
  $: if (newPassword) errors.newPassword = '';
  $: if (confirmPassword) errors.confirmPassword = '';
  
  // Clear form error when any field changes
  $: if (newPassword || confirmPassword) {
    errors.form = '';
    successMessage = '';
  }
  
  // Handle form submission
  async function handleSubmit() {
    // Reset messages
    errors = {
      token: token ? '' : 'Reset token is missing',
      newPassword: '',
      confirmPassword: '',
      form: ''
    };
    successMessage = '';
    
    // Validate form
    let isValid = true;
    
    if (!token) {
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
      await resetPassword(token, newPassword);
      
      // Show success message
      successMessage = 'Password has been reset successfully';
      
      // Redirect to login page after a delay
      setTimeout(() => {
        goto('/login');
      }, 3000);
    } catch (error) {
      errors.form = error.message || 'Failed to reset password';
    }
  }
</script>

<div class="reset-password-form">
  <h2>Reset Your Password</h2>
  
  {#if errors.token}
    <div class="error-message">
      {errors.token}
      <p>Please request a new password reset link.</p>
      <a href="/forgot-password" class="link-button">Request New Link</a>
    </div>
  {:else}
    {#if errors.form}
      <div class="error-message">{errors.form}</div>
    {/if}
    
    {#if successMessage}
      <div class="success-message">
        {successMessage}
        <p>You will be redirected to the login page shortly.</p>
      </div>
    {/if}
    
    <form on:submit|preventDefault={handleSubmit}>
      <div class="form-group">
        <label for="newPassword">New Password</label>
        <input 
          type="password" 
          id="newPassword" 
          bind:value={newPassword} 
          disabled={$isLoading || !!errors.token || !!successMessage}
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
          disabled={$isLoading || !!errors.token || !!successMessage}
          required
        />
        {#if errors.confirmPassword}
          <span class="field-error">{errors.confirmPassword}</span>
        {/if}
      </div>
      
      <div class="form-actions">
        <button 
          type="submit" 
          disabled={$isLoading || !!errors.token || !!successMessage}
        >
          {$isLoading ? 'Resetting Password...' : 'Reset Password'}
        </button>
      </div>
      
      <div class="form-footer">
        <a href="/login">Back to Login</a>
      </div>
    </form>
  {/if}
</div>

<style>
  .reset-password-form {
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
  
  .form-footer {
    text-align: center;
    margin-top: 15px;
  }
  
  .form-actions {
    margin-top: 20px;
  }
  
  .link-button {
    display: inline-block;
    margin-top: 10px;
    padding: 8px 16px;
    background-color: #3498db;
    color: white;
    text-decoration: none;
    border-radius: 4px;
  }
  
  .link-button:hover {
    background-color: #2980b9;
  }
</style>