// src/hooks.server.ts
import type { Handle } from '@sveltejs/kit';
import { redirect } from '@sveltejs/kit';

export const handle: Handle = async ({ event, resolve }) => {
	// Example: read a cookie to check login status
	const session = event.cookies.get('session');

	// Allow access to login page even if not logged in
	//if (!session && event.url.pathname !== '/login') {
		// Redirect to login
	//	throw redirect(302, '/login');
	//}

	// Continue to requested page
	return resolve(event);
};
