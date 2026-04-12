const API_URL = import.meta.env.VITE_API_URL
const API_KEY = import.meta.env.VITE_API_KEY

const headers = {
    'Content-type': 'application/json',
    'x-api-key': API_KEY,
}

export async function getSets() {
    const response = await fetch(`${API_URL}/sets`, { headers });
    if(!response.ok) throw new Error('Failed to fetch sets');
    return response.json();
}

export async function getSetsById(id: string) {
    const response = await fetch(`${API_URL}/sets/${id}`, { headers });
    if(!response.ok) throw new Error('Failed to fetch set by id');
    return response.json();
}

export async function getCategoriesForSet(setId: string) {
    const response = await fetch(`${API_URL}/sets/${setId}/categories`, { headers });
    if(!response.ok) throw new Error('Failed to fetch categories for set');
    return response.json();
}

export async function createSet(url: string, djName: string) {
    const data = {
        'url': url,
        'dj_name': djName
    }
    const response = await fetch(`${API_URL}/sets`, { 
        method: 'POST', 
        headers, 
        body: JSON.stringify(data) 
    });
    if(!response.ok) throw new Error('Failed to create set');
    return response.json();
}

export async function deleteSetById(id: string) {
    const response = await fetch(`${API_URL}/sets/${id}`, { 
        method: 'DELETE',
        headers
    });
    if(!response.ok) throw new Error('Failed to delete set');
    return response.json();
}

export async function getCategories() {
    const response = await fetch(`${API_URL}/categories`, { headers });
    if(!response.ok) throw new Error('Failed to fetch categories');
    return response.json();
}

export async function createCategory(name: string) {
    const data = {
        'name': name
    }
    const response = await fetch(`${API_URL}/categories`, {
        method: 'POST',
        headers,
        body: JSON.stringify(data)
    });
    if(!response.ok) throw new Error('Failed to create category');
    return response.json();
}

export async function addCategoryToSet(setId: string, category: string) {
    const data = {
        'category': category
    }
    const response = await fetch(`${API_URL}/sets/${setId}/categories`, {
        method: 'POST',
        headers,
        body: JSON.stringify(data)
    });
    if(!response.ok) throw new Error('Failed to add category');
    return response.text();
}

export async function searchSets(query: string) {
    const response = await
    fetch(`${API_URL}/sets/search?q=${query}`, { headers });
    if(!response.ok) throw new Error('Failed to search sets');
    return response.json();
}