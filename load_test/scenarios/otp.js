import http from 'k6/http';
import { checkResponse, logResponse } from '../utils/helpers.js';
import { config } from '../utils/config.js';
import { generateHeaders } from '../utils/data-generator.js';

/**
 * Generates an OTP for a given voter.
 * @param {string} voterId - The ID of the voter.
 * @param {string} phoneNumber - The phone number of the voter.
 * @returns {object} The response from the server.
 */
export function generateOTP(voterId, phoneNumber) {
    const headers = generateHeaders();
    const payload = JSON.stringify({
        voter_id: voterId,
        purpose: 'vote_cast',
        phone_number: phoneNumber,
    });

    const response = http.post(`${config.base_url}/v1/otp/generate`, payload, { headers });
    // The checkResponse function will log errors if the status is not 200
    checkResponse(response, 200, 'generate otp');
    return response;
}

/**
 * Retrieves the OTP for a given voter by calling a special test endpoint.
 * NOTE: This requires you to implement a test-only endpoint on your backend
 * that can retrieve the latest OTP for a voter from your Redis cache.
 * @param {string} voterId - The ID of the voter.
 * @returns {string|null} The OTP code.
 */
export function retrieveOTP(voterId) {
    // This function assumes you have created a test endpoint at `/v1/test/otp`
    const response = http.get(`${config.base_url}/v1/test/otp?voter_id=${voterId}&purpose=vote_cast`);

    if (!checkResponse(response, 200, 'retrieve test otp')) {
        // Error is logged by checkResponse
        return null;
    }

    try {
        const responseBody = JSON.parse(response.body);
        // This assumes your test endpoint returns JSON like: { "data": { "code": "123456" } }
        // Adjust the parsing based on your actual response structure.
        return responseBody.data.code;
    } catch (e) {
        console.log(`Failed to parse retrieve OTP response: ${e}`);
        logResponse(response, 'Parse OTP retrieval response failed')
        return null;
    }
}

/**
 * Verifies the OTP and returns an OTP token.
 * @param {string} voterId - The ID of the voter.
 * @param {string} otpCode - The OTP code to verify.
 * @returns {string|null} The OTP token if verification is successful, otherwise null.
 */
export function verifyOTP(voterId, otpCode) {
    const headers = generateHeaders();
    const payload = JSON.stringify({
        voter_id: voterId,
        purpose: 'vote_cast',
        code: otpCode,
    });

    const response = http.post(`${config.base_url}/v1/otp/verify`, payload, { headers });

    const success = checkResponse(response, 200, 'verify otp');
    if (!success) {
        return null;
    }

    try {
        const responseBody = JSON.parse(response.body);
        return responseBody.data.otp_token;
    } catch (e) {
        console.log(`Failed to parse verify otp response: ${e}`);
        logResponse(response, 'Parse OTP verification response failed')
        return null;
    }
}