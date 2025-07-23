<script lang="ts">
  import { onMount } from 'svelte';
  import { initializeAuth, logout } from '$lib/authService';
  import { user, isAuthenticated, isLoading, authError } from '$lib/store';
  import { syncUserToDb } from '$lib/userService'; 
  import type { AppUser } from '$lib/store';
  import { page } from '$app/stores';
  import { fly } from 'svelte/transition';
  import { sineInOut } from 'svelte/easing';
  import { goto } from '$app/navigation';
  import '../app.css'; 

  let initialAuthDone = false;
  let sidebarOpen = true;
  let currentUser: AppUser | null = null;

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

  onMount(async () => {
    try {
      isLoading.set(true);
      
      // Call initializeAuth to handle session checking and user fetching.
      // initializeAuth will update isAuthenticated, user, isLoading, and authError stores.
      await initializeAuth();

      // After initializeAuth, check if user is authenticated and sync to DB if needed.
      // Note: The 'currentUser' variable here might not yet have the role.
      // The 'user' store ($user) will be updated by syncUserToDb with the role.
      const unsubscribeUser = user.subscribe((value: AppUser | null) => { currentUser = value; });
      let currentAuthStatus: boolean = false;
      const unsubscribeAuth = isAuthenticated.subscribe((value: boolean) => { currentAuthStatus = value; });

      if (currentAuthStatus && currentUser) { // currentUser here is the basic user profile
        try {
          console.log('Layout: User is authenticated, attempting to sync to DB...');
          await syncUserToDb(); // Call syncUserToDb without arguments
          console.log('Layout: syncUserToDb completed. User store should now have role:', $user);
        } catch (syncError) {
          console.error('Layout: Failed to sync user to DB on layout load:', syncError);
        }
      }
      unsubscribeUser();
      unsubscribeAuth();
    } catch (e) {
      if (e instanceof Error) {
        authError.set(e);
      } else {
        authError.set(new Error('An unknown error occurred during auth initialization.'));
      }
      console.error('Auth initialization error:', e);
      isAuthenticated.set(false);
      user.set(null);
    } finally {
      isLoading.set(false);
      initialAuthDone = true;
    }
  });
  
  function handleLogin() {
    // Redirect to login page
    goto('/login');
  }

  async function handleLogout() {
    try {
      await logout();
      // Redirect to home after logout
      goto('/');
    } catch (error) {
      console.error('Logout error:', error);
    }
  }
</script>

{#if $isLoading && !initialAuthDone}
  <div>Loading authentication...</div>
{:else}
  <div class="flex h-screen bg-gray-100">
    {#if $page.url.pathname !== '/login'}
      <!-- Sidebar: Navigation links, User info, Logout button -->
      <aside class:w-64={sidebarOpen} class:w-20={!sidebarOpen} class="bg-gray-800 text-white flex flex-col transition-all duration-300 ease-in-out fixed top-0 left-0 h-full z-50">
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
            {#if link.href === '/tenders' || link.href === '/documents'}
              {#if $user?.role !== 'requester'}
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
              {/if}
            {:else} <!-- Regular links like Home and Requisitions -->
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
            {/if}
          {/each}
        </nav>

        <div class="border-t border-gray-700 p-0" class:p-4={sidebarOpen}>
          {#if $isAuthenticated && $user}
            {console.log('Sidebar user object:', $user)} 
            <div class="user-info my-2 flex flex-col items-center" class:items-start={sidebarOpen} class:px-4={sidebarOpen}>
              {#if $user.picture_url}
                <img src={$user.picture_url} alt="{$user && ($user.username || 'User')}'s profile picture" class="rounded-full mb-2" class:w-10={!sidebarOpen} class:h-10={!sidebarOpen} class:w-12={sidebarOpen} class:h-12={sidebarOpen} referrerpolicy="no-referrer" />
              {:else}
                <div class="rounded-full bg-gray-600 flex items-center justify-center text-white mb-2" class:w-10={!sidebarOpen} class:h-10={!sidebarOpen} class:w-12={sidebarOpen} class:h-12={sidebarOpen} class:text-lg={!sidebarOpen} class:text-xl={sidebarOpen}>
                  {($user && ($user.username || $user.email || 'U')).charAt(0).toUpperCase()}
                </div>
              {/if}
              {#if sidebarOpen}
                <span class="block text-sm font-semibold">{$user && ($user.username || $user.email)}</span>
                {#if $user && $user.role}
                  <span class="block text-xs text-gray-400 capitalize">{$user.role}</span>
                {/if}
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
    {/if}
    <!-- Main content area -->
    <div class:ml-64={sidebarOpen && $page.url.pathname !== '/login'} class:ml-20={!sidebarOpen && $page.url.pathname !== '/login'} class:ml-0={$page.url.pathname === '/login'} class="flex-1 flex flex-col overflow-hidden transition-all duration-300 ease-in-out">
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
