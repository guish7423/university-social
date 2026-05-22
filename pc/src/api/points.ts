import request from "./request"

export interface PointsLog {
  id: number
  amount: number
  reason: string
  created_at: string
}

export function getBalance() {
  return request.get<any, { balance: number }>("/points/balance")
}

export function claimDailyLogin() {
  return request.post<any, { message: string; amount: number }>("/points/daily-login")
}

export function getLogs() {
  return request.get<any, PointsLog[]>("/points/logs")
}

export function getLeaderboard() {
  return request.get<any, { id: number; nickname: string; avatar: string; points: number }[]>("/points/leaderboard")
}
