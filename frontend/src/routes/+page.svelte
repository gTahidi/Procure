<script lang="ts">
  import { user } from '$lib/store';
  import { getAccessTokenSilently, login } from '$lib/authService';
  import { PUBLIC_VITE_API_BASE_URL} from '$env/static/public';

  // Interfaces for our data structures
  interface DashboardStats {
    pendingApproval: number;
    readyForTender: number;
    activeTenders: number;
    recentlyClosed: number;
  }

  interface RecentRequisition {
    id: number;
    title: string;
    status: string;
    creationDate: string;
  }

  interface LiveTender {
    id: number;
    title: string;
    category: string;
    closingDate: string;
  }

  // State variables
  let dashboardStats: DashboardStats | null = null;
  let recentRequisitions: RecentRequisition[] | null = null;
  let liveTenders: LiveTender[] | null = null;
  let isLoadingData = false;
  let error: string | null = null;

  // State for requester's dashboard
  interface MyRequisitionStats {
    pending: number;
    approved: number;
    rejected: number;
  }

  interface MyRecentRequisition {
    id: number;
    status: string;
    created_at: string;
    updated_at: string;
  }

  let myStats: MyRequisitionStats | null = null;
  let myRecentRequisitions: MyRecentRequisition[] | null = null;
  let isLoadingMyData = false;
  let myError: string | null = null;

  // State for supplier's dashboard
  interface SupplierDashboardData {
    bids_submitted: number;
    bids_awarded: number;
    active_tenders: any[]; // Using any[] for now, can be typed to Tender model
    my_bids: any[]; // Using any[] for now, can be typed to Bid model
  }

  let supplierData: SupplierDashboardData | null = null;
  let isLoadingSupplierData = false;
  let supplierError: string | null = null;

  // Consolidated data fetching function
  async function fetchAllDashboardData() {
    if (isLoadingData) return;
    isLoadingData = true;
    error = null;

    try {
      const token = await getAccessTokenSilently();
      if (!token) throw new Error('Could not retrieve access token.');

      const headers = { Authorization: `Bearer ${token}` };

      // Fetch all data in parallel for efficiency
      const [statsResponse, requisitionsResponse, tendersResponse] = await Promise.all([
        fetch(`${PUBLIC_VITE_API_BASE_URL}/api/dashboard/requisition-stats`, { headers }),
        fetch(`${PUBLIC_VITE_API_BASE_URL}/api/dashboard/recent-requisitions`, { headers }),
        fetch(`${PUBLIC_VITE_API_BASE_URL}/api/dashboard/live-tenders`, { headers }),
      ]);

      // Check all responses before proceeding
      if (!statsResponse.ok) throw new Error(`Failed to fetch stats: ${statsResponse.statusText}`);
      if (!requisitionsResponse.ok) throw new Error(`Failed to fetch recent requisitions: ${requisitionsResponse.statusText}`);
      if (!tendersResponse.ok) throw new Error(`Failed to fetch live tenders: ${tendersResponse.statusText}`);

      // Parse all responses
      [dashboardStats, recentRequisitions, liveTenders] = await Promise.all([
        statsResponse.json(),
        requisitionsResponse.json(),
        tendersResponse.json(),
      ]);
    } catch (e) {
      if (e instanceof Error) error = e.message;
      else error = 'An unknown error occurred.';
      console.error('Error fetching dashboard data:', e);
    } finally {
      isLoadingData = false;
    }
  }

  async function fetchRequesterDashboardData() {
    if (isLoadingMyData) return;
    isLoadingMyData = true;
    myError = null;

    try {
      const token = await getAccessTokenSilently();
      if (!token) throw new Error('Could not retrieve access token.');
      const headers = { Authorization: `Bearer ${token}` };

      const [statsRes, recentRes] = await Promise.all([
        fetch(`${PUBLIC_VITE_API_BASE_URL}/api/dashboard/my-stats`, { headers }),
        fetch(`${PUBLIC_VITE_API_BASE_URL}/api/dashboard/my-recent-requisitions`, { headers }),
      ]);

      if (!statsRes.ok) throw new Error(`Failed to fetch your stats: ${statsRes.statusText}`);
      if (!recentRes.ok) throw new Error(`Failed to fetch your recent requisitions: ${recentRes.statusText}`);

      myStats = await statsRes.json();
      myRecentRequisitions = await recentRes.json();
    } catch (e) {
      if (e instanceof Error) myError = e.message;
      else myError = 'An unknown error occurred while fetching your data.';
      console.error('Error fetching requester dashboard data:', e);
    } finally {
      isLoadingMyData = false;
    }
  }

  // Reactively fetch data based on user role
  $: {
    if ($user) {
      const role = $user.role;
      if (role === 'procurement_officer') {
        if (!isLoadingData && !dashboardStats && !error) {
          fetchAllDashboardData();
        }
      } else if (role === 'supplier') {
        if (!isLoadingSupplierData && !supplierData && !supplierError) {
          fetchSupplierDashboardData();
        }
      } else { // Assumes requester for any other role
        if (!isLoadingMyData && !myStats && !myError) {
          fetchRequesterDashboardData();
        }
      }
    }
  }

  // Helper to format dates for display
  async function fetchSupplierDashboardData() {
    if (isLoadingSupplierData) return;
    isLoadingSupplierData = true;
    supplierError = null;

    try {
      const token = await getAccessTokenSilently();
      if (!token) throw new Error('Could not retrieve access token.');
      const headers = { Authorization: `Bearer ${token}` };

      const res = await fetch(`${PUBLIC_VITE_API_BASE_URL}/api/dashboard/supplier`, { headers });
      if (!res.ok) throw new Error(`Failed to fetch supplier dashboard: ${res.statusText}`);

      supplierData = await res.json();
    } catch (e) {
      if (e instanceof Error) supplierError = e.message;
      else supplierError = 'An unknown error occurred while fetching supplier data.';
      console.error('Error fetching supplier dashboard data:', e);
    } finally {
      isLoadingSupplierData = false;
    }
  }

  function formatDate(dateString: string) {
    if (!dateString) return 'N/A';
    const options: Intl.DateTimeFormatOptions = { year: 'numeric', month: 'short', day: 'numeric' };
    return new Date(dateString).toLocaleDateString('en-US', options);
  }
