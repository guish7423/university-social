import request from "./request"

export interface InviteCodeData {
  id: number
  code: string
  is_used: boolean
  used_by?: number
  created_at: string
}

export function createInviteCode() {
  return request.post<any, { code: string }>("/invite-codes")
}

export function myInviteCodes() {
  return request.get<any, InviteCodeData[]>("/invite-codes")
}

export function redeemInviteCode(code: string) {
  return request.post<any, { message: string }>("/invite-codes/redeem", { code })
}
