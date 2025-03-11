<script lang="ts">
	import { browser } from '$app/environment';
	import axios from "axios";
	import { onMount } from 'svelte';
	
	let formData = Array(11).fill("");
	let checkboxes = [false, false, false];
	let dropdownSelection = "";
	let remark = "";
	let selectedContacts: number[] = [];
	let allContacts: any[] = [];
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

	// Function to fetch all contacts when page loads
	async function fetchAllContacts() {
		loading = true;
		error = '';
		
		try {
			const response = await axios.get(`${API_URL}/contact/get?username=tom`, {
				headers: {
					"Authorization": "123456789",
				},
			});
			
			// Transform the API response format for our component
			allContacts = response.data.contacts.map((contact: any) => ({
				id: contact.id,
				anrede: contact.anrede,
				vorname: contact.name.split(' ')[0],
				nachname: contact.name.split(' ').slice(1).join(' '),
				name: contact.name,
				position: contact.position,
				telefon: contact.telefon,
				mobil: contact.mobil,
				email: contact.email,
				abteilung: contact.abteilung,
				geburtstag: contact.geburtstag,
				selected: false // Add a property to track selection state
			}));
			
			loading = false;
		} catch (err) {
			console.error("Error fetching contacts:", err);
			error = 'Failed to load contacts. Please try again later.';
			loading = false;
			
			// Fallback to example data if API fails
			allContacts = [
				{ id: 1, anrede: "Herr", vorname: "Max", nachname: "Mustermann", name: "Max Mustermann", position: "IT Support", telefon: "+49 123 456789", mobil: "+49 987 654321", email: "max.mustermann@example.com", abteilung: "Technik", geburtstag: "01.01.1990", selected: false },
				{ id: 2, anrede: "Frau", vorname: "Lisa", nachname: "Musterfrau", name: "Lisa Musterfrau", position: "Projektmanagerin", telefon: "+49 234 567890", mobil: "+49 876 543210", email: "lisa.musterfrau@example.com", abteilung: "Management", geburtstag: "15.05.1985", selected: false }
			];
		}
	}

	// Call API when page loads
	onMount(fetchAllContacts);

	// Toggle contact selection
	function toggleContactSelection(id: number) {
		const index = allContacts.findIndex(contact => contact.id === id);
		if (index !== -1) {
			allContacts[index].selected = !allContacts[index].selected;
			// Update selectedContacts array
			if (allContacts[index].selected) {
				selectedContacts = [...selectedContacts, id];
			} else {
				selectedContacts = selectedContacts.filter(contactId => contactId !== id);
			}
			// Log selected contacts to console
			console.log("Selected contacts:", selectedContacts);
			allContacts = [...allContacts]; // Trigger UI update
		}
	}

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
		
		// Add selected contacts IDs to the JSON data
		jsonData["contacts"] = selectedContacts;

		try {
			const response = await axios.post(`${API_URL}/firm/submit?username=tom`, jsonData, {
				headers: {
					"Authorization": "123456789",
					"Content-Type": "application/json",
				},
			});
			console.log("Form submitted successfully:", response.data);
			// Refresh contacts list and reset selections
			fetchAllContacts();
			selectedContacts = [];
			formData = Array(11).fill("");
			checkboxes = [false, false, false];
			remark = "";
			dropdownSelection = "";
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
				
				<!-- Contact Selection Area -->
				<div class="contact-selection">
					<h3 class="subsection-title">Kontakte zuordnen</h3>
					<p class="help-text">Wählen Sie einen oder mehrere Kontakte aus:</p>
					
					{#if loading}
						<div class="loading-state">
							<p>Kontakte werden geladen...</p>
						</div>
					{:else if error}
						<div class="error-state">
							<p>{error}</p>
						</div>
					{:else if allContacts.length === 0}
						<div class="empty-state">
							<p>Keine Kontakte gefunden.</p>
						</div>
					{:else}
						<div class="contacts-checkboxes">
							{#each allContacts as contact (contact.id)}
								<label class="contact-checkbox">
									<input 
										type="checkbox" 
										checked={contact.selected} 
										on:change={() => toggleContactSelection(contact.id)} 
									/>
									<span class="contact-name">{contact.anrede} {contact.name}</span>
									<span class="contact-role">({contact.position || 'Keine Angabe'})</span>
								</label>
							{/each}
						</div>
						
						<!-- Selected Contacts Preview -->
						{#if selectedContacts.length > 0}
							<div class="selected-preview">
								<h4>Ausgewählte Kontakte:</h4>
								<ul class="selected-list">
									{#each selectedContacts as id}
										{@const contact = allContacts.find(c => c.id === id)}
										{#if contact}
											<li>{contact.anrede} {contact.name} <button class="remove-btn" on:click={() => toggleContactSelection(id)}>✕</button></li>
										{/if}
									{/each}
								</ul>
							</div>
						{/if}
					{/if}
				</div>
		
				<!-- Submit Button -->
				<button 
					class="submit-button" 
					on:click={submitForm}
					disabled={loading}
				>
					Firma Speichern
				</button>
			</section>
  
			<!-- RIGHT SECTION: CONTACTS -->
			<section class="contacts-section">
				<h2 class="section-title">Kontaktauswahl</h2>
				
				{#if loading}
					<div class="loading-state">
						<p>Kontakte werden geladen...</p>
					</div>
				{:else if error}
					<div class="error-state">
						<p>{error}</p>
					</div>
				{:else if allContacts.length === 0}
					<div class="empty-state">
						<p>Keine Kontakte gefunden.</p>
					</div>
				{:else}
					<div class="contacts-list">
						{#each allContacts as contact (contact.id)}
							<div class="contact-card {contact.selected ? 'selected' : ''}" on:click={() => toggleContactSelection(contact.id)}>
								<div class="contact-header">
									<h3>{contact.anrede} {contact.vorname} {contact.nachname}</h3>
									<p class="contact-role">{contact.position} - {contact.abteilung}</p>
									{#if contact.selected}
										<div class="selected-badge">✓</div>
									{/if}
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
	
	.subsection-title {
		color: #2c3e50;
		margin-bottom: 10px;
		font-size: 1.1rem;
	}
	
	.help-text {
		color: #7f8c8d;
		font-size: 0.9rem;
		margin-bottom: 10px;
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
	
	/* Contact Selection */
	.contact-selection {
		background-color: #f8f9fa;
		border-radius: 8px;
		padding: 15px;
		margin-top: 15px;
	}
	
	.contacts-checkboxes {
		display: flex;
		flex-direction: column;
		gap: 8px;
		max-height: 200px;
		overflow-y: auto;
		padding: 5px;
		border: 1px solid #e9ecef;
		border-radius: 4px;
		background-color: white;
	}
	
	.contact-checkbox {
		display: flex;
		align-items: center;
		padding: 8px;
		border-radius: 4px;
		transition: background-color 0.2s;
		cursor: pointer;
	}
	
	.contact-checkbox:hover {
		background-color: #f1f5f9;
	}
	
	.contact-name {
		margin-left: 8px;
		font-weight: 500;
	}
	
	.contact-role {
		margin-left: 8px;
		font-size: 0.8rem;
		color: #7f8c8d;
	}
	
	/* Selected Contacts Preview */
	.selected-preview {
		margin-top: 15px;
		padding: 10px;
		background-color: #e1f5fe;
		border-radius: 6px;
	}
	
	.selected-preview h4 {
		margin-bottom: 8px;
		font-size: 0.9rem;
		color: #0277bd;
	}
	
	.selected-list {
		list-style: none;
	}
	
	.selected-list li {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 5px 0;
		border-bottom: 1px dashed #b3e5fc;
	}
	
	.remove-btn {
		background: none;
		border: none;
		color: #e57373;
		cursor: pointer;
		font-weight: bold;
	}

	.submit-button {
		background-color: #3498db;
		color: white;
		border: none;
		padding: 12px;
		border-radius: 4px;
		cursor: pointer;
		transition: background-color 0.3s ease;
		margin-top: 20px;
	}

	.submit-button:hover {
		background-color: #2980b9;
	}
	
	.submit-button:disabled {
		background-color: #95a5a6;
		cursor: not-allowed;
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
		max-height: 600px;
		overflow-y: auto;
	}

	.contact-card {
		background-color: white;
		border: 1px solid #e9ecef;
		border-radius: 6px;
		padding: 15px;
		box-shadow: 0 2px 4px rgba(0,0,0,0.05);
		cursor: pointer;
		transition: all 0.2s ease;
		position: relative;
	}
	
	.contact-card:hover {
		transform: translateY(-3px);
		box-shadow: 0 4px 8px rgba(0,0,0,0.1);
	}
	
	.contact-card.selected {
		border-color: #3498db;
		background-color: #ebf5fb;
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
	
	.selected-badge {
		position: absolute;
		top: 10px;
		right: 10px;
		width: 24px;
		height: 24px;
		background-color: #3498db;
		color: white;
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		font-weight: bold;
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