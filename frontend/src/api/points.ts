import { request } from "./request"

export function getPointsBalance() {
  return request.get('/points/balance')
}

export function claimDailyLogin() {
  return request.post('/points/daily-login')
}

export function getPointsLogs(offset = 0, limit = 20) {
  return request.get(`/points/logs?offset=${offset}&limit=${limit}`)
}

export function getPointsLeaderboard(limit = 20) {
  return request.get(`/points/leaderboard?limit=${limit}`)
}
