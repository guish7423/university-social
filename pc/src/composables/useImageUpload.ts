import { ref } from "vue"
import request from "@/api/request"

export function useImageUpload() {
  const uploading = ref(false)

  async function uploadImage(file: File): Promise<string> {
    const formData = new FormData()
    formData.append("file", file)
    const res = await request.post("/upload", formData, {
      headers: { "Content-Type": "multipart/form-data" },
    }) as { url: string }
    return res.url
  }

  async function uploadImages(files: File[]): Promise<string[]> {
    uploading.value = true
    try {
      return await Promise.all(files.map(uploadImage))
    } finally {
      uploading.value = false
    }
  }

  return { uploading, uploadImage, uploadImages }
}
