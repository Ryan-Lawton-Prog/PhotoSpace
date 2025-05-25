import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

export const useAuthStore = defineStore('authToken', () => {
  const token = ref<string | null>(null)
    const set = (newToken: string | null) => {
        token.value = newToken
    }
  const isAuthenticated = computed(() => token.value !== null && token.value !== '')

  return { token, isAuthenticated, set }
})
