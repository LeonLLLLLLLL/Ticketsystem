<script lang="ts">
	import { browser } from '$app/environment';
	import axios from "axios";
	import { onMount } from 'svelte';
	
	let formData = Array(9).fill("");
	let dropdownSelection = "";
	
	// Initialize firms with empty array, will be filled from API
	let firms = [];
	let loading = true;
	let error = '';

	let input_placeholders = ["Anrede", "Vorname", "Nachname", "Position", "Abteilung", "Telefon", "Mobil", "Email", "Geburtstag"];
	let input_ids = ["anrede", "vorname", "nachname", "position", "abteilung", "telefon", "mobil", "email", "geburtstag"];
	
	// API URL setup
	const API_URL = browser 
		? (window.location.hostname === "localhost"
			? "http://localhost:8000"
			: "http://address_module_backend:8000")
		: "http://address_module_backend:8000";

	// Function to fetch firms for a contact when page loads
	async function fetchFirms() {
		loading = true;
		error = '';
		
		try {
			const response = await axios.get(`${API_URL}/firm/get_by_id?username=tom&id=1`, {
				headers: {
					"Authorization": "123456789",
				},
			});
			
			// Transform the API response format to match our component's expected format
			firms = response.data.firms.map(firm => ({
				name: firm.name_1, // Assuming name_1 is the company name
				adresse: `${firm.straße}, ${firm.plz} ${firm.ort}`,
				telefon: firm.telefon,
				email: firm.email,
				website: firm.website,
				firma_typ: firm.firma_typ
			}));
			
			loading = false;
		} catch (err) {
			console.error("Error fetching firms:", err);
			error = 'Failed to load firms. Please try again later.';
			loading = false;
			
			// Fallback to example data if API fails
			firms = [
				{ name: "Musterfirma GmbH", adresse: "Hauptstr. 1, 10115 Berlin", telefon: "+49 30 12345678", email: "info@musterfirma.com", website: "www.musterfirma.com", firma_typ: "Hauptfirma" },
				{ name: "Example Corp.", adresse: "Musterweg 42, 20354 Hamburg", telefon: "+49 40 87654321", email: "contact@example.com", website: "www.example.com", firma_typ: "Niederlassung" }
			];
		}
	}

	// Call API when page loads
	onMount(fetchFirms);

	async function submitForm() {
		let jsonData: Record<string, any> = {};
		input_ids.forEach((key, index) => {
			jsonData[key] = formData[index];
		});
		jsonData["firma"] = dropdownSelection;

		try {
			const response = await axios.post(`${API_URL}/contact/submit?username=tom`, jsonData, {
				headers: {
					"Authorization": "123456789",
					"Content-Type": "application/json",
				},
			});
			console.log("Contact submitted successfully:", response.data);
			// Refresh firms list after successful submission
			fetchFirms();
		} catch (err) {
			console.error("Submission error:", err);
		}
	}
</script>

