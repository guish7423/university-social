import request from "./request"

export interface CampusEvent {
  id: number
  title: string
  date: string
  location: string
  description: string
}

export interface CampusContact {
  id: number
  name: string
  department: string
  phone: string
  email: string
  office: string
}

export interface EmptyRoom {
  id: number
  building: string
  room: string
  capacity: number
}

export interface Building {
  id: number
  name: string
}

export function listCalendar() {
  return request.get<any, CampusEvent[]>("/campus/calendar")
}

export function listDirectory() {
  return request.get<any, CampusContact[]>("/campus/directory")
}

export function listEmptyRooms(building?: string) {
  return request.get<any, EmptyRoom[]>("/campus/rooms", { params: { building } })
}

export function listBuildings() {
  return request.get<any, Building[]>("/campus/rooms/buildings")
}
