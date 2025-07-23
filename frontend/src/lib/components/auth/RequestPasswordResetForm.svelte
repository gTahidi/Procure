<script>
  import { requestPasswordReset } from '$lib/authService';
  import { isLoading } from '$lib/store';
  
  // Form data
  let email = '';
  
  // Form state
  let errors = {
    email: '',
    form: ''
  };
  let successMessage = '';
  
  // Reset form errors when input changes
  $: if (email) errors.email = '';
  
  // Clear form error when email changes
  $: if (email) {
    errors.form = '';
    successMessage = '';
  }
  
  // Handle form submission
  async function handleSubmit() {
    // Reset messages
    errors = {
      email: '',
      form: ''
    };
    successMessage = '';
    
    // Validate form
    let isValid = true;
    
    if (!email) {
      errors.email = 'Email is required';
      isValid = false;
    } else if (!/\S+@\S+\.\S+/.test(email)) {
      errors.email = 'Email is invalid';
      isValid = false;
    }
    
    if (!isValid) return;
    
    try {
      await requestPasswordReset(email);
      
      // Clear form on success
      email = '';
      
      // Show success message
      successMessage = 'Password reset instructions have been sent to your email';
    } catch (error) {
      // Don't reveal if email exists or not for security reasons
      // Just show a generic success message
      successMessage = 'If your email is registered, you will receive password reset instructions';
    }
  }
</script>

<div class="reset-password-form">
  <h2>Reset Password</h2>
  
  <p class="form-description">
    Enter your email address and we'll send you instructions to reset your password.
  </p>
  
  {#if errors.form}
    <div class="error-message">{errors.form}</div>
  {/if}
  
  {#if successMessage}
    <div class="success-message">{successMessage}</div>
  {/if}
  
  <form on:submit|preventDefault={handleSubmit}>
    <div class="form-group">
      <label for="email">Email</label>
      <input 
        type="email" 
        id="email" 
        bind:value={email} 
        disabled={$isLoading}
        required
      />
      {#if errors.email}
        <span class="field-error">{errors.email}</span>
      {/if}
    </div>
    
    <div class="form-actions">
      <button type="submit" disabled={$isLoading}>
        {$isLoading ? 'Sending...' : 'Send Reset Instructions'}
      </button>
    </div>
    
    <div class="form-footer">
      <a href="/login">Back to Login</a>
    </div>
  </form>
</div>

<style>
  .reset-password-form {
    max-width: 400px;
    margin: 0 auto;
    padding: 20px;
  }
  
  h2 {
    text-align: center;
    margin-bottom: 10px;
  }
  
  .form-description {
    text-align: center;
    margin-bottom: 20px;
    color: #666;
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
</style>