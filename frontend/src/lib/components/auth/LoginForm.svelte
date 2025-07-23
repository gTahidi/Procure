<script>
  import { login } from '$lib/authService';
  import { authError, isLoading } from '$lib/store';
  import { onMount } from 'svelte';
  
  // Form data
  let email = '';
  let password = '';
  
  // Form validation
  let errors = {
    email: '',
    password: '',
    form: ''
  };
  
  // Reset form errors when input changes
  $: if (email) errors.email = '';
  $: if (password) errors.password = '';
  
  // Clear form error when any field changes
  $: if (email || password) errors.form = '';
  
  // Handle form submission
  async function handleSubmit() {
    // Reset errors
    errors = {
      email: '',
      password: '',
      form: ''
    };
    
    // Validate form
    let isValid = true;
    
    if (!email) {
      errors.email = 'Email is required';
      isValid = false;
    }
    
    if (!password) {
      errors.password = 'Password is required';
      isValid = false;
    }
    
    if (!isValid) return;
    
    try {
      await login(email, password);
    } catch (error) {
      errors.form = error.message || 'Login failed';
    }
  }
  
  // Clear auth error when component mounts
  onMount(() => {
    authError.set(null);
  });
</script>

<div class="login-form">
  <h2>Log In</h2>
  
  {#if errors.form}
    <div class="error-message">{errors.form}</div>
  {/if}
  
  {#if $authError}
    <div class="error-message">{$authError.message}</div>
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
    
    <div class="form-group">
      <label for="password">Password</label>
      <input 
        type="password" 
        id="password" 
        bind:value={password} 
        disabled={$isLoading}
        required
      />
      {#if errors.password}
        <span class="field-error">{errors.password}</span>
      {/if}
    </div>
    
    <div class="form-actions">
      <button type="submit" disabled={$isLoading}>
        {$isLoading ? 'Logging in...' : 'Log In'}
      </button>
    </div>
    
    <div class="form-footer">
      <div class="forgot-password">
        <a href="/forgot-password">Forgot password?</a>
      </div>
      <div class="register-link">
        Don't have an account? <a href="/register">Register</a>
      </div>
    </div>
  </form>
</div>

<style>
  .login-form {
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
    margin-top: 15px;
    text-align: center;
  }
  
  .forgot-password {
    margin-bottom: 10px;
  }
  
  .form-actions {
    margin-top: 20px;
  }
</style>