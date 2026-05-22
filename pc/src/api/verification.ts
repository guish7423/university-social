import request from "./request"

export function sendCode(data: { phone?: string; email?: string }) {
  return request.post<any, { message: string }>("/verification/send-code", data)
}

export function verify(data: { code: string; phone?: string; email?: string }) {
  return request.post<any, { message: string }>("/verification/verify", data)
}

export function verificationStatus() {
  return request.get<any, { is_verified: boolean }>("/verification/status")
}
