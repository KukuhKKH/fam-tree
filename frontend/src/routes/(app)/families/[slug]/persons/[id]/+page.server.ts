import type { PageServerLoad } from './$types';
import { apiFetch } from '$lib/api';
import { error } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ params, fetch }) => {
    const { slug, id } = params;

    try {
        const [family, person, persons] = await Promise.all([
            apiFetch(`/families/${slug}`, { fetch }),
            apiFetch(`/families/${slug}/persons/${id}`, { fetch }),
            apiFetch(`/families/${slug}/persons`, { fetch })
        ]);

        return {
            family,
            person,
            persons,
            slug
        };
    } catch (err: any) {
        console.error(`Error loading person [${id}] in family [${slug}]:`, err.message);
        throw error(404, 'Data person tidak ditemukan');
    }
};
