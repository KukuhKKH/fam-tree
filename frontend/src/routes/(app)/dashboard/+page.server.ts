import type { PageServerLoad } from './$types';
import { INTERNAL_BACKEND_URL } from '$env/static/private';

export const load: PageServerLoad = async ({ fetch, cookies }) => {
	const sessionId = cookies.get('session_id');
	
	try {
		const response = await fetch(`${INTERNAL_BACKEND_URL}/families`, {
			headers: {
				'Cookie': `session_id=${sessionId}`
			}
		});

		if (response.ok) {
			const result = await response.json();
			return {
				families: result.data || []
			};
		}
	} catch (error) {
		console.error('Fetch Families Error:', error);
	}

	return {
		families: []
	};
};
