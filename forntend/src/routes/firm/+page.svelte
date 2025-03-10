<script lang="ts">
	import { browser } from '$app/environment';
	import axios from "axios";
	import { onMount } from 'svelte';
	
	let formData = Array(11).fill("");
	let checkboxes = [false, false, false];
	let dropdownSelection = "";
	let remark = "";
  
	// Initialize contacts with empty array, will be filled from API
	let contacts: any[] = [];
	let loading = true;
	let error = '';
  
	let input_placeholders = ["Anrede", "Name 1", "Name 2", "Name 3", "Straße", "Land", "PLZ", "Ort", "Telefon", "Email", "Website"];
	let input_ids = ["anrede", "name_1", "name_2", "name_3", "straße", "land", "plz", "ort", "telefon", "email", "website"];
	let checkbox_text = ["Kunde", "Lieferant", "Gesperrt"];
  
	// API URL setup
	const API_URL = browser 
		? (window.location.hostname === "localhost"
			? "http://localhost:8000"
			: "http://address_module_backend:8000")
		: "http://address_module_backend:8000";

	interface ApiContact {
		anrede: string;
		name: string;
		position: string;
		telefon: string;
		mobil: string;
		email: string;
		abteilung: string;
		geburtstag: string;
		// Add any other properties that might exist in the API response
	}

	// Function to fetch contacts for a firm when page loads
	async function fetchContacts() {
		loading = true;
		error = '';
		
		try {
			const response = await axios.get(`${API_URL}/contact/get_by_id?username=tom&id=1`, {
				headers: {
					"Authorization": "123456789",
				},
			});
			
			// Transform the API response format to match our component's expected format
			contacts = response.data.contacts.map((contact: ApiContact) => ({
				anrede: contact.anrede,
				vorname: contact.name.split(' ')[0],
				nachname: contact.name.split(' ').slice(1).join(' '),
				position: contact.position,
				telefon: contact.telefon,
				mobil: contact.mobil,
				email: contact.email,
				abteilung: contact.abteilung,
				geburtstag: contact.geburtstag
			}));
			
			loading = false;
		} catch (err) {
			console.error("Error fetching contacts:", err);
			error = 'Failed to load contacts. Please try again later.';
			loading = false;
			
			// Fallback to example data if API fails
			contacts = [
				{ anrede: "Herr", vorname: "Max", nachname: "Mustermann", position: "IT Support", telefon: "+49 123 456789", mobil: "+49 987 654321", email: "max.mustermann@example.com", abteilung: "Technik", geburtstag: "01.01.1990" },
				{ anrede: "Frau", vorname: "Lisa", nachname: "Musterfrau", position: "Projektmanagerin", telefon: "+49 234 567890", mobil: "+49 876 543210", email: "lisa.musterfrau@example.com", abteilung: "Management", geburtstag: "15.05.1985" }
			];
		}
	}

	// Call API when page loads
	onMount(fetchContacts);

	async function submitForm() {
		let jsonData: Record<string, any> = {};
		input_ids.forEach((key, index) => {
			jsonData[key] = formData[index];
		});
		checkbox_text.forEach((key, index) => {
			jsonData[key.toLowerCase()] = checkboxes[index];
		});
		jsonData["bemerkung"] = remark;
		jsonData["firma_typ"] = dropdownSelection;

		try {
			const response = await axios.post(`${API_URL}/firm/submit?username=tom`, jsonData, {
				headers: {
					"Authorization": "123456789",
					"Content-Type": "application/json",
				},
			});
			console.log("Form submitted successfully:", response.data);
			// Refresh contacts list after successful submission
			fetchContacts();
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
			<h1 class="header-title">Adresse/Firma Modul</h1>
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
				<h2 class="section-title">Hauptregister</h2>
				
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
		
				<!-- Checkboxes -->
				<div class="checkbox-group">
					{#each checkboxes as checked, i}
						<label class="checkbox-label">
							<input type="checkbox" bind:checked={checkboxes[i]} />
							<span>{checkbox_text[i]}</span>
						</label>
					{/each}
				</div>
		
				<!-- Textarea and Dropdown -->
				<div class="form-extras">
					<textarea 
						class="textarea" 
						placeholder="Bemerkung" 
						bind:value={remark}
					></textarea>
		
					<select 
						bind:value={dropdownSelection} 
						class="dropdown"
					>
						<option value="" disabled selected>Firmentyp wählen</option>
						<option>Hauptfirma</option>
						<option>Niederlassung</option>
					</select>
				</div>
		
				<!-- Submit Button -->
				<button 
					class="submit-button" 
					on:click={submitForm}
				>
					Firma Speichern
				</button>
			</section>
  
			<!-- RIGHT SECTION: CONTACTS -->
			<section class="contacts-section">
				<h2 class="section-title">Zugehörige Ansprechpartner</h2>
				
				{#if loading}
					<div class="loading-state">
						<p>Kontakte werden geladen...</p>
					</div>
				{:else if error}
					<div class="error-state">
						<p>{error}</p>
					</div>
				{:else if contacts.length === 0}
					<div class="empty-state">
						<p>Keine Ansprechpartner gefunden.</p>
					</div>
				{:else}
					<div class="contacts-list">
						{#each contacts as contact}
							<div class="contact-card">
								<div class="contact-header">
									<h3>{contact.anrede} {contact.vorname} {contact.nachname}</h3>
									<p class="contact-role">{contact.position} - {contact.abteilung}</p>
								</div>
								<div class="contact-details">
									<p>📞 {contact.telefon || 'Nicht angegeben'}</p>
									<p>📱 {contact.mobil || 'Nicht angegeben'}</p>
									<p>✉ {contact.email || 'Nicht angegeben'}</p>
									<p>🎂 {contact.geburtstag || 'Nicht angegeben'}</p>
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

	.checkbox-group {
		display: flex;
		gap: 15px;
		margin-bottom: 15px;
	}

	.checkbox-label {
		display: flex;
		align-items: center;
		gap: 5px;
	}

	.form-extras {
		display: flex;
		flex-direction: column;
		gap: 15px;
	}

	.textarea {
		min-height: 100px;
		padding: 10px;
		border: 1px solid #bdc3c7;
		border-radius: 4px;
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

	/* Contacts Section */
	.contacts-section {
		background-color: #f8f9fa;
		border-radius: 8px;
		padding: 20px;
	}

	.contacts-list {
		display: flex;
		flex-direction: column;
		gap: 15px;
	}

	.contact-card {
		background-color: white;
		border: 1px solid #e9ecef;
		border-radius: 6px;
		padding: 15px;
		box-shadow: 0 2px 4px rgba(0,0,0,0.05);
	}

	.contact-header {
		margin-bottom: 10px;
		border-bottom: 1px solid #e9ecef;
		padding-bottom: 10px;
	}

	.contact-header h3 {
		color: #2c3e50;
		margin-bottom: 5px;
	}

	.contact-role {
		color: #7f8c8d;
		font-size: 0.9rem;
	}

	.contact-details p {
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