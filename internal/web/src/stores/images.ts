import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

type imageResponseType = {
    photo_ids: Array<string>
}

export const useImageStore = defineStore('imageStore', () => {
    // TODO: Create data structure to better manage images
    const imageIds = ref<string[]>([])
    const imageUrls = ref<string[]>([])
    const addImage = (id: string, photo: string) => {
        imageIds.value.push(id)
        imageUrls.value.push(photo)
    }
    const removeImage = (image: string) => {
        imageIds.value = imageIds.value.filter((img) => img !== image)
    }
    const clearImages = () => {
        imageIds.value = []
        imageUrls.value = []
    }
    const getImages = computed(() => imageIds.value)

    const imageCount = computed(() => imageIds.value.length)
    const hasImages = computed(() => imageIds.value.length > 0)
    const getImageCount = () => imageCount.value

    const fetchThumbnail = async (token: string, id: string): Promise<string | void> => {
        return await fetch(`${import.meta.env.VITE_API_URI}/api/photo/thumbnail`, {
            method: 'GET',
            headers: {
                Authorization: `Bearer ${token}`,
                photo_id: id,
            },
        })
            .then(async (response) => {
                if (response.ok) {
                    // Handle successful login
                    const photo = await response.blob()
                    return URL.createObjectURL(photo)
                } else {
                    // Handle error response
                    throw new Error('Failed to fetch thumbnail')
                }
            })
            .catch((error) => {
                console.error('Error fetching thumbnail:', error)
                throw error
            })
    }

    const fetchImages = async (token: string): Promise<void> => {
        if (!token) {
            console.warn('No token provided for fetching images')
            return
        }
        console.log('Using token:', token)
        fetch(`${import.meta.env.VITE_API_URI}/api/photo/ids`, {
            method: 'GET',
            headers: {
                Authorization: `Bearer ${token}`,
            },
        })
            .then(async (response) => {
                if (response.ok) {
                    // Handle successful login
                    const photo_ids = ((await response.json()) as imageResponseType).photo_ids

                    // Remove existing images before adding new ones
                    clearImages()

                    for (const id of photo_ids) {
                        const photo = await fetchThumbnail(token, id)
                        console.log(`Fetched thumbnail for ID: ${id}`, photo)
                        if (photo === undefined) {
                            console.warn(`No photo found for ID: ${id} ${photo}`)
                            continue // Skip if no photo is found
                        }
                        addImage(id, photo)
                    }
                } else {
                    // Handle error response
                    throw new Error('Failed to fetch images')
                }
            })
            .catch((error) => {
                console.error('Error fetching images:', error)
                throw error
            })
    }

    const uploadImage = async (token: string, file: File): Promise<void> => {
        const formData = new FormData()
        formData.append('photo', file)

        return await fetch(`${import.meta.env.VITE_API_URI}/api/photo`, {
            method: 'POST',
            headers: {
                Authorization: `Bearer ${token}`,
            },
            body: formData,
        })
            .then(async (response) => {
                if (response.ok) {
                    // Handle successful upload
                    const photo = await response.json()
                    console.log('Image uploaded successfully:', photo)
                    imageIds.value.push(photo.photo_id)
                    const thumbnail = await fetchThumbnail(token, photo.photo_id)
                    if (thumbnail) {
                        imageUrls.value.push(thumbnail)
                    } else {
                        console.warn(`No thumbnail found for uploaded photo ID: ${photo.photo_id}`)
                    }
                } else {
                    // Handle error response
                    throw new Error('Failed to upload image')
                }
            })
            .catch((error) => {
                console.error('Error uploading image:', error)
                throw error
            })
    }

    return {
        imageIds,
        imageUrls,
        addImage,
        removeImage,
        clearImages,
        getImages,
        hasImages,
        getImageCount,
        fetchImages,
        uploadImage,
    }
})
