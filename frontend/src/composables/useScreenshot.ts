import { ref, type Ref } from 'vue'

import html2canvas from 'html2canvas'

export function useScreenshot(elementRef: Ref<HTMLElement>) {
  const screenshotUrl = ref<string | null>(null)
  const loading = ref<boolean>(false)
  const error = ref<string | null>(null)

  const options = {
    backgroundColor: null,
    scale: 1,
    useCORS: true,
    allowTaint: true,
  }

  async function capture() {
    if (!elementRef.value) {
      error.value = 'element not found'
    }

    loading.value = true
    error.value = null

    try {
      const canvas = await html2canvas(elementRef.value, options)
      const dataUrl = canvas.toDataURL('image/png')

      screenshotUrl.value = dataUrl
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'unknown error'
    } finally {
      loading.value = false
    }
  }

  function download(filename: string = 'screenshot') {
    if (!screenshotUrl.value) {
      error.value = 'no screenshot found'
      return
    }

    const link = document.createElement('a')
    link.href = screenshotUrl.value
    link.download = `${filename}.png`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
  }

  const clear = () => {
    if (screenshotUrl.value) {
      URL.revokeObjectURL(screenshotUrl.value)
      screenshotUrl.value = null
    }
    error.value = null
  }

  return {
    screenshotUrl,
    loading,
    error,

    capture,
    download,
    clear,
  }
}
