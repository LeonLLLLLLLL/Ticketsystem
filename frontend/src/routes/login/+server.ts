// src/routes/login/+server.ts
import { json } from '@sveltejs/kit';

const isDocker = process.env.DOCKER === 'true';

const API_URL = isDocker
	? 'http://address_module_backend:8000'
	: 'http://localhost:8000';

export async function POST({ request, cookies }) {
	const { identifier, password } = await request.json();

	const res = await fetch(`${API_URL}/auth/login`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ identifier, password })
	});

	if (!res.ok) {
		const errText = await res.text();
		return new Response(errText, { status: 401 });
	}

	const data = await res.json();

	cookies.set('session', data.token, {
		httpOnly: true,
		path: '/',
		sameSite: 'strict',
		secure: process.env.NODE_ENV !== 'development',
		maxAge: 60 * 60 * 24
	});

	return json({ success: true });
}