<div class="page-container">
	<!-- HEADER -->
	<header class="header">
		<div class="header-content">
			<div class="header-nav">
				<button class="nav-icon">🏠</button>
				<button class="nav-icon">❓</button>
				<button class="nav-icon">☰</button>
			</div>
			<h1 class="header-title">Adresse/Kontakte Modul</h1>
			<div class="header-actions">
				<button class="nav-icon">🔍</button>
				<button class="nav-icon">🔔</button>
				<button class="nav-icon">🚪</button>
			</div>
		</div>
	</header>

	<!-- MAIN CONTENT -->
	<main class="main-content">
		<div class="content-grid">
			<!-- LEFT SECTION: FORM -->
			<section class="form-section">
				<h2 class="section-title">Kontakt anlegen</h2>
				
				<div class="input-grid">
					{#each formData as input, i}
						<div class="input-wrapper">
							<label for="{input_ids[i]}" class="input-label">{input_placeholders[i]}</label>
							<input 
								bind:value={formData[i]} 
								id="{input_ids[i]}" 
								class="input-field" 
								placeholder="{input_placeholders[i]}" 
							/>
						</div>
					{/each}
				</div>
		
				<!-- Dropdown for firm selection -->
				<div class="form-extras">
					<select 
						bind:value={dropdownSelection} 
						class="dropdown"
					>
						<option value="" disabled selected>Firma zuordnen</option>
						{#each firms as firm}
							<option>{firm.name}</option>
						{/each}
					</select>
				</div>
		
				<!-- Submit Button -->
				<button 
					class="submit-button" 
					on:click={submitForm}
				>
					Kontakt Speichern
				</button>
			</section>

			<!-- RIGHT SECTION: FIRMS -->
			<section class="firms-section">
				<h2 class="section-title">Zugeordnete Firmen</h2>
				
				{#if loading}
					<div class="loading-state">
						<p>Firmen werden geladen...</p>
					</div>
				{:else if error}
					<div class="error-state">
						<p>{error}</p>
					</div>
				{:else if firms.length === 0}
					<div class="empty-state">
						<p>Keine Firmen gefunden.</p>
					</div>
				{:else}
					<div class="firms-list">
						{#each firms as firm}
							<div class="firm-card">
								<div class="firm-header">
									<h3>{firm.name}</h3>
									<p class="firm-type">{firm.firma_typ || 'Keine Angabe'}</p>
								</div>
								<div class="firm-details">
									<p>📍 {firm.adresse || 'Keine Angabe'}</p>
									<p>📞 {firm.telefon || 'Keine Angabe'}</p>
									<p>✉ {firm.email || 'Keine Angabe'}</p>
									<p>🌐 {firm.website || 'Keine Angabe'}</p>
								</div>
							</div>
						{/each}
					</div>
				{/if}
			</section>
		</div>
	</main>
</div>

<style>
	/* Reset and Base Styles */
	* {
		margin: 0;
		padding: 0;
		box-sizing: border-box;
	}

	/* Page Layout */
	.page-container {
		display: flex;
		flex-direction: column;
		min-height: 100vh;
		background-color: #f4f7fa;
		font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
	}

	/* Header Styles */
	.header {
		background-color: #2c3e50;
		color: white;
		padding: 15px 20px;
		box-shadow: 0 2px 4px rgba(0,0,0,0.1);
	}

	.header-content {
		display: flex;
		justify-content: space-between;
		align-items: center;
		max-width: 1200px;
		margin: 0 auto;
	}

	.header-title {
		font-size: 1.25rem;
		font-weight: 600;
	}

	.nav-icon {
		background: none;
		border: none;
		color: white;
		font-size: 1.25rem;
		cursor: pointer;
		margin: 0 10px;
		transition: color 0.3s ease;
	}

	.nav-icon:hover {
		color: #3498db;
	}

	/* Main Content */
	.main-content {
		flex-grow: 1;
		padding: 20px;
	}

	.content-grid {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 30px;
		max-width: 1200px;
		margin: 0 auto;
		background-color: white;
		border-radius: 10px;
		box-shadow: 0 4px 6px rgba(0,0,0,0.1);
		padding: 30px;
	}

	/* Form Section */
	.form-section {
		display: flex;
		flex-direction: column;
		gap: 20px;
	}

	.section-title {
		text-align: center;
		color: #2c3e50;
		border-bottom: 2px solid #3498db;
		padding-bottom: 10px;
		margin-bottom: 20px;
	}

	.input-grid {
		display: grid;
		grid-template-columns: repeat(2, 1fr);
		gap: 15px;
	}

	.input-wrapper {
		display: flex;
		flex-direction: column;
	}

	.input-label {
		margin-bottom: 5px;
		color: #34495e;
		font-size: 0.9rem;
	}

	.input-field {
		padding: 10px;
		border: 1px solid #bdc3c7;
		border-radius: 4px;
		transition: border-color 0.3s ease;
	}

	.input-field:focus {
		outline: none;
		border-color: #3498db;
	}

	.form-extras {
		display: flex;
		flex-direction: column;
		gap: 15px;
	}

	.dropdown {
		padding: 10px;
		border: 1px solid #bdc3c7;
		border-radius: 4px;
	}

	.submit-button {
		background-color: #3498db;
		color: white;
		border: none;
		padding: 12px;
		border-radius: 4px;
		cursor: pointer;
		transition: background-color 0.3s ease;
	}

	.submit-button:hover {
		background-color: #2980b9;
	}

	/* Firms Section */
	.firms-section {
		background-color: #f8f9fa;
		border-radius: 8px;
		padding: 20px;
	}

	.firms-list {
		display: flex;
		flex-direction: column;
		gap: 15px;
	}

	.firm-card {
		background-color: white;
		border: 1px solid #e9ecef;
		border-radius: 6px;
		padding: 15px;
		box-shadow: 0 2px 4px rgba(0,0,0,0.05);
	}

	.firm-header {
		margin-bottom: 10px;
		border-bottom: 1px solid #e9ecef;
		padding-bottom: 10px;
	}

	.firm-header h3 {
		color: #2c3e50;
		margin-bottom: 5px;
	}

	.firm-type {
		color: #7f8c8d;
		font-size: 0.9rem;
	}

	.firm-details p {
		margin-bottom: 5px;
		color: #34495e;
	}
	
	/* Loading and Error States */
	.loading-state, .error-state, .empty-state {
		padding: 20px;
		text-align: center;
		background-color: white;
		border-radius: 6px;
		border: 1px solid #e9ecef;
	}

	.loading-state {
		color: #3498db;
	}

	.error-state {
		color: #e74c3c;
	}

	.empty-state {
		color: #7f8c8d;
	}

	/* Responsive Adjustments */
	@media (max-width: 1024px) {
		.content-grid {
			grid-template-columns: 1fr;
		}
	}
</style>