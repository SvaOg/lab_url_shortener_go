<script setup>
import { ref } from 'vue'

const longUrl = ref('')
const shortUrl = ref(null)
const loading = ref(false)
const error = ref(null)

// Мы обращаемся к локалхосту, так как запускаем фронт без Докера
const API_URL = 'http://localhost:8000'

const shortenUrl = async () => {
  // Сброс состояний
  error.value = null
  shortUrl.value = null
  
  if (!longUrl.value) {
    error.value = "Please enter a URL"
    return
  }

  loading.value = true

  try {
    const response = await fetch(`${API_URL}/`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ url: longUrl.value })
    })

    if (!response.ok) {
      throw new Error('Network response was not ok')
    }

    const data = await response.json()
    // Формируем полную ссылку для пользователя
    shortUrl.value = `${API_URL}/${data.slug}`
    
  } catch (err) {
    console.error(err)
    error.value = 'Failed to shorten link. Is Backend running?'
  } finally {
    loading.value = false
  }
}

const copyToClipboard = () => {
  navigator.clipboard.writeText(shortUrl.value)
  alert('Copied!')
}
</script>

<template>
  <div class="min-h-screen bg-slate-900 flex items-center justify-center p-4">
    
    <div class="bg-white rounded-xl shadow-2xl p-8 max-w-md w-full">
      <h1 class="text-3xl font-bold text-slate-800 mb-2 text-center">URL Shortener</h1>
      <p class="text-slate-500 text-center mb-6">Modern & Fast</p>

      <div class="space-y-4">
        <div>
          <input 
            v-model="longUrl"
            type="url" 
            placeholder="Paste long URL here..." 
            class="w-full px-4 py-3 border border-slate-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:outline-none transition text-slate-700 placeholder-slate-400"
            @keyup.enter="shortenUrl"
          />
        </div>

        <button 
          @click="shortenUrl" 
          :disabled="loading"
          class="w-full bg-indigo-600 hover:bg-indigo-700 text-white font-bold py-3 rounded-lg transition cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <span v-if="loading">Processing...</span>
          <span v-else>Shorten URL</span>
        </button>
      </div>

      <div v-if="error" class="mt-4 p-3 bg-red-50 text-red-600 border border-red-200 rounded-lg text-sm text-center">
        {{ error }}
      </div>

      <div v-if="shortUrl" class="mt-6 p-4 bg-indigo-50 rounded-lg border border-indigo-100">
        <p class="text-xs text-indigo-500 font-bold uppercase mb-1">Success!</p>
        <div class="flex items-center justify-between gap-2">
          <a :href="shortUrl" target="_blank" class="text-indigo-700 font-medium truncate hover:underline">
            {{ shortUrl }}
          </a>
          <button 
            @click="copyToClipboard"
            class="text-indigo-600 hover:text-indigo-800 p-2 rounded hover:bg-indigo-200 transition cursor-pointer"
            title="Copy"
          >
            Copy
          </button>
        </div>
      </div>

    </div>
  </div>
</template>
