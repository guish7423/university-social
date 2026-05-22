import { request } from "./request"

export interface CalendarEvent {
  id: number
  title: string
  event_date: string
  event_type: string
}

export interface DirectoryEntry {
  id: number
  name: string
  department: string
  position: string
  phone: string
  email: string
  office: string
  category: string
}

export interface EmptyRoom {
  id: number
  building: string
  room: string
  capacity: number
  weekday: number
  start_time: string
  end_time: string
}

export function listCalendar() {
  return request<CalendarEvent[]>({ url: "/campus/calendar", method: "GET" })
}

export function listDirectory(department?: string) {
  return request<DirectoryEntry[]>({ url: "/campus/directory", method: "GET", data: department ? { department } : {} })
}

export function listBuildings() {
  return request<string[]>({ url: "/campus/rooms/buildings", method: "GET" })
}

export function listEmptyRooms(building?: string, weekday?: number) {
  return request<EmptyRoom[]>({ url: "/campus/rooms", method: "GET", data: { building, weekday } })
}
