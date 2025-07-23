<script>
  import { register } from '$lib/authService';
  import { authError, isLoading } from '$lib/store';
  import { onMount } from 'svelte';
  
  // Form data
  let email = '';
  let username = '';
  let password = '';
  let confirmPassword = '';
  let firstName = '';
  let lastName = '';
  
  // Form validation
  let errors = {
    email: '',
    username: '',
    password: '',
    confirmPassword: '',
    form: ''
  };
  
  // Reset form errors when input changes
  $: if (email) errors.email = '';
  $: if (username) errors.username = '';
  $: if (password) errors.password = '';
  $: if (confirmPassword) errors.confirmPassword = '';
  
  // Clear form error when any field changes
  $: if (email || username || password || confirmPassword) errors.form = '';
  
  // Handle form submission
  async function handleSubmit() {
    // Reset errors
    errors = {
      email: '',
      username: '',
      password: '',
      confirmPassword: '',
      form: ''
    };
    
    // Validate form
    let isValid = true;
    
    if (!email) {
      errors.email = 'Email is required';
      isValid = false;
    } else if (!/\S+@\S+\.\S+/.test(email)) {
      errors.email = 'Email is invalid';
      isValid = false;
    }
    
    if (!username) {
      errors.username = 'Username is required';
      isValid = false;
    }
    
    if (!password) {
      errors.password = 'Password is required';
      isValid = false;
    } else if (password.length < 8) {
      errors.password = 'Password must be at least 8 characters';
      isValid = false;
    }
    
    if (password !== confirmPassword) {
      errors.confirmPassword = 'Passwords do not match';
      isValid = false;
    }
    
    if (!isValid) return;
    
    try {
      await register(email, username, password, firstName, lastName);
    } catch (error) {
      errors.form = error.message || 'Registration failed';
    }
  }
  
  // Clear auth error when component mounts
  onMount(() => {
    authError.set(null);
  });
</script>

<div class="register-form">
  <h2>Create an Account</h2>
  
  {#if errors.form}
    <div class="error-message">{errors.form}</div>
  {/if}
  
  {#if $authError}
    <div class="error-message">{$authError.message}</div>
  {/if}
  
  <form on:submit|preventDefault={handleSubmit}>
    <div class="form-group">
      <label for="email">Email *</label>
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
      <label for="username">Username *</label>
      <input 
        type="text" 
        id="username" 
        bind:value={username} 
        disabled={$isLoading}
        required
      />
      {#if errors.username}
        <span class="field-error">{errors.username}</span>
      {/if}
    </div>
    
    <div class="form-group">
      <label for="firstName">First Name</label>
      <input 
        type="text" 
        id="firstName" 
        bind:value={firstName} 
        disabled={$isLoading}
      />
    </div>
    
    <div class="form-group">
      <label for="lastName">Last Name</label>
      <input 
        type="text" 
        id="lastName" 
        bind:value={lastName} 
        disabled={$isLoading}
      />
    </div>
    
    <div class="form-group">
      <label for="password">Password *</label>
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
    
    <div class="form-group">
      <label for="confirmPassword">Confirm Password *</label>
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
        {$isLoading ? 'Registering...' : 'Register'}
      </button>
    </div>
    
    <div class="form-footer">
      Already have an account? <a href="/login">Log in</a>
    </div>
  </form>
</div>

<style>
  .register-form {
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
    text-align: center;
    margin-top: 15px;
  }
  
  .form-actions {
    margin-top: 20px;
  }
</style>