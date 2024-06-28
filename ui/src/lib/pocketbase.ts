import PocketBase from 'pocketbase'

const pocketbase = new PocketBase('/')
pocketbase.autoCancellation(false)
export const pb = pocketbase
