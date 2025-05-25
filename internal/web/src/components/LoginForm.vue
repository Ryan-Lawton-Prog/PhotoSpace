<script setup lang="ts">

defineProps<{
}>()

import { ref } from 'vue';
const username = ref<string>('');
const password = ref<string>('');
import { useRouter } from 'vue-router';
const router = useRouter();

import { useAuthStore } from '@/stores/auth';
const auth = useAuthStore();

const errorMessage = ref<string | null>(null);

function submitLogin() {
  if (username && password.value) {
    fetch(`${import.meta.env.VITE_API_URI}/auth/sign-in`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({username: username.value, password: password.value})
    })
    .then(async response => {
      if (response.ok) {
        // Handle successful login
        const json = await response.json();
        auth.set(json.token);
        router.push('/dashboard');
      } else {
        // Handle login error
        errorMessage.value = `Error: ${response.statusText}`;
      }
    })
    .catch(error => {
      // Handle network error
      errorMessage.value = `Fatal Error: ${error.message}`;
    });
  } else {
    errorMessage.value = 'Username and password are required';
  }
}

</script>

<template>
  <form @submit.prevent="submitLogin">
    <input v-model="username" placeholder="Username">
    <input type="password" v-model="password" placeholder="Password">
    <button @click="submitLogin">Login</button>
    <button @click="router.push('/register')">Register</button>
  </form>

  <div v-if="errorMessage" class="error">{{ errorMessage }}</div>
  
</template>

<style scoped>
  form {
    display: flex;
    flex-direction: column;
    width: 300px;
    margin: auto;
    color: white;
  }
  input {
    margin: 10px;
    padding: 10px;
    border: none;
    border-radius: 20px;
  }
  button {
    padding: 10px;
    border: none;
    border-radius: 20px;
    background-color: transparent;
    color: white;
    cursor: pointer;
  }
  button:hover {
    background-color: gray;
  }
  input:focus {
    outline: none;
    border: 2px solid greenyellow;
  }
  input::placeholder {
    color: lightgray;
  }
  input:focus::placeholder {
    color: transparent;
  }

  .error {
    color: red;
    text-align: center;
    margin-top: 10px;
  }
</style>