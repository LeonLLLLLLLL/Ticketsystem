<script lang="ts">
  // 🔗 All your logic here — unchanged, from your original code
  import { onMount } from 'svelte';
  import { browser } from '$app/environment';

  const BASE_URL = 'http://localhost:8000';

  let devices: any[] = [];
  let deviceLinks: any[] = [];

  let error: string | null = null;
  let success: string | null = null;

  let showForm = false;
  let mode: 'create' | 'edit' = 'create';
  let selectedId: number | null = null;
  let searchQuery = '';

  let selectedFromId: number | null = null;
  let selectedToId: number | null = null;

  let showDevicesSection = false;
  let showLinksSection = false;

  function resetDevice() {
    return {
      name: '',
      hostname: '',
      ip: '',
      domain: '',
      manufacturer: '',
      model_type: '',
      serial_numbers: '',
      mac: '',
      description: '',
      equipment: '',
      function: '',
      settings: '',
      device_link: '',
      commissioning_date: '',
      origin: '',
      warranty_service_number: '',
      warranty_until: '',
      licenses: '',
      location_text: '',
      department: '',
      internal_contact: '',
      external_contact: '',
      map_link: '',
      software_interfaces: '',
      backup_method: '',
      backup_file_link: '',
      software_asset: '',
      password_link: '',
      internal_access: '',
      external_access: '',
      misc_links: '',
      externally_accessible: false,
      restart_how: '',
      restart_notes: '',
      restart_coordination: '',
      network_connection: '',
      patch_location: '',
      documents: '',
    };
  }

  let device = resetDevice();

  function formatDateToISOString(dateStr: string) {
    const date = new Date(dateStr);
    if (isNaN(date.getTime())) return null;
    return date.toISOString();
  }

  async function fetchDevices() {
    try {
      const res = await fetch(`${BASE_URL}/devices/list`);
      if (!res.ok) throw new Error('Failed to fetch devices');
      devices = await res.json();
    } catch (err) {
      error = err.message;
    }
  }

  async function fetchLinks() {
    try {
      const res = await fetch(`${BASE_URL}/device_links/list`);
      if (!res.ok) throw new Error('Failed to fetch device links');
      deviceLinks = await res.json();
    } catch (err) {
      error = err.message;
    }
  }

  async function createDevice() {
    try {
      const payload = {
        ...device,
        commissioning_date: formatDateToISOString(device.commissioning_date),
        warranty_until: formatDateToISOString(device.warranty_until),
      };

      const res = await fetch(`${BASE_URL}/devices/create`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
      });

      if (!res.ok) {
        const err = await res.json();
        throw new Error(err.message || 'Error creating device');
      }

      const newDevice = await res.json();
      devices = [newDevice, ...devices];
      success = 'Device created successfully ✅';
      error = null;
      device = resetDevice();
      showForm = false;
    } catch (err) {
      success = null;
      error = err.message;
    }
  }

  async function updateDevice() {
    try {
      const payload = {
        ...device,
        id: selectedId,
        commissioning_date: formatDateToISOString(device.commissioning_date),
        warranty_until: formatDateToISOString(device.warranty_until),
      };

      const res = await fetch(`${BASE_URL}/devices/update`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
      });

      if (!res.ok) {
        const err = await res.json();
        throw new Error(err.message || 'Error updating device');
      }

      success = 'Device updated successfully ✅';
      error = null;
      showForm = false;
      device = resetDevice();
      selectedId = null;
      mode = 'create';
      await fetchDevices();
    } catch (err) {
      success = null;
      error = err.message;
    }
  }

  async function deleteDevice(id: number) {
    if (!confirm("Are you sure you want to delete this device?")) return;

    try {
      const res = await fetch(`${BASE_URL}/devices/delete?id=${id}`, { method: 'DELETE' });
      if (!res.ok) {
        const err = await res.json();
        throw new Error(err.message || 'Failed to delete');
      }

      devices = devices.filter(d => d.id !== id);
      success = 'Device deleted ✅';
      error = null;
    } catch (err) {
      error = err.message;
      success = null;
    }
  }

  async function createDeviceLink() {
    if (!selectedFromId || !selectedToId) return alert("Select both devices");
    if (selectedFromId === selectedToId) return alert("Cannot link a device to itself");

    const payload = {
      from_device_id: Number(selectedFromId),
      to_device_id: Number(selectedToId)
    };

    try {
      const res = await fetch(`${BASE_URL}/device_links/create`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
      });

      if (!res.ok) {
        const err = await res.json();
        throw new Error(err.message || 'Error creating link');
      }

      const newLink = await res.json();
      deviceLinks = [...deviceLinks, newLink];
      success = 'Link created successfully ✅';
      selectedFromId = null;
      selectedToId = null;
    } catch (err) {
      error = err.message;
      success = null;
    }
  }

  async function deleteLink(id: number) {
    if (!confirm("Delete this link?")) return;
    try {
      const res = await fetch(`${BASE_URL}/device_links/delete?id=${id}`, { method: 'DELETE' });
      if (!res.ok) {
        const err = await res.json();
        throw new Error(err.message || 'Failed to delete link');
      }
      deviceLinks = deviceLinks.filter(l => l.id !== id);
      success = 'Link deleted ✅';
    } catch (err) {
      error = err.message;
    }
  }

  function getDeviceName(id) {
    const d = devices.find(dev => dev.id === id);
    return d ? d.name : `Device #${id}`;
  }

  function startEdit(d) {
    device = { ...d };
    if (d.commissioning_date) device.commissioning_date = d.commissioning_date.split('T')[0];
    if (d.warranty_until) device.warranty_until = d.warranty_until.split('T')[0];
    selectedId = d.id;
    mode = 'edit';
    showForm = true;
  }

  $: filteredDevices = devices.filter((d) => {
    const q = searchQuery.toLowerCase();
    return (
      d.name?.toLowerCase().includes(q) ||
      d.ip?.toLowerCase().includes(q) ||
      d.manufacturer?.toLowerCase().includes(q) ||
      d.model_type?.toLowerCase().includes(q) ||
      d.description?.toLowerCase().includes(q)
    );
  });

  onMount(async () => {
    if (browser) {
      await fetchDevices();
      await fetchLinks();
    }
  });
