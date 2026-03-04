import type { PageServerLoad } from './$types';
import { apiFetch } from '$lib/api';
import { error } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ params, fetch }) => {
    const { slug, id } = params;

    try {
        const [family, person, persons, relationships] = await Promise.all([
            apiFetch(`/families/${slug}`, { fetch }),
            apiFetch(`/families/${slug}/persons/${id}`, { fetch }),
            apiFetch(`/families/${slug}/persons`, { fetch }),
            apiFetch(`/families/${slug}/relationships`, { fetch })
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
