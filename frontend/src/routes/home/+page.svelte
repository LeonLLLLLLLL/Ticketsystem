<script lang="ts">
	import { goto } from '$app/navigation';
	let isMenuOpen = $state(false);
	let hoveredModule = $state('');

	let modules = [
		"Adresse Modul",
		"Ticket Modul",
		"User Modul",
		"Statistik Modul",
		"Einstellungen",
	];

	function handleModuleClick(module: string) {
		console.log(`Navigating to: ${module}`);
		if (module === "Adresse Modul") {
			goto('/address');
		} else if (module === "Ticket Modul") {
			goto('/ticket');
		} else if (module === "User Modul") {
			goto('/user');
		} else if (module === "Statistik Modul") {
			goto('/statistics');
		} else if (module === "Einstellungen") {
			goto('/settings');
		}
	}

	function handleSubClick(sub: string) {
		if (sub === "Contact") goto('/address/contact');
		if (sub === "Firm") goto('/address/firm');
	}
</script>

<div class="container">
	<header class="header">
		<div class="nav-left">
			<button title="Home">🏠</button>
			<button title="Help">❓</button>
			<button class="Menu" on:click={() => isMenuOpen = !isMenuOpen}>☰</button>
		</div>
		<h1>Ticketsystem</h1>
		<div class="nav-right">
			<button title="Search">🔍</button>
			<button title="Notifications">🔔</button>
			<button title="Logout">🚪</button>
		</div>
	</header>

	{#if isMenuOpen}
		<div class="dropdown">
			{#each modules as module}
				<div 
					class="module-item" 
					on:mouseenter={() => hoveredModule = module}
					on:mouseleave={() => hoveredModule = ''}
				>
					<button on:click={() => handleModuleClick(module)}>{module}</button>
					{#if module === 'Adresse Modul' && hoveredModule === 'Adresse Modul'}
						<div class="submenu">
							<button on:click={() => handleSubClick('Contact')}>Contact</button>
							<button on:click={() => handleSubClick('Firm')}>Firm</button>
						</div>
					{/if}
				</div>
			{/each}
		</div>
	{/if}
</div>

<style>
	.container {
		display: flex;
		flex-direction: column;
		min-height: 100vh;
	}
	header.header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		background: #34495e;
		color: white;
		padding: 15px 30px;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
	}
	.nav-left, .nav-right {
		display: flex;
		gap: 15px;
	}
	button {
		background: none;
		border: none;
		color: inherit;
		font-size: 1.2rem;
		cursor: pointer;
		transition: color 0.3s;
	}
	button:hover {
		color: #1abc9c;
	}
	h1 {
		font-size: 1.5rem;
		font-weight: bold;
	}
	.dropdown {
		background: #ecf0f1;
		padding: 20px;
		display: flex;
		flex-direction: column;
		gap: 10px;
		width: fit-content;
	}
	.module-item {
		position: relative;
	}
	.module-item > button {
		background: #1abc9c;
		color: white;
		padding: 10px 15px;
		border-radius: 6px;
		font-size: 1rem;
		box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
	}
	.submenu {
		position: absolute;
		top: 0;
		left: 110%;
		background: #fff;
		box-shadow: 0 0 5px rgba(0, 0, 0, 0.15);
		border-radius: 5px;
		padding: 10px;
		display: flex;
		flex-direction: column;
		gap: 5px;
	}
	.submenu button {
		background: #3498db;
		color: white;
		padding: 6px 12px;
		border-radius: 4px;
		font-size: 0.9rem;
	}
</style>