</script>

<!-- 💼 CLEAN CORPORATE UI -->
<style>
  :global(body) {
    font-family: system-ui, sans-serif;
    line-height: 1.5;
    color: #333;
    background-color: #f5f7fa;
    -webkit-font-smoothing: antialiased;
    margin: 0;
    padding: 2rem;
  }

  h1, h2, h3 {
    color: #2c3e50;
  }

  button {
    background-color: #0077cc;
    color: white;
    border: none;
    padding: 0.5rem 1rem;
    font-weight: bold;
    border-radius: 5px;
    cursor: pointer;
    transition: background-color 0.2s ease;
  }

  button:hover {
    background-color: #005fa3;
  }

  input, textarea, select {
    width: 100%;
    padding: 0.45rem 0.6rem;
    border: 1px solid #ccc;
    border-radius: 4px;
    font-size: 0.95rem;
    background: white;
    transition: border-color 0.2s ease;
  }

  input:focus, textarea:focus, select:focus {
    border-color: #0077cc;
    outline: none;
  }

  label {
    font-weight: 500;
    font-size: 0.9rem;
    color: #444;
    margin-bottom: 4px;
    display: inline-block;
  }

  .form-container {
    background-color: white;
    border-radius: 8px;
    padding: 1.5rem;
    border: 1px solid #ddd;
    margin-top: 1rem;
    box-shadow: 0 2px 6px rgba(0,0,0,0.03);
  }

  .grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
    gap: 1rem;
  }

  table {
    width: 100%;
    border-collapse: collapse;
    font-size: 0.95rem;
    background-color: white;
  }

  th, td {
    padding: 0.65rem;
    border: 1px solid #e2e8f0;
    text-align: left;
  }

  th {
    background-color: #f0f4f8;
    font-weight: 600;
  }

  tr:nth-child(even) {
    background-color: #f9fbfc;
  }

  tr:hover {
    background-color: #eef3f9;
  }

  details summary {
    background-color: #eef3f9;
    padding: 0.75rem;
    font-size: 1.1rem;
    border-radius: 6px;
    cursor: pointer;
    border: 1px solid #d0dce7;
    margin-top: 1.5rem;
  }

  details summary:hover {
    background-color: #e3ecf7;
  }

  .success {
    color: #2e7d32;
    background: #e6f4ea;
    padding: 0.5rem 1rem;
    border: 1px solid #a5d6a7;
    border-radius: 4px;
    margin-bottom: 1rem;
  }

  .error {
    color: #c62828;
    background: #fdecea;
    padding: 0.5rem 1rem;
    border: 1px solid #f5c6cb;
    border-radius: 4px;
    margin-bottom: 1rem;
  }

  #search {
    margin: 0.5rem 0 1rem 0;
    border-radius: 6px;
  }
</style>

<!-- ✅ UI -->
<h1>📟 Device Management</h1>