</script>

<svelte:head>
  <title>Dashboard - Procurement System</title>
  <meta name="description" content="Main dashboard for the Procurement System" />
</svelte:head>

{#if $user}
  {#if $user?.role === 'procurement_officer'}
    <!-- Procurement Officer Dashboard -->
    <div class="max-w-7xl mx-auto py-8 px-4 sm:px-6 lg:px-8">
            <div class="mb-12">
        <h1 class="text-3xl font-bold text-gray-900 mb-2">
          Welcome, {$user?.name || $user?.nickname || $user?.email || 'User'}!
        </h1>
        <p class="text-lg text-gray-600">Here is an overview of all procurement activities.</p>
      </div>

    {#if isLoadingData && !dashboardStats}
      <p class="text-center text-gray-500">Loading dashboard...</p>
    {:else if error}
      <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
        <strong class="font-bold">Error:</strong>
        <span class="block sm:inline">{error}</span>
      </div>
    {:else if dashboardStats}
      <!-- Stats Cards -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-12">
        <div class="bg-white overflow-hidden shadow rounded-lg p-5">
          <dt class="text-sm font-medium text-gray-500 truncate">Pending Approval</dt>
          <dd class="mt-1 text-3xl font-semibold text-gray-900">{dashboardStats.pendingApproval}</dd>
        </div>
        <div class="bg-white overflow-hidden shadow rounded-lg p-5">
          <dt class="text-sm font-medium text-gray-500 truncate">Ready for Tender</dt>
          <dd class="mt-1 text-3xl font-semibold text-gray-900">{dashboardStats.readyForTender}</dd>
        </div>
        <div class="bg-white overflow-hidden shadow rounded-lg p-5">
          <dt class="text-sm font-medium text-gray-500 truncate">Active Tenders</dt>
          <dd class="mt-1 text-3xl font-semibold text-gray-900">{dashboardStats.activeTenders}</dd>
        </div>
        <div class="bg-white overflow-hidden shadow rounded-lg p-5">
          <dt class="text-sm font-medium text-gray-500 truncate">Recently Closed</dt>
          <dd class="mt-1 text-3xl font-semibold text-gray-900">{dashboardStats.recentlyClosed}</dd>
        </div>
      </div>

      <!-- Tables Section -->
      <div class="space-y-12">
        <!-- Recent Requisitions Table -->
        <div>
          <h2 class="text-2xl font-bold text-gray-800 mb-4">Recent Requisitions</h2>
          {#if recentRequisitions && recentRequisitions.length > 0}
            <div class="bg-white shadow overflow-hidden rounded-lg">
              <table class="min-w-full divide-y divide-gray-200">
                <thead class="bg-gray-50">
                  <tr>
                    <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">ID</th>
                    <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Title</th>
                    <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
                    <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Created On</th>
                  </tr>
                </thead>
                <tbody class="bg-white divide-y divide-gray-200">
                  {#each recentRequisitions as req}
                    <tr>
                      <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{req.id}</td>
                      <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{req.title}</td>
                      <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{req.status}</td>
                      <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{formatDate(req.creationDate)}</td>
                    </tr>
                  {/each}
                </tbody>
              </table>
            </div>
          {:else}
            <p class="text-center text-gray-500 py-4">No recent requisitions found.</p>
          {/if}
        </div>

        <!-- Live Tenders Table -->
        <div>
          <h2 class="text-2xl font-bold text-gray-800 mb-4">Live Tenders</h2>
          {#if liveTenders && liveTenders.length > 0}
            <div class="bg-white shadow overflow-hidden rounded-lg">
              <table class="min-w-full divide-y divide-gray-200">
                <thead class="bg-gray-50">
                  <tr>
                    <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">ID</th>
                    <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Title</th>
                    <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Category</th>
                    <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Closing On</th>
                  </tr>
                </thead>
                <tbody class="bg-white divide-y divide-gray-200">
                  {#each liveTenders as tender}
                    <tr>
                      <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{tender.id}</td>
                      <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{tender.title}</td>
                      <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{tender.category}</td>
                      <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{formatDate(tender.closingDate)}</td>
                    </tr>
                  {/each}
                </tbody>
              </table>
            </div>
          {:else}
            <p class="text-center text-gray-500 py-4">No live tenders found.</p>
          {/if}
        </div>
      </div>
    {/if}
  </div>
{:else if $user?.role === 'supplier'}
  <!-- Supplier Dashboard -->
  <div class="max-w-7xl mx-auto py-8 px-4 sm:px-6 lg:px-8">
    <div class="mb-12">
      <h1 class="text-3xl font-bold text-gray-900 mb-2">
        Supplier Dashboard
      </h1>
      <p class="text-lg text-gray-600">Welcome, {$user?.name || $user?.nickname || $user?.email || 'Supplier'}!</p>
    </div>

    {#if isLoadingSupplierData && !supplierData}
      <p class="text-center text-gray-500">Loading supplier dashboard...</p>
    {:else if supplierError}
      <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
        <strong class="font-bold">Error:</strong>
        <span class="block sm:inline">{supplierError}</span>
      </div>
    {:else if supplierData}
      <!-- Supplier Stats -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-12">
        <div class="bg-white overflow-hidden shadow rounded-lg p-5">
          <dt class="text-sm font-medium text-gray-500 truncate">Active Tenders</dt>
          <dd class="mt-1 text-3xl font-semibold text-gray-900">{supplierData.active_tenders.length}</dd>
        </div>
        <div class="bg-white overflow-hidden shadow rounded-lg p-5">
          <dt class="text-sm font-medium text-gray-500 truncate">Bids Submitted</dt>
          <dd class="mt-1 text-3xl font-semibold text-gray-900">{supplierData.bids_submitted}</dd>
        </div>
        <div class="bg-white overflow-hidden shadow rounded-lg p-5">
          <dt class="text-sm font-medium text-gray-500 truncate">Bids Awarded</dt>
          <dd class="mt-1 text-3xl font-semibold text-gray-900">{supplierData.bids_awarded}</dd>
        </div>
      </div>

      <!-- Active Tenders Table -->
      <div class="mb-12">
        <h2 class="text-2xl font-bold text-gray-800 mb-4">Active Tenders</h2>
        <div class="bg-white shadow overflow-hidden rounded-lg">
          <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
              <tr>
                <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Tender ID</th>
                <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Title</th>
                <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Closing Date</th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              {#each supplierData.active_tenders as tender}
                <tr>
                  <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-blue-600 hover:underline"><a href={`/tenders/${tender.id}`}>TEN-{tender.id}</a></td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{tender.title}</td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{formatDate(tender.closing_date)}</td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      </div>

      <!-- My Bids Table -->
      <div>
        <h2 class="text-2xl font-bold text-gray-800 mb-4">My Submitted Bids</h2>
        <div class="bg-white shadow overflow-hidden rounded-lg">
          <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
              <tr>
                <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Bid ID</th>
                <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Tender Title</th>
                <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Amount</th>
                <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
                <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Submitted On</th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              {#each supplierData.my_bids as bid}
                <tr>
                  <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-blue-600 hover:underline"><a href={`/bids/${bid.id}`}>BID-{bid.id}</a></td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{bid.tender?.title || 'N/A'}</td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">${bid.bid_amount.toLocaleString()}</td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{bid.status}</td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{formatDate(bid.submission_date)}</td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      </div>
    {/if}
  </div>
{:else}
  <!-- Default/Requester Dashboard -->
  <div class="max-w-7xl mx-auto py-8 px-4 sm:px-6 lg:px-8">
    <div class="mb-12">
      <h1 class="text-3xl font-bold text-gray-900 mb-2">
        Welcome, {$user?.name || $user?.nickname || $user?.email || 'User'}!
      </h1>
      <p class="text-lg text-gray-600">Here's what's happening with your procurement activities today.</p>
    </div>

    {#if isLoadingMyData && !myStats}
      <p class="text-center text-gray-500">Loading your dashboard...</p>
    {:else if myError}
      <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
        <strong class="font-bold">Error:</strong>
        <span class="block sm:inline">{myError}</span>
      </div>
    {:else if myStats}
      <!-- Requester Stats Cards -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-12">
        <div class="bg-white overflow-hidden shadow rounded-lg p-5">
          <dt class="text-sm font-medium text-gray-500 truncate">Pending My Approval</dt>
          <dd class="mt-1 text-3xl font-semibold text-gray-900">{myStats.pending}</dd>
        </div>
        <div class="bg-white overflow-hidden shadow rounded-lg p-5">
          <dt class="text-sm font-medium text-gray-500 truncate">Approved Requisitions</dt>
          <dd class="mt-1 text-3xl font-semibold text-gray-900">{myStats.approved}</dd>
        </div>
        <div class="bg-white overflow-hidden shadow rounded-lg p-5">
          <dt class="text-sm font-medium text-gray-500 truncate">Rejected Requisitions</dt>
          <dd class="mt-1 text-3xl font-semibold text-gray-900">{myStats.rejected}</dd>
        </div>
        <div class="bg-white overflow-hidden shadow rounded-lg p-5 flex flex-col justify-center items-center">
			<a href="/requisitions/new" class="w-full text-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
				Create New Requisition
			</a>
        </div>
      </div>

      <!-- My Recent Requisitions Table -->
      <div>
        <h2 class="text-2xl font-bold text-gray-800 mb-4">My Recent Requisitions</h2>
        {#if myRecentRequisitions && myRecentRequisitions.length > 0}
          <div class="bg-white shadow overflow-hidden rounded-lg">
            <table class="min-w-full divide-y divide-gray-200">
              <thead class="bg-gray-50">
                <tr>
                  <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">ID</th>
                  <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
                  <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Created On</th>
                  <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Last Updated</th>
                </tr>
              </thead>
              <tbody class="bg-white divide-y divide-gray-200">
                {#each myRecentRequisitions as req}
                  <tr>
                    <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-blue-600 hover:underline">
						<a href={`/requisitions/${req.id}`}>REQ-{req.id}</a>
					</td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{req.status}</td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{formatDate(req.created_at)}</td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{formatDate(req.updated_at)}</td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        {:else}
          <div class="text-center border-2 border-dashed border-gray-300 rounded-lg p-12">
            <h3 class="text-lg font-medium text-gray-900">No requisitions yet!</h3>
            <p class="mt-2 text-sm text-gray-500">Get started by creating your first requisition.</p>
          </div>
        {/if}
      </div>
    {/if}
  </div>
  {/if}
{:else}
  <!-- Login Prompt for unauthenticated users -->
  <div class="text-center py-16 max-w-2xl mx-auto">
    <h1 class="text-3xl font-bold text-gray-800 mb-4">Welcome to the Procurement System</h1>
    <p class="text-lg text-gray-600 mb-8">
      Your central hub for managing requisitions, tenders, and bids.
    </p>
    <p class="text-gray-500 mb-8">Please log in to access your dashboard and manage your procurement activities.</p>
    <button
      on:click={() => login()}
      class="bg-blue-600 hover:bg-blue-700 text-white font-bold py-3 px-8 rounded-lg text-lg transition duration-300 ease-in-out transform hover:scale-105"
    >
      Log In
    </button>
  </div>
{/if}
