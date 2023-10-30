/* tslint:disable */
/* eslint-disable */
/**
 * Disclone Proxy API
 * Basic API proxy for Disclone
 *
 * The version of the OpenAPI document: 0.0.0
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */

import { exists, mapValues } from '../runtime';
import type { User } from './User';
import {
    UserFromJSON,
    UserFromJSONTyped,
    UserToJSON,
} from './User';

/**
 * 
 * @export
 * @interface ProfileGetResponseResult
 */
export interface ProfileGetResponseResult {
    /**
     * 
     * @type {User}
     * @memberof ProfileGetResponseResult
     */
    user?: User;
}

/**
 * Check if a given object implements the ProfileGetResponseResult interface.
 */
export function instanceOfProfileGetResponseResult(value: object): boolean {
    let isInstance = true;

    return isInstance;
}

export function ProfileGetResponseResultFromJSON(json: any): ProfileGetResponseResult {
    return ProfileGetResponseResultFromJSONTyped(json, false);
}

export function ProfileGetResponseResultFromJSONTyped(json: any, ignoreDiscriminator: boolean): ProfileGetResponseResult {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'user': !exists(json, 'user') ? undefined : UserFromJSON(json['user']),
    };
}

export function ProfileGetResponseResultToJSON(value?: ProfileGetResponseResult | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'user': UserToJSON(value.user),
    };
}

