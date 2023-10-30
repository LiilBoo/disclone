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
/**
 * 
 * @export
 * @interface UpdateProfileRequest
 */
export interface UpdateProfileRequest {
    /**
     * 
     * @type {string}
     * @memberof UpdateProfileRequest
     */
    displayName: string;
}

/**
 * Check if a given object implements the UpdateProfileRequest interface.
 */
export function instanceOfUpdateProfileRequest(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "displayName" in value;

    return isInstance;
}

export function UpdateProfileRequestFromJSON(json: any): UpdateProfileRequest {
    return UpdateProfileRequestFromJSONTyped(json, false);
}

export function UpdateProfileRequestFromJSONTyped(json: any, ignoreDiscriminator: boolean): UpdateProfileRequest {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'displayName': json['display_name'],
    };
}

export function UpdateProfileRequestToJSON(value?: UpdateProfileRequest | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'display_name': value.displayName,
    };
}