{#if error}<p class="error">❌ {error}</p>{/if}
{#if success}<p class="success">✅ {success}</p>{/if}

<button on:click={() => {
  showForm = !showForm;
  if (!showForm) {
    device = resetDevice();
    selectedId = null;
    mode = 'create';
  }
}}>
  {showForm && mode === 'create' ? 'Cancel' : '➕ Create Device'}
</button>

{#if showForm}
  <div class="form-container">
    <h2>{mode === 'create' ? 'Create New Device' : 'Update Device'}</h2>
    <form on:submit|preventDefault={mode === 'create' ? createDevice : updateDevice}>
      <div class="grid">
        <div><label>Name</label><input bind:value={device.name} required /></div>
        <div><label>Hostname</label><input bind:value={device.hostname} /></div>
        <div><label>IP</label><input bind:value={device.ip} /></div>
        <div><label>Domain</label><input bind:value={device.domain} /></div>
        <div><label>Manufacturer</label><input bind:value={device.manufacturer} /></div>
        <div><label>Model Type</label><input bind:value={device.model_type} /></div>
        <div><label>Serial Numbers</label><input bind:value={device.serial_numbers} /></div>
        <div><label>MAC</label><input bind:value={device.mac} /></div>
        <div><label>Description</label><textarea bind:value={device.description}></textarea></div>
        <div><label>Equipment</label><input bind:value={device.equipment} /></div>
        <div><label>Function</label><input bind:value={device.function} /></div>
        <div><label>Settings</label><input bind:value={device.settings} /></div>
        <div><label>Device Link</label><input bind:value={device.device_link} /></div>
        <div><label>Commissioning Date</label><input type="date" bind:value={device.commissioning_date} /></div>
        <div><label>Origin</label><input bind:value={device.origin} /></div>
        <div><label>Warranty Service Number</label><input bind:value={device.warranty_service_number} /></div>
        <div><label>Warranty Until</label><input type="date" bind:value={device.warranty_until} /></div>
        <div><label>Licenses</label><input bind:value={device.licenses} /></div>
        <div><label>Location Text</label><input bind:value={device.location_text} /></div>
        <div><label>Department</label><input bind:value={device.department} /></div>
        <div><label>Internal Contact</label><input bind:value={device.internal_contact} /></div>
        <div><label>External Contact</label><input bind:value={device.external_contact} /></div>
        <div><label>Map Link</label><input bind:value={device.map_link} /></div>
        <div><label>Software Interfaces</label><input bind:value={device.software_interfaces} /></div>
        <div><label>Backup Method</label><input bind:value={device.backup_method} /></div>
        <div><label>Backup File Link</label><input bind:value={device.backup_file_link} /></div>
        <div><label>Software Asset</label><input bind:value={device.software_asset} /></div>
        <div><label>Password Link</label><input bind:value={device.password_link} /></div>
        <div><label>Internal Access</label><input bind:value={device.internal_access} /></div>
        <div><label>External Access</label><input bind:value={device.external_access} /></div>
        <div><label>Misc Links</label><input bind:value={device.misc_links} /></div>
        <div><label>Externally Accessible</label><input type="checkbox" bind:checked={device.externally_accessible} /></div>
        <div><label>Restart How</label><input bind:value={device.restart_how} /></div>
        <div><label>Restart Notes</label><input bind:value={device.restart_notes} /></div>
        <div><label>Restart Coordination</label><input bind:value={device.restart_coordination} /></div>
        <div><label>Network Connection</label><input bind:value={device.network_connection} /></div>
        <div><label>Patch Location</label><input bind:value={device.patch_location} /></div>
        <div><label>Documents</label><input bind:value={device.documents} /></div>
      </div>

      <button type="submit" style="margin-top: 1rem;">
        {mode === 'create' ? '✅ Submit Device' : '💾 Update Device'}
      </button>
    </form>
  </div>
{/if}

<details bind:open={showDevicesSection}>
  <summary>📟 All Devices</summary>

  <label for="search">🔍 Filter Devices</label>
  <input
    id="search"
    type="text"
    placeholder="Search by name, IP, model, manufacturer..."
    bind:value={searchQuery}
  />

  <div style="max-height: 420px; overflow-y: auto;">
    {#if filteredDevices.length === 0}
      <p>No devices found.</p>
    {:else}
      <table>
        <thead>
          <tr>
            <th>Name</th>
            <th>IP</th>
            <th>Manufacturer</th>
            <th>Model</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each filteredDevices as d}
            <tr>
              <td>{d.name}</td>
              <td>{d.ip}</td>
              <td>{d.manufacturer}</td>
              <td>{d.model_type}</td>
              <td>
                <button on:click={() => startEdit(d)}>Edit</button>
                <button on:click={() => deleteDevice(d.id)}>Delete</button>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    {/if}
  </div>
</details>

<details bind:open={showLinksSection}>
  <summary>🔗 Device Links</summary>

  <div class="form-container">
    <h3>Create New Link</h3>
    <div class="grid">
      <div>
        <label>From Device</label>
        <select bind:value={selectedFromId}>
          <option value="">Select Device</option>
          {#each devices as d}
            <option value={d.id}>{d.name}</option>
          {/each}
        </select>
      </div>
      <div>
        <label>To Device</label>
        <select bind:value={selectedToId}>
          <option value="">Select Device</option>
          {#each devices as d}
            <option value={d.id}>{d.name}</option>
          {/each}
        </select>
      </div>
    </div>
    <button on:click={createDeviceLink}>➕ Create Link</button>
  </div>

  <div style="max-height: 300px; overflow-y: auto;">
    {#if deviceLinks.length === 0}
      <p>No device links found.</p>
    {:else}
      <table>
        <thead>
          <tr>
            <th>ID</th>
            <th>From</th>
            <th>To</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each deviceLinks as link}
            <tr>
              <td>{link.id}</td>
              <td>{getDeviceName(link.from_device_id)}</td>
              <td>{getDeviceName(link.to_device_id)}</td>
              <td><button on:click={() => deleteLink(link.id)}>🗑️ Delete</button></td>
            </tr>
          {/each}
        </tbody>
      </table>
    {/if}
  </div>
</details>
