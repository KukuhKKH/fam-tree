import type { PageServerLoad } from './$types';
import { apiFetch } from '$lib/api';
import { error } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ params, fetch, cookies }) => {
    const { slug, id } = params;
    const sessionId = cookies.get('session_id');
    const cookieString = sessionId ? `session_id=${sessionId}` : '';

    try {
        const [family, person, persons, relationships] = await Promise.all([
            apiFetch(`/families/${slug}`, { fetch, cookie: cookieString }),
            apiFetch(`/families/${slug}/persons/${id}`, { fetch, cookie: cookieString }),
            apiFetch(`/families/${slug}/persons`, { fetch, cookie: cookieString }),
            apiFetch(`/families/${slug}/relationships`, { fetch, cookie: cookieString })
        ]);

        return {
            family,
            person,
            persons,
            relationships,
            slug
        };
    } catch (err: any) {
        console.error(`Error loading person [${id}] in family [${slug}]:`, err.message);
        throw error(404, 'Data person tidak ditemukan');
    }
};
