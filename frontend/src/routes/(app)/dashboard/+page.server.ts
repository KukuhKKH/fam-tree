import type { PageServerLoad } from './$types';
import { env } from '$env/dynamic/private';
const INTERNAL_BACKEND_URL = env.INTERNAL_BACKEND_URL;

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
			const families = result.data || [];

			// Fetch tree data from all families to compute stats
			let totalPersons = 0;
			let totalRelationships = 0;

			const treePromises = families.map(async (family: any) => {
				try {
					const treeRes = await fetch(`${INTERNAL_BACKEND_URL}/families/${family.slug}/tree`, {
						headers: { 'Cookie': `session_id=${sessionId}` }
					});
					if (treeRes.ok) {
						const treeResult = await treeRes.json();
						const tree = treeResult.data || treeResult;
						totalPersons += tree.nodes?.length || 0;
						totalRelationships += tree.edges?.length || 0;
					}
				} catch { /* skip failed tree fetches */ }
			});

			await Promise.all(treePromises);

			return {
				families,
				totalPersons,
				totalRelationships,
			};
		}
	} catch (error) {
		console.error('Fetch Families Error:', error);
	}

	return {
		families: [],
		totalPersons: 0,
		totalRelationships: 0,
	};
};
