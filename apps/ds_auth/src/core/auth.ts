import {User} from "../types";

export function createSessionID(user: User) {
    return "123456"
}

export function createSignatureToken(user: User) {
    return "123456"
}

export function passwordMatches(password: string, hash: string) {
    return password === hash
}

