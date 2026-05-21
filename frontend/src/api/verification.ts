import { request } from "./request"

export interface VerificationStatus {
  verified: boolean
}

export async function sendVerificationCode(email: string, school: string): Promise<void> {
  await request("/verification/send-code", {
    method: "POST",
    data: { email, school },
  })
}

export async function verifyCode(email: string, code: string): Promise<void> {
  await request("/verification/verify", {
    method: "POST",
    data: { email, code },
  })
}

export async function getVerificationStatus(): Promise<VerificationStatus> {
  return request("/verification/status")
}
