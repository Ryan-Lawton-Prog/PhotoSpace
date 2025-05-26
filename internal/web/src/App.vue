<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'

import { useAuthStore } from '@/stores/auth'
import { ref } from 'vue';
import { useImageStore } from './stores/images';
const auth = useAuthStore()
const images = useImageStore()

function logout() {
    auth.set(null)
}

function handleFileChange(event: Event) {
    const input = event.target as HTMLInputElement;
    if (!input.files || input.files.length === 0) {
        return;
    }

    images.uploadImage( auth.token as string, input.files[0],)
        .then(() => {
            input.value = '';
        })
        .catch((error: any) => {
            console.error('Error uploading image:', error);
        });
}
</script>

<template>
    <header>
        <div class="wrapper">
            <img alt="Vue logo" class="logo" src="@/assets/logo.svg" width="25" height="25" />
            <nav>
                <RouterLink v-if="!auth.isAuthenticated" to="/login">Login</RouterLink>
                <RouterLink v-if="!auth.isAuthenticated" to="/register">Register</RouterLink>
                <RouterLink v-if="auth.isAuthenticated" to="/dashboard">Dashboard</RouterLink>
                <RouterLink v-if="auth.isAuthenticated" to="/upload">Upload</RouterLink>
                <RouterLink v-if="auth.isAuthenticated" to="/login" @click="logout"
                    >Logout</RouterLink
                >
            </nav>
            <div v-if="auth.isAuthenticated" class ="file-upload">
                <label for="file-input" class="file-label">
                    ➕
                </label>
                <input type="file" id="file-input" @change="handleFileChange" style="display: none;">
            </div>
        </div>
    </header>

    <div class="content">
        <RouterView />
    </div>

    <footer>
        <p>&copy; 2023 PhotoSpace. All rights reserved.</p>
    </footer>
</template>

<style scoped>
header {
    line-height: 1.5;
    position: fixed;
    top: 0;
    align-content: left;
    left: 2rem;
    width: 100%;
    padding-top: 1rem;
    padding-bottom: 1rem;
    background-color: var(--color-background);
}

.logo {
    margin: 0 auto 2rem;
}

nav {
    font-size: 12px;
    text-align: center;
    margin-top: 2re;
}

nav a.router-link-exact-active {
    color: var(--color-text);
}

nav a.router-link-exact-active:hover {
    background-color: transparent;
}

nav button {
    background-color: transparent;
    color: var(--color-text);
    border: none;
    padding: 0.5rem 1rem;
    cursor: pointer;
}

nav a {
    display: inline-block;
    padding: 0 1rem;
    border-left: 1px solid var(--color-border);
}

nav a:first-of-type {
    border: 0;
}

@media (min-width: 500px) {
    header {
        display: flex;
    }

    .logo {
        margin: 0 2rem 0 0;
    }

    header .wrapper {
        display: flex;
        flex-direction: row;
        justify-content: space-between;
    }

    nav {
        text-align: left;
        margin-left: -1rem;
        font-size: 1rem;
    }
}

.file-upload {
    align-self: flex-end;
    background-color: transparent;
    border-radius: 50%;
    padding: 5px;
    margin-left: auto;
    margin-right: 1rem;

}

.file-upload .file-label {
    background-color: var(--color-button-background);
    color: white;
    border-radius: 5px;
    cursor: pointer;
}

.file-upload:hover {
    background-color: gray;
    color: var(--color-button-hover-text);
    border-color: gray;
}

.content {
    margin-top: 57px;
    background-color: var(--color-background);
    color: var(--color-text);
}

.wrapper {
    display: flex;
    flex-direction: row;
    align-items: center;
    padding: 0 2rem;
    width: 100%;
}

footer {
    text-align: center;
    padding: 1rem;
    background-color: var(--color-background);
    color: var(--color-text);
    bottom: 0;
    left: 0;
    position: fixed;
}
</style>
