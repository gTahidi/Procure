<script lang="ts">
  import { onMount } from 'svelte';
  import { initializeAuth, login, logout } from '$lib/authService';
  import { isAuthenticated, user, isLoading, authError } from '$lib/store';
  import { syncUserToDb } from '$lib/userService'; 
  import { page } from '$app/stores'; // Import $page store
  import { fly } from 'svelte/transition'; // Import fly transition
  import { sineInOut } from 'svelte/easing'; // Import sineInOut easing
  import '../app.css'; 

  let initialAuthDone = false;
  let sidebarOpen = true;

  function toggleSidebar() {
    sidebarOpen = !sidebarOpen;
  }

  const navLinks = [
    { href: '/', label: 'Home', icon: '[H]' }, 
    { href: '/requisitions', label: 'Requisitions', icon: '[R]' }, 
    { href: '/tenders', label: 'Tenders', icon: '[T]' }, 
    { href: '/documents', label: 'Documents', icon: '[D]' } 
  ];

  $: activeLink = (href: string) => {
    return $page.url.pathname === href || ($page.url.pathname.startsWith(href) && href !== '/');
  };

  onMount(() => {
    let unsubscribeUser: (() => void) | undefined;

    async function setupAuthAndSubscribe() {
      await initializeAuth();
      initialAuthDone = true;

      unsubscribeUser = user.subscribe(currentUser => {
        if (currentUser && $isAuthenticated && initialAuthDone) {
          console.log('User authenticated, attempting to sync to DB:', currentUser);
          syncUserToDb(currentUser); 
        }
      });
    }

    setupAuthAndSubscribe();

    return () => {
      if (unsubscribeUser) {
        unsubscribeUser();
      }
    };
  });

  function handleLogin() {
    login();
  }

  function handleLogout() {
    logout();
  }
</script>

{#if $isLoading && !initialAuthDone}
  <div>Loading authentication...</div>
{:else}
  <div class="flex h-screen bg-gray-100">
    <!-- Sidebar -->
    <aside 
      class="bg-gray-800 text-white flex flex-col shadow-xl transition-all duration-300 ease-in-out"
      class:w-64={sidebarOpen}
      class:w-20={!sidebarOpen}
    >
      <div class="flex items-center justify-between p-4 border-b border-gray-700 h-16">
        {#if sidebarOpen}
          <a href="/" class="text-xl font-semibold whitespace-nowrap" transition:fly={{ x: -20, duration: 300, easing: sineInOut }}>Procurement System</a>
        {/if}
        <button on:click={toggleSidebar} class="p-2 rounded-md hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-white">
          {#if sidebarOpen}
            <span>&lt;</span> 
          {:else}
            <span>M</span> 
          {/if}
        </button>
      </div>

      <nav class="flex-1 mt-4 space-y-2 overflow-y-auto overflow-x-hidden">
        {#each navLinks as link}
          <a 
            href={link.href} 
            class="flex items-center py-2.5 rounded-lg hover:bg-gray-700 transition-colors duration-150"
            class:bg-gray-900={activeLink(link.href)}
            class:px-4={sidebarOpen}
            class:px-0={!sidebarOpen}
            class:justify-start={sidebarOpen}
            class:justify-center={!sidebarOpen}
            title={link.label}
          >
            <span class:ml-0={!sidebarOpen} class:mr-0={!sidebarOpen} class:mx-auto={!sidebarOpen} class:ml-4={sidebarOpen}>
              <span class="inline-block w-6 h-6 text-center">{link.icon}</span> 
            </span>
            {#if sidebarOpen}
              <span class="ml-3 whitespace-nowrap" transition:fly={{ x: -20, duration: 200, delay: 50, easing: sineInOut }}>{link.label}</span>
            {/if}
          </a>
        {/each}
      </nav>

      <div class="border-t border-gray-700 p-0" class:p-4={sidebarOpen}>
        {#if $isAuthenticated}
          <div class="user-info my-2" class:text-center={!sidebarOpen} class:px-4={sidebarOpen}>
            {#if sidebarOpen}
              <span class="block text-sm">{$user?.name || $user?.email}</span>
            {/if}
          </div>
          <button 
            on:click={handleLogout}
            class="flex items-center w-full py-2.5 rounded-lg hover:bg-gray-700 transition-colors duration-150"
            class:px-4={sidebarOpen}
            class:px-0={!sidebarOpen}
            class:justify-start={sidebarOpen}
            class:justify-center={!sidebarOpen}
            title="Logout"
          >
            <span class:ml-0={!sidebarOpen} class:mr-0={!sidebarOpen} class:mx-auto={!sidebarOpen} class:ml-4={sidebarOpen}>
              <span class="inline-block w-6 h-6 text-center">[LO]</span> <!-- Placeholder for LogoutIcon -->
            </span>
            {#if sidebarOpen}
              <span class="ml-3 whitespace-nowrap" transition:fly={{ x: -20, duration: 200, delay: 50, easing: sineInOut }}>Logout</span>
            {/if}
          </button>
        {:else}
          <button 
            on:click={handleLogin}
            class="flex items-center w-full py-2.5 rounded-lg hover:bg-gray-700 transition-colors duration-150"
            class:px-4={sidebarOpen}
            class:px-0={!sidebarOpen}
            class:justify-start={sidebarOpen}
            class:justify-center={!sidebarOpen}
            title="Login"
          >
            <span class:ml-0={!sidebarOpen} class:mr-0={!sidebarOpen} class:mx-auto={!sidebarOpen} class:ml-4={sidebarOpen}>
              <span class="inline-block w-6 h-6 text-center">[LI]</span> <!-- Placeholder for LoginIcon -->
            </span>
            {#if sidebarOpen}
              <span class="ml-3 whitespace-nowrap" transition:fly={{ x: -20, duration: 200, delay: 50, easing: sineInOut }}>Login</span>
            {/if}
          </button>
        {/if}
        {#if $authError}
          <p style="color: red; font-size: 0.8em;" class:px-4={sidebarOpen} class:text-center={!sidebarOpen} >Error: {$authError.message}</p>
        {/if}
      </div>
    </aside>

    <!-- Main content area -->
    <div class="flex-1 flex flex-col overflow-hidden">
      <main class="flex-1 overflow-x-hidden overflow-y-auto bg-gray-100 p-6">
        <slot />
      </main>

      <!-- Footer (optional, can be removed or kept simple) -->
      <footer class="bg-white shadow-inner p-4 text-center text-sm text-gray-600">
        &copy; {new Date().getFullYear()} Procurement System. All rights reserved.
      </footer>
    </div>
  </div>
{/if}
