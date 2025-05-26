<script setup lang="ts">
import { onMounted, ref } from 'vue'

import { useAuthStore } from '../stores/auth'
import { useImageStore } from '../stores/images'
import { useRouter } from 'vue-router'
const auth = useAuthStore()
const images = useImageStore()
const router = useRouter()

const loading = ref<boolean>(!images.hasImages)

onMounted(() => {
    if (!auth.isAuthenticated) {
        router.push('/login')
    } else if (!images.hasImages) {
        images.fetchImages(auth.token as string)
        .then(() => {
        })
        .catch((error: any) => {
            console.error('Error fetching images:', error)
        })
        .finally(() => {
            loading.value = false
        })
    }
})

function logout() {
    auth.set(null)
    router.push('/login')
}
</script>

<template>
    <main>
        <div class="content">
            <div class="loading" v-if="loading">
                <p>Loading your photos...</p>
            </div>

            <div v-else>
                <div class="gallery">
                    <div class="image" v-for="image in images.imageUrls" :key="image">
                        <img :src="image" />
                    </div>
                </div>
                <div v-if="!images.hasImages" class="no-images">
                    <h1>No Images Found</h1>
                    <p>It seems you haven't uploaded any images yet.</p>
                    <ul>
                        <li>Click the "Upload" button to add your first photo.</li>
                        <li>Check back later to see your uploaded photos.</li>
                        <li>Explore the app to discover more features.</li>
                    </ul>
                </div>
            </div>
        </div>
    </main>
</template>

<style scoped>
.gallery {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(15rem, 0fr));
    grid-gap: 1rem;
    grid-auto-rows: 200px;
    align-items: center;
    justify-items: center;
    width: 80vw;
    align-items: center;
    place-content: space-evenly;
}

.images img {
    width: 100%;
    height: 22vw;
    object-fit: cover;
    border-radius: 0.75rem;
}

.toolbar {
    display: flex;
    justify-content: space-between;
    padding: 10px;
    background-color: #f0f0f0;
}
button {
    padding: 10px;
    border: none;
    border-radius: 5px;
    background-color: #007bff;
    color: white;
    cursor: pointer;
}

button:hover {
    background-color: #0056b3;
}
.content {
    padding: 20px;
    background-color: #ffffff;
    border-radius: 5px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}
.footer {
    text-align: center;
    padding: 10px;
    background-color: #f0f0f0;
    position: fixed;
    width: 100%;
    bottom: 0;
}
.no-images {
    text-align: center;
    margin-top: 20px;
}
.no-images h1 {
    font-size: 2rem;
    color: #333;
}
.no-images p {
    font-size: 1.2rem;
    color: #666;
}
.no-images ul {
    list-style-type: none;
    padding: 0;
}
.no-images li {
    font-size: 1rem;
    color: #333;
}
.no-images li::before {
    content: '• ';
    color: #007bff;
}
.no-images li {
    margin: 5px 0;
}
</style>
