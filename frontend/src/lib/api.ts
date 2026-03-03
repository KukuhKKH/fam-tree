import { env } from '$env/dynamic/public';

interface FetchOptions extends RequestInit {
	params?: Record<string, string>;
}

export async function apiFetch<T = any>(path: string, options: FetchOptions = {}): Promise<T> {
	const baseUrl = env.PUBLIC_BACKEND_URL || 'http://localhost:8080';
	
	let url = `${baseUrl}${path.startsWith('/') ? '' : '/'}${path}`;
	
	if (options.params) {
		const searchParams = new URLSearchParams(options.params);
		url += `?${searchParams.toString()}`;
	}

	const response = await fetch(url, {
		...options,
		headers: {
			'Content-Type': 'application/json',
			...options.headers,
		},
		credentials: 'include', // Important for sending/receiving session cookies
	});

	if (!response.ok) {
		const errorData = await response.json().catch(() => ({}));
		throw new Error(errorData.messages?.[0] || response.statusText);
	}

	const result = await response.json();
	return result.data as T;
}
