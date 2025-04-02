import BASE_URL from './BASE_URL.js'

// Función para obtener la lista de series, aplicando filtros y ordenamiento
const getSeriesList = async (queryParams = {}) => {
  // Construir los parámetros de query (por ejemplo: search, status, sort)
  const params = new URLSearchParams(queryParams)
  const url = `${BASE_URL}/series?${params.toString()}`

  try {
    const response = await fetch(url)
    if (!response.ok) {
      throw new Error(`Error fetching series: ${response.statusText}`)
    }
    return await response.json()
  } catch (error) {
    console.error('Failed to fetch series list:', error)
    throw error
  }
}

export default getSeriesList
