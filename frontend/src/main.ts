// Vue
import { createApp } from 'vue'

// Vuetify
import 'vuetify/styles'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'

import '@mdi/font/css/materialdesignicons.css'
import '@/styles/base.css'

import App from '@/App.vue'
import { router } from '@/router'

createApp(App)
  .use(router)
  .use(
    createVuetify({
      theme: { defaultTheme: 'dark' },
      components,
      directives,
    }),
  )
  .mount('#app')
