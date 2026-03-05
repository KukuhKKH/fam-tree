import type { Handle } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';
const INTERNAL_BACKEND_URL = env.INTERNAL_BACKEND_URL;

export const handle: Handle = async ({ event, resolve }) => {
	const sessionId = event.cookies.get('session_id');

	if (sessionId) {
		try {
			const response = await event.fetch(`${INTERNAL_BACKEND_URL}/auth/me`, {
				headers: {
					'Cookie': `session_id=${sessionId}`
				}
			});

			if (response.ok) {
				const { data } = await response.json();
				event.locals.user = data;
			} else {
				event.locals.user = null;
			}
		} catch (error) {
			console.error('BFF Auth Error:', error);
			event.locals.user = null;
		}
	} else {
		event.locals.user = null;
	}

	// Protect app routes (dashboard and families management)
	const isProtectedRoute = event.url.pathname.startsWith('/dashboard') || 
		event.url.pathname.startsWith('/families');
	
	if (isProtectedRoute && !event.locals.user) {
		return Response.redirect(`${event.url.origin}/login`, 302);
	}

	return resolve(event);
};
